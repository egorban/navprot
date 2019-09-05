package general

import "github.com/ashirko/navprot/pkg/egts"

// FuelData is a general type for storing fuel level information
type FuelData struct {
	Type byte
	Fuel uint32
}

// ToEgtsSubrecord implement ToEGTS method of Subrecord interface
func (data FuelData) ToEgtsSubrecord() *egts.SubRecord {
	fuel := egts.FuelData{
		Type: data.Type,
		Fuel: data.Fuel,
	}
	sub := egts.SubRecord{
		Type: egts.EgtsSrLiquidLevelSensor,
		Data: &fuel,
	}
	return &sub
}
