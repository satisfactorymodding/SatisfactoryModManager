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
        :state="settingsState"
        style="user-select: none;"
        @settingsClicked="settingsClicked"
      />
      <ControlArea
        :selected-install="selectedInstall"
        :selected-config="selectedConfig"
        :selected-filters="filters"
        :installs="satisfactoryInstalls"
        :configs="configs"
        :mod-filters="modFilters"
        :sort-by="sortBy"
        :in-progress="inProgress"
        style="user-select: none;"
        @selectedInstallChanged="selectedInstall = $event"
        @selectedConfigChanged="selectedConfig = $event"
        @selectedFiltersChanged="filters = $event"
      />
      <ModsList
        :mods="filteredMods"
        :expanded-mod-id="expandedModId"
        :favorite-mod-ids="favoriteModIds"
        :in-progress="inProgress"
        @expandMod="expandMod"
        @favoriteMod="favoriteMod"
        @unfavoriteMod="unfavoriteMod"
        @switchMod="switchModInstalled"
      />
      <v-btn
        block
        tile
        color="primary"
        dark
        elevation="0"
        height="82px"
        style="font-size: 18px;"
        :disabled="!!inProgress.length"
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
      <ModDetails
        v-if="expandedModId"
        :mod="mods.find((mod) => mod.modInfo.mod_reference === expandedModId)"
        :is-favorite="favoriteModIds.includes(expandedModId)"
        :in-progress="inProgress"
        :configs="configs"
        @close="unexpandMod"
        @favoriteMod="favoriteMod"
        @unfavoriteMod="unfavoriteMod"
        @switchMod="switchModInstalled"
        @addToConfig="addModToConfig(expandedModId, $event)"
        @removeFromConfig="removeModFromConfig(expandedModId, $event)"
      />
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
            @click="error = ''"
          >
            Ok
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog
      hide-overlay
      persistent
      :value="inProgress.some((prog) => prog.id === '__loadingApp__')"
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
          v-if="inProgress.some((prog) => prog.id === '__loadingApp__')"
          class="text-center"
        >
          <v-progress-linear
            :value="Math.round(lastElement(inProgress.find((prog) => prog.id === '__loadingApp__').progresses).progress * 100)"
            :class="lastElement(inProgress.find((prog) => prog.id === '__loadingApp__').progresses).fastUpdate ? 'fast' : ''"
            background-color="#000000"
            color="#5bb71d"
            height="2"
            reactive
            :indeterminate="lastElement(inProgress.find((prog) => prog.id === '__loadingApp__').progresses).progress < 0"
          />
          {{ lastElement(inProgress.find((prog) => prog.id === '__loadingApp__').progresses).message || '&nbsp;' }}
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script>
// TODO: Loading screen
import {
  getAvailableMods, getModsCount, MODS_PER_PAGE, getAvailableSMLVersions,
  getInstalls, getConfigs,
  setDebug, addDownloadProgressCallback,
  loadCache,
} from 'satisfactory-mod-manager-api';
import { satisfies, coerce, valid } from 'semver';
import TitleBar from './TitleBar';
import ControlArea from './ControlArea';
import ModsList from './ModsList';
import ModDetails from './ModDetails';
import { getSetting, saveSetting } from '../settings';
import { lastElement } from '../utils';

