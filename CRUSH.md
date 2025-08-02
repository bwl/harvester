# CRUSH.md

Repo: Go (1.23+), module "bubbleRouge". CLI TUI using charmbracelet/bubbletea + lipgloss. No tests yet.

Build/Run
- Build: go build ./...
- Run: go run ./cmd/game
- Clean/cache: go clean -testcache

Lint/Format/Static analysis
- Format: go fmt ./...
- Imports: goimports -w .  (install: go install golang.org/x/tools/cmd/goimports@latest)
- Vet: go vet ./...
- Staticcheck (optional): staticcheck ./...  (install: go install honnef.co/go/tools/cmd/staticcheck@latest)
- Lint (optional): golangci-lint run  (install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)

Testing
- All tests: go test ./...
- Single package: go test ./path/to/pkg -v
- Single test: go test ./path/to/pkg -run ^TestName$ -v
- Benchmark: go test ./path/to/pkg -bench . -benchmem

Code style
- Formatting via go fmt and goimports; group imports: std, third-party, local. No unused imports.
- Types: prefer explicit types for exported APIs; use concrete types locally; avoid interface{}; keep zero-values idiomatic.
- Naming: CamelCase; exported identifiers start with caps and have doc-ready names; errors as err; receivers are short (m, t, p).
- Errors: return wrapped errors with fmt.Errorf("...: %w", err); never panic in library-like code; print to stderr for CLI failures.
- Control flow: keep functions small; avoid deep nesting; prefer early returns.
- Concurrency: use context when adding goroutines or timeouts; avoid data races; guard shared state.
- Logs/IO: do not log secrets; for CLI, write errors to os.Stderr and exit non-zero.

Project notes
- Entry point: main.go; uses bubbletea NewProgram; deterministic runs may need rand.Seed control.
- No Cursor or Copilot rules detected.

Agent tips
- Before running tools, ensure goimports/golangci-lint exist or gate usage. Use go env GOPATH to resolve tool installs when needed.