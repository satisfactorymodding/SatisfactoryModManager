<template>
  <v-card
    class="d-flex flex-column"
    height="100%"
  >
    <TitleBar
      v-if="!hasFrame"
      title="Satisfactory Mod Manager"
    />
    <v-card
      class="d-flex"
      height="100%"
    >
      <v-card
        height="100%"
        style="width: 550px; min-width: 550px; max-width: 550px; z-index: 1;"
        class="d-flex flex-column"
      >
        <MenuBar />
        <AnnouncementsBar />
        <ControlArea
          style="user-select: none;"
          :filters.sync="filters"
          :available-filters="availableFilters"
          :available-sorting="availableSorting"
        />
        <ModsList
          class="flex-grow-1 flex-shrink-1"
          style="height: 0px"
          :filters.sync="filters"
          @set-available-filters="availableFilters = $event"
          @set-available-sorting="availableSorting = $event"
        />
        <v-btn
          class="flex-grow-0 flex-shrink-0"
          block
          tile
          color="primary"
          elevation="0"
          style="font-size: 18px; min-height: 50px;"
          :style="launchButton && selectedInstall && selectedInstall.launchPath && !isGameRunning ? 'height: 98px' : ''"
          :disabled="!!inProgress.length || isGameRunning || (selectedInstall && !selectedInstall.launchPath)"
          :ripple="!launchButton"
          @click="() => !launchButton && launchSatisfactory()"
        >
          <template v-if="launchButton && selectedInstall && selectedInstall.launchPath && !isGameRunning">
            <img
              src="static/launch_fun.png"
              draggable="false"
            >
            <img
              :src="`static/launch_fun_button_${launchFunState}.png`"
              style="position: absolute;"
              :style="launchFunState === 'press' ? 'top: 1px' : ''"
              draggable="false"
              @click="launchFunPress"
              @mousedown="launchFunState = 'press'"
              @mouseup="launchFunState = 'over'"
              @mouseenter="launchFunState = 'over'"
              @mouseleave="launchFunState = 'normal'"
            >
          </template>
          <span
            v-else
          >{{ isGameRunning ? 'Game is running' : (selectedInstall && selectedInstall.launchPath ? 'Launch Satisfactory' : 'Cannot launch this install') }}</span>
        </v-btn>
      </v-card>
      <ModDetails v-if="expandedModId" />
    </v-card>
    <v-dialog
      v-model="errorDialog"
      :persistent="errorPersistent"
      max-width="550"
    >
      <v-card>
        <v-card-title class="headline">
          Error
        </v-card-title>

        <v-card-text style="white-space: pre-line;">
          <span>{{ error }}</span>
          <br>
          <span>Seems wrong? <a @click="exportDebugData">Generate debug info</a> and send it together with this error message
            over on <a @click="moddingDiscord">the modding discord</a> in #help-using-mods</span>
        </v-card-text>

        <v-card-actions v-if="!errorPersistent">
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
      v-model="installSetupErrorDialog"
      persistent
      max-width="550"
    >
      <v-card>
        <v-card-title class="headline">
          Error
        </v-card-title>

        <v-card-text style="white-space: pre-line;">
          <span>{{ installSetupError }}</span>
          <br>
          <span>Seems wrong? <a @click="exportDebugData">Generate debug info</a> and send it together with this error message
            over on <a @click="moddingDiscord">the modding discord</a> in #help-using-mods</span>
        </v-card-text>

        <v-card-actions>
          <v-btn
            color="primary"
            text
            @click="retryInstallSetup"
          >
            Retry
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog
      persistent
      :value="isLoadingAppInProgress"
      width="500"
      height="230"
    >
      <v-card color="loadingBackground !important">
        <v-row
          no-gutters
          justify="center"
        >
          <v-img
            class="mt-4"
            src="static/smm_icon.png"
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
            :class="currentLoadingAppProgress.fast ? 'fast' : ''"
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
    <v-dialog
      persistent
      :value="showUpdateDownloadProgress"
      width="500"
      height="230"
    >
      <v-card color="loadingBackground !important">
        <v-row
          no-gutters
          justify="center"
        >
          <v-img
            class="mt-4"
            src="static/smm_icon.png"
            max-height="82px"
            max-width="87px"
          />
        </v-row>
        <v-card-title class="loading-text-main">
          UPDATING SATISFACTORY MOD MANAGER
        </v-card-title>

        <v-card-text
          v-if="isUpdateDownloadInProgress"
          class="text-center"
        >
          <v-progress-linear
            :value="Math.round(currentUpdateDownloadProgress.progress * 100)"
            :class="currentUpdateDownloadProgress.fast ? 'fast' : ''"
            background-color="#000000"
            color="#5bb71d"
            height="2"
            reactive
            :indeterminate="currentUpdateDownloadProgress.progress < 0"
          />
          {{ currentUpdateDownloadProgress.message || '&nbsp;' }}
        </v-card-text>
      </v-card>
    </v-dialog>
    <v-dialog
      v-model="oldSMLauncherInstalled"
      width="500"
      height="230"
    >
      <v-card color="loadingBackground !important">
        <v-card-title class="loading-text-main">
          Old SMLauncher install
        </v-card-title>

        <v-card-text
          class="text-center"
        >
          The update from SMLauncher to SMM is installed as a new app because of the name change. Uninstall the old SMLauncher version
        </v-card-text>

        <v-card-actions>
          <v-btn
            color="primary"
            text
            @click="uninstallOldSMLauncher"
          >
            Uninstall old SMLauncher
          </v-btn>
          <v-btn
            color="text"
            text
            @click="oldSMLauncherInstalled = false"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script>
