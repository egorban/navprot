package ndtp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// NplData describes transport layer of NDTP protocol
type NplData struct {
	PeerAddress []byte
	DataType    byte
	ReqID       uint16
}

func (npl *NplData) String() string {
	if npl == nil {
		return "NPL: nil; "
	}
	return fmt.Sprintf("NPL: %+v;", *npl)
}

func parseNPL(message []byte) (npl *NplData, packet, restBuf []byte, err error) {
	first, last, restBuf, err := checkPacket(message)
	if err != nil {
		return
	}
	npl = new(NplData)
	npl.DataType = message[first+8]
	npl.ReqID = binary.LittleEndian.Uint16(message[first+13 : first+15])
	npl.PeerAddress = message[first+9 : first+13]
	packet = message[first:last]
	return
}

func simpleParseNPL(message []byte) (packet []byte, restBuf []byte, err error) {
	first, last, restBuf, err := checkPacket(message)
	if err != nil {
		restBuf = message
		return
	}
	return message[first:last], message[last:], nil
}

func checkPacket(message []byte) (first, last int, rest []byte, err error){
	first = bytes.Index(message, nplSignature)
	if first == -1 {
		err = errors.New("nplData signature not found")
		return
	}
	messageLen := len(message) - first
	if messageLen < nplHeaderLen {
		err = errors.New("messageLen is too short")
		rest = message[first:]
		return
	}
	dataLen := int(binary.LittleEndian.Uint16(message[first+2 : first+4]))
	if messageLen < nplHeaderLen+dataLen {
		err = errors.New("messageLen is too short")
		rest = message[first:]
		return
	}
	last = first + nplHeaderLen + dataLen
	if binary.LittleEndian.Uint16(message[first+4:first+6])&2 != 0 {
		crcHead := binary.BigEndian.Uint16(message[first+6 : first+8])
		crcCalc := crc16(message[first+nplHeaderLen : last])
		if crcHead != crcCalc {
			err = fmt.Errorf("crc incorrect: calc %d; receive: %d", crcCalc, crcHead)
			return
		}
	}
	rest = message[last:]
	return
}
