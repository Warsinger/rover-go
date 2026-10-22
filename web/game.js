import {startGoGame} from './shared/go-runtime.mjs';
startGoGame({
  wasmURL: new URL('game.wasm', import.meta.url),
  inputKeys: ['left', 'right', 'jump'],
  latchedKeys: ['jump'],
  preventKeys: ['Space', 'ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight']
});
