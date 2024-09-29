<script context="module" lang="ts">
  import { ficsitcli } from '$wailsjs/go/models';

  export const supportedProgressTypes = [
    // ficsitcli.Action.SELECT_INSTALL, // This is instant, so no need for a progress modal
    ficsitcli.Action.SELECT_PROFILE,
    ficsitcli.Action.TOGGLE_MODS,
    ficsitcli.Action.UPDATE,
    ficsitcli.Action.IMPORT_PROFILE,
    ficsitcli.Action.APPLY,
  ];
</script>

<script lang="ts">
  import { ProgressBar } from '@skeletonlabs/skeleton';

  import { getModalStore } from '$lib/skeletonExtensions';
  import { progress, progressMessage, progressPercent, progressTitle } from '$lib/store/ficsitCLIStore';

  // Skeleton passes the parent prop to the modal component, and we would get a warning if the prop is not present here
  export let parent: { onClose: () => void };

  // Just so that it's not unused
  $: parent;

  const modalStore = getModalStore();

  $: if(!$progress) {
    // We cannot use parent.onClose because we might not be the top modal
    // Also this can get triggered multiple times for some reason,
    // which would cause an error in skeleton, so the modal would not actually be closed
    close();
  }

  let closed = false;

  function close() {
    if (closed) {
      return;
    }
    closed = true;
    modalStore.close('progress');
  }
</script>

<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="w-[48rem] card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    {$progressTitle}
  </header>
  <section class="p-4">
    {#if $progress}
      <p>{$progressMessage}</p>
      <ProgressBar
        class="h-4 w-full"
        max={1}
        meter="bg-primary-600"
        value={$progressPercent}/>
    {/if}
  </section>
</div>
