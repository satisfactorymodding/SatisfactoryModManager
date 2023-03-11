<script lang="ts">
  import Button, { Label } from '@smui/button';
  import Select, { Option } from '@smui/select';
  import Dialog, { Title, Content, Actions } from '@smui/dialog';
  import TextField from '@smui/textfield'; 

  import { mdiCheckCircle, mdiCloseCircle, mdiDownload, mdiHelpCircle, mdiPencil, mdiPlusCircle, mdiTrashCan, mdiUpload, mdiWeb } from '@mdi/js';
  import { siDiscord, siGithub } from 'simple-icons/icons';
  import SvgIcon from '$lib/components/SVGIcon.svelte';
  
  import { addProfile, deleteProfile, importProfile, installs, profiles, canModify, renameProfile, selectedInstall, selectedProfile } from '$lib/store/ficsitCLIStore';
  import { error } from '$lib/store/generalStore';
  import { BrowserOpenURL, EventsOn } from '$wailsjs/runtime/runtime';
  import { OpenFileDialog } from '$wailsjs/go/bindings/App';

  import Settings from './Settings.svelte';
  import Updates from './Updates.svelte';
  import { bindings, ficsitcli_bindings } from '$wailsjs/go/models';
  import HelperText from '@smui/textfield/helper-text';
  import { ExportCurrentProfile, ReadExportedProfileMetadata, SetModsEnabled } from '$wailsjs/go/ficsitcli_bindings/FicsitCLI';
  import LaunchButton from './LaunchButton.svelte';

  $: modsEnabled = !$selectedInstall?.installation?.vanilla;

  async function setModsEnabled(enabled: boolean) {
    if ($selectedInstall) {
      await SetModsEnabled(enabled);
      if ($selectedInstall.installation) {
        $selectedInstall.installation.vanilla = !enabled;
      }
    }
  }

  let addProfileDialog = false;
  let newProfileName = '';
  async function finishAddProfile() {
    try {
      await addProfile(newProfileName);
      selectedProfile.set(newProfileName);
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
      await renameProfile(renameOldProfileName, renameNewProfileName);
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
      await deleteProfile(deleteProfileName);
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
  let importProfileMetadata: ficsitcli_bindings.ExportedProfileMetadata | null = null;
  async function pickImportProfileFile() {
    if(fileDialogOpen) {
      return;
    }
    fileDialogOpen = true;
    try {
      importProfileFilepath = await OpenFileDialog(bindings.OpenDialogOptions.createFrom({
        filters: [
          bindings.FileFilter.createFrom({
            displayName: 'SMM Profile (*.smmprofile)',
            pattern: '*.smmprofile',
          }),
        ],
      }));
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
      await importProfile(importProfileName, importProfileFilepath);
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
</script>

<div class="flex flex-col h-full p-4 left-bar w-[22rem] min-w-[22rem] ">
  <div class="flex flex-col">
    <span class="pl-4">Game version</span>
    <Select
      variant="filled"
      class="mt-2"
      menu$class="max-h-[32rem]"
      bind:value={$selectedInstall}
      ripple={false}
      disabled={!$canModify}
    >
      {#each $installs as install}
        <Option value={install}>
          <Label>{install?.info?.branch} ({install?.info?.launcher})</Label>
        </Option>
      {/each}
    </Select>
    <div class="flex w-full mt-2">
      <Button variant="unelevated" class="w-1/2 rounded-tr-none rounded-br-none mods-toggle-button {modsEnabled ? '' : 'mods-off'}" on:click={() => setModsEnabled(false)} disabled={!$canModify}>
        <Label>
          Mods off
        </Label>
        <div class="grow"/>
        <SvgIcon icon={mdiCloseCircle} class="h-5 w-5" />
      </Button>
      <Button variant="unelevated" class="w-1/2 rounded-tl-none rounded-bl-none mods-toggle-button {modsEnabled ? 'mods-on' : ''}" on:click={() => setModsEnabled(true)} disabled={!$canModify}>
        <Label>
          Mods on
        </Label>
        <div class="grow"/>
        <SvgIcon icon={mdiCheckCircle} class="h-5 w-5" />
      </Button>
    </div>
  </div>
  <div class="flex flex-col mt-8">
    <span class="pl-4">Profile</span>
    <Select
      variant="filled"
      class="mt-2"
      menu$class="max-h-[32rem]"
      bind:value={$selectedProfile}
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
      <Button class="w-1/3 pr-2 pl-4 profile-delete" on:click={() => { deleteProfileName = $selectedProfile ?? ''; deleteProfileDialog = true; }} disabled={!$canModify}>
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
  <div class="flex flex-col mt-8">
    <span class="pl-4">Updates</span>
    <Updates />
  </div>
  <div class="flex flex-col mt-8">
    <span class="pl-4">Other</span>
    <Settings />
    <Button variant="unelevated" class="w-full mt-2">
      <Label>
        Help
      </Label>
      <div class="grow" />
      <SvgIcon icon={mdiHelpCircle} class="h-5 w-5" />
    </Button>
  </div>
  <div class="flex flex-col mt-8">
    <span class="pl-4">Links</span>
    <Button variant="unelevated" class="w-full mt-2" on:click={() => BrowserOpenURL('https://ficsit.app')}>
      <Label>
        ficsit.app mod website
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