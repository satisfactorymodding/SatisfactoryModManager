import { writable } from 'svelte/store';

export const profileFilepath = writable<string>('');
export const profileName = writable<string>('');
