package main

import (
	"github.com/dradtke/gopher"
	"github.com/dradtke/gopher/bus"
	"github.com/dradtke/gopher/console"
	"github.com/dradtke/gopher/example/states/loading"
    "github.com/dradtke/gopher/example/events"
)

func onDot(msg string) {
    console.Info(msg)
}

func onCommand(cmd int) {
    println("hi")
}

func main() {
	gopher.Init(&loading.LoadingState{})

    bus.AddListener(events.DotNotifyEvent, onDot)
    bus.AddListener(bus.ConsoleCommandEvent, onCommand)
    defer bus.Clear()

	defer gopher.Cleanup()
	gopher.Loop()
}
