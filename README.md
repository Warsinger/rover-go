# Rover

A side scrolling game where you drive a rover, jump and avoid obstacles, and other fun things.

## Gameplay

### Movement

* Left and Right arrow keys to move
* Spacebar to jump, you can jump while in the air to jump hire

## Command line options

* `-speed N` Set the general game speed 1 - 60 ticks per second. Default is 30.
* `-debug` Print some debug information t the screen and the console.
* `-width N` Set the width of the screen. Default is 800.
* `-height N` Set the height of the screen. Default is 600.

## Graphics

Graphics by Matt. Nathan created the Lego builds that were used for the assets.

![Rover](assets/images/rover1.png)

## TODO

* Implement max jump height
  * fuel that runs out and requires recharge with item pickup or time
* Graphics
  * ~~fire rockets on jump~~
  * ~~flip vehicle direction~~
  * ~~scale background~~
* Obstacles
* Score
* Objective of the game?

## Build and dependency verification

Requires Go 1.27.3 or newer. Go's automatic toolchain selection can download the version declared in `go.mod`.

```sh
go mod download
go mod verify
go vet ./...
go test ./...
go build -o build/game .
```

Ebitengine's native Linux build requires its graphics development libraries and a display to run. See the [Ebitengine installation guide](https://ebitengine.org/en/documents/install.html).

A WebAssembly compilation check is also available:

```sh
GOOS=js GOARCH=wasm go build -o build/game.wasm .
```

In PowerShell, set `$env:GOOS = "js"` and `$env:GOARCH = "wasm"` before the build, then remove those environment variables afterward. This checks compatibility; the browser launcher and web controls are tracked separately in the web-version issue. There are currently no automated simulation tests; those are tracked in the testing issue.
