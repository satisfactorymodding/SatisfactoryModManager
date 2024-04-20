<script lang="ts">
  import { mdiFolder, mdiLoading, mdiServerNetwork, mdiSubdirectoryArrowLeft } from '@mdi/js';
  import _ from 'lodash';
  import { onDestroy } from 'svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { StartPicker, StopPicker, TryPick } from '$lib/generated/wailsjs/go/ficsitcli/serverPicker';


  export let basePath: string;
  export let disabled = false;
  
  let currentBasePath = '';
  let pickerId = '';

  async function stopPicker() {
    if (pickerId) {
      await StopPicker(pickerId);
      pickerId = '';
    }
  }
  
  const restartPicker = _.debounce(async () => {
    await stopPicker();
    if (!disabled) {
      pickerId = await StartPicker(basePath);
      currentBasePath = basePath;
    }
  // TODO: handle errors
  }, 1000);

  let pendingValidCheck = false;

  $: if (basePath !== currentBasePath && !disabled) {
    pendingValidCheck = true;
    displayedItems = [];
    restartPicker();
  }

  $: if (disabled) {
    stopPicker();
    currentBasePath = '';
  }

  onDestroy(() => {
    stopPicker();
  });

  interface PickerDirectory {
    path: string;
    name: string;
    isValidInstall: boolean;
  }

  let displayedItems: { path: string; name: string; isValidInstall: boolean }[] = [];

  export let path = '';
  export let valid = false;

  let validError: string | null = null;

  const checkValid = _.debounce(async () => {
    pendingValidCheck = true;
    const forPath = path;
    const forPickerId = pickerId;
    try {
      let { isValidInstall } = await TryPick(pickerId, path + '/');
      // If the path has changed since the request was made,
      // or the picker has been stopped, ignore the response
      if (path !== forPath || pickerId !== forPickerId) {
        return;
      }
      valid = isValidInstall;
      validError = null;
    } catch (e) {
      // If the path has changed since the request was made,
      // or the picker has been stopped, ignore the response
      if (path !== forPath || pickerId !== forPickerId) {
        return;
      }
      valid = false;
      validError = e as string;
    }
    pendingValidCheck = false;
  }, 250);

  $: if (pickerId && !disabled) {
    path;
    pendingValidCheck = true;
    checkValid();
  }

  $: trimmedPath = path.endsWith('/') ? path.slice(0, -1) : path;

  let displayedPath = '';
  $: if (!pendingValidCheck) {
    displayedPath = valid ? trimmedPath.split('/').slice(0, -1).join('/') : trimmedPath;
  }

  let pendingDisplay = false;

  let error: string | null = null;

  const updateDisplay = _.debounce(async () => {
    pendingDisplay = true;
    const forPath = displayedPath;
    const forPickerId = pickerId;
    try {
      let { items } = await TryPick(pickerId, displayedPath + '/');
      // If the path has changed since the request was made,
      // or the picker has been stopped, ignore the response
      if (displayedPath !== forPath || pickerId !== forPickerId) {
        return;
      }
      displayedItems = items.map((d) => ({
        path: displayedPath + '/' + d.name,
        name: d.name,
        isValidInstall: d.isValidInstall,
      }));
      error = null;
    } catch (e) {
      // If the path has changed since the request was made,
      // or the picker has been stopped, ignore the response
      if (path !== forPath || pickerId !== forPickerId) {
        return;
      }
      error = e as string;
    }
    pendingDisplay = false;
  }, 250);

  $: if(pickerId && !disabled) {
    displayedPath;
    pendingDisplay = true;
    updateDisplay();
  }

  function select(item: PickerDirectory) {
    valid = item.isValidInstall;
    path = item.path;
  }

  // If the path is root and it is a valid install, don't list contents of root
  $: actualDisplayedItems = (valid && displayedPath === trimmedPath) ? [{ path: '', name: '(root)', isValidInstall: true }] : displayedItems;
</script>

<div class="relative">
  <div class="flex flex-col w-full card bg-surface-200-700-token">
    <button class="w-full btn !scale-100" disabled={displayedPath.length <= 1 || pendingDisplay || pendingValidCheck} on:click={() => { path = displayedPath.split('/').slice(0, -1).join('/'); valid = false; }}>
      <SvgIcon
        class="h-5 w-5"
        icon={mdiSubdirectoryArrowLeft} />
      <div class="grow"/>
    </button>
    {#if !disabled}
      <div class="overflow-y-auto">
        {#each actualDisplayedItems as item}
          <button
            class="w-full btn !scale-100"
            class:variant-ghost-primary={path.startsWith(item.path) && valid}
            disabled={pendingDisplay || pendingValidCheck}
            on:click={() => select(item)}>
            <SvgIcon
              class="h-5 w-5"
              icon={item.isValidInstall ? mdiServerNetwork : mdiFolder} />
            <span>{item.name}</span>
            <div class="grow"/>
          </button>
        {/each}
      </div>
      {#if validError && !pendingValidCheck && !error}
        <div class="text-error-500 p-4">Failed to check if selected path is a valid server</div>
      {/if}
    {/if}
  </div>
  {#if (((pendingDisplay || pendingValidCheck) && !valid) || error) && !disabled}
    <div class="w-full h-full flex justify-center card items-center absolute top-0 !bg-surface-600/80">
      {#if ((pendingDisplay || pendingValidCheck) && !valid)}
        <SvgIcon
          class="h-10 w-10 animate-spin text-primary-600"
          icon={mdiLoading} />
      {:else}
        <div class="text-error-500">Failed to list directory</div>
      {/if}
    </div>
  {/if}
</div>
