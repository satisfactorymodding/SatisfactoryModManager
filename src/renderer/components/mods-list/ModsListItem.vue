<template>
  <div style="margin-top: 6px; margin-bottom: 6px; user-select: none;">
    <v-list-item
      class="mod"
      :class="{ 'expanded': isExpanded, 'in-progress': isModInProgress, 'disabled': disabled && !isModInProgress }"
      style="min-height: 45px; margin-top: 0; margin-bottom: 0; border-radius: 5px; z-index: 1;"
      two-line
    >
      <v-list-item-avatar
        tile
        style="margin-top: 0px; margin-bottom: 0px; margin-right: 8px; height: 45px; width: 45px;"
      >
        <img
          v-if="!mod.isCached || !isOffline"
          :src="icon"
        >
        <v-icon
          v-else
          style="font-size: 45px !important"
        >
          mdi-image-off-outline
        </v-icon>
      </v-list-item-avatar>
      <v-list-item-content
        style="cursor: pointer; user-select: none; padding: 0;"
        @click="expandClicked"
      >
        <v-tooltip
          ref="tooltip"
          v-model="tooltipShow"
          top
          color="background"
          :close-delay="100"
          :disabled="!errorTooltip && !disabled"
          content-class="mod-compatibility-tooltip"
        >
          <template #activator="{ on, attrs }">
            <v-list-item-title
              :class="{
                'text--text': compatibility === COMPATIBILITY_LEVEL.COMPATIBLE,
                'warning--text': compatibility === COMPATIBILITY_LEVEL.POSSIBLY_COMPATIBLE,
                'error--text': compatibility === COMPATIBILITY_LEVEL.INCOMPATIBLE,
              }"
              style="font-weight: 300"
              v-bind="attrs"
              v-on="on"
            >
              {{ mod.name }}
            </v-list-item-title>
          </template>
          <div
            @mouseenter="keepTooltipOpen = true; $refs.tooltip.clearDelay()"
            @mouseleave="keepTooltipOpen = false; $refs.tooltip.runDelay('close')"
          >
            <span v-if="errorTooltip && isOffline">You are currently offline.<br>A connection to ficsit.app is required to view details and install new mods.</span>
            <!-- eslint-disable-next-line vue/no-v-html -->
            <span v-else-if="errorTooltip"><span v-html="errorTooltip" /></span>
            <span v-if="disabled">This mod is disabled. Press the pause icon to enable it.<br></span>
          </div>
        </v-tooltip>
        <v-list-item-subtitle v-if="!isModInProgress">
          <div
            class="d-inline-flex align-center icon--text"
            style="width: 85px"
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
        <template v-if="(isEnabled || isInstalled) && !isModInProgress">
          <ModActionButton
            :disabled="isDependency || !modsEnabled || isGameRunning || inProgress.length > 0"
            background-normal-class=""
            background-hover-class="ficsitOrange"
            icon-normal-color="text"
            icon-hover-color="white"
            :normal-icon="isEnabled ? 'mdi-play' : 'mdi-pause'"
            :hover-icon="isEnabled ? 'mdi-pause' : 'mdi-play'"
            @click="isEnabled ? disable() : enable()"
          >
            <template #tooltip>
              <span v-if="isDependency">Dependency of {{ dependantsFriendly }}</span>
              <span v-else-if="!modsEnabled">Enable mods to be able to make changes</span>
              <span v-else-if="isGameRunning">Cannot install mods while game is running</span>
              <span v-else-if="inProgress.length > 0">Another operation is in progress</span>
            </template>
          </ModActionButton>
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
        <ModActionButton
          v-else
          :disabled="(versionCompatibility === COMPATIBILITY_LEVEL.INCOMPATIBLE && !isInstalled) || isDependency || !modsEnabled || isGameRunning "
          background-normal-class=""
          :background-hover-class="isInstalled || isEnabled ? 'red' : 'green'"
          :icon-normal-color="versionCompatibility === COMPATIBILITY_LEVEL.INCOMPATIBLE ? 'error' : (isInstalled || isEnabled ? 'green' : 'text')"
          icon-hover-color="white"
          :normal-icon="versionCompatibility === COMPATIBILITY_LEVEL.INCOMPATIBLE ? 'mdi-alert' : (isInstalled || isEnabled ? 'mdi-check-circle' : isInQueue ? 'mdi-tray-full' : 'mdi-download')"
          :hover-icon="isInstalled ? 'mdi-delete' : isInQueue ? 'mdi-tray-full' : isModBeingDownloaded ? 'mdi-tray-plus' : 'mdi-download'"
          @click="isInstalled ? uninstall() : install()"
        >
          <template #tooltip>
            <span v-if="isInQueue">Add to Download Queue</span>
            <span v-else-if="isDependency">Dependency of {{ dependantsFriendly }}</span>
            <span v-else-if="!modsEnabled">Enable mods to be able to make changes</span>
            <span v-else-if="isGameRunning">Cannot install mods while game is running</span>
            <span v-else-if="versionCompatibility === COMPATIBILITY_LEVEL.INCOMPATIBLE">This mod is incompatible with your game version</span>
          </template>
        </ModActionButton>
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
        :class="{ 'fast': currentModProgress.fast }"
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
import { mapGetters, mapState } from 'vuex';
import gql from 'graphql-tag';
import {
  lastElement, isCompatibleFast, COMPATIBILITY_LEVEL, smrCompatibilityToCompatibilityLevel, markdownAsHtmlText,
} from '@/utils';
import ModActionButton from './ModActionButton';

