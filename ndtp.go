package navprot

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// NDTP describes NDTP (Navigation Data Transfer Protocol)
type NDTP struct {
	Npl    *NplData
	Nph    *NphData
	Packet []byte
}

// NplData describes transport layer of NDPT protocol
type NplData struct {
	DataType    byte
	PeerAddress []byte
	ReqID       uint16
}

// NphData describes session layer of NDTP protocol
type NphData struct {
	ServiceID   uint16
	PacketType  uint16
	RequestFlag bool
	ReqID       uint32
	Data        interface{}
}

// NavData describes information of NPH_SRV_NAVDATA service
type NavData struct {
	Time    uint32
	Lon     float64
	Lat     float64
	Bearing uint16
	Speed   uint16
	Sos     bool
	// 0 - W; 1 - E
	Lohs int8
	// 0 - S; 1 - N
	Lahs  int8
	Valid bool
}

// ExtDevice describes information of NPH_SRV_EXTERNAL_DEVICE service
type ExtDevice struct {
	MesID   uint16
	PackNum uint16
	Res     uint32
}

var nplSignature = []byte{0x7E, 0x7E}

const (
	nplHeaderLen     = 15
	nphResult        = 0
	nphHeaderLen     = 10
	ndtpResultLen    = nphHeaderLen + nplHeaderLen + 4
	ndtpExtResultLen = nphHeaderLen + nplHeaderLen + 8

	// NPH service types

	// NphSrvGenericControls defines NPH_SRV_GENERIC_CONTROLS service
	NphSrvGenericControls = 0
	// NphSrvNavdata defines NPH_SRV_NAVDATA service
	NphSrvNavdata = 1
	// NphSrvExternalDevice defines NPH_SRV_EXTERNAL_DEVICE service
	NphSrvExternalDevice = 5

	// NphSrvGenericControls packets

	nphSgcConnRequest = 100
	// NphSgsConnRequest defines NPH_SGC_CONN_REQUEST packet type
	NphSgsConnRequest = "NPH_SGC_CONN_REQUEST"

	// NphSrvNavdata packets

	nphSndHistory = 100
	// NphSndHistory defines NPH_SND_HISTORY packet type
	NphSndHistory  = "NPH_SND_HISTORY"
	nphSndRealtime = 101
	// NphSndRealtime defines NPH_SND_REALTIME packet type
	NphSndRealtime = "NPH_SND_REALTIME"

	// NphSrvExternalDevice packets

	nphSedDeviceTitleData = 100
	// NphSedDeviceTitleData defines NPH_SED_DEVICE_TITLE_DATA packet type
	NphSedDeviceTitleData = "NPH_SED_DEVICE_TITLE_DATA"
	nphSedDeviceResult    = 102
	// NphSedDeviceResult defines NPH_SED_DEVICE_RESULT packet type
	NphSedDeviceResult = "NPH_SED_DEVICE_RESULT"

	// NDTP packet fields names

	// NplReqID defines NPL request id field
	NplReqID = "NplReqID"
	// NphReqID defines NPH request id field
	NphReqID = "NphReqID"
	// PacketType defines NPH type field
	PacketType = "PacketType"
)

// Parse NDTP packet. Parsed information is stored in variable with NDTP type.
func (packetData *NDTP) Parse(message []byte) (restBuf []byte, err error) {
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
	packetData.Npl = new(NplData)
	packetData.Npl.DataType = message[index+8]
	packetData.Npl.ReqID = binary.LittleEndian.Uint16(message[index+13 : index+15])
	err = parseNPH(message[index+nplHeaderLen:], packetData)
	packetData.Packet = message[index:packetLen]
	packetData.Npl.PeerAddress = packetData.Packet[index+9 : index+13]
	restBuf = append([]byte(nil), message[packetLen:]...)
	return
}

// Form generate NDTP binary packet.
func (packetData *NDTP) Form() []byte {
	//todo implement method
	panic("implement me")
}

// Print generate string with information about NDTP packet in readable format.
func (packetData NDTP) String() string {
	sNPL := printNPL(packetData.Npl)
	sNPH := printNPH(packetData.Nph)
	return sNPL + sNPH
}

// IsResult returns true, if packetData is a NPH_RESULT.
func (packetData *NDTP) IsResult() bool {
	return packetData.Nph.PacketType == 0
}

// ChangeAddress change Peer Address field of NPL layer
func (packetData *NDTP) ChangeAddress(newAddress []byte) {
	for i, j := 9, 0; i < 13; i, j = i+1, j+1 {
		packetData.Packet[i] = newAddress[j]
	}
}

