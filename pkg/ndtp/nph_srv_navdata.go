package ndtp

import (
	"encoding/binary"
	"github.com/ashirko/navprot/pkg/general"
)

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

// FuelData contains information about fuel level
type FuelData struct {
	Type byte
	Fuel uint16
}

func (data *NavData) parse(message []byte) {
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
}

func (data *NavData) toGeneral() *general.NavData {
	gen := &general.NavData{
		Time:    data.Time,
		Lon:     data.Lon,
		Lat:     data.Lat,
		Bearing: data.Bearing,
		Speed:   data.Speed,
		Valid:   data.Valid,
	}
	if data.Sos {
		gen.Source = 13
	}
	return gen
}

func (data *FuelData) parseUziM(message []byte) {
	levelMm := binary.LittleEndian.Uint16(message[3:5])
	levelL := binary.LittleEndian.Uint16(message[5:7])
	if message[2] == 0 {
		if levelL > 0 {
			data.Type = 2
			data.Fuel = levelL
		} else {
			data.Type = 0
			data.Fuel = levelMm
		}
	}
}

func (data *FuelData) parseM333(message []byte) {
	if binary.LittleEndian.Uint32(message[2:6]) != 0xFFFFFFFF {
		fuelLevel := binary.LittleEndian.Uint16(message[18:20])
		data.Type = byte(2 - (fuelLevel&0x8000)>>15)
		data.Fuel = fuelLevel & 0x7fff
	}
}