import { mapState } from 'vuex';
import { exec } from 'child_process';
import { getCacheFolder } from 'platform-folders';
import fs from 'fs';
import path from 'path';
import { lastElement, bytesToAppropriate } from '@/utils';
import { getSetting } from '~/settings';
import TitleBar from './TitleBar';
import MenuBar from './menu-bar/MenuBar';
import AnnouncementsBar from './AnnouncementsBar';
import ControlArea from './ControlArea';
import ModsList from './mods-list/ModsList';
import ModDetails from './mod-details/ModDetails';

const SMLauncherUninstallerPath = path.join(getCacheFolder(), 'Programs', 'satisfactory-mod-launcher-gui', 'Uninstall Satisfactory Mod Launcher.exe');

export default {
  components: {
    TitleBar,
    MenuBar,
    AnnouncementsBar,
    ControlArea,
    ModsList,
    ModDetails,
  },
  data() {
    return {
      smmUpdateDownloadProgress: {},
      updateDownloadFinished: false,
      showUpdateDownloadProgress: false,
      oldSMLauncherInstalled: false,
      filters: {
        modFilters: {},
        sortBy: '',
        search: '',
      },
      availableFilters: [],
      availableSorting: [],
      launchFunState: 'normal',
      launchFun: 0,
    };
  },
  computed: {
    ...mapState(
      [
        'selectedInstall',
        'expandedModId',
        'inProgress',
        'isGameRunning',
        'installSetupError',
        'error',
        'errorPersistent',
        'launchButton',
      ],
    ),
    errorDialog: {
      get() {
        return !!this.error;
      },
      set() {
        this.clearError();
      },
    },
    installSetupErrorDialog() {
      return !!this.installSetupError;
    },
    isLoadingAppInProgress() {
      return this.inProgress.some((prog) => prog.id === '__loadingApp__');
    },
    loadingAppProgress() {
      return this.inProgress.find((prog) => prog.id === '__loadingApp__');
    },
    currentLoadingAppProgress() {
      return lastElement(this.loadingAppProgress.progresses);
    },
    isUpdateDownloadInProgress() {
      return this.inProgress.some((prog) => prog.id === '__downloadingUpdate__');
    },
    updateDownloadProgress() {
      return this.inProgress.find((prog) => prog.id === '__downloadingUpdate__');
    },
    currentUpdateDownloadProgress() {
      return lastElement(this.updateDownloadProgress.progresses);
    },
    hasFrame() {
      return this.$electron.remote.getGlobal('frame');
    },
  },
  async mounted() {
    const keyQueue = [];
    const code = [38, 38, 40, 40, 37, 39, 37, 39, 66, 65];
    window.addEventListener('keydown', (event) => {
      keyQueue.push(event.keyCode);
      if (keyQueue.length > code.length) {
        keyQueue.shift();
      }
      if (keyQueue.length === code.length && keyQueue.every((val, idx) => code[idx] === val)) {
        this.$store.dispatch('konami');
      }
    });
    this.$electron.ipcRenderer.send('unexpand');
    this.$electron.ipcRenderer.on('installedMods', () => {
      this.$electron.ipcRenderer.send('installedMods', this.$store.state.installedMods);
    });
    this.$electron.ipcRenderer.on('updateAvailable', () => {
      this.smmUpdateDownloadProgress = {
        id: '__downloadingUpdate__',
        progresses: [{
          id: '', progress: -1, message: 'Downloading update', fast: false,
        }],
      };
      this.updateDownloadFinished = false;
      this.$electron.ipcRenderer.on('updateDownloadProgress', this.updateProgress);
      this.$electron.ipcRenderer.once('updateDownloaded', () => {
        this.inProgress.remove(this.smmUpdateDownloadProgress);
        this.updateDownloadFinished = true;
        this.$electron.ipcRenderer.off('updateDownloadProgress', this.updateProgress);
      });
    });
    this.$electron.ipcRenderer.on('openedByUrl', (_, url) => {
      const parsed = new URL(url);
      const command = parsed.pathname.replace(/^\/+|\/+$/g, '');
      if (command === 'install') {
        const modID = parsed.searchParams.get('modID');
        const version = parsed.searchParams.get('version');
        this.$store.dispatch('installModVersion', { modId: modID, version });
      }
    });
    this.$electron.ipcRenderer.on('autoUpdateError', (_, err) => {
      this.$store.dispatch('showError', `Error while checking for SMM updates: ${err}`);
    });
    const hasUpdate = await this.checkForUpdates();
    if (hasUpdate && getSetting('updateCheckMode', 'launch') === 'launch') {
      this.downloadUpdate();
      return;
    }
    this.$root.$emit('doneLaunchUpdateCheck');
    this.$root.$on('downloadUpdate', this.downloadUpdate);
    await this.$store.dispatch('initApp');

    this.$electron.ipcRenderer.send('vue-ready');
    this.oldSMLauncherInstalled = fs.existsSync(SMLauncherUninstallerPath);
  },
  methods: {
    async checkForUpdates() {
      this.$electron.ipcRenderer.send('checkForUpdates');
      return new Promise((resolve) => {
        this.$electron.ipcRenderer.once('updateAvailable', () => {
          resolve(true);
        });
        this.$electron.ipcRenderer.once('updateNotAvailable', () => {
          resolve(false);
        });
      });
    },
    downloadUpdate() {
      if (!this.updateDownloadFinished) {
        this.showUpdateDownloadProgress = true;
        this.inProgress.push(this.smmUpdateDownloadProgress);
        this.$electron.ipcRenderer.on('updateDownloaded', () => {
          setInterval(() => {
            if (this.inProgress.length === 0) {
              this.$electron.remote.getCurrentWindow().close();
            }
          }, 100);
        });
      } else {
        setInterval(() => {
          if (this.inProgress.length === 0) {
            this.$electron.remote.getCurrentWindow().close();
          }
        }, 100);
      }
    },
    updateProgress(e, info) {
      this.smmUpdateDownloadProgress.progresses[0].progress = info.percent / 100;
      this.smmUpdateDownloadProgress.progresses[0].message = `Downloading update ${Math.round(info.percent)}% (${bytesToAppropriate(info.transferred)}/${bytesToAppropriate(info.total)} - ${bytesToAppropriate(info.bytesPerSecond)}/s)`;
      this.smmUpdateDownloadProgress.progresses[0].fast = true;
    },
    clearError() {
      this.$store.dispatch('clearError');
    },
    retryInstallSetup() {
      this.$store.dispatch('clearInstallSetupError');
      this.$store.dispatch('setupInstalls');
    },
    async launchSatisfactory() {
      if (this.selectedInstall && !this.isGameRunning) {
        this.$store.commit('launchGame');
        exec(this.selectedInstall.launchPath).unref();
      }
    },
    uninstallOldSMLauncher() {
      if (fs.existsSync(SMLauncherUninstallerPath)) {
        exec(`start "" "${SMLauncherUninstallerPath}"`).unref();
      }
      this.oldSMLauncherInstalled = false;
    },
    exportDebugData() {
      this.$root.$emit('exportDebugData');
    },
    moddingDiscord() {
      this.$root.$emit('moddingDiscord');
    },
    launchFunPress() {
      this.launchFun += 1;
      if (this.launchFun === 15) {
        this.launchSatisfactory();
        this.launchFun = 0;
      }
    },
    lastElement,
  },
};
</script>

<style scoped>
</style>
