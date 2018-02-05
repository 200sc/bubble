package bubble

import (
	"image"
	"image/color"

	"github.com/oakmound/oak/render/particle"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
) 

func addMainScene() {
	oak.Add("bubble", func(string, interface{}) {

		particle.NewColorGenerator(snow).Generate(0)
		particle.NewColorGenerator(backSnow).Generate(0)

		oak.Background = image.NewUniform(color.RGBA{100, 100, 229, 255})

		NewPlayer() 

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
