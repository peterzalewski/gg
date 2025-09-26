# CRUSH.md - Development Guide for gg

## Build/Test/Lint Commands
```bash
# Build the project
go build -o gg .

# Run all tests
go test ./...

# Run tests for specific package
go test ./model
go test ./cmd

# Run single test function
go test -run TestFunctionName ./package

# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Run the built binary
./gg --help
```

## Code Style Guidelines

- **Package structure**: Commands in `cmd/`, models in `model/`
- **Imports**: Standard library first, then third-party, then local packages
- **Error handling**: Return errors explicitly, use custom error variables (e.g., `ErrNotAGitRepository`)
- **Naming**: Use camelCase for functions/variables, PascalCase for exported types
- **Types**: Use struct embedding for composition (e.g., `Blob` embeds `*Object`)
- **Context**: Pass repository via context in cobra commands
- **Flags**: Use cobra's persistent flags for global options like `--work-tree`
- **File operations**: Always defer `Close()` after opening files
- **String building**: Use `strings.Builder` for efficient string concatenation
- **Regex**: Compile regex patterns as package-level variables
- **Options pattern**: Use functional options for constructors (e.g., `RepositoryOption`)
