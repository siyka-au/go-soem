package plc

import (
	"context"
	"fmt"
	"time"

	"github.com/qmuntal/stateless"
)

type MultiClick struct {
	fsm               *stateless.StateMachine
	clickCount        uint8
	clickTimeout      time.Duration
	currentClickCount uint8
	cancel            chan bool
	Clicks            chan bool
}

type multiClickState uint8

const (
	multiClickStateInit multiClickState = iota
	multiClickStateIdle
	multiClickStatePartial
	multiClickStateClicked
)

type multiClickTrigger uint8

const (
	multiClickTriggerInit multiClickTrigger = iota
	multiClickTriggerClick
	multiClickTriggerTimeout
	multiClickTriggerClickCountReached
	multiClickTriggerReset
)

func NewMultiClick(clickCount uint8, clickTimeout time.Duration) *MultiClick {

	fsm := stateless.NewStateMachine(multiClickStateInit)
	mc := MultiClick{
		fsm:               fsm,
		clickCount:        clickCount,
		clickTimeout:      clickTimeout,
		currentClickCount: 0,
		cancel:            make(chan bool),
		Clicks:            make(chan bool),
	}

	fsm.Configure(multiClickStateInit).
		Permit(multiClickTriggerInit, multiClickStateIdle)

	fsm.Configure(multiClickStateIdle).
		OnEntry(func(_ context.Context, _ ...interface{}) error {
			mc.reset()
			return nil
		}).
		Permit(multiClickTriggerClick, multiClickStatePartial)

	fsm.Configure(multiClickStatePartial).
		OnEntry(func(_ context.Context, _ ...interface{}) error {
			mc.currentClickCount++
			go func() {
				timer := time.NewTimer(mc.clickTimeout)
				select {
				case <-timer.C:
					fmt.Println("Timed out")
					fsm.Fire(multiClickTriggerTimeout)
				case <-mc.cancel:
					if !timer.Stop() {
						<-timer.C
					}
				}
			}()
			return nil
		}).
		OnExit(func(_ context.Context, _ ...interface{}) error {
			select {
			case mc.cancel <- true:
			default:
			}
			return nil
		}).
		PermitReentry(multiClickTriggerClick, func(_ context.Context, _ ...interface{}) bool {
			return !mc.clickCountReached()
		}).
		Permit(multiClickTriggerClick, multiClickStateClicked, func(_ context.Context, _ ...interface{}) bool {
			return mc.clickCountReached()
		}).
		Permit(multiClickTriggerTimeout, multiClickStateIdle)

	fsm.Configure(multiClickStateClicked).
		OnEntry(func(_ context.Context, _ ...interface{}) error {
			fmt.Println("Clicked")
			select {
			case mc.Clicks <- true:
			default:
			}
			fsm.Fire(multiClickTriggerReset)
			return nil
		}).
		Permit(multiClickTriggerReset, multiClickStateIdle)

	fsm.Fire(multiClickTriggerInit)
	return &mc
}

func (mc *MultiClick) Click() {
	mc.fsm.Fire(multiClickTriggerClick)
}

func (mc *MultiClick) reset() {
	mc.currentClickCount = 0
}

func (mc *MultiClick) clickCountReached() bool {
	return mc.currentClickCount >= mc.clickCount
}
