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
        <v-list-item-title
          :class="isCompatible ? 'text--text' : 'error--text'"
          style="font-weight: 300"
        >
          {{ mod.name }}
        </v-list-item-title>
        <v-list-item-subtitle v-if="!isModInProgress">
          <div
            class="d-inline-flex align-center icon--text"
            style="width: 25%"
          >
            <v-icon
              color="icon"
              style="padding-right: 4px"
            >
              mdi-eye
            </v-icon>
            {{ mod.views.toLocaleString() }}
          </div>
          <div class="d-inline-flex align-center icon--text">
            <v-icon
              color="icon"
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
        v-if="isInstalled || isEnabled"
        class="custom mod-button d-inline-flex align-center justify-center"
        style="width: 50px"
      >
        <template v-if="isEnabled">
          <v-hover v-slot="{ hover }">
            <v-icon
              v-if="!hover"
              color="text"
            >
              mdi-play
            </v-icon>
            <template v-else-if="isDependency">
              <v-tooltip
                color="background"
                left
              >
                <template #activator="{ on, attrs }">
                  <v-icon
                    color="text"
                    v-bind="attrs"
                    v-on="on"
                  >
                    mdi-play
                  </v-icon>
                </template>
                <span>Dependency of {{ dependantsFriendly }}</span>
              </v-tooltip>
            </template>
            <v-tooltip
              v-else-if="!modsEnabled"
              color="background"
              left
            >
              <template #activator="{ on, attrs }">
                <v-icon
                  color="text"
                  v-bind="attrs"
                  v-on="on"
                >
                  mdi-play
                </v-icon>
              </template>
              <span>Enable mods to be able to make changes</span>
            </v-tooltip>
            <v-tooltip
              v-else-if="isGameRunning"
              color="background"
              left
            >
              <template #activator="{ on, attrs }">
                <v-icon
                  color="text"
                  v-bind="attrs"
                  v-on="on"
                >
                  mdi-play
                </v-icon>
              </template>
              <span>Cannot install mods while game is running</span>
            </v-tooltip>
            <div
              v-else
              class="d-inline-flex align-center justify-center"
              style="height: 30px; width: 30px; background: var(--v-ficsitOrange-base) !important;"
            >
              <v-icon
                color="white"
                style="background-color: unset !important"
                @click="disable"
              >
                mdi-pause
              </v-icon>
            </div>
          </v-hover>
        </template>
        <template v-else-if="isInstalled">
          <v-hover v-slot="{ hover }">
            <v-icon
              v-if="!hover"
              color="text"
            >
              mdi-pause
            </v-icon>
            <v-tooltip
              v-else-if="!modsEnabled"
              color="background"
              left
            >
              <template #activator="{ on, attrs }">
                <v-icon
                  color="text"
                  v-bind="attrs"
                  v-on="on"
                >
                  mdi-pause
                </v-icon>
              </template>
              <span>Enable mods to be able to make changes</span>
            </v-tooltip>
            <v-tooltip
              v-else-if="isGameRunning"
              color="background"
              left
            >
              <template #activator="{ on, attrs }">
                <v-icon
                  color="text"
                  v-bind="attrs"
                  v-on="on"
                >
                  mdi-pause
                </v-icon>
              </template>
              <span>Cannot install mods while game is running</span>
            </v-tooltip>
            <div
              v-else
              class="d-inline-flex align-center justify-center primary"
              style="height: 30px; width: 30px"
            >
              <v-icon
                color="white"
                style="background-color: unset !important"
                @click="enable"
              >
                mdi-play
              </v-icon>
            </div>
          </v-hover>
        </template>
      </v-list-item-action>
      <v-list-item-action
        class="custom mod-button d-inline-flex align-center justify-center"
        style="width: 50px"
      >
        <v-icon
          v-if="isModInProgress"
          color="warning"
        >
          mdi-sync
        </v-icon>
        <template v-else-if="isEnabled && isDependency">
          <v-tooltip
            color="background"
            left
          >
            <template #activator="{ on, attrs }">
              <v-icon
                color="green"
                v-bind="attrs"
                v-on="on"
              >
                mdi-check-circle
              </v-icon>
            </template>
            <span>Dependency of {{ dependantsFriendly }}</span>
          </v-tooltip>
        </template>
        <template v-else-if="isInstalled">
          <v-hover v-slot="{ hover }">
            <v-icon
              v-if="!hover"
              color="green"
            >
              mdi-check-circle
            </v-icon>
            <v-tooltip
              v-else-if="!modsEnabled"
              color="background"
              left
            >
              <template #activator="{ on, attrs }">
                <v-icon
                  color="green"
                  v-bind="attrs"
                  v-on="on"
                >
                  mdi-check-circle
                </v-icon>
              </template>
              <span>Enable mods to be able to make changes</span>
            </v-tooltip>
            <v-tooltip
              v-else-if="isGameRunning"
              color="background"
              left
            >
              <template #activator="{ on, attrs }">
                <v-icon
                  color="green"
                  v-bind="attrs"
                  v-on="on"
                >
                  mdi-check-circle
                </v-icon>
              </template>
              <span>Cannot install mods while game is running</span>
            </v-tooltip>
            <div
              v-else
              class="d-inline-flex align-center justify-center red"
              style="height: 30px; width: 30px"
            >
              <v-icon
                color="text"
                style="background-color: unset !important"
                @click="uninstall"
              >
                mdi-delete
              </v-icon>
            </div>
          </v-hover>
        </template>
        <template v-else>
          <v-icon
            v-if="!isCompatible"
            color="error"
          >
            mdi-alert
          </v-icon>
          <v-tooltip
            v-else-if="!modsEnabled"
            color="background"
            left
          >
            <template #activator="{ on, attrs }">
              <v-icon
                color="text"
                v-bind="attrs"
                v-on="on"
              >
                mdi-download
              </v-icon>
            </template>
            <span>Enable mods to be able to make changes</span>
          </v-tooltip>
          <v-tooltip
            v-else-if="isGameRunning"
            color="background"
            left
          >
            <template #activator="{ on, attrs }">
              <v-icon
                class="icon"
                v-bind="attrs"
                v-on="on"
              >
                mdi-download
              </v-icon>
            </template>
            <span>Cannot install mods while game is running</span>
          </v-tooltip>
          <v-icon
            v-else-if="inProgress.length > 0"
            color="icon"
          >
            mdi-download
          </v-icon>
          <div
            v-else
            class="d-inline-flex align-center justify-center hover-green"
            style="height: 30px; width: 30px"
            @click="install"
          >
            <v-icon
              class="icon"
            >
              mdi-download
            </v-icon>
          </div>
        </template>
      </v-list-item-action>
      <v-list-item-action
        class="mod-button custom d-inline-flex align-center justify-center"
        style="width: 50px"
      >
        <div
          class="d-inline-flex align-center justify-center hover-yellow"
          style="height: 30px; width: 30px"
          @click="favoriteClicked"
        >
          <v-icon
            class="icon"
            :color="isFavorite ? 'yellow' : 'icon'"
          >
            mdi-star
          </v-icon>
        </div>
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
import { mapState } from 'vuex';
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
      'modsEnabled',
      'isGameRunning',
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
    manifestItem() {
      return this.$store.state.manifestItems ? this.$store.state.manifestItems.find((item) => item.id === this.mod.mod_reference) : undefined;
    },
    isInstalled() {
      return !!this.manifestItem;
    },
    isEnabled() {
      return !!this.$store.state.installedMods[this.mod.mod_reference];
    },
    dependants() {
      return Object.entries(this.$store.state.installedMods || {}).filter(([, data]) => !!data.dependencies[this.mod.mod_reference]).map(([item]) => item);
    },
    isDependency() {
      return this.dependants.length > 0;
    },
  },
  asyncComputed: {
    isCompatible: {
      get() {
        if (!this.$store.state.selectedInstall) return false;
        if (this.mod.hidden && !this.isDependency) return false;
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
    install() {
      this.$store.dispatch('installMod', this.mod.mod_reference);
    },
    uninstall() {
      this.$store.dispatch('uninstallMod', this.mod.mod_reference);
    },
    enable() {
      this.$store.dispatch('enableMod', this.mod.mod_reference);
    },
    disable() {
      this.$store.dispatch('disableMod', this.mod.mod_reference);
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

.mod, .mod * {
  background: var(--v-backgroundModsList-base) !important;
}
.mod.in-progress, .mod.in-progress * {
  background: transparent !important;
}
.expanded, .expanded * {
  background-color: var(--v-backgroundModsList-darken1) !important;
}

.v-application .icon {
  color: var(--v-icon-base);
  background-color: unset !important;
}

.hover-green:hover {
  background-color: #4caf50 !important;
}
.hover-green:hover .icon {
  color: var(--v-text-base) !important;
}

.hover-yellow:hover {
  background-color: #ffeb3b !important;
}
.hover-yellow:hover .icon {
  color: var(--v-text-base) !important;
}

.hover-red:hover {
  background-color: #f44336 !important;
}
</style>
