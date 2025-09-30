package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"time"

	config_manager "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/config_manager/load"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/dbmanager"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine/attacks"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/session_manager"
	_ "github.com/lib/pq"
	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/cobra"
	"github.com/xconstruct/go-pushbullet"
)

const DEFAULT_TTL = 1 * time.Minute
const SYNC_REPETITIONS = 1

// Refactor into enums/config file (?)
var AttackModes = []string{"NAIVE", "GEOMETRIC", "MAJORITY"}
var sessionManager *session_manager.SessionManager
var datamanager *dbmanager.DBManager[dbmanager.AttackSessionLog]
var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Run a simulation in CLI mode (no SSE)",
	Run: func(cmd *cobra.Command, args []string) {
		config_manager.InitEnv()
		maxSimulations := config_manager.GetMaxSimulations()
		// simulationPool := pool.New().WithMaxGoroutines(maxSimulations).WithErrors().WithFirstError()
		simulationPool := pool.New().WithMaxGoroutines(maxSimulations)
		envConfig := config_manager.LoadDBEnv()
		connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", //TODO: add SSL toggle in env
			envConfig.User, envConfig.Pass, envConfig.Host, envConfig.Port, envConfig.Name)
		db, err := sql.Open("postgres", connString)
		if err != nil {
			panic(err)
		}
		err = db.Ping()
		if err != nil {
			panic(err)
		}
		log.Println("Connected to Postgres!")

		datamanager = dbmanager.NewDBManager(db, dbmanager.InsertAttackSessions, 500, 2*time.Second)
		defer datamanager.Close(context.Background())

		pbApiKey := config_manager.LoadPBApiKey()
		pb := pushbullet.New(pbApiKey)
		devs, err := pb.Devices()
		if err != nil {
			panic(err)
		}

		sessionManager = session_manager.NewSessionManager(DEFAULT_TTL)
		batchGroup, err := config_manager.ScanAndLoadBatchSettings("./config")
		if err != nil {
			log.Fatalf("Error occurred: %v", err)
		}
		for _, batchCollected := range batchGroup {
			if batchCollected.Err != nil {
				continue
			}

			logMsg := fmt.Sprintf("File Path: %s \nConfigurations loaded: %d\n", batchCollected.Path, len(batchCollected.SettingsList))
			err = pb.PushNote(devs[0].Iden, fmt.Sprintf("Starting simulations for config file %s", filepath.Base(batchCollected.Path)), logMsg)
			if err != nil {
				panic(err)
			}
			fmt.Println(logMsg)
			for _, mtpmSettings := range batchCollected.SettingsList {
				simulationPool.Go(func() { RunInstance(mtpmSettings) })
			}

			logMsg = fmt.Sprintf("File %s has finished simulating %d configurations", batchCollected.Path, len(batchCollected.SettingsList))
			err = pb.PushNote(devs[0].Iden, fmt.Sprintf("Config file %s finished", filepath.Base(batchCollected.Path)), logMsg)
			if err != nil {
				panic(err)
			}
			fmt.Println(logMsg)
		}

		simulationPool.Wait()
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
	trackedState := engine.NewTrackedState(settings) //Move outside and maybe new type for attacks?
	sessionManager.AddMTPM(trackedState.UID, trackedState)
	defer sessionManager.DeleteMTPM(trackedState.UID)

	for _, attack_type := range AttackModes {
		for i := 0; i < SYNC_REPETITIONS; i++ {
			result := attacks.RunTrackedAttack(trackedState, attack_type)
			sessionLog, err := dbmanager.NewAttackSessionLog(result)
			if err != nil {
				panic("Fatal error when creating a new attack session log: " + err.Error())
			}
			datamanager.Add(sessionLog)
		}
	}

}
