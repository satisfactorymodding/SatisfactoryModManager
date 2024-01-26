<script lang="ts">
  import { Label } from '@smui/common';
  import Select, { Option } from '@smui/select';
  import Dialog, { Title, Content, Actions } from '@smui/dialog';
  import TextField from '@smui/textfield'; 
  import Tooltip, { Wrapper } from '@smui/tooltip';
  import { mdiAlert, mdiCheckCircle, mdiCloseCircle, mdiDownload, mdiFolderOpen, mdiHelpCircle, mdiPencil, mdiPlusCircle, mdiTrashCan, mdiUpload, mdiWeb } from '@mdi/js';
  import { siDiscord, siGithub } from 'simple-icons/icons';
  import HelperText from '@smui/textfield/helper-text';
  import LinearProgress from '@smui/linear-progress';

  import Settings from './Settings.svelte';
  import Updates from './Updates.svelte';
  import LaunchButton from './LaunchButton.svelte';
  import ServerManager from './ServerManager.svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { installs, profiles, canModify, selectedInstallMetadata, selectedInstall, selectedProfile, modsEnabled, progress, installsMetadata } from '$lib/store/ficsitCLIStore';
  import { error, siteURL } from '$lib/store/generalStore';
  import { BrowserOpenURL, EventsOn } from '$wailsjs/runtime/runtime';
  import { OpenExternal, OpenFileDialog } from '$wailsjs/go/bindings/App';
  import { common, type ficsitcli } from '$wailsjs/go/models';
  import { AddProfile, DeleteProfile, RenameProfile, ImportProfile, ExportCurrentProfile, ReadExportedProfileMetadata } from '$wailsjs/go/ficsitcli/FicsitCLI';
  
  const selectedInstallPathInit = selectedInstall.isInit;
  const selectedProfileInit = selectedProfile.isInit;

  async function installSelectChanged({ detail: { value } }: CustomEvent<{value?: string}>) {
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

  async function profileSelectChanged({ detail: { value } }: CustomEvent<{value?: string}>) {
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
</script>

<div class="flex flex-col h-full p-4 left-bar w-[22rem] min-w-[22rem] ">
  <div class="flex flex-col gap-2">
    <span class="pl-4">Game version</span>
    <Select
      variant="filled"
      menu$class="max-h-[32rem]"
      value={$selectedInstall}
      on:SMUISelect:change={installSelectChanged}
      ripple={false}
      disabled={!$canModify}
    >
      {#each $installs as install}
        <Option value={install}>
          <Label class="mdc-deprecated-list-item__text">
            {#if $installsMetadata[install]?.branch && $installsMetadata[install]?.type}
              {$installsMetadata[install]?.branch}{$installsMetadata[install]?.type !== common.InstallType.WINDOWS ? ' - DS' : ''}
            {:else}
              Unknown
            {/if}
            ({$installsMetadata[install]?.launcher})
          </Label>
          <div class="py-4 !m-0 !ml-auto !h-full" on:click={(e) => {
            e.stopPropagation();
            OpenExternal(install);
          }}>
            {#if $installsMetadata[install]?.branch && $installsMetadata[install]?.type}
              <Wrapper>
                <SvgIcon icon={mdiFolderOpen} class="!w-full !h-full"/>
                <Tooltip surface$class="max-w-lg text-base">
                  {install}
                </Tooltip>
              </Wrapper>
            {:else}
              <Wrapper>
                <SvgIcon icon={mdiAlert} class="!w-full !h-full text-red-500" />
                <Tooltip surface$class="max-w-lg text-base">
                  Failed to connect to server, retry connection in the server manager
                </Tooltip>
              </Wrapper>
            {/if}
          </div>
        </Option>
      {/each}
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
      variant="filled"
      menu$class="max-h-[32rem]"
      value={$selectedProfile}
      on:SMUISelect:change={profileSelectChanged}
      ripple={false}
      disabled={!$canModify}
    >
      {#each $profiles as profile}
        <Option value={profile}>
          <Label>{profile}</Label>
        </Option>
      {/each}
    </Select>
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
      <LinearProgress progress={$progress.progress} indeterminate={$progress.progress === -1} class="h-4 w-full rounded-lg"/>
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
