# Project specification

## Purpose and current scope

A desktop and browser Go/Ebitengine rover prototype with a scrolling background, horizontal movement and repeated jump boosts. Obstacles, scoring, fuel, a jump ceiling and win/loss objectives are future work. The current game starts immediately.

## Architecture

- main.go: command-line options and native game loop.
- game/: Ebitengine adapter, drawing, keyboard/fullscreen/resize handling.
- internal/simulation: validated initial configuration, input state, motion, animation frame and tick timing; no Ebitengine or device dependencies.
- assets/: embedded PNG/WAV files and graphics/font/audio loading. Sounds are loaded but no gameplay path currently calls PlaySound.
- docs/: specification, coding and testing guidance.
- .github/workflows/build.yml: native/Wasm verification and headless simulation coverage.

Each game owns one simulation state. The adapter forwards left/right held flags and a just-pressed jump flag on every update, then advances simulation. Rendering reads state and keeps the rover horizontally centered while the background scrolls. Asset attribution remains in README.

## Configuration

Defaults: viewport 800x600 pixels, nominal speed 30, debug off. Startup dimensions must be 1..8192 pixels and speed 0..60; invalid settings return an error before assets/window initialization. Zero means pause. Previously out-of-range speeds were silently clamped; they now produce a useful startup error.

## Preserved behavior

Initial position is (width/2,100), animation frame 1, facing right. Left wins when both horizontal keys are held. Horizontal movement is 15 world pixels per simulation step; release stops horizontal motion. Animation advances on steps with nonzero velocity.

Jump adds -15 to current vertical velocity and resets jump time to 1, including in midair. Each physics step uses current y/velocity with the accumulated jump time: y becomes round(y + velocity*time + time*time), velocity increases by 2*time, then time increments. Landing clamps y to height-75 and clears jump time and vertical velocity. Resize changes this ground level. These historical formulas are regression-tested; this is not a redesign of jump physics.

The original scheduler steps only when its counter is strictly greater than TPS/speed, then resets. Speed is nominal, not an exact physics-step rate. Pause stops motion but still accepts direction/jump input and accumulates the tick counter; resume can move immediately. These semantics are preserved.

## Follow-up

Issue #2 tracks a browser launcher, web controls for all command-line options, Cloudflare deployment and landing-page integration. A Wasm build alone does not complete it. Obstacles, score, fuel and physics improvements remain separate gameplay decisions.

## Web implementation

The browser version uses the same simulation and drawing code with platform-specific query arguments, ready reporting and on-screen input plus fixed logical viewport dimensions. The page launches immediately. All desktop flag values are exposed in validated settings; restart creates a fresh isolated iframe runtime. Rendering scales to the page while preserving the logical configuration. Browser Q is ignored. See BROWSER.md and DEPLOYMENT.md for implemented hosting and loader details; production verification and landing-page publication complete issue #2.
