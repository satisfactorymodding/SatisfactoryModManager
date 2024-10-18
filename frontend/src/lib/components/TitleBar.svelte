<script lang="ts">
  import { mdiWindowClose, mdiWindowMaximize, mdiWindowMinimize, mdiWindowRestore } from '@mdi/js';
  import { onDestroy } from 'svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { version } from '$lib/store/settingsStore';
  import { Quit, WindowIsMaximised, WindowMinimise, WindowToggleMaximise } from '$wailsjs/runtime';

  let isMaximized = false;

  function minimize() {
    WindowMinimise();
  }

  function toggleMaximize() {
    WindowToggleMaximise();
    isMaximized = !isMaximized;
  }

  function close() {
    Quit();
  }

  function updateMaximized() {
    WindowIsMaximised().then((maximized) => {
      isMaximized = maximized;
    });
  }

  const updateMaximizedInterval = setInterval(updateMaximized, 100);

  onDestroy(() => {
    clearInterval(updateMaximizedInterval);
  });
</script>

<div class="flex items-center h-9">
  <!-- System level keybinds can be used instead -->
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div class="dragregion grow flex items-center" on:click={toggleMaximize}>
    <img class="h-7 pl-4 pr-2" alt="SMM Icon" src="/images/smm_icon_small.png" />
    <div>
      <span class="app-title pt-3 text-base">Satisfactory Mod Manager v{$version}</span>
    </div>
  </div>

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
