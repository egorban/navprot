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
		t.Error("For packet", packet, "\n",
			"expected: ", egtsExpected.Print(), "\n",
			"got:      ", egts.Print())
	}
}
