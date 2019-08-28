package ndtp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// NplData describes transport layer of NDPT protocol
type NplData struct {
	DataType    byte
	PeerAddress []byte
	ReqID       uint16
}

func (npl *NplData) String() string {
	if npl == nil {
		return "NPL: nil; "
	}
	return fmt.Sprintf("NPL: %+v;", *npl)
}

func parseNPL(message []byte) (npl *NplData, packet, restBuf []byte, err error) {
	index := bytes.Index(message, nplSignature)
	if index == -1 {
		err = errors.New("nplData signature not found")
		return
	}
	messageLen := len(message) - index
	if messageLen < nplHeaderLen {
		restBuf = append([]byte(nil), message...)
		err = errors.New("messageLen is too short")
		return
	}
	dataLen := int(binary.LittleEndian.Uint16(message[index+2 : index+4]))
	if messageLen < nplHeaderLen+dataLen {
		restBuf = append([]byte(nil), message...)
		err = errors.New("messageLen is too short")
		return
	}
	packetLen := index + nplHeaderLen + dataLen
	if binary.LittleEndian.Uint16(message[index+4:index+6])&2 != 0 {
		crcHead := binary.BigEndian.Uint16(message[index+6 : index+8])
		crcCalc := crc16(message[index+nplHeaderLen : packetLen])
		if crcHead != crcCalc {
			err = fmt.Errorf("crc incorrect: calc %d; receive: %d", crcCalc, crcHead)
			return
		}
	}
	npl = new(NplData)
	npl.DataType = message[index+8]
	npl.ReqID = binary.LittleEndian.Uint16(message[index+13 : index+15])
	npl.PeerAddress = message[index+9 : index+13]
	packet = message[index:packetLen]
	restBuf = append([]byte(nil), message[packetLen:]...)
	return
}
