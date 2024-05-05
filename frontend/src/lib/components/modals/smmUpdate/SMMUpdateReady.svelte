<script lang="ts">
  import { isUpdateOnStart } from './smmUpdate';

  import { UpdateAndRestart } from '$lib/generated/wailsjs/go/autoupdate/autoUpdate';
  import { smmUpdate } from '$lib/store/smmUpdateStore';

  export let parent: { onClose: () => void };
</script>

<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="w-[48rem] card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    SMM Update Ready - {$smmUpdate?.version}
  </header>
  <section class="p-4 flex-auto">
    <span class="text-base">Update ready to install</span>
  </section>
  <footer class="card-footer">
    {#if !$isUpdateOnStart}
      <button
        class="btn"
        on:click={parent.onClose}>
        Cancel
      </button>
    {/if}
    <button
      class="btn text-primary-600"
      on:click={() => UpdateAndRestart()}>
      <span>Update and Restart</span>
    </button>
  </footer>
</div>
