<script lang="ts">
  import Button, { Label } from '@smui/button';
  import Select, { Option } from '@smui/select';

  import { mdiCheckCircle, mdiCloseCircle } from '@mdi/js';
  import MdiIcon from '$lib/components/MDIIcon.svelte';
  
  import type { main } from '../../../wailsjs/go/models';
  import { GetInstallationsInfo, GetProfiles, SelectInstall, SetProfile } from '../../../wailsjs/go/main/FicsitCLI';

  let installs: main.InstallationInfo[] = [];
  let selectedInstall: main.InstallationInfo | null = null;

  let profiles: string[] = [];
  let selectedProfile: string | null = null;
  
  if(typeof window !== 'undefined') {
    Promise.all([
      GetInstallationsInfo().then((installations) => {
        installs = installations;
      }),
      GetProfiles().then((p) => {
        profiles = p;
      })
    ]).then(() => {
      selectedInstall = installs[0];
    });
  }

  $: {
    const path = selectedInstall?.info?.path;
    if(path) {
      SelectInstall(path);
      updateSelectedProfile();
    }
  }

  $: {
    if(selectedProfile) {
      SetProfile(selectedProfile);
      updateInstallProfile();
    }
  }

  function updateSelectedProfile() {
    if(selectedInstall && selectedInstall.installation) {
      selectedProfile = selectedInstall.installation.profile;
    }
  }

  function updateInstallProfile() {
    if(selectedInstall && selectedInstall.installation && selectedProfile) {
      selectedInstall.installation.profile = selectedProfile;
    }
  }

  let modsEnabled = true;
</script>

<div class="flex flex-col h-full p-4">
  <div class="flex flex-col">
    <span class="pl-4">Game version</span>
    <Select
      variant="filled"
      class="left-bar-select pt-2"
      bind:value={selectedInstall}
      ripple={false}
    >
      {#each installs as install}
        <Option value={install}>
          <Label>{install?.info?.branch} ({install?.info?.launcher}) - CL{install?.info?.version}</Label>
        </Option>
      {/each}
    </Select>
    <div class="flex w-full pt-2">
      <Button variant="unelevated" class={'w-1/2 rounded-tr-none rounded-br-none mods-toggle-button ' + (modsEnabled ? '' : 'mods-off')} on:click={() => modsEnabled = false}>
        <Label>
          Mods off
        </Label>
        <MdiIcon icon={mdiCloseCircle} />
      </Button>
      <Button variant="unelevated" class={'w-1/2 rounded-tl-none rounded-bl-none mods-toggle-button ' + (modsEnabled ? 'mods-on' : '')} on:click={() => modsEnabled = true}>
        <Label>
          Mods on
        </Label>
        <MdiIcon icon={mdiCheckCircle} />
      </Button>
    </div>
  </div>
  <div class="flex flex-col pt-8">
    <span class="pl-4">Profile</span>
    <Select
      variant="filled"
      class="left-bar-select pt-2"
      bind:value={selectedProfile}
      ripple={false}
    >
      {#each profiles as profile}
        <Option value={profile}>
          <Label>{profile}</Label>
        </Option>
      {/each}
    </Select>
  </div>
  <div class="grow"/>
  <center>
    <Button variant="unelevated" class="h-12 w-full">
      <Label>Launch Satisfactory</Label>
    </Button>
  </center>
</div>