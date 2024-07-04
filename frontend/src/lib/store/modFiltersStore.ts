import type { Client } from '@urql/svelte';
import { get, writable } from 'svelte/store';

import { bindingTwoWayNoExcept } from './wailsStoreBindings';

import { CompatibilityState, type GetModsQuery } from '$lib/generated';
import { queuedMods } from '$lib/store/actionQueue';
import { favoriteMods, lockfileMods, manifestMods, selectedInstallMetadata } from '$lib/store/ficsitCLIStore';
import { getCompatibility } from '$lib/utils/modCompatibility';
import { installTypeToTargetName } from '$lib/wailsTypesExtensions';
import { GetModFiltersFilter, GetModFiltersOrder, SetModFiltersFilter, SetModFiltersOrder } from '$wailsjs/go/settings/settings';

export type OrderByField = 'name' | 'last-updated' | 'popularity' | 'hotness' | 'views' | 'downloads';
export type FilterField = 'all' | 'compatible' | 'favorite' | 'queued' | 'installed' | 'dependency' | 'not-installed' | 'enabled' | 'disabled';

export interface OrderBy {
  id: OrderByField;
  func: (mod: PartialMod) => unknown,
}

export interface Filter {
  id: FilterField;
  func: (mod: PartialMod, urqlClient: Client) => Promise<boolean> | boolean,
}

export const orderByOptions: OrderBy[] = [
  { id: 'name', func: (mod: PartialMod) => mod.name.trim() },
  { id: 'last-updated', func: (mod: PartialMod) => 'last_version_date' in mod ? Date.now() - Date.parse(mod.last_version_date) : 0 },
  { id: 'popularity', func: (mod: PartialMod) => 'popularity' in mod ? -mod.popularity : 0 },
  { id: 'hotness', func: (mod: PartialMod) => 'hotness' in mod ? -mod.hotness : 0 },
  { id: 'views', func: (mod: PartialMod) => 'views' in mod ? -mod.views : 0 },
  { id: 'downloads', func: (mod: PartialMod) => 'downloads' in mod ? -mod.downloads : 0 },
];

export const filterOptions: Filter[] = [
  { id: 'all', func: () => true },
  { 
    id: 'compatible',
    func: async (mod: PartialMod, urqlClient: Client) => { 
      const installInfo = get(selectedInstallMetadata).info;
      if(!installInfo) {
        return false;
      }
      const compatibility = await getCompatibility(mod.mod_reference, installInfo.branch, installInfo.version, installTypeToTargetName(installInfo.type), urqlClient);
      return compatibility.state !== CompatibilityState.Broken;
    }, 
  },
  { id: 'favorite', func: (mod: PartialMod) => get(favoriteMods).includes(mod.mod_reference) },
  { id: 'queued', func: (mod: PartialMod) => get(queuedMods).some((q) => q.mod === mod.mod_reference) },
  { id: 'installed', func: (mod: PartialMod) => mod.mod_reference in get(manifestMods) },
  { id: 'dependency', func: (mod: PartialMod) => !(mod.mod_reference in get(manifestMods)) && mod.mod_reference in get(lockfileMods) },
  { id: 'not-installed', func: (mod: PartialMod) => !(mod.mod_reference in get(manifestMods)) },
  { id: 'enabled', func: (mod: PartialMod) => get(manifestMods)[mod.mod_reference]?.enabled ?? mod.mod_reference in get(lockfileMods) },
  { id: 'disabled', func: (mod: PartialMod) => mod.mod_reference in get(manifestMods) && !(mod.mod_reference in get(lockfileMods)) },
];

export type PartialSMRMod = GetModsQuery['getMods']['mods'][number];
export interface OfflineMod {
  offline: true;
  mod_reference: string;
  name: string;
  logo?: string;
  authors: string[];
}
export interface MissingMod {
  missing: true;
  mod_reference: string;
  name: string;
  logo?: string;
  authors: string[];
}
export type PartialMod = PartialSMRMod | OfflineMod | MissingMod;

export const search = writable('');
export const order = bindingTwoWayNoExcept(orderByOptions[1], { 
  initialGet: async () => GetModFiltersOrder().then((i) => orderByOptions.find((o) => o.name === i) || orderByOptions[1]),
}, {
  updateFunction: async (o) => SetModFiltersOrder(o.name),
});
export const filter = bindingTwoWayNoExcept(filterOptions[0], {
  initialGet: async () => GetModFiltersFilter().then((i) => filterOptions.find((o) => o.name === i) || filterOptions[0]),
}, {
  updateFunction: async (f) => SetModFiltersFilter(f.name),
});
