<script lang="ts">
  import Button, { Label } from '@smui/button';
  import Select, { Option } from '@smui/select';
  import Dialog, { Title, Content, Actions } from '@smui/dialog';
  import TextField from '@smui/textfield'; 
  import Tooltip, { Wrapper } from '@smui/tooltip';
  import { mdiCheckCircle, mdiCloseCircle, mdiDownload, mdiHelpCircle, mdiPencil, mdiPlusCircle, mdiTrashCan, mdiUpload, mdiWeb } from '@mdi/js';
  import { siDiscord, siGithub } from 'simple-icons/icons';
  import HelperText from '@smui/textfield/helper-text';
  import LinearProgress from '@smui/linear-progress';

  import Settings from './Settings.svelte';
  import Updates from './Updates.svelte';
  import LaunchButton from './LaunchButton.svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { installs, profiles, canModify, selectedInstall, selectedInstallPath, selectedProfile, modsEnabled, progress } from '$lib/store/ficsitCLIStore';
  import { error, siteURL } from '$lib/store/generalStore';
  import { BrowserOpenURL, EventsOn } from '$wailsjs/runtime/runtime';
  import { OpenFileDialog } from '$wailsjs/go/bindings/App';
  import type { ficsitcli } from '$wailsjs/go/models';
  import { AddProfile, DeleteProfile, RenameProfile, ImportProfile, ExportCurrentProfile, ReadExportedProfileMetadata } from '$wailsjs/go/ficsitcli/FicsitCLI';
  
  const selectedInstallPathInit = selectedInstallPath.isInit;
  const selectedProfileInit = selectedProfile.isInit;

  async function installSelectChanged({ detail: { value } }: CustomEvent<{value?: string}>) {
    if (!value) {
      return;
    }
    if (!$selectedInstallPathInit) {
      return;
    }
    try {
      await selectedInstallPath.asyncSet(value);
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
    if ($selectedInstall) {
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
  <div class="flex flex-col">
    <span class="pl-4">Game version</span>
    <Select
      variant="filled"
      class="mt-2"
      menu$class="max-h-[32rem]"
      value={$selectedInstallPath}
      on:SMUISelect:change={installSelectChanged}
      ripple={false}
      disabled={!$canModify}
    >
      {#each $installs as install}
        <Wrapper>
          <Option value={install.path}>
            <Label>{install?.branch} ({install?.launcher})</Label>
          </Option>
          
          <Tooltip surface$class="max-w-lg text-base">
            {install?.path}
          </Tooltip>
        </Wrapper>
      {/each}
    </Select>
    <div class="flex w-full mt-2">
      <Button variant="unelevated" class="w-1/2 rounded-tr-none rounded-br-none mods-toggle-button {$modsEnabled ? '' : 'mods-off'}" on:click={() => setModsEnabled(false)} disabled={!$canModify}>
        <Label>
          Mods off
        </Label>
        <div class="grow"/>
        <SvgIcon icon={mdiCloseCircle} class="h-5 w-5" />
      </Button>
      <Button variant="unelevated" class="w-1/2 rounded-tl-none rounded-bl-none mods-toggle-button {$modsEnabled ? 'mods-on' : ''}" on:click={() => setModsEnabled(true)} disabled={!$canModify}>
        <Label>
          Mods on
        </Label>
        <div class="grow"/>
        <SvgIcon icon={mdiCheckCircle} class="h-5 w-5" />
      </Button>
    </div>
  </div>
  <div class="flex flex-col mt-4 h-md:mt-8">
    <span class="pl-4">Profile</span>
    <Select
      variant="filled"
      class="mt-2"
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
    <div class="flex w-full mt-2">
      <Button class="w-1/3 pr-2 pl-5 profile-add" on:click={() => addProfileDialog = true} disabled={!$canModify}>
        <Label>
          Add
        </Label>
        <div class="grow"/>
        <SvgIcon icon={mdiPlusCircle} />
      </Button>
      <Button class="w-1/3 mx-2 pr-0 profile-edit" on:click={() => { renameOldProfileName = $selectedProfile ?? ''; renameProfileDialog = true; }} disabled={!$canModify}>
        <Label>
          Rename
        </Label>
        <div class="grow"/>
        <SvgIcon icon={mdiPencil} />
      </Button>
      <Button class="w-1/3 pr-2 pl-4 profile-delete" on:click={() => { deleteProfileName = $selectedProfile ?? ''; deleteProfileDialog = true; }} disabled={!$canModify || $profiles.length === 1}>
        <Label>
          Delete
        </Label>
        <div class="grow"/>
        <SvgIcon icon={mdiTrashCan} />
      </Button>
    </div>
    <div class="flex w-full mt-2">
      <Button class="w-1/2 pr-2 pl-5 mr-1 profile-import" on:click={() => importProfileDialog = true} disabled={!$canModify}>
        <Label>
          Import
        </Label>
        <div class="grow"/>
        <SvgIcon icon={mdiDownload} />
      </Button>
      <Button class="w-1/2 pr-2 pl-4 ml-1 profile-export" on:click={() => { exportCurrentProfile(); }} disabled={!$canModify}>
        <Label>
          Export
        </Label>
        <div class="grow"/>
        <SvgIcon icon={mdiUpload} />
      </Button>
    </div>
  </div>
  <div class="flex flex-col mt-4 h-md:mt-8">
    <span class="pl-4">Updates</span>
    <Updates />
  </div>
  <div class="flex flex-col mt-4 h-md:mt-8">
    <span class="pl-4">Other</span>
    <Settings />
    <Button variant="unelevated" class="w-full mt-2" on:click={() => BrowserOpenURL('https://docs.ficsit.app/satisfactory-modding/latest/ForUsers/SatisfactoryModManager.html')}>
      <Label>
        Help
      </Label>
      <div class="grow" />
      <SvgIcon icon={mdiHelpCircle} class="h-5 w-5" />
    </Button>
  </div>
  <div class="flex flex-col mt-4 h-md:mt-8">
    <span class="pl-4">Links</span>
    <Button variant="unelevated" class="w-full mt-2" on:click={() => BrowserOpenURL($siteURL)}>
      <Label>
        ficsit.app (Mod Repository)
      </Label>
      <div class="grow" />
      <SvgIcon icon={mdiWeb} class="h-5 w-5" />
    </Button>
    <Button variant="unelevated" class="w-full mt-2" on:click={() => BrowserOpenURL('https://discord.gg/xkVJ73E')}>
      <Label>
        Satisfactory Modding Discord
      </Label>
      <div class="grow" />
      <SvgIcon icon={siDiscord.path} class="h-5 w-5" />
    </Button>
    <Button variant="unelevated" class="w-full mt-2" on:click={() => BrowserOpenURL('https://github.com/satisfactorymodding/SatisfactoryModManager')} >
      <Label>
        SMM GitHub
      </Label>
      <div class="grow" />
      <SvgIcon icon={siGithub.path} class="h-5 w-5" />
    </Button>
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
    <Button on:click={() => addProfileDialog = false}>
      <Label>Cancel</Label>
    </Button>
    <Button on:click={finishAddProfile} disabled={!newProfileName}>
      <Label>Add</Label>
    </Button>
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
    <Button on:click={() => renameProfileDialog = false}>
      <Label>Cancel</Label>
    </Button>
    <Button on:click={finishRenameProfile} disabled={!renameNewProfileName}>
      <Label>Rename</Label>
    </Button>
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
    <Button on:click={() => deleteProfileDialog = false}>
      <Label>Cancel</Label>
    </Button>
    <Button on:click={finishDeleteProfile}>
      <Label>Delete</Label>
    </Button>
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
    <Button on:click={() => importProfileDialog = false}>
      <Label>Cancel</Label>
    </Button>
    <Button on:click={finishImportProfile} disabled={!importProfileName || !importProfileFilepath || !!importProfileError}>
      <Label>Import</Label>
    </Button>
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
