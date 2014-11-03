// +build ignore

package main

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/config"
	"github.com/dradtke/allegory/example/state/game"
	"github.com/dradtke/allegory/example/state/loading"
)

func onLoad() {
	allegory.NewState(new(game.PlayingState))
}

/* -- Main -- */

func main() {
	config.SetWindowTitle("Let's Go!")
	config.SetWindowIcons(loading.IMG_DIR + "/standing.png")

	loading := &loading.LoadingState{*new(allegory.BaseGameState), onLoad}
	allegory.Run(loading)
}
