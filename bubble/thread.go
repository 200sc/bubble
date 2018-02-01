package bubble

import (
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"

	"image/color"
	"image/draw"
)

type Threads struct {
	render.LayeredPoint
	points   []floatgeom.Point2
	velocity physics.Vector
	Color    color.Color
}

func NewThreads(c color.Color, size int, vel, anchor physics.Vector, layer int) *Threads {
	return &Threads{
		LayeredPoint: render.LayeredPoint{
			Vector: anchor,
			Layer:  render.NewLayer(layer),
		},
		Color:    c,
		velocity: vel,
		points:   make([]floatgeom.Point2, size),
	}
}

func (ts *Threads) Draw(buff draw.Image) {
	ts.DrawOffset(buff, 0, 0)
}
func (ts *Threads) DrawOffset(buff draw.Image, xOff, yOff float64) {
	ts.update()
	for _, p := range ts.points {
		buff.Set(int(p.X()+ts.X()), int(p.Y()+ts.Y()), ts.Color)
	}
}

func (ts *Threads) update() {
	// move pixels in the direction of velocity
	// occupied := make(map[floatgeom.Point2]struct{})
	// for len(occupied) < len(ts.points) {

	// }
}
