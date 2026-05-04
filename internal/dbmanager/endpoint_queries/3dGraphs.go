package endpointqueries

import (
	"context"
	"fmt"
	"strings"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/dbmanager"
)

// Allowed value whitelists — kept here so validation lives next to the query.

var AllowedAxes = map[string]bool{
	"network_size": true,
	"first_k":      true,
	"first_n":      true,
	"last_k":       true,
	"last_n":       true,
	"l":            true,
	"m":            true,
	"h":            true,
}

var AllowedGroupBy = map[string]bool{
	"learn_rule": true,
	"scenario":   true,
}

var AllowedStatuses = map[string]bool{
	"LIMIT_REACHED":  true,
	"ON_SYNC":        true,
	"ATTACK_SUCCESS": true,
}

var AllowedTables = map[string]bool{
	"attack_sessions": true,
	"sync_sessions":   true,
}

var DefaultStatuses = []string{"LIMIT_REACHED", "ON_SYNC", "ATTACK_SUCCESS"}

// Graph3DParams holds the validated query parameters.
type Graph3DParams struct {
	XAxis    string
	YAxis    string
	GroupBy  string // empty means no grouping
	Statuses []string
	Table    string
}

// Graph3DRow is a single result row returned from the database.
type Graph3DRow struct {
	GroupVal                     string
	X, Y                         float64
	AvgStim, MinStim, MaxStim    float64
	AvgLearn, MinLearn, MaxLearn float64
}

// FetchGraph3D runs the aggregation query and returns raw rows.
func FetchGraph3D(ctx context.Context, qm *dbmanager.QueryManager, p Graph3DParams) ([]Graph3DRow, error) {
	selectCols := []string{
		p.XAxis,
		p.YAxis,
		"AVG(stimulate_iterations) AS avg_stimulate",
		"MIN(stimulate_iterations) AS min_stimulate",
		"MAX(stimulate_iterations) AS max_stimulate",
		"AVG(learn_iterations)     AS avg_learn",
		"MIN(learn_iterations)     AS min_learn",
		"MAX(learn_iterations)     AS max_learn",
	}
	groupCols := []string{p.XAxis, p.YAxis}

	if p.GroupBy != "" {
		selectCols = append([]string{p.GroupBy}, selectCols...)
		groupCols = append([]string{p.GroupBy}, groupCols...)
	}

	placeholders := make([]string, len(p.Statuses))
	args := make([]interface{}, len(p.Statuses))
	for i, s := range p.Statuses {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = s
	}

	query := fmt.Sprintf(
		`SELECT %s
		 FROM %s
		 WHERE session_status IN (%s)
		 GROUP BY %s
		 ORDER BY %s`,
		strings.Join(selectCols, ", "),
		p.Table,
		strings.Join(placeholders, ", "),
		strings.Join(groupCols, ", "),
		strings.Join(groupCols, ", "),
	)

	rows, err := qm.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("graph3d query: %w", err)
	}
	defer rows.Close()

	var result []Graph3DRow
	for rows.Next() {
		var row Graph3DRow
		var targets []interface{}

		if p.GroupBy != "" {
			targets = []interface{}{
				&row.GroupVal,
				&row.X, &row.Y,
				&row.AvgStim, &row.MinStim, &row.MaxStim,
				&row.AvgLearn, &row.MinLearn, &row.MaxLearn,
			}
		} else {
			row.GroupVal = "all"
			targets = []interface{}{
				&row.X, &row.Y,
				&row.AvgStim, &row.MinStim, &row.MaxStim,
				&row.AvgLearn, &row.MinLearn, &row.MaxLearn,
			}
		}

		if err := rows.Scan(targets...); err != nil {
			return nil, fmt.Errorf("graph3d scan: %w", err)
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("graph3d rows: %w", err)
	}

	return result, nil
}
