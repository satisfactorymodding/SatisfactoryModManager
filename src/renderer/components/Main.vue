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
        :disabled="!!inProgress.id"
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
      :value="inProgress.id === '__loadingApp__'"
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

        <v-card-text class="text-center">
          <v-progress-linear
            :value="Math.round(inProgress.progress * 100)"
            background-color="#000000"
            color="#5bb71d"
            height="2"
            reactive
          />
          {{ inProgress.message || '&nbsp;' }}
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script>
// TODO: Loading screen
import {
  getAvailableMods, getInstalls, getConfigs, setDebug, addDownloadProgressCallback, getAvailableSMLVersions,
  loadCache,
} from 'satisfactory-mod-manager-api';
import { satisfies, coerce, valid } from 'semver';
import TitleBar from './TitleBar';
import ControlArea from './ControlArea';
import ModsList from './ModsList';
import ModDetails from './ModDetails';
import { getSetting, saveSetting } from '../settings';

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
      modFilters: [{ name: 'All mods', mods: 0 }, { name: 'Compatible', mods: 0 }, { name: 'Favourite', mods: 0 }],
      sortBy: ['Name', 'Last updated', 'Popularity', 'Hotness', 'Views', 'Downloads'],
      satisfactoryInstalls: [],
      selectedInstall: {},
      smlVersions: [],
      mods: [],
      expandedModId: '',
      favoriteModIds: [],
      installedMods: [],
      manifestMods: [],
      inProgress: { id: '', progress: 0, message: '' },
      currentDownloadProgress: {},
      error: '',
      fakeProgress: 0,
    };
  },
  computed: {
    filteredMods() {
      let filtered;
      if (this.filters.modFilters === this.modFilters[1]) filtered = this.mods.filter((mod) => mod.isCompatible);
      else if (this.filters.modFilters === this.modFilters[2]) filtered = this.mods.filter((mod) => this.favoriteModIds.includes(mod.modInfo.mod_reference));
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
      if (this.inProgress.id !== '__loadingApp__') {
        this.inProgress.id = '__loadingApp__';
        const savedConfigName = this.getSavedConfigName();
        this.selectedConfig = this.configs.find((conf) => conf.name === savedConfigName);
        this.selectedInstall.setConfig(savedConfigName).then(() => {
          this.refreshModsInstalledCompatible();
          this.inProgress.id = '';
        }).catch((e) => {
          this.showError(e);
          this.inProgress.id = '';
        });
        saveSetting('selectedInstall', this.selectedInstall.installLocation);
      }
    },
    selectedConfig(newValue) {
      if (this.inProgress.id !== '__loadingApp__') {
        this.inProgress.id = '__loadingApp__';
        this.selectedInstall.setConfig(newValue.name).then(() => {
          this.refreshModsInstalledCompatible();
          this.inProgress.id = '';
          this.saveSelectedConfig();
        }).catch((e) => {
          this.showError(e);
          this.inProgress.id = '';
        });
      }
    },
  },
  mounted() {
    this.inProgress.id = '__loadingApp__';
    setDebug(true);
    addDownloadProgressCallback(this.downloadProgress);
    if (this.hasUpdate) {
      this.settingsState = 'notify';
    }
    this.favoriteModIds = getSetting('favoriteMods', []);
    this.configs = getConfigs();

    const savedFilters = getSetting('filters', { modFilters: this.modFilters[0].name, sortBy: this.filters.sortBy[0] });
    console.log(savedFilters);
    this.filters.modFilters = this.modFilters.find((modFilter) => modFilter.name === savedFilters.modFilters) || this.modFilters[0];
    this.filters.sortBy = this.sortBy.find((item) => item === savedFilters.sortBy) || this.sortBy[0];
    Promise.all([
      loadCache(),
      this.getSMLVersions(),
      this.getAllMods(),
      getInstalls().then((installs) => {
        this.satisfactoryInstalls = installs;
        const savedLocation = getSetting('selectedInstall');
        this.selectedInstall = this.satisfactoryInstalls.find((install) => install.installLocation === savedLocation) || this.satisfactoryInstalls[0];
        const savedConfigName = this.getSavedConfigName();
        this.selectedConfig = this.configs.find((conf) => conf.name === savedConfigName);

        return this.selectedInstall.setConfig(this.getSavedConfigName());
      }),
    ]).catch((e) => {
      this.showError(e);
    }).then(() => {
      this.modFilters[0].mods = this.mods.length;
      this.modFilters[2].mods = this.mods.filter((mod) => this.favoriteModIds.includes(mod.modInfo.mod_reference)).length;
      this.refreshModsInstalledCompatible();
      this.$electron.ipcRenderer.send('vue-ready');
      this.inProgress.id = '';
    });
    this.unexpandMod();
  },
  methods: {
    getSMLVersions() {
      return getAvailableSMLVersions().then((versions) => {
        this.smlVersions = versions;
      });
    },
    getMods(page) {
      return getAvailableMods(page || 0).then((mods) => {
        if (mods.length > 0) {
          return this.getMods((page || 0) + 1).then((otherMods) => mods.concat(otherMods));
        }
        return mods;
      });
    },
    getAllMods() {
      return this.getMods().then((mods) => {
        this.mods = mods.map((mod) => ({
          modInfo: mod,
          isInstalled: false,
          isCompatible: true,
          isDependency: false,
        }));
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
    },
    switchModInstalled(modId) {
      if (this.inProgress.id) {
        this.showError('Another operation is currently in progress');
        return;
      }
      let message = '';
      // fake progress because manifest mutations don't return progress
      this.currentDownloadProgress = {};
      this.fakeProgress = 0;
      const interval = setInterval(() => {
        if (Object.keys(this.currentDownloadProgress).length === 0) {
          this.fakeProgress += (1 - this.fakeProgress) * 0.0025;
          this.inProgress.progress = this.fakeProgress;
          this.inProgress.message = message;
        }
      }, 50);
      if (this.mods.find((mod) => mod.modInfo.mod_reference === modId).isInstalled) {
        message = 'Checking for mods that are no longer needed';
        this.inProgress = { id: modId, progress: 0, message };
        this.selectedInstall.uninstallMod(modId).then(() => {
          this.inProgress.progress = 1;
          setTimeout(() => {
            this.inProgress = { id: '', progress: 0 };
          }, 500);
          clearInterval(interval);
          this.refreshModsInstalledCompatible();
        }).catch((e) => {
          this.inProgress = { id: '', progress: 0 };
          clearInterval(interval);
          this.showError(e);
        });
      } else {
        message = 'Finding the best version to install';
        this.inProgress = { id: modId, progress: 0, message };
        this.selectedInstall.installMod(modId).then(() => {
          this.inProgress.progress = 1;
          setTimeout(() => {
            this.inProgress = { id: '', progress: 0 };
          }, 500);
          clearInterval(interval);
          this.refreshModsInstalledCompatible();
        }).catch((e) => {
          this.inProgress = { id: '', progress: 0 };
          clearInterval(interval);
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
      this.currentDownloadProgress[url] = { friendlyName, progress: progress.percent };
      const newProgress = (Object.keys(this.currentDownloadProgress)
        .map((key) => this.currentDownloadProgress[key].progress).reduce((prev, current) => prev + current)
      / Object.keys(this.currentDownloadProgress).length);
      if (Math.abs(newProgress - this.inProgress.progress) >= 0.005 || newProgress >= 0.95) {
        this.inProgress.progress = newProgress;
      }
      this.inProgress.message = `Downloading ${Object.keys(this.currentDownloadProgress).map((key) => `${this.currentDownloadProgress[key].friendlyName} ${Math.round(this.currentDownloadProgress[key].progress * 100)}%`).join(', ')}`;
      if (progress.percent === 1) {
        setTimeout(() => {
          delete this.currentDownloadProgress[url];
        }, 100);
      }
    },
    showError(e) {
      this.error = typeof e === 'string' ? e : e.message;
    },
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
