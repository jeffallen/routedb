// automatically generated, do not modify

package route

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Route struct {
	_tab flatbuffers.Table
}

func GetRootAsRoute(buf []byte, offset flatbuffers.UOffsetT) *Route {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Route{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Route) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Route) Country() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Route) City() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Route) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Route) Path(obj *GeoPoint, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 8
		if obj == nil {
			obj = new(GeoPoint)
		}
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Route) PathLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func RouteStart(builder *flatbuffers.Builder) { builder.StartObject(4) }
func RouteAddCountry(builder *flatbuffers.Builder, country flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(country), 0)
}
func RouteAddCity(builder *flatbuffers.Builder, city flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(city), 0)
}
func RouteAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(name), 0)
}
func RouteAddPath(builder *flatbuffers.Builder, path flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(path), 0)
}
func RouteStartPathVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(8, numElems, 4)
}
func RouteEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
