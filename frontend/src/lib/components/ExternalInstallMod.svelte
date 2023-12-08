<script lang="ts">
  import Button from '@smui/button';
  import Dialog, { Actions, Content } from '@smui/dialog';
  import { getContextClient, queryStore } from '@urql/svelte';

  import { GetModSummaryDocument } from '$lib/generated';
  import { manifestMods } from '$lib/store/ficsitCLIStore';
  import { addQueuedModAction, queuedMods, removeQueuedModAction } from '$lib/store/actionQueue';
  import { offline } from '$lib/store/settingsStore';
  import { error } from '$lib/store/generalStore';
  import { InstallMod } from '$wailsjs/go/ficsitcli/FicsitCLI';
  import { EventsOn } from '$wailsjs/runtime/runtime';

  let modId: string | undefined;
  let version: string | undefined;

  const client = getContextClient();

  $: modQuery = queryStore(
    {
      query: GetModSummaryDocument,
      client,
      pause: !modId || !!$offline,
      variables: {
        modReference: modId ?? '',
      },
    },
  );

  $: mod = $modQuery.fetching ? null : $modQuery.data?.mod;
  
  $: queued = $queuedMods.some((q) => q.mod === modId);
  $: isInstalled = !!modId && modId in $manifestMods;

  EventsOn('externalInstallMod', (m: string, v: string) => {
    modId = m ? m : undefined;
    version = v ? v : undefined;
  });

  function cancel() {
    modId = undefined;
    version = undefined;
  }

  function install() {
    if(!modId) {
      return;
    }
    const modReference = modId;
    const action = async () => InstallMod(modReference).catch((e) => $error = e);
    const actionName = 'install';
    if(queued) {
      removeQueuedModAction(modReference);
      return;
    }
    modId = undefined;
    return addQueuedModAction(
      modReference,
      actionName,
      action,
    );
  }

  $: renderedLogo = mod?.logo || 'https://ficsit.app/images/no_image.webp';
</script>

<Dialog open={!!modId} scrimClickAction="" escapeKeyAction="">
  <Content>
    {#if mod}
      <p class="text-lg">Install mod</p>
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
        <img src={renderedLogo} alt="{mod.name} Logo" class="logo h-24 w-24 mx-2" />
      </div>
    {:else if $modQuery.fetching}
      <p>Loading...</p>
    {:else if $modQuery.error}
      <p>Error loading mod details</p>
    {/if}
  </Content>
  <Actions>
    <Button on:click={cancel}>Cancel</Button>
    <Button on:click={install} disabled={isInstalled}>{ isInstalled ? 'Already installed' : 'Install' }</Button>
  </Actions>
</Dialog>
