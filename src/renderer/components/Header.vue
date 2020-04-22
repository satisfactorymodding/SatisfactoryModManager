<template>
  <v-system-bar
    window
    fixed
    class="header-drag disable-selection minw-500 maxw-500 color-1"
  >
    <div
      class="header-action icon header-no-drag"
      style="width: 32px"
      @click="showMenu()"
    >
      <v-icon color="outine-color">
        $settingIcon
      </v-icon>
    </div>
    <span class="header-action header-no-drag">
      <router-link
        exact-active-class="active"
        to="/"
      >
        Satisfactory Mod Launcher
      </router-link>
    </span>
    <v-spacer />
    <div
      class="header-action icon header-no-drag"
      style="width: 32px"
      @click="minimizeApp()"
    >
      <v-icon>
        $minimizeIcon
      </v-icon>
    </div>
    <div
      class="header-action icon header-no-drag"
      style="width: 32px"
      @click="closeApp()"
    >
      <v-icon>
        $closeIcon
      </v-icon>
    </div>
  </v-system-bar>
</template>

<script>
import EventBus from '../plugins/eventBus';

export default {
  name: 'Header',
  data: () => ({
    isDarkMode: false,
  }),
  methods: {
    showMenu() {
      this.$electron.ipcRenderer.send('display-app-menu');
    },
    minimizeApp() {
      EventBus.$emit('minimize-app');
    },
    closeApp() {
      EventBus.$emit('close-app');
    },
  },
};
</script>
