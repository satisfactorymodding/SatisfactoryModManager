<script lang="ts">
  import './_global.postcss';
  import TitleBar from '$lib/components/TitleBar.svelte';
  import ModsList from '$lib/components/mods-list/ModsList.svelte';
  import { initializeGraphQLClient } from '$lib/core/graphql';
  import { setClient } from '@urql/svelte';
  import { Environment } from '$wailsjs/runtime/runtime';
  import ModDetails from '$lib/components/mod-details/ModDetails.svelte';
  import { ExpandMod, UnexpandMod } from '$wailsjs/go/bindings/App';
  import LeftBar from '$lib/components/LeftBar.svelte';
  import { installs, invalidInstalls } from '$lib/ficsitCLIStore';
  import Dialog, { Actions, Content, Title } from '@smui/dialog';
  import Button, { Label } from '@smui/button';
  import { GenerateDebugInfo } from '$wailsjs/go/bindings/DebugInfo';

  Environment().then((env) => {
    if (env.buildType !== 'dev') {
      document.addEventListener('contextmenu', (event) => event.preventDefault());
    }
  });

  setClient(initializeGraphQLClient());

  let selectedModId: string | null = null;

  let windowExpanded = false;

  $: if (selectedModId) {
    ExpandMod().then(() => { windowExpanded = true; });
  } else {
    windowExpanded = false;
    setTimeout(() => {
      UnexpandMod();
    }, 100);
  }

  $: pendingExpand = selectedModId && !windowExpanded;

  $: modsListCompact = windowExpanded;

  let invalidInstallsDialog = false;
  let noInstallsDialog = false;

  $: if(installs.isInit && $installs.length === 0) {
    if($invalidInstalls.length > 0) {
      invalidInstallsDialog = true;
    } else {
      noInstallsDialog = true;
    }
  }
</script>

<div class="flex flex-col h-screen w-screen">
  <TitleBar />
  <div class="flex grow h-0 select-none">
    <LeftBar />
    <div class:normal={!selectedModId || pendingExpand} class:compact={windowExpanded}>
      <ModsList bind:selectedMod={selectedModId} bind:compact={modsListCompact}/>
    </div>
    {#if selectedModId}
      <div class="grow">
        <ModDetails id={selectedModId} on:close={() => selectedModId = null}/>
      </div>
    {/if}
  </div>
</div>

<Dialog
  bind:open={invalidInstallsDialog}  
  scrimClickAction=""
  escapeKeyAction=""
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
>
  <Title>{ $invalidInstalls.length } invalid Satisfactory installs found</Title>
  <Content>
    {#each $invalidInstalls as invalidInstall}
      <span>"{ invalidInstall }"</span><br>
    {/each}
    <br>
    Seems wrong? Click the button below and send the generated zip file on the modding discord in #help-using-mods.
  </Content>
  <Actions>
    <Button action="" on:click={GenerateDebugInfo}>
      <Label>Generate debug info</Label>
    </Button>
  </Actions>
</Dialog>

<Dialog
  bind:open={noInstallsDialog}  
  scrimClickAction=""
  escapeKeyAction=""
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
>
  <Title>No Satisfactory installs found</Title>
  <Content>
    Seems wrong? Click the button below and send the generated zip file on the modding discord in #help-using-mods.
  </Content>
  <Actions>
    <Button action="" on:click={GenerateDebugInfo}>
      <Label>Generate debug info</Label>
    </Button>
  </Actions>
</Dialog>

<style>
  .normal {
    width: 610px !important;
    min-width: 610px !important;
  }
  .compact {
    width: 470px;
    min-width: 470px;
  }
</style>
