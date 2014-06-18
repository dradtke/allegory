package main

import (
	"flag"
	al "github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/dialog"
	"github.com/dradtke/go-allegro/allegro/font"
	prim "github.com/dradtke/go-allegro/allegro/primitives"
	"github.com/dradtke/gopher"
	"github.com/dradtke/gopher/config"
	"github.com/dradtke/gopher/subsystems/console"
	"github.com/dradtke/gopher/states/loading"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            panic(err)
        }
        defer f.Close()

        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

    {
        if err := al.Install(); err != nil {
            panic(err)
        }
        defer al.Uninstall()

        if err := dialog.Install(); err != nil {
            panic(err)
        }

        if err := prim.Install(); err != nil {
            gopher.Fatal(err)
        }
        defer prim.Uninstall()

        font.Install()
        defer font.Uninstall()
    }

    gopher.Init()
    defer gopher.Cleanup()

	// Initialize subsystems.
    // TODO: figure out why this breaks things when debugging...
	//console.Init(gopher.EventQueue())

	// Set the screen to black.
	al.ClearToColor(config.BlankColor())
	al.FlipDisplay()

    gopher.NewState(&loading.LoadingState{})
    gopher.Loop()

	al.ClearToColor(config.BlankColor())
	al.FlipDisplay()

	gopher.Display().SetWindowTitle("Shutting down...")
	console.Save()
}
