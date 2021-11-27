<template>
  <div class="titlebar">
    <div
      class="d-inline-flex align-items-center"
    >
      <img
        src="static/smm_icon_small.png"
        class="app-icon"
      >
    </div>
    <div class="bar">
      <div class="dragregion">
        <span class="app-title">{{ title }}</span>
      </div>
      <div
        class="button minimize"
        @click="minimize"
      >
        <v-icon
          color="text"
          style="margin-left: calc(50% - 9.5px);"
        >
          mdi-window-minimize
        </v-icon>
      </div>
      <div
        class="button maximize"
        @click="maximize"
      >
        <v-icon
          v-if="isMaximized"
          color="text"
          style="margin-left: calc(50% - 9.5px);"
        >
          mdi-window-restore
        </v-icon>
        <v-icon
          v-else
          color="text"
          style="margin-left: calc(50% - 9.5px);"
        >
          mdi-window-maximize
        </v-icon>
      </div>
      <div
        class="button close"
        @click="close"
      >
        <v-icon
          color="text"
          style="margin-left: calc(50% - 9.5px);"
        >
          mdi-window-close
        </v-icon>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex';

export default {
  props: {
    title: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      isMaximized: false,
    };
  },
  computed: {
    ...mapState([
      'inProgress',
    ]),
  },
  created() {
    this.$electron.ipcRenderer.on('maximize', this.onMaximize);
    this.$electron.ipcRenderer.on('unmaximize', this.onUnmaximize);
  },
  async mounted() {
    this.isMaximized = await this.$electron.ipcRenderer.invoke('isMaximized');
  },
  destroyed() {
    this.$electron.ipcRenderer.off('maximize', this.onMaximize);
    this.$electron.ipcRenderer.off('unmaximize', this.onUnmaximize);
  },
  methods: {
    minimize() {
      this.$electron.ipcRenderer.invoke('minimize');
    },
    maximize() {
      if (!this.isMaximized) {
        this.$electron.ipcRenderer.invoke('maximize');
      } else {
        this.$electron.ipcRenderer.invoke('unmaximize');
      }
    },
    close() {
      if (this.inProgress.length === 0) {
        this.$electron.ipcRenderer.invoke('close');
      }
    },
    onMaximize() {
      this.isMaximized = true;
    },
    onUnmaximize() {
      this.isMaximized = false;
    },
  },
};
</script>

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
.titlebar, .titlebar > * {
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
  display: flex;
  align-items: center;
}
.dragregion>span {
  flex-grow: 1;
  margin-top: -4px;
}
.button {
  width: 44px;
  text-align: center;
  font-size: 12pt;
  display: flex;
  align-items: center;
  cursor: pointer;
  color: var(--v-text-base);
}
.button>span {
  flex-grow: 1;
  user-select: none;
}
.button:hover {
  background-color: gray;
  color: white !important;
}
.close:hover {
  background-color: red;
  color: white !important;
}
.button>span.dash {
  vertical-align: sub;
  margin-top: 0px;
}
</style>
