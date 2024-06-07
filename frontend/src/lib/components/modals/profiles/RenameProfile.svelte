<script lang="ts">
  import { getTranslate } from '@tolgee/svelte';

  import { newProfileName } from './renameProfile';

  import T from '$lib/components/T.svelte';
  import { profiles } from '$lib/store/ficsitCLIStore';
  import { error } from '$lib/store/generalStore';
  import { RenameProfile } from '$wailsjs/go/ficsitcli/ficsitCLI';

  export let parent: { onClose: () => void };

  export let profile: string;

  const { t } = getTranslate();

  $: newProfileNameExists = $profiles.includes($newProfileName);

  async function finishRenameProfile() {
    try {
      await RenameProfile(profile, $newProfileName);
      $newProfileName = '';
      parent.onClose();
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
</script>

<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="w-[40rem] card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    <T defaultValue="Rename profile" keyName="profiles.rename.title" />
  </header>
  <section class="p-4 grow space-y-2">
    <label class="label w-full">
      <span><T defaultValue="Old profile name" keyName="profiles.rename.old-profile" /></span>
      <input
        class="input px-4 py-2"
        placeholder={$t('profiles.rename.old-profile-placeholder', 'Old Profile')}
        readonly
        type="text"
        value={profile}/>
    </label>
    <label class="label w-full">
      <span><T defaultValue="New profile name" keyName="profiles.rename.new-profile" /></span>
      <input
        class="input px-4 py-2"
        class:input-error={newProfileNameExists}
        placeholder={$t('profiles.rename.new-profile-placeholder', 'New Profile')}
        type="text"
        bind:value={$newProfileName}/>
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
      disabled={!$newProfileName || newProfileNameExists}
      on:click={finishRenameProfile}>
      <T defaultValue="Rename" keyName="common.rename" />
    </button>
  </footer>
</div>
