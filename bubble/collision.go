package bubble

import "github.com/oakmound/oak/collision"

const (
	// The only collision label we need for this demo is 'ground',
	// indicating something we shouldn't be able to fall or walk through
	Ground collision.Label = 1
)
