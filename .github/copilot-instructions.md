<HighLevelDetails>
- **Project**: TypeGo is an embedded TypeScript runtime for Go, allowing users to script Go applications with TS logic.
- **Nature**: It uses `goja` (JS engine in Go) and `esbuild` (bundler) to provide a Node.js-like experience but as a single static binary.
- **Tech Stack**:
  - **Go**: v1.24+
  - **Runtime**: `dop251/goja`
  - **Bundler**: `evanw/esbuild`
  - **CLI**: `spf13/cobra`
</HighLevelDetails>

<StyleGuidelines>
**Adhere to the "Senior Janitor" Persona.**
Refer to `.github/ai-style.md` for the complete list of forbidden phrases, comment rules, and generating guidelines.

**Core Directives**:
1. **Be Laconic**: Output code directly. No "Here is the code".
2. **Comments**: Write "Why", not "What".
</StyleGuidelines>

<BuildInstructions>
**Always run these commands from the repository root.**

**Sequence**:
1.  **Setup**: `go mod download`
    - *Precondition*: Go 1.24+ installed.
2.  **Lint**: `golangci-lint run`
    - *Note*: Must be clean before commit.
3.  **Test**: `go test -v -race ./...`
    - *Mandatory*: Must pass with race detection enabled.
4.  **Build**: `go build -v -o typego.exe ./cmd/typego`
    - *Verification*: `./typego.exe run examples/01-hello-world.ts`

**Dev Mode**:
- `go run ./cmd/typego dev <file.ts>` for hot-reload.
</BuildInstructions>

<ProjectLayout>
**Major Elements**:
- **`cmd/typego/`**: Main CLI entry point. Configuration for flags/commands here.
- **`engine/`**: Core runtime wrapper around `goja`. Handles EventLoop.
- **`bridge/`**: Binding layer. `bridge/modules` contains stdlib implementations.
- **`compiler/`**: Esbuild wrapper. Handles TS->JS compilation and caching.

**Key Files**:
- `Makefile`: Task runner aliases.
- `go.mod`: Dependency definitions.
- `CONTRIBUTING.md`: Detailed contribution rules.

**Validation Strategy**:
- CI runs `go test -race ./...` and `golangci-lint`.
- Feature branches only. No direct commits to main.

**Search Policy**:
- Trust these instructions. Only search if you encounter an error not covered here.
</ProjectLayout>
