<script lang="ts">
  import { getClient } from '@urql/svelte';
  import { GetModsDocument, GetModCountDocument } from '$lib/generated';
  import VirtualList from '$lib/components/mods-list/VirtualModList.svelte';
  import ModsListItem from '$lib/components/mods-list/ModsListItem.svelte';
  import _ from 'lodash';
  import Fuse from 'fuse.js';
  import ModListFilters from './ModsListFilters.svelte';
  import { filter, order, search, type PartialMod } from '$lib/modFiltersStore';
  import { favouriteMods, lockfileMods, manifestMods } from '$lib/ficsitCLIStore';

  let mods: PartialMod[] = [];

  const MODS_PER_PAGE = 50;

  const urqlClient = getClient();

  async function fetchAllMods() {
    const result = await urqlClient.query(GetModCountDocument).toPromise();
    const count = result.data?.getMods.count;
    if (count) {
      const pages = Math.ceil(count / MODS_PER_PAGE);

      mods = (await Promise.all(Array.from({length: pages}).map(async (_, i) => {
        const offset = i * MODS_PER_PAGE;
        const modsPage = await urqlClient.query(GetModsDocument, { offset, limit: MODS_PER_PAGE }).toPromise();
        return modsPage.data?.getMods.mods ?? [];
      }))).flat();
    }
  }

  fetchAllMods();

  $: filteredMods = () => {
    // Watch the required store states
    $manifestMods;
    $lockfileMods;
    $favouriteMods;
    
    const filteredMods = mods.filter($filter.func);
    const sortedMods = _.sortBy(filteredMods, $order.func) as PartialMod[];
    if(!$search) {
      return sortedMods;
    }
    
    const modifiedSearchString = $search.replace(/(?:author:"(.+?)"|author:([^\s"]+))/g, '="$1$2"');
    
    const fuse = new Fuse(sortedMods, {
      keys: [
        {
          name: 'name',
          weight: 2,
        },
        {
          name: 'short_description',
          weight: 1,
        },
        {
          name: 'full_description',
          weight: 0.75,
        },
        {
          name: 'authors.user.username',
          weight: 0.4,
        },
      ],
      useExtendedSearch: true,
      threshold: 0.2,
      ignoreLocation: true,
    });
    return fuse.search(modifiedSearchString).map((result) => result.item);
  };

  export let selectedMod: string | null = null;
  export let compact: boolean;
</script>

<div class="h-full flex flex-col">
  <div class="flex-none">
    <ModListFilters bind:compact />
  </div>
  <div class="py-4 grow h-0 mods-list" style="position: relative;">
    <div class="ml-5 mr-3 h-full">
      <VirtualList items={filteredMods()} let:item={mod}>
        <ModsListItem {mod} on:click={() => selectedMod = mod.mod_reference} bind:compact={compact} selected={selectedMod == mod.mod_reference}/>
      </VirtualList>
    </div>
  </div>
</div>