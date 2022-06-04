import { readable } from 'svelte/store';
import { cli, main } from '../../wailsjs/go/models';
import { EventsOff, EventsOn } from '../../wailsjs/runtime/runtime';

export type ProfileMods = Dictionary<string, cli.ProfileMod>;

export const manifestMods = readable({} as ProfileMods, typeof window !== 'undefined' ? ((set) => {
  EventsOn('manifestMods', (data) => {
    set(data ?? {});
  });

  return () => {
    EventsOff('manifestMods');
  };
}) : undefined);

export interface LockedMod {
  version: string;
  hash: string;
  link: string;
  dependencies: Dictionary<string, string>;
}

export type LockFile = Dictionary<string, LockedMod>;

export const lockfileMods = readable({} as LockFile, typeof window !== 'undefined' ? ((set) => {
  EventsOn('lockfileMods', (data) => {
    set(data ?? {});
  });

  return () => {
    EventsOff('lockfileMods');
  };
}) : undefined);

export const progress = readable(null as main.Progress | null, typeof window !== 'undefined' ? ((set) => {
  EventsOn('progress', (data) => {
    set(data);
  });

  return () => {
    EventsOff('progress');
  };
}) : undefined);