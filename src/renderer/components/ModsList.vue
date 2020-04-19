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
          v-for="(mod, index) in mods"
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
                <span>{{ mod.name }}</span>
              </v-list-item-content>
              <v-list-item-icon>
                <div
                  class="d-inline-flex align-items-center mod-button"
                  :class="expandedModId === mod.id ? 'active' : ''"
                  fill-height
                  @click="expandClicked(mod)"
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
                  :class="favoriteModIds.includes(mod.id) ? 'active' : ''"
                  fill-height
                  @click="favoriteClicked(mod)"
                >
                  <v-icon
                    class="ma-1 icon"
                    color="#ffc107"
                  >
                    mdi-star
                  </v-icon>
                </div>
                <v-switch
                  inset
                  dense
                  color="primary"
                  class="custom"
                  :disabled="!mod.isCompatible"
                  @change="switchClicked(mod)"
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
  },
  data() {
    return {
      selectedModIdx: -1,
    };
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
  opacity: 1;
}
.mod-button:hover {
  opacity: 1;
}

.list-shadow {
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0px;
  box-shadow: inset 0px 60px 30px -20px rgba(0,0,0,0.35), inset 0px -60px 30px -20px rgba(0,0,0,0.35);
  z-index: 1;
  background: transparent !important;
  pointer-events:none;
}
</style>
