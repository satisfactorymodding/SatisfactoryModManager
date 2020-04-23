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
        :compatibility="compatibility"
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
        :mod="mods.find((mod) => mod.modInfo.id === expandedModId)"
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
  </v-card>
</template>

<script>
import {
  getAvailableMods, getInstalls, getConfigs, setDebug, addDownloadProgressCallback,
} from 'satisfactory-mod-manager-api';
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
          compatibility: {},
          sortBy: '',
        },
      },
      configs: [{ name: 'vanilla', items: [] }, { name: 'modded', items: [] }, { name: 'development', items: [] }],
      compatibility: [{ name: 'All mods', mods: 50 }, { name: 'Compatible', mods: 30 }],
      sortBy: ['Name, alphanumerical', 'Latest', 'Last update', 'Most popular', 'Favourite'],
      satisfactoryInstalls: [],
      selectedInstall: {},
      mods: [],
      expandedModId: '',
      favoriteModIds: [],
      installedMods: [],
      inProgress: '',
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
    addDownloadProgressCallback(this.downloadProgress);
    if (this.hasUpdate) {
      this.settingsState = 'notify';
    }
    this.favoriteModIds = getSetting('favoriteMods', []);
    this.configs = getConfigs().map((name) => ({ name, items: [] }));
    const savedConfigName = getSetting('selectedConfig', 'modded');
    this.controlData.config = this.configs.find((conf) => conf.name === savedConfigName);
    [this.controlData.filters.compatibility] = this.compatibility;
    [this.controlData.filters.sortBy] = this.sortBy;
    Promise.all([
      this.getMods(),
      getInstalls().then((installs) => {
        this.satisfactoryInstalls = installs;
        const savedLocation = getSetting('selectedInstall');
        this.selectedInstall = this.satisfactoryInstalls.find((install) => install.installLocation === savedLocation) || this.satisfactoryInstalls[0];
      }),
    ]).then(() => {
      this.refreshInstalledMods();
      this.$electron.ipcRenderer.send('vue-ready');
    });
    this.unexpandMod();
    setDebug(true);
  },
  methods: {
    getMods() {
      return getAvailableMods().then((mods) => {
        this.mods = mods.map((mod) => ({ modInfo: mod, isInstalled: false, isCompatible: true }));
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
    refreshInstalledMods() {
      this.installedMods = Object.keys(this.selectedInstall.mods);
      for (let i = 0; i < this.mods.length; i += 1) {
        this.mods[i].isInstalled = this.installedMods.includes(this.mods[i].modInfo.id);
      }
    },
    switchModInstalled(modId) {
      this.inProgress = modId;
      if (this.mods.find((mod) => mod.modInfo.id === modId).isInstalled) {
        this.selectedInstall.uninstallMod(modId).then(() => {
          this.inProgress = '';
          this.refreshInstalledMods();
        });
      } else {
        this.inProgress = modId;
        this.selectedInstall.installMod(modId).then(() => {
          this.inProgress = '';
          this.refreshInstalledMods();
        });
      }
    },
    addModToConfig(mod, config) {
      console.log(config);
    },
    removeModFromConfig(mod, config) {
      console.log(config);
    },
    downloadProgress(url, progress) {
      console.log(`Downloading ${url}: ${progress.percent}`);
    },
  },
};
</script>

<style scoped>
</style>
