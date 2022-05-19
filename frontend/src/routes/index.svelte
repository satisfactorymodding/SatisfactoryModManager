<script lang="ts">
  import './_global.postcss';
  import TitleBar from '$lib/components/TitleBar.svelte';
  import ModsList from '$lib/components/mods-list/ModsList.svelte';
  import { initializeGraphQLClient } from '$lib/core/graphql';
  import { setClient } from '@urql/svelte';
  import { Environment, type EnvironmentInfo } from '../../wailsjs/runtime/runtime';
  import ModDetails from '$lib/components/mod-details/ModDetails.svelte';
  import { ExpandMod, UnexpandMod } from '../../wailsjs/go/main/App';
  import LeftBar from '$lib/components/LeftBar.svelte';
  
  // Workaround for incorrect wails type definition
  async function WailsEnv(): Promise<EnvironmentInfo & { buildtype: string }> {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    return Environment();
  }

  if (typeof window !== 'undefined') {
    WailsEnv().then((env) => {
      if (env.buildtype !== 'dev') {
        document.addEventListener('contextmenu', (event) => event.preventDefault());
      }
    });
  }

  setClient(initializeGraphQLClient());

  let selectedModId: string | null = null;

  let windowExpanded = false;

  $: if(typeof window !== 'undefined') { 
    if (selectedModId) {
      ExpandMod().then(() => { windowExpanded = true; });
    } else {
      windowExpanded = false;
      setTimeout(() => {
        UnexpandMod();
      }, 200);
    }
  }

  $: pendingExpand = selectedModId && !windowExpanded;

  $: modsListCompact = windowExpanded;
</script>

<div class="flex flex-col h-screen w-screen">
  <TitleBar />
  <div class="flex grow h-0 select-none">
    <div class="w-96">
      <LeftBar />
    </div>
    <div class:normal={!selectedModId || pendingExpand} class:compact={windowExpanded}>
      <ModsList bind:selectedMod={selectedModId} bind:compact={modsListCompact}/>
    </div>
    {#if selectedModId}
      <div class="flex-1" class:pendingExpand>
        <ModDetails id={selectedModId} on:close={() => selectedModId = null}/>
      </div>
    {/if}
  </div>
</div>

<style>
  .normal {
    width: 610px !important;
  }
  .compact {
    width: 470px;
  }
  .pendingExpand {
    display: none;
  }
</style>
