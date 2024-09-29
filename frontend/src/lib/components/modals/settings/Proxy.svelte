<script lang="ts">
  import T from '$lib/components/T.svelte';
  import { progress } from '$lib/store/ficsitCLIStore';
  import { proxy } from '$lib/store/settingsStore';
  import { CloseAndRestart } from '$wailsjs/go/app/app';

  export let parent: { onClose: () => void };
  
  let proxyURL = $proxy;
  let saving = false;

  $: canChange = !$progress && !saving;
  $: canSave = proxyURL !== $proxy && canChange;

  async function setProxy() {
    saving = true;
    $proxy = proxyURL;
    setTimeout(() => {
      CloseAndRestart();
    }, 1000);
  }
</script>

<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="w-[60rem] card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    <T defaultValue="Set Proxy" keyName="settings.proxy.title" />
  </header>
  <section class="p-4 grow">
    <label class="label">
      <span><T defaultValue="Proxy" keyName="settings.proxy.proxy" /></span>
      <div class="flex items-baseline">
        <input
          class="input px-4 py-2 grow"
          type="text"
          bind:value={proxyURL}
        />
        <button
          class="btn shrink-0 text-primary-600"
          disabled={!canChange}
          on:click={() => { proxyURL = ''; setProxy(); }}>
          <span class="flex-auto text-start">
            <T defaultValue="Remove proxy" keyName="settings.proxy.remove" />
          </span>
        </button>
      </div>
    </label>
  </section>
  <footer class="card-footer">
    <button
      class="btn text-primary-600 variant-ringed"
      on:click={parent.onClose}>
      <span>
        <T defaultValue="Close" keyName="common.close" />
      </span>
    </button>
    <button
      class="btn shrink-0 text-primary-600"
      disabled={!canSave}
      on:click={() => setProxy()}>
      <span><T defaultValue="Save and restart" keyName="settings.proxy.save" /></span>
    </button>
  </footer>
</div>

