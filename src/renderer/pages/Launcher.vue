<template>
  <v-card class="h-100 d-flex">
    <v-card class="h-100 minw-500 maxw-500">
      <Header />
      <v-card
        tile
        class="main-launcher launcher-wrapper d-flex flex-column"
      >
        <ModConfig
          v-show="showConfig"
        />
        <div class="container d-flex">
          <SearchBar class="flex-grow-1" />
          <div class="hide-config-button align-self-center p-1">
            <v-btn
              text
              class="fs-075 p-0"
              style="min-width: 35px;"
              @click="showConfig = !showConfig"
            >
              <v-icon>{{ showConfig ? '$visibleIcon' : '$visibleOffIcon' }}</v-icon>
            </v-btn>
          </div>
        </div>
        <div class="mod-list">
          <ModList />
        </div>
        <div class="launcher-button w-100">
          <v-btn
            tile
            block
            bottom
            height="70"
            :disabled="inProgress.length > 0 || configLoadInProgress"
            @click="launchSatisfactory"
          >
            Launch
          </v-btn>
        </div>
        <v-card
          v-show="inProgress.length > 0 || configLoadInProgress"
          class="position-absolute w-100"
        >
          <v-progress-linear
            color="success"
            indeterminate
            height="2"
          ></v-progress-linear>
        </v-card>
        <v-card
          v-if="isLoading"
          class="loading fixed-top w-100 h-100 d-flex align-items-center justify-content-center"
        >
          <transition name="loading">
            <div>
              Loading<span class="ldot">.</span><span class="ldot">.</span><span class="ldot">.</span>
              <div
                class="header-action icon header-no-drag position-fixed"
                style="width: 32px; top: 5px; right: 0; color: white"
                @click="closeApp()"
              >
                <v-icon>
                  $closeIcon
                </v-icon>
              </div>
            </div>
          </transition>
        </v-card>
        <Alert />
      </v-card>
    </v-card>
    <v-card
      tile
      flat
      class="color-2 h-100 w-100 overflow-auto launcher-scroll-2"
    >
      <ModDetail v-if="sidePanelData" />
    </v-card>
  </v-card>
</template>

<script>
import { mapState, mapActions } from 'vuex';
import {
  toggleDebug,
  isDebug,
  clearCache,
} from 'satisfactory-mod-manager-api';
import Header from '../components/Header';
import ModConfig from '../components/ModConfig';
import ModDetail from '../components/ModDetail';
import SearchBar from '../components/SearchBar';
import ModList from '../components/ModList';
import Alert from '../components/Alert';

export default {
  name: 'Launcher',
  components: {
    Header,
    ModConfig,
    ModDetail,
    SearchBar,
    ModList,
    Alert,
  },
  data() {
    return {
      showConfig: true,
      isLoading: true,
    };
  },
  computed: {
    ...mapState([
      'configLoadInProgress',
      'inProgress',
      'sidePanelData',
    ]),
  },
  mounted() {
    this.$electron.ipcRenderer.on('openedByUrl', (e, url) => {
      this.$store.dispatch('openedByUrl', url);
    });
    this.$electron.ipcRenderer.on('toggleDebug', () => {
      toggleDebug();
      if (isDebug()) {
        this.$electron.ipcRenderer.send('openDevTools');
        this.$electron.ipcRenderer.send('open-side-panel');
      }
    });
    this.$electron.ipcRenderer.on('clearCache', () => {
      clearCache();
      if (this.$store.state.selectedSatisfactoryInstall) {
        this.$store.state.selectedSatisfactoryInstall.clearCache();
      }
    });
    this.$electron.ipcRenderer.on('update-available', (e, updateInfo) => {
      this.$store.state.availableUpdate = updateInfo;
      this.$eventBus.$emit(
        'popup-alert',
        'info',
        `New update available: ${this.$store.state.availableUpdate ? this.$store.state.availableUpdate.version : ''}`,
        'Update will be installed when the app closes',
        this.$store.state.availableUpdate ? this.$store.state.availableUpdate.releaseNotes : '',
      );
    });
  },
  created() {
    this.initLauncher().then(() => {
      setTimeout(() => {
        this.isLoading = false;
      }, 2000);
    });
  },
  methods: {
    ...mapActions([
      'initLauncher',
      'launchSatisfactory',
    ]),
    closeApp() {
      this.$eventBus.$emit('close-app');
    },
  },
};
</script>
