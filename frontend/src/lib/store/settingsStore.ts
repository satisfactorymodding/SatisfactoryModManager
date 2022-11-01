import type { LaunchButtonType, ViewType } from '$lib/wailsTypesExtensions';
import { GetStartView, SetStartView, GetKonami, SetKonami, GetLaunchButton, SetLaunchButton } from '$wailsjs/go/bindings/Settings';
import { writableBindingSync } from './wailsStoreBindings';

export const startView = writableBindingSync<ViewType | null>(null, { 
  initialGet: GetStartView,
  updateFunction: (value) => value ? SetStartView(value) : Promise.resolve(),
});

export const konami = writableBindingSync(false, { initialGet: GetKonami, updateFunction: SetKonami });

export const launchButton = writableBindingSync<LaunchButtonType>('normal', { initialGet: GetLaunchButton, updateFunction: SetLaunchButton });