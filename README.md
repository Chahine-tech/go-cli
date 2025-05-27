# LogAnalyzer - Distributed Log Analysis Tool

A powerful command-line tool written in Go for concurrent analysis of log files from multiple sources. LogAnalyzer helps system administrators centralize and analyze logs in parallel with robust error handling capabilities.

## 🚀 Features

- **Concurrent Processing**: Analyzes multiple log files in parallel using goroutines
- **Custom Error Handling**: Implements custom error types with proper `errors.Is()` and `errors.As()` handling
- **JSON Configuration**: Uses JSON files for flexible log configuration
- **JSON Reporting**: Exports analysis results to JSON format
- **CLI Interface**: Built with Cobra framework for intuitive command-line usage
- **Real-time Progress**: Shows real-time analysis progress and results
- **Modular Architecture**: Clean separation of concerns with internal packages

## 📋 Requirements

- Go 1.20 or higher
- Access to log files you want to analyze

## 🛠 Installation

### Option 1: Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd loganalyzer

# Install dependencies
go mod tidy

# Build the application
go build -o loganalyzer main.go

# Run the application
./loganalyzer --help
```

### Option 2: Direct Go Install

```bash
go install <repository-url>@latest
```

## 📖 Usage

### Basic Usage

```bash
# Analyze logs using a configuration file
loganalyzer analyze --config config.json

# Analyze logs and save results to a file
loganalyzer analyze --config config.json --output report.json

# Using short flags
loganalyzer analyze -c config.json -o report.json
```

### Help and Documentation

```bash
# Show general help
loganalyzer --help

# Show help for the analyze command
loganalyzer analyze --help

# Show version
loganalyzer --version
```

## 📁 Configuration File Format

The configuration file is a JSON array containing log file specifications:

```json
[
  {
    "id": "web-server-1",
    "path": "/var/log/nginx/access.log",
    "type": "nginx access"
  },
  {
    "id": "app-backend-2",
    "path": "/var/log/my_app/errors.log",
    "type": "custom application"
  },
  {
    "id": "system-logs",
    "path": "/var/log/syslog",
    "type": "system log"
  }
]
```

### Configuration Fields

- **id**: Unique identifier for the log file (required)
- **path**: Absolute or relative path to the log file (required)
- **type**: Description of the log type (required)

## 📊 Output Format

### Console Output

The tool provides real-time progress updates and a summary:

```
Loading configuration from: config.json
Loaded configuration with 3 log files
Starting analysis of 3 log files...
Processing log: web-server-1 (/var/log/nginx/access.log)
Processing log: app-backend-2 (/var/log/my_app/errors.log)
Processing log: system-logs (/var/log/syslog)
✓ Completed analysis of log: web-server-1
✗ File error for log app-backend-2: file not found or inaccessible: /var/log/my_app/errors.log
✓ Completed analysis of log: system-logs

=== Analysis Summary ===
✓ [web-server-1] /var/log/nginx/access.log: Analysis completed successfully.
✗ [app-backend-2] /var/log/my_app/errors.log: File not found.
   Error: file not found or inaccessible: /var/log/my_app/errors.log
✓ [system-logs] /var/log/syslog: Analysis completed successfully.

Total: 3 logs analyzed (2 successful, 1 failed)
Analysis results saved to: report.json
Analysis completed successfully!
```

### JSON Report Format

When using the `--output` flag, results are saved in JSON format:

```json
[
  {
    "log_id": "web-server-1",
    "file_path": "/var/log/nginx/access.log",
    "status": "OK",
    "message": "Analysis completed successfully.",
    "error_details": ""
  },
  {
    "log_id": "app-backend-2",
    "file_path": "/var/log/my_app/errors.log",
    "status": "FAILURE",
    "message": "File not found.",
    "error_details": "file not found or inaccessible: /var/log/my_app/errors.log"
  }
]
```

## 🏗 Architecture

The project follows a clean, modular architecture:

```
loganalyzer/
├── main.go                # Application entry point
├── cmd/                   # CLI commands
│   ├── root.go            # Root command definition
│   └── analyze.go         # Analyze command implementation
├── internal/              # Internal packages
│   ├── config/            # Configuration handling
│   │   └── config.go
│   ├── analyzer/          # Log analysis and error handling
│   │   ├── analyzer.go    # Main analysis logic
│   │   └── errors.go      # Custom error types
│   └── reporter/          # Result reporting
│       └── reporter.go
├── examples/              # Example files
│   ├── config.json        # Sample configuration
│   └── test-config.json   # Test configuration
└── README.md              # This file
```

### Package Responsibilities

- **`cmd/`**: CLI command definitions using Cobra framework
- **`internal/config/`**: JSON configuration file loading and validation
- **`internal/analyzer/`**: Log file analysis, concurrency management, and custom error handling
- **`internal/reporter/`**: Result collection and output formatting

## 🔧 Key Technical Features

### Concurrency

- Uses goroutines for parallel log file processing
- Implements `sync.WaitGroup` for synchronization
- Channels for safe result collection between goroutines

### Error Handling

- Custom error types: `FileNotFoundError` and `ParseError`
- Proper error wrapping and unwrapping
- Uses `errors.Is()` and `errors.As()` for type-safe error handling

### Simulation Features

- Random processing time simulation (50-200ms per file)
- 10% chance of simulated parsing errors for testing
- File accessibility verification

## 🧪 Testing

### Quick Test

Use the provided test configuration:

```bash
# Build the application
go build -o loganalyzer main.go

# Run with test configuration
./loganalyzer analyze -c examples/test-config.json -o test-results.json
```

### Creating Test Files

```bash
# Create some test log files
echo "Test log content" > test1.log
echo "Another test log" > test2.log

# Create a test configuration
cat > test.json << EOF
[
  {
    "id": "test1",
    "path": "./test1.log",
    "type": "test log"
  },
  {
    "id": "test2", 
    "path": "./test2.log",
    "type": "test log"
  },
  {
    "id": "missing",
    "path": "./missing.log",
    "type": "test missing"
  }
]
EOF

# Run analysis
./loganalyzer analyze -c test.json -o results.json
```

## 📚 Learning Objectives Covered

This project demonstrates:

1. **Concurrency**: Goroutines, WaitGroups, and channels for parallel processing
2. **Error Handling**: Custom errors with `errors.Is()` and `errors.As()`
3. **CLI Development**: Cobra framework with subcommands and flags
4. **JSON Processing**: Import/export of JSON data
5. **Code Organization**: Modular design with internal packages
6. **Documentation**: Comprehensive code and usage documentation
