<script lang="ts">
  import './_global.postcss';
  import { setContextClient } from '@urql/svelte';
  import Dialog, { Actions, Content, Title } from '@smui/dialog';
  import Button, { Label } from '@smui/button';
  import LinearProgress from '@smui/linear-progress';

  import TitleBar from '$lib/components/TitleBar.svelte';
  import ModsList from '$lib/components/mods-list/ModsList.svelte';
  import { initializeGraphQLClient } from '$lib/core/graphql';
  import { Environment } from '$wailsjs/runtime';
  import ModDetails from '$lib/components/mod-details/ModDetails.svelte';
  import { ExpandMod, UnexpandMod } from '$wailsjs/go/bindings/App';
  import LeftBar from '$lib/components/left-bar/LeftBar.svelte';
  import { installs, invalidInstalls, progress, selectedInstall, selectedProfile } from '$lib/store/ficsitCLIStore';
  import { konami } from '$lib/store/settingsStore';
  import { expandedMod, error, siteURL } from '$lib/store/generalStore';
  import { GenerateDebugInfo } from '$wailsjs/go/bindings/DebugInfo';
  import ExternalInstallMod from '$lib/components/ExternalInstallMod.svelte';

  let frameless = false;
  Environment().then((env) => {
    if (env.buildType !== 'dev') {
      document.addEventListener('contextmenu', (event) => event.preventDefault());
    }
    if (env.platform === 'windows') {
      frameless = true;
    }
  });

  export let apiEndpointURL!: string;
  export let siteEndpointURL!: string;
  
  $: $siteURL = siteEndpointURL;

  setContextClient(initializeGraphQLClient(apiEndpointURL));

  let windowExpanded = false;

  $: if ($expandedMod) {
    ExpandMod().then(() => {
      setTimeout(() => {
        windowExpanded = true;
      }, 100);
    });
  } else {
    windowExpanded = false;
    setTimeout(() => {
      UnexpandMod();
    }, 100);
  }

  $: pendingExpand = $expandedMod && !windowExpanded;

  let invalidInstallsDialog = false;
  let noInstallsDialog = false;
  let focusOnEntry: HTMLSpanElement;

  const installsInit = installs.isInit;
  const invalidInstallsInit = invalidInstalls.isInit;

  $: if($installsInit && $invalidInstallsInit && $installs.length === 0) {
    if($invalidInstalls.length > 0) {
      invalidInstallsDialog = true;
    } else {
      noInstallsDialog = true;
    }
  }

  $: installProgress = $progress?.item === '__select_install__';
  $: profileProgress = $progress?.item === '__select_profile__';
  $: modsEnabledProgress = $progress?.item === '__toggle_mods__';
  $: updateProgress = $progress?.item === '__update__';
  
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
  {#if frameless}
    <TitleBar />
  {/if}
  <div class="flex grow h-0">
    <LeftBar />
    <div class="grow w-1/2 min-w-[400px] md:min-w-[420px] lg:min-w-[445px] {$expandedMod ? 'max-w-[600px]' : ''}">
      <ModsList on:expandedMod={() => {
        focusOnEntry.focus();
      }}/>
    </div>
    {#if $expandedMod}
      <div class:grow={!pendingExpand} class:w-0={pendingExpand}>
        <ModDetails bind:focusOnEntry/>
      </div>
    {/if}
  </div>
</div>

<ExternalInstallMod />

<Dialog
  bind:open={installProgress}
  scrimClickAction=""
  escapeKeyAction=""
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
>
  <Title>Checking install {$selectedInstall?.branch} ({$selectedInstall?.launcher}) - CL{$selectedInstall?.version}</Title>
  <Content>
    {#if $progress}
      <p>{$progress.message}</p>
      <LinearProgress progress={$progress.progress} indeterminate={$progress.progress === -1} class="h-4 w-full rounded-lg"/>
    {/if}
  </Content>
</Dialog>

<Dialog
  bind:open={profileProgress}
  scrimClickAction=""
  escapeKeyAction=""
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
>
  <Title>Selecting profile {$selectedProfile}</Title>
  <Content>
    {#if $progress}
      <p>{$progress.message}</p>
      <LinearProgress progress={$progress.progress} indeterminate={$progress.progress === -1} class="h-4 w-full rounded-lg"/>
    {/if}
  </Content>
</Dialog>

<Dialog
  bind:open={modsEnabledProgress}
  scrimClickAction=""
  escapeKeyAction=""
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
>
  <Title>Togging mods</Title>
  <Content>
    {#if $progress}
      <p>{$progress.message}</p>
      <LinearProgress progress={$progress.progress} indeterminate={$progress.progress === -1} class="h-4 w-full rounded-lg"/>
    {/if}
  </Content>
</Dialog>

<Dialog
  bind:open={updateProgress}
  scrimClickAction=""
  escapeKeyAction=""
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
>
  <Title>Updating mods</Title>
  <Content>
    {#if $progress}
      <p>{$progress.message}</p>
      <LinearProgress progress={$progress.progress} indeterminate={$progress.progress === -1} class="h-4 w-full rounded-lg"/>
    {/if}
  </Content>
</Dialog>

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
    Seems wrong? Click the button below and send the generated zip file on the <a class="color-primary underline" href="https://discord.gg/xkVJ73E">modding discord</a> in #help-using-mods.
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
    Seems wrong? Click the button below and send the generated zip file on the <a class="color-primary underline" href="https://discord.gg/xkVJ73E">modding discord</a> in #help-using-mods.
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
    <p class="pt-4">Seems wrong? Click the button below and send the generated zip file on the <a class="color-primary underline" href="https://discord.gg/xkVJ73E">modding discord</a> in #help-using-mods.</p>
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
