package plc

import (
	"context"
	"time"

	"github.com/qmuntal/stateless"
)

type AutoManual struct {
	fsm     *stateless.StateMachine
	timeout time.Duration
	cancel  chan bool
}

type autoManualState uint8

const (
	autoManualStateInit autoManualState = iota
	autoManualStateAutomatic
	autoManualStateManual
)

type autoManualTrigger uint8

const (
	autoManualTriggerInit autoManualTrigger = iota
	autoManualTriggerEnterManual
	autoManualTriggerManualTimedOut
	autoManualTriggerCancelManual
)

func NewAutoManual() *AutoManual {

	fsm := stateless.NewStateMachine(autoManualStateInit)
	am := AutoManual{
		fsm:     fsm,
		timeout: 3 * time.Second,
		cancel:  make(chan bool),
	}

	fsm.Configure(autoManualStateInit).
		Permit(autoManualTriggerInit, autoManualStateAutomatic)

	fsm.Configure(autoManualStateAutomatic).
		Permit(autoManualTriggerEnterManual, autoManualStateManual)

	fsm.Configure(autoManualStateManual).
		OnEntry(func(_ context.Context, _ ...interface{}) error {
			go func() {
				timer := time.NewTimer(am.timeout)
				select {
				case <-timer.C:
					fsm.Fire(autoManualTriggerManualTimedOut)
				case <-am.cancel:
					if !timer.Stop() {
						<-timer.C
					}
					fsm.Fire(autoManualTriggerCancelManual)
				}
			}()
			return nil
		}).
		Permit(autoManualTriggerManualTimedOut, autoManualStateAutomatic).
		Permit(autoManualTriggerCancelManual, autoManualStateAutomatic)

	fsm.Fire(autoManualTriggerInit)
	return &am
}

func (am *AutoManual) StartManual(timeout time.Duration) {
	am.timeout = timeout
	am.fsm.Fire(autoManualTriggerEnterManual)
}

func (am *AutoManual) CancelManual() {
	select {
	case am.cancel <- true:
	default:
	}
}

func (am *AutoManual) IsAutomatic() bool {
	return am.fsm.MustState() == autoManualStateAutomatic
}

func (am *AutoManual) IsManual() bool {
	return am.fsm.MustState() == autoManualStateManual
}
