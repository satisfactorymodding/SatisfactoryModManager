<script lang="ts">
  import { GenerateDebugInfo } from '$lib/generated/wailsjs/go/app/app';

  export let parent: {onClose: () => void};

  export let error: string;
</script>

<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="w-[48rem] card flex flex-col gap-6">
  <header class="card-header font-bold text-2xl text-center">
    Something went wrong
  </header>
  <section class="px-4 overflow-y-auto">
    <p>{error}</p>
  </section>
  <section class="px-4">
    <p>Seems wrong? Click the button below and send the generated zip file on the <a class="text-primary-600 underline" href="https://discord.gg/xkVJ73E">modding discord</a> in #help-using-mods.</p>
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
