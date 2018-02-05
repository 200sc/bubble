package bubble

import (
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/shape"
	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"

	"image/color"
)

var (
	snow = particle.And(
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

	backSnow = particle.And(
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
)