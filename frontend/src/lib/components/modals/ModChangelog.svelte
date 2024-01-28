<script lang="ts">
  import { getContextClient, queryStore } from '@urql/svelte';
  import { gt, lte } from 'semver';

  import { GetChangelogDocument } from '$lib/generated';
  import { offline } from '$lib/store/settingsStore';

  export let parent: {onClose: () => void};

  export let mod: string;
  export let versionRange: string | { from: string, to: string };

  const client = getContextClient();

  $: modVersionChangelogStore = queryStore({
    query: GetChangelogDocument,
    client,
    pause: !!$offline,
    variables: {
      modReference: mod,
    },
  });

  $: modData = $modVersionChangelogStore.data?.getModByReference;

  $: versions = ($offline === null || $offline) ? [] : modData?.versions;

  function isVersionInRange(version: string) {
    if(typeof versionRange === 'string') {
      return version === versionRange;
    } else {
      return gt(version, versionRange.from) && lte(version, versionRange.to);
    }
  }

  $: changelogs = versions ? versions.filter((v) => isVersionInRange(v.version)) : [];
</script>

<div style="width: 500px; max-width: calc(100vw - 32px);" class="card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    {modData?.name ?? 'Loading...'} Changelog
  </header>
  <section class="p-4">
    {#each changelogs as changelog}
      <div class="text-lg font-semibold">v{changelog.version}</div>
      <div class="text-base">{changelog.changelog}</div>
    {/each}
  </section>
  <footer class="card-footer">
    <button
      class="btn"
      on:click={parent.onClose}>
      Close
    </button>
  </footer>
</div>
