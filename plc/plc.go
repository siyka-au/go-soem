package plc

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
