/*
Package ndtp provides functions for parsing and forming NDTP packets.
*/
package ndtp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ashirko/navprot/pkg/general"
)

// Packet contains information about about NDTP (Navigation Data Transfer Protocol) packet
type Packet struct {
	Npl    *NplData
	Nph    *Nph
	Packet []byte
}

const (
	nplHeaderLen     = 15
	nphResult        = 0
	nphHeaderLen     = 10
	ndtpResultLen    = nphHeaderLen + nplHeaderLen + 4
	ndtpExtResultLen = nphHeaderLen + nplHeaderLen + 8
	navDataCellLen   = 28

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
func (packetData *Packet) Parse(message []byte) (restBuf []byte, err error) {
	packetData.Npl, packetData.Packet, restBuf, err = parseNPL(message)
	if err != nil {
		return
	}
	packetData.Nph = new(Nph)
	err = packetData.Nph.parse(packetData.Packet[nplHeaderLen:])
	return
}

// String generate string with information about NDTP packet in readable format.
func (packetData Packet) String() string {
	sNPL := packetData.Npl.String()
	sNPH := packetData.Nph.String()
	packet := fmt.Sprintf("; Packet: %v", packetData.Packet)
	return sNPL + sNPH + packet
}

// IsResult returns true, if packetData is a NPH_RESULT.
func (packetData *Packet) IsResult() bool {
	return packetData.Nph.isResult()
}

// Service returns value of packet's service type.
func (packetData *Packet) Service() int {
	return packetData.Nph.service()
}

// PacketType returns name of NDTP packet type.
func (packetData *Packet) PacketType() (ptype string) {
	return packetData.Nph.packetType()
}

// ChangeAddress change Peer Address field of NPL layer
func (packetData *Packet) ChangeAddress(newAddress []byte) {
	for i, j := 9, 0; i < 13; i, j = i+1, j+1 {
		packetData.Packet[i] = newAddress[j]
	}
}

// GetID returns ID of terminal, which is included only in NPH_SGC_CONN_REQUEST packets
func (packetData *Packet) GetID() (id int, err error) {
	if packetData.PacketType() == NphSgsConnRequest {
		id = int(packetData.Nph.Data.(uint32))
	} else {
		err = errors.New("incorrect packet type")
	}
	return
}

// NeedReply returns true if packet needs reply
func (packetData *Packet) NeedReply() (flag bool) {
	return packetData.Nph.RequestFlag
}

// Reply creates NPH_RESULT packet for packetData.Packet
func (packetData *Packet) Reply(result uint32) []byte {
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

// ReplyExt creates NPH_SED_DEVICE_RESULT packet
func (packetData *Packet) ReplyExt(result uint32) ([]byte, error) {
	if packetData.Service() == NphSrvExternalDevice {
		reply := packetData.Nph.Data.(*ExtDevice).reply(packetData.Packet, result)
		return reply, nil
	}
	return nil, errors.New("incorrect packet service")
}

// ChangePacket changes values of some fields in NDTP packet
func (packetData *Packet) ChangePacket(changes map[string]int) {
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

// ToGeneral form general subrecords from NDTP packet
func (packetData *Packet) ToGeneral() (subrecords []general.Subrecord, err error) {
	if packetData.Service() == NphSrvNavdata {
		for _, sub := range packetData.Nph.Data.([]*NavData) { //TODO fix type assertion
			gen := sub.toGeneral()
			if packetData.PacketType() == NphSndRealtime {
				gen.RealTime = true
			}
			subrecords = append(subrecords, gen)
		}
	} else {
		err = errors.New("incorrect packet type")
	}
	return
}