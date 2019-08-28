/*
Package general provides general interfaces ant types for working with telematic data.
*/
package general

import "github.com/ashirko/navprot/pkg/egts"

// NavProtocol is an interface for arbitrary navigation protocol.
type NavProtocol interface {
	Parse([]byte) ([]byte, error)
	ToGeneral() ([]Subrecord, error)
	String() string
}

// Subrecord is an interface for data that can be converted into EGTS subrecord
type Subrecord interface {
	ToEgtsSubrecord() *egts.SubRecord
}
