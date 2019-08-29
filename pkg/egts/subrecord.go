package egts

import (
	"encoding/binary"
	"fmt"
	"math"
)

// SubRecord describes subrecord of EGTS_PT_SIGNED_APPDATA packet
type SubRecord struct {
	// Subrecord Type
	Type byte
	// Subrecord Data
	Data interface{}
}

// Confirmation describes confirmation subrecord
type Confirmation struct {
	// Confirmed Record Number
	CRN uint16
	// Record Status
	RST byte
}

// PosData describes EGTS_SR_POS_DATA subrecord
type PosData struct {
	Time     uint32
	Lon      float64
	Lat      float64
	Bearing  uint16
	Speed    uint16
	Lohs     byte
	Lahs     byte
	Mv       byte
	RealTime byte
	Valid    byte
	Source   byte
}

type FuelData struct {
	Type byte
	Fuel uint32
}

func (subData *SubRecord) parse(service byte, buff []byte) []byte {
	subData.Type = buff[0]
	srl := binary.LittleEndian.Uint16(buff[1:3])
	subEnd := 3 + srl
	if subData.Type == egtsPtResponse {
		subData.parseResponce(buff[3:subEnd])
	}
	if service == EgtsTeledataService {
		subData.parseTeledataService(buff[3:subEnd])
	}
	return buff[subEnd:]
}

func (subData *SubRecord) parseResponce(buff []byte) {
	conf := new(Confirmation)
	conf.CRN = binary.LittleEndian.Uint16(buff[:2])
	conf.RST = buff[2]
	subData.Data = conf
}

func (subData *SubRecord) parseTeledataService(buff []byte) {
	if subData.Type == EgtsSrPosData {
		subData.parseSrPosData(buff)
	}
	//todo handle errors
}

func (subData *SubRecord) parseSrPosData(buff []byte) {
	data := new(PosData)
	lahs := buff[12] >> 5 & 1
	lohs := buff[12] >> 6 & 1
	if buff[12]&1 != 0 {
		data.Valid = 1
	}
	data.Time = binary.LittleEndian.Uint32(buff[:4])
	data.Lat = float64(binary.LittleEndian.Uint32(buff[4:8])) * 90 / 0xffffffff * (1 - 2*float64(lahs))
	data.Lon = float64(binary.LittleEndian.Uint32(buff[8:12])) * 180 / 0xffffffff * (1 - 2*float64(lohs))
	spdHi := buff[14] & 63
	spdLo := buff[13]
	data.Speed = uint16(spdHi)*256 + uint16(spdLo)
	dirHi := buff[14] >> 7
	dirLo := buff[15]
	data.Bearing = uint16(dirHi)*256 + uint16(dirLo)
	data.Source = buff[20]
	subData.Data = data
}

func (subData *SubRecord) form(service byte) (sub []byte, err error) {
	switch t := subData.Data.(type) {
	case *PosData:
		sub = subData.formPosData()
	case *Confirmation:
		sub = subData.formResponce()
	default:
		err = fmt.Errorf("subrecord type %T is not implemented", t)
	}
	return
}

func (subData *SubRecord) formPosData() (subrec []byte) {
	data := subData.Data.(*PosData)
	subrec = make([]byte, egtsSubrecDataLen+3)
	subrec[0] = byte(EgtsSrPosData)
	binary.LittleEndian.PutUint16(subrec[1:3], uint16(egtsSubrecDataLen))
	binary.LittleEndian.PutUint32(subrec[3:7], data.Time)
	lat := uint32(math.Abs(data.Lat) / 90 * 0xffffffff)
	lon := uint32(math.Abs(data.Lon) / 180 * 0xffffffff)
	binary.LittleEndian.PutUint32(subrec[7:11], lat)
	binary.LittleEndian.PutUint32(subrec[11:15], lon)
	flags := data.Lohs*64 | data.Lahs*32 | data.Mv*16 | data.RealTime*8 | 0x02 | data.Valid
	spdHi := data.Speed * 10 / 256
	spdLo := data.Speed * 10 % 256
	bearHi := data.Bearing / 256
	bearLo := data.Bearing % 256
	flags2 := ((bearHi << 0x07) | (spdHi & 0x3F)) & 0xBF //bearHi:1,0:1,spdHi:6
	subrec = append(subrec[:15], flags, byte(spdLo), byte(flags2), byte(bearLo), 0, 0, 0, 0, data.Source)
	return
}

func (subData *SubRecord) formResponce() (subrec []byte) {
	data := subData.Data.(*Confirmation)
	subrec = make([]byte, 6)
	binary.LittleEndian.PutUint16(subrec[1:3], uint16(3))
	binary.LittleEndian.PutUint16(subrec[3:5], data.CRN)
	subrec[5] = data.RST
	return
}

func (subData *SubRecord) String() string {
	header := fmt.Sprintf("{SubType: %d,", subData.Type)
	var data string
	switch subData.Data.(type) {
	case *Confirmation:
		data = subData.Data.(*Confirmation).String()
	case *PosData:
		data = subData.Data.(*PosData).String()
	default:
		data = fmt.Sprintf("%v", data)
	}
	return header + data + "}"
}

func (sub *PosData) String() string {
	return stringDefault(*sub)
}

func (sub *Confirmation) String() string {
	return stringDefault(*sub)
}

func stringDefault(v interface{}) string {
	return fmt.Sprintf("%+v", v)
}