<script lang="ts">
  import Dialog, { Title, Content, Actions } from '@smui/dialog';
  import TextField from '@smui/textfield'; 
  import { mdiAlert, mdiCheckCircle, mdiCloseCircle, mdiDownload, mdiFolderOpen, mdiHelpCircle, mdiPencil, mdiPlusCircle, mdiTrashCan, mdiUpload, mdiWeb } from '@mdi/js';
  import { siDiscord, siGithub } from 'simple-icons/icons';
  import HelperText from '@smui/textfield/helper-text';
  import { popup, type PopupSettings, ProgressBar } from '@skeletonlabs/skeleton';

  import Tooltip from '../Tooltip.svelte';

  import Settings from './Settings.svelte';
  import Updates from './Updates.svelte';
  import LaunchButton from './LaunchButton.svelte';
  import ServerManager from './ServerManager.svelte';

  import Select from '$lib/components/Select.svelte';
  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { installs, profiles, canModify, selectedInstallMetadata, selectedInstall, selectedProfile, modsEnabled, progress, installsMetadata } from '$lib/store/ficsitCLIStore';
  import { error, siteURL } from '$lib/store/generalStore';
  import { BrowserOpenURL, EventsOn } from '$wailsjs/runtime/runtime';
  import { OpenExternal, OpenFileDialog } from '$wailsjs/go/bindings/App';
  import { common, type ficsitcli } from '$wailsjs/go/models';
  import { AddProfile, DeleteProfile, RenameProfile, ImportProfile, ExportCurrentProfile, ReadExportedProfileMetadata } from '$wailsjs/go/ficsitcli/FicsitCLI';
  
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
    if ($selectedInstallMetadata) {
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
  }

  let addProfileDialog = false;
  let newProfileName = '';
  async function finishAddProfile() {
    try {
      await AddProfile(newProfileName);
      await selectedProfile.asyncSet(newProfileName);
      newProfileName = '';
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

  let renameProfileDialog = false;
  let renameOldProfileName = '';
  let renameNewProfileName = '';
  async function finishRenameProfile() {
    try {
      await RenameProfile(renameOldProfileName, renameNewProfileName);
      renameOldProfileName = '';
      renameNewProfileName = '';
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

  let deleteProfileDialog = false;
  let deleteProfileName = '';
  async function finishDeleteProfile() {
    try {
      await DeleteProfile(deleteProfileName);
      deleteProfileName = '';
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

  let importProfileDialog = false;
  let importProfileName = '';
  let importProfileFilepath = '';
  let importProfileError = '';
  let fileDialogOpen = false;
  let importProfileMetadata: ficsitcli.ExportedProfileMetadata | null = null;
  async function pickImportProfileFile() {
    if(fileDialogOpen) {
      return;
    }
    fileDialogOpen = true;
    try {
      importProfileFilepath = await OpenFileDialog({
        filters: [
          {
            displayName: 'SMM Profile (*.smmprofile)',
            pattern: '*.smmprofile',
          },
        ],
      });
      importProfileMetadata = await ReadExportedProfileMetadata(importProfileFilepath);
    } catch (e) {
      fileDialogOpen = false;
      if(e instanceof Error) {
        importProfileError = e.message;
      } else if (typeof e === 'string') {
        importProfileError = e;
      } else {
        importProfileError = 'Unknown error';
      }
    }
    fileDialogOpen = false;
  }
  async function finishImportProfile() {
    try {
      await ImportProfile(importProfileName, importProfileFilepath);
    } catch(e) {
      if(e instanceof Error) {
        importProfileError = e.message;
      } else if (typeof e === 'string') {
        importProfileError = e;
      } else {
        importProfileError = 'Unknown error';
      }
    }
  }
  EventsOn('externalImportProfile', async (path: string) => {
    importProfileFilepath = path;
    importProfileMetadata = await ReadExportedProfileMetadata(importProfileFilepath);
    importProfileDialog = true;
  });

  $: importProfileProgress = $progress?.item === '__import_profile__';

  function installPathPopupId(install: string) {
    return `install-path-${install.replace(/[^a-zA-Z0-9]/g, '-')}`;
  }

  $: installPathPopups = $installs.map((i) => [i, {
    event: 'hover',
    target: installPathPopupId(i),
    middleware: {
      offset: 4,
    },
    placement: 'right',
  } as PopupSettings]).reduce((acc, [k, v]) => ({ ...acc, [k as string]: v as PopupSettings }), {} as Record<string, PopupSettings>);

  function installWarningPopupId(install: string) {
    return `install-warning-${install}`;
  }

  $: installWarningPopups = $installs.map((i) => [i, {
    event: 'hover',
    target: installWarningPopupId(i),
    middleware: {
      offset: 4,
    },
    placement: 'right',
  } as PopupSettings]).reduce((acc, [k, v]) => ({ ...acc, [k as string]: v as PopupSettings }), {} as Record<string, PopupSettings>);
</script>

<div class="flex flex-col h-full p-4 left-bar w-[22rem] min-w-[22rem] ">
  <div class="flex flex-col gap-2">
    <span class="pl-4">Game version</span>
    <Select
      name="installsCombobox"
      class="w-full h-8"
      buttonClass="bg-surface-200-700-token px-4 text-sm"
      itemClass="bg-surface-50-900-token"
      itemActiveClass="!bg-surface-300/20"
      disabled={!$canModify}
      items={$installs}
      value={$selectedInstall ?? ''}
      on:change={installSelectChanged}
    >
      <svelte:fragment slot="item" let:item>
        <span>
          {#if $installsMetadata[item]?.branch && $installsMetadata[item]?.type}
            {$installsMetadata[item]?.branch}{$installsMetadata[item]?.type !== common.InstallType.WINDOWS ? ' - DS' : ''}
          {:else}
            Unknown
          {/if}
          ({$installsMetadata[item]?.launcher})
        </span>
      </svelte:fragment>
      <svelte:fragment slot="itemTrail" let:item>
        {#if $installsMetadata[item]?.branch && $installsMetadata[item]?.type}
          <Tooltip popupId={installPathPopupId(item)}>
            {item}
          </Tooltip>
          <div use:popup={installPathPopups[item]} class="!w-5 !h-5" on:click={(e) => {
            e.stopPropagation();
            OpenExternal(item);
          }} >
            <SvgIcon icon={mdiFolderOpen} class="!w-full !h-full"/>
          </div>
        {:else}
          <div use:popup={installWarningPopups[item]} class="!w-5 !h-5">
            <SvgIcon icon={mdiAlert} class="!w-full !h-full"/>
          </div>
          <Tooltip popupId={installWarningPopupId(item)}>
            {item}
          </Tooltip>
        {/if}
      </svelte:fragment>
      <svelte:fragment slot="selected" let:item>
        {#if $installsMetadata[item]?.branch && $installsMetadata[item]?.type}
          {$installsMetadata[item]?.branch}{$installsMetadata[item]?.type !== common.InstallType.WINDOWS ? ' - DS' : ''}
        {:else}
          Unknown
        {/if}
        ({$installsMetadata[item]?.launcher})
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
  <div class="flex flex-col gap-2 mt-4 h-md:mt-8">
    <span class="pl-4">Profile</span>
    
    <Select
      name="profileCombobox"
      class="w-full h-8"
      buttonClass="bg-surface-200-700-token px-4 text-sm"
      itemClass="bg-surface-50-900-token"
      itemActiveClass="!bg-surface-300/20"
      disabled={!$canModify}
      items={$profiles}
      value={$selectedProfile ?? ''}
      on:change={profileSelectChanged}
    />

    <div class="flex w-full gap-1">
      <button
        class="btn w-1/3 bg-surface-200-700-token px-4 h-8 text-sm"
        disabled={!$canModify}
        on:click={() => addProfileDialog = true}>
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
        on:click={() => { renameOldProfileName = $selectedProfile ?? ''; renameProfileDialog = true; }}>
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
        on:click={() => { deleteProfileName = $selectedProfile ?? ''; deleteProfileDialog = true; }}>
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
        on:click={() => importProfileDialog = true}
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
        on:click={() => { exportCurrentProfile(); }}
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
  <div class="flex flex-col gap-2 mt-4 h-md:mt-8">
    <span class="pl-4">Updates</span>
    <Updates />
  </div>
  <div class="flex flex-col gap-2 mt-4 h-md:mt-8">
    <span class="pl-4">Other</span>
    <ServerManager />
    <Settings />
    <button
      class="btn w-full bg-surface-200-700-token px-4 h-8 text-sm"
      disabled={!$canModify}
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
  <div class="flex flex-col gap-2 mt-4 h-md:mt-8">
    <span class="pl-4">Links</span>
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
  <div class="grow"/>
  <LaunchButton />
</div>

<Dialog
  bind:open={addProfileDialog}
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
>
  <Title>Add profile</Title>
  <Content>
    <TextField
      bind:value={newProfileName}
      label="Profile name"
      class="w-full"
    />
  </Content>
  <Actions>
    <button
      class="btn"
      on:click={() => addProfileDialog = false}>
      Cancel
    </button>
    <button
      class="btn text-primary-600"
      disabled={!newProfileName}
      on:click={finishAddProfile}>
      Add
    </button>
  </Actions>
</Dialog>

<Dialog
  bind:open={renameProfileDialog}
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
>
  <Title>Rename profile</Title>
  <Content>
    <TextField
      bind:value={renameOldProfileName}
      label="Old profile name"
      class="w-full"
      disabled
    />
    <TextField
      bind:value={renameNewProfileName}
      label="New profile name"
      class="w-full"
    />
  </Content>
  <Actions>
    <button
      class="btn"
      on:click={() => renameProfileDialog = false}>
      Cancel
    </button>
    <button
      class="btn text-primary-600"
      disabled={!renameNewProfileName}
      on:click={finishRenameProfile}>
      Rename
    </button>
  </Actions>
</Dialog>

<Dialog
  bind:open={deleteProfileDialog}
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
>
  <Title>Delete profile</Title>
  <Content>
    <TextField
      bind:value={deleteProfileName}
      label="Profile name"
      class="w-full"
      disabled
    />
  </Content>
  <Actions>
    <button
      class="btn"
      on:click={() => deleteProfileDialog = false}>
      Cancel
    </button>
    <button
      class="btn text-error-500"
      on:click={finishDeleteProfile}>
      Delete
    </button>
  </Actions>
</Dialog>

<Dialog
  bind:open={importProfileDialog}
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
>
  <Title>Import profile</Title>
  <Content>
    <TextField
      bind:value={importProfileName}
      label="Profile name"
      class="w-full"
    />
    <TextField
      bind:value={importProfileFilepath}
      invalid={!!importProfileError}
      label="Profile file"
      class="w-full"
      input$readonly
      on:click={() => pickImportProfileFile()}
    >
      <HelperText validationMsg slot="helper">
        { importProfileError }
      </HelperText>
    </TextField>
  </Content>
  <Actions>
    <button
      class="btn"
      on:click={() => importProfileDialog = false}>
      Cancel
    </button>
    <button
      class="btn text-primary-600"
      disabled={!importProfileName || !importProfileFilepath || !!importProfileError}
      on:click={finishImportProfile}>
      Import
    </button>
  </Actions>
</Dialog>

<Dialog
  bind:open={importProfileProgress}
  scrimClickAction=""
  escapeKeyAction=""
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
>
  <Title>Importing profile {importProfileName}</Title>
  <Content>
    {#if $progress}
      <p>{$progress.message}</p>
      <ProgressBar value={$progress.progress === -1 ? undefined : $progress.progress} max={1} class="h-4 w-full" meter="bg-primary-600"/>
    {/if}
  </Content>
</Dialog>

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
