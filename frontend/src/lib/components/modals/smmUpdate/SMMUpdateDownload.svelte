<script lang="ts">
  import { ProgressBar } from '@skeletonlabs/skeleton';
  import { compare } from 'semver';

  import { isUpdateOnStart } from './smmUpdate';

  import Markdown from '$lib/components/Markdown.svelte';
  import { getModalStore } from '$lib/skeletonExtensions';
  import { smmUpdate, smmUpdateProgress, smmUpdateProgressStats, smmUpdateReady } from '$lib/store/smmUpdateStore';
  import { bytesToAppropriate, secondsToAppropriate } from '$lib/utils/dataFormats';

  export let parent: { onClose: () => void };

  const modalStore = getModalStore();

  function filterChangelog(changelog: string) {
    const changelogStart = changelog.indexOf('## Changelog');
    if (changelogStart === -1) {
      return changelog;
    }
    const startIndex = changelog.indexOf('\n', changelogStart + 1);
    return changelog.slice(startIndex);
  }

  $: speed = $smmUpdateProgressStats?.speed ?? 0;
  $: eta = $smmUpdateProgressStats?.eta ?? 0;

  $: changelogs = $smmUpdate ? Object.entries($smmUpdate.changelogs).map(([version, changelog]) => ({ version, changelog })).sort((a, b) => -compare(a.version, b.version)) : [];

  $: if($smmUpdateReady) {
    modalStore.trigger({
      type: 'component',
      component: 'smmUpdateReady',
    });
    parent.onClose();
  }
</script>

<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="w-[48rem] card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    SMM Update Available - {$smmUpdate?.version}
  </header>
  {#if !$smmUpdateReady && $smmUpdateProgress}
    <section class="p-4">
      <div>Downloading in background</div>
      <ProgressBar
        class="h-4 w-full"
        max={$smmUpdateProgress.total}
        meter="bg-primary-600"
        value={$smmUpdateProgress.total ? $smmUpdateProgress.current : undefined}/>
      <div class="text-base">Downloading update: {bytesToAppropriate($smmUpdateProgress.current)} / {bytesToAppropriate($smmUpdateProgress.total)}, {bytesToAppropriate(speed)}/s, ETA {eta >= 0 ? secondsToAppropriate(eta) : 'soon™'}</div>
    </section>
  {/if}
  <section class="p-4 overflow-y-auto">
    {#each changelogs ?? [] as changelog}
      <div class="text-xl font-semibold">{changelog.version}</div>
      <Markdown markdown={filterChangelog(changelog.changelog)}/>
      <hr />
    {/each}
  </section>
  {#if !$isUpdateOnStart}
    <footer class="card-footer">
      <button
        class="btn"
        on:click={parent.onClose}>
        Close
      </button>
    </footer>
  {/if}
</div>
