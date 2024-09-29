<script lang="ts">
  import { mdiCheckCircle, mdiSync } from '@mdi/js';

  import Marquee from '$lib/components/Marquee.svelte';
  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import T from '$lib/components/T.svelte';
  import { getModalStore } from '$lib/skeletonExtensions';
  import { checkForUpdates, progress, unignoredUpdates, updateCheckInProgress, updates } from '$lib/store/ficsitCLIStore';
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
  <Marquee class="flex-auto text-start">
    {#if $smmUpdate}
      <T defaultValue="SMM update available" keyName="updates.smm-update-available"/>
    {:else if $unignoredUpdates.length > 0}
      <T defaultValue={'{updates} mod {updates, plural, one {update} other {updates}} available'} keyName="updates.mod-update-available" params={{ updates: $unignoredUpdates.length }}/>
    {:else}
      <T defaultValue="No mod/SMM updates right now" keyName="updates.no-updates"/>
    {/if}
  </Marquee>
  <SvgIcon
    class="h-5 w-5"
    icon={mdiCheckCircle} />
</button>

<button
  class="btn w-full bg-surface-200-700-token px-4 h-8 text-sm"
  disabled={!!$progress || $updateCheckInProgress}
  on:click={checkForAllUpdates}>
  <Marquee class="flex-auto text-start">
    {#if $updateCheckInProgress}
      <T defaultValue="Checking for updates..." keyName="updates.checking-for-updates"/>
    {:else}
      <T defaultValue="Check for updates" keyName="updates.check-for-updates"/>
    {/if}
  </Marquee>
  <SvgIcon
    class="h-5 w-5 {$updateCheckInProgress ? 'update-check' : ''}"
    icon={mdiSync} />
</button>
