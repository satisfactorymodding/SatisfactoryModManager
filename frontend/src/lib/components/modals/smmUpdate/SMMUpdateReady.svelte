<script lang="ts">
  import { UpdateAndRestart } from '$lib/generated/wailsjs/go/bindings/Update';
  import { smmUpdate } from '$lib/store/smmUpdateStore';

  export let parent: {onClose: () => void};
</script>

<div class="card flex flex-col gap-2" style="width: 500px; max-width: calc(100vw - 32px);">
  <header class="card-header font-bold text-2xl text-center">
    SMM Update Ready - {$smmUpdate?.newVersion}
  </header>
  <section class="p-4 grow">
    <span class="text-base">Update ready to install</span>
  </section>
  <footer class="card-footer">
    <button
      class="btn"
      on:click={parent.onClose}>
      Cancel
    </button>
    <button
      class="btn text-primary-600"
      on:click={() => UpdateAndRestart()}>
      <span>Update and Restart</span>
    </button>
  </footer>
</div>
