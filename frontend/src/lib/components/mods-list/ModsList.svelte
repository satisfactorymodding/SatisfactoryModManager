<script lang="ts">
  import { getClient } from '@urql/svelte';
  import { GetModsDocument, GetModCountDocument } from '$lib/generated';
  import VirtualList from '$lib/components/mods-list/VirtualModList.svelte';
  import ModsListItem from '$lib/components/mods-list/ModsListItem.svelte';
  import _ from 'lodash';
  import Fuse from 'fuse.js';
  import ModListFilters from './ModsListFilters.svelte';
  import { filterOptions, orderByOptions, type Filter, type OrderBy, type PartialMod } from '$lib/components/mods-list/modFilters';
  import { lockfileMods, manifestMods } from '$lib/store';

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
        const result = await urqlClient.query(GetModsDocument, { offset, limit: MODS_PER_PAGE }).toPromise();
        return result.data?.getMods.mods ?? [];
      }))).flat();
    }
  }

  fetchAllMods();

  let searchString = '';
  let order: OrderBy = orderByOptions[1];
  let filter: Filter = filterOptions[0];

  $: filteredMods = () => {
    // Watch the required store states
    $manifestMods;
    $lockfileMods;
    
    const filteredMods = mods.filter(filter.func);
    const sortedMods = _.sortBy(filteredMods, order.func) as PartialMod[];
    if(!searchString) {
      return sortedMods;
    }
    
    const modifiedSearchString = searchString.replace(/(?:author:"(.+?)"|author:([^\s"]+))/g, '="$1$2"');
    
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
    <ModListFilters bind:search={searchString} bind:order={order}  bind:filter={filter} bind:compact={compact} />
  </div>
  <div class="py-4 grow h-0 mods-list" style="position: relative;">
    <div class="ml-5 mr-3 h-full">
      <VirtualList items={filteredMods()} let:item={mod}>
        <ModsListItem {mod} on:click={() => selectedMod = mod.mod_reference} bind:compact={compact} selected={selectedMod == mod.mod_reference}/>
      </VirtualList>
    </div>
  </div>
</div>