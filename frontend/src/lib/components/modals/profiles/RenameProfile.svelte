<script lang="ts">
  import { newProfileName } from './renameProfile';

  import { RenameProfile } from '$lib/generated/wailsjs/go/ficsitcli/ficsitCLI';
  import { error } from '$lib/store/generalStore';
  import { profiles } from '$lib/store/ficsitCLIStore';

  export let parent: {onClose: () => void};

  export let profile: string;

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
    Rename profile
  </header>
  <section class="p-4 grow space-y-2">
    <label class="label w-full">
      <span>Old profile name</span>
      <input
        class="input px-4 py-2"
        placeholder="Old Profile"
        readonly
        type="text"
        value={profile}/>
    </label>
    <label class="label w-full">
      <span>New profile name</span>
      <input
        class="input px-4 py-2"
        class:input-error={newProfileNameExists}
        placeholder="My New Profile"
        type="text"
        bind:value={$newProfileName}/>
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
      disabled={!$newProfileName || newProfileNameExists}
      on:click={finishRenameProfile}>
      Rename
    </button>
  </footer>
</div>
