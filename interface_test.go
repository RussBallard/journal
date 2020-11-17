package journal_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/shestakovda/errx"
	"github.com/shestakovda/fdbx"
	"github.com/shestakovda/journal"
	"github.com/shestakovda/journal/crash"
	"github.com/shestakovda/typex"
	"github.com/stretchr/testify/suite"
)

func TestFactory(t *testing.T) {
	suite.Run(t, new(journal.FdbSuite))
}

func TestProvider(t *testing.T) {
	suite.Run(t, new(journal.ProviderSuite))
}

func TestInterface(t *testing.T) {
	suite.Run(t, new(InterfaceSuite))
}

type InterfaceSuite struct {
	suite.Suite

	fdb fdbx.Conn
	crp crash.Provider
}

func (s *InterfaceSuite) SetupTest() {
	var err error

	if testing.Short() {
		s.T().SkipNow()
	}

	s.fdb, err = fdbx.NewConn(crash.DatabaseAPI, fdbx.ConnVersion610)
	s.Require().NoError(err)
	s.Require().NoError(s.fdb.ClearDB())

	s.crp = crash.NewTestProvider()
	s.crp.Register(http.StatusForbidden, journal.TestNum, journal.TestTitle, errx.ErrForbidden)
}

func (s *InterfaceSuite) TestWorkflow() {
	log := journal.NewProvider(1, s.crp, journal.NewDriverFDB(s.fdb), nil, "")
	log2 := log.Clone()
	log3 := log.Clone()

	// Не должны сюда попасть
	if log.V(3) {
		s.True(false)
	}

	// Должны попасть и записать строку
	if log.V(0) {
		log.Print("ololo %s %d", "test1", 41)
		log2.Print("ololo %s %d", "test2", 42)
		log3.Print("ololo %s %d", "test3", 43)
	}

	mt := &mType{
		id:   24,
		name: "event",
	}

	// Должны записать данные парочке моделей
	if log.V(0) {
		log.Model(mt, "", "empty %s", "id")
		log.Model(mt, "eventID", "some %s", "comment1")
		log2.Model(mt, "eventID", "some %s", "comment2")
		log3.Model(mt, "eventID", "some %s", "comment3")
	}

	// Должны записать данные о модели с ошибкой
	rep := log.Crash(journal.ErrTest.WithReason(errx.ErrForbidden))

	// Должны записать в лог и вызвать сохранение, несколько раз желательно, чтобы курсор работал
	entry := log.Close()
	entry2 := log2.Close()
	entry3 := log3.Close()

	var cur journal.Cursor

	// Где-то в другом месте его можно получить по айдишке
	s.Require().NoError(s.fdb.Tx(func(db fdbx.DB) (exp error) {
		var mod journal.Model

		fac := journal.NewFactoryFDB(s.fdb, db)

		if mod, exp = fac.ByID(entry.ID); exp != nil {
			return
		}

		// По крайней мере, их внешние представления должны совпадать
		if row, err := mod.Export(true); s.NoError(err) {
			s.Equal(entry, row)

			// Сравним, как это выгружается в формате API
			api := mod.ExportAPI(log)
			s.Equal(row.ID, api.ID)
			s.Equal(row.Total, api.Total)
			s.Equal("ololo test1 41", api.Name)
		}

		mtc := &mType{
			id:   crash.ModelTypeCrash,
			name: "crash",
		}

		// Попробуем найти по модели ошибки
		if mods, exp := fac.ByModel(mtc, rep.ID); s.NoError(exp) && s.Len(mods, 1) {
			if row, err := mods[0].Export(true); s.NoError(err) {
				s.Equal(entry, row)
			}
		}

		// Попробуем найти по дате
		from := time.Now().Add(-time.Hour)
		to := time.Now().Add(time.Hour)
		if cur, exp = fac.ByDate(from, to, 10); s.NoError(exp) {
			if mods, exp := cur.NextPage(10); s.NoError(exp) && s.Len(mods, 3) {
				if row, err := mods[0].Export(true); s.NoError(err) {
					s.Equal(entry3, row)
				}
				if row, err := mods[1].Export(true); s.NoError(err) {
					s.Equal(entry2, row)
				}
				if row, err := mods[2].Export(true); s.NoError(err) {
					s.Equal(entry, row)
				}
			}
		}

		// Попробуем найти по модели
		if cur, exp = fac.ByModelDate(mt, "eventID", from, to, 10); s.NoError(exp) {
			if mods, exp := cur.NextPage(1); s.NoError(exp) && s.Len(mods, 1) {
				if row, err := mods[0].Export(true); s.NoError(err) {
					s.Equal(entry3, row)
				}
			}
		}

		return nil
	}))

	s.False(cur.Empty())

	// В след. раз загружаем этот курсор и смотрим, чот там есть
	s.Require().NoError(s.fdb.Tx(func(db fdbx.DB) (exp error) {
		fac := journal.NewFactoryFDB(s.fdb, db)

		if _, exp = fac.Cursor(typex.NewUUID().Hex()); s.Error(exp) {
			s.True(errors.Is(exp, errx.ErrNotFound))
		}

		if cur, exp = fac.Cursor(cur.ID()); s.NoError(exp) {
			if mods, exp := cur.NextPage(5); s.NoError(exp) && s.Len(mods, 2) {
				if row, err := mods[0].Export(true); s.NoError(err) {
					s.Equal(entry2, row)
				}

				if row, err := mods[1].Export(true); s.NoError(err) {
					s.Equal(entry, row)
				}
			}
		}

		return nil
	}))

	s.True(cur.Empty())
}

type mType struct {
	id   int
	name string
}

func (mt mType) ID() int        { return mt.id }
func (mt mType) String() string { return mt.name }
