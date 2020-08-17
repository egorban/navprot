package test

import (
	"reflect"
	"testing"

	"github.com/egorban/navprot/pkg/convertation"
	"github.com/egorban/navprot/pkg/egts"
	"github.com/egorban/navprot/pkg/ndtp"
)

type args struct {
	packet *ndtp.Packet
	id     uint32
	packID uint16
	recID  uint16
}

type argsBin struct {
	packet []byte
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

func TestNDTPtoEGTSBin(t *testing.T) {
	tests := []struct {
		name     string
		args     argsBin
		wantEgts []byte
		wantErr  bool
	}{
		{name: "navData", args: ndtpNavArgsBin(), wantEgts: egtsNavBin(), wantErr: false},
		//{name: "fuelData", args: fuelArgs(), wantEgts: egtsFuelPacket(), wantErr: false},
		//{name: "navAndFuelData", args: navAndFlueArgs(), wantEgts: egtsNavAndFuelPacket(), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ndtpData := new(ndtp.Packet)
			_, err := ndtpData.Parse(tt.args.packet)
			//fmt.Printf("parsed NDTP: %v\n", ndtpData)
			if (err != nil) != tt.wantErr {
				t.Errorf("ndtpData.parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			egtsData, err := convertation.ToEGTS(ndtpData, tt.args.id, tt.args.packID, tt.args.recID)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertation.ToEGTS() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				return
			}
			res, err := egtsData.Form()
			//e := new(egts.Packet)
			//_, err = e.Parse(res)
			//fmt.Printf("parsed NDTP: %v, %v\n", err, e)
			if (err != nil) != tt.wantErr {
				t.Errorf("egtsData.Form() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.wantEgts, res) {
				t.Error("got:      ", res, "\nexpected: ", tt.wantEgts)
			}
		})
	}
}

func ndtpNavArgsBin() argsBin {
	return argsBin{
		packet: ndtpNavBin(),
		id:     0,
		packID: 0,
		recID:  0,
	}
}

func ndtpNavBin() []byte {
	return []byte{0, 80, 86, 161, 44, 216, 192, 140, 96, 196, 138, 54, 8, 0, 69, 0, 0, 129, 102, 160, 64, 0, 125, 6,
		18, 51, 10, 68, 41, 150, 10, 176, 70, 26, 236, 153, 35, 56, 151, 147, 73, 96, 98, 94, 76, 40, 80,
		24, 1, 2, 190, 27, 0, 0, 126, 126, 74, 0, 2, 0, 107, 210, 2, 0, 0, 0, 0, 0, 0, 1, 0, 101, 0, 1, 0, 171,
		20, 0, 0, 0, 0, 36, 141, 198, 90, 87, 110, 119, 22, 201, 186, 64, 33, 224, 203, 0, 0, 0, 0, 83, 1, 0,
		0, 220, 0, 4, 0, 2, 0, 22, 0, 67, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 167, 97, 0, 0, 31, 6, 0, 0, 8,
		0, 2, 0, 0, 0, 0, 0, 1, 2, 3}
}

func egtsNavBin() []byte {
	return []byte{1, 0, 0, 11, 0, 45, 0, 0, 0, 1, 47, 34, 0, 0, 0, 1, 0, 0, 0, 0, 2, 2, 16, 21, 0, 36, 82, 137, 15, 2,
		84, 176, 158, 238, 114, 155, 53, 11, 0, 128, 83, 0, 0, 0, 0, 0, 27, 7, 0, 66, 1, 0, 0, 0, 0, 0, 99, 77}
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
