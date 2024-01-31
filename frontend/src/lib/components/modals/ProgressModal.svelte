<script context="module" lang="ts">
  export const supportedProgressTypes = [
    '__select_install__',
    '__select_profile__',
    '__toggle_mods__',
    '__update__',
    '__import_profile__',
  ];
</script>

<script lang="ts">
  import { ProgressBar } from '@skeletonlabs/skeleton';

  import { progress, selectedInstallMetadata, selectedProfile } from '$lib/store/ficsitCLIStore';

  export let parent: { onClose: () => void };

  $: if(!$progress) {
    parent.onClose();
  }

  let title = '';

  $: title = (() => {
    switch ($progress?.item) {
      case '__select_install__':
        return `Selecting install ${$selectedInstallMetadata?.branch} (${$selectedInstallMetadata?.launcher}) - CL${$selectedInstallMetadata?.version}`;
      case '__select_profile__':
        return `Selecting profile ${$selectedProfile}`;
      case '__toggle_mods__':
        return 'Toggling mods';
      case '__update__':
        return 'Updating mods';
      case '__import_profile__':
        return `Importing profile ${$selectedProfile}`;
    }
    return '';
  })();
</script>

<div style="width: 500px; max-width: calc(100vw - 32px);" class="card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    {title}
  </header>
  <section class="p-4">
    {#if $progress}
      <p>{$progress.message}</p>
      <ProgressBar
        class="h-4 w-full"
        max={1}
        meter="bg-primary-600"
        value={$progress.progress === -1 ? undefined : $progress.progress}/>
    {/if}
  </section>
</div>
