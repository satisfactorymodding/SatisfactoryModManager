<template>
  <v-card class="h-100 d-flex">
    <v-card
      class="h-100"
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
  </v-card>
</template>

<script>
import 'satisfactory-mod-manager-api';
import TitleBar from './TitleBar';
import ControlArea from './ControlArea';
import ModsList from './ModsList';

export default {
  components: {
    TitleBar,
    ControlArea,
    ModsList,
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
      configs: [{ name: 'vanilla' }, { name: 'modded' }, { name: 'development' }],
      compatibility: [{ name: 'All mods', mods: 50 }, { name: 'Compatible', mods: 30 }],
      sortBy: ['Name, alphanumerical', 'Latest', 'Last update', 'Most popular', 'Favourite'],
      mods: [
        {
          name: 'Test0', id: 'Test0', isCompatible: false,
        },
        {
          name: 'Test1', id: 'Test1', isCompatible: true,
        },
        {
          name: 'Test2', id: 'Test2', isCompatible: true,
        },
        {
          name: 'Test3', id: 'Test3', isCompatible: false,
        },
        {
          name: 'Test4', id: 'Test4', isCompatible: true,
        },
        {
          name: 'Test5', id: 'Test5', isCompatible: false,
        },
        {
          name: 'Test6', id: 'Test6', isCompatible: true,
        },
        {
          name: 'Test7', id: 'Test7', isCompatible: false,
        },
        {
          name: 'Test8', id: 'Test8', isCompatible: true,
        },
        {
          name: 'Test9', id: 'Test9', isCompatible: true,
        },
        {
          name: 'Testa', id: 'Testa', isCompatible: true,
        },
      ],
      expandedModId: '',
      favoriteModIds: [],
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
    if (this.hasUpdate) {
      this.settingsState = 'notify';
    }
    [this.controlData.config] = this.configs;
    [this.controlData.filters.compatibility] = this.compatibility;
    [this.controlData.filters.sortBy] = this.sortBy;
    this.$electron.ipcRenderer.send('vue-ready');
  },
  methods: {
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
      }
    },
    unfavoriteMod(modId) {
      this.favoriteModIds.remove(modId);
    },
  },
};
</script>

<style scoped>
</style>
