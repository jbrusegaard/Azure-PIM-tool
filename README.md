# Azure PIM Tool

A Go-based command-line tool that automates Azure Privileged Identity Management (PIM) role activation using Playwright for browser automation and token acquisition.

## Overview

This tool simplifies the process of activating Azure PIM roles by automating the authentication flow and role activation process. It uses Playwright to handle browser interactions and obtain Azure authentication tokens, then activates the specified PIM roles with configurable duration and reason.

## Features

- **Automated Authentication**: Uses Playwright to automate Azure login and token acquisition
- **PIM Role Activation**: Activate PIM roles for groups, resources, or specific roles
- **Flexible Duration**: Configurable activation duration (default: 8 hours)
- **Reason Tracking**: Optional reason logging for audit purposes
- **CLI Interface**: Clean command-line interface built with Cobra
- **Cross-Platform**: Works on Windows, macOS, and Linux

## Prerequisites

- Go 1.24.6 or higher
- Playwright browsers (automatically installed with the tool)
- Azure account with PIM access
- Appropriate permissions to activate PIM roles

## Installation

### From Source

1. Clone the repository:
```bash
git clone https://github.com/yourusername/Azure-PIM-tool.git
cd Azure-PIM-tool
```

2. Install dependencies:
```bash
go mod download
```

3. Install Playwright browsers:
```bash
go run github.com/playwright-community/playwright-go/cmd/playwright@latest install
```

4. Build the tool:
```bash
go build -o azure-pim-tool
```

### Binary Installation

Download the latest release binary for your platform from the releases page.

## Usage

### Basic Commands

```bash
# Activate a PIM role
./azure-pim-tool activate <type> <filter> [flags]

# List available options
./azure-pim-tool list

# Get help
./azure-pim-tool --help
```

### Examples

```bash
# Activate a group role for 4 hours
./azure-pim-tool activate group "MyGroup" --duration 4 --reason "Emergency access"

# Activate a resource role
./azure-pim-tool activate resource "MyResource" --duration 12 --reason "Maintenance window"

# Activate a specific role
./azure-pim-tool activate role "Contributor" --duration 8 --reason "Project work"
```

### Command Flags

- `--duration, -d`: Duration of activation in hours (default: 8)
- `--reason, -r`: Reason for activation (optional)
- `--help, -h`: Show help for the command

## Project Structure

```
Azure-PIM-tool/
├── cmd/                    # Command definitions
│   ├── activateCMD.go     # PIM activation command
│   ├── list.go           # List command
│   └── root.go           # Root command and CLI setup
├── constants/             # Application constants
│   └── constants.go
├── src/                   # Core functionality
│   ├── activate.go       # PIM activation logic
│   ├── fileHandler.go    # File operations
│   ├── pimUtils.go       # PIM utility functions
│   ├── settings.go       # Configuration management
│   └── Token.go          # Token handling
├── main.go               # Application entry point
├── go.mod                # Go module dependencies
└── go.sum                # Go module checksums
```

## Configuration

The tool uses configuration files to store settings and credentials. Configuration files are typically stored in the user's home directory.

### Environment Variables

- `AZURE_PIM_CONFIG_PATH`: Path to configuration file
- `AZURE_PIM_LOG_LEVEL`: Logging level (debug, info, warn, error)

## Dependencies

- **github.com/spf13/cobra**: CLI framework
- **github.com/playwright-community/playwright-go**: Browser automation
- **github.com/charmbracelet/log**: Structured logging

## Development

### Building

```bash
# Development build
go build -o azure-pim-tool

# Production build
go build -ldflags="-s -w" -o azure-pim-tool
```

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Code Style

This project follows Go standard formatting and linting practices. Use `gofmt` and `golint` to ensure code quality.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Security Considerations

- **Token Storage**: Azure tokens are stored securely and should not be shared
- **Permissions**: Only activate PIM roles that you have legitimate need for
- **Audit Logging**: All activations are logged for compliance purposes
- **Network Security**: Ensure secure network connections when using the tool

## Troubleshooting

### Common Issues

1. **Playwright Browser Not Found**: Run `go run github.com/playwright-community/playwright-go/cmd/playwright@latest install`
2. **Authentication Failed**: Verify Azure credentials and permissions
3. **PIM Activation Failed**: Check if the role is eligible for activation

### Debug Mode

Enable debug logging to troubleshoot issues:

```bash
export AZURE_PIM_LOG_LEVEL=debug
./azure-pim-tool activate <type> <filter>
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions:
- Open an issue on GitHub
- Check the troubleshooting section
- Review Azure PIM documentation

## Roadmap

- [ ] Support for conditional access policies
- [ ] Bulk role activation
- [ ] Integration with Azure DevOps pipelines
- [ ] Enhanced audit logging
- [ ] Configuration file validation
- [ ] Unit and integration tests

---

**Note**: This is a placeholder README. Please update with actual project details, screenshots, and specific configuration examples before publishing. 