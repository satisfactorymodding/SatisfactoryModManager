<template>
  <div style="margin-top: 6px; margin-bottom: 6px; user-select: none;">
    <v-list-item
      class="mod"
      :class="(isExpanded ? 'expanded' : '') + ' ' + (isModInProgress ? 'in-progress' : '')"
      style="min-height: 45px; margin-top: 0; margin-bottom: 0; border-radius: 5px; z-index: 1;"
    >
      <v-list-item-avatar
        tile
        style="margin-top: 0px; margin-bottom: 0px; margin-right: 8px; height: 45px; width: 45px;"
      >
        <img :src="icon">
      </v-list-item-avatar>
      <v-list-item-content
        style="cursor: pointer; user-select: none; padding: 0;"
        @click="expandClicked"
      >
        <v-list-item-title :class="mod.isCompatible || 'error--text'">
          {{ mod.modInfo.name }}
        </v-list-item-title>
        <v-list-item-subtitle v-if="!isModInProgress">
          <v-row style="padding-left: 12px">
            <v-col
              style="padding: 0; flex: 0 0 20%; max-width: 20%;"
            >
              <div class="d-inline-flex align-center">
                <v-icon
                  color="text"
                  style="padding-right: 4px"
                >
                  mdi-eye
                </v-icon>
                {{ mod.modInfo.views.toLocaleString() }}
              </div>
            </v-col>
            <v-col style="padding: 0">
              <div class="d-inline-flex align-center">
                <v-icon
                  color="text"
                  style="padding-right: 4px"
                >
                  mdi-download
                </v-icon>
                {{ mod.modInfo.downloads.toLocaleString() }}
              </div>
            </v-col>
          </v-row>
        </v-list-item-subtitle>
        <v-list-item-subtitle
          v-else
          style="z-index: 2;"
        >
          <div class="d-inline-flex align-center">
            <v-icon
              color="warning"
              style="padding-right: 4px"
            >
              mdi-information
            </v-icon>
            <span style="font-size: 11px">{{ currentModProgress.message }}</span>
          </div>
        </v-list-item-subtitle>
      </v-list-item-content>
      <v-list-item-action
        :class="isFavorite ? 'active' : ''"
        class="mod-button toggle custom"
      >
        <v-icon
          class="icon"
          color="warning"
          @click="favoriteClicked"
        >
          mdi-star
        </v-icon>
      </v-list-item-action>
      <v-list-item-action
        class="custom mod-button"
      >
        <v-icon
          v-if="!mod.isCompatible"
          color="error"
        >
          mdi-alert
        </v-icon>
        <v-icon
          v-else-if="isModInProgress"
          color="warning"
        >
          mdi-sync
        </v-icon>
        <v-icon
          v-else-if="mod.isInstalled && (mod.isDependency || !canInstallMods || inProgress.length > 0)"
          color="green"
        >
          mdi-checkbox-marked
        </v-icon>
        <v-icon
          v-else-if="mod.isInstalled"
          color="green"
          @click="switchInstalled"
        >
          mdi-checkbox-marked
        </v-icon>
        <v-icon
          v-else-if="!canInstallMods || inProgress.length > 0"
          color="text"
        >
          mdi-checkbox-blank-outline
        </v-icon>
        <v-icon
          v-else
          color="text"
          @click="switchInstalled"
        >
          mdi-checkbox-blank-outline
        </v-icon>
      </v-list-item-action>
    </v-list-item>
    <v-list-item
      v-if="isModInProgress"
      style="height: 0px; min-height: 0px; padding: 0; z-index: 0"
    >
      <v-progress-linear
        :value="Math.round(currentModProgress.progress * 100)"
        :class="currentModProgress.fast ? 'fast' : ''"
        color="mod-list-progress-bar"
        height="45"
        reactive
        style="position: relative; top: -22.5px; margin-left: 45px"
        :indeterminate="currentModProgress.progress < 0"
      />
    </v-list-item>
  </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex';
import { lastElement } from '@/utils';

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
    isExpanded() {
      return this.expandedModId === this.mod.modInfo.mod_reference;
    },
    icon() {
      return this.mod.modInfo.logo || 'https://ficsit.app/static/assets/images/no_image.png';
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
    switchInstalled() {
      this.$store.dispatch('switchModInstalled', this.mod.modInfo.mod_reference);
    },
  },
};
</script>

<style>
.mod-list-progress-bar {
  background-color: var(--v-backgroundModsList-darken1) !important;
  border-color: var(--v-backgroundModsList-darken1) !important;
}
</style>

<style scoped>
.v-icon {
  transition: background-color 0ms !important;
}
.v-list-item__subtitle .v-icon {
  font-size: 14px !important
}
.custom.v-list-item__action {
  margin-top: 0;
  margin-bottom: 0;
}
.mod-button {
  margin-top: 5px;
  margin-left: 0 !important;
  margin-right: 5px;
}
.mod-button.toggle {
  opacity: 0.75;
}
.mod-button.toggle:not(:hover):not(.active)>.v-icon {
  color: var(--v-backgroundModsList-lighten2) !important;
}
.mod-button.active {
  opacity: 1 !important;
}
.mod-button:hover {
  opacity: 0.65;
}

.mod, .mod * {
  background: var(--v-backgroundModsList-base) !important;
}
.mod.in-progress, .mod.in-progress * {
  background: transparent !important;
}
.expanded, .expanded *{
  background-color: var(--v-backgroundModsList-darken1) !important;
}
</style>
