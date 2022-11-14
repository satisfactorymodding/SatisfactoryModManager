<script lang="ts">
  import { GetChangelogDocument } from '$lib/generated';
  import type { ficsitcli_bindings } from '$wailsjs/go/models';
  import Dialog, { Content, Title } from '@smui/dialog';
  import { getContextClient, queryStore } from '@urql/svelte';
  import { gt, lte } from 'semver';

  export let update: ficsitcli_bindings.Update | null = null;

  const client = getContextClient();

  $: modVersionChangelogStore = queryStore({
    query: GetChangelogDocument,
    client,
    pause: !update,
    variables: {
      modReference: update?.item ?? '',
    },
  });

  $: versions = $modVersionChangelogStore.data?.getModByReference?.versions;

  $: changelogs = versions ? versions.filter((v) => update && gt(v.version, update.currentVersion) && lte(v.version, update.newVersion)) : [];
</script>

<Dialog
  open={!!update}
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
  class="updates-dialog"
  >
  <Title>Changelog {update?.item} {update?.currentVersion} -> {update?.newVersion}</Title>
  <Content>
    {#each changelogs as changelog}
      <div class="text-lg font-semibold">v{changelog.version}</div>
      <div class="text-base">{changelog.changelog}</div>
    {/each}
  </Content>
</Dialog>