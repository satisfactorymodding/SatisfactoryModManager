<script lang="ts">
  import { getTranslate } from '@tolgee/svelte';

  import { profileFilepath, profileName } from './importProfile';

  import T from '$lib/components/T.svelte';
  import { profiles, selectedInstallMetadata } from '$lib/store/ficsitCLIStore';
  import { error } from '$lib/store/generalStore';
  import { OpenFileDialog } from '$wailsjs/go/app/app';
  import { ImportProfile, ReadExportedProfileMetadata } from '$wailsjs/go/ficsitcli/ficsitCLI';
  import type { ficsitcli } from '$wailsjs/go/models';

  export let parent: { onClose: () => void };

  const { t } = getTranslate();

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

<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="w-[40rem] card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    <T defaultValue="Import profile" keyName="profiles.import.title" />
  </header>
  <section class="p-4 grow space-y-2">
    <label class="label w-full">
      <span><T defaultValue="Profile name" keyName="profiles.import.profile-name" /></span>
      <input
        class="input px-4 py-2"
        placeholder={$t('profiles.import.profile-name-placeholder', 'New Profile Name')}
        type="text"
        bind:value={$profileName}/>
    </label>
    <label class="label w-full">
      <span><T defaultValue="Profile file" keyName="profiles.import.profile-file" /></span>
      <input
        class="input px-4 py-2 hover:!cursor-pointer"
        class:input-error={!!pickerError || newProfileNameExists}
        readonly
        type="text" 
        value={$profileFilepath}
        on:click={() => pickImportProfileFile()}
      />
      {#if importProfileMetadata}
        {#if importProfileMetadata.gameVersion < ($selectedInstallMetadata?.info?.version ?? 0)}
          <p>
            <T defaultValue="This profile was created with a newer version of the game. It may not be compatible with this version." keyName="profiles.import.profile-version-warning" />
          </p>
        {/if}
      {/if}
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
      <T defaultValue="Cancel" keyName="common.cancel" />
    </button>
    <button
      class="btn text-primary-600"
      disabled={!$profileName || !$profileFilepath || !!pickerError || newProfileNameExists}
      on:click={finishImportProfile}>
      <T defaultValue="Import" keyName="common.import" />
    </button>
  </footer>
</div>
