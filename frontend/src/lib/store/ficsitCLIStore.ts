import { writable, derived } from 'svelte/store';

import { binding, bindingTwoWay } from './wailsStoreBindings';
import { isLaunchingGame } from './generalStore';

import { cli, ficsitcli } from '$wailsjs/go/models';
import { CheckForUpdates, GetInstallationsInfo, GetInvalidInstalls, GetProfiles, GetSelectedInstall, GetSelectedProfile, SelectInstall, SetProfile, GetModsEnabled, SetModsEnabled, GetSelectedInstallProfileMods, GetSelectedInstallLockfileMods } from '$wailsjs/go/ficsitcli/FicsitCLI';
import { GetFavoriteMods } from '$wailsjs/go/bindings/Settings';

export const invalidInstalls = binding([], { initialGet: GetInvalidInstalls });

export const installs = binding([], { initialGet: GetInstallationsInfo, updateEvent: 'installs' });
export const selectedInstallPath = bindingTwoWay(null, { initialGet: () => GetSelectedInstall().then((i) => i?.path ?? null), updateEvent: 'selectedInstall' }, { updateFunction: SelectInstall });
export const selectedInstall = derived([installs, selectedInstallPath], ([$installs, $selectedInstallPath]) => {
  return $installs.find((i) => i.path === $selectedInstallPath) ?? null;
});

export const profiles = binding([], { initialGet: GetProfiles, updateEvent: 'profiles' });
export const selectedProfile = bindingTwoWay(null, { initialGet: GetSelectedProfile, updateEvent: 'selectedProfile', allowNull: false }, { updateFunction: SetProfile });

export const modsEnabled = bindingTwoWay(true, { initialGet: GetModsEnabled, updateEvent: 'modsEnabled', allowNull: false }, { updateFunction: SetModsEnabled });

export type ProfileMods = { [name: string]: cli.ProfileMod };

export const manifestMods = binding<ProfileMods>({}, { initialGet: GetSelectedInstallProfileMods, updateEvent: 'manifestMods', allowNull: false });

export interface LockedMod {
  version: string;
  hash: string;
  link: string;
  dependencies: { [id: string]: string };
}

export const lockfileMods = binding({}, { initialGet: GetSelectedInstallLockfileMods, updateEvent: 'lockfileMods', allowNull: false });

export interface Progress {
  item: string;
  progress: number;
  message: string;
}

export const progress = binding<Progress | null>(null, { updateEvent: 'progress' });

export const favoriteMods = binding<string[]>([], { updateEvent: 'favoriteMods', initialGet: GetFavoriteMods });

export const isGameRunning = binding(false, { updateEvent: 'isGameRunning', allowNull: false });

export const canModify = derived([isGameRunning, progress, isLaunchingGame], ([$isGameRunning, $progress, $isLaunchingGame]) => {
  return !$isGameRunning && !$progress && !$isLaunchingGame;
});

export const updates = writable<ficsitcli.Update[]>([]);
export const updateCheckInProgress = writable(false);

export async function checkForUpdates() {
  updateCheckInProgress.set(true);
  const result = await CheckForUpdates();
  updateCheckInProgress.set(false);
  if(result instanceof Error) {
    throw result;
  }
  updates.set(result ?? []);
}

setInterval(checkForUpdates, 1000 * 60 * 5); // Check for updates every 5 minutes
