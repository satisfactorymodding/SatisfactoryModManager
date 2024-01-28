<script lang="ts">
  import { profileFilepath, profileName } from './importProfile';

  import { OpenFileDialog } from '$lib/generated/wailsjs/go/bindings/App';
  import { ImportProfile, ReadExportedProfileMetadata } from '$lib/generated/wailsjs/go/ficsitcli/FicsitCLI';
  import type { ficsitcli } from '$lib/generated/wailsjs/go/models';
  import { error } from '$lib/store/generalStore';
  import { profiles } from '$lib/store/ficsitCLIStore';

  export let parent: {onClose: () => void};

  export let filepath = '';

  if (!$profileFilepath) {
    $profileFilepath = filepath;
  }

  $: newProfileNameExists = $profiles.includes($profileName);

  let fileDialogOpen = false;
  let importProfileMetadata: ficsitcli.ExportedProfileMetadata | null = null;
  let pickerError: string | null = null;
  async function pickImportProfileFile() {
    if(fileDialogOpen) {
      return;
    }
    fileDialogOpen = true;
    try {
      $profileFilepath = await OpenFileDialog({
        filters: [
          {
            displayName: 'SMM Profile (*.smmprofile)',
            pattern: '*.smmprofile',
          },
        ],
      });
      if (!$profileFilepath) {
        fileDialogOpen = false;
        return;
      }
      importProfileMetadata = await ReadExportedProfileMetadata($profileFilepath);
    } catch (e) {
      fileDialogOpen = false;
      if(e instanceof Error) {
        pickerError = e.message;
      } else if (typeof e === 'string') {
        pickerError = e;
      } else {
        pickerError = 'Unknown error';
      }
    }
    fileDialogOpen = false;
  }

  async function finishImportProfile() {
    try {
      await ImportProfile($profileName, $profileFilepath);
      $profileName = '';
      $profileFilepath = '';
      parent.onClose();
    } catch(e) {
      if(e instanceof Error) {
        $error = e.message;
      } else if (typeof e === 'string') {
        $error = e;
      } else {
        $error = 'Unknown error';
      }
    }
  }
</script>

<div style="width: 500px; max-width: calc(100vw - 32px);" class="card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    Import profile
  </header>
  <section class="p-4 grow space-y-2">
    <label class="label w-full">
      <span>Profile name</span>
      <input
        class="input px-4 py-2"
        placeholder="My New Profile"
        type="text"
        bind:value={$profileName}/>
    </label>
    <label class="label w-full">
      <span>Profile file</span>
      <input
        class="input px-4 py-2 hover:!cursor-pointer"
        class:input-error={!!pickerError || newProfileNameExists}
        readonly
        type="text" 
        value={$profileFilepath}
        on:click={() => pickImportProfileFile()}
      />
      {#if pickerError}
        <p>
          {pickerError}
        </p>
      {/if}
    </label>
  </section>
  <footer class="card-footer">
    <button
      class="btn"
      on:click={parent.onClose}>
      Cancel
    </button>
    <button
      class="btn text-primary-600"
      disabled={!$profileName || !$profileFilepath || !!pickerError || newProfileNameExists}
      on:click={finishImportProfile}>
      Import
    </button>
  </footer>
</div>
