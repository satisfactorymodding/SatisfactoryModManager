<template>
  <div>
    <v-list-item>
      <v-list-item-avatar
        tile
        style="margin-top: 0px; margin-bottom: 0px"
      >
        <v-img :src="icon" />
      </v-list-item-avatar>
      <v-list-item-content
        style="cursor: pointer; user-select: none; padding: 0;"
        @click="expandClicked"
      >
        <v-list-item-title :class="mod.isCompatible || 'error--text'">
          {{ mod.modInfo.name }}
        </v-list-item-title>
        <v-list-item-subtitle>
          <v-row style="padding-left: 12px">
            <v-col
              cols="4"
              style="padding: 0"
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
div {
  background: var(--v-backgroundModsList-base) !important;
}
.mod-button {
  margin-top: 5px;
  margin-right: 5px;
  font-size: 25px;
  opacity: 0.75;
}
.mod-button:not(:hover):not(.active)>.v-icon {
  color: var(--v-backgroundModsList-lighten2) !important;
}
.mod-button.active {
  opacity: 1 !important;
}
.mod-button:hover {
  opacity: 0.65;
}
</style>
