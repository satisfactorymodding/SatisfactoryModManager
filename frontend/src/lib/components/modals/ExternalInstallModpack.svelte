<script lang="ts">
  import { getContextClient, queryStore } from '@urql/svelte';

  import T from '$lib/components/T.svelte';
  import { GetModpackSummaryDocument } from '$lib/generated';
  import { addQueuedModAction } from '$lib/store/actionQueue';
  import { error } from '$lib/store/generalStore';
  import { offline } from '$lib/store/settingsStore';
  import { InstallModpackRelease } from '$wailsjs/go/ficsitcli/ficsitCLI';

  export let parent: { onClose: () => void };

  export let modpackReference: string;
  export let version: string;

  const client = getContextClient();

  $: modpackQuery = queryStore(
    {
      query: GetModpackSummaryDocument,
      client,
      pause: !!$offline,
      variables: {
        modpackID: modpackReference,
      },
    },
  );

  $: modpack = $modpackQuery.fetching ? null : $modpackQuery.data?.modpack;

  function install() {
    if (!modpack) return;
    const modpackName = modpack.name;
    const action = async () => (InstallModpackRelease(modpackReference, version, modpackName)).catch((e) => $error = e);
    const actionName = 'install';
    addQueuedModAction(
      modpackReference,
      actionName,
      action,
    );
    parent.onClose();
  }

  $: renderedLogo = modpack?.logo || 'https://ficsit.app/images/no_image.webp';
</script>

<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="w-[48rem] card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    <T defaultValue="Install modpack" keyName="external-install-modpack.title" />
  </header>
  <section class="p-4 overflow-y-auto">
    {#if modpack}
      <div class="flex">
        <div class="grow">
          <p>{modpack.name}</p>
          {#if version}
            <p><T defaultValue={'Version {version}'} keyName="external-install-modpack.version" params={{ version }} /></p>
          {:else}
            <p><T defaultValue="Latest version" keyName="external-install-modpack.latest-version" /></p>
          {/if}
          <p>{modpack.short_description}</p>
        </div>
        <img class="logo h-24 w-24 mx-2" alt="{modpack.name} Logo" src={renderedLogo} />
      </div>
    {:else if $modpackQuery.fetching}
      <p><T defaultValue="Loading..." keyName="common.loading" /></p>
    {:else if $modpackQuery.error}
      <p><T defaultValue="Error loading modpack details" keyName="external-install-modpack.error-loading" /></p>
    {/if}
  </section>
  <footer class="card-footer">
    <button
      class="btn text-primary-600 variant-ringed"
      on:click={install}>
      <T defaultValue="Install" keyName="external-install-modpack.install" />
    </button>
    <button
      class="btn"
      on:click={parent.onClose}>
      <T defaultValue="Cancel" keyName="common.cancel" />
    </button>
  </footer>
</div>
