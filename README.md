# Stock-exchange-simulator

A project to simulate and optimize task execution based on limited stock and process dependencies. Built in Go.

## Structure

- `cmd/stock_exchange`: Main program entry
- `cmd/checker`: Checker tool for validating execution logs
- `internal/parser`: Config file parsing logic
- `internal/scheduler`: Core scheduling logic
- `internal/checker`: Log validation logic
- `examples/`: Example input and output files

## Quick Start

### Building the Programs

```bash
# Build the main simulator
go build -o stock_exchange ./cmd/stock_exchange

# Build the checker
go build -o checker ./cmd/checker
```

### Running the Simulator

```bash
./stock_exchange <config_file>
```

This will:
1. Parse the configuration file
2. Run the simulation
3. Generate a `.log` file with the execution trace
4. Print final results to the terminal

### Running the Checker

```bash
./checker -config <config_file> -log <log_file>
```

This will validate that the log file represents a valid execution according to the configuration.

## Configuration File Format

Configuration files use a simple text format with three sections:

```
# Stocks
iron_ore:100
coal:50
wood:25

# Processes
smelting:(iron_ore:2;coal:1):(iron_ingot:1):30
crafting:(iron_ingot:3;wood:5):(sword:1):45
mining:(coal:1):(iron_ore:2):20

# Optimize
optimize:(time;sword)
```

### Stocks Section
- Format: `item_name:quantity`
- Defines initial stock levels for all items

### Processes Section
- Format: `process_name:(inputs):(outputs):duration`
- Inputs/outputs: `item:quantity;item:quantity`
- Duration: seconds (e.g., `30`) or Go duration format (e.g., `30s`, `1m`)

### Optimize Section
- Format: `optimize:(target1;target2)`
- Special target `time` optimizes for minimum execution time
- Other targets optimize for maximum production of that item

## Log File Format

The simulator generates log files in the format:
```
0:mining
0:smelting
1:mining
1:smelting
2:crafting
```

Each line represents: `cycle:process_name`

## Checker Program

The checker validates execution logs against configuration files to ensure simulations were executed correctly.

### Usage

```bash
./checker -config <config_file> -log <log_file>
```

### Command Line Options

- `-config string`: Path to the configuration file (required)
- `-log string`: Path to the log file to check (required)
- `-h`: Show help message

### What the Checker Validates

1. **Log Format**: Ensures each line follows `cycle:process_name` format
2. **Process Existence**: Verifies all processes exist in the configuration
3. **Cycle Order**: Checks that cycles are in non-decreasing order
4. **Stock Availability**: Simulates execution to verify sufficient stock at each step
5. **Timing**: Properly accounts for process durations and when outputs are produced

### Example Usage

#### Valid Execution
```bash
$ ./checker -config test_config.txt -log test_config.log
✓ Log validation successful!
Final cycle: 44
Total processes executed: 55

Final stock levels after input consumption:
  coal: 0
  iron_ingot: -15
  iron_ore: 50
  wood: 0
```

#### Invalid Execution (Insufficient Stock)
```bash
$ ./checker -config test_config.txt -log invalid.log
Checker error: Line 51, Cycle 50, Process smelting: not enough coal (need 1, have 0)
```

#### Invalid Execution (Cycle Order)
```bash
$ ./checker -config test_config.txt -log invalid_order.log
Checker error: Line 3, Cycle 1, Process mining: cycle 1 is less than previous cycle 2
```

### Error Types Detected

- **Insufficient Stock**: Process cannot start due to lack of required inputs
- **Invalid Cycle Order**: Cycles go backwards in time
- **Unknown Process**: Log references a process not defined in configuration
- **Invalid Format**: Malformed log entries or cycle numbers
- **Missing Files**: Configuration or log file cannot be opened

## Testing

### Running Unit Tests

```bash
# Test all packages
go test ./...

# Test specific package with verbose output
go test ./internal/checker -v
go test ./internal/parser -v
go test ./internal/scheduler -v
```

### Test Coverage

The checker includes comprehensive tests for:
- Valid execution traces
- Various error conditions (insufficient stock, invalid cycles, etc.)
- Edge cases (empty logs, malformed entries)
- Integration with real simulator output

## Example Workflow

1. **Create a configuration file** (e.g., `my_config.txt`)
2. **Run the simulator**:
   ```bash
   ./stock_exchange my_config.txt
   ```
   This generates `my_config.log`
3. **Validate the execution**:
   ```bash
   ./checker -config my_config.txt -log my_config.log
   ```
4. **Analyze results** or debug any validation errors

## Development

### Project Structure
```
├── cmd/
│   ├── stock_exchange/    # Main simulator binary
│   └── checker/           # Checker binary
├── internal/
│   ├── parser/           # Configuration file parsing
│   ├── scheduler/        # Simulation engine
│   └── checker/          # Log validation logic
└── README.md
```

### Adding New Features

1. **Parser**: Extend `internal/parser/` for new configuration formats
2. **Scheduler**: Modify `internal/scheduler/` for new optimization strategies
3. **Checker**: Update `internal/checker/` for additional validation rules
