package bubble

import (
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/timing"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak"

	"image/color"
	"math"
	"sync"
	"time"
)

type Player struct {
	*entities.Moving
	*Threads
	fallspeed float64
}

func (p *Player) Init() event.CID {
	return event.NextID(p)
}

func NewPlayer() *Player {
	p := new(Player)
	p.Moving = entities.NewMoving(100, 100, 10, 12,
		charR,
		nil, p.Init(), 0) 

	render.Draw(p.R, 0, 1)

	p.Speed = physics.NewVector(.8, 4)

	p.fallspeed = .1

	thdsPos := physics.NewVector(0, 0).Attach(p.Point.Vector, 5)

	p.Threads = NewThreads(
		//color.RGBA{255, 127, 237, 255},
		color.RGBA{191,95, 178, 191},
		8,
		p.Delta,
		thdsPos,
		2,
		3,
	)
	
	p.Threads.StaticPoints = 4
	render.Draw(p.Threads, 0)

	p.Bind(func(id int, nothing interface{}) int {
		p := event.GetEntity(id).(*Player)

		// Move left and right with A and D
		if oak.IsDown(key.A) {
			p.Delta.SetX(-p.Speed.X())
		} else if oak.IsDown(key.D) {
			p.Delta.SetX(p.Speed.X())
		} else {
			p.Delta.SetX(0)
		}
		if oak.IsDown(key.LeftShift) || oak.IsDown(key.RightShift) {
			p.Delta.SetX(p.Delta.X() * 2)
		}
		oldX, oldY := p.GetPos()
		p.ShiftPos(p.Delta.X(), p.Delta.Y())

		aboveGround := false

		hit := collision.HitLabel(p.Space, Ground)

		// If we've moved in y value this frame and in the last frame,
		// we were below what we're trying to hit, we are still falling
		if hit != nil && p.Delta.Y() > 0 && !(oldY != p.Y() && oldY+p.H > hit.Y()) {
			// Correct our y if we started falling into the ground
			p.SetY(hit.Y() - p.H)
			// Stop falling
			p.Delta.SetY(0)
			// Jump with Space when on the ground
			if oak.IsDown(key.Spacebar) {
				p.Delta.ShiftY(-p.Speed.Y())
			}
			aboveGround = true
		} else {
			// Fall if there's no ground
			p.Delta.ShiftY(p.fallspeed)
		}

		if hit != nil {
			// If we walked into a piece of ground, move back
			xover, yover := p.Space.Overlap(hit)
			// We, perhaps unintuitively, need to check the Y overlap, not
			// the x overlap
			// if the y overlap exceeds a superficial value, that suggests
			// we're in a state like
			//
			// G = Ground, C = player
			//
			// GG C
			// GG C
			//
			// moving to the left
			if math.Abs(yover) > 1 {
				p.SetX(oldX)
				if p.Delta.Y() < 0 {
					p.Delta.SetY(0)
				}
			}

			// If we're below what we hit and we have significant xoverlap, by contrast,
			// then we're about to jump from below into the ground, and we
			// should stop the pacter.
			if !aboveGround && math.Abs(xover) > 1 {
				// We add a buffer so this doesn't retrigger immediately
				p.SetY(oldY + 1)
				p.Delta.SetY(p.fallspeed)
			}

		}

		p.UpdateAnim()

		 
		p.Threads.Points[0] =	floatgeom.Point2{p.Moving.X(), p.Moving.Y()}
		p.Threads.Points[1] =	floatgeom.Point2{p.Moving.X(), p.Moving.Y()+1}
		p.Threads.Points[2] =	floatgeom.Point2{p.Moving.X()+9, p.Moving.Y()+1}
		p.Threads.Points[3] =	floatgeom.Point2{p.Moving.X()+9, p.Moving.Y()}
		p.Threads.BaseTargets[7] = p.Threads.Points[0].Add(floatgeom.Point2{0,5})
		p.Threads.BaseTargets[6] = p.Threads.Points[1].Add(floatgeom.Point2{0,6})
		p.Threads.BaseTargets[5] = p.Threads.Points[2].Add(floatgeom.Point2{0,8})
		p.Threads.BaseTargets[4] = p.Threads.Points[3].Add(floatgeom.Point2{0,5})

		return 0
	}, event.Enter)

	hasTarget := false
	targetLock := sync.Mutex{}
	event.GlobalBind(func(_ int, me interface{}) int {
		mevent := me.(mouse.Event)
		targetLock.Lock()
		defer targetLock.Unlock()			
		if hasTarget {
			return 0
		}
		p.Threads.SetTargets(floatgeom.Point2{mevent.X(), mevent.Y()})
		hasTarget = true
		timing.DoAfter(600 * time.Millisecond, func(){
			p.Threads.ResetTargets()
			hasTarget = false
		})
		return 0
	}, mouse.Press)

	
	return p
}

func (p *Player) UpdateAnim() {
	sw := p.R.(*render.Switch)
	if p.Delta.X() < 0 {
		if math.Abs(p.Delta.Y()) < .2 {
			if p.Delta.X() < -1.5 {
				sw.Set("leftRun")
			} else {
				sw.Set("leftWalk")
			}
		} else {
			if p.Delta.Y() > 0 {
				sw.Set("leftJumpDown")
			} else {
				sw.Set("leftJumpUp")
			}
		}
	} else if p.Delta.X() > 0 {
		if math.Abs(p.Delta.Y()) < .2 {
			if p.Delta.X() > 1.4 {
				sw.Set("rightRun")
			} else {
				sw.Set("rightWalk")
			}
			sw.Set("rightWalk")
		} else {
			if p.Delta.Y() > 0 {
				sw.Set("rightJumpDown")
			} else {
				sw.Set("rightJumpUp")
			}
		}
	} else {
		cur := sw.Get()
		if cur[0] == 'l' {
			if math.Abs(p.Delta.Y()) < .2 {
				sw.Set("leftStand")
			} else if p.Delta.Y() > 0 {
				sw.Set("leftJumpDown")
			} else {
				sw.Set("leftJumpUp")
			}
		} else {
			if math.Abs(p.Delta.Y()) < .2 {
				sw.Set("rightStand")
			} else if p.Delta.Y() > 0 {
				sw.Set("rightJumpDown")
			} else {
				sw.Set("rightJumpUp")
			}
		}
	}
}