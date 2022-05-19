import type { Mod } from '$lib/generated';

export interface OrderBy {
  name: string;
  func: (mod: Mod) => unknown,
}

export const orderByOptions: OrderBy[] = [
  { name: 'Name', func: (mod: Mod) => mod.name.trim() },
  { name: 'Last updated', func: (mod: Mod) => Date.now() - Date.parse(mod.last_version_date) },
];

export type PartialMod = Pick<Mod, 'id' | 'mod_reference' | 'name' | 'logo' | 'short_description' | 'views' | 'downloads' | 'last_version_date'> & { authors: { user: Pick<User, 'username'> }[] };