// GetID returns ID of terminal, which is included only in NPH_SGC_CONN_REQUEST packets
func (packetData *NDTP) GetID() (id int, err error) {
	if packetData.PacketType() == NphSgsConnRequest {
		id = int(packetData.Nph.Data.(uint32))
	} else {
		err = errors.New("incorrect packet type")
	}
	return
}

// PacketType returns name of NDTP packet type.
func (packetData *NDTP) PacketType() (ptype string) {
	switch packetData.Nph.ServiceID {
	case NphSrvGenericControls:
		if packetData.Nph.PacketType == nphSgcConnRequest {
			ptype = NphSgsConnRequest
		}
	case NphSrvNavdata:
		if packetData.Nph.PacketType == nphSndHistory {
			ptype = NphSndHistory
		} else if packetData.Nph.PacketType == nphSndRealtime {
			ptype = NphSndRealtime
		}
	case NphSrvExternalDevice:
		if packetData.Nph.PacketType == nphSedDeviceTitleData {
			ptype = NphSedDeviceTitleData
		} else if packetData.Nph.PacketType == nphSedDeviceResult {
			ptype = NphSedDeviceResult
		}
	}
	return
}

// Service returns value of packet's service type.
func (packetData *NDTP) Service() int {
	return int(packetData.Nph.ServiceID)
}

// NeedReply returns true if packet needs reply
func (packetData *NDTP) NeedReply() (flag bool) {
	return packetData.Nph.RequestFlag
}

// Reply creates NPH_RESULT packet for packetData.Packet
func (packetData *NDTP) Reply(result uint32) []byte {
	reply := make([]byte, ndtpResultLen)
	copy(reply, packetData.Packet[:nplHeaderLen+nphHeaderLen])
	for i := nplHeaderLen + 2; i < nplHeaderLen+6; i++ {
		reply[i] = 0
	}
	binary.LittleEndian.PutUint32(reply[nplHeaderLen+nphHeaderLen:], result)
	binary.LittleEndian.PutUint16(reply[2:], uint16(ndtpResultLen-nplHeaderLen))
	crc := crc16(reply[nplHeaderLen:])
	binary.BigEndian.PutUint16(reply[6:], crc)
	return reply
}

// ReplyExt creates NPH_SED_DEVICE_RESULT  packet
func (packetData *NDTP) ReplyExt(result uint32) ([]byte, error) {
	if packetData.Service() == NphSrvExternalDevice {
		reply := make([]byte, ndtpExtResultLen)
		copy(reply, packetData.Packet[:nplHeaderLen+nphHeaderLen])
		for i := nplHeaderLen + 4; i < nplHeaderLen+6; i++ {
			reply[i] = 0
		}
		ext := packetData.Nph.Data.(ExtDevice)
		binary.LittleEndian.PutUint16(reply[nplHeaderLen+nphHeaderLen:], ext.PackNum)
		binary.LittleEndian.PutUint32(reply[nplHeaderLen+nphHeaderLen+2:], result)
		binary.LittleEndian.PutUint16(reply[nplHeaderLen+nphHeaderLen+6:], ext.MesID)
		binary.LittleEndian.PutUint16(reply[nplHeaderLen+2:], uint16(nphSedDeviceResult))
		binary.LittleEndian.PutUint16(reply[2:], uint16(ndtpExtResultLen-nplHeaderLen))
		crc := crc16(reply[nplHeaderLen:])
		binary.BigEndian.PutUint16(reply[6:], crc)
		return reply, nil
	}
	return nil, errors.New("incorrect packet service")

}

// ChangePacket changes values of some fields in NDTP packet
func (packetData *NDTP) ChangePacket(changes map[string]int) {
	if nplReqID, ok := changes[NplReqID]; ok {
		binary.LittleEndian.PutUint16(packetData.Packet[13:], uint16(nplReqID))
	}
	if nphReqID, ok := changes[NphReqID]; ok {
		binary.LittleEndian.PutUint32(packetData.Packet[nplHeaderLen+6:], uint32(nphReqID))
	}
	if packetType, ok := changes[PacketType]; ok {
		binary.LittleEndian.PutUint16(packetData.Packet[nplHeaderLen+2:], uint16(packetType))
	}
	crc := crc16(packetData.Packet[nplHeaderLen:])
	binary.BigEndian.PutUint16(packetData.Packet[6:], crc)
}

