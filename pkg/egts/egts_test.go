package egts

import (
	"reflect"
	"testing"
)

func TestEGTS_Form(t *testing.T) {
	packetExpected := []byte{1, 0, 0, 11, 0, 35, 0, 0, 0, 1, 153, 24, 0, 0, 0, 1, 239, 0, 0, 0, 2, 2,
		16, 21, 0, 210, 49, 43, 16, 79, 186, 58, 158, 210, 39, 188, 53, 3, 0, 0, 178, 0, 0, 0, 0, 0, 106, 141}
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
	egts := Packet{
		Type:    EgtsPtAppdata,
		ID:      0,
		Records: []*Record{rec},
		Data:    nil,
	}
	packet, err := egts.Form()
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(packetExpected, packet) {
		t.Error("\nexpected: ", packetExpected, "\n",
			"\ngot:      ", packet)
	}
}

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
