package bubble

import (
	"path/filepath"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/oak/scene"
)

var (
	eyes   [][]*render.Sprite
	bodies [][]*render.Sprite
	heads  [][]*render.Sprite
	shoes  [][]*render.Sprite
	charR  render.Renderable
)

func addSetupScene() {
	oak.Add("setup", func(string, interface{}) {
		eyes = render.GetSheet(filepath.Join("4x2", "eyes.png"))
		shoes = render.GetSheet(filepath.Join("8x5", "shoes.png"))
		bodies = render.GetSheet(filepath.Join("10x7", "bodies.png"))
		heads = render.GetSheet(filepath.Join("10x7", "heads.png"))

		walk := render.NewSequence(18,
			shoes[0][0].Copy(),
			shoes[1][0].Copy(),
			shoes[1][1].Copy(),
			shoes[1][2].Copy(),
			shoes[0][1].Copy(),
			shoes[0][2].Copy(),
		)

		run := walk.Copy().(*render.Sequence)
		run.SetFPS(48)

		right := render.NewCompositeM(
			bodies[0][0].Copy(),
			heads[0][0].Copy(),
			eyes[0][0].Copy(),
			shoes[0][0].Copy(),
		)

		right.SetOffsets(
			floatgeom.Point2{0, 5},
			floatgeom.Point2{0, 0},
			floatgeom.Point2{5, 3},
			floatgeom.Point2{1, 7},
		)

		left := right.Copy().Modify(mod.FlipX).(*render.CompositeM)
		left.AddOffset(2, floatgeom.Point2{1, 3})

		leftWalk := left.Copy().(*render.CompositeM)
		leftWalk.SetIndex(3, walk.Copy())
		leftWalk.AddOffset(3, floatgeom.Point2{1, 7})

		rightWalk := right.Copy().(*render.CompositeM)
		rightWalk.SetIndex(3, walk.Copy().Modify(mod.FlipX))
		rightWalk.AddOffset(3, floatgeom.Point2{1, 7})

		leftRun := left.Copy().(*render.CompositeM)
		leftRun.SetIndex(3, run.Copy())
		leftRun.AddOffset(3, floatgeom.Point2{1, 7})

		rightRun := right.Copy().(*render.CompositeM)
		rightRun.SetIndex(3, run.Copy().Modify(mod.FlipX))
		rightRun.AddOffset(3, floatgeom.Point2{1, 7})

		rightJumpUp := right.Copy().(*render.CompositeM)
		rightJumpUp.SetIndex(2, eyes[1][1].Copy())
		rightJumpUp.AddOffset(2, floatgeom.Point2{5, 3})
		rightJumpUp.AddOffset(3, floatgeom.Point2{1, 6})

		rightJumpDown := right.Copy().(*render.CompositeM)
		rightJumpDown.SetIndex(2, eyes[0][1].Copy())
		rightJumpDown.AddOffset(2, floatgeom.Point2{5, 3})
		rightJumpDown.AddOffset(3, floatgeom.Point2{1, 8})

		leftJumpUp := rightJumpUp.Copy().Modify(mod.FlipX).(*render.CompositeM)
		leftJumpUp.AddOffset(2, floatgeom.Point2{1, 3})

		leftJumpDown := rightJumpDown.Copy().Modify(mod.FlipX).(*render.CompositeM)
		leftJumpDown.AddOffset(2, floatgeom.Point2{1, 3})

		swt := render.NewSwitch(
			"rightStand",
			map[string]render.Modifiable{
				"leftStand":     left,
				"rightStand":    right,
				"leftWalk":      leftWalk,
				"rightWalk":     rightWalk,
				"leftRun":       leftRun,
				"rightRun":      rightRun,
				"leftJumpUp":    leftJumpUp,
				"leftJumpDown":  leftJumpDown,
				"rightJumpUp":   rightJumpUp,
				"rightJumpDown": rightJumpDown,
			},
		)

		charR = swt

	}, func() bool {
		return false
	}, func() (string, *scene.Result) {
		return "bubble", nil
	})
}
