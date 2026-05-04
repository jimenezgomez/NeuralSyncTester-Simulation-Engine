package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	config_manager "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/config_manager/load"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/dbmanager"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/session_manager"
	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/cobra"
	"github.com/xconstruct/go-pushbullet"
)

var rootCmd = &cobra.Command{
	Use:   "neural-sync",
	Short: "NeuralSyncTester CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		config_manager.InitEnv()

		// 2. Load YAML configs
		simPath := filepath.Join("config", "simulation.yaml")
		trackPath := filepath.Join("config", "tracking.yaml")

		var err error
		GlobalSimulationConfig, err = config_manager.LoadSimulationConfig(simPath)
		if err != nil {
			return fmt.Errorf("simulation config: %w", err)
		}

		GlobalTrackingConfig, err = config_manager.LoadTrackingConfig(trackPath)
		if err != nil {
			return fmt.Errorf("tracking config: %w", err)
		}

		//Load configuration and set-up the simulation
		maxSimulations := config_manager.GetMaxSimulations()
		simulationPool = pool.New().WithMaxGoroutines(maxSimulations)
		dbEnv := config_manager.LoadDBEnv()

		//Connect to postgres
		connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", //TODO: add SSL toggle in env
			dbEnv.User, dbEnv.Pass, dbEnv.Host, dbEnv.Port, dbEnv.Name)
		db, err := sql.Open("postgres", connString)
		if err != nil {
			return fmt.Errorf("open DB connection: %w", err)
		}
		err = db.Ping()
		if err != nil {
			return fmt.Errorf("DB ping: %w", err)
		}
		log.Println("Connected to Postgres!")
		//Setup dbmanager to insert data every 2 seconds
		//Now setup is done depending on the command - review later
		queryManager = dbmanager.NewQueryManager(db)
		attDataManager = dbmanager.NewDBManager(db, dbmanager.InsertAttackSessions, 500, 2*time.Second)
		syncDataManager = dbmanager.NewDBManager(db, dbmanager.InsertSyncSessions, 500, 2*time.Second)
		//defer datamanager.Close(context.Background())
		//Note: the database is closed on PersistentPostRunE

		//Set-up PushBullet
		pbApiKey := config_manager.LoadPBApiKey()

		if strings.TrimSpace(pbApiKey) == "" {
			pbClient = nil
		} else {
			pbClient = pushbullet.New(pbApiKey)
			pbDevices, err = pbClient.Devices()
			if err != nil {
				pbClient = nil
			}
		}

		//Create sim session manager
		sessionManager = session_manager.NewSessionManager(GlobalTrackingConfig.ParsedTTL)
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		// Close DB after command is done
		attDataManager.Close(context.Background())
		return nil
	},
}

var (
	GlobalSimulationConfig config_manager.SimulationConfig
	GlobalTrackingConfig   config_manager.TrackingConfig
	sessionManager         *session_manager.SessionManager
	simulationPool         *pool.Pool
	queryManager           *dbmanager.QueryManager
	attDataManager         *dbmanager.DBManager[dbmanager.AttackSessionLog]
	syncDataManager        *dbmanager.DBManager[dbmanager.SyncSessionLog]
	pbClient               *pushbullet.Client
	pbDevices              []*pushbullet.Device
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(attackCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(sizeSolverCmd)
}
