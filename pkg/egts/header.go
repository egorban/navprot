package egts

import (
	"bytes"
	"encoding/binary"
	"errors"
)

func (packetData *Packet) parseHeader(message []byte) (body, restBuf []byte, err error) {
	index := bytes.IndexByte(message, prvSignature)
	if index == -1 {
		err = errors.New("prvSignature signature not found")
		return
	}
	messageLen := len(message) - index
	if messageLen < minEgtsHeaderLen {
		restBuf = append([]byte(nil), message...)
		err = errors.New("message is too short")
		return
	}
	headerLen := int(message[index+3])
	if messageLen < headerLen {
		restBuf = append([]byte(nil), message...)
		err = errors.New("message is too short")
		return
	}
	if headerLen < minEgtsHeaderLen {
		err = errors.New("headerLen is too short")
		return
	}
	header := message[index : index+headerLen]
	headerCrc := header[headerLen-1]
	headerCrcCalc := crc8EGTS(header[:headerLen-1])
	if uint(headerCrc) != headerCrcCalc {
		err = errors.New("incorrect header crc")
		return
	}
	startBody := index + headerLen
	bodyLen := int(binary.LittleEndian.Uint16(header[5:7]))
	if len(message[startBody:]) < bodyLen+2 {
		err = errors.New("message is too short")
		return
	}
	body = message[startBody : startBody+bodyLen]
	bodyCrc := binary.LittleEndian.Uint16(message[startBody+bodyLen : startBody+bodyLen+2])
	bodyCrcCalc := crc16EGTS(body)
	if bodyCrc != bodyCrcCalc {
		err = errors.New("incorrect parseHeader crc")
		return
	}
	packetData.Type = message[index+9]
	restBuf = append([]byte(nil), message[index+headerLen+bodyLen+2:]...)
	return
}
