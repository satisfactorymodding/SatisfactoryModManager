import { GetStartView, SetStartView, GetKonami, SetKonami, GetLaunchButton, SetLaunchButton } from '$wailsjs/go/bindings/Settings';
import { writableBindingSync } from './wailsStoreBindings';

export type View = 'compact' | 'expanded';

export const startView = writableBindingSync<View | null>(null, { 
  initialGet: GetStartView,
  updateFunction: (value) => value ? SetStartView(value) : Promise.resolve(),
});

export const konami = writableBindingSync(false, { initialGet: GetKonami, updateFunction: SetKonami });

export type LaunchButton = 'normal' | 'cat' | 'button';

export const launchButton = writableBindingSync<LaunchButton>('normal', { initialGet: GetLaunchButton, updateFunction: SetLaunchButton });