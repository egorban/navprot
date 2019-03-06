package navprot

import (
	"reflect"
	"testing"
)

func TestEGTS_Parse(t *testing.T) {
	packet := []byte{1, 0, 3, 11, 0, 16, 0, 6, 0, 0, 22, 6, 0, 0, 6, 0, 6, 0, 24, 2, 2, 0, 3, 0, 6, 0, 0, 24, 29}
	data := EgtsResponce{6, 0}
	egtsExpected := EGTS{0, 0, 0, &data}
	var egts EGTS
	_, err := egts.Parse(packet)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(egtsExpected, egts) {
		t.Error("\nexpected: ", egtsExpected.Print(), "\n",
			"got:      ", egts.Print())
	}
}

func TestEGTS_Form(t *testing.T) {
	packetExpected := []byte{1, 0, 0, 11, 0, 35, 0, 0, 0, 1, 153, 24, 0, 0, 0, 1, 239, 0, 0, 0, 2, 2,
		16, 21, 0, 210, 49, 43, 16, 79, 186, 58, 158, 210, 39, 188, 53, 3, 0, 0, 178, 0, 0, 0, 0, 0, 106, 141}
	subrec := &PosData{1533570258-timestamp20100101utc,37.782409656276556, 55.62752532903746, 178, 0, 0, 0, 0, 0, 1, 0}
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
