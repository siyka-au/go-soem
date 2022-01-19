package soem

/*
#cgo LDFLAGS: -lsoem

#include <stdio.h>
#include <stdlib.h>
#include <soem/ethercat.h>

*/
import "C"
import (
	"fmt"
)

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
	EC_STATE_ACK         EtherCATState = 0x10 // TODO This is baffling; why the heck do we have duplicate elements?
	EC_STATE_ERROR       EtherCATState = 0x10 // TODO This is baffling; why the heck do we have duplicate elements?
)

func (e EtherCATState) String() string {
	switch e {
	case EC_STATE_NONE:
		return "EC_STATE_NONE"
	case EC_STATE_INIT:
		return "EC_STATE_INIT"
	case EC_STATE_PRE_OP:
		return "EC_STATE_PRE_OP"
	case EC_STATE_BOOT:
		return "EC_STATE_BOOT"
	case EC_STATE_SAFE_OP:
		return "EC_STATE_SAFE_OP"
	case EC_STATE_OPERATIONAL:
		return "EC_STATE_OPERATIONAL"
	// case EC_STATE_ACK:
	// 	return "EC_STATE_ACK"
	case EC_STATE_ERROR:
		return "EC_STATE_ERROR"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}
