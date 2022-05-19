<script lang="ts">
  import MDIIcon from './MDIIcon.svelte';
  import { mdiWindowMinimize, mdiWindowRestore, mdiWindowMaximize, mdiWindowClose } from '@mdi/js';
  import { Quit, WindowMinimise, WindowToggleMaximise } from '../../../wailsjs/runtime';

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

<div class="titlebar">
  <div class="inline-flex items-center">
    <img src="/images/smm_icon_small.png" class="app-icon" alt="SMM Icon" />
  </div>
  <div class="bar">
    <div class="dragregion" data-wails-drag>
      <span class="app-title">Satisfactory Mod Manager</span>
    </div>
    <div class="button minimize grid place-items-center" on:click={minimize}>
      <MDIIcon icon={mdiWindowMinimize} />
    </div>
    <div class="button maximize grid place-items-center" on:click={toggleMaximize}>
      {#if isMaximized}
        <MDIIcon icon={mdiWindowRestore} />
      {:else}
        <MDIIcon icon={mdiWindowMaximize} />
      {/if}
    </div>
    <div class="button close grid place-items-center" on:click={close}>
      <MDIIcon icon={mdiWindowClose} />
    </div>
  </div>
</div>

<style scoped>
  .app-title {
    font-size: 15px !important;
    color: var(--v-text-lighten2);
  }
  .app-icon {
    margin: 4px 0px 4px 10px;
    height: 25px;
  }
  .titlebar {
    display: flex;
    user-select: none;
    z-index: 9999;
  }
  .titlebar,
  .titlebar > * {
    color: var(--v-text2-base) !important;
    background-color: var(--v-background-base);
  }
  .bar {
    flex-grow: 1;
    display: flex;
  }
  .dragregion {
    flex-grow: 1;
    text-align: left;
    vertical-align: middle;
    margin: 4px 3px 0px 10px;
    -webkit-app-region: drag;
  }
  .dragregion > span {
    flex-grow: 1;
    margin-top: -4px;
  }
  .button {
    width: 44px;
    text-align: center;
    font-size: 12pt;
    cursor: pointer;
    color: var(--v-text-base);
  }
  .button:hover {
    background-color: gray;
    color: white !important;
  }
  .close:hover {
    background-color: red;
    color: white !important;
  }
</style>
