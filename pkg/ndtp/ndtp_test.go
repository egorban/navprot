package ndtp

import (
	"reflect"
	"testing"
)

func TestNDTP_Parse(t *testing.T) {
	tests := []struct {
		name        string
		message     []byte
		wantRestBuf []byte
		wantNDTP    *Packet
		wantErr     bool
	}{
		{"navigation", packetNav(), []byte{1, 2, 3}, ndtpNav(), false},
		{"extTitle", packetExtTitle(), nil, ndtpExtTitle(), false},
		{"extResult", packetExtResult(), nil, ndtpExtResult(), false},
		{"incorrectCS", ndtpIncorrectCS(), nil, new(Packet), true},
		{"shortPacket", ndtpShort(), nil, new(Packet), true},
		{"shortVeryPacket", ndtpVeryShort(), nil, new(Packet), true},
		{"signatureNotFound", ndtpWithoutSignature(), nil, new(Packet), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ndtp := new(Packet)
			restBuf, err := ndtp.Parse(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("NDTP.parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(restBuf, tt.wantRestBuf) {
				t.Errorf("NDTP.parse() = %v, want %v", restBuf, tt.wantRestBuf)
			}
			if !reflect.DeepEqual(ndtp, tt.wantNDTP) {
				t.Error("\ngot     ", ndtp, "\nexpected", tt.wantNDTP)
			}
		})
	}
}

func ndtpIncorrectCS() []byte {
	return []byte{0, 80, 86, 161, 44, 216, 192, 140, 96, 196, 138, 54, 8, 0, 69, 0, 0, 129, 102, 160, 64, 0, 125, 6,
		18, 51, 10, 68, 41, 150, 10, 176, 70, 26, 236, 153, 35, 56, 151, 147, 73, 96, 98, 94, 76, 40, 80,
		24, 1, 2, 190, 27, 0, 0, 126, 127, 74, 0, 2, 0, 107, 210, 2, 0, 0, 0, 0, 0, 0, 1, 0, 101, 0, 1, 0, 171,
		20, 0, 0, 0, 0, 36, 141, 198, 90, 87, 110, 119, 22, 201, 186, 64, 33, 224, 203, 0, 0, 0, 0, 83, 1, 0,
		0, 220, 0, 4, 0, 2, 0, 22, 0, 67, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 167, 97, 0, 0, 31, 6, 0, 0, 8,
		0, 2, 0, 0, 0, 0, 0, 1, 2, 3}
}

func ndtpShort() []byte {
	return []byte{0, 80, 86, 161, 44, 216, 192, 140, 96, 196, 138, 54, 8, 0, 69, 0, 0, 129, 102, 160, 64, 0, 125, 6,
		18, 51, 10, 68, 41, 150, 10, 176, 70, 26, 236, 153, 35, 56, 151, 147, 73, 96, 98, 94, 76, 40, 80,
		24, 1, 2, 190, 27, 0, 0, 126, 127, 74, 0, 2, 0, 107, 210, 2, 0, 0, 0, 0, 0, 0, 1, 0, 101, 0, 1, 0, 171,
		20, 0, 0, 0, 0, 36, 141, 198, 90, 87, 110, 119, 22, 201, 186, 64, 33, 224, 203, 0, 0, 0, 0, 83, 1, 0,
		0, 220, 0, 4, 0, 2, 0, 22, 0, 67, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 167, 97, 0, 0, 31, 6, 0, 0, 8,
		0, 2, 0, 0}
}

func ndtpVeryShort() []byte {
	return []byte{0, 80, 86}
}

func ndtpWithoutSignature() []byte {
	return []byte{4, 4, 4, 4, 4, 5, 6, 34, 3, 4, 5, 6, 4, 3, 3, 2, 2, 5, 6, 2, 3, 5, 6, 6, 2, 1, 4, 5, 6, 2, 1, 1}
}

