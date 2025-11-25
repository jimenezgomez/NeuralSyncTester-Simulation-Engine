package cmd

import (
	"fmt"
	"log"

	config_manager "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/config_manager/load"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/dbmanager"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine/attacks"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Run a simulation in CLI mode (no SSE)",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Using config directory: ", GlobalSimulationConfig.BatchPath)
		batchGroup, err := config_manager.ScanAndLoadBatchSettings(GlobalSimulationConfig.BatchPath)
		if err != nil {
			log.Fatalf("Error occurred: %v", err)
		}
		for _, batchCollected := range batchGroup {
			if batchCollected.Err != nil {
				continue
			}

			logMsg := fmt.Sprintf("File Path: %s \nConfigurations loaded: %d\n", batchCollected.Path, len(batchCollected.SettingsList))
			// err = pb.PushNote(devs[0].Iden, fmt.Sprintf("Starting simulations for config file %s", filepath.Base(batchCollected.Path)), logMsg)
			if err != nil {
				panic(err)
			}
			fmt.Println(logMsg)
			for _, mtpmSettings := range batchCollected.SettingsList {
				simulationPool.Go(func() { RunInstance(mtpmSettings) })
			}
		}

		simulationPool.Wait()
		// err = pb.PushNote(devs[0].Iden, "All config files have finished", "All files have finished simulating attacks.")
		if err != nil {
			panic(err)
		}
		fmt.Println("All configuration files finished.")
		// if err := simulationPool.Wait(); err != nil {
		// 	log.Fatalf("Error occurred: %v", err)
		// }

		// trackedState.Subscribe(ch)
		// defer trackedState.Unsubscribe(ch)
		// for msg := range ch {
		// 	var snapshot engine.SimulationInstance
		// 	if err := json.Unmarshal(msg, &snapshot); err != nil {
		// 		fmt.Println("failed to unmarshal snapshot:", err)
		// 		continue
		// 	}

		// 	fmt.Printf("%s\n", trackedState.StartTime.Format(time.ANSIC))
		// 	fmt.Printf("%s", snapshot.PrettyPrint())
		// }

	},
}

func RunInstance(settings engine.MTPMSettings) {
	trackedState := engine.NewTrackedSession(settings, GlobalSimulationConfig.SyncRepetitions) //Move outside and maybe new type for attacks?
	sessionManager.AddMTPM(trackedState.UID, trackedState)
	defer sessionManager.DeleteMTPM(trackedState.UID)

	for _, attack_type := range GlobalSimulationConfig.AttackModes {
		for i := 0; i < GlobalSimulationConfig.SyncRepetitions; i++ {
			result := attacks.RunTrackedAttack(trackedState, attack_type, GlobalSimulationConfig, GlobalTrackingConfig)
			sessionLog, err := dbmanager.NewAttackSessionLog(result)
			if err != nil {
				panic("Fatal error when creating a new attack session log: " + err.Error())
			}
			datamanager.Add(sessionLog)
		}
	}

}
