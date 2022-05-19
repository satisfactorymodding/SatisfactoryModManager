<script lang="ts">
  import { mdiDownload, mdiEye } from '@mdi/js';
  import MDIIcon from '$lib/components/MDIIcon.svelte';
  import { createEventDispatcher } from 'svelte';
  import type { PartialMod } from './modFilters';
  
  export let mod: PartialMod;

  const dispatch = createEventDispatcher();

  function click() {
    dispatch('click');
  }
  
  export let compact: boolean;
  export let selected: boolean;

  $: renderedLogo = mod.logo || 'https://ficsit.app/images/no_image.webp';
  $: author = mod.authors[0].user.username;
</script>

<div class="my-2 px-0 flex w-full" readonly on:click={click} class:selected>
  <img src={renderedLogo} alt="{mod.name} Logo" class="logo h-24 w-24" />
  <div class="ml-2 flex flex-col grow w-0">
    <div>
      <span class="text-xl font-medium">{mod.name}</span>
      {#if !compact}
        <span class="pl-1">by</span>
        <span class="color-primary">{author}</span>
      {/if}
    </div>
    <span class="truncate w-full">{mod.short_description}</span>
    <div>
      #tags
    </div>
    <div class="flex">
      <div class="w-24 flex">
        <div class="pr-1 inline-flex items-center justify-items-center"><MDIIcon icon={mdiEye}/></div>{mod.views.toLocaleString()}
      </div>
      <div class="w-24 flex">
        <div class="pr-1 inline-flex items-center justify-items-center"><MDIIcon icon={mdiDownload}/></div>{mod.downloads.toLocaleString()}
      </div>
    </div>
  </div>
</div>

<style>
  .selected {
    background-color: #1c1c1c;
    border-radius: 16px;
  }
</style>