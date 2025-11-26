# Configuration files for batch execution
The idea behind batch execution is to test the influence of different parameters in the synchronization time (measured in iterations). This way, the user can test the impact of adding new hidden layers, extending the range of synaptic weights, testing different learn rules, etc. 

All batch execution files work for synchronization and attack simulation.

## Creating new files
Creating your own files for batch execution may be tricky at first, here is what you need to know:
### Key differences between overlap scenarios
**For MTPMs with overlap:**
- You need to specify the amount of **neurons** on **each layer**.
- You need to specify the amount of **inputs** that each neuron receives on the **first layer**.

Example using 6 neurons on the first layer and 3 neurons on the second layer; the first layer has 2 inputs on each neuron:
```json
...
 "scenario": "partial_overlap",
 "n0_configs": [
    2
  ],
  "k_configs": [
    [
      6,
      3
    ]
  ]
...
```
*Note: Here the scenario can be `full_overlap` or `partial_overlap`. The amount of inputs between layers is calculated depending on the chosen scenario.*

**For MTPMs without overlap:**
- You need to specify the amount of **inputs** on **each layer**.
- You need to specify the amount of **neurons** that each neuron receives on the **last layer**.

Example using 3 neurons on the last layer; the first and second layer both have 2 inputs on each neuron:
```json
...
 "scenario": "no_overlap",
 "klast_configs": [
    3
  ],
  "n_configs":[
    [
      2,
      2
    ]
  ]
...
```
*Note: MTPMs without overlap are created from the last layer to the first. In the previous example we start with 3 neurons on the last layer that need 2 inputs, so the previous layer needs at least 3\*2=6 neurons.*

### Common parameters
The following parameters are common for all scenarios:
```json
 "m_configs": [
    1
  ],
  "l_configs": [
    3
  ],
  "learn_rules": [
    "Hebbian",
    "Anti-Hebbian",
    "Random-Walk"
  ],
```
- `m_configs`: The range of the input stimulus. Must be over 0. Using a number over 1 means non-binary inputs on the first layer. 
- `l_configs`: The range of the synaptic weights. Must be over 0.
- `learn_rules`: All of the specified learn rules will be tested. If you want to add new rules, make sure to add the logic to parse the configuration file.


## Additional notes

### Calculating the amount of inputs for MTPMs with Partial Overlap
The algorithm for calculating the amount of inputs that each neuron receives on a hidden layer is calculated as follows:

$$N_{h} = K_{h-1} - K_{h} + 1$$

With $N_h$ the amount of stimulus of each neuron on the layer $h$, $K_h$ the amount of neurons of the layer $h$ and $K_{h-1}$ the amount of neurons in the previous layer.

This means that, for partial overlap scenarios, **the amount of neurons in a hidden layer must always be less than the previous layer**.

### Parsing the MTPM architecures
The batch configuration loader logic is in the file `internal/config_manager/load/batchConfigLoader.go`.

The function `LoadBatchSettingsFromFile(filename string)` takes the filename as a parameter and returns a list of ALL the possible combinations. It will first parse the architectures of the file: the amount of neurons in each layer (`K`) and the amount of inputs that each neuron needs in each layer (`N`). 
Then it will create all possible combinations for the common parameters (`L`,`M`,and the `Learning Rules`).

### Parameters with no effect on the simulation
The following parameters have no effect, as they have been moved to the other config files:

- "max_session_count": moved to `simulation.yaml` -> sync_repetitions
- "max_iterations": moved to `simulation.yaml` -> iteration_limit
- "max_worker_count": moved to `.env` -> MAX_SIMULATIONS

_NOTE: MAX_SIMULATIONS means max **concurrent** simulations, the name may be a bit misleading_
