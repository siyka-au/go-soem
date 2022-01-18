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
	master, err := NewSOEMMaster("eno1")
	if err != nil {
		fmt.Println(err)
		return
	}
	master.ConfigInit()
	fmt.Printf("Found %d attached slaves\n", master.SlaveCount)

	master.ConfigMap(1024)

	slaves, _ := master.GetSlaves()
	for i, slave := range slaves {
		fmt.Printf(
			"Slave %d Name %s\n"+
				"  Vendor ID 0x%08x\n"+
				"  Product Code 0x%08x\n"+
				"  Revision 0x%08x\n"+
				"  Configured Address 0x%04x\n"+
				"  Alias Address 0x%04x\n"+
				"  Input Bits %d\n"+
				"  Input Bytes %d\n"+
				"  Output Bits %d\n"+
				"  Output Bytes %d\n"+
				"  Configured Address 0x\n",
			i, slave.Name, slave.VendorID, slave.ProductCode, slave.Revision,
			slave.ConfiguredAddress, slave.AliasAddress,
			slave.InputBits, slave.InputBytes,
			slave.OutputBits, slave.OutputBytes)
	}

	master.SendProcessData()
	fmt.Printf("WKC: %d\n", master.ReceiveProcessData(EC_TIMEOUTRET))
	fmt.Printf("??: %d\n", master.WriteState(0))
	fmt.Printf("State: %d\n", master.CheckState(0, EC_STATE_OPERATIONAL, EC_TIMEOUTSTATE))
}

const (
	/** return value no frame returned */
	EC_NOFRAME = -1
	/** return value unknown frame received */
	EC_OTHERFRAME = -2
	/** return value general error */
	EC_ERROR = -3
	/** return value too many slaves */
	EC_SLAVECOUNTEXCEEDED = -4
	/** return value request timeout */
	EC_TIMEOUT = -5
	/** maximum EtherCAT frame length in bytes */
	EC_MAXECATFRAME = 1518
	/** maximum EtherCAT LRW frame length in bytes */
	/* MTU - Ethernet header - length - datagram header - WCK - FCS */
	EC_MAXLRWDATA = (EC_MAXECATFRAME - 14 - 2 - 10 - 2 - 4)
	/** size of DC datagram used in first LRW frame */
	EC_FIRSTDCDATAGRAM = 20
	/** standard frame buffer size in bytes */
	EC_BUFSIZE = EC_MAXECATFRAME
	/** datagram type EtherCAT */
	EC_ECATTYPE = 0x1000
	/** number of frame buffers per channel (tx, rx1 rx2) */
	EC_MAXBUF = 16
	/** timeout value in us for tx frame to return to rx */
	EC_TIMEOUTRET = 2000
	/** timeout value in us for safe data transfer, max. triple retry */
	EC_TIMEOUTRET3 = (EC_TIMEOUTRET * 3)
	/** timeout value in us for return "safe" variant (f.e. wireless) */
	EC_TIMEOUTSAFE = 20000
	/** timeout value in us for EEPROM access */
	EC_TIMEOUTEEP = 20000
	/** timeout value in us for tx mailbox cycle */
	EC_TIMEOUTTXM = 20000
	/** timeout value in us for rx mailbox cycle */
	EC_TIMEOUTRXM = 700000
	/** timeout value in us for check statechange */
	EC_TIMEOUTSTATE = 2000000
	/** size of EEPROM bitmap cache */
	EC_MAXEEPBITMAP = 128
	/** size of EEPROM cache buffer */
	EC_MAXEEPBUF = EC_MAXEEPBITMAP << 5
	/** default number of retries if wkc <= 0 */
	EC_DEFAULTRETRIES = 3
	/** default group size in 2^x */
	EC_LOGGROUPOFFSET = 16
)

type EtherCATError uint16

const (
	/** No error */
	EC_ERR_OK EtherCATError = iota
	/** Library already initialized. */
	EC_ERR_ALREADY_INITIALIZED
	/** Library not initialized. */
	EC_ERR_NOT_INITIALIZED
	/** Timeout occurred during execution of the function. */
	EC_ERR_TIMEOUT
	/** No slaves were found. */
	EC_ERR_NO_SLAVES
	/** Function failed. */
	EC_ERR_NOK
)

type EtherCATErrorType uint16

