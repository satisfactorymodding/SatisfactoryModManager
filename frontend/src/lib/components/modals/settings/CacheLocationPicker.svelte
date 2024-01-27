<script lang="ts">
  import { newCacheLocation } from './cacheLocationPicker';

  import { OpenDirectoryDialog } from '$lib/generated/wailsjs/go/bindings/App';
  import { cacheDir } from '$lib/store/settingsStore';

  export let parent: {onClose: () => void};

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

<div class="card flex flex-col gap-2 !min-w-[800px] min-h-[200px]" style="max-height: calc(100vh - 128px); max-width: calc(100vw - 128px);">
  <header class="card-header font-bold text-2xl text-center">
    Change download cache location
  </header>
  <section class="p-4 grow">
    <label class="label">
      <span>Cache location</span>
      <div class="flex items-baseline">
        <div class="grow">
          <input type="text"
            class="input px-4 py-2 hover:!cursor-pointer"
            class:input-error={cacheError}
            value={$newCacheLocation}
            readonly
            on:click={() => pickCacheLocation()}
          />
          <p>
            {#if cacheError }
              { cacheError }
            {/if}
          </p>
        </div>
        <button
          class="btn mr-4 shrink-0 text-primary-600"
          disabled={cacheMoveInProgress}
          on:click={() => resetCacheLocation()}>
          <span>Reset to default</span>
          <div class="grow" />
        </button>
        <button
          class="btn shrink-0 text-primary-600"
          disabled={cacheMoveInProgress}
          on:click={() => setCacheLocation()}>
          <span>Save and move</span>
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
      <span>Close</span>
    </button>
  </footer>
</div>

