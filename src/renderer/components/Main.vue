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
        <template v-if="(!launchButton && !launchCat) || !(selectedInstall && selectedInstall.launchPath && !isGameRunning)">
          <v-tooltip
            top
            color="background"
            :disabled="incompatibleInstalledCount === 0 && possiblyCompatibleInstalledCount === 0"
          >
            <template #activator="{ on, attrs }">
              <v-btn
                class="flex-grow-0 flex-shrink-0"
                block
                tile
                :color="buttonColor"
                elevation="0"
                style="font-size: 18px; min-height: 50px;"
                :style="{ 'height: 98px' : launchButton && selectedInstall && selectedInstall.launchPath && !isGameRunning }"
                :disabled="!!inProgress.length || isGameRunning || (selectedInstall && !selectedInstall.launchPath)"
                :ripple="!launchButton"
                v-bind="attrs"
                @click="() => !launchButton && launchSatisfactory()"
                v-on="on"
              >
                {{ launchButtonText }}
              </v-btn>
            </template>
            <span
              v-if="incompatibleInstalledCount !== 0 && possiblyCompatibleInstalledCount !== 0"
            >
              You have {{ possiblyCompatibleInstalledCount }} mod{{ possiblyCompatibleInstalledCount > 1 ? 's' : '' }}
              that {{ possiblyCompatibleInstalledCount > 1 ? 'are' : 'is' }} likely incompatible and
              {{ incompatibleInstalledCount }} mod{{ incompatibleInstalledCount > 1 ? 's' : '' }} that are incompatible with your game.
              These will either not load or crash your game.
              Are you sure you want to launch?
            </span>
            <span v-else-if="incompatibleInstalledCount !== 0">
              You have {{ incompatibleInstalledCount }} incompatible mod{{ incompatibleInstalledCount > 1 ? 's' : '' }} which will either not load or crash your game.
              Are you sure you want to launch?
            </span>
            <span v-else>
              You have {{ possiblyCompatibleInstalledCount }} mod{{ possiblyCompatibleInstalledCount > 1 ? 's' : '' }}
              that {{ possiblyCompatibleInstalledCount > 1 ? 'are' : 'is' }} likely incompatible and can crash your game.
              Are you sure you want to launch?
            </span>
          </v-tooltip>
        </template>
        <template v-else-if="launchButton">
          <div style="height: 98px">
            <img
              src="static/launch/fun/launch_fun.png"
              draggable="false"
            >
            <img
              :src="`static/launch/fun/launch_fun_button_${launchFunState}.png`"
              style="position: relative; left: 227px"
              :style="launchFunState === 'press' ? 'top: -105px' : 'top: -106px'"
              draggable="false"
              @click="launchFunPress"
              @mousedown="launchFunState = 'press'"
              @mouseup="launchFunState = 'over'"
              @mouseenter="launchFunState = 'over'"
              @mouseleave="launchFunState = 'normal'"
            >
          </div>
        </template>
        <template v-else-if="launchCat">
          <div
            style="height: 60px"
            @mouseup="catPressed = false"
            @mousemove="catMouseMove"
          >
            <img
              src="static/launch/cat/bg.png"
              draggable="false"
            >
            <img
              :src="`static/launch/cat/cat_full.png`"
              style="position: relative; top: -61px"
              :style="`left: calc(-450px + ${catPosition * 82}%)`"
              draggable="false"
              @mousedown="catMouseDown"
            >
          </div>
        </template>
      </v-card>
      <ModDetails v-if="expandedModId" />
    </v-card>
    <v-dialog
      v-model="uriInstallDialog"
      max-width="550"
    >
      <v-card>
        <div class="d-flex flex-no-wrap justify-space-between">
          <div>
            <v-card-title
              class="text-h5"
            >
              {{ uriInstallModData ? uriInstallModData.name : '(Loading name...)' }}
            </v-card-title>

            <v-card-subtitle>
              {{ uriInstallModData ? uriInstallModData.short_description : '(Loading description...)' }}
            </v-card-subtitle>
            <v-card-subtitle>
              Version: {{ uriInstallModVersion || 'latest' }}
            </v-card-subtitle>

            <v-card-actions>
              <v-btn
                color="primary"
                text
                @click="uriInstallMod"
              >
                Install
              </v-btn>
              <v-btn
                text
                @click="uriInstallDialog = false"
              >
                Cancel
              </v-btn>
            </v-card-actions>
          </div>

          <v-avatar
            class="ma-3"
            size="125"
            tile
          >
            <v-img :src="(uriInstallModData ? uriInstallModData.logo : null) || 'https://ficsit.app/images/no_image.webp'" />
          </v-avatar>
        </div>
      </v-card>
    </v-dialog>
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
            :class="{ 'fast': currentLoadingAppProgress.fast }"
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
            :class="{ 'fast': currentUpdateDownloadProgress.fast }"
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
    <OutdatedModsDialog />
  </v-card>
</template>

