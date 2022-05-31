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
  <Select
    variant="filled"
    class="select-rounded"
    bind:value={selectedInstall}
    ripple={false}
  >
    {#each installs as install}
      <Option value={install}>
        <Label>{install?.info?.branch} ({install?.info?.launcher}) - CL{install?.info?.version}</Label>
      </Option>
    {/each}
  </Select>
  <div class="flex w-full pt-4">
    <Button variant="unelevated" class={'w-1/2 rounded-tr-none rounded-br-none ' + (!modsEnabled ? 'mods-toggle-button-off' : '')} color={!modsEnabled ? undefined : 'secondary'} on:click={() => modsEnabled = false}>
      <Label>
        Mods off
        <MdiIcon icon={mdiCloseCircle} />
      </Label>
    </Button>
    <Button variant="unelevated" class="w-1/2 rounded-tl-none rounded-bl-none" color={modsEnabled ? 'primary' : 'secondary'} on:click={() => modsEnabled = true}>
      <Label>
        Mods on
        <MdiIcon icon={mdiCheckCircle} />
      </Label>
    </Button>
  </div>
  <Select
    variant="filled"
    class="select-rounded pt-4"
    bind:value={selectedProfile}
    ripple={false}
  >
    {#each profiles as profile}
      <Option value={profile}>
        <Label>{profile}</Label>
      </Option>
    {/each}
  </Select>
  <div class="grow"/>
  <center>
    <Button variant="unelevated" class="h-12 w-full">
      <Label>Launch Satisfactory</Label>
    </Button>
  </center>
</div>

<style>
  * :global(.mods-toggle-button-off) {
    background-color: #e51c22;
  }

  * :global(.select-rounded),
  * :global(.select-rounded .mdc-select__anchor) {
    border-radius: 8px;
  }
</style>