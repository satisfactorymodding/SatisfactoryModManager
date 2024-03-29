import { derived, writable } from 'svelte/store';

import { isLaunchingGame } from './generalStore';
import { ignoredUpdates } from './settingsStore';
import { binding, bindingTwoWay } from './wailsStoreBindings';

import { CheckForUpdates, GetInstallations, GetInstallationsMetadata, GetInvalidInstalls, GetModsEnabled, GetProfiles, GetRemoteInstallations, GetSelectedInstall, GetSelectedInstallLockfileMods, GetSelectedInstallProfileMods, GetSelectedProfile, SelectInstall, SetModsEnabled, SetProfile } from '$wailsjs/go/ficsitcli/ficsitCLI';
import { type cli, ficsitcli } from '$wailsjs/go/models';
import { GetFavoriteMods } from '$wailsjs/go/settings/settings';

export const invalidInstalls = binding([], { initialGet: GetInvalidInstalls });

export const installs = binding([], { initialGet: GetInstallations, updateEvent: 'installations', allowNull: false });
export const installsMetadata = binding({}, { initialGet: GetInstallationsMetadata, updateEvent: 'installationsMetadata', allowNull: false });
export const selectedInstall = bindingTwoWay(null, { initialGet: () => GetSelectedInstall().then((i) => i?.path ?? null), updateEvent: 'selectedInstallation' }, { updateFunction: SelectInstall });
export const selectedInstallMetadata = derived([installsMetadata, selectedInstall], ([$installsMetadata, $selectedInstallPath]) => {
  return $installsMetadata[$selectedInstallPath ?? '__invalid__install__'] ?? null;
});

export const remoteServers = binding([], { initialGet: () => GetRemoteInstallations(), updateEvent: 'remoteServers', allowNull: false });

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

export const canModify = derived([isGameRunning, progress, isLaunchingGame, installs, selectedInstallMetadata], ([$isGameRunning, $progress, $isLaunchingGame, $installs, $selectedInstallMetadata]) => {
  return !$isGameRunning && !$progress && !$isLaunchingGame && $installs.length > 0 && $selectedInstallMetadata?.state === ficsitcli.InstallState.VALID;
});

export const canChangeInstall = derived([isGameRunning, progress, isLaunchingGame, installs], ([$isGameRunning, $progress, $isLaunchingGame, $installs]) => {
  return !$isGameRunning && !$progress && !$isLaunchingGame && $installs.length > 0;
});

export const canInstallMods = derived([isGameRunning, isLaunchingGame, installs, selectedInstallMetadata], ([$isGameRunning, $isLaunchingGame, $installs, $selectedInstallMetadata]) => {
  return !$isGameRunning && !$isLaunchingGame && $installs.length > 0 && $selectedInstallMetadata?.state === ficsitcli.InstallState.VALID;
});

export const updates = writable<ficsitcli.Update[]>([]);
export const unignoredUpdates = derived([updates, ignoredUpdates], ([$updates, $ignoredUpdates]) => $updates.filter((u) => !$ignoredUpdates[u.item]?.includes(u.newVersion)));
export const updateCheckInProgress = writable(false);

export async function checkForUpdates() {
  updateCheckInProgress.set(true);
  try {
    const result = await CheckForUpdates();
    updates.set(result ?? []);
  } finally {
    updateCheckInProgress.set(false);
  }
}

setInterval(() => checkForUpdates().catch(console.error), 1000 * 60 * 5); // Check for updates every 5 minutes