<script>
import { mapState } from 'vuex';
import { exec } from 'child_process';
import { getCacheFolder } from 'platform-folders';
import fs from 'fs';
import path from 'path';
import gql from 'graphql-tag';
import {
  lastElement, bytesToAppropriate, isCompatibleFast, COMPATIBILITY_LEVEL,
} from '@/utils';
import { getSetting } from '~/settings';
import TitleBar from './TitleBar';
import MenuBar from './menu-bar/MenuBar';
import AnnouncementsBar from './announcements-bar/AnnouncementsBar';
import ControlArea from './ControlArea';
import ModsList from './mods-list/ModsList';
import ModDetails from './mod-details/ModDetails';
import OutdatedModsDialog from './OutdatedModsDialog';

const SMLauncherUninstallerPath = path.join(getCacheFolder(), 'Programs', 'satisfactory-mod-launcher-gui', 'Uninstall Satisfactory Mod Launcher.exe');

export default {
  components: {
    TitleBar,
    MenuBar,
    AnnouncementsBar,
    ControlArea,
    ModsList,
    ModDetails,
    OutdatedModsDialog,
  },
  data() {
    return {
      uriInstallDialog: false,
      uriInstallModReference: null,
      uriInstallModVersion: null,
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
      catPosition: 0,
      catOffset: 0,
      catPressed: false,
    };
  },
  asyncComputed: {
    hasFrame: {
      async get() {
        return this.$electron.ipcRenderer.invoke('hasFrame');
      },
      default: true,
    },
    modStates: {
      async get() {
        return Promise.all(Object.keys(this.installedMods).map(async (modReference) => {
          if (modReference === 'SML' || modReference === 'bootstrapper') {
            return { modReference, name: modReference, compatible: true };
          }
          const { mod } = (await this.$apollo.query({
            query: gql`            
              query checkOutdatedMod($modReference: ModReference!) {
                mod: getModByReference(modReference: $modReference) {
                  id,
                  mod_reference,
                  name,
                  versions(filter: { limit: 100 }) {
                    id,
                    sml_version,
                  }
                }
              }
            `,
            variables: {
              modReference,
            },
          })).data;
          if (!mod) {
            return { modReference, name: modReference, compatible: false };
          }
          return { modReference, name: mod.name, compatible: await isCompatibleFast(mod, this.$store.state.selectedInstall.version) };
        }));
      },
      default: [],
    },
    possiblyCompatibleInstalledCount: {
      async get() {
        return (await this.modStates).filter(({ compatible }) => compatible === COMPATIBILITY_LEVEL.POSSIBLY_COMPATIBLE).length;
      },
      default: 0,
    },
    incompatibleInstalledCount: {
      async get() {
        return (await this.modStates).filter(({ compatible }) => compatible === COMPATIBILITY_LEVEL.INCOMPATIBLE).length;
      },
      default: 0,
    },
    buttonColor: {
      async get() {
        if (this.incompatibleInstalledCount > 0) {
          return 'error';
        }
        if (this.possiblyCompatibleInstalledCount > 0) {
          return 'warning';
        }
        return 'primary';
      },
      default: 'primary',
    },
  },
  computed: {
    ...mapState(
      [
        'selectedInstall',
        'installedMods',
        'expandedModId',
        'inProgress',
        'isGameRunning',
        'installSetupError',
        'error',
        'errorPersistent',
        'launchButton',
        'launchCat',
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
    launchButtonText() {
      if (this.isGameRunning) {
        return 'Game is running';
      }
      if (!this.selectedInstall) {
        return 'Loading...';
      }
      if (this.selectedInstall.launchPath) {
        return 'Launch Satisfactory';
      }
      return 'Cannot launch this install';
    },
  },
  apollo: {
    uriInstallModData: {
      query: gql`
        query GetURIInstallMod($modReference: ModReference!){
          uriInstallModData: getModByReference(modReference: $modReference) {
            id,
            name,
            short_description,
            logo
          }
        }
      `,
      variables() {
        return {
          modReference: this.uriInstallModReference,
        };
      },
      skip() {
        return !this.uriInstallModReference;
      },
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
        this.uriInstallModReference = parsed.searchParams.get('modID');
        this.uriInstallVersion = parsed.searchParams.get('version') || null;
        this.uriInstallDialog = true;
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
              this.$electron.ipcRenderer.invoke('close');
            }
          }, 100);
        });
      } else {
        setInterval(() => {
          if (this.inProgress.length === 0) {
            this.$electron.ipcRenderer.invoke('close');
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
    uriInstallMod() {
      this.$store.dispatch('installModVersion', { modId: this.uriInstallModReference, version: this.uriInstallVersion });
      this.uriInstallDialog = false;
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
    catMouseDown(e) {
      this.catOffset = 550 - e.offsetX;
      this.catPressed = true;
    },
    catMouseMove(e) {
      if (this.catPressed) {
        this.catPosition = (e.clientX - 100 + this.catOffset) / 450;
        this.catPosition = Math.min(1, Math.max(-0.1, this.catPosition));
        if (this.catPosition === 1) {
          this.catPressed = false;
          setTimeout(() => {
            this.launchSatisfactory();
          }, 1000);
        }
      }
    },
    lastElement,
  },
};
</script>

<style scoped>
</style>