const (
	EC_ERR_TYPE_SDO_ERROR           EtherCATErrorType = 0
	EC_ERR_TYPE_EMERGENCY           EtherCATErrorType = 1
	EC_ERR_TYPE_PACKET_ERROR        EtherCATErrorType = 3
	EC_ERR_TYPE_SDOINFO_ERROR       EtherCATErrorType = 4
	EC_ERR_TYPE_FOE_ERROR           EtherCATErrorType = 5
	EC_ERR_TYPE_FOE_BUF2SMALL       EtherCATErrorType = 6
	EC_ERR_TYPE_FOE_PACKETNUMBER    EtherCATErrorType = 7
	EC_ERR_TYPE_SOE_ERROR           EtherCATErrorType = 8
	EC_ERR_TYPE_MBX_ERROR           EtherCATErrorType = 9
	EC_ERR_TYPE_FOE_FILE_NOTFOUND   EtherCATErrorType = 10
	EC_ERR_TYPE_EOE_INVALID_RX_DATA EtherCATErrorType = 11
)

type EtherCATCommandType uint16

const (
	/** No operation */
	EC_CMD_NOP EtherCATCommandType = iota
	/** Auto Increment Read */
	EC_CMD_APRD
	/** Auto Increment Write */
	EC_CMD_APWR
	/** Auto Increment Read Write */
	EC_CMD_APRW
	/** Configured Address Read */
	EC_CMD_FPRD
	/** Configured Address Write */
	EC_CMD_FPWR
	/** Configured Address Read Write */
	EC_CMD_FPRW
	/** Broadcast Read */
	EC_CMD_BRD
	/** Broadcast Write */
	EC_CMD_BWR
	/** Broadcast Read Write */
	EC_CMD_BRW
	/** Logical Memory Read */
	EC_CMD_LRD
	/** Logical Memory Write */
	EC_CMD_LWR
	/** Logical Memory Read Write */
	EC_CMD_LRW
	/** Auto Increment Read Multiple Write */
	EC_CMD_ARMW
	/** Configured Read Multiple Write */
	EC_CMD_FRMW
	/** Reserved */
)

type EtherCATEEPROMCommandType uint16

const (
	EC_ECMD_NOP    EtherCATEEPROMCommandType = 0x0000
	EC_ECMD_READ   EtherCATEEPROMCommandType = 0x0100
	EC_ECMD_WRITE  EtherCATEEPROMCommandType = 0x0201
	EC_ECMD_RELOAD EtherCATEEPROMCommandType = 0x0300
)

type EtherCATDataType uint16

const (
	ECT_BOOLEAN         EtherCATDataType = 0x0001
	ECT_INTEGER8        EtherCATDataType = 0x0002
	ECT_INTEGER16       EtherCATDataType = 0x0003
	ECT_INTEGER32       EtherCATDataType = 0x0004
	ECT_UNSIGNED8       EtherCATDataType = 0x0005
	ECT_UNSIGNED16      EtherCATDataType = 0x0006
	ECT_UNSIGNED32      EtherCATDataType = 0x0007
	ECT_REAL32          EtherCATDataType = 0x0008
	ECT_VISIBLE_STRING  EtherCATDataType = 0x0009
	ECT_OCTET_STRING    EtherCATDataType = 0x000A
	ECT_UNICODE_STRING  EtherCATDataType = 0x000B
	ECT_TIME_OF_DAY     EtherCATDataType = 0x000C
	ECT_TIME_DIFFERENCE EtherCATDataType = 0x000D
	ECT_DOMAIN          EtherCATDataType = 0x000F
	ECT_INTEGER24       EtherCATDataType = 0x0010
	ECT_REAL64          EtherCATDataType = 0x0011
	ECT_INTEGER64       EtherCATDataType = 0x0015
	ECT_UNSIGNED24      EtherCATDataType = 0x0016
	ECT_UNSIGNED64      EtherCATDataType = 0x001B
	ECT_BIT1            EtherCATDataType = 0x0030
	ECT_BIT2            EtherCATDataType = 0x0031
	ECT_BIT3            EtherCATDataType = 0x0032
	ECT_BIT4            EtherCATDataType = 0x0033
	ECT_BIT5            EtherCATDataType = 0x0034
	ECT_BIT6            EtherCATDataType = 0x0035
	ECT_BIT7            EtherCATDataType = 0x0036
	ECT_BIT8            EtherCATDataType = 0x0037
)

