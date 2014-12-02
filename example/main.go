package main

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/example/g"
	"github.com/dradtke/allegory/example/playing"
)

func main() {
	allegory.DefState(g.PLAYING).
		Init(playing.Init).
		Update(playing.Update).
		HandleEvent(playing.HandleEvent)

	allegory.Run(g.PLAYING)
}
