package general

import (
	"github.com/ashirko/navprot/pkg/egts"
)

// NavData is a general type for storing navigation information
type NavData struct {
	Time     uint32
	Lon      float64
	Lat      float64
	Bearing  uint16
	Speed    uint16
	RealTime bool
	Valid    bool
	// Source 13 - sos
	Source byte
}

// ToEgtsSubrecord implement ToEGTS method of Subrecord interface
func (data NavData) ToEgtsSubrecord() *egts.SubRecord {
	nav := egts.PosData{
		Time:    data.Time - egts.Timestamp20100101utc,
		Lon:     data.Lon,
		Lat:     data.Lat,
		Bearing: data.Bearing,
		Speed:   data.Speed,
		Source:  data.Source,
	}
	if data.Lat < 0 {
		nav.Lahs = 1
	}
	if data.Lon < 0 {
		nav.Lohs = 1
	}
	if data.Speed > 0 {
		nav.Mv = 1
	}
	if data.RealTime {
		nav.RealTime = 1
	}
	if data.Valid {
		nav.Valid = 1
	}
	sub := egts.SubRecord{
		Type: egts.EgtsSrPosData,
		Data: &nav,
	}
	return &sub
}
