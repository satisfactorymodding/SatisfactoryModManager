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
    if (this.getBrowserWindow().getSize()[0] > 500) {
      this.$store.state.isSidePanelOpen = true;
    }
    this.$eventBus.$on('minimize-app', () => {
      this.minimizeApp();
    });
    this.$eventBus.$on('close-app', () => {
      this.closeApp();
    });
    this.$eventBus.$on('open-side-panel', (data) => {
      this.openSidePanel(data);
    });
    this.$eventBus.$on('close-side-panel', () => {
      this.closeSidePanel();
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
    openSidePanel(data) {
      if (!this.$store.state.isSidePanelOpen) {
        this.$electron.ipcRenderer.send('open-side-panel');
      }
      this.$store.state.isSidePanelOpen = true;
      this.$store.state.sidePanelData = data;
    },
    closeSidePanel() {
      if (this.$store.state.isSidePanelOpen) {
        this.$electron.ipcRenderer.send('close-side-panel');
      }
      this.$store.state.isSidePanelOpen = false;
      this.$store.state.sidePanelData = null;
    },
  },
};
</script>
