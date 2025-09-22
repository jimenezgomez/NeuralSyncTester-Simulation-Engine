package dbmanager

import (
	"context"
	"database/sql"
	"encoding/json"
)

func InsertSessions(ctx context.Context, db *sql.DB, sessions []SyncSessionLog) error {
	if len(sessions) == 0 {
		return nil
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO sessions (
			network_size, first_k, first_n, last_k, last_n,
			start_time, end_time, stimulate_iterations, learn_iterations,
			k, n, l, m, h, learn_rule, scenario,
			initial_state, final_state, session_status
		) VALUES (
			$1,$2,$3,$4,$5,
			$6,$7,$8,$9,
			$10,$11,$12,$13,$14,$15,$16,
			$17,$18,$19
		)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, s := range sessions {
		kBytes, _ := json.Marshal(s.K)
		nBytes, _ := json.Marshal(s.N)
		initBytes, _ := json.Marshal(s.InitialState)
		finalBytes, _ := json.Marshal(s.FinalState)

		_, err = stmt.ExecContext(ctx,
			s.NetworkSize, s.FirstK, s.FirstN, s.LastK, s.LastN,
			s.StartTime, s.EndTime, s.StimulateIterations, s.LearnIterations,
			kBytes, nBytes, s.L, s.M, s.H, s.LearnRule, s.Scenario,
			initBytes, finalBytes, s.SessionStatus,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
