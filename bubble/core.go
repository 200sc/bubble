package bubble

import (
	"image"
	"image/color"
	"math"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/shape"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

func addMainScene() {
	oak.Add("bubble", func(string, interface{}) {

		snow := particle.And(
			particle.Angle(floatrange.NewLinear(-20, 20)),
			particle.Layer(func(physics.Vector) int { return 2 }),
			particle.Pos(330, 120),
			particle.Spread(10, 130),
			particle.Size(intrange.NewLinear(1, 3)),
			particle.Color(
				color.RGBA{240, 240, 240, 255},
				color.RGBA{15, 15, 15, 0},
				color.RGBA{240, 240, 240, 255},
				color.RGBA{15, 15, 15, 0},
			),
			particle.Duration(particle.Inf),
			particle.LifeSpan(floatrange.NewLinear(100, 101)),
			particle.Speed(floatrange.NewLinear(3, 8)),
			particle.NewPerFrame(floatrange.NewLinear(2, 10)),
			particle.Shape(shape.Square),
		)

		backSnow := particle.And(
			snow,
			particle.Layer(func(physics.Vector) int { return 0 }),
			particle.Color(
				color.RGBA{140, 140, 140, 255},
				color.RGBA{15, 15, 15, 0},
				color.RGBA{140, 140, 140, 255},
				color.RGBA{15, 15, 15, 0},
			),
			particle.Speed(floatrange.NewLinear(2, 5)),
			particle.NewPerFrame(floatrange.NewLinear(1, 5)),
		)

		particle.NewColorGenerator(snow).Generate(0)
		particle.NewColorGenerator(backSnow).Generate(0)

		oak.Background = image.NewUniform(color.RGBA{100, 100, 229, 255})

		char := entities.NewMoving(100, 100, 10, 12,
			charR,
			nil, 0, 0)

		render.Draw(char.R, 0, 1)

		char.Speed = physics.NewVector(.8, 4)

		fallSpeed := .1

		thdsPos := physics.NewVector(0, 0).Attach(char.Point.Vector, 5)

		thds := NewThreads(
			color.RGBA{255, 0, 0, 255},
			20,
			char.Delta,
			thdsPos,
			2,
		)

		render.Draw(thds, 0)

		char.Bind(func(id int, nothing interface{}) int {
			char := event.GetEntity(id).(*entities.Moving)

			// Move left and right with A and D
			if oak.IsDown(key.A) {
				char.Delta.SetX(-char.Speed.X())
			} else if oak.IsDown(key.D) {
				char.Delta.SetX(char.Speed.X())
			} else {
				char.Delta.SetX(0)
			}
			if oak.IsDown(key.LeftShift) || oak.IsDown(key.RightShift) {
				char.Delta.SetX(char.Delta.X() * 2)
			}
			oldX, oldY := char.GetPos()
			char.ShiftPos(char.Delta.X(), char.Delta.Y())

			aboveGround := false

			hit := collision.HitLabel(char.Space, Ground)

			// If we've moved in y value this frame and in the last frame,
			// we were below what we're trying to hit, we are still falling
			if hit != nil && char.Delta.Y() > 0 && !(oldY != char.Y() && oldY+char.H > hit.Y()) {
				// Correct our y if we started falling into the ground
				char.SetY(hit.Y() - char.H)
				// Stop falling
				char.Delta.SetY(0)
				// Jump with Space when on the ground
				if oak.IsDown(key.Spacebar) {
					char.Delta.ShiftY(-char.Speed.Y())
				}
				aboveGround = true
			} else {
				// Fall if there's no ground
				char.Delta.ShiftY(fallSpeed)
			}

			if hit != nil {
				// If we walked into a piece of ground, move back
				xover, yover := char.Space.Overlap(hit)
				// We, perhaps unintuitively, need to check the Y overlap, not
				// the x overlap
				// if the y overlap exceeds a superficial value, that suggests
				// we're in a state like
				//
				// G = Ground, C = Character
				//
				// GG C
				// GG C
				//
				// moving to the left
				if math.Abs(yover) > 1 {
					char.SetX(oldX)
					if char.Delta.Y() < 0 {
						char.Delta.SetY(0)
					}
				}

				// If we're below what we hit and we have significant xoverlap, by contrast,
				// then we're about to jump from below into the ground, and we
				// should stop the character.
				if !aboveGround && math.Abs(xover) > 1 {
					// We add a buffer so this doesn't retrigger immediately
					char.SetY(oldY + 1)
					char.Delta.SetY(fallSpeed)
				}

			}

			sw := char.R.(*render.Switch)
			if char.Delta.X() < 0 {
				if math.Abs(char.Delta.Y()) < .2 {
					if char.Delta.X() < -1.5 {
						sw.Set("leftRun")
					} else {
						sw.Set("leftWalk")
					}
				} else {
					if char.Delta.Y() > 0 {
						sw.Set("leftJumpDown")
					} else {
						sw.Set("leftJumpUp")
					}
				}
			} else if char.Delta.X() > 0 {
				if math.Abs(char.Delta.Y()) < .2 {
					if char.Delta.X() > 1.4 {
						sw.Set("rightRun")
					} else {
						sw.Set("rightWalk")
					}
					sw.Set("rightWalk")
				} else {
					if char.Delta.Y() > 0 {
						sw.Set("rightJumpDown")
					} else {
						sw.Set("rightJumpUp")
					}
				}
			} else {
				cur := sw.Get()
				if cur[0] == 'l' {
					if math.Abs(char.Delta.Y()) < .2 {
						sw.Set("leftStand")
					} else if char.Delta.Y() > 0 {
						sw.Set("leftJumpDown")
					} else {
						sw.Set("leftJumpUp")
					}
				} else {
					if math.Abs(char.Delta.Y()) < .2 {
						sw.Set("rightStand")
					} else if char.Delta.Y() > 0 {
						sw.Set("rightJumpDown")
					} else {
						sw.Set("rightJumpUp")
					}
				}
			}

			return 0
		}, event.Enter)

		platforms := []floatgeom.Rect2{
			floatgeom.NewRect2WH(0, 200, 150, 10),
			floatgeom.NewRect2WH(50, 125, 20, 10),
			floatgeom.NewRect2WH(170, 150, 50, 10),
		}

		for _, p := range platforms {
			ground := entities.NewSolid(p.Min.X(), p.Min.Y(), p.W(), p.H(),
				render.NewColorBox(int(p.W()), int(p.H()), color.RGBA{180, 180, 180, 255}),
				nil, 0)
			ground.UpdateLabel(Ground)

			render.Draw(ground.R, 0, 1)
		}

	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "bubble", nil
	})
}
