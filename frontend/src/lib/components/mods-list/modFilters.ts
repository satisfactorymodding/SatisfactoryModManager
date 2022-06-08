import type { Mod } from '$lib/generated';
import { favouriteMods, lockfileMods, manifestMods } from '$lib/store';
import { get } from 'svelte/store';

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

export type PartialMod = Pick<Mod, 'id' | 'mod_reference' | 'name' | 'logo' | 'short_description' | 'views' | 'downloads' | 'last_version_date'> & { authors: { user: Pick<User, 'username'> }[] };