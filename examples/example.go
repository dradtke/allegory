// +build ignore

package main

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/config"
	"github.com/dradtke/allegory/examples/state/game"
	"github.com/dradtke/allegory/examples/state/loading"
)

func onLoad() {
	allegory.NewState(new(game.PlayingState))
}

/* -- Main -- */

func main() {
	config.SetWindowTitle("Let's Go!")
	config.SetWindowIcons(loading.IMG_DIR + "/standing.png")

	loading := new(loading.LoadingState)
	loading.OnLoad = onLoad
	allegory.Run(loading)
}
