<template>
  <v-app>
    <v-theme-provider root>
      <router-view class="content" />
    </v-theme-provider>
  </v-app>
</template>

<script>
export default {
  name: 'SatisfactoryModLauncherGUI',
  created() {
    this.validateSidePanel();
    this.$eventBus.$on('minimize-app', () => {
      this.minimizeApp();
    });
    this.$eventBus.$on('close-app', () => {
      this.closeApp();
    });
    this.$eventBus.$on('open-close-side-panel', (data) => {
      this.openCloseSidePanel(data);
    });
  },
  methods: {
    getBrowserWindow() {
      return this.$electron.remote.getCurrentWindow();
    },
    minimizeApp() {
      const browserWindow = this.getBrowserWindow();
      if (browserWindow.minimizable) {
        browserWindow.minimize();
      }
    },
    closeApp() {
      const browserWindow = this.getBrowserWindow();
      browserWindow.close();
    },
    validateSidePanel() {
      this.$store.state.isSidePanelOpen = this.getBrowserWindow().getSize()[0] > 500;
    },
    openCloseSidePanel(data) {
      this.validateSidePanel();
      const eventName = (this.$store.state.isSidePanelOpen && !data) ? 'close-side-panel' : 'open-side-panel';
      this.$electron.ipcRenderer.send(eventName);
      this.$store.state.isSidePanelOpen = !this.$store.state.isSidePanelOpen;
      this.$store.state.sidePanelData = data;
    },
  },
};
</script>
