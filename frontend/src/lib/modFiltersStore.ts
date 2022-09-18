import type { GetModsQuery } from '$lib/generated';
import { favouriteMods, lockfileMods, manifestMods } from '$lib/ficsitCLIStore';
import { get, writable } from 'svelte/store';
import { writableBinding } from './utils/wailsStoreBindings';
import { GetModFilters, SetModFilters } from '$wailsjs/go/bindings/Settings';

export interface OrderBy {
  name: string;
  func: (mod: PartialMod) => unknown,
}

export interface Filter {
  name: string;
  func: (mod: PartialMod) => boolean,
}

export const orderByOptions: OrderBy[] = [
  { name: 'Name', func: (mod: PartialMod) => mod.name.trim() },
  { name: 'Last updated', func: (mod: PartialMod) => Date.now() - Date.parse(mod.last_version_date) },
  { name: 'Popularity', func: (mod: PartialMod) => -mod.popularity },
  { name: 'Hotness', func: (mod: PartialMod) => -mod.hotness },
  { name: 'Views', func: (mod: PartialMod) => -mod.views },
  { name: 'Downloads', func: (mod: PartialMod) => -mod.downloads },
];

export const filterOptions: Filter[] = [
  { name: 'All mods', func: () => true },
  { name: 'Favourite', func: (mod: PartialMod) => get(favouriteMods).includes(mod.mod_reference) },
  { name: 'Installed', func: (mod: PartialMod) => mod.mod_reference in get(manifestMods) },
  { name: 'Not installed', func: (mod: PartialMod) => !(mod.mod_reference in get(manifestMods)) },
  { name: 'Enabled', func: (mod: PartialMod) => mod.mod_reference in get(lockfileMods) },
  { name: 'Disabled', func: (mod: PartialMod) => mod.mod_reference in get(manifestMods) && !(mod.mod_reference in get(lockfileMods)) },
];

export type PartialMod = GetModsQuery['getMods']['mods'][number];

export const search = writable('');
export const order = writableBinding(orderByOptions[1], { 
  initialGet: async () => {
    const { order } = await GetModFilters();
    return orderByOptions.find(o => o.name === order) || orderByOptions[1];
  }
});
export const filter = writableBinding(filterOptions[0], { 
  initialGet: async () => {
    const { filter } = await GetModFilters();
    return filterOptions.find(o => o.name === filter) || filterOptions[0];
  }
});

order.subscribe((o) => {
  if(order.isInit) {
    SetModFilters(o.name, get(filter).name);
  }
});

filter.subscribe((f) => {
  if(filter.isInit) {
    SetModFilters(get(order).name, f.name);
  }
});