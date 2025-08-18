# Azure PIM Tool

This tool helps you manage Azure Privileged Identity Management (PIM) roles through a command-line interface.

### One-Step Installation

Copy and paste the following into your terminal to set up everything automatically:

```bash
cat > setup.sh << 'EOF'
#!/bin/bash

if [ -d "Azure-PIM-tool" ]; then
    echo "Repository already exists locally"
else
    echo "Cloning Repository"
    git clone "https://github.com/jbrusegaard/Azure-PIM-tool.git"
fi

if ! command -v go &> /dev/null; then
    echo "Installing go"
    brew install go
else
    echo "Go is already installed"
fi

cd Azure-PIM-tool

echo "Downloading Go Dependencies"
go mod download

echo "Installing Playwright..."
go run github.com/playwright-community/playwright-go/cmd/playwright@latest install

echo "Building application"
go build -o azure-pim-tool

echo "Setup complete! The application has been built as 'azure-pim-tool'"
EOF

chmod +x setup.sh
./setup.sh
```

### What This Does

The setup script will:

- Clone the repository if needed
- Install Go if not already installed
- Download Go dependencies
- Install Playwright
- Build the application

## Usage

### List Available Roles

To list all available PIM roles you can activate:

```bash
go run main.go list
```

### Activate a Role

You can activate one or multiple PIM roles simultaneously:

```bash
go run main.go activate [role_name] [role_name] -r "reason for activation"
```

Replace each `[role_name]` with the name of the role you want to activate and provide a reason for the activation after the `-r` flag.

Example:

One role

```bash
go run main.go activate "Global_Admin" -r "Needed to complete task"
```

Multiple roles

```bash
go run main.go activate "Global_Admin" "Security_Admin" -r "Needed to complete task"
```
