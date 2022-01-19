package soem

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

type Master struct {
	SlaveCount uint16
	Slaves     []*Slave

	context   C.ecx_contextt
	ioMap     unsafe.Pointer
	ioMapSize C.int
}

// TODO Work out
func NewSOEMMaster(ifname string) (*Master, error) {
	soem := new(Master)
	cifname := C.CString(ifname)
	defer C.free(unsafe.Pointer(cifname))

	soem.context = C.ecx_context // no idea why this is needed or even works

	if C.ecx_init(&soem.context, cifname) <= 0 {
		return nil, fmt.Errorf("error opening interface %s", ifname)
	}

	return soem, nil
}

func (m *Master) Close() {
	C.ecx_close(&m.context)
	C.free(unsafe.Pointer(m.ioMap))
}

func (m *Master) ConfigInit() {
	m.SlaveCount = uint16(C.ecx_config_init(&m.context, 0))
	m.Slaves = make([]*Slave, m.SlaveCount)

	for i := 0; i < int(m.SlaveCount); i++ {
		// Do stuff to update the slaves
		slave := new(Slave)
		cslave := C.ec_slave[i+1]

		slave.VendorID = uint32(cslave.eep_man)
		slave.ProductCode = uint32(cslave.eep_id)
		slave.Revision = uint32(cslave.eep_rev)
		slave.Name = C.GoString(&cslave.name[0])

		slave.State = EtherCATState(cslave.state)
		slave.ALStatusCode = uint16(cslave.ALstatuscode)
		slave.AliasAddress = uint16(cslave.aliasadr)
		slave.ConfiguredAddress = uint16(cslave.configadr)

		m.Slaves[i] = slave
	}
}

func (m *Master) ConfigMapWithGroup(group uint8, size uint) {
	m.ioMap = C.malloc(C.size_t(size))
	m.ioMapSize = C.ecx_config_map_group(&m.context, m.ioMap, C.uchar(group))

	for i, s := range m.Slaves {
		cslave := C.ec_slave[i+1]

		pdo := SlavePDO{
			uint16(cslave.Ibits),
			uint32(cslave.Ibytes),
			uint16(cslave.Obits),
			uint32(cslave.Obytes),
			cslave.inputs,
			cslave.outputs}

		m.Slaves[i].PDO = &pdo
		fmt.Println(s)
	}
}

func (m *Master) ConfigMap(size uint) {
	m.ConfigMapWithGroup(0, size)
}

func (m *Master) ReadState() int {
	return int(C.ecx_readstate(&m.context))
}

func (m *Master) SetState(state EtherCATState) (uint, error) {
	C.ec_slave[0].state = C.ushort(state)
	ret := C.ecx_writestate(&m.context, 0)
	if ret < 0 {
		switch ret {
		case EC_NOFRAME:
			return 0, errors.New("EC_NOFRAME")
		default:
			return 0, errors.New("undefined error")
		}
	}

	return uint(ret), nil
}

func (m *Master) CheckState(slave uint16, expectedState EtherCATState, timeout int) (EtherCATState, error) {
	state := EtherCATState(int(C.ecx_statecheck(&m.context,
		C.ushort(slave),
		C.ushort(expectedState),
		C.int(timeout))))

	if state != expectedState {
		return state, fmt.Errorf("current state %s not as expected state %s", state, expectedState)
	}

	return state, nil
}

func (m *Master) SendProcessDataWithGroup(group uint8) {
	C.ecx_send_processdata(&m.context)
}

func (m *Master) SendProcessData() {
	m.SendProcessDataWithGroup(0)
}

func (m *Master) ReceiveProcessDataWithGroup(group uint8, timeout int) uint {
	return uint(C.ecx_receive_processdata(&m.context, C.int(timeout)))
}

func (m *Master) ReceiveProcessData(timeout int) uint {
	return m.ReceiveProcessDataWithGroup(0, timeout)
}

// func (m *SOEM) RecoverSlave(uint16 slave, int timeout) error {
// 	if C.ecx_recover_slave(&m.context) <= 0 {
// 		return errors.New("error recovering slave")
// 	}
// 	return nil
// }

func (m *Master) isError() bool {
	return int(C.ecx_iserror(&m.context)) > 0
}

func marshalSlave(cslave C.ec_slavet) *Slave {
	slave := new(Slave)

	ptr := unsafe.Pointer(cslave.inputs)
	fmt.Print("Inputs Pointer: ")
	fmt.Println(ptr)

	ptr = unsafe.Pointer(cslave.outputs)
	fmt.Print("Outputs Pointer: ")
	fmt.Println(ptr)
	// z := (*[1 << 30]C.uint)(ptr)
	// fmt.Println(z[0:1]

	return slave
}
