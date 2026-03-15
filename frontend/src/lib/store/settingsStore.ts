import { binding, bindingTwoWay, bindingTwoWayNoExcept } from './wailsStoreBindings';

import type { LaunchButtonType, ViewType } from '$lib/wailsTypesExtensions';
import { GetVersion } from '$wailsjs/go/app/app';
import { GetOffline, SetOffline } from '$wailsjs/go/ficsitcli/ficsitCLI';
import { settings } from '$wailsjs/go/models';
import {
  GetCacheDir,
  GetDebug,
  GetIgnoredUpdates,
  GetKonami,
  GetLanguage,
  GetLaunchButton,
  GetModFiltersTagSearchMode,
  GetProxy,
  GetQueueAutoStart,
  GetRestoreWindowPosition,
  GetStartView,
  GetUpdateCheckMode,
  GetViewedAnnouncements,
  SetCacheDir,
  SetDebug,
  SetKonami,
  SetLanguage,
  SetLaunchButton,
  SetModFiltersTagSearchMode,
  SetProxy,
  SetQueueAutoStart, SetRestoreWindowPosition,
  SetStartView,
  SetUpdateCheckMode,
} from '$wailsjs/go/settings/settings';

export const startView = bindingTwoWayNoExcept<ViewType | null>(null, { initialGet: GetStartView }, { updateFunction: SetStartView });

export const tagSearchMode = bindingTwoWayNoExcept<settings.TagSearchMode>(
  settings.TagSearchMode.ANY,
  { initialGet: GetModFiltersTagSearchMode },
  { updateFunction: SetModFiltersTagSearchMode },
);

export const saveWindowPosition = bindingTwoWayNoExcept(true, { initialGet: GetRestoreWindowPosition }, { updateFunction: SetRestoreWindowPosition });

export const konami = bindingTwoWayNoExcept(false, { initialGet: GetKonami }, { updateFunction: SetKonami });

export const launchButton = bindingTwoWayNoExcept<LaunchButtonType>('normal', { initialGet: () => GetLaunchButton().then((l) => l as LaunchButtonType) }, { updateFunction: SetLaunchButton });

export const queueAutoStart = bindingTwoWayNoExcept(true, { initialGet: GetQueueAutoStart }, { updateFunction: SetQueueAutoStart });

export const offline = bindingTwoWayNoExcept<boolean>(false, { initialGet: GetOffline }, { updateFunction: SetOffline });

export const proxy = bindingTwoWayNoExcept<string>('', { initialGet: GetProxy }, { updateFunction: SetProxy });

export const updateCheckMode = bindingTwoWayNoExcept<'launch'|'exit'|'ask'>('launch', { initialGet: GetUpdateCheckMode }, { updateFunction: SetUpdateCheckMode });

export const viewedAnnouncements = binding<string[]>([], { initialGet: GetViewedAnnouncements, updateEvent: 'viewedAnnouncements' });

export const ignoredUpdates = binding<Record<string, string[]>>({}, { initialGet: GetIgnoredUpdates, updateEvent: 'ignoredUpdates' });

export const cacheDir = bindingTwoWay<string, null>(null, { initialGet: GetCacheDir, updateEvent: 'cacheDir' }, { updateFunction: SetCacheDir });

export const version = binding<string>('0.0.0', { initialGet: GetVersion });

export const debug = bindingTwoWayNoExcept<boolean>(false, { initialGet: GetDebug }, { updateFunction: SetDebug });

export const language = bindingTwoWayNoExcept<string>('en', { initialGet: () => GetLanguage().then((l) => l ? l : 'en'), allowNull: false }, { updateFunction: SetLanguage });
