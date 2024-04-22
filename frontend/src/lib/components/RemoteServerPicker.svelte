<script lang="ts">
  import { mdiFolder, mdiLoading, mdiServerNetwork, mdiSubdirectoryArrowLeft } from '@mdi/js';
  import _ from 'lodash';
  import { onDestroy } from 'svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { GetPathSeparator, StartPicker, StopPicker, TryPick } from '$lib/generated/wailsjs/go/ficsitcli/serverPicker';
  import type { ficsitcli } from '$lib/generated/wailsjs/go/models';


  export let basePath: string;
  export let disabled = false;
  
  let currentBasePath: string | null = null;
  let pickerId = '';

  let platformPathSeparator = '/';

  GetPathSeparator().then((sep) => {
    platformPathSeparator = sep;
  });

  $: isLocalPath = (() => {
    try {
      const url = new URL(basePath);
      if (url.protocol === 'ftp:' || url.protocol === 'sftp:') {
        return false;
      }
      return true;
    } catch {
      // If not parsable as a URL, it is definitely a local path
      return true;
    }
  })();

  $: pathSeparator = isLocalPath ? platformPathSeparator : '/';

  function parentPath(path: string) {
    return path.split(pathSeparator).slice(0, -1).join(pathSeparator);
  }

  async function stopPicker() {
    if (pickerId) {
      await StopPicker(pickerId);
      pickerId = '';
    }
  }

  let setupError: string | null = null;
  
  const restartPicker = _.debounce(async () => {
    await stopPicker();
    if (!disabled) {
      try {
        pickerId = await StartPicker(basePath);
        currentBasePath = basePath;
        setupError = null;
      } catch (e) {
        setupError = e as string;
      }
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
    currentBasePath = null;
  }

  onDestroy(() => {
    stopPicker();
  });

  let displayedItems: { path: string; name: string; isValidInstall: boolean }[] = [];

  export let path = '';
  export let valid = false;

  let validError: string | null = null;

  const checkValid = _.debounce(async () => {
    pendingValidCheck = true;
    const forPath = path;
    const forPickerId = pickerId;
    try {
      let { isValidInstall } = await TryPick(pickerId, path + pathSeparator);
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

  $: trimmedPath = path.endsWith(pathSeparator) ? path.slice(0, -1) : path;

  let displayedPath = '';
  $: if (!pendingValidCheck) {
    displayedPath = valid ? parentPath(trimmedPath) : trimmedPath;
  }

  let pendingDisplay = false;

  let error: string | null = null;

  const updateDisplay = _.debounce(async () => {
    pendingDisplay = true;
    const forPath = displayedPath;
    const forPickerId = pickerId;
    try {
      let { items } = await TryPick(pickerId, displayedPath + pathSeparator);
      // If the path has changed since the request was made,
      // or the picker has been stopped, ignore the response
      if (displayedPath !== forPath || pickerId !== forPickerId) {
        return;
      }
      displayedItems = items;
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

  function select(item: ficsitcli.PickerDirectory) {
    valid = item.isValidInstall;
    path = item.path;
  }

  // If the path is root and it is a valid install, don't list contents of root
  $: actualDisplayedItems = (valid && displayedPath === trimmedPath) ? [{ path: '', name: '(root)', isValidInstall: true }] : displayedItems;
</script>

<div class="relative">
  <div class="flex flex-col w-full card bg-surface-200-700-token">
    <button class="w-full btn !scale-100" disabled={displayedPath.length <= 1 || pendingDisplay || pendingValidCheck} on:click={() => { path = parentPath(displayedPath); valid = false; }}>
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
  {#if (((pendingDisplay || pendingValidCheck) && !valid) || error || setupError) && !disabled}
    <div class="w-full h-full flex justify-center card items-center absolute top-0 !bg-surface-600/80">
      {#if ((pendingDisplay || pendingValidCheck) && !valid)}
        <SvgIcon
          class="h-10 w-10 animate-spin text-primary-600"
          icon={mdiLoading} />
      {:else if setupError}
        <div class="text-error-500">{setupError}</div>
      {:else}
        <div class="text-error-500">Failed to list directory</div>
      {/if}
    </div>
  {/if}
</div>
