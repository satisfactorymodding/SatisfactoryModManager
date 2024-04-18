import { derived, get, writable } from 'svelte/store';

import { isLaunchingGame } from './generalStore';
import { ignoredUpdates } from './settingsStore';
import { binding, bindingTwoWay } from './wailsStoreBindings';

import { bytesToAppropriate, secondsToAppropriate } from '$lib/utils/dataFormats';
import { timeSeries } from '$lib/utils/timeSeries';
import { CheckForUpdates, GetInstallations, GetInstallationsMetadata, GetInvalidInstalls, GetModsEnabled, GetProfiles, GetRemoteInstallations, GetSelectedInstall, GetSelectedInstallLockfileMods, GetSelectedInstallProfileMods, GetSelectedProfile, SelectInstall, SetModsEnabled, SetProfile } from '$wailsjs/go/ficsitcli/ficsitCLI';
import { type cli, common, ficsitcli } from '$wailsjs/go/models';
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

export const progress = binding<ficsitcli.Progress | null>(null, { updateEvent: 'progress' });

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

export const progressTitle = derived(progress, ($progress) => {
  if (!$progress) return '';
  switch ($progress.action) {
    case ficsitcli.Action.SELECT_INSTALL: {
      const install = get(installsMetadata)[$progress.item.name];
      return `Selecting install ${install?.info?.branch} (${install?.info?.launcher}) - CL${install?.info?.version}`;
    }
    case ficsitcli.Action.SELECT_PROFILE:
      return `Selecting profile ${$progress.item.name}`;
    case ficsitcli.Action.TOGGLE_MODS:
      return `Turning mods ${$progress.item.name === 'true' ? 'on' : 'off'}`;
    case ficsitcli.Action.UPDATE:
      return 'Updating mods';
    case ficsitcli.Action.IMPORT_PROFILE:
      return `Importing profile ${$progress.item.name}`;
  }
});

const downloadTimeSeries = timeSeries(5000);
const extractTimeSeries = timeSeries(5000);

const STATS_UPDATE_INTERVAL = { speed: 0, eta: 500 };
const downloadStats = writable({ speed: 0, eta: 0 as number | undefined });
const extractStats = writable({ speed: 0, eta: 0 as number | undefined });
const lastStatsUpdate = { speed: 0, eta: 0 };

progress.subscribe(($progress) => {
  if (!$progress) {
    downloadTimeSeries.clear();
    extractTimeSeries.clear();
  } else {
    const { 
      download,
      extract,
    } = getTasksTotal($progress.tasks);
  
    downloadTimeSeries.addValue(download.current);
    extractTimeSeries.addValue(extract.current);

    const downloadSpeed = downloadTimeSeries.getDerivative() ?? 0;
    const extractSpeed = extractTimeSeries.getDerivative() ?? 0;
    const downloadETA = downloadSpeed !== 0 ? (download.total - download.current) / downloadSpeed : undefined;
    const extractETA = extractSpeed !== 0 ? (extract.total - extract.current) / extractSpeed : undefined;

    if (Date.now() - lastStatsUpdate.speed > STATS_UPDATE_INTERVAL.speed) {
      downloadStats.set({ speed: downloadSpeed, eta: get(downloadStats).eta });
      extractStats.set({ speed: extractSpeed, eta: get(extractStats).eta });
      lastStatsUpdate.speed = Date.now();
    }
    if (Date.now() - lastStatsUpdate.eta > STATS_UPDATE_INTERVAL.eta) {
      downloadStats.set({ speed: get(downloadStats).speed, eta: downloadETA });
      extractStats.set({ speed: get(extractStats).speed, eta: extractETA });
      lastStatsUpdate.eta = Date.now();
    }
  }
});

export const progressMessage = derived(progress, ($progress) => {
  if (!$progress) return '';
  
  const { 
    download,
    extract,
    downloadingMods,
    extractingMods,
  } = getTasksTotal($progress.tasks);

  if (download.total === 0 && extract.total === 0) {
    // Not downloading and not extracting, so nothing started yet
    const isRemoteInstall = get(installsMetadata)[$progress.item.name]?.info?.location === common.LocationType.REMOTE;
    switch ($progress.action) {
      case ficsitcli.Action.INSTALL:
      case ficsitcli.Action.ENABLE:
        return 'Finding the best version to install';
      case ficsitcli.Action.UNINSTALL:
      case ficsitcli.Action.DISABLE:
        return 'Checking for mods that are no longer needed';
      case ficsitcli.Action.SELECT_INSTALL:
      case ficsitcli.Action.SELECT_PROFILE:
      case ficsitcli.Action.IMPORT_PROFILE:
        return `Validating install... ${isRemoteInstall ? '(this may take a while for remote servers)' : ''}`;
      case ficsitcli.Action.UPDATE:
        return 'Updating...';
      case ficsitcli.Action.TOGGLE_MODS:
        if ($progress.item.name === 'true') {
          return 'Restoring mods...';
        } else {
          return 'Removing mods...';
        }
    }
  }
  
  downloadTimeSeries.addValue(download.current);
  extractTimeSeries.addValue(extract.current);

  if (download.current !== download.total) {
    // Downloading something, prioritize that
    const completeMods = downloadingMods.filter((m) => m.complete);
    const { speed, eta } = get(downloadStats);
    return `Downloading \
            ${completeMods.length}/${downloadingMods.length} mods: \
            ${bytesToAppropriate(download.current)}/${bytesToAppropriate(download.total)}, \
            ${bytesToAppropriate(speed)}/s, \
            ${eta !== undefined ? (eta !== 0 ? secondsToAppropriate(eta) : 'soon™') : '...'}`;
  }
  // Not downloading anything
  const completeMods = extractingMods.filter((m) => m.complete);
  const { speed, eta } = get(extractStats);
  return `Extracting \
          ${completeMods.length}/${extractingMods.length} mods: \
          ${bytesToAppropriate(extract.current)}/${bytesToAppropriate(extract.total)}, \
          ${bytesToAppropriate(speed)}/s, \
          ${eta !== undefined ? (eta !== 0 ? secondsToAppropriate(eta) : 'soon™') : '...'}`;
});

export const progressPercent = derived(progress, ($progress) => {
  if (!$progress) return undefined;
  const { download, extract } = getTasksTotal($progress.tasks);
  if (download.total === 0 && extract.total === 0) {
    // Not downloading and not extracting, so nothing started yet
    return undefined;
  }
  if (download.current !== download.total) {
    // Downloading something, prioritize that
    return download.current / download.total;
  }
  // Not downloading anything
  return extract.current / extract.total;
});

function getTasksTotal(tasks: ficsitcli.Progress['tasks']) {
  const download = { current: 0, total: 0 } as ficsitcli.ProgressTask;
  const extract = { current: 0, total: 0 } as ficsitcli.ProgressTask;
  const downloadingMods = [] as { name: string; version: string; complete: boolean }[];
  const extractingMods = [] as { name: string; version: string; complete: boolean }[];
  for (const [modVersionTask, status] of Object.entries(tasks)) {
    const [modVersion, task] = modVersionTask.split(':');
    const [name, version] = modVersion.split('@');
    if (task === 'download') {
      download.current += status.current;
      download.total += Math.max(status.current, status.total);
      downloadingMods.push({ name, version, complete: status.current === status.total && status.total !== 0 });
    } else if (task === 'extract') {
      extract.current += status.current;
      extract.total += Math.max(status.current, status.total);
      extractingMods.push({ name, version, complete: status.current === status.total && status.total !== 0 });
    }
  }
  return { download, extract, downloadingMods, extractingMods };
}
