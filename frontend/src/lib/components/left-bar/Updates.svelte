<script lang="ts">
  import { mdiCheckCircle, mdiSync } from '@mdi/js';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { checkForUpdates, progress, unignoredUpdates, updateCheckInProgress, updates } from '$lib/store/ficsitCLIStore';
  import { getModalStore } from '$lib/store/skeletonExtensions';
  import { smmUpdate, smmUpdateReady } from '$lib/store/smmUpdateStore';
  import { CheckForUpdates as CheckForSMMUpdates } from '$wailsjs/go/autoupdate/autoUpdate';

  const modalStore = getModalStore();

  function checkForAllUpdates() {
    checkForUpdates().catch(console.error);
    if(!$smmUpdate || $smmUpdateReady) {
      CheckForSMMUpdates().catch(console.error);
    }
  }

  function showUpdateDialog() {
    if($smmUpdate) {
      if ($smmUpdateReady) {
        modalStore.trigger({
          type: 'component',
          component: 'smmUpdateReady',
        });
      } else {
        modalStore.trigger({
          type: 'component',
          component: 'smmUpdateDownload',
        });
      }
    } else if ($updates.length > 0) {
      modalStore.trigger({
        type: 'component',
        component: 'modUpdates',
      });
    }
  }
</script>

<button
  class="btn w-full bg-surface-200-700-token px-4 h-8 text-sm"
  class:!bg-primary-600={$smmUpdate || $unignoredUpdates.length > 0}
  on:click={() => showUpdateDialog()}>
  <span>
    {#if $smmUpdate}
      SMM update available
    {:else if $unignoredUpdates.length > 0}
      {$unignoredUpdates.length} mod updates available
    {:else}
      No mod/SMM updates right now
    {/if}
  </span>
  <div class="grow" />
  <SvgIcon
    class="h-5 w-5"
    icon={mdiCheckCircle} />
</button>

<button
  class="btn w-full bg-surface-200-700-token px-4 h-8 text-sm"
  disabled={!!$progress || $updateCheckInProgress}
  on:click={checkForAllUpdates}>
  <span>
    {#if $updateCheckInProgress}
      Checking for updates...
    {:else}
      Check for updates
    {/if}
  </span>
  <div class="grow" />
  <SvgIcon
    class="h-5 w-5 {$updateCheckInProgress ? 'update-check' : ''}"
    icon={mdiSync} />
</button>
