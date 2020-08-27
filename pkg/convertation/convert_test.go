package convertation

import (
	"reflect"
	"testing"

	"github.com/egorban/navprot/pkg/egts"
	"github.com/egorban/navprot/pkg/general"
	"github.com/egorban/navprot/pkg/ndtp"
)

type args struct {
	packet general.NavProtocol
	id     uint32
	packID uint16
	recID  uint16
}

func TestToEGTS(t *testing.T) {

	tests := []struct {
		name    string
		args    args
		want    *egts.Packet
		wantErr bool
	}{
		{name: "navdata", args: navArgs(), want: navEgtsWant(), wantErr: false},
		{name: "fueldata", args: fuelArgs(), want: fuelEgtsWant(), wantErr: false},
		{name: "navAndFueldata", args: navAndFuelArgs(), want: navAndFuelEgtsWant(), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToEGTS(tt.args.packet, tt.args.id, tt.args.packID, tt.args.recID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToEGTS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToEGTS() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func navEgtsWant() *egts.Packet {
	posData := egts.PosData{
		Time:     260657700,
		Lon:      37.6925783,
		Lat:      55.7890249,
		Bearing:  339,
		RealTime: 1,
		Valid:    1,
		Source:   13,
	}
	subrec := egts.SubRecord{
		Type: egts.EgtsSrPosData,
		Data: &posData,
	}
	rec := egts.Record{
		RecNum:  0,
		ID:      0,
		Service: egts.EgtsTeledataService,
		Data:    []*egts.SubRecord{&subrec},
	}
	return &egts.Packet{
		Type:    egts.EgtsPtAppdata,
		ID:      0,
		Records: []*egts.Record{&rec},
		Data:    nil,
	}
}

func navArgs() args {
	return args{
		packet: ndtpNavPacket(),
		id:     0,
		packID: 0,
		recID:  0,
	}
}

func ndtpNavPacket() *ndtp.Packet {
	data := ndtp.NavData{
		Time:    1522961700,
		Lon:     37.6925783,
		Lat:     55.7890249,
		Bearing: 339,
		Sos:     true,
		Lohs:    1,
		Lahs:    1,
		Valid:   true,
	}
	nph := ndtp.Nph{
		ServiceID:   1,
		PacketType:  101,
		RequestFlag: true,
		ReqID:       5291,
		Data:        []ndtp.Subrecord{&data},
	}
	npl := ndtp.NplData{
		DataType:    2,
		PeerAddress: make([]byte, 4),
		ReqID:       0,
	}
	return &ndtp.Packet{
		Npl:    &npl,
		Nph:    &nph,
		Packet: []byte(nil),
	}
}

func fuelEgtsWant() *egts.Packet {
	flueData := egts.FuelData{
		Type: 1,
		Fuel: 20,
	}
	subrec := egts.SubRecord{
		Type: egts.EgtsSrLiquidLevelSensor,
		Data: &flueData,
	}
	rec := egts.Record{
		RecNum:  0,
		ID:      0,
		Service: egts.EgtsTeledataService,
		Data:    []*egts.SubRecord{&subrec},
	}
	return &egts.Packet{
		Type:    egts.EgtsPtAppdata,
		ID:      0,
		Records: []*egts.Record{&rec},
		Data:    nil,
	}
}

func fuelArgs() args {
	return args{
		packet: ndtpFuelPacket(),
		id:     0,
		packID: 0,
		recID:  0,
	}
}

func ndtpFuelPacket() *ndtp.Packet {
	data := ndtp.FuelData{
		Type: 1,
		Fuel: 20,
	}
	nph := ndtp.Nph{
		ServiceID:   1,
		PacketType:  101,
		RequestFlag: true,
		ReqID:       5291,
		Data:        []ndtp.Subrecord{&data},
	}
	npl := ndtp.NplData{
		DataType:    2,
		PeerAddress: make([]byte, 4),
		ReqID:       0,
	}
	return &ndtp.Packet{
		Npl:    &npl,
		Nph:    &nph,
		Packet: []byte(nil),
	}
}

func navAndFuelEgtsWant() *egts.Packet {
	posData := egts.PosData{
		Time:     260657700,
		Lon:      37.6925783,
		Lat:      55.7890249,
		Bearing:  339,
		RealTime: 1,
		Valid:    1,
		Source:   13,
	}
	subrecNav := egts.SubRecord{
		Type: egts.EgtsSrPosData,
		Data: &posData,
	}
	flueData := egts.FuelData{
		Type: 1,
		Fuel: 20,
	}
	subrecFlue := egts.SubRecord{
		Type: egts.EgtsSrLiquidLevelSensor,
		Data: &flueData,
	}
	rec := egts.Record{
		RecNum:  0,
		ID:      0,
		Service: egts.EgtsTeledataService,
		Data:    []*egts.SubRecord{&subrecNav, &subrecFlue},
	}
	return &egts.Packet{
		Type:    egts.EgtsPtAppdata,
		ID:      0,
		Records: []*egts.Record{&rec},
		Data:    nil,
	}
}

func navAndFuelArgs() args {
	return args{
		packet: ndtpNavAndFuelPacket(),
		id:     0,
		packID: 0,
		recID:  0,
	}
}

func ndtpNavAndFuelPacket() *ndtp.Packet {
	dataNav := ndtp.NavData{
		Time:    1522961700,
		Lon:     37.6925783,
		Lat:     55.7890249,
		Bearing: 339,
		Sos:     true,
		Lohs:    1,
		Lahs:    1,
		Valid:   true,
	}
	dataFlue := ndtp.FuelData{
		Type: 1,
		Fuel: 20,
	}
	nph := ndtp.Nph{
		ServiceID:   1,
		PacketType:  101,
		RequestFlag: true,
		ReqID:       5291,
		Data:        []ndtp.Subrecord{&dataNav, &dataFlue},
	}
	npl := ndtp.NplData{
		DataType:    2,
		PeerAddress: make([]byte, 4),
		ReqID:       0,
	}
	return &ndtp.Packet{
		Npl:    &npl,
		Nph:    &nph,
		Packet: []byte(nil),
	}
}
