package navprot

import (
	"reflect"
	"testing"
)

func TestNDTPtoEGTS(t *testing.T) {
	data := NavData{1522961700, 37.6925783, 55.7890249, 339, 0, true, 1, 1, true}
	nph := NphData{1, 101, true, 5291, &data}
	npl := NplData{0x02, make([]byte, 4), 0x00}
	ndtp := NDTP{&npl, &nph, []byte(nil)}
	egts, err := NDTPtoEGTS(ndtp, 1)
	if err != nil {
		t.Error(err)
		return
	}
	subrec := PosData{260657700, 37.6925783, 55.7890249, 339, 0, 0, 0, 0, 1, 1, 13}
	rec := EgtsRecord{0, egtsTeledataService, egtsSrPosData, &subrec}
	egtsExpected := EGTS{egtsPtAppdata, 0, 1, &rec}
	if !reflect.DeepEqual(egtsExpected, *egts) {
		t.Error("\nexpected: ", egtsExpected.Print(), "\n",
			"got:     ", egts.Print())
	}
}
