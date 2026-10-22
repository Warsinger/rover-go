import {startLauncher} from './shared/launcher.mjs';
import {validateOptions} from './options.mjs';
startLauncher({
  validateOptions: values => validateOptions(values, 'rover'),
  readOptions: form => ({...Object.fromEntries(new FormData(form)), debug: form.elements.debug.checked}),
  readyMessage: 'Ready. Use arrow keys and Space, or the on-screen controls.',
  heldKeys: ['left', 'right']
});
