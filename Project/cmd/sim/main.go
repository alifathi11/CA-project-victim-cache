package main

import (
	"fmt"

	"victimcacheproject/internal/config"
	"victimcacheproject/internal/system"
)

func main() {
	cfg := config.Default()
	sim := system.New(cfg)

	fmt.Println("Victim Cache Project 6 — scaffold")
	fmt.Printf("Config: %+v\n", cfg)

	if err := sim.Validate(); err != nil {
		panic(err)
	}

	// TODO(stage-2): Build Akita components and connections.
	// TODO(stage-3): Load benchmark requests.
	// TODO(stage-4): Run the Akita engine and collect metrics.
	fmt.Println("Scaffold is valid. Simulation is not implemented yet.")
}
