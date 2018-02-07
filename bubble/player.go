package bubble

import (
	"fmt"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/timing"

	"image/color"
	"math"
	"sync"
	"time"
)

type Player struct {
	*entities.Moving
	*Threads
	fallspeed float64
	MaxSpeed  floatgeom.Point2
	State
	Facing
	hasTarget  bool
	Target     floatgeom.Point2
	SwingAccel floatgeom.Point2
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

	p.Speed = physics.NewVector(.15, 4)
	p.MaxSpeed = floatgeom.Point2{4, 15}
	p.fallspeed = .1

	thdsPos := physics.NewVector(0, 0).Attach(p.Point.Vector, 5)

	p.Threads = NewThreads(
		//color.RGBA{255, 127, 237, 255},
		color.RGBA{191, 95, 178, 191},
		9,
		p.Delta,
		thdsPos,
		2,
		3,
	)
	p.Threads.Swinger = p.CID

	p.Threads.StaticPoints = 4
	render.Draw(p.Threads, 0, 3)

	p.Bind(func(id int, nothing interface{}) int {
		p := event.GetEntity(id).(*Player)

		// Move left and right with A and D
		speed := p.Speed
		if p.CanWalk() {
			if p.State == Grounded {
				if oak.IsDown(key.LeftShift) || oak.IsDown(key.RightShift) {
					speed = p.Speed.Copy().Scale(2)
				}
			} else if p.State == InAir {
				speed = p.Speed.Copy().Scale(.7)
			}

			if oak.IsDown(key.A) {
				p.Delta.ShiftX(-speed.X())
			} else if oak.IsDown(key.D) {
				p.Delta.ShiftX(speed.X())
			} else {
				p.Delta.SetX(p.Delta.X() * .8)
			}
		}
		if p.State == Swinging {
			pos := floatgeom.Point2{p.X(), p.Y()}
			magn := pos.Sub(p.Target).Magnitude()
			angle := pos.RadiansTo(p.Target)
			for angle > math.Pi*2 {
				angle -= math.Pi * 2
			}
			c := math.Cos(angle) * -1
			s := math.Sin(angle)
			p.SwingAccel = p.SwingAccel.Add(floatgeom.Point2{c, s})
			fmt.Println(magn, angle, p.SwingAccel)

			p.Delta.ShiftX(p.SwingAccel.X() * p.fallspeed * p.fallspeed)
			p.Delta.ShiftY(p.SwingAccel.Y() * p.fallspeed * p.fallspeed)

			dMgn := p.Delta.Magnitude()
			circum := magn * math.Pi * 2
			arcPercentage := dMgn / circum
			arcAngle := arcPercentage * 2 * math.Pi
			if c > 0 {
				arcAngle *= -1
			}
			goalAngle := angle + arcAngle
			goalPos := p.Target.Add(floatgeom.RadianPoint(goalAngle).MulConst(magn))
			p.Delta.SetX(goalPos.X() - p.X())
			p.Delta.SetY(goalPos.Y() - p.Y())

			// correct delta to move dMgn but end up on circle somewhere
			//
			//   T
			//
			//        C
			//
			//
			//

			if oak.IsDown(key.Spacebar) {
				p.State = InAir
				// push off?
			}
		}
		if p.State == Grounded || p.State == Swinging {
			p.Delta.SetX(p.Delta.X() * .95)
			if math.Abs(p.Delta.X()) < .05 {
				p.Delta.SetX(0)
			}
		}

		if math.Abs(p.Delta.X()) > p.MaxSpeed.X() {
			if p.Delta.X() < 0 {
				p.Delta.SetX(-p.MaxSpeed.X())
			} else {
				p.Delta.SetX(p.MaxSpeed.X())
			}
		}

		// Jump with Space
		if p.CanJump() {
			if oak.IsDown(key.Spacebar) {
				p.Delta.ShiftY(-p.Speed.Y())
			}
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
			p.State = Grounded
			aboveGround = true
		} else {
			if hit == nil && math.Abs(p.Delta.Y()) > 0 && p.State != Swinging {
				p.State = InAir
			}
			// Fall if there's no ground
			if p.State != Swinging {
				p.Delta.ShiftY(p.fallspeed)
			}
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

		p.Threads.Points[0] = floatgeom.Point2{p.Moving.X(), p.Moving.Y()}
		p.Threads.Points[1] = floatgeom.Point2{p.Moving.X(), p.Moving.Y() + 1}
		p.Threads.Points[2] = floatgeom.Point2{p.Moving.X() + 9, p.Moving.Y() + 1}
		p.Threads.Points[3] = floatgeom.Point2{p.Moving.X() + 9, p.Moving.Y()}
		if p.Facing == Left {
			p.Threads.BaseTargets[8] = p.Threads.Points[0].Add(floatgeom.Point2{0, 1})
			p.Threads.BaseTargets[7] = p.Threads.Points[1].Add(floatgeom.Point2{3, 3})
			p.Threads.BaseTargets[6] = p.Threads.Points[2].Add(floatgeom.Point2{-3, 5})
			p.Threads.BaseTargets[5] = p.Threads.Points[3].Add(floatgeom.Point2{0, 8})
			p.Threads.BaseTargets[4] = p.Threads.Points[3].Add(floatgeom.Point2{2, 12})
		} else {
			p.Threads.BaseTargets[8] = p.Threads.Points[0].Add(floatgeom.Point2{-2, 12})
			p.Threads.BaseTargets[7] = p.Threads.Points[0].Add(floatgeom.Point2{0, 8})
			p.Threads.BaseTargets[6] = p.Threads.Points[1].Add(floatgeom.Point2{3, 5})
			p.Threads.BaseTargets[5] = p.Threads.Points[2].Add(floatgeom.Point2{-3, 3})
			p.Threads.BaseTargets[4] = p.Threads.Points[3].Add(floatgeom.Point2{0, 1})
		}

		// remove this
		if p.Y() > float64(oak.ScreenHeight) {
			p.Moving.ShiftY(float64(-oak.ScreenHeight))
		}
		return 0
	}, event.Enter)

	targetLock := sync.Mutex{}
	p.Bind(func(id int, me interface{}) int {
		p := event.CID(id).E().(*Player)
		if p.CanSwing() {
			fmt.Println("p can swing", p.State)
			mevent := me.(mouse.Event)
			targetLock.Lock()
			defer targetLock.Unlock()
			if p.hasTarget {
				return 0
			}
			p.Threads.SetTargets(floatgeom.Point2{mevent.X(), mevent.Y()})
			p.hasTarget = true
			go timing.DoAfter(600*time.Millisecond, func() {
				p.Threads.ResetTargets()
				p.hasTarget = false
			})
		}
		return 0
	}, mouse.Press)

	p.Bind(func(id int, hit interface{}) int {
		if p.hasTarget {
			cp := hit.(*collision.Space)
			fmt.Println("Threads hit", hit)
			x, y := cp.GetCenter()
			p.Target = floatgeom.Point2{x, y}
			p.State = Swinging
		}
		return 0
	}, "SwingHit")

	return p
}

func (p *Player) UpdateAnim() {
	sw := p.R.(*render.Switch)
	if p.Delta.X() < 0 {
		p.Facing = Left
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
		p.Facing = Right
		if math.Abs(p.Delta.Y()) < .2 {
			if p.Delta.X() > 1.4 {
				sw.Set("rightRun")
			} else {
				sw.Set("rightWalk")
			}
		} else {
			if p.Delta.Y() > 0 {
				sw.Set("rightJumpDown")
			} else {
				sw.Set("rightJumpUp")
			}
		}
	} else {
		if p.Facing == Left {
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
