package ftrpc

import (
	"encoding/binary"
	"errors"
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/core/trpc/trpc"
	"google.golang.org/protobuf/proto"
)

type DataFrameType uint8

type StreamFrameType uint8

const (
	// MagicValue 笑死了
	MagicValue          uint16          = 0x930
	FrameHeaderLength   int             = 16
	DataFrameUnary      DataFrameType   = 0x0
	DataFrameStream     DataFrameType   = 0x1
	StreamFrameUnary    StreamFrameType = 0x0
	StreamFrameInit     StreamFrameType = 0x1
	StreamFrameData     StreamFrameType = 0x2
	StreamFrameFeedback StreamFrameType = 0x3
	StreamFrameClose    StreamFrameType = 0x4
)

type FrameHeader struct {
	DataFrameType   DataFrameType
	StreamFrameType StreamFrameType
	TotalLen        uint32
	HeadLen         uint16
	StreamID        uint32
}

type TrpcProtocol struct {
	FrameHeader *FrameHeader
	HeaderBuf   []byte
	BodyBuf     []byte
}

var (
	ParsePayloadError = errors.New("failed to parse tRPC")
)

func readBytes(p []byte, n int) (remain []byte, buf []byte, err error) {
	if len(p) < n {
		return nil, nil, ParsePayloadError
	}
	return p[n:], p[:n], nil
}

func readUint8(p []byte) (remain []byte, v uint8, err error) {
	var buf []byte
	if remain, buf, err = readBytes(p, 1); err != nil {
		return nil, 0, err
	}
	return remain, buf[0], nil
}

func readUint16(p []byte) (remain []byte, v uint16, err error) {
	var buf []byte
	if remain, buf, err = readBytes(p, 2); err != nil {
		return nil, 0, err
	}
	return remain, binary.BigEndian.Uint16(buf), nil
}

func readUint32(p []byte) (remain []byte, v uint32, err error) {
	var buf []byte
	if remain, buf, err = readBytes(p, 4); err != nil {
		return nil, 0, err
	}
	return remain, binary.BigEndian.Uint32(buf), nil
}

func Read(p []byte) (tp *TrpcProtocol, err error) {
	var magicValue uint16
	p, magicValue, err = readUint16(p)
	if magicValue != MagicValue {
		return nil, err
	}
	fh := &FrameHeader{}
	p, dataFrameType, err := readUint8(p)
	if err != nil {
		return nil, err
	}
	p, streamFrameType, err := readUint8(p)
	if err != nil {
		return nil, err
	}

	if p, fh.TotalLen, err = readUint32(p); err != nil {
		return nil, err
	}
	if p, fh.HeadLen, err = readUint16(p); err != nil {
		return nil, err
	}
	if p, fh.StreamID, err = readUint32(p); err != nil {
		return nil, err
	}

	fh.DataFrameType = DataFrameType(dataFrameType)
	fh.StreamFrameType = StreamFrameType(streamFrameType)

	// 读取并忽略保留位
	if p, _, err = readBytes(p, 2); err != nil {
		return nil, err
	}

	tp = &TrpcProtocol{FrameHeader: fh}
	if p, tp.HeaderBuf, err = readBytes(p, int(fh.HeadLen)); err != nil {
		return nil, err
	}
	if p, tp.BodyBuf, err = readBytes(p, int(fh.TotalLen)-FrameHeaderLength-int(fh.HeadLen)); err != nil {
		return nil, err
	}
	return tp, nil
}

func HandleUnaryReq(tp *TrpcProtocol) (p *trpc.RequestProtocol, err error) {
	p = &trpc.RequestProtocol{}
	if err = proto.Unmarshal(tp.HeaderBuf, p); err != nil {
		return nil, err
	}
	return p, nil
}

func HandleUnaryResp(tp *TrpcProtocol) (p *trpc.ResponseProtocol, err error) {
	p = &trpc.ResponseProtocol{}
	if err = proto.Unmarshal(tp.HeaderBuf, p); err != nil {
		return nil, err
	}
	return p, nil
}
