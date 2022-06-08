<script lang="ts">
  import Button, { Label } from '@smui/button';
  import Select, { Option } from '@smui/select';
  import Dialog, { Title, Content, Actions } from '@smui/dialog';
  import TextField from '@smui/textfield'; 

  import { mdiCheckCircle, mdiCloseCircle, mdiPencil, mdiPlusCircle, mdiTrashCan } from '@mdi/js';
  import MdiIcon from '$lib/components/MDIIcon.svelte';
  
  import { addProfile, deleteProfile, installs, profiles, renameProfile, selectedInstall, selectedProfile } from '$lib/ficsitCLIStore';

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

<div class="flex flex-col h-full p-4">
  <div class="flex flex-col">
    <span class="pl-4">Game version</span>
    <Select
      variant="filled"
      class="left-bar-select pt-2"
      menu$class="max-h-[32rem]"
      bind:value={$selectedInstall}
      ripple={false}
    >
      {#each $installs as install}
        <Option value={install}>
          <Label>{install?.info?.branch} ({install?.info?.launcher}) - CL{install?.info?.version}</Label>
        </Option>
      {/each}
    </Select>
    <div class="flex w-full pt-2">
      <Button variant="unelevated" class="w-1/2 rounded-tr-none rounded-br-none mods-toggle-button {modsEnabled ? '' : 'mods-off'}" on:click={() => modsEnabled = false}>
        <Label>
          Mods off
        </Label>
        <div class="grow"/>
        <MdiIcon icon={mdiCloseCircle} class="h-5" />
      </Button>
      <Button variant="unelevated" class="w-1/2 rounded-tl-none rounded-bl-none mods-toggle-button {modsEnabled ? 'mods-on' : ''}" on:click={() => modsEnabled = true}>
        <Label>
          Mods on
        </Label>
        <div class="grow"/>
        <MdiIcon icon={mdiCheckCircle} class="h-5" />
      </Button>
    </div>
  </div>
  <div class="flex flex-col pt-8">
    <span class="pl-4">Profile</span>
    <Select
      variant="filled"
      class="left-bar-select pt-2"
      menu$class="max-h-[32rem]"
      bind:value={$selectedProfile}
      ripple={false}
    >
      {#each $profiles as profile}
        <Option value={profile}>
          <Label>{profile}</Label>
        </Option>
      {/each}
    </Select>
    <div class="flex w-full pt-2">
      <Button class="w-1/3 pr-0 profile-add" on:click={() => addProfileDialog = true}>
        <Label>
          Add
        </Label>
        <div class="grow"/>
        <MdiIcon icon={mdiPlusCircle} />
      </Button>
      <Button class="w-1/3 mx-2 pr-0 profile-edit" on:click={() => { renameOldProfileName = $selectedProfile ?? ''; renameProfileDialog = true; }}>
        <Label>
          Rename
        </Label>
        <div class="grow"/>
        <MdiIcon icon={mdiPencil} />
      </Button>
      <Button class="w-1/3 pr-0 profile-delete" on:click={() => { deleteProfileName = $selectedProfile ?? ''; deleteProfileDialog = true; }}>
        <Label>
          Delete
        </Label>
        <div class="grow"/>
        <MdiIcon icon={mdiTrashCan} />
      </Button>
    </div>
  </div>
  <div class="grow"/>
  <center>
    <Button variant="unelevated" class="h-12 w-full">
      <Label>Launch Satisfactory</Label>
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