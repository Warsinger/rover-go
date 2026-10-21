# Testing

Run from the repository root with the Go version declared in go.mod:

```sh
go test ./...
go vet ./...
go mod verify
go mod tidy -diff
go build -o build/game .
go test -race -count=1 -coverprofile=build/coverage.out ./internal/simulation
go tool cover -func=build/coverage.out
go tool cover -html=build/coverage.out -o build/coverage.html
```

Create build/ first if it does not exist. In PowerShell, quote arguments containing equals signs or dots, for example '-coverprofile=build/coverage.out'. The race detector requires a supported platform and C toolchain; regular go test works without that additional toolchain for the pure simulation package.

The simulation package has no display, graphics, audio or filesystem dependency. CI runs its tests with the race detector on Linux and publishes raw/HTML coverage as an artifact and a summary. A Windows job also runs all package tests, vet, module verification and native/Wasm builds. Simulation coverage does not include the Ebitengine adapter, keyboard events, render output or asset loading.

WebAssembly compilation (POSIX shell):

```sh
GOOS=js GOARCH=wasm go build -o build/game.wasm .
```

PowerShell:

```powershell
$env:GOOS = 'js'
$env:GOARCH = 'wasm'
try { go build -o build/game.wasm . } finally { Remove-Item Env:GOOS; Remove-Item Env:GOARCH }
```

This is the low-level compilation check; the complete browser launcher and browser-specific manual checks are documented in BROWSER.md. Native Linux packages importing Ebitengine also require graphics libraries and a display; use Xvfb when testing those packages on a headless Linux host. See the Ebitengine installation guide linked in README.

## Automated behavior

Tests cover initial state, left/right precedence, stopping, sprite frame advancement, the strict tick threshold, pause with queued input, falling, the historical jump trajectory, repeated midair boosts, landing, ground changes on resize, independent instances and startup limits.

## Manual checklist

Run go run .; move in both directions and release the keys; jump and boost again midair; land; pause with minus and resume with equals; toggle D; toggle fullscreen with F and restore; resize smaller/larger; verify the background repeats and the rover/flame face correctly; quit with Q. Test invalid width/height/speed and confirm an error. Assets and audio startup still need a real device check. There are embedded sounds, but gameplay currently does not play them.


## Browser checks

Run `node --test scripts/*.test.mjs`, `go run ./cmd/web`, and `node scripts/package-pages.mjs <fresh-output-directory>`. CI uploads the complete Pages artifact after these pass. See BROWSER.md for the actual Chromium checks performed and remaining browser/touch coverage limits. Verify settings after restart, initial focus, graceful failed loading, gzip bundle startup, and that the preview URL appears on the PR and in the workflow summary.
