<script lang="ts">
  import { mdiDownload } from '@mdi/js';
  import { getTranslate } from '@tolgee/svelte';
  import { getContextClient, queryStore } from '@urql/svelte';
  import { onMount } from 'svelte';

  import { selectedUpdates, showIgnored } from './updatesStore';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import T from '$lib/components/T.svelte';
  import ModChangelog from '$lib/components/modals/ModChangelog.svelte';
  import { GetModNamesDocument } from '$lib/generated';
  import { getModalStore } from '$lib/skeletonExtensions';
  import { canModify, unignoredUpdates, updateCheckInProgress, updates } from '$lib/store/ficsitCLIStore';
  import { error } from '$lib/store/generalStore';
  import { offline } from '$lib/store/settingsStore';
  import { OfflineGetModsByReferences, UpdateMods } from '$wailsjs/go/ficsitcli/ficsitCLI';
  import type { ficsitcli } from '$wailsjs/go/models';
  import { SetUpdateIgnore, SetUpdateUnignore } from '$wailsjs/go/settings/settings';

  export let parent: { onClose: () => void };

  const { t } = getTranslate();

  const modalStore = getModalStore();

  const client = getContextClient();

  $: modNamesQuery = queryStore({
    query: GetModNamesDocument,
    client,
    pause: !!$offline,
    variables: {
      modReferences: $updates.map((u) => u.item),
    },
  });

  let modNamesQueryResult: { mod_reference: string; name: string; }[] | undefined;

  $: modNames = modNamesQueryResult?.reduce((acc, mod) => {
    if(mod) {
      acc[mod.mod_reference] = mod.name;
    }
    return acc;
  }, {} as Record<string, string>) ?? {};

  $: if($offline) {
    OfflineGetModsByReferences($updates.map((u) => u.item)).then((mods) => modNamesQueryResult = mods);
  } else {
    modNamesQueryResult = $modNamesQuery.data?.getMods?.mods;
  }


  $: updatesToDisplay = $showIgnored ? $updates : $unignoredUpdates;

  async function updateAll() {
    if(updatesToDisplay.length > 0) {
      try {
        await UpdateMods(updatesToDisplay.map((u) => u.item));
        $updates = $updates.filter((u) => !updatesToDisplay.includes(u));
      } catch(e) {
        if (e instanceof Error) {
          $error = e.message;
        } else if (typeof e === 'string') {
          $error = e;
        } else {
          $error = 'Unknown error';
        }
      }
    }
  }

  async function updateSelected() {
    if($selectedUpdates.length > 0) {
      try {
        await UpdateMods($selectedUpdates.map((u) => u.item));
        $updates = $updates.filter((u) => !$selectedUpdates.includes(u));
      } catch(e) {
        if (e instanceof Error) {
          $error = e.message;
        } else if (typeof e === 'string') {
          $error = e;
        } else {
          $error = 'Unknown error';
        }
      }
    }
  }

  function toggleSelected(update: ficsitcli.Update) {
    if($selectedUpdates.includes(update)) {
      $selectedUpdates = $selectedUpdates.filter((u) => u !== update);
    } else {
      $selectedUpdates = [...$selectedUpdates, update];
    }
  }

  $: () => {
    $updates;
    $selectedUpdates = [];
  };

  function toggleIgnoreUpdate(update: ficsitcli.Update) {
    if($unignoredUpdates.includes(update)) {
      SetUpdateIgnore(update.item, update.newVersion);
      $selectedUpdates = $selectedUpdates.filter((u) => u !== update);
    } else {
      SetUpdateUnignore(update.item, update.newVersion);
    }
  }

  onMount(() => {
    if (!$unignoredUpdates.length) {
      parent.onClose();
    }
  });
</script>

<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="w-[48rem] card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    <T defaultValue="Updates" keyName="updates.title" />
  </header>
  <section class="px-4">
    <button
      class="btn"
      on:click={() => $showIgnored = !$showIgnored}>
      {$showIgnored ? $t('updates.hide-ignored', 'Hide ignored') : $t('updates.show-ignored', 'Show ignored')}
    </button>
  </section>
  <section class="px-4 flex-auto grid grid-cols-12 overflow-y-auto">
    {#each updatesToDisplay as update}
      <button class="btn p-2 col-span-8 text-left space-x-2" on:click={() => toggleSelected(update)}>
        <div class="h-full w-6">
          {#if $selectedUpdates.includes(update)}
            <SvgIcon class="h-full w-8 mx-auto" icon={mdiDownload} />
          {/if}
        </div>
        <div class="h-full flex-auto flex flex-col content-center">
          <span>{modNames[update.item] ?? update.item}</span>
          <span>{update.currentVersion} -> {update.newVersion}</span>
        </div>
      </button>
      <button
        class="btn col-span-2"
        on:click={() => modalStore.trigger({ type:'component', component:{ ref: ModChangelog, props:{ mod:update.item, versionRange:{ from:update.currentVersion, to:update.newVersion } } } }, true)}>
        <T defaultValue="Changelog" keyName="updates.changelog" />
      </button>
      <button
        class="btn col-span-2"
        on:click={() => toggleIgnoreUpdate(update)}>
        {$unignoredUpdates.includes(update) ? $t('updates.ignore', 'Ignore') : $t('updates.unignore', 'Unignore')}
      </button>
    {/each}
  </section>
  <footer class="card-footer">
    <button
      class="btn"
      on:click={parent.onClose}>
      <T defaultValue="Close" keyName="common.close" />
    </button>
    <button
      class="btn"
      disabled={!$canModify || $updateCheckInProgress || !updatesToDisplay.length}
      on:click={() => updateAll()}>
      <T defaultValue="Update All" keyName="updates.update-all" />
    </button>
    <button
      class="btn"
      disabled={!$canModify || $updateCheckInProgress || !updatesToDisplay.length || !$selectedUpdates.length}
      on:click={() => updateSelected()}>
      <T defaultValue="Update Selected" keyName="updates.update-selected" />
    </button>
  </footer>
</div>
