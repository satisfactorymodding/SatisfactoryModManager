<template>
  <v-card
    tile
  >
    <v-card
      tile
      flat
      class="mod-list-wrapper launcher-scroll overflow-y-auto"
    >
      <Mod
        v-for="mod in searchMods"
        :key="mod.id"
        :mod="mod"
      />
    </v-card>
    <v-card
      v-if="isLoading || configLoadInProgress"
      class="mod-loading loading fixed-top w-100 h-100 d-flex align-items-center justify-content-center"
    >
      <transition name="loading">
        <v-progress-circular
          indeterminate
        />
      </transition>
    </v-card>
    <ModUpdate />
    <ModInstallUninstall />
  </v-card>
</template>

<script>
import arrayEqual from 'array-equal';
import Mod from './Mod';
import ModUpdate from './dialogs/ModUpdate';
import ModInstallUninstall from './dialogs/ModInstallUninstall';

export default {
  name: 'ModList',
  components: {
    Mod,
    ModUpdate,
    ModInstallUninstall,
  },
  data() {
    return {
      selectedIndex: 0,
      isLoading: true,
    };
  },
  computed: {
    availableMods() {
      return this.$store.state.availableMods;
    },
    searchMods() {
      return this.$store.state.searchMods;
    },
    selectedMod: {
      get() {
        return this.$store.state.selectedMod;
      },
      set(value) {
        this.$store.state.selectedMod = value;
      },
    },
    configLoadInProgress() {
      return this.$store.state.configLoadInProgress;
    },
  },
  watch: {
    searchMods(newObjects, oldObjects) {
      if (!arrayEqual(newObjects, oldObjects)) {
        if (this.searchMods.length > 0) {
          const index = Math.min(Math.max(this.selectedIndex, 0), this.searchMods.length - 1);
          this.selectMod(this.searchMods[index]);
          this.isLoading = false;
        } else {
          this.selectMod(null);
        }
      }
    },
  },
  methods: {
    selectMod(mod) {
      this.selectedMod = mod;
      this.selectedIndex = this.searchMods.indexOf(mod);
    },
  },
};
</script>
