package engine

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_core"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_learnRules"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_stimHandlers"
)

type MTPMSettings struct {
	K                   []int
	N                   []int
	L                   int
	M                   int
	H                   int
	LearnRule           string
	Scenario            string
	stimulationHandlers tpm_stimHandlers.TPMStimulationHandlers
	learnRuleHandler    tpm_learnRules.TPMLearnRuleHandler
}

func NewMTPMSettings(K, N []int, L, M, H int, learnRule, scenario string) MTPMSettings {
	var stimHandler tpm_stimHandlers.TPMStimulationHandlers
	var ruleHandler tpm_learnRules.TPMLearnRuleHandler

	scenario = strings.ToUpper(strings.TrimSpace(scenario))
	learnRule = strings.ToUpper(strings.TrimSpace(learnRule))

	switch scenario {
	case "NO_OVERLAP", "NO OVERLAP":
		stimHandler = tpm_stimHandlers.NoOverlapTPM{}
	case "PARTIAL_OVERLAP", "PARTIAL OVERLAP", "PARTIALLY_CONNECTED", "PARTIALLY CONNECTED":
		stimHandler = tpm_stimHandlers.PartialOverlapTPM{}
	case "FULL_OVERLAP", "FULL OVERLAP", "FULLY_CONNECTED", "FULLY CONNECTED":
		stimHandler = tpm_stimHandlers.FullOverlapTPM{}
	default:
		panic("No valid scenario was selected.")

	}

	switch learnRule {
	case "HEBBIAN":
		ruleHandler = tpm_learnRules.HebbianLearnRule{}
	case "ANTI-HEBBIAN", "ANTIHEBBIAN", "ANTI HEBBIAN":
		ruleHandler = tpm_learnRules.AntiHebbianLearnRule{}
	case "RANDOM-WALK", "RANDOM", "RANDOMWALK", "RANDOM WALK":
		ruleHandler = tpm_learnRules.RandomWalkLearnRule{}
	default:
		panic("No valid learning rule was selected.")
	}

	return MTPMSettings{
		H:                   H,
		K:                   K,
		N:                   N,
		L:                   L,
		M:                   M,
		LearnRule:           learnRule,
		Scenario:            scenario,
		stimulationHandlers: stimHandler,
		learnRuleHandler:    ruleHandler,
	}
}

// Core network state, reusable across different types
type MTPMState struct {
	Weights       [][][]int
	InputBuffer   [][][]int
	OutputBuffer  [][]int
	NetworkOutput int
}

func NewMTPMState(settings MTPMSettings) MTPMState {
	newWeights := make([][][]int, settings.H)
	for layer := 0; layer < settings.H; layer++ {
		newWeights[layer] = tpm_core.CreateRandomLayerWeightsArray(settings.K[layer], settings.N[layer], settings.L)
	}
	return MTPMState{
		Weights:       newWeights,
		InputBuffer:   make([][][]int, settings.H),
		OutputBuffer:  make([][]int, settings.H),
		NetworkOutput: 0,
	}
}

func (mtpmState *MTPMState) DeepCopy() *MTPMState {
	copyState := &MTPMState{
		NetworkOutput: mtpmState.NetworkOutput,
	}

	// Deep copy Weights
	if mtpmState.Weights != nil {
		copyState.Weights = make([][][]int, len(mtpmState.Weights))
		for i := range mtpmState.Weights {
			copyState.Weights[i] = make([][]int, len(mtpmState.Weights[i]))
			for j := range mtpmState.Weights[i] {
				copyState.Weights[i][j] = make([]int, len(mtpmState.Weights[i][j]))
				copy(copyState.Weights[i][j], mtpmState.Weights[i][j])
			}
		}
	}

	// Deep copy InputBuffer
	if mtpmState.InputBuffer != nil {
		copyState.InputBuffer = make([][][]int, len(mtpmState.InputBuffer))
		for i := range mtpmState.InputBuffer {
			copyState.InputBuffer[i] = make([][]int, len(mtpmState.InputBuffer[i]))
			for j := range mtpmState.InputBuffer[i] {
				copyState.InputBuffer[i][j] = make([]int, len(mtpmState.InputBuffer[i][j]))
				copy(copyState.InputBuffer[i][j], mtpmState.InputBuffer[i][j])
			}
		}
	}

	// Deep copy OutputBuffer
	if mtpmState.OutputBuffer != nil {
		copyState.OutputBuffer = make([][]int, len(mtpmState.OutputBuffer))
		for i := range mtpmState.OutputBuffer {
			copyState.OutputBuffer[i] = make([]int, len(mtpmState.OutputBuffer[i]))
			copy(copyState.OutputBuffer[i], mtpmState.OutputBuffer[i])
		}
	}

	return copyState
}

func (mtpmState *MTPMState) String() string {
	return fmt.Sprintf(
		"MTPMState{\n  Weights: %v,\n  InputBuffer: %v,\n  OutputBuffer: %v,\n  NetworkOutput: %d\n}",
		mtpmState.Weights, mtpmState.InputBuffer, mtpmState.OutputBuffer, mtpmState.NetworkOutput,
	)
}