type EtherCATState uint8

const (
	EC_STATE_NONE        EtherCATState = 0x00
	EC_STATE_INIT        EtherCATState = 0x01
	EC_STATE_PRE_OP      EtherCATState = 0x02
	EC_STATE_BOOT        EtherCATState = 0x03
	EC_STATE_SAFE_OP     EtherCATState = 0x04
	EC_STATE_OPERATIONAL EtherCATState = 0x08
	EC_STATE_ACK         EtherCATState = 0x10
	EC_STATE_ERROR       EtherCATState = 0x10
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
	// input bits
	InputBits uint16
	// input bytes, if Ibits < 8 then Ibytes = 0
	InputBytes uint32
	// output bits
	OutputBits uint16
	// output bytes, if Obits < 8 then Obytes = 0
	OutputBytes uint32
}

type SOEM struct {
	SlaveCount uint

	context   C.ecx_contextt
	ioMap     unsafe.Pointer
	ioMapSize C.int
}

func NewSOEMMaster(ifname string) (*SOEM, error) {
	soem := new(SOEM)
	cifname := C.CString(ifname)
	defer C.free(unsafe.Pointer(cifname))

	soem.context = C.ecx_context // no idea why this is needed or even works

	if int(C.ecx_init(&soem.context, cifname)) <= 0 {
		return nil, fmt.Errorf("error opening interface %s", ifname)
	}

	return soem, nil
}

func (m *SOEM) ConfigInit() {
	m.SlaveCount = uint(C.ecx_config_init(&m.context, 0))
}

func (m *SOEM) ConfigMapWithGroup(group uint8, size uint) {
	m.ioMap = C.malloc(C.size_t(size))
	m.ioMapSize = C.ecx_config_map_group(&m.context, m.ioMap, C.uchar(group))
}

func (m *SOEM) ConfigMap(size uint) {
	m.ConfigMapWithGroup(0, size)
}

func (m *SOEM) GetSlave(index uint) (*Slave, error) {
	if index >= m.SlaveCount {
		return nil, errors.New("slave index out of range")
	}
	return marshalSlave(C.ec_slave[index+1]), nil
}

func (m *SOEM) GetSlaves() ([]*Slave, error) {
	slaves := make([]*Slave, m.SlaveCount)
	for i := uint(0); i < m.SlaveCount; i++ {
		slave, err := m.GetSlave(i)
		if err != nil {
			return nil, err
		}
		slaves[i] = slave
	}
	return slaves, nil
}

func (m *SOEM) GetIOMap() []byte {
	return C.GoBytes(m.ioMap, m.ioMapSize)
}

func (m *SOEM) ReadState() int {
	return int(C.ecx_readstate(&m.context))
}

func (m *SOEM) WriteState(slave uint16) int {
	return int(C.ecx_writestate(&m.context, C.ushort(slave)))
}

func (m *SOEM) CheckState(slave uint16, requestedState EtherCATState, timeout int) int {
	return int(C.ecx_statecheck(&m.context,
		C.ushort(slave),
		C.ushort(requestedState),
		C.int(timeout)))
}

func (m *SOEM) SendProcessDataWithGroup(group uint8) {
	C.ecx_send_processdata(&m.context)
}

func (m *SOEM) SendProcessData() {
	m.SendProcessDataWithGroup(0)
}

func (m *SOEM) ReceiveProcessDataWithGroup(group uint8, timeout int) int {
	return int(C.ecx_receive_processdata(&m.context, C.int(timeout)))
}

func (m *SOEM) ReceiveProcessData(timeout int) int {
	return m.ReceiveProcessDataWithGroup(0, timeout)
}

// func (m *SOEM) RecoverSlave(uint16 slave, int timeout) error {
// 	if C.ecx_recover_slave(&m.context) <= 0 {
// 		return errors.New("error recovering slave")
// 	}
// 	return nil
// }

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
	slave.InputBits = uint16(cslave.Ibits)
	slave.InputBytes = uint32(cslave.Ibytes)
	slave.OutputBits = uint16(cslave.Obits)
	slave.OutputBytes = uint32(cslave.Obytes)

	return slave
}
