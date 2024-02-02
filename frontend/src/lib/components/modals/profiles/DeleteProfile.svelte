<script lang="ts">
  import { DeleteProfile } from '$lib/generated/wailsjs/go/ficsitcli/ficsitCLI';
  import { error } from '$lib/store/generalStore';

  export let parent: {onClose: () => void};

  export let profile: string;

  async function finishDeleteProfile() {
    try {
      await DeleteProfile(profile);
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
    Delete profile
  </header>
  <section class="p-4 grow space-y-2">
    <label class="label w-full">
      <span>Profile name</span>
      <input
        class="input px-4 py-2"
        readonly
        type="text"
        value={profile}/>
    </label>
  </section>
  <footer class="card-footer">
    <button
      class="btn"
      on:click={parent.onClose}>
      Cancel
    </button>
    <button
      class="btn text-error-500"
      on:click={finishDeleteProfile}>
      Delete
    </button>
  </footer>
</div>
