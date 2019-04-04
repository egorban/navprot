package navprot

import (
	"reflect"
	"testing"
)

func TestEGTS_Form(t *testing.T) {
	packetExpected := []byte{1, 0, 0, 11, 0, 35, 0, 0, 0, 1, 153, 24, 0, 0, 0, 1, 239, 0, 0, 0, 2, 2,
		16, 21, 0, 210, 49, 43, 16, 79, 186, 58, 158, 210, 39, 188, 53, 3, 0, 0, 178, 0, 0, 0, 0, 0, 106, 141}
	subrec := &PosData{1533570258 - timestamp20100101utc, 37.782409656276556, 55.62752532903746, 178, 0, 0, 0, 0, 0, 1, 0}
	rec := &EgtsRecord{0, egtsTeledataService, egtsSrPosData, subrec}
	egts := EGTS{egtsPtAppdata, 0, 239, rec}
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
		wantEgts    *EGTS
		wantErr     bool
	}{
		{"result", packetResult(), []byte{1, 2, 3}, egtsRes(), false},
		{"posData", packetPosData(), nil, egtsPosData(), false},
		{"signatureNotFound", egtsWithoutSignature(), nil, new(EGTS), true},
		{"shortVeryPacket", egtsVeryShort(), egtsVeryShort(), new(EGTS), true},
		{"shortHeader", egtsShortHeader(), nil, new(EGTS), true},
		{"shortBody", egtsShortBody(), nil, new(EGTS), true},
		{"incorrectHeaderCrc", egtsIncorrectHeaderCrc(), nil, new(EGTS), true},
		{"incorrectBodyCrc", egtsIncorrectBodyCrc(), nil, new(EGTS), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			egts := new(EGTS)
			restBuf, err := egts.Parse(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("EGTS.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(restBuf, tt.wantRestBuf) {
				t.Errorf("EGTS.Parse() = %v, want %v", restBuf, tt.wantRestBuf)
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

func egtsRes() *EGTS {
	data := EgtsResponce{RecID: 6}
	return &EGTS{Data: &data}
}

func egtsPosData() *EGTS {
	subrec := &PosData{Time: 1533570258 - timestamp20100101utc, Lon: 37.782409656276556, Lat: 55.62752532903746, Bearing: 178, Valid: 1}
	rec := &EgtsRecord{Service: egtsTeledataService, SubType: egtsSrPosData, Sub: subrec}
	return &EGTS{PacketType: egtsPtAppdata, ID: 239, Data: rec}
}