func packetNav() []byte {
	return []byte{0, 80, 86, 161, 44, 216, 192, 140, 96, 196, 138, 54, 8, 0, 69, 0, 0, 129, 102, 160, 64, 0, 125, 6,
		18, 51, 10, 68, 41, 150, 10, 176, 70, 26, 236, 153, 35, 56, 151, 147, 73, 96, 98, 94, 76, 40, 80,
		24, 1, 2, 190, 27, 0, 0, 126, 126, 74, 0, 2, 0, 107, 210, 2, 0, 0, 0, 0, 0, 0, 1, 0, 101, 0, 1, 0, 171,
		20, 0, 0, 0, 0, 36, 141, 198, 90, 87, 110, 119, 22, 201, 186, 64, 33, 224, 203, 0, 0, 0, 0, 83, 1, 0,
		0, 220, 0, 4, 0, 2, 0, 22, 0, 67, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 167, 97, 0, 0, 31, 6, 0, 0, 8,
		0, 2, 0, 0, 0, 0, 0, 1, 2, 3}
}

func packetExtTitle() []byte {
	return []byte{126, 126, 90, 1, 2, 0, 33, 134, 2, 0, 4, 0, 0, 144, 7, 5, 0, 100, 0, 0, 0, 1, 0, 0, 0, 18, 0, 0, 128, 0, 0, 0, 0, 1, 0, 0, 0, 60, 78, 65, 86, 83, 67, 82, 32, 118, 101, 114, 61, 49, 46, 48, 62, 60, 73, 68, 62, 49, 56, 60, 47, 73, 68, 62, 60, 70, 82, 79, 77, 62, 83, 69, 82, 86, 69, 82, 60, 47, 70, 82, 79, 77, 62, 60, 84, 79, 62, 85, 83, 69, 82, 60, 47, 84, 79, 62, 60, 84, 89, 80, 69, 62, 81, 85, 69, 82, 89, 60, 47, 84, 89, 80, 69, 62, 60, 77, 83, 71, 32, 116, 105, 109, 101, 61, 54, 48, 32, 98, 101, 101, 112, 61, 49, 32, 116, 121, 112, 101, 61, 98, 97, 99, 107, 103, 114, 111, 117, 110, 100, 62, 60, 98, 114, 47, 62, 60, 98, 114, 47, 62, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 194, 251, 32, 236, 229, 237, 255, 32, 241, 235, 251, 248, 232, 242, 229, 63, 60, 98, 114, 47, 62, 60, 98, 114, 47, 62, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 60, 98, 116, 110, 49, 62, 196, 224, 60, 47, 98, 116, 110, 49, 62, 60, 98, 114, 47, 62, 60, 98, 114, 47, 62, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 60, 98, 116, 110, 50, 62, 205, 229, 242, 60, 47, 98, 116, 110, 50, 62, 60, 98, 114, 47, 62, 60, 47, 77, 83, 71, 62, 60, 47, 78, 65, 86, 83, 67, 82, 62}
}

func packetExtResult() []byte {
	return []byte{126, 126, 18, 0, 2, 0, 4, 58, 2, 0, 4, 0, 0, 7, 1, 5, 0, 102, 0, 0, 0, 7, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}
}

func ndtpNav() *Packet {
	data := []interface{}{&NavData{1522961700, 37.6925783, 55.7890249, 339, 0, false, 1, 1, true}}
	nph := Nph{1, 101, true, 5291, data}
	npl := NplData{0x02, make([]byte, 4), 0x00}
	packExpected := []byte{126, 126, 74, 0, 2, 0, 107, 210, 2, 0, 0, 0, 0, 0, 0, 1, 0, 101, 0, 1, 0, 171,
		20, 0, 0, 0, 0, 36, 141, 198, 90, 87, 110, 119, 22, 201, 186, 64, 33, 224, 203, 0, 0, 0, 0, 83, 1, 0,
		0, 220, 0, 4, 0, 2, 0, 22, 0, 67, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 167, 97, 0, 0, 31, 6, 0, 0, 8,
		0, 2, 0, 0, 0, 0, 0}
	return &Packet{&npl, &nph, packExpected}
}

