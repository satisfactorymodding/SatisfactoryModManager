<script lang="ts">
  import { mdiAlert, mdiCheckCircle, mdiCloseCircle, mdiDownload, mdiFolderOpen, mdiHelpCircle, mdiLoading, mdiPencil, mdiPlusCircle, mdiServerNetwork, mdiTrashCan, mdiUpload, mdiWeb } from '@mdi/js';
  import { type PopupSettings, popup } from '@skeletonlabs/skeleton';
  import _ from 'lodash';
  import { siDiscord, siGithub } from 'simple-icons/icons';

  import Tooltip from '../Tooltip.svelte';
  import DeleteProfile from '../modals/profiles/DeleteProfile.svelte';
  import RenameProfile from '../modals/profiles/RenameProfile.svelte';

  import LaunchButton from './LaunchButton.svelte';
  import Settings from './Settings.svelte';
  import Updates from './Updates.svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import Select from '$lib/components/Select.svelte';
  import { canChangeInstall, canModify, installs, installsMetadata, modsEnabled, profiles, selectedInstall, selectedProfile } from '$lib/store/ficsitCLIStore';
  import { error, siteURL } from '$lib/store/generalStore';
  import { getModalStore } from '$lib/store/skeletonExtensions';
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
</script>

<div class="flex flex-col h-full p-4 space-y-4 h-md:space-y-8 left-bar w-[22rem] min-w-[22rem] ">
  <div class="flex flex-col flex-auto h-full w-full space-y-4 h-md:space-y-8 overflow-y-auto">
    <div class="flex flex-col gap-2">
      <span class="pl-4 sticky top-0 z-[1] bg-surface-50-900-token">Game version</span>
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
          <span>
            {#if $installsMetadata[item]?.state === ficsitcli.InstallState.VALID}
              {$installsMetadata[item].info?.branch}{$installsMetadata[item].info?.type !== common.InstallType.WINDOWS ? ' - DS' : ''}
              ({$installsMetadata[item]?.info?.launcher})
            {:else if $installsMetadata[item]?.state === ficsitcli.InstallState.LOADING}
              Loading...
            {:else if $installsMetadata[item]?.state === ficsitcli.InstallState.INVALID}
              Invalid
            {:else}
              Unknown
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
                <span>Status: Loading...</span>
              {:else if $installsMetadata[item]?.state === ficsitcli.InstallState.INVALID}
                <span>Status: SMM cannot manage this install</span>
              {:else}
                <span>Status: Could not get information about this install</span>
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
        <svelte:fragment slot="selected" let:item>
          {#if $installsMetadata[item]?.state === ficsitcli.InstallState.VALID}
            {$installsMetadata[item].info?.branch}{$installsMetadata[item].info?.type !== common.InstallType.WINDOWS ? ' - DS' : ''}
            ({$installsMetadata[item]?.info?.launcher})
          {:else if $installsMetadata[item]?.state === ficsitcli.InstallState.LOADING}
            Loading...
          {:else if $installsMetadata[item]?.state === ficsitcli.InstallState.INVALID}
            Invalid
          {:else}
            Unknown
          {/if}
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
              Mods off
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
              Mods on
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
      <span class="pl-4 sticky top-0 z-[1] bg-surface-50-900-token">Profile</span>
      
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
            Add
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
            Rename
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
            Delete
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
            Import
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
            Export
          </span>
          <div class="grow"/>
          <SvgIcon
            class="h-5 w-5"
            icon={mdiUpload} />
        </button>
      </div>
    </div>
    <div class="flex flex-col gap-2">
      <span class="pl-4 sticky top-0 z-[1] bg-surface-50-900-token">Updates</span>
      <Updates />
    </div>
    <div class="flex flex-col gap-2">
      <span class="pl-4 sticky top-0 z-[1] bg-surface-50-900-token">Other</span>
      <button
        class="btn px-4 h-8 w-full text-sm bg-surface-200-700-token"
        on:click={() => modalStore.trigger({ type: 'component', component: 'serverManager' })}>
        <span>Manage Servers</span>
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
          Help
        </span>
        <div class="grow"/>
        <SvgIcon
          class="h-5 w-5"
          icon={mdiHelpCircle} />
      </button>
    </div>
    <div class="flex flex-col gap-2">
      <span class="pl-4 sticky top-0 z-[1] bg-surface-50-900-token">Links</span>
      <button
        class="btn w-full bg-surface-200-700-token px-4 h-8 text-sm"
        on:click={() => BrowserOpenURL($siteURL)}>
        <span>
          ficsit.app (Mod Repository)
        </span>
        <div class="grow" />
        <SvgIcon
          class="h-5 w-5"
          icon={mdiWeb} />
      </button>
      <button
        class="btn w-full bg-surface-200-700-token px-4 h-8 text-sm"
        on:click={() => BrowserOpenURL('https://discord.gg/xkVJ73E')}>
        <span>
          Satisfactory Modding Discord
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
          SMM GitHub
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
