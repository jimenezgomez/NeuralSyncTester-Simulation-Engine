package cmd

import (
	"fmt"
	"log"

	config_manager "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/config_manager/load"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/dbmanager"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine/sync"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Run sync simulations in CLI mode (no Endpoints available)",
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
				simulationPool.Go(func() { RunSyncInstance(mtpmSettings) })
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

func RunSyncInstance(settings engine.MTPMSettings) {
	trackedState := engine.NewTrackedSession(settings, GlobalSimulationConfig.SyncRepetitions) //Move outside and maybe new type for attacks?
	sessionManager.AddMTPM(trackedState.UID, trackedState)
	defer sessionManager.DeleteMTPM(trackedState.UID)

	for i := 0; i < GlobalSimulationConfig.SyncRepetitions; i++ {
		sync.SimulateTrackedSync(trackedState, GlobalSimulationConfig, GlobalTrackingConfig)
		result := sync.SimulateTrackedSync(trackedState, GlobalSimulationConfig, GlobalTrackingConfig)
		sessionLog, err := dbmanager.NewLogFromResult(result)
		if err != nil {
			panic("Fatal error when creating a new attack session log: " + err.Error())
		}
		syncDataManager.Add(sessionLog)
	}

}
