<script lang="ts">
  import { mdiAlert, mdiCheckCircle, mdiCloseCircle, mdiDownload, mdiFolderOpen, mdiHelp, mdiHelpCircle, mdiLoading, mdiMonitor, mdiPencil, mdiPlusCircle, mdiServer, mdiServerNetwork, mdiTrashCan, mdiUpload, mdiWeb } from '@mdi/js';
  import { type PopupSettings, popup } from '@skeletonlabs/skeleton';
  import _ from 'lodash';
  import { siDiscord, siGithub } from 'simple-icons/icons';

  import LaunchButton from './LaunchButton.svelte';
  import Settings from './Settings.svelte';
  import Updates from './Updates.svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import Select from '$lib/components/Select.svelte';
  import T from '$lib/components/T.svelte';
  import Tooltip from '$lib/components/Tooltip.svelte';
  import DeleteProfile from '$lib/components/modals/profiles/DeleteProfile.svelte';
  import RenameProfile from '$lib/components/modals/profiles/RenameProfile.svelte';
  import { getModalStore } from '$lib/skeletonExtensions';
  import { queuedMods } from '$lib/store/actionQueue';
  import { canChangeInstall, canModify, installs, installsMetadata, modsEnabled, profiles, selectedInstall, selectedProfile } from '$lib/store/ficsitCLIStore';
  import { error, siteURL } from '$lib/store/generalStore';
  import { OpenExternal } from '$wailsjs/go/app/app';
  import { ExportCurrentProfile } from '$wailsjs/go/ficsitcli/ficsitCLI';
  import { common, ficsitcli } from '$wailsjs/go/models';
  import { BrowserOpenURL } from '$wailsjs/runtime/runtime';
  
  const modalStore = getModalStore();

  const selectedInstallPathInit = selectedInstall.isInit;
  const selectedProfileInit = selectedProfile.isInit;

  async function installSelectChanged({ detail: value }: CustomEvent<string>) {
    if (!value) {
      return;
    }
    if (!$selectedInstallPathInit) {
      return;
    }
    try {
      await selectedInstall.asyncSet(value);
    } catch(e) {
      if (e instanceof Error) {
        $error = e.message;
      } else if (typeof e === 'string') {
        $error = e;
      } else {
        $error = 'Unknown error';
      }
    }
  }

  async function profileSelectChanged({ detail: value }: CustomEvent<string>) {
    if (!value) {
      return;
    }
    if (!$selectedProfileInit) {
      return;
    }
    try {
      await selectedProfile.asyncSet(value);
    } catch(e) {
      if (e instanceof Error) {
        $error = e.message;
      } else if (typeof e === 'string') {
        $error = e;
      } else {
        $error = 'Unknown error';
      }
    }
  }

  async function setModsEnabled(enabled: boolean) {
    if (!$selectedInstallPathInit) {
      return;
    }
    try {
      await modsEnabled.asyncSet(enabled);
    } catch(e) {
      if (e instanceof Error) {
        $error = e.message;
      } else if (typeof e === 'string') {
        $error = e;
      } else {
        $error = 'Unknown error';
      }
    }
  }

  async function exportCurrentProfile() {
    try {
      await ExportCurrentProfile();
    } catch(e) {
      if (e instanceof Error) {
        $error = e.message;
      } else if (typeof e === 'string') {
        $error = e;
      } else {
        $error = 'Unknown error';
      }
    }
  }

  function installOptionPopupId(install: string) {
    return `install-path-${install.replace(/[^a-zA-Z0-9]/g, '-')}`;
  }

  $: installOptionPopups = $installs.map((i) => [i, {
    event: 'hover',
    target: installOptionPopupId(i),
    middleware: {
      offset: 4,
    },
    placement: 'right',
  } as PopupSettings]).reduce((acc, [k, v]) => ({ ...acc, [k as string]: v as PopupSettings }), {} as Record<string, PopupSettings>);


  function iconForInstallType(type?: common.InstallType) {
    switch (type) {
      case common.InstallType.WINDOWS:
        return mdiMonitor;
      case common.InstallType.WINDOWS_SERVER:
      case common.InstallType.LINUX_SERVER:
        return mdiServer;
      default:
        return mdiHelp;
    }
  }

  const queueLockoutPopupId = 'queue-lockout-popup';
  const queueLockoutPopup = {
    event: 'hover',
    target: queueLockoutPopupId,
    middleware: {
      offset: 4,
    },
    placement: 'bottom-end',
  } satisfies PopupSettings;
