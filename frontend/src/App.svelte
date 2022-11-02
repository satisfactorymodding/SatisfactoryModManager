<script lang="ts">
  import './_global.postcss';
  import TitleBar from '$lib/components/TitleBar.svelte';
  import ModsList from '$lib/components/mods-list/ModsList.svelte';
  import { initializeGraphQLClient } from '$lib/core/graphql';
  import { setClient } from '@urql/svelte';
  import { Environment } from '$wailsjs/runtime/runtime';
  import ModDetails from '$lib/components/mod-details/ModDetails.svelte';
  import { ExpandMod, UnexpandMod } from '$wailsjs/go/bindings/App';
  import LeftBar from '$lib/components/left-bar/LeftBar.svelte';
  import { installs, invalidInstalls } from '$lib/store/ficsitCLIStore';
  import { konami } from '$lib/store/settingsStore';
  import { expandedMod, error } from '$lib/store/generalStore';
  import Dialog, { Actions, Content, Title } from '@smui/dialog';
  import Button, { Label } from '@smui/button';
  import { GenerateDebugInfo } from '$wailsjs/go/bindings/DebugInfo';

  Environment().then((env) => {
    if (env.buildType !== 'dev') {
      document.addEventListener('contextmenu', (event) => event.preventDefault());
    }
  });

  setClient(initializeGraphQLClient());

  let windowExpanded = false;

  $: if ($expandedMod) {
    ExpandMod().then(() => { windowExpanded = true; });
  } else {
    windowExpanded = false;
    setTimeout(() => {
      UnexpandMod();
    }, 100);
  }

  $: pendingExpand = $expandedMod && !windowExpanded;

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
  
  const code = [38, 38, 40, 40, 37, 39, 37, 39, 66, 65];
  const keyQueue: number[] = [];
  window.addEventListener('keydown', (event) => {
    keyQueue.push(event.keyCode);
    if (keyQueue.length > code.length) {
      keyQueue.shift();
    }
    if (keyQueue.length === code.length && keyQueue.every((val, idx) => code[idx] === val)) {
      $konami = !$konami;
    }
  });
</script>

<div class="flex flex-col h-screen w-screen select-none">
  <TitleBar />
  <div class="flex grow h-0">
    <LeftBar />
    <div class:normal={!$expandedMod || pendingExpand} class:compact={windowExpanded}>
      <ModsList bind:compact={modsListCompact}/>
    </div>
    {#if $expandedMod}
      <div class="grow">
        <ModDetails />
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

<Dialog
  open={!!$error}  
  scrimClickAction=""
  escapeKeyAction=""
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
>
  <Title>Error</Title>
  <Content>
    <p>{ $error }</p>
    <p class="pt-4">Seems wrong? Click the button below and send the generated zip file on the modding discord in #help-using-mods.</p>
  </Content>
  <Actions>
    <Button action="" on:click={() => $error = null}>
      <Label>Close</Label>
    </Button>
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
