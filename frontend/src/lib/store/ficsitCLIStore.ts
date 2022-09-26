import { get, writable } from 'svelte/store';
import { cli, bindings } from '$wailsjs/go/models';
import { AddProfile, CheckForUpdates, DeleteProfile, GetInstallationsInfo, GetInvalidInstalls, GetProfiles, ImportProfile, RenameProfile, SelectInstall, SetProfile } from '$wailsjs/go/bindings/FicsitCLI';
import { GetFavouriteMods } from '$wailsjs/go/bindings/Settings';
import { readableBinding, writableBinding } from './wailsStoreBindings';
import { tick } from 'svelte';

export const invalidInstalls = readableBinding<(Error & {path?: string})[]>([], { initialGet: GetInvalidInstalls });

export const installs = readableBinding<bindings.InstallationInfo[]>([], { initialGet: GetInstallationsInfo });
export const selectedInstall = writable(null as bindings.InstallationInfo | null);

export const profiles = writableBinding<string[]>([], { initialGet: GetProfiles });
export const selectedProfile = writable(null as string | null);

Promise.all([installs.waitForInit, profiles.waitForInit]).then(() => {
  const i = get(installs);
  if(i.length > 0) {
    selectedInstall.set(get(installs)[0]);
  }
});

selectedInstall.subscribe((i) => {
  const path = i?.info?.path;
  if(path) {
    SelectInstall(path);
    if(i.installation) {
      selectedProfile.set(i.installation.profile);
    }
    checkForUpdates();
  }
});

selectedProfile.subscribe((p) => {
  if(p) {
    SetProfile(p);
    const install = get(selectedInstall);
    if(install && install.installation) {
      install.installation.profile = p;
    }
    checkForUpdates();
  }
});

export async function addProfile(name: string) {
  const err = await AddProfile(name);
  if(err) {
    throw err;
  }
  const newProfiles = get(profiles);
  if(!newProfiles.includes(name)) {
    newProfiles.push(name);
    profiles.set(newProfiles);
  }
}

export async function renameProfile(oldName: string, newName: string) {
  const err = await RenameProfile(oldName, newName);
  if(err) {
    throw err;
  }
  const newProfiles = get(profiles);
  if(newProfiles.includes(oldName)) {
    const idx = newProfiles.indexOf(oldName);
    newProfiles[idx] = newName;
    profiles.set(newProfiles);
  }
  get(installs).forEach((i) => { if(i.installation.profile === oldName) { i.installation.profile = newName; } });
  if(get(selectedProfile) === oldName) {
    selectedProfile.set(newName);
  }
}

export async function deleteProfile(name: string) {
  const err = await DeleteProfile(name);
  if(err) {
    throw err;
  }
  const newProfiles = get(profiles);
  if(newProfiles.includes(name)) {
    const idx = newProfiles.indexOf(name);
    newProfiles.splice(idx, 1);
    profiles.set(newProfiles);
  }
  get(installs).forEach((i) => { if(i.installation.profile === name) { i.installation.profile = 'Default'; } });
  if(get(selectedProfile) === name) {
    selectedProfile.set('Default');
  }
}

export async function importProfile(name: string, filepath: string) {
  const err = await ImportProfile(name, filepath);
  if(err) {
    throw err;
  }
  const newProfiles = get(profiles);
  if(!newProfiles.includes(name)) {
    newProfiles.push(name);
    profiles.set(newProfiles);
    tick().then(() => {
      selectedProfile.set(name);
    });
  }
}

export type ProfileMods = { [name: string]: cli.ProfileMod };

export const manifestMods = readableBinding<ProfileMods>({}, { allowNull: false, updateEvent: 'manifestMods'}); // Event will be

export interface LockedMod {
  version: string;
  hash: string;
  link: string;
  dependencies: { [id: string]: string };
}

export type LockFile = { [name: string]: LockedMod };

export const lockfileMods = readableBinding<LockFile>({}, { allowNull: false, updateEvent: 'lockfileMods'});

export const progress = readableBinding<bindings.Progress | null>(null, { updateEvent: 'progress'});

export const favouriteMods = readableBinding<string[]>([], { updateEvent: 'favouriteMods', initialGet: GetFavouriteMods});

export const updates = writable<bindings.Update[]>([]);
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