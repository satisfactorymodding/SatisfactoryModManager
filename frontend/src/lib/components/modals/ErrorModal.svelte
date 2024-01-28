<script lang="ts">
  import { GenerateDebugInfo } from '$lib/generated/wailsjs/go/bindings/DebugInfo';

  export let parent: {onClose: () => void};

  export let error: string;
</script>

<div class="card flex flex-col gap-2" style="width: 500px; max-width: calc(100vw - 32px);">
  <header class="card-header font-bold text-2xl text-center">
    Something went wrong
  </header>
  <section class="p-4">
    <p>{error}</p>
    <p class="pt-4">Seems wrong? Click the button below and send the generated zip file on the <a class="text-primary-600 underline" href="https://discord.gg/xkVJ73E">modding discord</a> in #help-using-mods.</p>
  </section>
  <footer class="card-footer">
    <button
      class="btn text-primary-600 variant-ringed"
      on:click={GenerateDebugInfo}>
      Generate debug info
    </button>
    <button
      class="btn"
      on:click={parent.onClose}>
      Close
    </button>
  </footer>
</div>
