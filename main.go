package main

/*
#cgo LDFLAGS: -lsoem

#include <stdio.h>
#include <stdlib.h>
#include <soem/ethercat.h>

*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

func main() {
	master, err := NewSOEMMaster("enp6s0f0")
	if err != nil {
		fmt.Println(err)
		return
	}
	master.ConfigInit()
	fmt.Printf("Found %d attached slaves\n", master.SlaveCount)
}

type Slave struct {
	// Manufacturer from EEprom
	VendorID uint32
	// ID from EEprom
	ProductCode uint32
	// revision from EEprom
	Revision uint32

	Name string

	// state of slave
	State uint16
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
	// output bits
	OutputBits uint16
	// output bytes, if Obits < 8 then Obytes = 0
	OutputBytes uint32
}

type SOEM struct {
	SlaveCount int
	Slaves     []Slave
	context    C.ecx_contextt
}

func NewSOEMMaster(ifname string) (*SOEM, error) {
	soem := new(SOEM)
	cifname := C.CString(ifname)
	defer C.free(unsafe.Pointer(cifname))

	soem.context = C.ecx_context

	if int(C.ecx_init(&soem.context, cifname)) <= 0 {
		return nil, fmt.Errorf("error opening interface %s", ifname)
	}

	return soem, nil
}

func (m *SOEM) ConfigInit() error {
	m.SlaveCount = int(C.ecx_config_init(&m.context, 0))
	if m.isError() {
		return errors.New("error when initialising slaves")
	}

	return nil
}

func (m *SOEM) ConfigMap() error {
	m.SlaveCount = int(C.ecx_config_init(&m.context, 0))
	if m.isError() {
		return errors.New("error when initialising slaves")
	}

	return nil
}

func (m *SOEM) getSlave(pos uint) (*Slave, error) {
	return nil, nil
}

func (m *SOEM) isError() bool {
	return int(C.ecx_iserror(&m.context)) > 0
}

func marshalSlave(cslave C.ec_slavet) *Slave {
	slave := new(Slave)
	slave.VendorID = uint32(cslave.eep_man)
	slave.ProductCode = uint32(cslave.eep_id)
	slave.Revision = uint32(cslave.eep_rev)
	slave.Name = C.GoString(&cslave.name[0])

	slave.State = uint16(cslave.state)
	slave.ALStatusCode = uint16(cslave.ALstatuscode)
	slave.AliasAddress = uint16(cslave.aliasadr)
	slave.ConfiguredAddress = uint16(cslave.configadr)
	return slave
}