// Optional: more human-readable nested printing
func (mtpmState *MTPMState) PrettyPrint() string {
	out := fmt.Sprintf("NetworkOutput: %d\n", mtpmState.NetworkOutput)

	out += "Weights:\n"
	for i, layer := range mtpmState.Weights {
		out += fmt.Sprintf(" Layer %d:\n", i)
		for j, row := range layer {
			out += fmt.Sprintf("  Row %d: %v\n", j, row)
		}
	}
	// out += "InputBuffer:\n"
	// for i, layer := range s.InputBuffer {
	// 	out += fmt.Sprintf(" Layer %d:\n", i)
	// 	for j, row := range layer {
	// 		out += fmt.Sprintf("  Row %d: %v\n", j, row)
	// 	}
	// }
	// out += "OutputBuffer:\n"
	// for i, row := range s.OutputBuffer {
	// 	out += fmt.Sprintf(" Row %d: %v\n", i, row)
	// }

	return out
}

type SimulationState struct {
	StateA MTPMState
	StateB MTPMState
}

type SimulationProgress struct {
	StimulateIterations int
	LearnIterations     int
}

type SimulationInstance struct {
	SimulationNumber int
	SimulationState
	SimulationProgress
}

func (s *SimulationInstance) DeepCopy() *SimulationInstance {
	copyInstance := &SimulationInstance{
		SimulationProgress: SimulationProgress{
			StimulateIterations: s.StimulateIterations,
			LearnIterations:     s.LearnIterations,
		},
	}

	// Deep copy both states
	copyInstance.StateA = *s.StateA.DeepCopy()
	copyInstance.StateB = *s.StateB.DeepCopy()

	return copyInstance
}

func (s *SimulationInstance) String() string {
	return fmt.Sprintf(
		"SimulationInstance{\n  StateA: %v,\n  StateB: %v,\n  StimulateIterations: %d,\n  LearnIterations: %d\n}",
		s.StateA, s.StateB, s.StimulateIterations, s.LearnIterations,
	)
}

// Optional human-readable
func (s *SimulationInstance) PrettyPrint() string {
	out := "StateA:\n" + s.StateA.PrettyPrint() + "\n"
	out += "StateB:\n" + s.StateB.PrettyPrint() + "\n"
	out += fmt.Sprintf("StimulateIterations: %d\nLearnIterations: %d\n", s.StimulateIterations, s.LearnIterations)
	return out
}

type SimulationResult struct {
	Settings      MTPMSettings
	InitialState  SimulationState
	FinalState    SimulationState
	Iterations    SimulationProgress
	SessionStatus string
	StartTime     time.Time
	EndTime       time.Time
}

type TrackedMTPMSession struct {
	UID            string
	StartTime      time.Time
	Settings       MTPMSettings
	MaxSimCount    int
	snapshotBuffer *SimulationBuffer
	subs           map[chan []byte]struct{}
	subCount       atomic.Int64
	subsLock       sync.RWMutex
	simProgress    atomic.Int64
}

func NewTrackedSession(settings MTPMSettings, maxSimCount, maxBufferCount int) *TrackedMTPMSession {
	uid, err := generateToken(16)
	if err != nil {
		uid = time.Now().Local().String()
	}
	ts := &TrackedMTPMSession{
		UID:            uid,
		Settings:       settings,
		MaxSimCount:    maxSimCount,
		snapshotBuffer: NewSimulationBuffer(maxBufferCount),
		subs:           make(map[chan []byte]struct{}),
	}
	return ts
}

func (ts *TrackedMTPMSession) Subscribe(ch chan []byte) {

	ts.subsLock.Lock()
	ts.subs[ch] = struct{}{}
	ts.subsLock.Unlock()

	ts.subCount.Add(1)
}

func (ts *TrackedMTPMSession) Unsubscribe(ch chan []byte) {
	ts.subsLock.Lock()
	delete(ts.subs, ch)
	ts.subsLock.Unlock()

	ts.subCount.Add(-1)
}

func (ts *TrackedMTPMSession) UpdateSnapshot(newSnap *SimulationInstance) {
	ts.snapshotBuffer.Push(*newSnap)
}

func (ts *TrackedMTPMSession) GetSnapshotRaw() []SimulationInstance {
	snap := ts.snapshotBuffer.Flush()
	return snap
}

func (ts *TrackedMTPMSession) UpdateAllSubscribers() {
	ts.subsLock.RLock()
	defer ts.subsLock.RUnlock()

	rawData := ts.GetSnapshotRaw()
	marshalledData, err := json.Marshal(rawData)
	if err != nil {

	}
	for subscriber := range ts.subs {
		select {
		case subscriber <- marshalledData:
		default: // skip slow/broken client
		}
	}
}

func (ts *TrackedMTPMSession) GetSubCount() int64 {
	return ts.subCount.Load()
}

func (ts *TrackedMTPMSession) GetSettings() MTPMSettings {
	return ts.Settings
}

func (ts *TrackedMTPMSession) GetSessionProgress() int64 {
	return ts.simProgress.Load()
}

func (ts *TrackedMTPMSession) AddProgress() int64 {
	return ts.simProgress.Add(1)
}

type SimulationBuffer struct {
	mu       sync.Mutex
	frames   []SimulationInstance
	capacity int
}

func NewSimulationBuffer(capacity int) *SimulationBuffer {
	return &SimulationBuffer{
		frames:   make([]SimulationInstance, 0, capacity),
		capacity: capacity,
	}
}

func (b *SimulationBuffer) Push(instance SimulationInstance) (full bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.frames = append(b.frames, instance)
	return len(b.frames) >= b.capacity
}

func (b *SimulationBuffer) Flush() []SimulationInstance {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]SimulationInstance, len(b.frames))
	copy(out, b.frames)
	b.frames = b.frames[:0] // reset without realloc
	return out
}

func generateToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
