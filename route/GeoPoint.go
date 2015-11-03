// automatically generated, do not modify

package route

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type GeoPoint struct {
	_tab flatbuffers.Struct
}

func (rcv *GeoPoint) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *GeoPoint) Lat() int32 { return rcv._tab.GetInt32(rcv._tab.Pos + flatbuffers.UOffsetT(0)) }
func (rcv *GeoPoint) Lon() int32 { return rcv._tab.GetInt32(rcv._tab.Pos + flatbuffers.UOffsetT(4)) }

func CreateGeoPoint(builder *flatbuffers.Builder, lat int32, lon int32) flatbuffers.UOffsetT {
    builder.Prep(4, 8)
    builder.PrependInt32(lon)
    builder.PrependInt32(lat)
    return builder.Offset()
}
