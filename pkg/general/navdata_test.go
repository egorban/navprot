package general

import (
	"reflect"
	"testing"

	"github.com/egorban/navprot/pkg/egts"
)

func TestNavData_ToEgtsSubrecord(t *testing.T) {
	type fields struct {
		Time     uint32
		Lon      float64
		Lat      float64
		Bearing  uint16
		Speed    uint16
		RealTime bool
		Valid    bool
		Source   byte
	}
	tests := []struct {
		name   string
		fields fields
		want   *egts.SubRecord
	}{
		{name: "navData", fields: fields{
			Time:     1522961700,
			Lon:      37.6925783,
			Lat:      55.7890249,
			Bearing:  339,
			RealTime: true,
			Valid:    true,
			Source:   13,
		}, want: egtsExpected()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := NavData{
				Time:     tt.fields.Time,
				Lon:      tt.fields.Lon,
				Lat:      tt.fields.Lat,
				Bearing:  tt.fields.Bearing,
				Speed:    tt.fields.Speed,
				RealTime: tt.fields.RealTime,
				Valid:    tt.fields.Valid,
				Source:   tt.fields.Source,
			}
			if got := data.ToEgtsSubrecord(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToEgtsSubrecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func egtsExpected() *egts.SubRecord {
	posData := egts.PosData{
		Time:     260657700,
		Lon:      37.6925783,
		Lat:      55.7890249,
		Bearing:  339,
		RealTime: 1,
		Valid:    1,
		Source:   13,
	}
	return &egts.SubRecord{
		Type: egts.EgtsSrPosData,
		Data: &posData,
	}
}