</script>

<div class="flex flex-col h-full p-4 space-y-4 h-md:space-y-8 left-bar w-[22rem] min-w-[22rem] ">
  <div class="flex flex-col flex-auto h-full w-full space-y-4 h-md:space-y-8 overflow-y-auto">
    <div class="h-md:space-y-8" use:popup={queueLockoutPopup}>

      <div class="flex flex-col gap-2">
        <span class="pl-4 sticky top-0 z-[1] bg-surface-50-900-token">
          <T defaultValue="Select Game Installation" keyName="left-bar.game-version"/>
        </span>
        <Select
          name="installsCombobox"
          class="w-full h-8"
          buttonClass="bg-surface-200-700-token px-4 text-sm"
          disabled={!$canChangeInstall}
          itemActiveClass="!bg-surface-300/20"
          itemClass="bg-surface-50-900-token"
          items={_.orderBy($installs)}
          value={$selectedInstall ?? ''}
          on:change={installSelectChanged}
        >
          <svelte:fragment slot="item" let:item>
            <span class="flex items-center min-w-0">
              {#if $installsMetadata[item]?.state === ficsitcli.InstallState.VALID}
                <SvgIcon class="!w-5 !h-5 mr-2 shrink-0" icon={iconForInstallType($installsMetadata[item].info?.type)}/>
                <span class="truncate min-w-0">
                  {$installsMetadata[item]?.info?.launcher}
                </span>
                <span class="shrink-0 ml-1">({#if $installsMetadata[item]?.info?.location === common.LocationType.LOCAL}{$installsMetadata[item]?.info?.branch}{:else}CL{$installsMetadata[item]?.info?.version}{/if})</span>
              {:else if $installsMetadata[item]?.state === ficsitcli.InstallState.LOADING}
                <T defaultValue="Loading..." keyName="left-bar.install-loading"/>
              {:else if $installsMetadata[item]?.state === ficsitcli.InstallState.INVALID}
                <T defaultValue="Invalid" keyName="left-bar.install-invalid"/>
              {:else}
                <T defaultValue="Unknown" keyName="left-bar.install-unknown"/>
              {/if}
            </span>
          </svelte:fragment>
          <svelte:fragment slot="itemTrail" let:item>
            <Tooltip popupId={installOptionPopupId(item)}>
              <div class="flex flex-col">
                <span>{item}</span>
                {#if $installsMetadata[item]?.state === ficsitcli.InstallState.VALID}
                  <!-- nothing extra -->
                {:else if $installsMetadata[item]?.state === ficsitcli.InstallState.LOADING}
                  <span>
                    <T defaultValue="Status: Loading..." keyName="left-bar.install-loading-tooltip"/>
                  </span>
                {:else if $installsMetadata[item]?.state === ficsitcli.InstallState.INVALID}
                  <span>
                    <T defaultValue="Status: SMM cannot manage this install" keyName="left-bar.install-invalid-tooltip"/>
                  </span>
                {:else}
                  <span>
                    <T defaultValue="Status: Could not get information about this install" keyName="left-bar.install-unknown-tooltip"/>
                  </span>
                {/if}
              </div>
            </Tooltip>
            <button
              class="!w-5 !h-5"
              on:click|stopPropagation={() => OpenExternal(item)}
              use:popup={installOptionPopups[item]} >
              {#if $installsMetadata[item]?.state === ficsitcli.InstallState.VALID}
                <SvgIcon class="!w-full !h-full" icon={mdiFolderOpen}/>
              {:else if $installsMetadata[item]?.state === ficsitcli.InstallState.LOADING}
                <SvgIcon class="!w-full !h-full animate-spin text-primary-600" icon={mdiLoading}/>
              {:else if $installsMetadata[item]?.state === ficsitcli.InstallState.INVALID}
                <SvgIcon class="!w-full !h-full" icon={mdiAlert}/>
              {:else}
                <SvgIcon class="!w-full !h-full" icon={mdiAlert}/>
              {/if}
            </button>
          </svelte:fragment>
        </Select>

        <div class="flex w-full">
          <div class="btn-group bg-surface-200-700-token w-full text-xl">
            <button
              class="w-1/2 !btn-sm !px-4"
              class:!bg-error-900={!$modsEnabled}
              disabled={!$canModify}
              on:click={() => setModsEnabled(false)}
            >
              <span>
                <T defaultValue="Mods off" keyName="left-bar.mods-off"/>
              </span>
              <div class="grow"/>
              <SvgIcon
                class="h-5 w-5"
                icon={mdiCloseCircle} />
            </button>
            <button
              class="w-1/2 !btn-sm !px-4"
              class:!bg-primary-900={$modsEnabled}
              disabled={!$canModify}
              on:click={() => setModsEnabled(true)}>
              <span>
                <T defaultValue="Mods on" keyName="left-bar.mods-on"/>
              </span>
              <div class="grow"/>
              <SvgIcon
                class="h-5 w-5"
                icon={mdiCheckCircle} />
            </button>
          </div>
        </div>
      </div>
      <div class="flex flex-col gap-2">
        <span class="pl-4 sticky top-0 z-[1] bg-surface-50-900-token">
          <T defaultValue="Profile" keyName="left-bar.profile"/>
        </span>
      
        <Select
          name="profileCombobox"
          class="w-full h-8"
          buttonClass="bg-surface-200-700-token px-4 text-sm"
          disabled={!$canModify}
          itemActiveClass="!bg-surface-300/20"
          itemClass="bg-surface-50-900-token"
          items={_.orderBy($profiles)}
          value={$selectedProfile ?? ''}
          on:change={profileSelectChanged}
        >
          <svelte:fragment slot="itemTrail" let:item>
            <button
              disabled={!$canModify}
              on:click|stopPropagation={() => modalStore.trigger({ type:'component', component: { ref: RenameProfile, props: { profile: item } } })}
            >
              <SvgIcon class="!w-5 !h-5 text-warning-500" icon={mdiPencil}/>
            </button>
            <button
              disabled={!$canModify || $profiles.length === 1}
              on:click|stopPropagation={() => modalStore.trigger({ type:'component', component: { ref: DeleteProfile, props: { profile: item } } })}
            >
              <SvgIcon class="!w-5 !h-5 text-error-700" icon={mdiTrashCan}/>
            </button>
          </svelte:fragment>
        </Select>
  
        <div class="flex w-full gap-1">
          <button
            class="btn w-1/3 bg-surface-200-700-token px-4 h-8 text-sm"
            disabled={!$canModify}
            on:click={() => modalStore.trigger({ type:'component', component: 'addProfile' })}>
            <span>
              <T defaultValue="Add" keyName="common.add"/>
            </span>
            <div class="grow"/>
            <SvgIcon
              class="h-5 w-5 text-primary-600"
              icon={mdiPlusCircle} />
          </button>
          <button
            class="btn w-1/3 bg-surface-200-700-token px-2 h-8 text-sm"
            disabled={!$canModify}
            on:click={() => modalStore.trigger({ type:'component', component: { ref: RenameProfile, props: { profile: $selectedProfile } } })}>
            <span>
              <T defaultValue="Rename" keyName="common.rename"/>
            </span>
            <div class="grow"/>
            <SvgIcon
              class="h-5 w-5 text-warning-500"
              icon={mdiPencil} />
          </button>
          <button
            class="btn w-1/3 bg-surface-200-700-token px-3 h-8 text-sm"
            disabled={!$canModify || $profiles.length === 1}
            on:click={() => modalStore.trigger({ type:'component', component: { ref: DeleteProfile, props: { profile: $selectedProfile } } })}>
            <span>
              <T defaultValue="Delete" keyName="common.delete"/>
            </span>
            <div class="grow"/>
            <SvgIcon
              class="h-5 w-5 text-error-700"
              icon={mdiTrashCan} />
          </button>
        </div>
        <div class="flex w-full gap-1">
          <button
            class="btn w-1/2 bg-surface-200-700-token px-4 h-8 text-sm"
            disabled={!$canModify}
            on:click={() => modalStore.trigger({ type: 'component', component: 'importProfile' })}
          >
            <span>
              <T defaultValue="Import" keyName="common.import"/>
            </span>
            <div class="grow"/>
            <SvgIcon
              class="h-5 w-5"
              icon={mdiDownload} />
          </button>
          <button
            class="btn w-1/2 bg-surface-200-700-token px-4 h-8 text-sm"
            disabled={!$canModify}
            on:click={() => exportCurrentProfile()}
          >
            <span>
              <T defaultValue="Export" keyName="left-bar.export"/>
            </span>
            <div class="grow"/>
            <SvgIcon
              class="h-5 w-5"
              icon={mdiUpload} />
          </button>
        </div>
      </div>

    </div>
    <Tooltip disabled={$queuedMods.length === 0} popupId={queueLockoutPopupId}>
      <span class="text-base">
        <T defaultValue={'You have {number} {number, plural, one {action} other {actions}} queued. Apply or cancel {number, plural, one {it} other {them}} before switching installs or profiles.'} keyName="left-bar.queue-blocking-switching-tooltip" params={{ number: $queuedMods.length }} />
      </span>
    </Tooltip>

    <div class="flex flex-col gap-2">
      <span class="pl-4 sticky top-0 z-[1] bg-surface-50-900-token">
        <T defaultValue="Updates" keyName="left-bar.updates"/>
      </span>
      <Updates />
    </div>
    <div class="flex flex-col gap-2">
      <span class="pl-4 sticky top-0 z-[1] bg-surface-50-900-token">
        <T defaultValue="Other" keyName="left-bar.other"/>
      </span>
      <button
        class="btn px-4 h-8 w-full text-sm bg-surface-200-700-token"
        on:click={() => modalStore.trigger({ type: 'component', component: 'serverManager' })}>
        <span>
          <T defaultValue="Manage Servers" keyName="left-bar.manage-servers"/>
        </span>
        <div class="grow" />
        <SvgIcon
          class="h-5 w-5"
          icon={mdiServerNetwork} />
      </button>
      <Settings />
      <button
        class="btn w-full bg-surface-200-700-token px-4 h-8 text-sm"
        on:click={() => BrowserOpenURL('https://docs.ficsit.app/satisfactory-modding/latest/ForUsers/SatisfactoryModManager.html')}
      >
        <span>
          <T defaultValue="Help" keyName="left-bar.help"/>
        </span>
        <div class="grow"/>
        <SvgIcon
          class="h-5 w-5"
          icon={mdiHelpCircle} />
      </button>
    </div>
    <div class="flex flex-col gap-2">
      <span class="pl-4 sticky top-0 z-[1] bg-surface-50-900-token">
        <T defaultValue="Links" keyName="left-bar.links"/>
      </span>
      <button
        class="btn w-full bg-surface-200-700-token px-4 h-8 text-sm"
        on:click={() => BrowserOpenURL($siteURL)}>
        <span>
          <T defaultValue="ficsit.app (Mod Repository)" keyName="left-bar.ficsit-app"/>
        </span>
        <div class="grow" />
        <SvgIcon
          class="h-5 w-5"
          icon={mdiWeb} />
      </button>
      <button
        class="btn w-full bg-surface-200-700-token px-4 h-8 text-sm"
        on:click={() => BrowserOpenURL('https://discord.ficsit.app/')}>
        <span>
          <T defaultValue="Satisfactory Modding Discord" keyName="left-bar.satisfactory-modding-discord"/>
        </span>
        <div class="grow" />
        <SvgIcon
          class="h-5 w-5"
          icon={siDiscord.path} />
      </button>
      <button
        class="btn w-full bg-surface-200-700-token px-4 h-8 text-sm"
        on:click={() => BrowserOpenURL('https://github.com/satisfactorymodding/SatisfactoryModManager')} >
        <span>
          <T defaultValue="SMM GitHub" keyName="left-bar.smm-github"/>
        </span>
        <div class="grow" />
        <SvgIcon
          class="h-5 w-5"
          icon={siGithub.path} />
      </button>
    </div>
  </div>
  <LaunchButton />
</div>

<style>
  :global(.update-check) {
    animation: spin 4s linear infinite;
  }
  @keyframes spin {
    100% {
      transform: rotate(-360deg);
    }
  }
</style>
