package convertation

import (
	"github.com/ashirko/navprot/pkg/egts"
	"github.com/ashirko/navprot/pkg/ndtp"
	"reflect"
	"testing"
)

func TestNDTPtoEGTS(t *testing.T) {
	ndtpPacket := ndtpNavPacket()
	egtsPacket, err := ToEGTS(ndtpPacket, 0, 0, 0)
	if err != nil {
		t.Error(err)
		return
	}
	egtsExpected := egtsNavPacket()
	if !reflect.DeepEqual(egtsExpected, egtsPacket) {
		t.Error("\nexpected: ", egtsExpected, "\n",
			"got:     ", egtsPacket)
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
		Data:        []*ndtp.NavData{&data},
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
	//return egts.Packet{egts.EgtsPtAppdata, 0, 1, &rec}
	return &egts.Packet{
		Type:    egts.EgtsPtAppdata,
		ID:      0,
		Records: []*egts.Record{&rec},
		Data:    nil,
	}
}
