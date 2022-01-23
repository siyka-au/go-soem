package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/qmuntal/stateless"
)

type ControllerTrigger string

const (
	triggerInit           ControllerTrigger = "Initialise"
	triggerStart          ControllerTrigger = "Start"
	triggerStartCompleted ControllerTrigger = "Start Completed"
	triggerStop           ControllerTrigger = "Stop"
	triggerStopCompleted  ControllerTrigger = "Stop Completed"
	triggerAbort          ControllerTrigger = "Abort"
)

type ControllerState string

const (
	stateInit     ControllerState = "Initialising"
	stateIdle     ControllerState = "Idle"
	stateStarting ControllerState = "Starting"
	stateRunning  ControllerState = "Running"
	stateStopping ControllerState = "Stopping"
	stateAborted  ControllerState = "Aborted"
)

type Controller struct {
	fsm  *stateless.StateMachine
	trig chan string
}

func NewController() Controller {

	fsm := stateless.NewStateMachine(stateInit)
	trig := make(chan string, 5)
	ctrl := Controller{fsm, trig}

	fsm.Configure(stateInit).
		Permit(triggerInit, stateIdle)

	fsm.Configure(stateIdle).
		OnEntry(func(_ context.Context, _ ...interface{}) error {
			fmt.Println("Idle")
			go func() {
				<-ctrl.trig
				fsm.Fire(triggerStart)
			}()
			return nil
		}).
		Permit(triggerStart, stateStarting)

	fsm.Configure(stateStarting).
		OnEntry(func(_ context.Context, _ ...interface{}) error {
			fmt.Println("Starting")
			go func() {
				timer1 := time.NewTimer(2 * time.Second)
				<-timer1.C
				fsm.Fire(triggerStartCompleted)
			}()
			return nil
		}).
		Permit(triggerStartCompleted, stateRunning)

	fsm.Configure(stateRunning).
		OnEntry(func(_ context.Context, _ ...interface{}) error {
			ctrl.startRunning()
			go ctrl.running()
			return nil
		}).
		OnExit(func(_ context.Context, _ ...interface{}) error {
			ctrl.stopRunning()
			return nil
		}).
		Permit(triggerStop, stateStopping)

	fsm.Configure(stateStopping).
		OnEntry(func(_ context.Context, _ ...interface{}) error {
			fmt.Println("Stopping")
			go func() {
				timer1 := time.NewTimer(2 * time.Second)
				<-timer1.C
				fsm.Fire(triggerStopCompleted)
			}()
			return nil
		}).
		Permit(triggerStopCompleted, stateIdle)

	fsm.Fire(triggerInit)
	return ctrl
}

func (c *Controller) Bump() {
	c.trig <- "bump"
}

func (c *Controller) startRunning() {
	fmt.Println("Start running")

}

func (c *Controller) running() {
	fmt.Println("Running")
	<-c.trig
	fmt.Println("Received trigger")
	c.fsm.Fire(triggerStop)
}

func (c *Controller) stopRunning() {
	fmt.Println("Stop running")
}
