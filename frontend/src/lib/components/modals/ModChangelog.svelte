<script lang="ts">
  import { getContextClient, queryStore } from '@urql/svelte';
  import { gt, lte } from 'semver';

  import Markdown from '$lib/components/Markdown.svelte';
  import T, { translationElementPart } from '$lib/components/T.svelte';
  import { GetChangelogDocument } from '$lib/generated';
  import { offline } from '$lib/store/settingsStore';

  export let parent: { onClose: () => void };

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

<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="w-[48rem] card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    <T 
      defaultValue={'<1>{mod}</1> Changelog'}
      keyName="mod-changelog.title"
      params={{ mod: modData?.name ?? ' ' }}
      parts={[
        translationElementPart('span', { class: !modData?.name ? 'animate-pulse placeholder' : '' }),
      ]}
    />
  </header>
  <section class="p-4 overflow-y-auto">
    {#each changelogs as changelog}
      <div class="text-lg font-semibold">v{changelog.version}</div>
      <Markdown markdown={changelog.changelog}/>
      <hr />
    {/each}
  </section>
  <footer class="card-footer">
    <button
      class="btn"
      on:click={parent.onClose}>
      <T defaultValue="Close" keyName="common.close"/>
    </button>
  </footer>
</div>
