<template>
  <v-card
    tile
    flat
    :class="{'color-2':isSelected(mod), 'color-3 mod-disable':!isCompatibleMod(mod), 'mod-error':hasUpdateMod(mod)}"
  >
    <div
      class="d-flex d-flex align-items-center p-2"
      @click="selectMod(mod)"
    >
      <div class="fs-075">
        {{ mod.name }}
      </div>
      <v-spacer />
      <v-tooltip bottom>
        <template v-slot:activator="{ on }">
          <v-icon
            v-show="isSelected(mod) || isModDetailOpened()"
            class="clickable-icon ml-1"
            :class="isModDetailOpened() ? 'active' : ''"
            :disable="isModDetailOpened()"
            v-on="on"
            @click="openModDetail(mod)"
          >
            $showMoreIcon
          </v-icon>
        </template>
        <span>{{ 'Mod Detail' }}</span>
      </v-tooltip>
      <v-tooltip bottom>
        <template v-slot:activator="{ on }">
          <v-icon
            v-show="isSelected(mod)"
            class="clickable-icon ml-1"
            v-on="on"
          >
            $starOffIcon
          </v-icon>
        </template>
        <span>{{ 'Mod Detail' }}</span>
      </v-tooltip>
      <v-switch
        v-model="installStatus"
        inset
        hide-details
        value
        color="success warning"
        class="p-0 m-0 ml-2"
        :loading="(inProgress.length > 0 && inProgress[0].id === mod.id) || configLoadInProgress"
        :disabled="!mod.versions[0] || (!isInstalledMod(mod) && !isCompatibleMod(mod)) || inProgress.length > 0 || configLoadInProgress || selectedConfig === 'vanilla'"
      />
    </div>
    <v-divider />
  </v-card>
</template>

<script>
import { mapState } from 'vuex';
import {
  isModInstalled,
  isModCompatible,
  isModHasUpdate,
} from '../utils/helper';

export default {
  name: 'Mod',
  props: {
    mod: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      installStatus: false,
    };
  },
  computed: {
    ...mapState([
      'selectedConfig',
      'inProgress',
      'configLoadInProgress',
      'selectedMod',
      'isSidePanelOpen',
      'sidePanelData',
    ]),
  },
  watch: {
    installStatus(value) {
      if ((value && !this.isInstalledMod(this.mod)) || (!value && this.isInstalledMod(this.mod))) {
        this.$store.dispatch('installUninstallUpdate', this.mod);
      }
    },
  },
  created() {
    this.installStatus = this.isInstalledMod(this.mod) && !this.hasUpdateMod(this.mod);
  },
  updated() {
    this.installStatus = this.isInstalledMod(this.mod) && !this.hasUpdateMod(this.mod);
  },
  methods: {
    selectMod(mod) {
      if (!this.isSelected(mod)) {
        this.$store.state.selectedMod = mod;
      }
    },
    openModDetail(mod) {
      this.$eventBus.$emit('open-side-panel', mod);
    },
    isInstalledMod(mod) {
      return isModInstalled(this.$store.state, mod);
    },
    isCompatibleMod(mod) {
      return isModCompatible(mod);
    },
    hasUpdateMod(mod) {
      return isModHasUpdate(this.$store.state, mod);
    },
    isSelected(mod) {
      return this.selectedMod.id === mod.id;
    },
    isModDetailOpened() {
      return this.isSidePanelOpen && this.sidePanelData && this.sidePanelData.id === this.mod.id;
    },
  },
};
</script>
