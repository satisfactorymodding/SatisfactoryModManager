<script lang="ts">
  import T, { translationElementPart } from '$lib/components/T.svelte';
  import { GenerateDebugInfo } from '$wailsjs/go/app/app';

  export let parent: { onClose: () => void };

  export let error: string;
</script>

<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="w-[48rem] card flex flex-col gap-6">
  <header class="card-header font-bold text-2xl text-center">
    <T defaultValue="Something went wrong" keyName="error.title" />
  </header>
  <section class="px-4 overflow-y-auto">
    <p>{error}</p>
  </section>
  <section class="px-4">
    <T
      defaultValue="Seems wrong? Click the button below and send the generated zip file on the <1>modding discord</1> in #help-using-mods."
      keyName="error.help"
      parts={[
        translationElementPart('a', {
          href: 'https://discord.gg/xkVJ73E',
          class: 'text-primary-600 underline',
        }),
      ]}
    />
  </section>
  <footer class="card-footer">
    <button
      class="btn text-primary-600 variant-ringed"
      on:click={GenerateDebugInfo}>
      <T defaultValue="Generate debug info" keyName="error.generate_debug_info" />
    </button>
    <button
      class="btn"
      on:click={parent.onClose}>
      <T defaultValue="Close" keyName="common.close" />
    </button>
  </footer>
</div>
