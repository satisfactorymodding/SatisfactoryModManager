<script lang="ts">
  import { getContextClient } from '@urql/svelte';
  import _ from 'lodash';
  import Fuse from 'fuse.js';
  import { createEventDispatcher } from 'svelte';

  import AnnouncementsBar from '../announcements/AnnouncementsBar.svelte';

  import ModListFilters from './ModsListFilters.svelte';

  import { GetModsDocument, GetModCountDocument } from '$lib/generated';
  import VirtualList from '$lib/components/mods-list/VirtualModList.svelte';
  import ModsListItem from '$lib/components/mods-list/ModsListItem.svelte';
  import { filter, order, search, type PartialMod, type OfflineMod } from '$lib/store/modFiltersStore';
  import { favoriteMods, lockfileMods, manifestMods } from '$lib/store/ficsitCLIStore';
  import { queuedMods } from '$lib/store/actionQueue';
  import { offline, startView } from '$lib/store/settingsStore';
  import { expandedMod } from '$lib/store/generalStore';
  import { OfflineGetMods } from '$wailsjs/go/ficsitcli/FicsitCLI';

  const dispatch = createEventDispatcher();

  const MODS_PER_PAGE = 50;

  const client = getContextClient();

  let onlineMods: PartialMod[] = [];
  async function fetchAllModsOnline() {
    const result = await client.query(GetModCountDocument, {}).toPromise();
    const count = result.data?.getMods.count;
    if (count) {
      const pages = Math.ceil(count / MODS_PER_PAGE);

      onlineMods = (await Promise.all(Array.from({ length: pages }).map(async (_, i) => {
        const offset = i * MODS_PER_PAGE;
        const modsPage = await client.query(GetModsDocument, { offset, limit: MODS_PER_PAGE }).toPromise();
        return modsPage.data?.getMods.mods ?? [];
      }))).flat();
    }
  }

  let offlineMods: PartialMod[] = [];
  async function fetchAllModsOffline() {
    offlineMods = (await OfflineGetMods()).map((mod) => ({
      ...mod,
      offline: true,
    } as OfflineMod));
  }
  
  $: if($offline !== null) {
    fetchAllModsOffline();
    if(!$offline) {
      fetchAllModsOnline();
    }
  }

  $: knownMods = $offline ? offlineMods : onlineMods;

  $: unknownModReferences = Object.keys($manifestMods)
    .filter((modReference) => !knownMods.find((knownMod) => knownMod.mod_reference === modReference) && modReference !== 'SML');

  $: unknownMods = unknownModReferences.map((modReference) => {
    const offlineMod = offlineMods.find((mod) => mod.mod_reference === modReference);
    const mod = {
      mod_reference: modReference,
      name: offlineMod ? offlineMod.name : modReference,
      logo: offlineMod ? offlineMod.logo : undefined,
      authors: offlineMod ? offlineMod.authors : ['N/A'],
      missing: true,
    } as PartialMod;
    return mod;
  });

  $: mods = [...knownMods, ...unknownMods];

  let filteredMods: PartialMod[] = [];
  $: {
    // Watch the required store states
    $manifestMods;
    $lockfileMods;
    $favoriteMods;
    $queuedMods;
    
    Promise.all(mods.map((mod) => $filter.func(mod, client))).then((results) => {
      filteredMods = mods.filter((_, i) => results[i]);
    });
  }

  let sortedMods: PartialMod[] = [];
  $: {
    // Watch the required store states
    $manifestMods;
    $lockfileMods;
    $favoriteMods;
    $queuedMods;
    
    sortedMods = _.sortBy(filteredMods, $order.func) as PartialMod[];
  }

  let displayMods: PartialMod[] = [];
  $: {
    if(!$search) {
      displayMods = sortedMods;
    } else {
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
            name: $offline ? 'authors' : 'authors.user.username',
            weight: 0.4,
          },
        ],
        useExtendedSearch: true,
        threshold: 0.2,
        ignoreLocation: true,
      });
      displayMods = fuse.search(modifiedSearchString).map((result) => result.item);
    }
  }

  let hasCheckedStartView = false;
  $: if($startView && mods.length > 0 && !hasCheckedStartView) {
    hasCheckedStartView = true;
    if($startView === 'expanded') {
      if(displayMods.length > 0) {
        $expandedMod = displayMods[0].mod_reference;
      }
    }
  }

  export let hideMods: boolean = false;
</script>

<div class="h-full flex flex-col">
  <div class="flex-none z-[1]">
    <ModListFilters />
  </div>
  <AnnouncementsBar />
  {#if hideMods}
    <slot />
  {:else}
    <div class="py-4 grow h-0 mods-list @container/mods-list" style="position: relative;">
      <div class="ml-5 mr-3 h-full">
        <VirtualList items={displayMods} let:item={mod}>
          <ModsListItem
            {mod}
            on:click={() => {
              $expandedMod = mod.mod_reference;
              dispatch('expandedMod', mod.mod_reference);
            }}
            selected={$expandedMod == mod.mod_reference}
          />
        </VirtualList>
      </div>
    </div>
  {/if}
</div>
