package graphendpoints

import (
	"encoding/json"
	"net/http"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/dbmanager"
	eq "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/dbmanager/endpoint_queries"
)

// --- Request / Response types ---

type graph3DRequest struct {
	XAxis    string   `json:"x_axis"`
	YAxis    string   `json:"y_axis"`
	GroupBy  string   `json:"group_by,omitempty"`
	Statuses []string `json:"statuses,omitempty"`
	Table    string   `json:"table"`
}

// ECharts scatter3D / bar3D consume [x, y, z] tuples directly.
// One series entry per group (or a single "all" entry when ungrouped).
//
//	{
//	  "axes": { "x": "network_size", "y": "first_k" },
//	  "series": [
//	    {
//	      "name": "hebbian",
//	      "stimulate_avg": [[x,y,z], ...],
//	      "stimulate_min": [[x,y,z], ...],
//	      "stimulate_max": [[x,y,z], ...],
//	      "learn_avg":     [[x,y,z], ...],
//	      "learn_min":     [[x,y,z], ...],
//	      "learn_max":     [[x,y,z], ...]
//	    }
//	  ]
//	}
type graph3DSeries struct {
	Name     string       `json:"name"`
	StimAvg  [][3]float64 `json:"stimulate_avg"`
	StimMin  [][3]float64 `json:"stimulate_min"`
	StimMax  [][3]float64 `json:"stimulate_max"`
	LearnAvg [][3]float64 `json:"learn_avg"`
	LearnMin [][3]float64 `json:"learn_min"`
	LearnMax [][3]float64 `json:"learn_max"`
}

type graph3DResponse struct {
	Axes   map[string]string `json:"axes"`
	Series []graph3DSeries   `json:"series"`
}

// --- Handler ---

func Graph3DHandler(qm *dbmanager.QueryManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req graph3DRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// --- Validate ---
		if !eq.AllowedAxes[req.XAxis] {
			http.Error(w, "invalid x_axis", http.StatusBadRequest)
			return
		}
		if !eq.AllowedAxes[req.YAxis] {
			http.Error(w, "invalid y_axis", http.StatusBadRequest)
			return
		}
		if req.XAxis == req.YAxis {
			http.Error(w, "x_axis and y_axis must differ", http.StatusBadRequest)
			return
		}
		if req.GroupBy != "" && !eq.AllowedGroupBy[req.GroupBy] {
			http.Error(w, "invalid group_by", http.StatusBadRequest)
			return
		}
		if !eq.AllowedTables[req.Table] {
			http.Error(w, "invalid table", http.StatusBadRequest)
			return
		}

		statuses := req.Statuses
		if len(statuses) == 0 {
			statuses = eq.DefaultStatuses
		}
		for _, s := range statuses {
			if !eq.AllowedStatuses[s] {
				http.Error(w, "invalid status: "+s, http.StatusBadRequest)
				return
			}
		}

		// --- Query ---
		params := eq.Graph3DParams{
			XAxis:    req.XAxis,
			YAxis:    req.YAxis,
			GroupBy:  req.GroupBy,
			Statuses: statuses,
			Table:    req.Table,
		}

		rawRows, err := eq.FetchGraph3D(r.Context(), qm, params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// --- Shape into ECharts series ---
		seriesMap := make(map[string]*graph3DSeries)
		order := []string{}

		for _, row := range rawRows {
			if _, exists := seriesMap[row.GroupVal]; !exists {
				order = append(order, row.GroupVal)
				seriesMap[row.GroupVal] = &graph3DSeries{Name: row.GroupVal}
			}
			s := seriesMap[row.GroupVal]
			pt := func(z float64) [3]float64 { return [3]float64{row.X, row.Y, z} }

			s.StimAvg = append(s.StimAvg, pt(row.AvgStim))
			s.StimMin = append(s.StimMin, pt(row.MinStim))
			s.StimMax = append(s.StimMax, pt(row.MaxStim))
			s.LearnAvg = append(s.LearnAvg, pt(row.AvgLearn))
			s.LearnMin = append(s.LearnMin, pt(row.MinLearn))
			s.LearnMax = append(s.LearnMax, pt(row.MaxLearn))
		}

		series := make([]graph3DSeries, 0, len(order))
		for _, name := range order {
			series = append(series, *seriesMap[name])
		}

		resp := graph3DResponse{
			Axes:   map[string]string{"x": req.XAxis, "y": req.YAxis},
			Series: series,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
