package bubble

type State int

const (
	Grounded State = iota
	InAir    State = iota
	Swinging State = iota
)

var (
	canWalk = map[State]bool{
		Grounded: true,
		InAir:    true,
		Swinging: false,
	}
	canJump = map[State]bool{
		Grounded: true,
		InAir:    false,
		Swinging: false,
	}
	canSwing = map[State]bool{
		Grounded: false,
		InAir:    true,
		Swinging: true,
	}
)

func (s State) CanWalk() bool {
	can, ok := canWalk[s]
	if !ok {
		return true
	}
	return can
}

func (s State) CanJump() bool {
	can, ok := canJump[s]
	if !ok {
		return false
	}
	return can
}

func (s State) CanSwing() bool {
	can, ok := canSwing[s]
	if !ok {
		return true
	}
	return can
}
