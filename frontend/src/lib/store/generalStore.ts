import { writable } from 'svelte/store';

export const expandedMod = writable(null as string | null);
export const error = writable<string|null>(null);