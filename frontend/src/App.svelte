<script lang="ts">
  import './_global.postcss';
  import { arrow, autoUpdate, computePosition, flip, offset, shift, size } from '@floating-ui/dom';
  import { Modal, initializeStores, storePopup } from '@skeletonlabs/skeleton';
  import { FormatIcu } from '@tolgee/format-icu';
  import { DevTools, FormatSimple, Tolgee, TolgeeProvider } from '@tolgee/svelte';
  import { setContextClient } from '@urql/svelte';

  import T from '$lib/components/T.svelte';
  import TitleBar from '$lib/components/TitleBar.svelte';
  import LeftBar from '$lib/components/left-bar/LeftBar.svelte';
  import ModDetails from '$lib/components/mod-details/ModDetails.svelte';
  import ErrorDetails from '$lib/components/modals/ErrorDetails.svelte';
  import ErrorModal from '$lib/components/modals/ErrorModal.svelte';
  import ExternalInstallMod from '$lib/components/modals/ExternalInstallMod.svelte';
  import { supportedProgressTypes } from '$lib/components/modals/ProgressModal.svelte';
  import { modalRegistry } from '$lib/components/modals/modalsRegistry';
  import ImportProfile from '$lib/components/modals/profiles/ImportProfile.svelte';
  import { isUpdateOnStart } from '$lib/components/modals/smmUpdate/smmUpdate';
  import ModsList from '$lib/components/mods-list/ModsList.svelte';
  import { initializeGraphQLClient } from '$lib/core/graphql';
  import { i18n } from '$lib/generated';
  import { getModalStore, initializeModalStore } from '$lib/skeletonExtensions';
  import { installs, invalidInstalls, progress } from '$lib/store/ficsitCLIStore';
  import { error, expandedMod, siteURL } from '$lib/store/generalStore';
  import { konami, language, updateCheckMode } from '$lib/store/settingsStore';
  import { smmUpdate, smmUpdateReady } from '$lib/store/smmUpdateStore';
  import { ExpandMod, UnexpandMod } from '$wailsjs/go/app/app';
  import { Environment, EventsOn } from '$wailsjs/runtime';

  initializeStores();
  initializeModalStore();

  storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow, size });
  
  const tolgee = Tolgee()
    .use(DevTools())
    .use(FormatSimple())
    .use(FormatIcu())
    .init({
      language: 'en',
      fallbackLanguage: 'en',
      
      apiUrl: import.meta.env.VITE_TOLGEE_API_URL,
      apiKey: import.meta.env.VITE_TOLGEE_API_KEY,

      staticData: i18n,
    });
  
  $: tolgee.changeLanguage($language);

  let frameless = false;
  Environment().then((env) => {
    if (env.platform === 'windows') {
      frameless = true;
    }
  });

  export let apiEndpointURL!: string;
  export let siteEndpointURL!: string;
  
  $: $siteURL = siteEndpointURL;

  setContextClient(initializeGraphQLClient(apiEndpointURL));

  let windowStateChanging = false;
  let windowExpanded = false;

  function setExpanded(value: boolean) {
    if (windowExpanded === value) return;
    if (windowStateChanging) return;
    windowStateChanging = true;
    const op = value ? ExpandMod : UnexpandMod;
    op().then(() => {
      // wait a bit to prevent flickering
      setTimeout(() => {
        windowExpanded = value;
        windowStateChanging = false;
      }, 100);
    });
  }

  $: setExpanded(!!$expandedMod);

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
  } else {
    invalidInstallsError = false;
    noInstallsError = false;
  }

  const modalStore = getModalStore();

  $: if($progress && supportedProgressTypes.includes($progress.action)) {
    modalStore.triggerUnique({
      type: 'component',
      component: 'progress',
      meta: {
        persistent: true,
      },
    }, true);
  }

  let checkStart: boolean | undefined = undefined;
  const updateCheckModeInit = updateCheckMode.isInit;
  $: if ($updateCheckModeInit && checkStart === undefined) {
    checkStart = $updateCheckMode === 'launch';
  }

  const smmUpdateInit = smmUpdate.isInit;
  $: if ($smmUpdateInit && checkStart) {
    checkStart = false;
    if ($smmUpdate) {
      $isUpdateOnStart = true;
      if ($smmUpdateReady) {
        modalStore.trigger({
          type: 'component',
          component: 'smmUpdateReady',
          meta: {
            persistent: true,
          },
        });
      } else {
        modalStore.trigger({
          type: 'component',
          component: 'smmUpdateDownload',
          meta: {
            persistent: true,
          },
        });      
      }
    }
  }
  
  $: if ($smmUpdateReady && $updateCheckMode === 'ask') {
    modalStore.trigger({
      type: 'component',
      component: 'smmUpdateReady',
    });
  }

  $: if($error) {
    modalStore.trigger({
      type: 'component',
      component: {
        ref: ErrorModal,
        props: {
          error: $error,
        },
      },
    }, true);
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

<TolgeeProvider {tolgee}>
  <div class="flex flex-col h-screen w-screen select-none">
    {#if frameless}
      <TitleBar />
    {/if}
    <div class="flex grow h-0">
      <LeftBar />
      <div class="flex flex-auto @container/mod-list-wrapper z-[1]">
        <div class="{$expandedMod && !windowStateChanging ? 'w-2/5 hidden @3xl/mod-list-wrapper:block @3xl/mod-list-wrapper:flex-auto' : 'w-full'}" class:max-w-[42.5rem]={!!$expandedMod}>
          <ModsList
            hideMods={noInstallsError || invalidInstallsError}
            on:expandedMod={() => {
              focusOnEntry.focus();
            }}>
            <div class="card my-auto mr-4">
              <ErrorDetails
                error={''}
                fullPageMode={true}
                parent={{ onClose: () => {} }}
              >
                <!-- Svelte slots don't support dynamic passing, so this is kinda ugly. Switch to snippets once in svelte 5 https://github.com/sveltejs/svelte/issues/7651-->
                <T
                  slot="title"
                  defaultValue={noInstallsError ? 'No Satisfactory installs found' : '{invalidInstalls} invalid Satisfactory {invalidInstalls, plural, one {install} other {installs}} found'}
                  keyName={noInstallsError ? 'error.no_installs' : 'error.invalid_installs'}
                  params={{ invalidInstalls: $invalidInstalls.length }}
                />
              </ErrorDetails>
            </div>
          </ModsList>
        </div>
        <div class="w-3/5" class:grow={!windowStateChanging} class:hidden={!$expandedMod || windowStateChanging}>
          <ModDetails bind:focusOnEntry/>
        </div>
      </div>
    </div>
  </div>

  <Modal components={modalRegistry} />
</TolgeeProvider>

<!--
  skeleton modals don't provide a way to make them persistent (i.e. ignore mouse clicks outside and escape key)
  but we can capture the events and stop them if the modal has the persistent meta flag set, and the event would have closed the modal
-->
<svelte:window
  on:keydown|capture|nonpassive={modalKeyDown}
  on:mousedown|capture|nonpassive={modalMouseDown} />
