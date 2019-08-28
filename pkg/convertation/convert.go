/*
Package convertation provides functions for converting between different navigation protocols.
Currently, convertation of general NavPtotocol to EGTS packet is supported.
*/
package convertation

import (
	"github.com/ashirko/navprot/pkg/egts"
	"github.com/ashirko/navprot/pkg/general"
)

// ToEGTS convert packet implemented NavProtocol iterface to egts.Packet type
func ToEGTS(packet general.NavProtocol, id uint32, packID, recID uint16) (*egts.Packet, error) {
	data, err := packet.ToGeneral()
	if err != nil {
		return nil, err
	}
	subrecords := egtsSubrecords(data)
	egtsPacket := formEgts(subrecords, id, packID, recID)
	return egtsPacket, nil
}

func egtsSubrecords(data []general.Subrecord) []*egts.SubRecord {
	subrecords := make([]*egts.SubRecord, 0, 1)
	for _, sub := range data {
		egtsSub := sub.ToEgtsSubrecord()
		subrecords = append(subrecords, egtsSub)
	}
	return subrecords
}

func formEgts(subrecords []*egts.SubRecord, id uint32, packID, recID uint16) *egts.Packet {
	record := egtsRecord(subrecords, id, recID)
	packetData := egtsPacketData(record, packID)
	return packetData
}

func egtsPacketData(record *egts.Record, packID uint16) *egts.Packet {
	return &egts.Packet{
		Type:    egts.EgtsPtAppdata,
		ID:      packID,
		Records: []*egts.Record{record},
		Data:    nil,
	}
}

func egtsRecord(subrecords []*egts.SubRecord, id uint32, recID uint16) *egts.Record {
	return &egts.Record{
		RecNum:  recID,
		ID:      id,
		Service: egts.EgtsTeledataService,
		Data:    subrecords,
	}
}
