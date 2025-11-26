# Neural Sync Tester — Simulation Engine

This repository contains the **Simulation Engine** for the Neural Sync Tester Suite.  
Its purpose is to automate large-scale synchronization experiments of **Tree Parity Machines (TPM)** and **Multilayer Tree Parity Machines (MTPM)**, generating reproducible results for research and analysis.

## Features
- Modular simulation pipeline for TPM and MTPM architectures  
- Batch execution with configurable parameters (Examples available in `config`)
- Automatic result persistence (PostgreSQL)  
- Extensible design for new architectures and learning rules

## Quick Start

### Requirements
- **Go** (v1.22+ recommended)
- **PostgreSQL** (Docker compose file included for quick deployment)

### Configuring simulations
There are three types of configurations:
- Enviroment variables used by the database and the Simulation Engine (Found in the `.env` file at the root of the project)
- Engine configuration in the `config` directory (`simulation.yaml` and `tracking.yaml`)
- MTPM Architecture configuration for batch execution (default path is `config/batches`, can be changed in `simulation.yaml`)

The provided examples feature three architectures for all overlap scenarios (no overlap, partial overlap and full overlap MTPMs).
![MTPMs available in the `config/batches` directory](docs/mtpm_examples.drawio.svg)


Additional notes on how MTPMs are loaded from the configuration files is available at the `docs` directory.
### Build
On the root of the project run:
```sh
go build -o simulation-engine .
```
You can also run the project without building a binary (full example below)

### Run a Sync Simulation
Sync simulations are stored on the table `sync_sessions`
```sh
./simulation-engine sync
```
Or running without creating a binary file:
```sh
go run . sync
```

### Run an Attack Simulation
Attack simulations are stored on the table `attack_sessions`
```sh
./simulation-engine attack
```
Or running without creating a binary file:
```sh
go run . attack
```

