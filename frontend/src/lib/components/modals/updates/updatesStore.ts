import { writable } from 'svelte/store';

// Because skeleton only ever keeps one modal loaded, we want to store this state between modals
export const selectedUpdates = writable<string[]>([]);
export const showIgnored = writable(false);