func parseNPH(message []byte, packetData *NDTP) (err error) {
	packetData.Nph = new(NphData)
	packetData.Nph.ServiceID = binary.LittleEndian.Uint16(message[:2])
	packetData.Nph.PacketType = binary.LittleEndian.Uint16(message[2:4])
	if binary.LittleEndian.Uint16(message[4:6]) == 1 {
		packetData.Nph.RequestFlag = true
	}
	packetData.Nph.ReqID = binary.LittleEndian.Uint32(message[6:10])
	if packetData.IsResult() {
		packetData.Nph.Data = binary.LittleEndian.Uint32(message[nphHeaderLen : nphHeaderLen+4])
		return
	}
	switch packetData.Service() {
	case NphSrvGenericControls:
		parseGenControl(packetData, message[nphHeaderLen:])
	case NphSrvNavdata:
		err = parseNavData(packetData, message[nphHeaderLen:])
	case NphSrvExternalDevice:
		err = parseExtDevice(packetData, message[nphHeaderLen:])
	default:
		err = errors.New("unknown service")
	}
	return
}

func parseNavData(packetData *NDTP, message []byte) (err error) {
	if message[0] == 0 {
		if len(message) >= 28 {
			data := new(NavData)
			data.Time = binary.LittleEndian.Uint32(message[2:6])
			lon := binary.LittleEndian.Uint32(message[6:10])
			lat := binary.LittleEndian.Uint32(message[10:14])
			if message[14]&128 != 0 {
				data.Valid = true
			}
			if message[14]&64 != 0 {
				data.Lohs = 1
			}
			if message[14]&32 != 0 {
				data.Lahs = 1
			}
			data.Lon = float64((2*int(data.Lohs)-1)*int(lon)) / 10000000.0
			data.Lat = float64((2*int(data.Lahs)-1)*int(lat)) / 10000000.0
			if message[14]&4 != 0 {
				data.Sos = true
			}
			data.Speed = binary.LittleEndian.Uint16(message[16:18])
			data.Bearing = binary.LittleEndian.Uint16(message[20:22])
			packetData.Nph.Data = data
		} else {
			err = errors.New("NavData type 0 is too short")
		}
	}
	return
}

func parseGenControl(packetData *NDTP, message []byte) {
	if packetData.PacketType() == NphSgsConnRequest {
		packetData.Nph.Data = binary.LittleEndian.Uint32(message[6:10])
	}
	return
}

func parseExtDevice(packetData *NDTP, message []byte) (err error) {
	switch packetData.PacketType() {
	case NphSedDeviceTitleData:
		ext := new(ExtDevice)
		ext.MesID = binary.LittleEndian.Uint16(message[:2])
		ext.PackNum = binary.LittleEndian.Uint16(message[2:4])
		packetData.Nph.Data = ext
	case NphSedDeviceResult:
		ext := new(ExtDevice)
		ext.PackNum = binary.LittleEndian.Uint16(message[:2])
		ext.Res = binary.LittleEndian.Uint32(message[2:6])
		ext.MesID = binary.LittleEndian.Uint16(message[6:8])
		packetData.Nph.Data = ext
	default:
		err = fmt.Errorf("parseExtDevice unknown NPHType: %d", packetData.Nph.PacketType)
	}
	return
}

func printNPL(npl *NplData) string {
	//return fmt.Sprintf("dataType: %d; peerAddress: %d; reqID: %d\n", npl.DataType, npl.PeerAddress, npl.ReqID)
	return fmt.Sprintf("NPL: %+v;", *npl)
}

func printNPH(nph *NphData) string {
	sNPH := fmt.Sprintf(" NPH: %+v;", *nph)
	sData := printData(nph.Data)
	return sNPH + sData
}

func printData(data interface{}) (sdata string) {
	prefix := " Data: "
	switch data.(type) {
	case nil:
		sdata = "nil"
	case int:
		sdata = fmt.Sprintf("%d", data)
	case *ExtDevice:
		ext := data.(*ExtDevice)
		sdata = fmt.Sprintf("%+v", *ext)
	case *NavData:
		nav := data.(*NavData)
		sdata = fmt.Sprintf("%+v", *nav)
	}
	return prefix + sdata
}

func crc16(bs []byte) (crc uint16) {
	l := len(bs)
	crc = 0xFFFF
	for i := 0; i < l; i++ {
		crc = (crc >> 8) ^ crc16tab[(crc&0xff)^uint16(bs[i])]
	}
	return
}

