import { binding, bindingTwoWay, bindingTwoWayNoExcept } from './wailsStoreBindings';

import type { LaunchButtonType, ViewType } from '$lib/wailsTypesExtensions';
import { GetOffline, SetOffline } from '$wailsjs/go/ficsitcli/ficsitCLI';
import { GetStartView, SetStartView, GetKonami, SetKonami, GetLaunchButton, SetLaunchButton, GetQueueAutoStart, SetQueueAutoStart, GetUpdateCheckMode, SetUpdateCheckMode, GetViewedAnnouncements, GetIgnoredUpdates, GetCacheDir, SetCacheDir } from '$wailsjs/go/settings/settings';
import { GetVersion } from '$lib/generated/wailsjs/go/app/app';

export const startView = bindingTwoWayNoExcept<ViewType | null>(null, { initialGet: GetStartView }, { updateFunction: SetStartView });

export const konami = bindingTwoWayNoExcept(false, { initialGet: GetKonami }, { updateFunction: SetKonami });

export const launchButton = bindingTwoWayNoExcept<LaunchButtonType>('normal', { initialGet: () => GetLaunchButton().then((l) => l as LaunchButtonType) } , { updateFunction: SetLaunchButton });

export const queueAutoStart = bindingTwoWayNoExcept(true, { initialGet: GetQueueAutoStart }, { updateFunction: SetQueueAutoStart });

export const offline = bindingTwoWayNoExcept<boolean>(false, { initialGet: GetOffline }, { updateFunction: SetOffline });

export const updateCheckMode = bindingTwoWayNoExcept<'launch'|'exit'|'ask'>('launch', { initialGet: GetUpdateCheckMode }, { updateFunction: SetUpdateCheckMode });

export const viewedAnnouncements = binding<string[]>([], { initialGet: GetViewedAnnouncements, updateEvent: 'viewedAnnouncements' });

export const ignoredUpdates = binding<Record<string, string[]>>({}, { initialGet: GetIgnoredUpdates, updateEvent: 'ignoredUpdates' });

export const cacheDir = bindingTwoWay<string, null>(null, { initialGet: GetCacheDir, updateEvent: 'cacheDir' }, { updateFunction: SetCacheDir });

export const version = binding<string>('0.0.0', { initialGet: GetVersion });
