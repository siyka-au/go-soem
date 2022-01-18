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

type Context struct {
}
type Slave struct {
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
