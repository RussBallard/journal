// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package models

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type FdbxJournalT struct {
	Start   int64
	Total   int64
	Chain   []*FdbxStageT
	Service string
}

func (t *FdbxJournalT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	chainOffset := flatbuffers.UOffsetT(0)
	if t.Chain != nil {
		chainLength := len(t.Chain)
		chainOffsets := make([]flatbuffers.UOffsetT, chainLength)
		for j := 0; j < chainLength; j++ {
			chainOffsets[j] = t.Chain[j].Pack(builder)
		}
		FdbxJournalStartChainVector(builder, chainLength)
		for j := chainLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(chainOffsets[j])
		}
		chainOffset = builder.EndVector(chainLength)
	}
	serviceOffset := builder.CreateString(t.Service)
	FdbxJournalStart(builder)
	FdbxJournalAddStart(builder, t.Start)
	FdbxJournalAddTotal(builder, t.Total)
	FdbxJournalAddChain(builder, chainOffset)
	FdbxJournalAddService(builder, serviceOffset)
	return FdbxJournalEnd(builder)
}

func (rcv *FdbxJournal) UnPackTo(t *FdbxJournalT) {
	t.Start = rcv.Start()
	t.Total = rcv.Total()
	chainLength := rcv.ChainLength()
	t.Chain = make([]*FdbxStageT, chainLength)
	for j := 0; j < chainLength; j++ {
		x := FdbxStage{}
		rcv.Chain(&x, j)
		t.Chain[j] = x.UnPack()
	}
	t.Service = string(rcv.Service())
}

func (rcv *FdbxJournal) UnPack() *FdbxJournalT {
	if rcv == nil {
		return nil
	}
	t := &FdbxJournalT{}
	rcv.UnPackTo(t)
	return t
}

type FdbxJournal struct {
	_tab flatbuffers.Table
}

func GetRootAsFdbxJournal(buf []byte, offset flatbuffers.UOffsetT) *FdbxJournal {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &FdbxJournal{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *FdbxJournal) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *FdbxJournal) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *FdbxJournal) Start() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FdbxJournal) MutateStart(n int64) bool {
	return rcv._tab.MutateInt64Slot(4, n)
}

func (rcv *FdbxJournal) Total() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FdbxJournal) MutateTotal(n int64) bool {
	return rcv._tab.MutateInt64Slot(6, n)
}

func (rcv *FdbxJournal) Chain(obj *FdbxStage, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *FdbxJournal) ChainLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *FdbxJournal) Service() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func FdbxJournalStart(builder *flatbuffers.Builder) {
	builder.StartObject(4)
}
func FdbxJournalAddStart(builder *flatbuffers.Builder, start int64) {
	builder.PrependInt64Slot(0, start, 0)
}
func FdbxJournalAddTotal(builder *flatbuffers.Builder, total int64) {
	builder.PrependInt64Slot(1, total, 0)
}
func FdbxJournalAddChain(builder *flatbuffers.Builder, chain flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(chain), 0)
}
func FdbxJournalStartChainVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func FdbxJournalAddService(builder *flatbuffers.Builder, service flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(service), 0)
}
func FdbxJournalEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
