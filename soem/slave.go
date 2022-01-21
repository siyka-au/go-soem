package soem

/*
#cgo LDFLAGS: -lsoem

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <soem/ethercat.h>

*/
import "C"
import (
	"fmt"
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

	PDO *SlavePDO

	HasDC bool

	D uint8
}

type SlavePDO struct {
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
	if s.PDO != nil {
		l := s.PDO.InputBytes
		if s.PDO.InputBytes < 1 && s.PDO.InputBits > 0 {
			l = 1
		}
		return C.GoBytes(unsafe.Pointer(s.PDO.inputBuffer), C.int(l))
	}
	return nil
}

func (s *Slave) Write(data []byte) {
	if s.PDO != nil {
		l := s.PDO.OutputBytes
		if s.PDO.OutputBytes < 1 && s.PDO.OutputBits > 0 {
			l = 1
		}
		l = l

		C.memcpy(unsafe.Pointer(s.PDO.outputBuffer), unsafe.Pointer(&data[0]), C.size_t(len(data)))
	}
}

func (slave *Slave) String() string {
	return fmt.Sprintf(
		"Name %s\n"+
			"  Vendor ID 0x%08x\n"+
			"  Product Code 0x%08x\n"+
			"  Revision 0x%08x\n"+
			"  Configured Address 0x%04x\n"+
			"  Alias Address 0x%04x\n"+
			"  Input Bits %d\n"+
			"  Input Bytes %d\n"+
			"  Output Bits %d\n"+
			"  Output Bytes %d\n"+
			"  Has DC %s (%d)\n",
		slave.Name, slave.VendorID, slave.ProductCode, slave.Revision,
		slave.ConfiguredAddress, slave.AliasAddress,
		slave.PDO.InputBits, slave.PDO.InputBytes,
		slave.PDO.OutputBits, slave.PDO.OutputBytes,
		stringSelect(slave.HasDC, "Yes", "No"), slave.D)
}

func stringSelect(selector bool, trueStr, falseStr string) string {
	if selector {
		return trueStr
	}
	return falseStr
}
