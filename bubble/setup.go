package bubble

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/scene"
)

func addSetupScene() {
	oak.Add("setup", func(string, interface{}) {
		SetupPlayer()
	}, func() bool {
		return false
	}, func() (string, *scene.Result) {
		return "bubble", nil
	}) 
}
