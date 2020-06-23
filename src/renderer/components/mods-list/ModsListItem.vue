<template>
  <div
    class="px-3"
  >
    <v-list-item>
      <v-list-item-content
        style="cursor: pointer; user-select: none;"
        @click="expandClicked"
      >
        <span :class="mod.isCompatible || 'error--text'">{{ mod.modInfo.name }}</span>
      </v-list-item-content>
      <v-list-item-action
        class="mod-button custom"
        :class="expandedModId === mod.modInfo.mod_reference ? 'active' : ''"
      >
        <v-icon
          class="ma-1 icon"
          color="primary"
          @click="expandClicked"
        >
          mdi-import
        </v-icon>
      </v-list-item-action>
      <v-list-item-action
        :class="isFavorite ? 'active' : ''"
        class="mod-button custom"
      >
        <v-icon
          class="ma-1 icon"
          color="warning"
          @click="favoriteClicked"
        >
          mdi-star
        </v-icon>
      </v-list-item-action>
      <v-list-item-action
        class="custom"
        style="margin-right: 10px"
      >
        <v-switch
          v-model="mod.isInstalled"
          inset
          dense
          color="primary"
          class="custom pr-1"
          :class="mod.isCompatible ? '' : 'incompatible'"
          flat
          :disabled="!mod.isCompatible || mod.isDependency || !!inProgress.length || !canInstallMods"
          @click.stop.prevent="switchClicked"
        />
      </v-list-item-action>
    </v-list-item>
    <v-list-item
      v-if="isModInProgress"
      style="height: 0px; min-height: 0px; padding: 0;"
    >
      <v-progress-linear
        :value="Math.round(currentModProgress.progress * 100)"
        :class="currentModProgress.fast ? 'fast' : ''"
        color="warning"
        height="49"
        reactive
        style="position: relative; top: -24.5px;"
        :indeterminate="currentModProgress.progress < 0"
      >
        <strong>{{ currentModProgress.message }}</strong>
      </v-progress-linear>
    </v-list-item>
  </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex';
import { lastElement } from '../../utils';

export default {
  props: {
    mod: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {

    };
  },
  computed: {
    ...mapState([
      'favoriteModIds',
      'expandedModId',
      'inProgress',
    ]),
    ...mapGetters([
      'canInstallMods',
    ]),
    isFavorite() {
      return this.favoriteModIds.includes(this.mod.modInfo.mod_reference);
    },
    isModInProgress() {
      return !!this.modProgress;
    },
    modProgress() {
      return this.inProgress.find((prog) => prog.id === this.mod.modInfo.mod_reference);
    },
    currentModProgress() {
      return lastElement(this.modProgress.progresses);
    },
  },
  methods: {
    expandClicked() {
      this.$store.dispatch('expandMod', this.mod.modInfo.mod_reference);
    },
    favoriteClicked() {
      this.$store.dispatch('toggleModFavorite', this.mod.modInfo.mod_reference);
    },
    switchClicked() {
      this.$store.dispatch('switchModInstalled', this.mod.modInfo.mod_reference);
    },
  },
};
</script>

<style scoped>
.custom.v-list-item__action {
  margin-top: 0;
  margin-bottom: 0;
}
.v-divider {
  border-color: var(--v-background-darken3) !important;
}
div {
  background: var(--v-backgroundSecondary-base) !important;
}
.mod-button {
  margin-top: 5px;
  margin-right: 5px;
  font-size: 25px;
  opacity: 0.75;
}
.mod-button:not(:hover):not(.active)>.v-icon {
  color: var(--v-backgroundSecondary-lighten2) !important;
}
.mod-button.active {
  opacity: 1 !important;
}
.mod-button:hover {
  opacity: 0.65;
}

.list-shadow-top, .list-shadow-bottom {
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  z-index: 1;
  background: transparent !important;
  pointer-events: none;
}

.list-shadow-top {
  box-shadow: inset 0px 45px 20px -20px rgba(0,0,0,0.3);
}
.list-shadow-bottom {
  box-shadow: inset 0px -45px 20px -20px rgba(0,0,0,0.3);
}
</style>
