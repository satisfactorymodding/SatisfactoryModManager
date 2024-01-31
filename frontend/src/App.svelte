<script lang="ts">
  import './_global.postcss';
  import { setContextClient } from '@urql/svelte';
  import { storePopup , initializeStores, Modal } from '@skeletonlabs/skeleton';
  import { computePosition, autoUpdate, offset, shift, flip, arrow, size } from '@floating-ui/dom';

  import TitleBar from '$lib/components/TitleBar.svelte';
  import ModsList from '$lib/components/mods-list/ModsList.svelte';
  import { initializeGraphQLClient } from '$lib/core/graphql';
  import { Environment, EventsOn } from '$wailsjs/runtime';
  import ModDetails from '$lib/components/mod-details/ModDetails.svelte';
  import { ExpandMod, UnexpandMod } from '$wailsjs/go/bindings/App';
  import LeftBar from '$lib/components/left-bar/LeftBar.svelte';
  import { installs, invalidInstalls, progress } from '$lib/store/ficsitCLIStore';
  import { konami } from '$lib/store/settingsStore';
  import { expandedMod, error, siteURL } from '$lib/store/generalStore';
  import { initializeModalStore, getModalStore } from '$lib/store/skeletonExtensions';
  import { GenerateDebugInfo } from '$wailsjs/go/bindings/DebugInfo';
  import ExternalInstallMod from '$lib/components/modals/ExternalInstallMod.svelte';
  import { modalRegistry } from '$lib/components/modals/modalsRegistry';
  import ErrorModal from '$lib/components/modals/ErrorModal.svelte';
  import ImportProfile from '$lib/components/modals/profiles/ImportProfile.svelte';
  import { supportedProgressTypes } from '$lib/components/modals/ProgressModal.svelte';

  initializeStores();
  initializeModalStore();

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

  const modalStore = getModalStore();

  $: (() => {
    const currentProgressModalIndex = $modalStore.findIndex((m) => m.component && m.component === 'progress');
    if ($progress && supportedProgressTypes.includes($progress.item) && currentProgressModalIndex === -1){
      // Add it at the beginning instead of the end so it's on top
      $modalStore = [{
        type: 'component',
        component: 'progress',
        meta: {
          persistent: true,
        },
      }, ...$modalStore];
    }
  })();

  $: if($error) {
    $modalStore = [{
      type: 'component',
      component: {
        ref: ErrorModal,
        props: {
          error: $error,
        },
      },
    }, ...$modalStore];
    $error = null;
  }

  EventsOn('externalInstallMod', (modReference: string, version: string) => {
    if (!modReference) return;
    modalStore.trigger({
      type: 'component',
      component: {
        ref: ExternalInstallMod,
        props: {
          modReference,
          version,
        },
      },
    });
  });

  EventsOn('externalImportProfile', async (path: string) => {
    if (!path) return;
    modalStore.trigger({
      type: 'component',
      component: {
        ref: ImportProfile,
        props: {
          filepath: path,
        },
      },
    });
  });

  $: isPersistentModal = $modalStore.length > 0 && $modalStore[0].meta?.persistent;

  function modalMouseDown(event: MouseEvent) {
    if (!isPersistentModal) return;
    if (!(event.target instanceof Element)) return;
    const classList = event.target.classList;
    if (classList.contains('modal-backdrop') || classList.contains('modal-transition')) {
      event.stopImmediatePropagation();
    }
  }

  function modalKeyDown(event: KeyboardEvent) {
    if (!isPersistentModal) return;
    if (event.key === 'Escape') {
      event.stopImmediatePropagation();
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
  {#if frameless}
    <TitleBar />
  {/if}
  <div class="flex grow h-0">
    <LeftBar />
    <div class="flex w-0 grow">
      <div class="{$expandedMod ? 'max-w-[42.5rem]' : ''} w-2/5 grow">
        <ModsList
          hideMods={noInstallsError || invalidInstallsError}
          on:expandedMod={() => {
            focusOnEntry.focus();
          }}>
          <div class="card my-auto mr-4">
            <header class="card-header font-bold text-2xl text-center">
              {#if noInstallsError}
                No Satisfactory installs found
              {:else}
                {$invalidInstalls.length} invalid Satisfactory install{$invalidInstalls.length !== 1 ? 's' : ''} found
              {/if}
            </header>
            <section class="p-4">
              <p class="text-base text-center">
                Seems wrong? Click the button below and send the generated zip file on the <a class="text-primary-600 underline" href="https://discord.gg/xkVJ73E">modding discord</a> in #help-using-mods.
              </p>
            </section>
            <footer class="card-footer">
              <button
                class="btn text-primary-600 w-full"
                on:click={GenerateDebugInfo}>
                Generate debug info
              </button>
            </footer>
          </div>
        </ModsList>
      </div>
      {#if $expandedMod}
        <div class="{pendingExpand ? 'w-0' : 'w-3/5'}" class:grow={!pendingExpand}>
          <ModDetails bind:focusOnEntry/>
        </div>
      {/if}
    </div>
  </div>
</div>

<!--
  skeleton modals don't provide a way to make them persistent (i.e. ignore mouse clicks outside and escape key)
  but we can capture the events and stop them if the modal has the persistent meta flag set, and the event would have closed the modal
-->
<svelte:window
  on:keydown|capture|nonpassive={modalKeyDown}
  on:mousedown|capture|nonpassive={modalMouseDown} />
<Modal components={modalRegistry} />
