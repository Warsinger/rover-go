# Browser application

The page starts immediately with the same defaults as the desktop app. There is no start screen. All command-line options have labeled web controls. Apply / restart validates the settings and creates a fresh iframe/game instance; the rover returns to its initial position. The configured logical viewport is fixed while the page scales it to fit.

The game canvas receives focus after its first actual Go update. Use Focus game after interacting with settings, or click the game. Browser Q is ignored; desktop Q still exits. Fullscreen is available from the page button or the game F key when supported by the browser. Arrow keys move and Space jumps/boosts. On-screen buttons also support pointer, touch, Space and Enter input. Held movement is released on blur/visibility changes; jump is latched until one Go update consumes it.

## Local build

Run from the repository root:

```sh
npm ci --ignore-scripts
go run ./cmd/web -serve -port 8080
```

Requires Node 24+ for the pinned browser package and browser checks. Open http://localhost:8080/. Without -serve, the command writes build/web and exits. Serve over HTTP rather than opening HTML as a local file. The builder copies wasm_exec.js from the exact Go toolchain that compiled the binary. It replaces inherited GOOS/GOARCH values for the Wasm subprocess without changing the native caller.

Settings are passed as allowlisted query arguments and validated both in JavaScript and by Go before game setup. They are not saved between visits. Browser rules and native simulation use the same Go implementation. No multiplayer, persistence, editing interface or new physics is added.

## Shared browser shell

The package is installed from an exact web-shared Git commit tarball with a committed integrity lock. cmd/web copies its browser modules and styles into build/web/shared/; Pages packaging includes and hashes them with the app. Playing requires no cross-site JS or CSS requests. web/launcher.js supplies Rover's settings conversion, ready message and held controls; web/game.js selects input/latched keys and the Wasm URL; web/styles.css supplies color overrides. HTML content, options validation, footer links/attribution and Go behavior remain here.

To update, select a reviewed web-shared commit, update its tarball dependency and regenerate package-lock.json, run npm ci --ignore-scripts plus browser checks, and verify a preview. Shared releases reach production only through a reviewed consumer update/deployment. Never use mutable branch URLs or edit node_modules. Each Go build still supplies its own matching wasm_exec.js.

## Loading and security boundaries

The launcher receives same-origin messages only from its own iframe and current session token, discarding previous-run and foreign messages. It reports ready only after Go's first update; download/runtime failures remain visible with a retry path. Runtime failures do not become successful loading messages. A slow-load notice does not prevent late success. Fullscreen denial is recoverable.

The local loader accepts raw Wasm. The Pages packager generates gzip Wasm and the same loader recognizes gzip magic bytes, avoiding double decompression when a browser/server already decodes the response. A current browser with WebAssembly, WebGL and DecompressionStream is required for the hosted bundle.

## Verification

Run npm ci --ignore-scripts, then node --test scripts/*.test.mjs plus the Go checks in TESTING.md. JavaScript tests cover option conversion/limits, real HTML/runtime references, reproducible hashed assets, gzip byte integrity, asset limits, immutable caching, and inclusion of shared browser modules/styles. Loader success/failure, lifecycle and input message trust/release tests live in web-shared. Coverage of actual rendering remains a browser check.

Manually verified in the in-app Chromium browser: automatic startup, keyboard focus, arrow movement without an initial click, pointer movement and facing. Review preview deployment on other browsers and real touch devices before claiming a full support matrix.
