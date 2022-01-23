package plc

import (
	"context"
	"fmt"
	"time"

	"github.com/qmuntal/stateless"
)

type RisingEdge struct {
	currentState bool
	lastState    bool
}

func NewRisingEdge() RisingEdge {
	r := RisingEdge{false, false}
	return r
}

func (r *RisingEdge) Run(state bool) bool {
	r.lastState = r.currentState
	r.currentState = state
	return !r.lastState && r.currentState
}

type FallingEdge struct {
	currentState bool
	lastState    bool
}

func NewFallingEdge() RisingEdge {
	r := RisingEdge{false, false}
	return r
}

func (r *FallingEdge) Run(state bool) bool {
	r.lastState = r.currentState
	r.currentState = state
	return r.lastState && !r.currentState
}

type AutoManual struct {
	fsm            *stateless.StateMachine
	manualDuration time.Duration
	trig           chan autoManualTrigger
}

type autoManualState uint8

const (
	stateInit autoManualState = iota
	stateAutomatic
	stateManual
)

type autoManualTrigger uint8

const (
	triggerInit autoManualTrigger = iota
	triggerEnterManual
	triggerManualTimedOut
	triggerCancelManual
)

func NewAutoManual() AutoManual {

	fsm := stateless.NewStateMachine(stateInit)
	trig := make(chan autoManualTrigger, 5)
	am := AutoManual{fsm, 10 * time.Second, trig}

	fsm.Configure(stateInit).
		Permit(triggerInit, stateAutomatic)

	fsm.Configure(stateAutomatic).
		OnEntry(func(_ context.Context, _ ...interface{}) error {
			fmt.Println("Automatic Mode")
			return nil
		}).
		Permit(triggerEnterManual, stateManual)

	fsm.Configure(stateManual).
		OnEntry(func(_ context.Context, _ ...interface{}) error {
			fmt.Println("Manual Mode")
			go func() {
				timer1 := time.NewTimer(am.manualDuration)
				<-timer1.C
				fsm.Fire(triggerManualTimedOut)
			}()
			return nil
		}).
		Permit(triggerManualTimedOut, stateAutomatic).
		Permit(triggerCancelManual, stateAutomatic)

	fsm.Fire(triggerInit)
	return am
}

func (am *AutoManual) StartManual(t time.Duration) {
	am.fsm.Fire(triggerEnterManual)
}

func (am *AutoManual) CancelManual() {
	am.fsm.Fire(triggerCancelManual)
}

func (am *AutoManual) IsAutomatic() bool {
	return true
}

func (am *AutoManual) IsManual() bool {
	return false
}
