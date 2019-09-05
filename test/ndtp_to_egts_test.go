package test

import (
	"github.com/ashirko/navprot/pkg/convertation"
	"github.com/ashirko/navprot/pkg/egts"
	"github.com/ashirko/navprot/pkg/ndtp"
	"reflect"
	"testing"
)

type args struct {
	packet *ndtp.Packet
	id     uint32
	packID uint16
	recID  uint16
}

func TestNDTPtoEGTS(t *testing.T) {
	tests := []struct {
		name     string
		args     args
		wantEgts *egts.Packet
		wantErr  bool
	}{
		{name: "navData", args: navArgs(), wantEgts: egtsNavPacket(), wantErr: false},
		{name: "fuelData", args: fuelArgs(), wantEgts: egtsFuelPacket(), wantErr: false},
		{name: "navAndFuelData", args: navAndFlueArgs(), wantEgts: egtsNavAndFuelPacket(), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			egtsPacket, err := convertation.ToEGTS(tt.args.packet, tt.args.id, tt.args.packID, tt.args.recID)
			if (err != nil) != tt.wantErr {
				t.Errorf("EGTS.parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.wantEgts, egtsPacket) {
				t.Error("got:      ", egtsPacket, "\nexpected: ", tt.wantEgts)
			}
		})
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

func egtsNavPacket() *egts.Packet {
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

func fuelArgs() args {
	return args{
		packet: ndtpFluePacket(),
		id:     0,
		packID: 0,
		recID:  0,
	}
}

func ndtpFluePacket() *ndtp.Packet {
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

func egtsFuelPacket() *egts.Packet {
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

func navAndFlueArgs() args {
	return args{
		packet: ndtpNavAndFluePacket(),
		id:     0,
		packID: 0,
		recID:  0,
	}
}

func ndtpNavAndFluePacket() *ndtp.Packet {
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

func egtsNavAndFuelPacket() *egts.Packet {
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
