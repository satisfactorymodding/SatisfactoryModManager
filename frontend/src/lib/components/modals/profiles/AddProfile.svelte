<script lang="ts">
  import { getTranslate } from '@tolgee/svelte';

  import { newProfileName } from './addProfile';

  import T from '$lib/components/T.svelte';
  import { getModalStore } from '$lib/skeletonExtensions';
  import { selectedProfile } from '$lib/store/ficsitCLIStore';
  import { error } from '$lib/store/generalStore';
  import { AddProfile } from '$wailsjs/go/ficsitcli/ficsitCLI';

  export let parent: { onClose: () => void };

  const { t } = getTranslate();

  const modalStore = getModalStore();

  async function finishAddProfile() {
    try {
      await AddProfile($newProfileName);
      await selectedProfile.asyncSet($newProfileName);

      $newProfileName = '';

      modalStore.close('addProfile');
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
    <T defaultValue="Add profile" keyName="profiles.add.title" />
  </header>
  <section class="p-4 grow">
    <label class="label w-full">
      <span><T defaultValue="Profile name" keyName="profiles.add.profile-name" /></span>
      <input
        class="input px-4 py-2"
        placeholder={$t('profiles.add.profile-name-placeholder', 'New Profile Name')}
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
      disabled={!newProfileName}
      on:click={finishAddProfile}>
      <T defaultValue="Add" keyName="common.add" />
    </button>
  </footer>
</div>
