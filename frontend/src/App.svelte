<script lang="ts">
  import './_global.postcss';
  import TitleBar from '$lib/components/TitleBar.svelte';
  import ModsList from '$lib/components/mods-list/ModsList.svelte';
  import { initializeGraphQLClient } from '$lib/core/graphql';
  import { setClient } from '@urql/svelte';
  import { Environment } from '$wailsjs/runtime/runtime';
  import ModDetails from '$lib/components/mod-details/ModDetails.svelte';
  import { ExpandMod, UnexpandMod } from '$wailsjs/go/bindings/App';
  import LeftBar from '$lib/components/LeftBar.svelte';

  Environment().then((env) => {
    if (env.buildType !== 'dev') {
      document.addEventListener('contextmenu', (event) => event.preventDefault());
    }
  });

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
      }, 100);
    }
  }

  $: pendingExpand = selectedModId && !windowExpanded;

  $: modsListCompact = windowExpanded;
</script>

<div class="flex flex-col h-screen w-screen">
  <TitleBar />
  <div class="flex grow h-0 select-none">
    <div class="left-bar">
      <LeftBar />
    </div>
    <div class:normal={!selectedModId || pendingExpand} class:compact={windowExpanded}>
      <ModsList bind:selectedMod={selectedModId} bind:compact={modsListCompact}/>
    </div>
    {#if selectedModId}
      <div class="grow">
        <ModDetails id={selectedModId} on:close={() => selectedModId = null}/>
      </div>
    {/if}
  </div>
</div>

<style>
  .normal {
    width: 610px !important;
    min-width: 610px !important;
  }
  .compact {
    width: 470px;
    min-width: 470px;
  }
  .left-bar {
    width: 24rem;
    min-width: 24rem;
  }
  /* .pendingExpand {
    display: none;
  } */
</style>
