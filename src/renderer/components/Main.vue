<template>
  <v-card
    class="d-flex"
    height="100%"
  >
    <v-card
      height="100%"
      style="width: 500px; min-width: 500px; max-width: 500px; z-index: 1"
    >
      <TitleBar
        title="Satisfactory Mod Manager"
        style="user-select: none;"
      />
      <ControlArea style="user-select: none;" />
      <ModsList />
      <v-btn
        block
        tile
        color="primary"
        dark
        elevation="0"
        height="82px"
        style="font-size: 18px;"
        :disabled="!!inProgress.length || isGameRunning"
        @click="launchSatisfactory"
      >
        <b>LAUNCH GAME</b>
      </v-btn>
    </v-card>
    <v-card
      tile
      flat
      class="color-2overflow-auto"
      width="100%"
      height="100%"
    >
      <ModDetails v-if="expandedModId" />
    </v-card>
    <v-dialog
      :value="!!error"
      max-width="290"
      persistent
    >
      <v-card>
        <v-card-title class="headline">
          Error
        </v-card-title>

        <v-card-text>
          {{ error }}
        </v-card-text>

        <v-card-actions>
          <v-btn
            color="primary"
            text
            @click="clearError"
          >
            Ok
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog
      hide-overlay
      persistent
      :value="isLoadingAppInProgress"
      width="500"
      height="230"
    >
      <v-card
        color="#2f3136 !important"
      >
        <v-row
          no-gutters
          justify="center"
        >
          <v-img
            class="mt-4"
            src="/static/smm_icon.png"
            max-height="82px"
            max-width="87px"
          />
        </v-row>
        <v-card-title class="loading-text-main">
          SATISFACTORY MOD MANAGER IS LOADING
        </v-card-title>

        <v-card-text
          v-if="isLoadingAppInProgress"
          class="text-center"
        >
          <v-progress-linear
            :value="Math.round(currentLoadingAppProgress.progress * 100)"
            :class="currentLoadingAppProgress.fastUpdate ? 'fast' : ''"
            background-color="#000000"
            color="#5bb71d"
            height="2"
            reactive
            :indeterminate="currentLoadingAppProgress.progress < 0"
          />
          {{ currentLoadingAppProgress.message || '&nbsp;' }}
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script>
import { mapState } from 'vuex';
import { exec } from 'child_process';
import TitleBar from './TitleBar';
import ControlArea from './ControlArea';
import ModsList from './ModsList';
import ModDetails from './ModDetails';
import { lastElement } from '../utils';

export default {
  components: {
    TitleBar,
    ControlArea,
    ModsList,
    ModDetails,
  },
  computed: {
    ...mapState(
      [
        'selectedInstall',
        'expandedModId',
        'inProgress',
        'isGameRunning',
        'error',
      ],
    ),
    isLoadingAppInProgress() {
      return this.inProgress.some((prog) => prog.id === '__loadingApp__');
    },
    loadingAppProgress() {
      return this.inProgress.find((prog) => prog.id === '__loadingApp__');
    },
    currentLoadingAppProgress() {
      return lastElement(this.loadingAppProgress.progresses);
    },
  },
  async mounted() {
    this.$electron.ipcRenderer.send('unexpand');
    await this.$store.dispatch('initApp');
    this.$electron.ipcRenderer.send('vue-ready');
  },
  methods: {
    clearError() {
      this.$store.dispatch('clearError');
    },
    launchSatisfactory() {
      if (this.selectedInstall && !this.isGameRunning) {
        exec(`start "" "${this.selectedInstall.launchPath}"`).unref();
        this.$store.commit('launchGame');
      }
    },
    lastElement,
  },
};
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css?family=Open+Sans');
.loading-text-main {
  color: #e6e6e4;
  text-align: center;
  font-family: 'Open Sans';
  font-size: 16px !important;
  font-weight: 600;
  display: block !important;
}
</style>
