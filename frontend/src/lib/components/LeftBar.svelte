<script lang="ts">
  import Button, { Label } from '@smui/button';
  import Select, { Option } from '@smui/select';
  import Dialog, { Title, Content, Actions } from '@smui/dialog';
  import TextField from '@smui/textfield'; 

  import { mdiBug, mdiCheckCircle, mdiChevronRight, mdiClipboard, mdiCloseCircle, mdiDiscord, mdiDownload, mdiGithub, mdiHelpCircle, mdiPencil, mdiPlusCircle, mdiSync, mdiTrashCan, mdiTune, mdiWeb } from '@mdi/js';
  import MdiIcon from '$lib/components/MDIIcon.svelte';
  
  import { addProfile, checkForUpdates, deleteProfile, installs, profiles, progress, renameProfile, selectedInstall, selectedProfile, updateCheckInProgress, updates } from '$lib/ficsitCLIStore';
  import { UpdateAllMods } from '$wailsjs/go/bindings/FicsitCLI';
  import { BrowserOpenURL } from '$wailsjs/runtime/runtime';
  import Menu, { type MenuComponentDev } from '@smui/menu';
  import List, { Item, PrimaryText, Text } from '@smui/list';
  import { GenerateDebugInfo } from '$wailsjs/go/bindings/DebugInfo';

  import { manifestMods, lockfileMods } from '$lib/ficsitCLIStore';
  import { GetModNameDocument } from '$lib/generated';
  import { getClient } from '@urql/svelte';

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

  async function updateAll() {
    if($updates.length > 0) {
      const err = await UpdateAllMods();
      if(!err) {
        $updates = [];
      }
    }
  }
  
  const urqlClient = getClient();

  async function copyModList() {
    // Generate mod entries
    const modList = await Promise.all(Object.keys($manifestMods).map(async (modReference) => {
      let modName = modReference;
      if(modReference === 'SML') {
        modName = 'Satisfactory Mod Loader';
      } else {
        const result = await urqlClient.query(GetModNameDocument, { modReference }).toPromise();
        if(result?.data?.getModByReference?.name) {
          modName = result.data.getModByReference.name;
        }
      }
      return {
        friendlyName: modName,
        modReference,
        version: $lockfileMods[modReference].version,
      };
    }));
      // Sort by Friendly Name
    modList.sort((a, b) => {
      const x = a.friendlyName.toLowerCase();
      const y = b.friendlyName.toLowerCase();
      return x.localeCompare(y);
    });
    // Get max lengths to use for padding
    const maxFriendlyNameLen = Math.max(...modList.map((mod) => mod.friendlyName.length));
    const maxModReferenceLen = Math.max(...modList.map((mod) => mod.modReference.length));
    // Create header and add all mods to string
    let modListString = `${'Mod Name'.padEnd(maxFriendlyNameLen + 1) + 'Mod Reference'.padEnd(maxModReferenceLen + 1)}Version\n`;
    modList.forEach((mod) => {
      mod.friendlyName = mod.friendlyName.padEnd(maxFriendlyNameLen, ' ');
      mod.modReference = mod.modReference.padEnd(maxModReferenceLen, ' ');
      modListString += `${mod.friendlyName} ${mod.modReference} ${mod.version}\n`;
    });
    navigator.clipboard.writeText(modListString.trim());
  }

  let settingsMenu: MenuComponentDev;
  let debugMenu: MenuComponentDev;
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
    <Button variant="unelevated" class="w-full mt-2 update-button {$updates.length > 0 ? 'has-update' : ''}" on:click={updateAll} disabled={!canInstall}>
      <Label>
        {#if $updates.length === 0}
          No updates right now
        {:else}
          {$updates.length} updates available
        {/if}
      </Label>
      <div class="grow" />
      <MdiIcon icon={mdiCheckCircle} class="h-5" />
    </Button>
    <Button variant="unelevated" class="w-full mt-2" on:click={checkForUpdates} disabled={!canInstall}>
      <Label>
        Check for updates
      </Label>
      <div class="grow" />
      <MdiIcon icon={mdiSync} class="h-5 {$updateCheckInProgress ? 'update-check' : ''}" />
    </Button>
  </div>
  <div class="flex flex-col mt-8">
    <span class="pl-4">Other</span>
    <div>
      <Button variant="unelevated" class="w-full mt-2" on:click={() => settingsMenu.setOpen(true)}>
        <Label>
          SMM settings
        </Label>
        <div class="grow" />
        <MdiIcon icon={mdiTune} class="h-5" />
      </Button>
      <Menu bind:this={settingsMenu} class="w-full max-h-[32rem] overflow-visible" anchorCorner="TOP_RIGHT">
        <List>
          <div>
            <Item on:click={() => debugMenu.setOpen(true)}>
              <MdiIcon icon={mdiBug} class="h-5" />
              <!-- <div class="w-7"/> -->
              <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
                <PrimaryText class="text-base">Debug</PrimaryText>
              </Text>
              <div class="grow" />
              <MdiIcon icon={mdiChevronRight} class="h-5" />
            </Item>
            <Menu bind:this={debugMenu} class="w-full max-h-[32rem] overflow-visible" anchorCorner="TOP_RIGHT">
              <List>
                <Item on:click={() => GenerateDebugInfo()}>
                  <MdiIcon icon={mdiDownload} class="h-5" />
                  <!-- <div class="w-7"/> -->
                  <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
                    <PrimaryText class="text-base">Generate debug info</PrimaryText>
                  </Text>
                  <div class="grow" />
                </Item>
                <Item on:click={() => copyModList()}>
                  <MdiIcon icon={mdiClipboard} class="h-5" />
                  <!-- <div class="w-7"/> -->
                  <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
                    <PrimaryText class="text-base">Copy mods list</PrimaryText>
                  </Text>
                  <div class="grow" />
                </Item>
              </List>
            </Menu>
          </div>
        </List>
      </Menu>
    </div>
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