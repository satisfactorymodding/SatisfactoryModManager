<script lang="ts">
  import '@tolgee/svelte'; // Import so that the tolgee cli parses this file

  import { newCacheLocation } from './cacheLocationPicker';

  import T from '$lib/components/T.svelte';
  import { cacheDir } from '$lib/store/settingsStore';
  import { OpenDirectoryDialog } from '$wailsjs/go/app/app';

  export let parent: { onClose: () => void };

  let cacheError: string | null = null;
  
  let fileDialogOpen = false;
  async function pickCacheLocation() {
    if(fileDialogOpen) {
      return;
    }
    fileDialogOpen = true;
    try {
      let result = await OpenDirectoryDialog({
        defaultDirectory: $newCacheLocation ?? undefined,
      });
      if (result) {
        $newCacheLocation = result;
      }
    } catch (e) {
      if(e instanceof Error) {
        cacheError = e.message;
      } else if (typeof e === 'string') {
        cacheError = e;
      } else {
        cacheError = 'Unknown error';
      }
    } finally {
      fileDialogOpen = false;
    }
  }

  let cacheMoveInProgress = false;

  async function setCacheLocation() {
    try {
      cacheMoveInProgress = true;
      await cacheDir.asyncSet($newCacheLocation ?? '');
      cacheError = null;
    } catch(e) {
      if (e instanceof Error) {
        cacheError = e.message;
      } else if (typeof e === 'string') {
        cacheError = e;
      } else {
        cacheError = 'Unknown error';
      }
    } finally {
      cacheMoveInProgress = false;
    }
  }

  async function resetCacheLocation() {
    try {
      cacheMoveInProgress = true;
      await cacheDir.asyncSet('');
      cacheError = null;
    } catch(e) {
      if (e instanceof Error) {
        cacheError = e.message;
      } else if (typeof e === 'string') {
        cacheError = e;
      } else {
        cacheError = 'Unknown error';
      }
    } finally {
      cacheMoveInProgress = false;
    }
  }
</script>

<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="w-[60rem] card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    <T defaultValue="Change download cache location" keyName="settings.cache.title" />
  </header>
  <section class="p-4 grow">
    <label class="label">
      <span><T defaultValue="Cache location" keyName="settings.cache.cache-location" /></span>
      <div class="flex items-baseline">
        <div class="grow">
          <input
            class="input px-4 py-2 hover:!cursor-pointer"
            class:input-error={cacheError}
            readonly
            type="text"
            value={$newCacheLocation}
            on:click={() => pickCacheLocation()}
          />
          <p>
            {#if cacheError}
              {cacheError}
            {/if}
          </p>
        </div>
        <button
          class="btn mr-4 shrink-0 text-primary-600"
          disabled={cacheMoveInProgress}
          on:click={() => resetCacheLocation()}>
          <span><T defaultValue="Reset to default" keyName="settings.cache.reset" /></span>
          <div class="grow" />
        </button>
        <button
          class="btn shrink-0 text-primary-600"
          disabled={cacheMoveInProgress}
          on:click={() => setCacheLocation()}>
          <span><T defaultValue="Save and move" keyName="settings.cache.save" /></span>
          <div class="grow" />
        </button>
      </div>
    </label>
  </section>
  <footer class="card-footer">
    <button
      class="btn text-primary-600 variant-ringed"
      disabled={cacheMoveInProgress}
      on:click={parent.onClose}>
      <span><T defaultValue="Close" keyName="common.close" /></span>
    </button>
  </footer>
</div>