func ndtpExtTitle() *Packet {
	data := ExtDevice{18, 32768, 0}
	nph := Nph{NphSrvExternalDevice, nphSedDeviceTitleData, false, 1, &data}
	npl := NplData{0x02, []byte{0, 4, 0, 0}, 1936}
	packExpected := []byte{126, 126, 90, 1, 2, 0, 33, 134, 2, 0, 4, 0, 0, 144, 7, 5, 0, 100, 0, 0, 0, 1, 0, 0, 0, 18, 0, 0, 128, 0, 0, 0, 0, 1, 0, 0, 0, 60, 78, 65, 86, 83, 67, 82, 32, 118, 101, 114, 61, 49, 46, 48, 62, 60, 73, 68, 62, 49, 56, 60, 47, 73, 68, 62, 60, 70, 82, 79, 77, 62, 83, 69, 82, 86, 69, 82, 60, 47, 70, 82, 79, 77, 62, 60, 84, 79, 62, 85, 83, 69, 82, 60, 47, 84, 79, 62, 60, 84, 89, 80, 69, 62, 81, 85, 69, 82, 89, 60, 47, 84, 89, 80, 69, 62, 60, 77, 83, 71, 32, 116, 105, 109, 101, 61, 54, 48, 32, 98, 101, 101, 112, 61, 49, 32, 116, 121, 112, 101, 61, 98, 97, 99, 107, 103, 114, 111, 117, 110, 100, 62, 60, 98, 114, 47, 62, 60, 98, 114, 47, 62, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 194, 251, 32, 236, 229, 237, 255, 32, 241, 235, 251, 248, 232, 242, 229, 63, 60, 98, 114, 47, 62, 60, 98, 114, 47, 62, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 60, 98, 116, 110, 49, 62, 196, 224, 60, 47, 98, 116, 110, 49, 62, 60, 98, 114, 47, 62, 60, 98, 114, 47, 62, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 38, 110, 98, 115, 112, 59, 60, 98, 116, 110, 50, 62, 205, 229, 242, 60, 47, 98, 116, 110, 50, 62, 60, 98, 114, 47, 62, 60, 47, 77, 83, 71, 62, 60, 47, 78, 65, 86, 83, 67, 82, 62}
	return &Packet{&npl, &nph, packExpected}
}

func ndtpExtResult() *Packet {
	data := ExtDevice{MesID: 1}
	nph := Nph{NphSrvExternalDevice, nphSedDeviceResult, false, 263, &data}
	npl := NplData{DataType: 0x02, PeerAddress: []byte{0, 4, 0, 0}, ReqID: 263}
	return &Packet{&npl, &nph, packetExtResult()}
}

func TestPacket_String(t *testing.T) {
	tests := []struct {
		name       string
		packetData *Packet
		want       string
	}{
		{name: "ndtp", packetData: ndtpNav(), want: wantNdtpString()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.packetData.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func wantNdtpString() string {
	return "NPL: {DataType:2 PeerAddress:[0 0 0 0] ReqID:0}; NPH: {ServiceID:1, PacketType:101, RequestFlag:true, ReqID:5291}; Data: [ &{Time:1522961700 Lon:37.6925783 Lat:55.7890249 Bearing:339 Speed:0 Sos:false Lohs:1 Lahs:1 Valid:true} ]; Packet: [126 126 74 0 2 0 107 210 2 0 0 0 0 0 0 1 0 101 0 1 0 171 20 0 0 0 0 36 141 198 90 87 110 119 22 201 186 64 33 224 203 0 0 0 0 83 1 0 0 220 0 4 0 2 0 22 0 67 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 167 97 0 0 31 6 0 0 8 0 2 0 0 0 0 0]"
}

//NPL: {DataType:2 PeerAddress:[0 0 0 0] ReqID:0}; NPH: {ServiceID:1 PacketType:101 RequestFlag:true ReqID:5291 Data:[0xc0000581e0]}; Data: [ &{Time:1522961700 Lon:37.6925783 Lat:55.7890249 Bearing:339 Speed:0 Sos:false Lohs:1 Lahs:1 Valid:true} ]; [126 126 74 0 2 0 107 210 2 0 0 0 0 0 0 1 0 101 0 1 0 171 20 0 0 0 0 36 141 198 90 87 110 119 22 201 186 64 33 224 203 0 0 0 0 83 1 0 0 220 0 4 0 2 0 22 0 67 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 167 97 0 0 31 6 0 0 8 0 2 0 0 0 0 0]
