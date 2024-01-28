<script lang="ts">
  import { ProgressBar, getModalStore } from '@skeletonlabs/skeleton';

  import { smmUpdate , smmUpdateProgress, smmUpdateReady } from '$lib/store/smmUpdateStore';
  import { markdown } from '$lib/utils/markdown';
  import { bytesToAppropriate, secondsToAppropriate } from '$lib/utils/dataFormats';

  export let parent: {onClose: () => void};

  const modalStore = getModalStore();

  function filterChangelog(changelog: string) {
    const changelogStart = changelog.indexOf('## Changelog');
    if (changelogStart === -1) {
      return changelog;
    }
    const startIndex = changelog.indexOf('\n', changelogStart + 1);
    return changelog.slice(startIndex);
  }

  $: eta = ($smmUpdateProgress && $smmUpdateProgress.speed) ? ($smmUpdateProgress.total - $smmUpdateProgress.downloaded) / $smmUpdateProgress.speed : -1;

  $: if($smmUpdateReady) {
    modalStore.trigger({
      type: 'component',
      component: 'smmUpdateReady',
    });
    parent.onClose();
  }
</script>

<div style="width: 500px; max-width: calc(100vw - 32px);" class="card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    SMM Update Available - {$smmUpdate?.newVersion}
  </header>
  <section class="p-4 grow">
    {#if !$smmUpdateReady}
      {#if $smmUpdateProgress}
        <div>Downloading in background</div>
        <ProgressBar
          class="h-4 w-full"
          max={$smmUpdateProgress.total}
          meter="bg-primary-600"
          value={$smmUpdateProgress.total ? $smmUpdateProgress.downloaded : undefined}/>
        <div class="text-base">Downloading update: {bytesToAppropriate($smmUpdateProgress.downloaded)} / {bytesToAppropriate($smmUpdateProgress.total)}, {bytesToAppropriate($smmUpdateProgress.speed)}/s, ETA {eta >= 0 ? secondsToAppropriate(eta) : 'N/A'}</div>
      {/if}
    {/if}
    <div class="max-h-[500px] overflow-y-auto">
      {#each $smmUpdate?.changelogs ?? [] as changelog}
        <div class="text-xl font-semibold">{changelog.version}</div>
        <!-- eslint-disable-next-line svelte/no-at-html-tags -->
        <div class="markdown-content">{@html markdown(filterChangelog(changelog.changelog))}</div>
      {/each}
    </div>
  </section>
  <footer class="card-footer">
    <button
      class="btn"
      on:click={parent.onClose}>
      Cancel
    </button>
  </footer>
</div>
