<script lang="ts">
  import { mdiWindowClose, mdiWindowMaximize, mdiWindowMinimize, mdiWindowRestore } from '@mdi/js';

  import SvgIcon from './SVGIcon.svelte';

  import { version } from '$lib/store/settingsStore';
  import { Quit, WindowMinimise, WindowToggleMaximise } from '$wailsjs/runtime';

  function minimize() {
    WindowMinimise();
  }

  function toggleMaximize() {
    WindowToggleMaximise();
  }

  function close() {
    Quit();
  }

  let isMaximized = false;
</script>

<div class="flex items-center h-9">
  <div class="dragregion grow flex items-center">
    <img class="h-7 pl-4 pr-2" alt="SMM Icon" src="/images/smm_icon_small.png" />
    <div>
      <span class="app-title pt-3 text-base">Satisfactory Mod Manager v{$version}</span>
    </div>
  </div>

  <!-- System level keybinds can be used instead -->
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div class="hover:bg-gray-500 p-1.5 w-11 cursor-default grid justify-center items-center h-full" on:click={minimize}>
    <SvgIcon class="w-full h-full" icon={mdiWindowMinimize}/>
  </div>
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div class="hover:bg-gray-500 p-1.5 w-11 cursor-default grid justify-center items-center h-full" on:click={toggleMaximize}>
    <SvgIcon class="w-full h-full" icon={isMaximized ? mdiWindowRestore : mdiWindowMaximize}/>
  </div>
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div class="hover:bg-red-600 p-1.5 w-11 cursor-default grid justify-center items-center h-full" on:click={close}>
    <SvgIcon class="w-full h-full" icon={mdiWindowClose}/>
  </div>
</div>

<style scoped>
  .dragregion {
    --webkit-app-region: drag;
    --wails-draggable: drag;
  }
</style>
