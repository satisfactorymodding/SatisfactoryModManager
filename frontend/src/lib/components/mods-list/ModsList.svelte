<script lang="ts">
  import { getContextClient } from '@urql/svelte';
  import _ from 'lodash';
  import Fuse from 'fuse.js';

  import AnnouncementsBar from '../announcements/AnnouncementsBar.svelte';

  import ModListFilters from './ModsListFilters.svelte';

  import { GetModsDocument, GetModCountDocument } from '$lib/generated';
  import VirtualList from '$lib/components/mods-list/VirtualModList.svelte';
  import ModsListItem from '$lib/components/mods-list/ModsListItem.svelte';
  import { filter, order, search, type PartialMod } from '$lib/store/modFiltersStore';
  import { favoriteMods, lockfileMods, manifestMods, queuedMods } from '$lib/store/ficsitCLIStore';
  import { offline, startView } from '$lib/store/settingsStore';
  import { expandedMod } from '$lib/store/generalStore';
  import { OfflineGetMods } from '$wailsjs/go/ficsitcli_bindings/FicsitCLI';

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
    }));
  }
  
  $: if($offline !== null) {
    if($offline) {
      fetchAllModsOffline();
    } else {
      fetchAllModsOnline();
    }
  }

  $: mods = $offline ? offlineMods : onlineMods;

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
</script>

<div class="h-full flex flex-col">
  <div class="flex-none">
    <ModListFilters />
  </div>
  <AnnouncementsBar />
  <div class="py-4 grow h-0 mods-list" style="position: relative;">
    <div class="ml-5 mr-3 h-full">
      <VirtualList items={displayMods} let:item={mod}>
        <ModsListItem {mod} on:click={() => $expandedMod = mod.mod_reference} selected={$expandedMod == mod.mod_reference}/>
      </VirtualList>
    </div>
  </div>
</div>