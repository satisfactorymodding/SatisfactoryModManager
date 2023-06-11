import { writableBinding, writableBindingSync } from './wailsStoreBindings';

import type { LaunchButtonType, ViewType } from '$lib/wailsTypesExtensions';
import { GetOffline, SetOffline } from '$wailsjs/go/ficsitcli/FicsitCLI';
import { GetStartView, SetStartView, GetKonami, SetKonami, GetLaunchButton, SetLaunchButton, GetQueueAutoStart, SetQueueAutoStart, GetUpdateCheckMode, SetUpdateCheckMode, GetViewedAnnouncements } from '$wailsjs/go/bindings/Settings';

export const startView = writableBindingSync<ViewType | null>(null, { 
  initialGet: GetStartView,
  updateFunction: (value) => value ? SetStartView(value) : Promise.resolve(),
});

export const konami = writableBindingSync(false, { initialGet: GetKonami, updateFunction: SetKonami });

export const launchButton = writableBindingSync<LaunchButtonType>('normal', { initialGet: GetLaunchButton, updateFunction: SetLaunchButton });

export const queueAutoStart = writableBindingSync(true, { initialGet: GetQueueAutoStart, updateFunction: SetQueueAutoStart });

export const offline = writableBindingSync<boolean|null>(null, { initialGet: GetOffline, updateFunction: SetOffline });

export const updateCheckMode = writableBindingSync<'launch'|'exit'|'ask'>('launch', { initialGet: GetUpdateCheckMode, updateFunction: SetUpdateCheckMode });

export const viewedAnnouncements = writableBinding<string[]>([], { initialGet: GetViewedAnnouncements });