import { writable } from 'svelte/store';

import type { ficsitcli } from '$wailsjs/go/models';

// Because skeleton only ever keeps one modal loaded, we want to store this state between modals
export const selectedUpdates = writable<ficsitcli.Update[]>([]);
export const showIgnored = writable(false);
