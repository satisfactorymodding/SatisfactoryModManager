<template>
  <div
    class="py-4"
    style="position: relative;"
  >
    <div
      v-if="topShadow"
      class="list-shadow-top"
    />
    <div
      v-if="bottomShadow"
      class="list-shadow-bottom"
    />
    <div
      ref="modsList"
      style="overflow-y: scroll; height: 100%;"
      class="mx-4"
      @scroll="modsListScrolled"
    >
      <v-list
        class="pt-1 mt-n4 custom"
      >
        <template
          v-for="({ modInfo, isInstalled, isCompatible, isDependency }, index) in mods"
        >
          <div
            :key="index"
            class="px-3"
          >
            <v-divider
              v-if="index != 0"
            />

            <v-list-item>
              <v-list-item-content
                style="cursor: pointer; user-select: none;"
                @click="expandClicked(modInfo)"
              >
                <span :class="isCompatible || 'error--text'">{{ modInfo.name }}</span>
              </v-list-item-content>
              <v-list-item-action
                class="mod-button custom"
                :class="expandedModId === modInfo.mod_reference ? 'active' : ''"
              >
                <v-icon
                  class="ma-1 icon"
                  color="primary"
                  @click="expandClicked(modInfo)"
                >
                  mdi-import
                </v-icon>
              </v-list-item-action>
              <v-list-item-action
                :class="favoriteModIds.includes(modInfo.mod_reference) ? 'active' : ''"
                class="mod-button custom"
              >
                <v-icon
                  class="ma-1 icon"
                  color="warning"
                  @click="favoriteClicked(modInfo)"
                >
                  mdi-star
                </v-icon>
              </v-list-item-action>
              <v-list-item-action
                class="custom"
                style="margin-right: 10px"
              >
                <v-switch
                  v-model="mods[index].isInstalled"
                  inset
                  dense
                  color="primary"
                  class="custom pr-1"
                  :class="isCompatible ? '' : 'incompatible'"
                  flat
                  :disabled="!isCompatible || isDependency || !!inProgress.length || !canInstallMods"
                  @click.stop.prevent="switchClicked(modInfo)"
                />
              </v-list-item-action>
            </v-list-item>
            <v-list-item
              v-if="inProgress.some((prog) => prog.id === modInfo.mod_reference)"
              style="height: 0px; min-height: 0px; padding: 0;"
            >
              <v-progress-linear
                :value="Math.round(currentModProgress(modInfo.mod_reference).progress * 100)"
                :class="currentModProgress(modInfo.mod_reference).fast ? 'fast' : ''"
                color="warning"
                height="49"
                reactive
                style="position: relative; top: -24.5px;"
                :indeterminate="currentModProgress(modInfo.mod_reference).progress < 0"
              >
                <strong>{{ currentModProgress(modInfo.mod_reference).message }}</strong>
              </v-progress-linear>
            </v-list-item>
          </div>
        </template>
      </v-list>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex';
import { lastElement } from '../utils';

export default {
  data() {
    return {
      topShadow: false,
      bottomShadow: true,
    };
  },
  computed: {
    ...mapState([
      'favoriteModIds',
      'expandedModId',
      'inProgress',
    ]),
    ...mapGetters({
      mods: 'filteredMods',
      canInstallMods: 'canInstallMods',
    }),
  },
  watch: {
    mods() {
      setTimeout(() => this.modsListScrolled(), 1);
    },
  },
  methods: {
    expandClicked(mod) {
      this.$store.dispatch('expandMod', mod.mod_reference);
    },
    favoriteClicked(mod) {
      this.$store.dispatch('toggleModFavorite', mod.mod_reference);
    },
    switchClicked(mod) {
      this.$store.dispatch('switchModInstalled', mod.mod_reference);
    },
    modProgress(mod) {
      return this.inProgress.find((prog) => prog.id === mod);
    },
    currentModProgress(mod) {
      return lastElement(this.modProgress(mod).progresses);
    },
    modsListScrolled() {
      this.topShadow = this.$refs.modsList.scrollTop > 0;
      this.bottomShadow = this.$refs.modsList.scrollTop + this.$refs.modsList.offsetHeight < this.$refs.modsList.scrollHeight;
    },
    lastElement,
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
