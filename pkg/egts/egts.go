/*
Package egts provides functions for parsing and forming EGTS packets.
*/
package egts

import (
	"encoding/binary"
	"fmt"
)

const egtsPtResponse byte = 0

const (
	prvSignature          = 0x01
	minEgtsHeaderLen      = 11
	egtsRecordHeaderLen   = 11
	egtsSubrecDataLen     = 21
	egtsSubrecFuelDataLen = 7

	// EgtsPtResponse defines EGTS_PT_RESPONSE packet type
	EgtsPtResponse = 0
	// EgtsPtAppdata defines EGTS_PT_APPDATA packet type
	EgtsPtAppdata = 1
	// EgtsTeledataService defines EGTS_TELEDATA_SERVICE
	EgtsTeledataService = 2
	// EgtsSrPosData defines EGTS_SR_POS_DATA subrecord
	EgtsSrPosData = 16
	// EgtsSrLiquidLevelSensor defines EGTS_SR_LIQUID_LEVEL_SENSOR subrecord
	EgtsSrLiquidLevelSensor = 27
	// EgtsSrResponse defines EGTS_SR_RECORD_RESPONSE subrecord
	EgtsSrResponse = 0
	// Timestamp20100101utc is EGTS initial time
	Timestamp20100101utc = 1262304000
)

// Packet contains information about about EGTS protocol (ERA GLONASS Telematics Standard) packet
type Packet struct {
	// Packet Type
	Type byte
	// Packet Identifier
	ID uint16
	// Service Data Records
	Records []*Record
	// Additional Data (optional)
	Data interface{}
}

// Response describes EGTS_PT_RESPONSE packet
type Response struct {
	// Response Packet ID
	RPID uint16
	// Processing Result
	ProcRes byte
}

// Parse EGTS packet. Parsed information is stored in variable with EGTS type.
func (packetData *Packet) Parse(message []byte) (restBuf []byte, err error) {
	body, restBuf, err := packetData.parseHeader(message)
	if err != nil {
		return
	}
	switch packetData.Type {
	case egtsPtResponse:
		packetData.parseResponce(body)
	case EgtsPtAppdata:
		packetData.parseAppData(body)
	default:
		err = fmt.Errorf("packet type %d not implemented", packetData.Type)
		return
	}
	return
}

// Form generate EGTS binary packet.
func (packetData *Packet) Form() (data []byte, err error) {
	data, err = formData(packetData)
	header := []byte{0x01, 0x00, 0x00, byte(minEgtsHeaderLen), 0x00, 0x00, 0x00, 0x00, 0x00, packetData.Type}
	if packetData.Type == EgtsPtResponse {
		header[2] = 0x03
	}
	binary.LittleEndian.PutUint16(header[5:7], uint16(len(data)))
	binary.LittleEndian.PutUint16(header[7:9], packetData.ID)
	crcPacket := crc8EGTS(header)
	header = append(header, byte(crcPacket))
	crcRec := make([]byte, 2)
	binary.LittleEndian.PutUint16(crcRec, crc16EGTS(data))
	data = append(data, crcRec...)
	data = append(header, data...)
	return
}

func formData(packetData *Packet) (data []byte, err error) {
	switch packetData.Type {
	case EgtsPtAppdata:
		data, err = packetData.formAppData()
	case EgtsPtResponse:
		data, err = packetData.formResponse()
	default:
		err = fmt.Errorf("data type %d not implemented", packetData.Type)
	}
	return
}

// Print generate string with information about EGTS packet in readable format.
func (packetData Packet) String() string {
	h := fmt.Sprintf("Header: {PacketType:%d; ID:%d}; ", packetData.Type, packetData.ID)
	b := packetData.data2String()
	return h + b
}

func (packetData *Packet) parseResponce(body []byte) {
	recp := new(Response)
	recp.RPID = binary.LittleEndian.Uint16(body[:2])
	recp.ProcRes = body[2]
	packetData.Records = parseRecords(body[3:])
	packetData.Data = recp
}

func (packetData *Packet) parseAppData(body []byte) {
	packetData.Records = parseRecords(body)
}

func (packetData *Packet) formAppData() ([]byte, error) {
	var packet []byte
	for _, rec := range packetData.Records {
		recBin, err := rec.form()
		if err != nil {
			return nil, err
		}
		packet = append(packet, recBin...)
	}
	return packet, nil
}

func (packetData *Packet) formResponse() ([]byte, error) {
	packet := make([]byte, 3)
	binary.LittleEndian.PutUint16(packet[0:2], packetData.Data.(*Response).RPID)
	for _, rec := range packetData.Records {
		recBin, err := rec.formResponse()
		if err != nil {
			return nil, err
		}
		packet = append(packet, recBin...)
	}
	return packet, nil
}

func parseRecords(body []byte) []*Record {
	records := make([]*Record, 0, 1)
	restBuff := body
	for len(restBuff) > 0 {
		recData := new(Record)
		restBuff = recData.parseRecord(restBuff)
		records = append(records, recData)
	}
	return records
}

func (packetData *Packet) data2String() (body string) {
	prefix := ""
	if packetData.Type == egtsPtResponse {
		prefix = prefix + "{Confirmation: " + fmt.Sprintf("%+v", *packetData.Data.(*Response)) + "}; "
	}
	prefix = prefix + "Records: "
	for _, rec := range packetData.Records {
		s := rec.String()
		body += s
	}
	return prefix + body
}
