import { writable } from 'svelte/store';

import type { ficsitcli } from '$lib/generated/wailsjs/go/models';

// Because skeleton only ever keeps one modal loaded, we want to store this state between modals
export const selectedUpdates = writable<ficsitcli.Update[]>([]);
export const showIgnored = writable(false);
