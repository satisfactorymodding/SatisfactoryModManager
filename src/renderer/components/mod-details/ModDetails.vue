<!-- eslint-disable vue/no-v-html -->
<template>
  <v-card
    tile
    class="d-flex"
    width="100%"
    height="100%"
    style="padding: 0; box-shadow: inset 10px 0 10px -10px rgba(0,0,0,1); max-width: calc(100vw - 550px);"
  >
    <v-overlay
      v-if="!$apollo.queries.mod.loading && (!mod || mod.mod_reference !== expandedModId)"
      absolute
    >
      <v-btn
        color="backgroundModsList-darken1"
        raised
        width="100%"
        @click="close"
      >
        <v-icon
          left
          dark
        >
          mdi-flip-h mdi-import
        </v-icon>
        <span>Offline</span>
      </v-btn>
    </v-overlay>
    <v-overlay
      v-else-if="!mod || mod.mod_reference !== expandedModId"
      absolute
    >
      Loading...
    </v-overlay>
    <ModDetailsInfo
      v-if="mod"
      :mod="mod"
      :is-installed="isInstalled"
      @close="close"
      @install-version="installVersion"
    />
    <v-card
      v-if="mod"
      height="100%"
      width="100%"
      class="mod-description-card"
    >
      <v-container
        fluid
        class="d-flex flex-column"
        style="height: 100%; padding: 20px 12x 15px 20px;"
      >
        <v-row style="flex-basis: 0; overflow: auto;">
          <div
            ref="modDescription"
            class="mod-description"
            v-html="modDescription"
          />
        </v-row>
        <v-row
          style="flex-grow: 0;"
          class="pb-2"
        >
          <v-spacer />
          <v-col style="flex-grow: 0; padding-bottom: 0">
            <v-btn
              raised
              width="150px"
              @click="favoriteClicked"
            >
              <span style="vertical-align: middle;">
                {{ isFavorite ? 'Unfavourite' : 'Favourite' }}
              </span>
              <v-spacer />
              <v-icon
                right
                :color="isFavorite ? 'warning' : 'text'"
              >
                mdi-star
              </v-icon>
            </v-btn>
          </v-col>
          <v-col
            v-if="isCompatible || isPossiblyCompatible"
            style="flex-grow: 0; padding-bottom: 0"
            class="d-flex"
          >
            <v-btn
              v-if="isEnabled"
              width="150px"
              raised
              :disabled="!!inProgress.length || isGameRunning || !modsEnabled"
              @click="disable"
            >
              <span style="vertical-align: middle;">
                Disable mod
              </span>
              <v-spacer />
              <v-icon
                right
                color="red"
              >
                mdi-close-circle
              </v-icon>
            </v-btn>
            <template v-else-if="isInstalled">
              <v-btn
                width="150px"
                raised
                :disabled="!!inProgress.length || isGameRunning || !modsEnabled"
                class="mr-2"
                @click="uninstall"
              >
                <span style="vertical-align: middle;">
                  Remove mod
                </span>
                <v-spacer />
                <v-icon
                  right
                  color="red"
                >
                  mdi-delete
                </v-icon>
              </v-btn>
              <v-btn
                width="150px"
                raised
                :disabled="!!inProgress.length || isGameRunning || !modsEnabled"
                @click="enable"
              >
                <span style="vertical-align: middle;">
                  Enable mod
                </span>
                <v-spacer />
                <v-icon
                  right
                  color="green"
                >
                  mdi-check-circle-outline
                </v-icon>
              </v-btn>
            </template>
            <v-btn
              v-else
              width="150px"
              raised
              :disabled="!!inProgress.length || isGameRunning || !modsEnabled"
              @click="install"
            >
              <span style="vertical-align: middle;">
                Install mod
              </span>
              <v-spacer />
              <v-icon
                right
                color="green"
              >
                mdi-download
              </v-icon>
            </v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-card>
    <v-dialog
      v-model="enlargedImageDialog"
      max-width="70%"
    >
      <v-img :src="enlargedImage" />
    </v-dialog>
  </v-card>
</template>

<script>
import { mapState } from 'vuex';
import gql from 'graphql-tag';
import { markdownAsElement, isCompatibleFast, COMPATIBILITY_LEVEL } from '@/utils';
import ModDetailsInfo from './ModDetailsInfo';