export default {
  components: {
    ModActionButton,
  },
  props: {
    mod: {
      type: Object,
      default: () => ({}),
    },
    isOffline: {
      type: Boolean,
      required: true,
    },
  },
  data() {
    return {
      COMPATIBILITY_LEVEL,
      tooltipShow: false,
      keepTooltipOpen: false,
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
    ...mapGetters([
      'branch',
    ]),
    isFavorite() {
      return this.favoriteModIds.includes(this.mod.mod_reference);
    },
    isExpanded() {
      return this.expandedModId === this.mod.mod_reference;
    },
    icon() {
      return this.mod.logo || 'https://ficsit.app/images/no_image.webp';
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
    isInQueue() {
      return this.$store.state.downloadQueue.includes(this.mod.mod_reference);
    },
    isModBeingDownloaded() {
      return this.$store.state.inProgress.length > 0
    },
    isEnabled() {
      return !!this.$store.state.installedMods[this.mod.mod_reference];
    },
    dependants() {
      return Object.entries(this.$store.state.installedMods || {}).filter(([id, data]) => !!data.dependencies[this.mod.mod_reference] && id !== this.mod.mod_reference).map(([item]) => item);
    },
    isDependency() {
      return this.dependants.length > 0;
    },
    disabled() {
      return !this.isEnabled && this.isInstalled;
    },
    reportedCompatibility() {
      if (!this.branch) return null;
      if (!this.mod.compatibility) return null;
      if (!this.mod.compatibility[this.branch]) return null;
      return smrCompatibilityToCompatibilityLevel(this.mod.compatibility[this.branch].state);
    },
    compatibility() {
      return this.reportedCompatibility ?? this.versionCompatibility;
    },
  },
  asyncComputed: {
    versionCompatibility: {
      async get() {
        if (!this.$store.state.selectedInstall) return COMPATIBILITY_LEVEL.INCOMPATIBLE;
        if (this.mod.hidden && !this.isDependency) return COMPATIBILITY_LEVEL.INCOMPATIBLE;
        return isCompatibleFast(this.mod, this.$store.state.selectedInstall.version);
      },
      default: COMPATIBILITY_LEVEL.COMPATIBLE,
    },
    errorTooltip: {
      get() {
        if (this.branch && this.reportedCompatibility) {
          switch (this.reportedCompatibility) {
            case COMPATIBILITY_LEVEL.INCOMPATIBLE:
              if (this.mod.compatibility[this.branch].note) {
                return `This mod has been reported as Broken on this game version:<br>${markdownAsHtmlText(this.mod.compatibility[this.branch].note)}`;
              }
              return 'This mod has been reported as Broken on this game version.';
            case COMPATIBILITY_LEVEL.POSSIBLY_COMPATIBLE:
              if (this.mod.compatibility[this.branch].note) {
                return `This mod has been reported as Damaged on this game version:<br>${markdownAsHtmlText(this.mod.compatibility[this.branch].note)}`;
              }
              return 'This mod has been reported as Damaged on this game version.';
            default:
              break;
          }
        }
        switch (this.versionCompatibility) {
          case COMPATIBILITY_LEVEL.COMPATIBLE:
            return null;
          case COMPATIBILITY_LEVEL.POSSIBLY_COMPATIBLE:
            return 'This mod is likely incompatible with your game version and may cause crashes.';
          case COMPATIBILITY_LEVEL.INCOMPATIBLE:
            if (!this.isDependency && this.mod.hidden) {
              return 'This mod was hidden by the author.';
            }
            return 'This mod is incompatible with your game version.';
          default:
            return null;
        }
      },
      default: null,
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
      if (!this.isOffline) {
        this.$store.dispatch('expandMod', this.mod.mod_reference);
      }
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
.mod-compatibility-tooltip {
  pointer-events: auto !important;
}
.mod-compatibility-tooltip p {
  margin-bottom: 0px !important;
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

.mod.disabled .v-avatar, .mod.disabled .v-list-item__content {
  opacity: 0.3;
  color: var(--v-backgroundModsList-lighten2) !important;
}

.mod.disabled img {
  filter: grayscale(1);
  animation-play-state: paused !important;
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
