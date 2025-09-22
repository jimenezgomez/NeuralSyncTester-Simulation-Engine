package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/dbmanager"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/session_manager"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var timeToLive = 1 * time.Minute

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Run a simulation in CLI mode (no SSE)",
	Run: func(cmd *cobra.Command, args []string) {

		// settings_list, baseBatchSettings, err := config_manager.LoadBatchSettingsFromFile("./simulation_settings.json")
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println(settings_list)
		// fmt.Println(baseBatchSettings)
		// return
		connString := "postgres://mtpm_user:mtpm_pass@localhost:5432/mtpm_db?sslmode=disable"

		db, err := sql.Open("postgres", connString)
		if err != nil {
			panic(err)
		}

		err = db.Ping()
		if err != nil {
			panic(err)
		}

		fmt.Println("Connected to Postgres!")

		manager := dbmanager.NewDBManager(db, dbmanager.InsertSessions, 500, 2*time.Second)
		defer manager.Close(context.Background())
		settings := engine.NewMTPMSettings([]int{3}, []int{3}, 3, 1, 1, "HEBBIAN", "NO_OVERLAP")
		trackedState := engine.NewTrackedState(settings)
		sessionManager := session_manager.NewSessionManager(timeToLive)
		sessionManager.AddMTPM(trackedState.UID, trackedState)
		ch := make(chan []byte, 10) // small buffer, prevent blocking
		trackedState.Subscribe(ch)
		defer trackedState.Unsubscribe(ch)

		go func() {
			for i := 0; i < 1_000_000; i++ {
				result := engine.SimulateTrackedSync(trackedState)
				sessionLog := dbmanager.NewLogFromResult(result)
				manager.Add(sessionLog)
			}
		}()

		for msg := range ch {
			var snapshot engine.SimulationInstance
			if err := json.Unmarshal(msg, &snapshot); err != nil {
				fmt.Println("failed to unmarshal snapshot:", err)
				continue
			}

			// Now you can access fields directly
			fmt.Printf("%s\n", trackedState.StartTime.Format(time.ANSIC))
			fmt.Printf("%s", snapshot.PrettyPrint())
		}

	},
}
