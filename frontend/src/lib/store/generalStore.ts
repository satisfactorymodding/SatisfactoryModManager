import { writable } from 'svelte/store';

export const expandedMod = writable(null as string | null);
export const error = writable<string|null>(null);
export const isLaunchingGame = writable(false);
export const viewedAnnouncements = writable<string[]>([]);