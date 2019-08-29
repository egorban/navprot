package convertation

import (
	"github.com/ashirko/navprot/pkg/egts"
	"github.com/ashirko/navprot/pkg/general"
	"reflect"
	"testing"
)

type args struct {
	packet general.NavProtocol
	id     uint32
	packID uint16
	recID  uint16
}

func TestToEGTS(t *testing.T) {

	tests := []struct {
		name    string
		args    args
		want    *egts.Packet
		wantErr bool
	}{
		//{name: "navdata", args: navArgs(), want: navEgtsWant(), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToEGTS(tt.args.packet, tt.args.id, tt.args.packID, tt.args.recID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToEGTS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToEGTS() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func navEgtsWant() *egts.Packet {
	panic("implement me")
}

func navArgs() args {
	panic("implement me")
}
