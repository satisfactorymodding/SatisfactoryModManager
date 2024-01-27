import { writable } from 'svelte/store';

import { cacheDir } from '$lib/store/settingsStore';

export const newCacheLocation = writable<string | null>();
cacheDir.subscribe((value) => {
  newCacheLocation.set(value);
});
