package main

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/example/playing"
	"github.com/dradtke/allegory/example/playing/paused"
)

func main() {
	playing.Register()
	paused.Register()

	allegory.Run("playing")
}
