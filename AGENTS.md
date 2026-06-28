# AGENTS.md

## Project
- Go CLI module: `github.com/anton-fuji/gitviz`.
- `cmd/main.go` defines the `-add` and `-graph` flags.
- `internal/scan.go` finds Git repositories under a folder and stores paths in `~/.gitlocalstats`.
- `internal/stats.go` reads registered repositories, counts commits for an author email over the last 183 days, and prints a GitHub-style terminal graph.

## Commands
- Test: `go test ./...`
- Run scanner: `go run ./cmd -add /path/to/projects`
- Run graph: `go run ./cmd -graph user@example.com`

## Guidelines
- Keep changes small and idiomatic Go; run `gofmt` on edited Go files.
- Preserve the existing CLI flags and the `~/.gitlocalstats` file format unless explicitly changing behavior.
- Avoid tests that write to the real home directory; prefer temp files or pure helper coverage.
- Be careful with recursive scans and Git history traversal in tests; keep fixtures small.
