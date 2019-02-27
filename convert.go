package navprot

import (
	"fmt"
)

// NDTPtoEGTS converts NDTP structure to EGTS structure
func NDTPtoEGTS(ndtp NDTP, id uint32) (egts *EGTS, err error) {
	egts = new(EGTS)
	egts.ID = id
	switch t := ndtp.Nph.Data.(type) {
	case *NavData:
		toEgtsPosData(ndtp, egts)
	default:
		err = fmt.Errorf("type %v of data is not implemented", t)
	}
	return
}

func toEgtsPosData(ndtp NDTP, egts *EGTS) {
	navData := ndtp.Nph.Data.(*NavData)
	subrec := new(PosData)
	subrec.Time = navData.Time - timestamp20100101utc
	subrec.Lat = navData.Lat
	subrec.Lon = navData.Lon
	subrec.Bearing = navData.Bearing
	if navData.Lon < 0 {
		subrec.Lohs = 1
	}
	if navData.Lat < 0 {
		subrec.Lahs = 1
	}
	if navData.Speed > 0 {
		subrec.Mv = 1
	}
	if !(PacketType(&ndtp) == NphSndHistory) {
		subrec.RealTime = 1
	}
	if navData.Valid {
		subrec.Valid = 1
	}
	if navData.Sos {
		subrec.Source = 13
	}
	record := new(EgtsRecord)
	record.Sub = subrec
	record.Service = byte(egtsTeledataService)
	record.SubType = byte(egtsSrPosData)
	egts.Data = record
	egts.PacketType = byte(egtsPtAppdata)
	return
}
