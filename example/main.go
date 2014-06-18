package main

import (
	"github.com/dradtke/gopher/example/states/loading"
	"flag"
	"os"
	"runtime/pprof"
	"github.com/dradtke/gopher"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

// TODO: refactor all references to Allegro into Gopher
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

	gopher.Init(&loading.LoadingState{})
	defer gopher.Cleanup()

	gopher.Loop()
}
