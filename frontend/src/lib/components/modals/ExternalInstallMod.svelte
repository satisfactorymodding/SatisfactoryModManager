<script lang="ts">
  import { getContextClient, queryStore } from '@urql/svelte';

  import { GetModSummaryDocument } from '$lib/generated';
  import { manifestMods } from '$lib/store/ficsitCLIStore';
  import { addQueuedModAction, queuedMods } from '$lib/store/actionQueue';
  import { offline } from '$lib/store/settingsStore';
  import { error } from '$lib/store/generalStore';
  import { InstallMod, InstallModVersion } from '$wailsjs/go/ficsitcli/FicsitCLI';

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

<div style="width: 500px; max-width: calc(100vw - 32px);" class="card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    Install mod
  </header>
  <section class="p-4">
    {#if mod}
      <div class="flex">
        <div class="grow">
          <p>{mod.name}</p>
          {#if version}
            <p>Version {version}</p>
          {:else}
            <p>Latest version</p>
          {/if}
          <p>{mod.short_description}</p>
        </div>
        <img class="logo h-24 w-24 mx-2" alt="{mod.name} Logo" src={renderedLogo} />
      </div>
    {:else if $modQuery.fetching}
      <p>Loading...</p>
    {:else if $modQuery.error}
      <p>Error loading mod details</p>
    {/if}
  </section>
  <footer class="card-footer">
    <button
      class="btn text-primary-600 variant-ringed"
      disabled={isInstalled || queued}
      on:click={install}>
      {#if queued}
        In queue
      {:else if isInstalled}
        Already installed
      {:else}
        Install
      {/if}
    </button>
    <button
      class="btn"
      on:click={parent.onClose}>
      Cancel
    </button>
  </footer>
</div>
