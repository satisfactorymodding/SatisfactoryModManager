<script lang="ts">
  import Dialog, { Actions, Content, Title } from '@smui/dialog';
  import LinearProgress from '@smui/linear-progress';

  import { smmUpdate, smmUpdateProgress, smmUpdateReady } from '$lib/store/smmUpdateStore';
  import { bytesToAppropriate, secondsToAppropriate } from '$lib/utils/dataFormats';
  import { markdown } from '$lib/utils/markdown';
  import { UpdateAndRestart } from '$wailsjs/go/bindings/Update';

  let updateAvailableDialogOpen = false;
  let updateReadyDialogOpen = false;

  $: if($smmUpdate) {
    updateAvailableDialogOpen = true;
  }

  $: if($smmUpdateReady) {
    updateReadyDialogOpen = true;
  }

  export function show() {
    if(!$smmUpdateReady) {
      updateAvailableDialogOpen = true;
    } else {
      updateReadyDialogOpen = true;
    }
  }

  function filterChangelog(changelog: string) {
    const changelogStart = changelog.indexOf('## Changelog');
    if (changelogStart === -1) {
      return changelog;
    }
    const startIndex = changelog.indexOf('\n', changelogStart + 1);
    return changelog.slice(startIndex);
  }

  $: eta = ($smmUpdateProgress && $smmUpdateProgress.speed) ? ($smmUpdateProgress.total - $smmUpdateProgress.downloaded) / $smmUpdateProgress.speed : -1;
</script>

<Dialog bind:open={updateAvailableDialogOpen}>
  <Title>SMM Update Available - {$smmUpdate?.newVersion}</Title>
  <Content>
    {#if !$smmUpdateReady}
      {#if $smmUpdateProgress}
        <div>Downloading in background</div>
        <LinearProgress progress={($smmUpdateProgress.downloaded) / ($smmUpdateProgress.total)} indeterminate={!$smmUpdateProgress.total} class="transition-none" />
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
  </Content>
  <Actions>
    <button
      class="btn"
      on:click={() => updateAvailableDialogOpen = false}>
      <span>Close</span>
      <div class="grow" />
    </button>
  </Actions>
</Dialog>

<Dialog bind:open={updateReadyDialogOpen}>
  <Title>SMM Update Ready - {$smmUpdate?.newVersion}</Title>
  <Content>
    <div class="text-base">Update ready to install</div>
  </Content>
  <Actions>
    <button
      class="btn"
      on:click={() => updateReadyDialogOpen = false}>
      <span>Update on Exit</span>
      <div class="grow" />
    </button>
    <button
      class="btn text-primary-600"
      on:click={() => UpdateAndRestart()}>
      <span>Update and Restart</span>
      <div class="grow" />
    </button>
  </Actions>
</Dialog>

<style>
  :global(.transition-none .mdc-linear-progress__bar) {
    transition-property: none;
  }
</style>
