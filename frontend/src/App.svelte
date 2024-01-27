<script lang="ts">
  import './_global.postcss';
  import { setContextClient } from '@urql/svelte';
  import Dialog, { Actions, Content, Title } from '@smui/dialog';
  import Card, { Content as CardContent, PrimaryAction } from '@smui/card';
  import { ProgressBar, storePopup } from '@skeletonlabs/skeleton';
  import { computePosition, autoUpdate, offset, shift, flip, arrow, size } from '@floating-ui/dom';

  import TitleBar from '$lib/components/TitleBar.svelte';
  import ModsList from '$lib/components/mods-list/ModsList.svelte';
  import { initializeGraphQLClient } from '$lib/core/graphql';
  import { Environment } from '$wailsjs/runtime';
  import ModDetails from '$lib/components/mod-details/ModDetails.svelte';
  import { ExpandMod, UnexpandMod } from '$wailsjs/go/bindings/App';
  import LeftBar from '$lib/components/left-bar/LeftBar.svelte';
  import { installs, invalidInstalls, progress, selectedInstallMetadata, selectedProfile } from '$lib/store/ficsitCLIStore';
  import { konami } from '$lib/store/settingsStore';
  import { expandedMod, error, siteURL } from '$lib/store/generalStore';
  import { GenerateDebugInfo } from '$wailsjs/go/bindings/DebugInfo';
  import ExternalInstallMod from '$lib/components/ExternalInstallMod.svelte';

  storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow, size });

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

  let invalidInstallsError = false;
  let noInstallsError = false;
  let focusOnEntry: HTMLSpanElement;

  const installsInit = installs.isInit;
  const invalidInstallsInit = invalidInstalls.isInit;

  $: if($installsInit && $invalidInstallsInit && $installs.length === 0) {
    if($invalidInstalls.length > 0) {
      invalidInstallsError = true;
    } else {
      noInstallsError = true;
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
    <div class="flex w-0 grow">
      <div class="{$expandedMod ? 'max-w-[600px]' : ''} w-2/5 grow">
        <ModsList on:expandedMod={() => {
          focusOnEntry.focus();
        }} hideMods={noInstallsError || invalidInstallsError}>
          {#if noInstallsError}
            <Card variant="outlined" class="my-auto mr-2">
              <CardContent>
                <p class="font-bold text-2xl text-center">No Satisfactory installs found</p>
                <p class="text-lg text-center">
                  Seems wrong? Click the button below and send the generated zip file on the <a class="color-primary underline" href="https://discord.gg/xkVJ73E">modding discord</a> in #help-using-mods.
                </p>
              </CardContent>
              <PrimaryAction>
                <button
                  class="btn text-primary-600"
                  on:click={GenerateDebugInfo}>
                  Generate debug info
                </button>
              </PrimaryAction>
            </Card>
          {/if}
          {#if invalidInstallsError}
            <Card variant="outlined" class="my-auto mr-2">
              <CardContent>
                <p class="font-bold text-2xl text-center">{ $invalidInstalls.length } invalid Satisfactory installs found</p>
                <p class="text-lg text-center">
                  Seems wrong? Click the button below and send the generated zip file on the <a class="color-primary underline" href="https://discord.gg/xkVJ73E">modding discord</a> in #help-using-mods.
                </p>
              </CardContent>
              <PrimaryAction>
                <button
                  class="btn text-primary-600"
                  on:click={GenerateDebugInfo}>
                  Generate debug info
                </button>
              </PrimaryAction>
            </Card>
          {/if}
        </ModsList>
      </div>
      {#if $expandedMod}
        <div class:grow={!pendingExpand} class="{pendingExpand ? 'w-0' : 'w-3/5'}">
          <ModDetails bind:focusOnEntry/>
        </div>
      {/if}
    </div>
  </div>
</div>

<ExternalInstallMod />

<Dialog
  bind:open={installProgress}
  scrimClickAction=""
  escapeKeyAction=""
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
>
  <Title>Checking install {$selectedInstallMetadata?.branch} ({$selectedInstallMetadata?.launcher}) - CL{$selectedInstallMetadata?.version}</Title>
  <Content>
    {#if $progress}
      <p>{$progress.message}</p>
      <ProgressBar value={$progress.progress === -1 ? undefined : $progress.progress} max={1} class="h-4 w-full" meter="bg-primary-600"/>
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
      <ProgressBar value={$progress.progress === -1 ? undefined : $progress.progress} max={1} class="h-4 w-full" meter="bg-primary-600"/>
    {/if}
  </Content>
</Dialog>

<Dialog
  bind:open={modsEnabledProgress}
  scrimClickAction=""
  escapeKeyAction=""
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
>
  <Title>Toggling mods</Title>
  <Content>
    {#if $progress}
      <p>{$progress.message}</p>
      <ProgressBar value={$progress.progress === -1 ? undefined : $progress.progress} max={1} class="h-4 w-full" meter="bg-primary-600"/>
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
      <ProgressBar value={$progress.progress === -1 ? undefined : $progress.progress} max={1} class="h-4 w-full" meter="bg-primary-600"/>
    {/if}
  </Content>
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
    <button
      class="btn"
      on:click={() => $error = null}>
      Close
    </button>
    <button
      class="btn text-primary-600"
      on:click={GenerateDebugInfo}>
      Generate debug info
    </button>
  </Actions>
</Dialog>
