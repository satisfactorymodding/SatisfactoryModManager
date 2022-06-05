import { get, readable, writable } from 'svelte/store';
import { cli, bindings } from '$wailsjs/go/models';
import { EventsOff, EventsOn } from '$wailsjs/runtime/runtime';
import { GetInstallationsInfo, GetProfiles, SelectInstall, SetProfile } from '$wailsjs/go/bindings/FicsitCLI';

function readableBinding<T>(options: {
  defaultValue: T,
  updateEvent?: string,
  allowNull?: boolean,
  initialGet?: () => Promise<T>
}) {
  const { defaultValue, updateEvent, allowNull, initialGet } = {
    allowNull: true,
    ...options
  };
  return readable(defaultValue, typeof window !== 'undefined' ? ((set) => {
    const setData = (data) => {
      if(data === null && !allowNull) {
        set(defaultValue);
      } else {
        set(data);
      }
    };

    EventsOn(updateEvent, setData);

    if(initialGet) {
      initialGet().then(setData);
    }

    return () => {
      EventsOff(updateEvent);
    };
  }) : undefined);
}

export const installs = readableBinding<bindings.InstallationInfo[]>({ defaultValue: [], initialGet: GetInstallationsInfo});
export const selectedInstall = writable(null as bindings.InstallationInfo | null);

export const profiles = readableBinding<string[]>({ defaultValue: [], initialGet: GetProfiles});
export const selectedProfile = writable(null as string | null);

const installsLoadDone = installs.subscribe((i) => {
  if(i.length > 0) {
    selectedInstall.set(i[0]);
    installsLoadDone();
  }
});

selectedInstall.subscribe((i) => {
  const path = i?.info?.path;
  if(path) {
    SelectInstall(path);
    if(i.installation) {
      selectedProfile.set(i.installation.profile);
    }
  }
});

selectedProfile.subscribe((p) => {
  if(p) {
    SetProfile(p);
    const install = get(selectedInstall);
    if(install && install.installation) {
      install.installation.profile = p;
    }
  }
});

export type ProfileMods = Dictionary<string, cli.ProfileMod>;

export const manifestMods = readableBinding<ProfileMods>({ defaultValue: {}, allowNull: false, updateEvent: 'manifestMods'}); // Event will be

export interface LockedMod {
  version: string;
  hash: string;
  link: string;
  dependencies: Dictionary<string, string>;
}

export type LockFile = Dictionary<string, LockedMod>;

export const lockfileMods = readableBinding<LockFile>({ defaultValue: {}, allowNull: false, updateEvent: 'lockfileMods'});

export const progress = readableBinding<bindings.Progress | null>({ defaultValue: null, updateEvent: 'progress'});
