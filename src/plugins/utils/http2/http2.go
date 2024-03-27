package http2

import (
	"encoding/binary"
	"errors"
)

type FrameType uint8

type Flags uint8

func (f Flags) Has(v Flags) bool {
	return (f & v) == v
}

const (
	GRPC                uint8     = 41
	FrameHeaderLen      int       = 9
	FrameTypeMax        FrameType = 0x9
	FrameData           FrameType = 0x0
	FrameHeaders        FrameType = 0x1
	FrameSetting        FrameType = 0x4
	FramePushPromise    FrameType = 0x5
	FramePing           FrameType = 0x6
	FrameWindowUpdate   FrameType = 0x8
	FrameContinuation   FrameType = 0x9
	FlagDataPadded      Flags     = 0x8
	FlagDataEndStream   Flags     = 0x1
	FlagHeadersPadded   Flags     = 0x8
	FlagHeadersPriority Flags     = 0x20
)

var (
	ParsePayloadError = errors.New("failed to l7 payload")
)

type FrameHeader struct {
	Length uint32
	Type   FrameType
	Flags  Flags
}

func readByte(p []byte) (remain []byte, b byte, err error) {
	if len(p) == 0 {
		return nil, 0, ParsePayloadError
	}
	return p[1:], p[0], nil
}

func readUint32(p []byte) (remain []byte, v uint32, err error) {
	if len(p) < 4 {
		return nil, 0, ParsePayloadError
	}
	return p[4:], binary.BigEndian.Uint32(p[:4]), nil
}

// frameOverflow 帧溢出检测
func frameOverflow(p []byte, fh FrameHeader) bool {
	if int(fh.Length) > len(p) {
		return true
	}
	return false
}

// readFrameHeader 读取帧头部信息
func readFrameHeader(p []byte) ([]byte, FrameHeader, error) {
	if len(p) < FrameHeaderLen {
		return nil, FrameHeader{}, ParsePayloadError
	}

	fh := FrameHeader{
		Type: FrameType(p[3]),
	}
	if fh.Type > FrameTypeMax {
		return nil, fh, ParsePayloadError
	}

	fh.Length = uint32(p[0])<<16 | uint32(p[1])<<8 | uint32(p[2])
	fh.Flags = Flags(p[4])

	return p[FrameHeaderLen:], fh, nil
}

// readHeaderBlockFragment 读取 HEADERS 帧的 Header Block Fragment
func readHeaderBlockFragment(p []byte, fh FrameHeader) (remain []byte, bf []byte, err error) {
	if isOverflow := frameOverflow(p, fh); isOverflow {
		return nil, nil, ParsePayloadError
	}
	var padLength uint8
	if fh.Flags.Has(FlagHeadersPadded) {
		if p, padLength, err = readByte(p); err != nil {
			return nil, nil, err
		}
	}
	if fh.Flags.Has(FlagHeadersPriority) {
		// read E + Stream Dependency
		if p, _, err = readUint32(p); err != nil {
			return nil, nil, err
		}
		// read Weight
		if p, _, err = readByte(p); err != nil {
			return nil, nil, err
		}
	}

	if len(p)-int(padLength) < 0 || int(padLength) > int(fh.Length) {
		return nil, nil, ParsePayloadError
	}
	return p[fh.Length:], p[:int(fh.Length)-int(padLength)], nil
}

// readDataBlockFragment 读取 DATA 帧 Payload
func readDataBlockFragment(p []byte, fh FrameHeader) (remain []byte, bf []byte, err error) {
	if isOverflow := frameOverflow(p, fh); isOverflow {
		return nil, nil, ParsePayloadError
	}
	var padLength uint8
	if fh.Flags.Has(FlagDataPadded) {
		if p, padLength, err = readByte(p); err != nil {
			return nil, nil, err
		}
	} else {
		padLength = 0
	}
	if len(p)-int(padLength) < 0 || int(padLength) > int(fh.Length) {
		return nil, nil, ParsePayloadError
	}
	return p[fh.Length:], p[:int(fh.Length)-int(padLength)], nil
}

// readFramePayload 读取一个帧
func readFramePayload(p []byte, fh FrameHeader) (remain []byte, fp []byte, err error) {
	if isOverflow := frameOverflow(p, fh); isOverflow {
		return nil, nil, ParsePayloadError
	}
	return p[fh.Length:], p[:fh.Length], nil
}

type FrameHook func(fh FrameHeader, bf []byte) error

type Parser struct {
	hooks map[FrameType]FrameHook
}

func (p *Parser) RegisterHook(frameType FrameType, hook FrameHook) {
	p.hooks[frameType] = hook
}

func (p *Parser) Do(payload []byte) {
	var err error
	var bf []byte
	var fh FrameHeader
	for len(payload) >= FrameHeaderLen {
		if payload, fh, err = readFrameHeader(payload); err != nil {
			break
		}
		switch fh.Type {
		case FrameData:
			payload, bf, err = readDataBlockFragment(payload, fh)
		case FrameHeaders:
			payload, bf, err = readHeaderBlockFragment(payload, fh)
		default:
			payload, bf, err = readFramePayload(payload, fh)
		}
		if err != nil {
			break
		}
		if hook, ok := p.hooks[fh.Type]; ok {
			if err = hook(fh, bf); err != nil {
				break
			}
		}
	}
}

func NewParser() *Parser {
	return &Parser{
		hooks: make(map[FrameType]FrameHook),
	}
}
