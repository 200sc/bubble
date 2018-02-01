package main

import (
	"github.com/200sc/bubble/bubble"
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/render"
)

func main() {
	bubble.AddScenes()

	oak.SetupConfig.BatchLoad = true
	oak.SetupConfig.Screen.Scale = 2
	oak.SetupConfig.Screen.Height = 240
	oak.SetupConfig.Screen.Width = 320
	oak.SetupConfig.DrawFrameRate = 120

	render.SetDrawStack(
		render.NewHeap(false),
		render.NewDrawFPS(),
		render.NewLogicFPS(),
	)

	oak.Init("setup")
}
