<script lang="ts">
  import Button, { Label } from '@smui/button';
  import Select, { Option } from '@smui/select';
  import Dialog, { Title, Content, Actions } from '@smui/dialog';
  import TextField from '@smui/textfield'; 

  import { mdiCheckCircle, mdiCloseCircle, mdiDiscord, mdiGithub, mdiHelpCircle, mdiPencil, mdiPlusCircle, mdiTrashCan, mdiWeb } from '@mdi/js';
  import MdiIcon from '$lib/components/MDIIcon.svelte';
  
  import { addProfile, deleteProfile, installs, profiles, progress, renameProfile, selectedInstall, selectedProfile } from '$lib/ficsitCLIStore';
  import { BrowserOpenURL } from '$wailsjs/runtime/runtime';

  import Settings from './Settings.svelte';
  import Updates from './Updates.svelte';

  $: canInstall = !$progress;

  let modsEnabled = true;

  let addProfileDialog = false;
  let newProfileName = '';
  async function finishAddProfile() {
    await addProfile(newProfileName);
    selectedProfile.set(newProfileName);
    newProfileName = '';
  }

  let renameProfileDialog = false;
  let renameOldProfileName = '';
  let renameNewProfileName = '';
  async function finishRenameProfile() {
    await renameProfile(renameOldProfileName, renameNewProfileName);
    renameOldProfileName = '';
    renameNewProfileName = '';
  }

  let deleteProfileDialog = false;
  let deleteProfileName = '';
  async function finishDeleteProfile() {
    await deleteProfile(deleteProfileName);
    deleteProfileName = '';
  }

</script>

<div class="flex flex-col h-full p-4 left-bar w-[24rem] min-w-[24rem] ">
  <div class="flex flex-col">
    <span class="pl-4">Game version</span>
    <Select
      variant="filled"
      class="mt-2"
      menu$class="max-h-[32rem]"
      bind:value={$selectedInstall}
      ripple={false}
      disabled={!canInstall}
    >
      {#each $installs as install}
        <Option value={install}>
          <Label>{install?.info?.branch} ({install?.info?.launcher}) - CL{install?.info?.version}</Label>
        </Option>
      {/each}
    </Select>
    <div class="flex w-full mt-2">
      <Button variant="unelevated" class="w-1/2 rounded-tr-none rounded-br-none mods-toggle-button {modsEnabled ? '' : 'mods-off'}" on:click={() => modsEnabled = false} disabled={!canInstall}>
        <Label>
          Mods off
        </Label>
        <div class="grow"/>
        <MdiIcon icon={mdiCloseCircle} class="h-5" />
      </Button>
      <Button variant="unelevated" class="w-1/2 rounded-tl-none rounded-bl-none mods-toggle-button {modsEnabled ? 'mods-on' : ''}" on:click={() => modsEnabled = true} disabled={!canInstall}>
        <Label>
          Mods on
        </Label>
        <div class="grow"/>
        <MdiIcon icon={mdiCheckCircle} class="h-5" />
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
      disabled={!canInstall}
    >
      {#each $profiles as profile}
        <Option value={profile}>
          <Label>{profile}</Label>
        </Option>
      {/each}
    </Select>
    <div class="flex w-full mt-2">
      <Button class="w-1/3 pr-2 pl-5 profile-add" on:click={() => addProfileDialog = true} disabled={!canInstall}>
        <Label>
          Add
        </Label>
        <div class="grow"/>
        <MdiIcon icon={mdiPlusCircle} />
      </Button>
      <Button class="w-1/3 mx-2 pr-0 profile-edit" on:click={() => { renameOldProfileName = $selectedProfile ?? ''; renameProfileDialog = true; }} disabled={!canInstall}>
        <Label>
          Rename
        </Label>
        <div class="grow"/>
        <MdiIcon icon={mdiPencil} />
      </Button>
      <Button class="w-1/3 pr-2 pl-4 profile-delete" on:click={() => { deleteProfileName = $selectedProfile ?? ''; deleteProfileDialog = true; }} disabled={!canInstall}>
        <Label>
          Delete
        </Label>
        <div class="grow"/>
        <MdiIcon icon={mdiTrashCan} />
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
      <MdiIcon icon={mdiHelpCircle} class="h-5" />
    </Button>
  </div>
  <div class="flex flex-col mt-8">
    <span class="pl-4">Links</span>
    <Button variant="unelevated" class="w-full mt-2" on:click={() => BrowserOpenURL('https://ficsit.app')}>
      <Label>
        ficsit.app mod website
      </Label>
      <div class="grow" />
      <MdiIcon icon={mdiWeb} class="h-5" />
    </Button>
    <Button variant="unelevated" class="w-full mt-2" on:click={() => BrowserOpenURL('https://discord.gg/xkVJ73E')}>
      <Label>
        Satisfactory Modding Discord
      </Label>
      <div class="grow" />
      <MdiIcon icon={mdiDiscord} class="h-5" />
    </Button>
    <Button variant="unelevated" class="w-full mt-2" on:click={() => BrowserOpenURL('https://github.com/satisfactorymodding/SatisfactoryModManager')} >
      <Label>
        SMM Source
      </Label>
      <div class="grow" />
      <MdiIcon icon={mdiGithub} class="h-5" />
    </Button>
  </div>
  <div class="grow"/>
  <center>
    <Button variant="unelevated" class="h-12 w-full launch-game" disabled={!canInstall}>
      <Label>Play Satisfactory</Label>
      <div class="grow" />
    </Button>
  </center>
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
    <Button on:click={finishAddProfile}>
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
    <Button on:click={finishRenameProfile}>
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