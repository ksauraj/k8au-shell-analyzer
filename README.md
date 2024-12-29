# K8au Shell Analyzer

🚀 An interactive TUI tool to analyze your shell history and provide insights about your command-line usage patterns.

## Table of Contents
- [Features](#features)
- [Installation](#installation)
  - [Pre-built Binaries](#pre-built-binaries)
  - [Quick Install Script](#quick-install-script)
  - [Manual Installation](#manual-installation)
  - [Package Managers](#package-managers)
- [Configuration](#configuration)
- [Usage](#usage)
- [Development](#development)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Features
- 📊 Shell history analysis
- 🛠️ Tech stack detection
- 📈 Productivity metrics
- 🔄 Work pattern analysis
- 🎯 Tool usage statistics
- 🎬 Year-in-review style wrap-up

## Installation

### Pre-built Binaries

Download the latest release for your platform:

| Platform | Architecture | Download Link |
|----------|-------------|---------------|
| Linux    | amd64       | [Download](https://github.com/ksauraj/k8au-shell-analyzer/releases/latest/download/k8au-shell-analyser-linux-amd64) |
| Linux    | arm64       | [Download](https://github.com/ksauraj/k8au-shell-analyzer/releases/latest/download/k8au-shell-analyser-linux-arm64) |
| macOS    | amd64       | [Download](https://github.com/ksauraj/k8au-shell-analyzer/releases/latest/download/k8au-shell-analyser-darwin-amd64) |
| macOS    | arm64       | [Download](https://github.com/ksauraj/k8au-shell-analyzer/releases/latest/download/k8au-shell-analyser-darwin-arm64) |
| Windows  | amd64       | [Download](https://github.com/ksauraj/k8au-shell-analyzer/releases/latest/download/k8au-shell-analyser-windows-amd64.exe) |
| Windows  | arm64       | [Download](https://github.com/ksauraj/k8au-shell-analyzer/releases/latest/download/k8au-shell-analyser-windows-arm64.exe) |

### Quick Install Script

#### Linux/macOS (One-line installer)
```bash
curl -L https://raw.githubusercontent.com/ksauraj/k8au-shell-analyzer/master/setup.sh | bash
```

#### Using wget
```bash
wget -qO - https://raw.githubusercontent.com/ksauraj/k8au-shell-analyzer/master/setup.sh | bash
```

### Manual Installation

```bash
# Linux/macOS
wget https://github.com/ksauraj/k8au-shell-analyzer/releases/latest/download/k8au-shell-analyser-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)
chmod +x k8au-shell-analyser-*
./k8au-shell-analyser-*

# Windows PowerShell
Invoke-WebRequest -Uri "https://github.com/ksauraj/k8au-shell-analyzer/releases/latest/download/k8au-shell-analyser-windows-amd64.exe" -OutFile "k8au-shell-analyser.exe"
```

## Configuration

### Build from Source

Requirements:
- Go 1.20 or higher
- Gemini API Key

```bash
# Clone repository
git clone https://github.com/ksauraj/k8au-shell-analyzer.git
cd k8au-shell-analyzer

# Build with API key
make build GEMINI_API_KEY=your_api_key_here

# Or using go build directly
go build -ldflags "-X github.com/ksauraj/k8au-shell-analyzer/internal/gemini.apiKey=YOUR_API_KEY" ./cmd/k8au-shell-analyzer
```

## Usage

### Basic Usage
```bash
./k8au-shell-analyser
```

### Navigation Keys
| Key           | Action                |
|---------------|----------------------|
| `Tab`         | Switch between views |
| `←/→`         | Navigate slides      |
| `q`           | Quit application     |

### Available Views
1. **Overview**: General statistics
2. **Tech Profile**: Technical expertise analysis
3. **Work Patterns**: Productivity patterns
4. **Tool Usage**: Developer tools usage
5. **Wrapped**: Year-in-review summary

## Development

### Setup Development Environment
```bash
# Clone repository
git clone https://github.com/ksauraj/k8au-shell-analyzer.git
cd k8au-shell-analyzer

# Install dependencies
go mod download

# Run tests
go test ./...

# Run with hot reload (using air)
air
```

## Troubleshooting

### Common Issues

1. **Permission Denied**
```bash
chmod +x k8au-shell-analyser
```

2. **Binary Not Found**
```bash
export PATH=$PATH:$(pwd)
```

3. **API Key Issues**
Ensure you're building with the correct Gemini API key:
```bash
make build GEMINI_API_KEY=your_api_key_here
```

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Author

**Ksauraj** - [GitHub](https://github.com/ksauraj)

Project Link: [https://github.com/ksauraj/k8au-shell-analyzer](https://github.com/ksauraj/k8au-shell-analyzer)
