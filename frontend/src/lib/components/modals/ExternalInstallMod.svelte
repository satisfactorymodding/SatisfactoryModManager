<script lang="ts">
  import { getContextClient, queryStore } from '@urql/svelte';

  import T from '$lib/components/T.svelte';
  import { GetModSummaryDocument } from '$lib/generated';
  import { addQueuedModAction, queuedMods } from '$lib/store/actionQueue';
  import { manifestMods } from '$lib/store/ficsitCLIStore';
  import { error } from '$lib/store/generalStore';
  import { offline } from '$lib/store/settingsStore';
  import { InstallMod, InstallModVersion } from '$wailsjs/go/ficsitcli/ficsitCLI';

  export let parent: { onClose: () => void };

  export let modReference: string;
  export let version: string | undefined;

  const client = getContextClient();

  $: modQuery = queryStore(
    {
      query: GetModSummaryDocument,
      client,
      pause: !!$offline,
      variables: {
        modReference,
      },
    },
  );

  $: mod = $modQuery.fetching ? null : $modQuery.data?.mod;
  
  $: queued = $queuedMods.some((q) => q.mod === modReference);
  $: isInstalled = !!modReference && modReference in $manifestMods;

  function install() {
    const action = async () => (version ? InstallModVersion(modReference, version) : InstallMod(modReference)).catch((e) => $error = e);
    const actionName = 'install';
    addQueuedModAction(
      modReference,
      actionName,
      action,
    );
    parent.onClose();
  }

  $: renderedLogo = mod?.logo || 'https://ficsit.app/images/no_image.webp';
</script>

<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="w-[48rem] card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    <T defaultValue="Install mod" keyName="external-install-mod.title" />
  </header>
  <section class="p-4 overflow-y-auto">
    {#if mod}
      <div class="flex">
        <div class="grow">
          <p>{mod.name}</p>
          {#if version}
            <p><T defaultValue={'Version {version}'} keyName="external-install-mod.version" params={{ version }} /></p>
          {:else}
            <p><T defaultValue="Latest version" keyName="external-install-mod.latest-version" /></p>
          {/if}
          <p>{mod.short_description}</p>
        </div>
        <img class="logo h-24 w-24 mx-2" alt="{mod.name} Logo" src={renderedLogo} />
      </div>
    {:else if $modQuery.fetching}
      <p><T defaultValue="Loading..." keyName="common.loading" /></p>
    {:else if $modQuery.error}
      <p><T defaultValue="Error loading mod details" keyName="external-install-mod.error-loading" /></p>
    {/if}
  </section>
  <footer class="card-footer">
    <button
      class="btn text-primary-600 variant-ringed"
      disabled={isInstalled || queued}
      on:click={install}>
      {#if queued}
        <T defaultValue="In queue" keyName="external-install-mod.in-queue" />
      {:else if isInstalled}
        <T defaultValue="Already installed" keyName="external-install-mod.already-installed" />
      {:else}
        <T defaultValue="Install" keyName="external-install-mod.install" />
      {/if}
    </button>
    <button
      class="btn"
      on:click={parent.onClose}>
      <T defaultValue="Cancel" keyName="common.cancel" />
    </button>
  </footer>
</div>
