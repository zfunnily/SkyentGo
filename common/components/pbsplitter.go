package components

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"pro2d/common"
)

type PBHead struct {
	Length   uint32
	Cmd      uint32
	ErrCode  int32
	PreField uint32
}

func (h *PBHead) GetDataLen() uint32 {
	return h.Length
}

func (h *PBHead) GetMsgID() uint32 {
	return h.Cmd
}

func (h *PBHead) GetErrCode() int32 {
	return h.ErrCode
}

func (h *PBHead) GetPreserve() uint32 {
	return h.PreField
}

type PBMessage struct {
	IMessage
	Head IHead
	Body []byte

	ID uint32
}

func (m *PBMessage) GetHeader() IHead {
	return m.Head
}

func (m *PBMessage) SetHeader(header IHead) {
	m.Head = header
}
func (m *PBMessage) GetData() []byte {
	return m.Body
}

func (m *PBMessage) SetData(b []byte) {
	m.Body = b
}

func (m *PBMessage) SetSID(id uint32) {
	m.ID = id
}

func (m *PBMessage) GetSID() uint32 {
	return m.ID
}

type PBSplitter struct {
	encipher IEncipher
}

func NewPBSplitter(encipher IEncipher) ISplitter {
	return &PBSplitter{
		encipher,
	}
}

func (m *PBSplitter) GetHeadLen() uint32 {
	return uint32(binary.Size(PBHead{}))
}

func (m *PBSplitter) ParseMsg(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// 表示我们已经扫描到结尾了
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if !atEOF && len(data) >= int(m.GetHeadLen()) { //4字节数据包长度  4字节指令
		length := int32(0)
		binary.Read(bytes.NewReader(data[0:4]), binary.BigEndian, &length)
		if length <= 0 {
			return 0, nil, fmt.Errorf("length is 0")
		}

		if length > common.MaxPacketLength {
			return 0, nil, fmt.Errorf("length exceeds maximum length")
		}
		if int(length) <= len(data) {
			return int(length), data[:int(length)], nil
		}
		return 0, nil, nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func (m *PBSplitter) Pack(cmd uint32, data []byte, errcode int32, preserve uint32) ([]byte, error) {
	buf := &bytes.Buffer{}
	h := &PBHead{
		Length:   m.GetHeadLen(),
		Cmd:      cmd,
		ErrCode:  errcode,
		PreField: preserve,
	}
	var dataEn []byte
	var err error
	if m.encipher != nil {
		dataEn, err = m.encipher.Encrypt(data)
		if err != nil {
			return nil, err
		}
	} else {
		dataEn = data
	}

	h.Length += uint32(len(dataEn))

	err = binary.Write(buf, binary.BigEndian, h)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, dataEn)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *PBSplitter) UnPack(data []byte) (IMessage, error) {
	h := &PBHead{}
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, h)
	if err != nil {
		return nil, err
	}

	var dataDe []byte
	if m.encipher != nil {
		dataDe, err = m.encipher.Decrypt(data[m.GetHeadLen():])
		if err != nil {
			return nil, err
		}
	} else {
		dataDe = data[m.GetHeadLen():]
	}

	return &PBMessage{
		Head: h,
		Body: dataDe,
	}, nil
}
