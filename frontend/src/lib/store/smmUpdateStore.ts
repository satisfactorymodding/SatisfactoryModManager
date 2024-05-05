import { compare } from 'semver';
import { derived, readable } from 'svelte/store';

import { binding } from './wailsStoreBindings';

import { progressStats } from '$lib/utils/progress';
import { PendingUpdate } from '$wailsjs/go/autoupdate/autoUpdate';
import type { autoupdate, utils } from '$wailsjs/go/models';
import { EventsOn } from '$wailsjs/runtime/runtime';
  
export const smmUpdate = binding<autoupdate.PendingUpdate | null>(null, { initialGet: PendingUpdate, updateEvent: 'updateAvailable' });

export const smmUpdateChangelogs = derived(smmUpdate, ($smmUpdate) => {
  return $smmUpdate ? Object.entries($smmUpdate.changelogs).map(([version, changelog]) => ({ version, changelog })).sort((a, b) => -compare(a.version, b.version)) : null;
});

export const smmUpdateProgress = binding<utils.Progress | null>(null, { updateEvent: 'updateDownloadProgress' });

export const smmUpdateProgressStats = progressStats(smmUpdateProgress);

export const smmUpdateReady = readable<boolean>(false, (set) => {
  EventsOn('updateReady', () => {
    set(true);
  });
  EventsOn('updateAvailable', () => {
    set(false);
  });
});