export default {
  components: {
    ModDetailsInfo,
  },
  data() {
    return {
      enlargedImage: '',
    };
  },
  computed: {
    ...mapState([
      'inProgress',
      'expandedModId',
      'isGameRunning',
      'modsEnabled',
    ]),
    isFavorite() {
      return this.$store.state.favoriteModIds.includes(this.expandedModId);
    },
    descriptionAsElement() {
      return markdownAsElement(this.mod.full_description || '');
    },
    modDescription() {
      const el = this.descriptionAsElement;
      const links = el.getElementsByTagName('a');
      for (let i = 0; i < links.length; i += 1) {
        links[i].target = '_blank';
      }
      return el.innerHTML;
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
      return Object.entries(this.$store.state.installedMods || {}).filter(([id, data]) => !!data.dependencies[this.mod.mod_reference] && id !== this.mod.mod_reference).map(([item]) => item);
    },
    isDependency() {
      return this.dependants.length > 0;
    },
    enlargedImageDialog: {
      get() {
        return !!this.enlargedImage;
      },
      set(value) {
        if (!value) {
          this.enlargedImage = '';
        }
      },
    },
  },
  asyncComputed: {
    isCompatible: {
      async get() {
        if (!this.$store.state.selectedInstall) return false;
        if (this.mod?.hidden && !this.isDependency) return false;
        return (await isCompatibleFast(this.mod, this.$store.state.selectedInstall.version)) === COMPATIBILITY_LEVEL.COMPATIBLE;
      },
      default: false,
    },
    isPossiblyCompatible: {
      async get() {
        if (!this.$store.state.selectedInstall) return false;
        if (this.mod?.hidden && !this.isDependency) return false;
        return (await isCompatibleFast(this.mod, this.$store.state.selectedInstall.version)) === COMPATIBILITY_LEVEL.POSSIBLY_COMPATIBLE;
      },
      default: false,
    },
  },
  apollo: {
    mod: {
      query: gql`
        query getModDetails($modReference: ModReference!) {
          mod: getModByReference(modReference: $modReference)
          {
            id,
            name,
            logo,
            mod_reference,
            full_description,
            created_at,
            last_version_date,
            downloads,
            views,
            hidden,
            authors {
              user {
                username,
                avatar
              }
              role
            },
            versions(filter: {limit: 100}) {
              id,
              version,
              sml_version,
              size
            }
          }
        }
      `,
      variables() {
        return {
          modReference: this.expandedModId,
        };
      },
      pollInterval: 5 * 60 * 1000,
    },
  },
  watch: {
    mod() {
      this.$nextTick(() => {
        const imgs = this.$refs.modDescription.getElementsByTagName('img');
        for (let i = 0; i < imgs.length; i += 1) {
          const img = imgs[i];
          img.onclick = () => {
            this.enlargedImage = img.src;
          };
        }
      });
    },
  },
  methods: {
    close() {
      this.$store.dispatch('unexpandMod');
    },
    install() {
      this.$store.dispatch('installMod', this.expandedModId);
    },
    uninstall() {
      this.$store.dispatch('uninstallMod', this.expandedModId);
    },
    enable() {
      this.$store.dispatch('enableMod', this.expandedModId);
    },
    disable() {
      this.$store.dispatch('disableMod', this.expandedModId);
    },
    favoriteClicked() {
      this.$store.dispatch('toggleModFavorite', this.expandedModId);
    },
    installVersion(version) {
      this.$store.dispatch('installModVersion', { modId: this.expandedModId, version });
    },
  },
};
</script>

<style>
.mod-description img {
  max-width: 100%;
}
.mod-description {
  font-size: 13.5px;
}
.v-application .mod-description p {
  margin-bottom: 5px;
}
</style>

<style scoped>
.v-icon {
  font-size: 20px !important;
}
.v-card {
  background-color: var(--v-backgroundModsList-base) !important;
}
.v-btn {
  color: var(--v-text-base) !important;
  font-size: 16px !important;
  padding-left: 8px !important;
  padding-right: 12px !important;
}

.mod-description {
  display: block;
  word-break: break-word;
  overflow-y: auto;
  box-shadow: none;
}
.mod-description-card {
  max-width: calc(100% - 296px);
}
@media (max-height: 700px) {
  .mod-description-card {
    max-width: calc(100% - 260px);
  }
}

.row {
  margin-left: 0 !important;
  margin-right: 0 !important;
}
</style>
