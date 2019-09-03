package ndtp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ashirko/navprot/pkg/general"
)

var nplSignature = []byte{0x7E, 0x7E}
var lenCells = map[byte]int{cellTypeNav: 28,
	cellTypeSensor:  28,
	cellTypeCorona:  16,
	cellTypeIRMA:    17,
	cellTypeKMD:     8,
	cellTypeCounter: 11,
	cellTypeDig:     3,
	cellTypeUziM:    8,
	cellTypeReg:     42,
	cellTypeM333:    39}

// Nph describes session layer of NDTP protocol
type Nph struct {
	ServiceID   uint16
	PacketType  uint16
	RequestFlag bool
	ReqID       uint32
	Data        interface{}
}

// Subrecord is an interface for data that can be converted into general.Subrecord
type Subrecord interface {
	toGeneral() general.Subrecord
}

func (nph *Nph) String() string {
	if nph == nil {
		return "NPH: nil;"
	}
	//sNPH := fmt.Sprintf(" NPH: %+v;", *nph)
	sNPH := fmt.Sprintf(" NPH: {ServiceID:%d, PacketType:%d, RequestFlag:%t, ReqID:%d};", nph.ServiceID, nph.PacketType, nph.RequestFlag, nph.ReqID)
	sData := sData(nph.Data)
	return sNPH + sData
}

// IsResult returns true, if packetData is a NPH_RESULT.
func (nph *Nph) isResult() bool {
	return nph.PacketType == 0
}

// Service returns value of packet's service type.
func (nph *Nph) service() int {
	return int(nph.ServiceID)
}

// PacketType returns name of NDTP packet type.
func (nph *Nph) packetType() (ptype string) {
	switch nph.ServiceID {
	case NphSrvGenericControls:
		if nph.PacketType == nphSgcConnRequest {
			ptype = NphSgsConnRequest
		}
	case NphSrvNavdata:
		if nph.PacketType == nphSndHistory {
			ptype = NphSndHistory
		} else if nph.PacketType == nphSndRealtime {
			ptype = NphSndRealtime
		}
	case NphSrvExternalDevice:
		if nph.PacketType == nphSedDeviceTitleData {
			ptype = NphSedDeviceTitleData
		} else if nph.PacketType == nphSedDeviceResult {
			ptype = NphSedDeviceResult
		}
	}
	return
}

func (nph *Nph) parse(message []byte) (err error) {
	nph.ServiceID = binary.LittleEndian.Uint16(message[:2])
	nph.PacketType = binary.LittleEndian.Uint16(message[2:4])
	if binary.LittleEndian.Uint16(message[4:6]) == 1 {
		nph.RequestFlag = true
	}
	nph.ReqID = binary.LittleEndian.Uint32(message[6:10])
	if nph.isResult() {
		nph.Data = binary.LittleEndian.Uint32(message[nphHeaderLen : nphHeaderLen+4])
		return
	}
	switch nph.service() {
	case NphSrvGenericControls:
		nph.parseGenControl(message[nphHeaderLen:])
	case NphSrvNavdata:
		err = nph.parseNavData(message[nphHeaderLen:])
	case NphSrvExternalDevice:
		err = nph.parseExtDevice(message[nphHeaderLen:])
	default:
		err = errors.New("unknown service")
	}
	return
}

func (nph *Nph) parseNavData(message []byte) (err error) {
	cellStart := 0
	allData := make([]interface{}, 0, 1)
	for message[cellStart] <= cellTypeM333 {
		cellType := message[cellStart]
		switch cellType {
		case cellTypeNav:
			if len(message[cellStart:]) >= lenCells[cellTypeNav] {
				data := new(NavData)
				data.parse(message[cellStart:])
				allData = append(allData, data)
				cellStart = cellStart + lenCells[cellTypeNav]
			} else {
				err = errors.New("NavData type 0 is too short")
				return
			}
		case cellTypeUziM:
			if len(message[cellStart:]) >= lenCells[cellTypeUziM] {
				data := new(FuelData)
				data.parse_UziM(message[cellStart:])
				allData = append(allData, data)
				cellStart = cellStart + lenCells[cellTypeUziM]
			} else {
				err = errors.New("NavData type 8 is too short")
				return
			}
		case cellTypeM333:
			if len(message[cellStart:]) >= lenCells[cellTypeM333] {
				data := new(FuelData)
				data.parse_M333(message[cellStart:])
				allData = append(allData, data)
				cellStart = cellStart + lenCells[cellTypeM333]
			} else {
				err = errors.New("NavData type 10 is too short")
				return
			}
		default:
			if len(message[cellStart:]) >= lenCells[cellType] {
				cellStart = cellStart + lenCells[cellType]
			} else {
				return
			}
		}
		if cellStart >= len(message) {
			break
		}
	}
	nph.Data = allData
	return
}

func (nph *Nph) parseGenControl(message []byte) {
	if nph.packetType() == NphSgsConnRequest {
		nph.Data = binary.LittleEndian.Uint32(message[6:10])
	}
	return
}

func (nph *Nph) parseExtDevice(message []byte) (err error) {
	ext := new(ExtDevice)
	err = ext.parse(nph.packetType(), message)
	if err == nil {
		nph.Data = ext
	}
	return
}

func sData(data interface{}) (sdata string) {
	prefix := " Data: "
	switch data.(type) {
	case int:
		sdata = fmt.Sprintf("%d", data)
	case *ExtDevice:
		ext := data.(*ExtDevice)
		sdata = fmt.Sprintf("%+v", *ext)
	case []interface{}:
		tmp := "["
		for _, e := range data.([]interface{}) {
			tmp = tmp + fmt.Sprintf(" %+v", e)
		}
		sdata = sdata + tmp + " ]"
	default:
		sdata = fmt.Sprintf("%v", data)
	}
	return prefix + sdata
}
