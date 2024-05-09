<script lang="ts">
  import { getContextClient } from '@urql/svelte';
  import Fuse from 'fuse.js';
  import _ from 'lodash';
  import { createEventDispatcher } from 'svelte';

  import ModListFilters from './ModsListFilters.svelte';

  import VirtualList from '$lib/components/VirtualList.svelte';
  import AnnouncementsBar from '$lib/components/announcements/AnnouncementsBar.svelte';
  import ModsListItem from '$lib/components/mods-list/ModsListItem.svelte';
  import { GetModCountDocument, GetModsDocument } from '$lib/generated';
  import { queuedMods } from '$lib/store/actionQueue';
  import { favoriteMods, lockfileMods, manifestMods } from '$lib/store/ficsitCLIStore';
  import { expandedMod, hasFetchedMods } from '$lib/store/generalStore';
  import { type OfflineMod, type PartialMod, filter, order, search } from '$lib/store/modFiltersStore';
  import { offline, startView } from '$lib/store/settingsStore';
  import { OfflineGetMods } from '$wailsjs/go/ficsitcli/ficsitCLI';

  const dispatch = createEventDispatcher();

  const MODS_PER_PAGE = 100;

  const client = getContextClient();

  let fetchingMods = false;
  let onlineMods: PartialMod[] = [];
  async function fetchAllModsOnline() {
    try {
      const result = await client.query(GetModCountDocument, {}, { requestPolicy: 'network-only' }).toPromise();
      const count = result.data?.getMods.count;
      if (count && count !== onlineMods.length) {
        fetchingMods = true;
        const pages = Math.ceil(count / MODS_PER_PAGE);

        onlineMods = (await Promise.all(Array.from({ length: pages }).map(async (_, i) => {
          const offset = i * MODS_PER_PAGE;
          const modsPage = await client.query(GetModsDocument, { offset, limit: MODS_PER_PAGE }, { requestPolicy: 'network-only' }).toPromise();
          return modsPage.data?.getMods.mods ?? [];
        }))).flat();
      }
    } finally {
      fetchingMods = false;
      $hasFetchedMods = true;
    }
  }

  let offlineMods: PartialMod[] = [];
  async function fetchAllModsOffline() {
    offlineMods = (await OfflineGetMods()).map((mod) => ({
      ...mod,
      offline: true,
    } as OfflineMod));
  }
  
  let onlineRefreshInterval: number | undefined;

  $: if($offline !== null) {
    fetchAllModsOffline();
    if (!onlineRefreshInterval) {
      clearInterval(onlineRefreshInterval);
      onlineRefreshInterval = undefined;
    }
    if(!$offline) {
      fetchAllModsOnline();
      // setInterval returns NodeJS.Timeout, but that's not the case for the browser
      // eslint-disable-next-line
      // @ts-ignore
      onlineRefreshInterval = setInterval(fetchAllModsOnline, 5 * 60 * 1000); // 5 minutes
    } else {
      $hasFetchedMods = true;
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
  $: if(!$search) {
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
    <div style="position: relative;" class="py-4 grow h-0 mods-list @container/mods-list bg-surface-200-700-token">
      <div class="mr-4 h-full flex flex-col">
        {#if fetchingMods}
          <div class="flex items-center justify-center">
            <div class="animate-spin rounded-full aspect-square h-8 border-t-2 border-b-2 border-primary-500"/>
          </div>
        {/if}
        <VirtualList containerClass="mx-4" items={displayMods} let:item={mod}>
          <ModsListItem
            {mod}
            selected={$expandedMod == mod.mod_reference}
            on:click={() => {
              $expandedMod = mod.mod_reference;
              dispatch('expandedMod', mod.mod_reference);
            }}
          />
        </VirtualList>
      </div>
    </div>
  {/if}
</div>
