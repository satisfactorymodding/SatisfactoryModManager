<template>
  <v-card
    class="d-flex"
    height="100%"
  >
    <v-card
      height="100%"
      style="width: 500px; min-width: 500px; max-width: 500px;"
    >
      <TitleBar
        title="Satisfactory Mod Manager"
        :state="settingsState"
        @settingsClicked="settingsClicked"
      />
      <ControlArea
        v-model="controlData"
        :configs="configs"
        :mod-filters="modFilters"
        :sort-by="sortBy"
      />
      <ModsList
        :mods="mods"
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
  </v-card>
</template>

<script>
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
      controlData: {
        config: {},
        filters: {
          modFilters: {},
          sortBy: '',
        },
      },
      configs: [{ name: 'vanilla', items: [] }, { name: 'modded', items: [] }, { name: 'development', items: [] }],
      modFilters: [{ name: 'All mods', mods: 50 }, { name: 'Compatible', mods: 30 }, { name: 'Favourite', mods: 30 }],
      sortBy: ['Name', 'Last updated', 'Popularity', 'Downloads'],
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
  watch: {
    controlData: {
      deep: true,
      handler: (newValue) => {
        console.log(newValue);
      },
    },
  },
  mounted() {
    setDebug(true);
    addDownloadProgressCallback(this.downloadProgress);
    if (this.hasUpdate) {
      this.settingsState = 'notify';
    }
    this.favoriteModIds = getSetting('favoriteMods', []);
    this.configs = getConfigs().map((name) => ({ name, items: [] }));
    const savedConfigName = getSetting('selectedConfig', 'modded');
    this.controlData.config = this.configs.find((conf) => conf.name === savedConfigName);
    [this.controlData.filters.modFilters] = this.modFilters;
    [this.controlData.filters.sortBy] = this.sortBy;
    this.inProgress.id = '__loadingApp__';
    Promise.all([
      loadCache(),
      this.getSMLVersions(),
      this.getAllMods(),
      getInstalls().then((installs) => {
        this.satisfactoryInstalls = installs;
        const savedLocation = getSetting('selectedInstall');
        this.selectedInstall = this.satisfactoryInstalls.find((install) => install.installLocation === savedLocation) || this.satisfactoryInstalls[0];

        return this.selectedInstall.manifestMutate([], [], []);
      }),
    ]).then(() => {
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
        saveSetting('favoriteMods', this.favoriteModIds);
      }
    },
    unfavoriteMod(modId) {
      this.favoriteModIds.remove(modId);
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
    addModToConfig(mod, config) {
      console.log(config);
    },
    removeModFromConfig(mod, config) {
      console.log(config);
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
</style>
