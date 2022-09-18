<script lang="ts">
  import MDIIcon from './MDIIcon.svelte';
  import { mdiWindowMinimize, mdiWindowRestore, mdiWindowMaximize, mdiWindowClose } from '@mdi/js';
  import { Quit, WindowMinimise, WindowToggleMaximise } from '$wailsjs/runtime';
  import { GetVersion } from '$wailsjs/go/bindings/App';

  function minimize() {
    WindowMinimise();
  }

  function toggleMaximize() {
    WindowToggleMaximise();
  }

  function close() {
    Quit();
  }

  let version = '0.0.0';
  GetVersion().then((v) => {
    version = v;
  });

  let isMaximized = false;
</script>

<div class="dragregion flex items-center h-9">
  <img src="/images/smm_icon_small.png" class="h-7 pl-4 pr-2" alt="SMM Icon" />
  <div>
    <span class="app-title pt-3">Satisfactory Mod Manager v{version}</span>
  </div>

  <div class="grow" />

  <div class="button w-11 cursor-default minimize grid justify-center items-center h-full" on:click={minimize}>
    <MDIIcon icon={mdiWindowMinimize}/>
  </div>
  <div class="button w-11 cursor-default maximize grid justify-center items-center h-full" on:click={toggleMaximize}>
    {#if isMaximized}
      <MDIIcon icon={mdiWindowRestore}/>
    {:else}
      <MDIIcon icon={mdiWindowMaximize}/>
    {/if}
  </div>
  <div class="button w-11 cursor-default close grid justify-center items-center h-full" on:click={close}>
    <MDIIcon icon={mdiWindowClose}/>
  </div>
</div>

<style scoped>
  .dragregion {
    --webkit-app-region: drag;
    --wails-draggable: drag;
  }
  .button:hover {
    background-color: gray;
  }
  .close:hover {
    background-color: red;
  }
</style>
