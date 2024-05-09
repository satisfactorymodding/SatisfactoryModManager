import { writable } from 'svelte/store';

export const expandedMod = writable(null as string | null);
export const error = writable<string|null>(null);
export const isLaunchingGame = writable(false);
export const siteURL = writable<string>('https://ficsit.app/');
export const hasFetchedMods = writable(false);
