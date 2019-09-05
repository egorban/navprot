package egts

import (
	"reflect"
	"testing"
)

func TestEGTS_Parse(t *testing.T) {
	tests := []struct {
		name        string
		message     []byte
		wantRestBuf []byte
		wantEgts    *Packet
		wantErr     bool
	}{
		{"result", packetResult(), []byte{1, 2, 3}, egtsRes(), false},
		{"posData", packetPosData(), nil, egtsPosData(), false},
		{"signatureNotFound", egtsWithoutSignature(), nil, new(Packet), true},
		{"shortVeryPacket", egtsVeryShort(), egtsVeryShort(), new(Packet), true},
		{"shortHeader", egtsShortHeader(), nil, new(Packet), true},
		{"shortBody", egtsShortBody(), nil, new(Packet), true},
		{"incorrectHeaderCrc", egtsIncorrectHeaderCrc(), nil, new(Packet), true},
		{"incorrectBodyCrc", egtsIncorrectBodyCrc(), nil, new(Packet), true},
		{"fuelData", packetFuelData(), nil, egtsFuelData(), false},
		{"posAndFuelData", packetPosAndFuelData(), nil, egtsPosAndFuelData(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			egts := new(Packet)
			restBuf, err := egts.Parse(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("EGTS.parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(restBuf, tt.wantRestBuf) {
				t.Errorf("EGTS.parse() = %v, want %v", restBuf, tt.wantRestBuf)
			}
			if !reflect.DeepEqual(tt.wantEgts, egts) {
				t.Error("got:      ", egts, "\nexpected: ", tt.wantEgts)
			}
		})
	}
}

func packetResult() []byte {
	return []byte{1, 0, 3, 11, 0, 16, 0, 6, 0, 0, 22, 6, 0, 0, 6, 0, 6, 0, 24, 2, 2, 0, 3, 0, 6, 0, 0, 24, 29, 1, 2, 3}
}

func egtsWithoutSignature() []byte {
	return []byte{4, 4, 4, 4, 4, 5, 6, 34, 3, 4, 5, 6, 4, 3, 3, 2, 2, 5, 6, 2, 3, 5, 6, 6, 2, 11, 4, 5, 6, 2, 41, 44}
}

func packetPosData() []byte {
	return []byte{1, 0, 0, 11, 0, 35, 0, 0, 0, 1, 153, 24, 0, 0, 0, 1, 239, 0, 0, 0, 2, 2,
		16, 21, 0, 210, 49, 43, 16, 79, 186, 58, 158, 210, 39, 188, 53, 3, 0, 0, 178, 0, 0, 0, 0, 0, 106, 141}
}

func packetFuelData() []byte {
	return []byte{1, 0, 0, 11, 0, 21, 0, 0, 0, 1, 149,
		10, 0, 0, 0, 1, 239, 0, 0, 0, 2, 2,
		27, 7, 0, 32, 0, 0, 20, 0, 0, 0,
		151, 47}
}

func packetPosAndFuelData() []byte {
	return []byte{1, 0, 0, 11, 0, 45, 0, 0, 0, 1, 47,
		34, 0, 0, 0, 1, 239, 0, 0, 0, 2, 2,
		16, 21, 0, 210, 49, 43, 16, 79, 186, 58, 158, 210, 39, 188, 53, 3, 0, 0, 178, 0, 0, 0, 0, 0,
		27, 7, 0, 32, 0, 0, 20, 0, 0, 0,
		148, 199}
}

func egtsIncorrectHeaderCrc() []byte {
	return []byte{1, 0, 0, 11, 0, 35, 0, 0, 0, 2, 153, 24, 0, 0, 0, 1, 239, 0, 0, 0, 2, 2,
		16, 21, 0, 210, 49, 43, 16, 79, 186, 58, 158, 210, 39, 188, 53, 3, 0, 0, 178, 0, 0, 0, 0, 0, 106, 141}
}

func egtsIncorrectBodyCrc() []byte {
	return []byte{1, 0, 0, 11, 0, 35, 0, 0, 0, 1, 153, 24, 0, 0, 0, 1, 239, 0, 0, 0, 2, 2,
		16, 21, 0, 210, 49, 43, 16, 79, 186, 58, 158, 211, 39, 188, 53, 3, 0, 0, 178, 0, 0, 0, 0, 0, 106, 141}
}

func egtsShortBody() []byte {
	return []byte{1, 0, 0, 11, 0, 36, 0, 0, 0, 1, 153, 24, 0, 0, 0, 1, 239, 0, 0, 0, 2, 2,
		16, 21, 0, 210, 49, 43, 16, 79, 186, 58, 158, 210, 39, 188, 53, 3, 0, 0, 178, 0, 0, 0, 0, 0, 106, 141}
}

func egtsVeryShort() []byte {
	return []byte{1, 0, 0}
}

func egtsShortHeader() []byte {
	return []byte{1, 0, 0, 10, 0, 35, 0, 0, 0, 1, 153, 24, 0, 0, 0, 1, 239, 0, 0, 0, 2, 2,
		16, 21, 0, 210, 49, 43, 16, 79, 186, 58, 158, 210, 39, 188, 53, 3, 0, 0, 178, 0, 0, 0, 0, 0, 106, 141}
}

func egtsRes() *Packet {
	data := Response{
		RPID:    6,
		ProcRes: 0,
	}
	subData := Confirmation{
		CRN: 6,
		RST: 0,
	}
	sub := SubRecord{
		Type: 0,
		Data: &subData,
	}
	rec := Record{
		RecNum:  6,
		ID:      0,
		Service: 2,
		Data:    []*SubRecord{&sub},
	}
	return &Packet{
		Type:    0,
		ID:      0,
		Records: []*Record{&rec},
		Data:    &data,
	}
}

func egtsPosData() *Packet {
	data := PosData{
		Time:    1533570258 - Timestamp20100101utc,
		Lon:     37.782409656276556,
		Lat:     55.62752532903746,
		Bearing: 178,
		Valid:   1,
	}
	sub := SubRecord{
		Type: EgtsSrPosData,
		Data: &data,
	}
	rec := Record{
		RecNum:  0,
		ID:      239,
		Service: EgtsTeledataService,
		Data:    []*SubRecord{&sub},
	}
	return &Packet{
		Type:    EgtsPtAppdata,
		ID:      0,
		Records: []*Record{&rec},
		Data:    nil,
	}
}

func egtsFuelData() *Packet {
	data := FuelData{
		Type: 2,
		Fuel: 2,
	}
	sub := SubRecord{
		Type: EgtsSrLiquidLevelSensor,
		Data: &data,
	}
	rec := Record{
		RecNum:  0,
		ID:      239,
		Service: EgtsTeledataService,
		Data:    []*SubRecord{&sub},
	}
	return &Packet{
		Type:    EgtsPtAppdata,
		ID:      0,
		Records: []*Record{&rec},
		Data:    nil,
	}
}

func egtsPosAndFuelData() *Packet {
	dataPos := PosData{
		Time:    1533570258 - Timestamp20100101utc,
		Lon:     37.782409656276556,
		Lat:     55.62752532903746,
		Bearing: 178,
		Valid:   1,
	}
	dataFule := FuelData{
		Type: 2,
		Fuel: 2,
	}
	subPos := SubRecord{
		Type: EgtsSrPosData,
		Data: &dataPos,
	}
	subFuel := SubRecord{
		Type: EgtsSrLiquidLevelSensor,
		Data: &dataFule,
	}
	rec := Record{
		RecNum:  0,
		ID:      239,
		Service: EgtsTeledataService,
		Data:    []*SubRecord{&subPos, &subFuel},
	}
	return &Packet{
		Type:    EgtsPtAppdata,
		ID:      0,
		Records: []*Record{&rec},
		Data:    nil,
	}
}

func TestPacket_Form(t *testing.T) {
	tests := []struct {
		name       string
		packetData *Packet
		wantData   []byte
		wantErr    bool
	}{
		{name: "navData", packetData: navPacket(), wantData: wantNavData(), wantErr: false},
		{name: "fuelData", packetData: fuelPacket(), wantData: wantFuelData(), wantErr: false},
		{name: "navAndFuelData", packetData: navAndFuelPacket(), wantData: wantNavAndFuelData(), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := tt.packetData.Form()
			if (err != nil) != tt.wantErr {
				t.Errorf("Form() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("Form() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func navPacket() *Packet {
	posData := &PosData{
		Time:    1533570258 - Timestamp20100101utc,
		Lon:     37.782409656276556,
		Lat:     55.62752532903746,
		Bearing: 178,
		Valid:   1,
	}
	subrec := &SubRecord{
		Type: EgtsSrPosData,
		Data: posData,
	}
	rec := &Record{
		RecNum:  0,
		ID:      239,
		Service: EgtsTeledataService,
		Data:    []*SubRecord{subrec},
	}
	return &Packet{
		Type:    EgtsPtAppdata,
		ID:      0,
		Records: []*Record{rec},
		Data:    nil,
	}
}

func wantNavData() []byte {
	return []byte{1, 0, 0, 11, 0, 35, 0, 0, 0, 1, 153, 24, 0, 0, 0, 1, 239, 0, 0, 0, 2, 2,
		16, 21, 0, 210, 49, 43, 16, 79, 186, 58, 158, 210, 39, 188, 53, 3, 0, 0, 178, 0, 0, 0, 0, 0, 106, 141}
}

func fuelPacket() *Packet {
	fuelData := &FuelData{
		Type: 2,
		Fuel: 2,
	}
	subrec := &SubRecord{
		Type: EgtsSrLiquidLevelSensor,
		Data: fuelData,
	}
	rec := &Record{
		RecNum:  0,
		ID:      239,
		Service: EgtsTeledataService,
		Data:    []*SubRecord{subrec},
	}
	return &Packet{
		Type:    EgtsPtAppdata,
		ID:      0,
		Records: []*Record{rec},
		Data:    nil,
	}
}

func wantFuelData() []byte {
	return []byte{1, 0, 0, 11, 0, 21, 0, 0, 0, 1, 149,
		10, 0, 0, 0, 1, 239, 0, 0, 0, 2, 2,
		27, 7, 0, 32, 0, 0, 20, 0, 0, 0,
		151, 47}
}

func navAndFuelPacket() *Packet {
	posData := &PosData{
		Time:    1533570258 - Timestamp20100101utc,
		Lon:     37.782409656276556,
		Lat:     55.62752532903746,
		Bearing: 178,
		Valid:   1,
	}

	fuelData := &FuelData{
		Type: 2,
		Fuel: 2,
	}

	subrec0 := &SubRecord{
		Type: EgtsSrPosData,
		Data: posData,
	}

	subrec1 := &SubRecord{
		Type: EgtsSrLiquidLevelSensor,
		Data: fuelData,
	}
	rec := &Record{
		RecNum:  0,
		ID:      239,
		Service: EgtsTeledataService,
		Data:    []*SubRecord{subrec0, subrec1},
	}
	return &Packet{
		Type:    EgtsPtAppdata,
		ID:      0,
		Records: []*Record{rec},
		Data:    nil,
	}
}

func wantNavAndFuelData() []byte {
	return []byte{1, 0, 0, 11, 0, 45, 0, 0, 0, 1, 47,
		34, 0, 0, 0, 1, 239, 0, 0, 0, 2, 2,
		16, 21, 0, 210, 49, 43, 16, 79, 186, 58, 158, 210, 39, 188, 53, 3, 0, 0, 178, 0, 0, 0, 0, 0,
		27, 7, 0, 32, 0, 0, 20, 0, 0, 0,
		148, 199}
}
