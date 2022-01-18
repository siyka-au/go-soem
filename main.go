package main

/*
#cgo LDFLAGS: -lsoem

#include <stdio.h>
#include <stdlib.h>
#include <soem/ethercat.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	master, err := NewSOEMMaster("enp6s0f0")
	if err != nil {
		return
	}
	print(master)
}

type SOEM struct {
	Name string
	Roll int
}

func NewSOEMMaster(ifname string) (*SOEM, error) {
	soem := new(SOEM)
	cifname := C.CString(ifname)
	defer C.free(unsafe.Pointer(cifname))

	if int(C.ec_init(cifname)) <= 0 {
		return nil, fmt.Errorf("error opening interface ifname = %s", ifname)
	}

	return soem, nil
}

func (m SOEM) slaveCount() int {
	return int(C.ec_slavecount())
}