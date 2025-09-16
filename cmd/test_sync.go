package main

import (
	"fmt"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
)

func main() {
	settings := engine.NewMTPMSettings(
		[]int{3},
		[]int{3},
		1,
		3,
		1,
		"HEBBIAN",
		"NO_OVERLAP")
	// result := engine.SimulateSimpleSync(settings)
	trackedState := engine.NewTrackedState(settings)
	result := engine.SimulateTrackedSync(trackedState)
	fmt.Println(result)

}
