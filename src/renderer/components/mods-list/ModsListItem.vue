<template>
  <div style="margin-top: 6px; margin-bottom: 6px; user-select: none;">
    <v-list-item
      class="mod"
      :class="(isExpanded ? 'expanded' : '') + ' ' + (isModInProgress ? 'in-progress' : '')"
      style="min-height: 45px; margin-top: 0; margin-bottom: 0; border-radius: 5px; z-index: 1;"
      two-line
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
        <v-list-item-title :class="isCompatible || 'error--text'">
          {{ mod.name }}
        </v-list-item-title>
        <v-list-item-subtitle v-if="!isModInProgress">
          <div
            class="d-inline-flex align-center"
            style="width: 20%"
          >
            <v-icon
              color="text"
              style="padding-right: 4px"
            >
              mdi-eye
            </v-icon>
            {{ mod.views.toLocaleString() }}
          </div>
          <div class="d-inline-flex align-center">
            <v-icon
              color="text"
              style="padding-right: 4px"
            >
              mdi-download
            </v-icon>
            {{ mod.downloads.toLocaleString() }}
          </div>
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
        <v-tooltip
          v-if="isInstalled && (isDependency || !canInstallMods || inProgress.length > 0)"
          color="background"
          left
        >
          <template #activator="{ on, attrs }">
            <v-icon
              color="green"
              v-bind="attrs"
              v-on="on"
            >
              mdi-checkbox-marked
            </v-icon>
          </template>
          <span>Dependency of {{ dependantsFriendly }}</span>
        </v-tooltip>
        <v-icon
          v-else-if="isInstalled"
          color="green"
          @click="switchInstalled"
        >
          mdi-checkbox-marked
        </v-icon>
        <v-icon
          v-else-if="!isCompatible"
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
import { lastElement, isCompatibleFast } from '@/utils';
import gql from 'graphql-tag';

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
      return this.favoriteModIds.includes(this.mod.mod_reference);
    },
    isExpanded() {
      return this.expandedModId === this.mod.mod_reference;
    },
    icon() {
      return this.mod.logo || 'https://ficsit.app/static/assets/images/no_image.png';
    },
    isModInProgress() {
      return !!this.modProgress;
    },
    modProgress() {
      return this.inProgress.find((prog) => prog.id === this.mod.mod_reference);
    },
    currentModProgress() {
      return lastElement(this.modProgress.progresses);
    },
    isInstalled() {
      return !!this.$store.state.installedMods[this.mod.mod_reference];
    },
    dependants() {
      return Object.entries(this.$store.state.installedMods || {}).filter(([, data]) => !!data.dependencies[this.mod.mod_reference]).map(([item]) => item);
    },
    isDependency() {
      return this.$store.state.manifestMods && !this.$store.state.manifestMods[this.mod.mod_reference] && this.dependants.length > 0;
    },
  },
  asyncComputed: {
    isCompatible: {
      get() {
        if (!this.$store.state.selectedInstall) return false;
        return isCompatibleFast(this.mod, this.$store.state.selectedInstall.version);
      },
      default: true,
    },
    dependantsFriendly: {
      async get() {
        const friendlyNames = await Promise.all(this.dependants.map(async (item) => (await this.$apollo.query({
          query: gql`
            query getModName($modReference: ModReference!) {
              mod: getModByReference(modReference: $modReference) {
                id,
                name,
              }
            }
          `,
          variables: {
            modReference: item,
          },
        })).data.mod.name));
        friendlyNames.sort();
        if (friendlyNames.length === 0) { return 'none'; }
        let finalString = friendlyNames[0];
        let i = 1;
        while (i < friendlyNames.length && finalString.length + `, ${friendlyNames[i]}`.length <= 40) {
          finalString += `, ${friendlyNames[i]}`;
          i += 1;
        }
        if (i < friendlyNames.length) {
          finalString += ` and ${friendlyNames.length - i} more`;
        }
        return finalString;
      },
      default: '',
    },
  },
  methods: {
    expandClicked() {
      this.$store.dispatch('expandMod', this.mod.mod_reference);
    },
    favoriteClicked() {
      this.$store.dispatch('toggleModFavorite', this.mod.mod_reference);
    },
    switchInstalled() {
      this.$store.dispatch('switchModInstalled', this.mod.mod_reference);
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
