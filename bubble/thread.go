package bubble

import (
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"

	"image/color"
	"image/draw"
)

type Threads struct {
	poly *render.Polygon
	render.LayeredPoint
	Points       []floatgeom.Point2
	StaticPoints int
	velocity     physics.Vector
	Targets      []floatgeom.Point2
	BaseTargets  []floatgeom.Point2
	Color        color.Color
	Swinger      event.CID
	Tree         *collision.Tree
}

func NewThreads(c color.Color, size int, vel, anchor physics.Vector, layer, staticPoints int) *Threads {
	return &Threads{
		LayeredPoint: render.LayeredPoint{
			Vector: anchor,
			Layer:  render.NewLayer(layer),
		},
		StaticPoints: staticPoints,
		Color:        c,
		velocity:     vel,
		Points:       make([]floatgeom.Point2, size),
		Targets:      make([]floatgeom.Point2, size),
		BaseTargets:  make([]floatgeom.Point2, size),
		Tree:         collision.DefTree,
	}
}

func (ts *Threads) Draw(buff draw.Image) {
	ts.DrawOffset(buff, 0, 0)
}
func (ts *Threads) DrawOffset(buff draw.Image, xOff, yOff float64) {
	ts.update()
	var err error
	ts.poly, err = render.NewPolygon(ts.Points...)
	if err != nil {
		dlog.Error(err)
		return
	}
	ts.poly.Fill(ts.Color)
	ts.poly.DrawOffset(buff, xOff, yOff)
}

func (ts *Threads) SetTargets(fs ...floatgeom.Point2) {
	j := 0
	for i := ts.StaticPoints; i < len(ts.Points); i++ {
		ts.Targets[i] = fs[j]
		j = (j + 1) % len(fs)
	}
}

func (ts *Threads) ResetTargets() {
	ts.Targets = make([]floatgeom.Point2, len(ts.BaseTargets))
}

func (ts *Threads) update() {
	zeroPoint := floatgeom.Point2{}
	tScale := .2
	vScale := .005
	// move pixels in the direction of target, and velocity
	for i := ts.StaticPoints; i < len(ts.Points); i++ {
		tar := ts.Targets[i]
		if tar == zeroPoint {
			tar = ts.BaseTargets[i]
		}
		delta := tar.Sub(ts.Points[i])
		ts.Points[i] = ts.Points[i].Add(delta.MulConst(tScale))
	}
	for i := ts.StaticPoints; i < len(ts.Points); i++ {
		delta := floatgeom.Point2{ts.velocity.X(), ts.velocity.Y()}.MulConst(-1)
		ts.Points[i] = ts.Points[i].Add(delta.MulConst(vScale))
		sp := collision.NewUnassignedSpace(ts.Points[i].X(), ts.Points[i].Y(), .1, .1)
		if hit := ts.Tree.HitLabel(sp, Swing); hit != nil {
			ts.Swinger.Trigger("SwingHit", hit)
		}
	}
}
