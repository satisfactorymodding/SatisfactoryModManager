<script lang="ts">
  import SvgIcon from './SVGIcon.svelte';
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

<div class="flex items-center h-9">
  <div class="dragregion grow flex items-center">
    <img src="/images/smm_icon_small.png" class="h-7 pl-4 pr-2" alt="SMM Icon" />
    <div>
      <span class="app-title pt-3">Satisfactory Mod Manager v{version}</span>
    </div>
  </div>

  <div class="hover:bg-gray-500 p-1.5 w-11 cursor-default grid justify-center items-center h-full" on:click={minimize}>
    <SvgIcon icon={mdiWindowMinimize}/>
  </div>
  <div class="hover:bg-gray-500 p-1.5 w-11 cursor-default grid justify-center items-center h-full" on:click={toggleMaximize}>
    {#if isMaximized}
      <SvgIcon icon={mdiWindowRestore}/>
    {:else}
      <SvgIcon icon={mdiWindowMaximize}/>
    {/if}
  </div>
  <div class="hover:bg-red-600 p-1.5 w-11 cursor-default grid justify-center items-center h-full" on:click={close}>
    <SvgIcon icon={mdiWindowClose}/>
  </div>
</div>

<style scoped>
  .dragregion {
    --webkit-app-region: drag;
    --wails-draggable: drag;
  }
</style>
