package ndtp

import (
	"encoding/binary"
	"fmt"
)

// ExtDevice describes information of NPH_SRV_EXTERNAL_DEVICE service
type ExtDevice struct {
	MesID   uint16
	PackNum uint16
	Res     uint32
}

func (ext *ExtDevice) parse(packetType string, message []byte) (err error) {
	switch packetType {
	case NphSedDeviceTitleData:
		ext.MesID = binary.LittleEndian.Uint16(message[:2])
		ext.PackNum = binary.LittleEndian.Uint16(message[2:4])
	case NphSedDeviceResult:
		ext.PackNum = binary.LittleEndian.Uint16(message[:2])
		ext.Res = binary.LittleEndian.Uint32(message[2:6])
		ext.MesID = binary.LittleEndian.Uint16(message[6:8])
	default:
		err = fmt.Errorf("parseExtDevice unknown NPHType: %s", packetType)
	}
	return
}

func (ext *ExtDevice) reply(packet []byte, result uint32) []byte {
	reply := make([]byte, ndtpExtResultLen)
	copy(reply, packet[:nplHeaderLen+nphHeaderLen])
	for i := nplHeaderLen + 4; i < nplHeaderLen+6; i++ {
		reply[i] = 0
	}
	binary.LittleEndian.PutUint16(reply[nplHeaderLen+nphHeaderLen:], ext.PackNum)
	binary.LittleEndian.PutUint32(reply[nplHeaderLen+nphHeaderLen+2:], result)
	binary.LittleEndian.PutUint16(reply[nplHeaderLen+nphHeaderLen+6:], ext.MesID)
	binary.LittleEndian.PutUint16(reply[nplHeaderLen+2:], uint16(nphSedDeviceResult))
	binary.LittleEndian.PutUint16(reply[2:], uint16(ndtpExtResultLen-nplHeaderLen))
	crc := crc16(reply[nplHeaderLen:])
	binary.BigEndian.PutUint16(reply[6:], crc)
	return reply
}