export default {
  components: {
    TitleBar,
    ControlArea,
    ModsList,
    ModDetails,
  },
  data() {
    return {
      settingsState: 'off',
      hasUpdate: false,
      selectedConfig: {},
      filters: {
        modFilters: {},
        sortBy: '',
        search: '',
      },
      configs: [],
      modFilters: [{ name: 'All mods', mods: 0 }, { name: 'Compatible', mods: 0 }, { name: 'Favourite', mods: 0 }, { name: 'Installed', mods: 0 }, { name: 'Not installed', mods: 0 }],
      sortBy: ['Name', 'Last updated', 'Popularity', 'Hotness', 'Views', 'Downloads'],
      satisfactoryInstalls: [],
      selectedInstall: {},
      smlVersions: [],
      mods: [],
      expandedModId: '',
      favoriteModIds: [],
      installedMods: [],
      manifestMods: [],
      inProgress: [], // { id: string, progresses: { id: string, progress: number, message: string, fast: boolean }[] }
      currentDownloadProgress: {},
      error: '',
    };
  },
  computed: {
    filteredMods() {
      let filtered;
      if (this.filters.modFilters === this.modFilters[1]) filtered = this.mods.filter((mod) => mod.isCompatible);
      else if (this.filters.modFilters === this.modFilters[2]) filtered = this.mods.filter((mod) => this.favoriteModIds.includes(mod.modInfo.mod_reference));
      else if (this.filters.modFilters === this.modFilters[3]) filtered = this.mods.filter((mod) => mod.isInstalled);
      else if (this.filters.modFilters === this.modFilters[4]) filtered = this.mods.filter((mod) => !mod.isInstalled);
      else filtered = [...this.mods];

      if (this.filters.search !== '') {
        filtered = filtered.filter((mod) => mod.modInfo.name.toLowerCase().includes(this.filters.search.toLowerCase())); // TODO: maybe search in description too
      }

      if (this.filters.sortBy === 'Name') filtered = filtered.sort((a, b) => a.modInfo.name.localeCompare(b.modInfo.name));
      if (this.filters.sortBy === 'Last updated') filtered = filtered.sort((a, b) => b.modInfo.last_version_date - a.modInfo.last_version_date);
      if (this.filters.sortBy === 'Popularity') filtered = filtered.sort((a, b) => b.modInfo.popularity - a.modInfo.popularity);
      if (this.filters.sortBy === 'Hotness') filtered = filtered.sort((a, b) => b.modInfo.hotness - a.modInfo.hotness);
      if (this.filters.sortBy === 'Views') filtered = filtered.sort((a, b) => b.modInfo.views - a.modInfo.views);
      if (this.filters.sortBy === 'Downloads') filtered = filtered.sort((a, b) => b.modInfo.downloads - a.modInfo.downloads);

      return filtered;
    },
  },
  watch: {
    filters: {
      deep: true,
      handler() {
        saveSetting('filters', { modFilters: this.filters.modFilters.name, sortBy: this.filters.sortBy });
      },
    },
    selectedInstall() {
      if (!this.inProgress.some((prog) => prog.id === '__loadingApp__')) {
        const loadProgress = {
          id: '__loadingApp__',
          progresses: [{
            id: '', progress: -1, message: 'Validating mod install', fast: false,
          }],
        };
        this.inProgress.push(loadProgress);
        const savedConfigName = this.getSavedConfigName();
        this.selectedConfig = this.configs.find((conf) => conf.name === savedConfigName);
        this.selectedInstall.setConfig(savedConfigName).then(() => {
          this.refreshModsInstalledCompatible();
          this.inProgress.remove(loadProgress);
        }).catch((e) => {
          this.showError(e);
          this.inProgress.remove(loadProgress);
        });
        saveSetting('selectedInstall', this.selectedInstall.installLocation);
      }
    },
    selectedConfig(newValue) {
      if (!this.inProgress.some((prog) => prog.id === '__loadingApp__')) {
        const loadProgress = {
          id: '__loadingApp__',
          progresses: [{
            id: '', progress: -1, message: 'Validating mod install', fast: false,
          }],
        };
        this.inProgress.push(loadProgress);
        this.selectedInstall.setConfig(newValue.name).then(() => {
          this.refreshModsInstalledCompatible();
          this.inProgress.remove(loadProgress);
          this.saveSelectedConfig();
        }).catch((e) => {
          this.showError(e);
          this.inProgress.remove(loadProgress);
        });
      }
    },
  },
  mounted() {
    const appLoadProgress = {
      id: '__loadingApp__',
      progresses: [{
        id: '', progress: -1, message: 'Loading', fast: false,
      }],
    };
    this.inProgress.push(appLoadProgress);
    setDebug(true);
    addDownloadProgressCallback(this.downloadProgress);
    if (this.hasUpdate) {
      this.settingsState = 'notify';
    }
    this.favoriteModIds = getSetting('favoriteMods', []);
    this.configs = getConfigs();

    const savedFilters = getSetting('filters', { modFilters: this.modFilters[0].name, sortBy: this.filters.sortBy[0] });
    this.filters.modFilters = this.modFilters.find((modFilter) => modFilter.name === savedFilters.modFilters) || this.modFilters[0];
    this.filters.sortBy = this.sortBy.find((item) => item === savedFilters.sortBy) || this.sortBy[0];
    Promise.all([
      loadCache().then(() => getInstalls().then((installs) => {
        const installValidateProgress = { id: 'validatingInstall', progress: -1, message: 'Validating mod install' };
        appLoadProgress.progresses.push(installValidateProgress);
        this.satisfactoryInstalls = installs;
        const savedLocation = getSetting('selectedInstall');
        this.selectedInstall = this.satisfactoryInstalls.find((install) => install.installLocation === savedLocation) || this.satisfactoryInstalls[0];
        const savedConfigName = this.getSavedConfigName();
        this.selectedConfig = this.configs.find((conf) => conf.name === savedConfigName);

        return this.selectedInstall.setConfig(this.getSavedConfigName()).then(() => appLoadProgress.progresses.remove(installValidateProgress));
      })),
      this.getSMLVersions(),
      this.getAllMods(appLoadProgress),
    ]).catch((e) => {
      this.showError(e);
    }).then(() => {
      this.modFilters[0].mods = this.mods.length;
      this.modFilters[2].mods = this.mods.filter((mod) => this.favoriteModIds.includes(mod.modInfo.mod_reference)).length;
      this.refreshModsInstalledCompatible();
      this.$electron.ipcRenderer.send('vue-ready');
      this.inProgress.remove(appLoadProgress);
    });
    this.unexpandMod();
  },
  methods: {
    getSMLVersions() {
      return getAvailableSMLVersions().then((versions) => {
        this.smlVersions = versions;
      });
    },
    getAllMods(progress) {
      const getModsProgress = { id: 'getMods', progress: -1, message: 'Getting available mods' };
      if (progress) {
        progress.progresses.push(getModsProgress);
      }
      let modsGot = 0;
      return getModsCount().then((modCount) => {
        getModsProgress.message = `Getting available mods (${modsGot}/${modCount})`;
        getModsProgress.progress = 0;
        const modPages = Math.ceil(modCount / MODS_PER_PAGE);
        return Promise.all(Array.from({ length: modPages }).map((_, i) => getAvailableMods(i).then((mods) => {
          modsGot += mods.length;
          getModsProgress.progress += 1 / modPages;
          getModsProgress.message = `Getting available mods (${modsGot}/${modCount})`;
          return mods;
        })))
          .then((pages) => pages.flat(1))
          .then((mods) => {
            this.mods = mods.map((mod) => ({
              modInfo: mod,
              isInstalled: false,
              isCompatible: true,
              isDependency: false,
            }));
            if (progress) {
              progress.progresses.remove(getModsProgress);
            }
          });
      });
    },
    settingsClicked() {
      if (this.settingsState !== 'on') {
        this.settingsState = 'on';
      } else if (this.hasUpdate) {
        this.settingsState = 'notify';
      } else {
        this.settingsState = 'off';
      }
    },
    expandMod(modId) {
      this.expandedModId = modId;
      this.$electron.ipcRenderer.send('expand');
    },
    unexpandMod() {
      this.expandedModId = '';
      this.$electron.ipcRenderer.send('unexpand');
    },
    favoriteMod(modId) {
      if (!this.favoriteModIds.includes(modId)) {
        this.favoriteModIds.push(modId);
        this.modFilters[2].mods = this.mods.filter((mod) => this.favoriteModIds.includes(mod.modInfo.mod_reference)).length;
        saveSetting('favoriteMods', this.favoriteModIds);
      }
    },
    unfavoriteMod(modId) {
      this.favoriteModIds.remove(modId);
      this.modFilters[2].mods = this.mods.filter((mod) => this.favoriteModIds.includes(mod.modInfo.mod_reference)).length;
      saveSetting('favoriteMods', this.favoriteModIds);
    },
    refreshModsInstalledCompatible() {
      this.installedMods = Object.keys(this.selectedInstall.mods);
      this.manifestMods = this.selectedInstall.manifestMods;
      for (let i = 0; i < this.mods.length; i += 1) {
        this.mods[i].isInstalled = this.installedMods.includes(this.mods[i].modInfo.mod_reference);
        this.mods[i].isDependency = this.installedMods.includes(this.mods[i].modInfo.mod_reference) && !this.manifestMods.includes(this.mods[i].modInfo.mod_reference);
        this.mods[i].isCompatible = this.mods[i].modInfo.versions.length > 0
        && !!this.mods[i].modInfo.versions.find((ver) => satisfies(ver.sml_version, '>=2.0.0')
              && this.smlVersions.some((smlVer) => valid(coerce(smlVer.version)) === valid(coerce(ver.sml_version)))
              && satisfies(valid(coerce(this.selectedInstall.version)), `>=${valid(coerce(this.smlVersions.find((smlVer) => valid(coerce(smlVer.version)) === valid(coerce(ver.sml_version))).satisfactory_version))}`));
      }
      this.modFilters[1].mods = this.mods.filter((mod) => mod.isCompatible).length;
      this.modFilters[3].mods = this.mods.filter((mod) => mod.isInstalled).length;
      this.modFilters[4].mods = this.mods.filter((mod) => !mod.isInstalled).length;
    },
    switchModInstalled(modId) {
      if (this.inProgress.length > 0) {
        this.showError('Another operation is currently in progress');
        return;
      }
      this.currentDownloadProgress = [];
      const modProgress = { id: modId, progresses: [] };
      this.inProgress.push(modProgress);
      const placeholderProgreess = {
        id: 'placeholder', progress: -1, message: '', fastUpdate: false,
      };
      modProgress.progresses.push(placeholderProgreess);
      if (this.mods.find((mod) => mod.modInfo.mod_reference === modId).isInstalled) {
        placeholderProgreess.message = 'Checking for mods that are no longer needed';
        this.selectedInstall.uninstallMod(modId).then(() => {
          placeholderProgreess.progress = 1;
          setTimeout(() => {
            this.inProgress.remove(modProgress);
          }, 500);
          this.refreshModsInstalledCompatible();
        }).catch((e) => {
          this.inProgress.remove(modProgress);
          this.showError(e);
        });
      } else {
        placeholderProgreess.message = 'Finding the best version to install';
        this.selectedInstall.installMod(modId).then(() => {
          this.inProgress.progress = 1;
          setTimeout(() => {
            this.inProgress.remove(modProgress);
          }, 500);
          this.refreshModsInstalledCompatible();
        }).catch((e) => {
          this.inProgress.remove(modProgress);
          this.showError(e);
        });
      }
    },
    // TODO
    addModToConfig(mod, config) {
      console.log(config);
    },
    removeModFromConfig(mod, config) {
      console.log(config);
    },
    saveSelectedConfig() {
      let current = getSetting('selectedConfig', {});
      if (typeof current !== 'object') { current = {}; }
      current[this.selectedInstall.installLocation] = this.selectedConfig.name;
      saveSetting('selectedConfig', current);
    },
    getSavedConfigName() {
      const current = getSetting('selectedConfig', {});
      return current[this.selectedInstall.installLocation] || 'modded';
    },
    downloadProgress(url, progress, friendlyName) {
      if (!this.currentDownloadProgress[url]) {
        this.currentDownloadProgress[url] = {
          id: `download_${url}`, progress: 0, message: '', fastUpdate: true,
        };
        this.inProgress[0].progresses.push(this.currentDownloadProgress[url]);
      }
      this.currentDownloadProgress[url].message = `Downloading ${friendlyName} ${Math.round(progress.percent * 100)}%`;
      this.currentDownloadProgress[url].progress = progress.percent;
      if (progress.percent === 1) {
        setTimeout(() => {
          this.inProgress[0].progresses.remove(this.currentDownloadProgress[url]);
          delete this.currentDownloadProgress[url];
        }, 100);
      }
    },
    showError(e) {
      this.error = typeof e === 'string' ? e : e.message;
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