var crc16tab = [256]uint16{
	0x0000, 0xC0C1, 0xC181, 0x0140, 0xC301, 0x03C0, 0x0280, 0xC241,
	0xC601, 0x06C0, 0x0780, 0xC741, 0x0500, 0xC5C1, 0xC481, 0x0440,
	0xCC01, 0x0CC0, 0x0D80, 0xCD41, 0x0F00, 0xCFC1, 0xCE81, 0x0E40,
	0x0A00, 0xCAC1, 0xCB81, 0x0B40, 0xC901, 0x09C0, 0x0880, 0xC841,
	0xD801, 0x18C0, 0x1980, 0xD941, 0x1B00, 0xDBC1, 0xDA81, 0x1A40,
	0x1E00, 0xDEC1, 0xDF81, 0x1F40, 0xDD01, 0x1DC0, 0x1C80, 0xDC41,
	0x1400, 0xD4C1, 0xD581, 0x1540, 0xD701, 0x17C0, 0x1680, 0xD641,
	0xD201, 0x12C0, 0x1380, 0xD341, 0x1100, 0xD1C1, 0xD081, 0x1040,
	0xF001, 0x30C0, 0x3180, 0xF141, 0x3300, 0xF3C1, 0xF281, 0x3240,
	0x3600, 0xF6C1, 0xF781, 0x3740, 0xF501, 0x35C0, 0x3480, 0xF441,
	0x3C00, 0xFCC1, 0xFD81, 0x3D40, 0xFF01, 0x3FC0, 0x3E80, 0xFE41,
	0xFA01, 0x3AC0, 0x3B80, 0xFB41, 0x3900, 0xF9C1, 0xF881, 0x3840,
	0x2800, 0xE8C1, 0xE981, 0x2940, 0xEB01, 0x2BC0, 0x2A80, 0xEA41,
	0xEE01, 0x2EC0, 0x2F80, 0xEF41, 0x2D00, 0xEDC1, 0xEC81, 0x2C40,
	0xE401, 0x24C0, 0x2580, 0xE541, 0x2700, 0xE7C1, 0xE681, 0x2640,
	0x2200, 0xE2C1, 0xE381, 0x2340, 0xE101, 0x21C0, 0x2080, 0xE041,
	0xA001, 0x60C0, 0x6180, 0xA141, 0x6300, 0xA3C1, 0xA281, 0x6240,
	0x6600, 0xA6C1, 0xA781, 0x6740, 0xA501, 0x65C0, 0x6480, 0xA441,
	0x6C00, 0xACC1, 0xAD81, 0x6D40, 0xAF01, 0x6FC0, 0x6E80, 0xAE41,
	0xAA01, 0x6AC0, 0x6B80, 0xAB41, 0x6900, 0xA9C1, 0xA881, 0x6840,
	0x7800, 0xB8C1, 0xB981, 0x7940, 0xBB01, 0x7BC0, 0x7A80, 0xBA41,
	0xBE01, 0x7EC0, 0x7F80, 0xBF41, 0x7D00, 0xBDC1, 0xBC81, 0x7C40,
	0xB401, 0x74C0, 0x7580, 0xB541, 0x7700, 0xB7C1, 0xB681, 0x7640,
	0x7200, 0xB2C1, 0xB381, 0x7340, 0xB101, 0x71C0, 0x7080, 0xB041,
	0x5000, 0x90C1, 0x9181, 0x5140, 0x9301, 0x53C0, 0x5280, 0x9241,
	0x9601, 0x56C0, 0x5780, 0x9741, 0x5500, 0x95C1, 0x9481, 0x5440,
	0x9C01, 0x5CC0, 0x5D80, 0x9D41, 0x5F00, 0x9FC1, 0x9E81, 0x5E40,
	0x5A00, 0x9AC1, 0x9B81, 0x5B40, 0x9901, 0x59C0, 0x5880, 0x9841,
	0x8801, 0x48C0, 0x4980, 0x8941, 0x4B00, 0x8BC1, 0x8A81, 0x4A40,
	0x4E00, 0x8EC1, 0x8F81, 0x4F40, 0x8D01, 0x4DC0, 0x4C80, 0x8C41,
	0x4400, 0x84C1, 0x8581, 0x4540, 0x8701, 0x47C0, 0x4680, 0x8641,
	0x8201, 0x42C0, 0x4380, 0x8341, 0x4100, 0x81C1, 0x8081, 0x4040}
