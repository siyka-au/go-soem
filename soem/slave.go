package soem

/*
#cgo LDFLAGS: -lsoem

#include <stdio.h>
#include <stdlib.h>
#include <soem/ethercat.h>

*/
import "C"
import (
	"unsafe"
)

type Slave struct {
	// Manufacturer from EEprom
	VendorID uint32
	// ID from EEprom
	ProductCode uint32
	// revision from EEprom
	Revision uint32

	Name string

	// state of slave
	State EtherCATState

	// AL status code
	ALStatusCode uint16
	// Configured address
	ConfiguredAddress uint16
	// Alias address
	AliasAddress uint16

	// Interface type
	InterfaceType uint16
	// Device type
	DeviceType uint16
	// input bits
	InputBits uint16
	// input bytes, if Ibits < 8 then Ibytes = 0
	InputBytes uint32
	// output bits
	OutputBits uint16
	// output bytes, if Obits < 8 then Obytes = 0
	OutputBytes uint32

	inputBuffer  *(C.uchar)
	outputBuffer *(C.uchar)
}

func (s *Slave) Read() []byte {
	return C.GoBytes(unsafe.Pointer(s.outputBuffer), C.int(s.OutputBytes))
}

func (s *Slave) Write(data []byte) []byte {
	return nil
}
