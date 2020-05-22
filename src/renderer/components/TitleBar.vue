<template>
  <div class="titlebar">
    <div
      class="d-inline-flex align-items-center"
      @click="settingsClicked"
    >
      <v-icon
        :color="getColorForState"
        class="ma-1 icon"
      >
        mdi-cog
      </v-icon>
    </div>
    <div class="bar">
      <div class="dragregion">
        <span class="app-title">{{ title }}</span>
      </div>
      <div
        class="button close"
        @click="onClose"
      >
        <span>&#10005;</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      menuOpen: false,
    };
  },
  computed: {
    state() {
      if (this.menuOpen) {
        return 'on';
      }
      if (this.$store.state.hasUpdate) {
        return 'notify';
      }
      return 'off';
    },
    getColorForState() {
      if (this.state === 'notify') {
        return '#ffc107';
      }
      if (this.state === 'on') {
        return 'primary';
      }
      return '#9e9e9e';
    },
  },
  methods: {
    onClose() {
      this.$electron.remote.getCurrentWindow().close();
    },
    settingsClicked() {
      this.menuOpen = !this.menuOpen;
    },
  },
};
</script>

<style scoped>
.app-title {
  font-size: 16px !important;
}
.titlebar {
  display: flex;
  height: var(--titlebar-height);
}
.titlebar, .titlebar > * {
  color: var(--v-text2-base) !important;
  background-color: var(--v-background-base);
}
.icon {
  font-size: 16px !important;
  padding: 3px 0px 0px 3px;
}
.bar {
  flex-grow: 1;
  display: flex;
}
.dragregion {
  flex-grow: 1;
  text-align: center;
  vertical-align: middle;
  margin: 3px 3px 0px 0px;
  -webkit-app-region: drag;
  display: flex;
  align-items: center;
}
.dragregion>span {
  flex-grow: 1;
  margin-top: -3px;
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
  margin-top: -3px;
  user-select: none;
}
.button:hover {
  background-color: var(--titlebar-button-hover-color);
}
.close:hover {
  background-color: red;
  color: white;
}
.button>span.dash {
  vertical-align: sub;
  margin-top: 0px;
}
#settingsIcon {
  height: 18px;
  margin: 4.5px;
}
</style>
