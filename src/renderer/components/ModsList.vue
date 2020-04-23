<template>
  <div
    class="py-4"
    style="position: relative;"
  >
    <div class="list-shadow" />
    <div
      style="overflow-y: scroll;"
      class="mx-2"
    >
      <v-list
        class="py-1 my-n4"
        height="475px"
      >
        <template
          v-for="({ modInfo, isInstalled, isCompatible }, index) in mods"
        >
          <div
            :key="index"
            class="px-3"
          >
            <v-divider
              v-if="index != 0"
            />

            <v-list-item>
              <v-list-item-content>
                <span>{{ modInfo.name }}</span>
              </v-list-item-content>
              <v-list-item-icon>
                <div
                  class="d-inline-flex align-items-center mod-button"
                  :class="expandedModId === modInfo.id ? 'active' : ''"
                  fill-height
                  @click="expandClicked(modInfo)"
                >
                  <v-icon
                    class="ma-1 icon"
                    color="primary"
                  >
                    mdi-import
                  </v-icon>
                </div>
                <div
                  class="d-inline-flex align-items-center mod-button"
                  :class="favoriteModIds.includes(modInfo.id) ? 'active' : ''"
                  style="margin-right: 15px"
                  fill-height
                  @click="favoriteClicked(modInfo)"
                >
                  <v-icon
                    class="ma-1 icon"
                    color="warning"
                  >
                    mdi-star
                  </v-icon>
                </div>
                <v-switch
                  :value="isInstalled"
                  inset
                  dense
                  color="primary"
                  class="custom"
                  :class="isCompatible ? '' : 'incompatible'"
                  flat
                  :disabled="!isCompatible || !!inProgress"
                  @click.stop.prevent="switchClicked(modInfo)"
                />
              </v-list-item-icon>
            </v-list-item>
          </div>
        </template>
      </v-list>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    mods: {
      type: Array,
      default() { return []; },
    },
    favoriteModIds: {
      type: Array,
      default() { return []; },
    },
    expandedModId: {
      type: String,
      default: '',
    },
    inProgress: {
      type: String,
      default: '',
    },
    progressPercent: {
      type: Number,
      default: 0,
    },
  },
  methods: {
    expandClicked(mod) {
      this.$emit('expandMod', mod.id);
    },
    favoriteClicked(mod) {
      if (!this.favoriteModIds.includes(mod.id)) {
        this.$emit('favoriteMod', mod.id);
      } else {
        this.$emit('unfavoriteMod', mod.id);
      }
    },
    switchClicked(mod) {
      this.$emit('switchMod', mod.id);
    },
  },
};
</script>

<style scoped>
div {
  background: var(--v-backgroundSecondary-base) !important;
}
.mod-button {
  margin-top: 5px;
  margin-right: 5px;
  font-size: 25px;
  opacity: 0.2;
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

.list-shadow {
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0px;
  box-shadow: inset 0px 45px 20px -20px rgba(0,0,0,0.3), inset 0px -45px 20px -20px rgba(0,0,0,0.3);
  z-index: 1;
  background: transparent !important;
  pointer-events: none;
}
</style>
