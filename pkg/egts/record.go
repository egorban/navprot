package egts

import (
	"encoding/binary"
	"fmt"
)

// Record describes record of EGTS_PT_SIGNED_APPDATA packet
type Record struct {
	// Record Data
	Data []*SubRecord
	// Object Identifier
	ID uint32
	// Record Number
	RecNum uint16
	// Source Service Type
	Service byte
}

func (recData *Record) form() (record []byte, err error) {
	subrec, err := recData.formSubrecords()
	if err != nil {
		return nil, err
	}
	headerRec := make([]byte, egtsRecordHeaderLen)
	binary.LittleEndian.PutUint16(headerRec[0:2], uint16(len(subrec)))
	binary.LittleEndian.PutUint16(headerRec[2:4], recData.RecNum)
	headerRec[4] = 0x01
	binary.LittleEndian.PutUint32(headerRec[5:9], recData.ID)
	headerRec = append(headerRec[:9], recData.Service, recData.Service)
	record = append(headerRec, subrec...)
	return record, nil
}

func (recData *Record) parseRecord(body []byte) []byte {
	dataLen := binary.LittleEndian.Uint16(body[:2])
	recData.RecNum = binary.LittleEndian.Uint16(body[2:4])
	tmfe := body[4] >> 2 & 1
	evfe := body[4] >> 1 & 1
	obfe := body[4] & 1
	if obfe != 0 {
		recData.ID = binary.LittleEndian.Uint32(body[5:9])
	}
	optLen := (tmfe + evfe + obfe) * 4
	headerLen := 7 + int(optLen)
	recordLen := headerLen + int(dataLen)
	recData.Service = body[5+optLen]
	sub := body[headerLen:recordLen]
	recData.parseSubRecords(sub)
	return body[recordLen:]
}

func (recData *Record) parseSubRecords(buff []byte) {
	restBuff := buff
	for len(restBuff) > 0 {
		sub := new(SubRecord)
		restBuff = sub.parse(recData.Service, restBuff)
		recData.Data = append(recData.Data, sub)
	}
}

func (recData *Record) formSubrecords() ([]byte, error) {
	var subrecords []byte
	for _, sub := range recData.Data {
		subBin, err := sub.form(recData.Service)
		if err != nil {
			return nil, err
		}
		subrecords = append(subrecords, subBin...)
	}
	return subrecords, nil
}

func (recData *Record) String() string {
	header := fmt.Sprintf("RecHeader: {Service:%d; ID:%d; RecNum:%d}, ", recData.Service, recData.ID, recData.RecNum)
	subrecords := ""
	for _, sub := range recData.Data {
		subrecords += sub.String()
	}
	return "{" + header + "[" + subrecords + "]}"
}
