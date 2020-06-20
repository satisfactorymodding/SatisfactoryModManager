<!-- eslint-disable vue/no-v-html -->
<template>
  <v-card
    tile
    class="d-flex flex-column"
    width="100%"
    height="100%"
    style="padding: 0; box-shadow: inset 10px 0 10px -10px rgba(0,0,0,1), inset 0 10px 10px -10px rgba(0,0,0,1); max-width: calc(100vw - 500px);"
  >
    <v-row
      no-gutters
      class="flex-nowrap"
      style="padding-top: 16px; flex: 0;"
    >
      <v-col cols="auto">
        <img
          :src="icon"
          class="mod-icon"
        >
      </v-col>
      <v-col style="max-width: calc(100% - 320px)">
        <v-row
          no-gutters
        >
          <v-col cols="2">
            <span class="header">By:</span>
          </v-col>
          <v-col cols="6">
            <span class="authors">{{ authors }}</span>
          </v-col>
        </v-row>
        <v-row
          no-gutters
          class="pt-6 pb-1"
        >
          <v-col cols="auto">
            <span class="header">{{ mod.modInfo.name }}</span>
          </v-col>
        </v-row>
        <v-row
          class="mod-description"
          :class="expandDetails ? 'expanded' : ''"
          :style="expandDetails ? `height: ${images && images.length > 0 ? windowHeight - 220 : windowHeight - 167}px;` : ''"
        >
          <div
            v-if="!expandDetails"
            class="expended-description-shadow"
          />
          <div v-html="modDescription" />
        </v-row>
        <v-row
          v-if="images && images.length > 0"
          class="expand-details-button"
          :class="expandDetails ? 'expanded' : ''"
          @click="toggleExpandDetails"
        >
          <v-col
            class="expand-details-text"
          >
            <span>{{ expandDetails ? 'See gallery' : 'Full description' }}<br><v-icon v-if="!expandDetails">mdi-chevron-down</v-icon><v-icon v-else>mdi-chevron-up</v-icon></span>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
    <v-row
      class="control-bar"
      align="center"
      style="flex: 0"
    >
      <v-col
        cols="auto"
        @click="close"
      >
        <v-icon>
          mdi-flip-h mdi-import
        </v-icon>
        <span
          class="mx-2"
        >Close</span>
      </v-col>
      <v-spacer />
      <v-col
        cols="auto"
      >
        <v-menu
          bottom
          offset-y
        >
          <template v-slot:activator="{ on }">
            <v-btn
              text
              :disabled="!!inProgress.length || !mod.isInstalled"
              v-on="on"
            >
              <span
                class="mx-2"
              >Versions</span>
              <v-icon>
                mdi-chevron-down
              </v-icon>
            </v-btn>
          </template>
          <v-list>
            <v-list-item
              @click="installVersion('')"
            >
              <v-list-item-action>
                <v-icon v-if="!mod.manifestVersion">
                  mdi-check
                </v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Latest</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <template
              v-for="(version, i) in mod.modInfo.versions"
            >
              <v-list-item
                :key="i"
                @click="installVersion(version.version)"
              >
                <v-list-item-action>
                  <v-icon v-if="validAndEq(mod.manifestVersion, version.version)">
                    mdi-check
                  </v-icon>
                </v-list-item-action>
                <v-list-item-content>
                  <v-list-item-title>{{ version.version }}</v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </template>
          </v-list>
        </v-menu>
      </v-col>
      <!--<v-col
        cols="auto"
      >
        <span
          class="mx-2"
        >Report mod</span>
        <v-icon color="warning">
          mdi-alert
        </v-icon>
      </v-col>-->
      <v-col
        cols="auto"
        align="center"
      >
        <v-row
          align="center"
        >
          <v-col
            cols="auto"
            align="center"
          >
            <span>Installed</span>
          </v-col>
          <v-col
            cols="auto"
            align="center"
          >
            <v-switch
              v-model="mod.isInstalled"
              inset
              dense
              color="primary"
              class="custom"
              :class="mod.isCompatible ? '' : 'incompatible'"
              flat
              :disabled="!mod.isCompatible || mod.isDependency || !!inProgress.length || !canInstallMods"
              @click.stop.prevent="switchClicked(mod.modInfo)"
            />
          </v-col>
        </v-row>
      </v-col>
      <!--<v-col
        cols="auto"
        @click="close"
      >
        <v-menu
          bottom
          offset-y
        >
          <template v-slot:activator="{ on }">
            <v-btn
              text
              :disabled="!!inProgress.length"
              v-on="on"
            >
              Add to profile&nbsp;
              <v-icon
                class="icon"
              >
                mdi-chevron-down
              </v-icon>
            </v-btn>
          </template>
          <v-list>
            <v-list-item
              v-for="(conf, i) in profiles"
              :key="i"
              @click="toggleIsInProfile(conf)"
            >
              <v-list-item-title>
                <v-icon v-if="conf.items.some((item) => item.id === mod.modInfo.mod_reference)">
                  mdi-check
                </v-icon><span v-else>&nbsp;&nbsp;&nbsp;&nbsp;</span>&nbsp;{{ conf.name }}
              </v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-col>-->
      <v-col
        cols="auto"
        align="center"
      >
        <v-row
          align="center"
        >
          <v-col
            cols="auto"
            align="center"
          >
            <span>Favorite</span>
          </v-col>
          <v-col
            cols="auto"
            align="center"
          >
            <v-switch
              v-model="isFavorite"
              inset
              dense
              color="primary"
              class="custom"
              flat
              @click.stop.prevent="favoriteClicked"
            />
          </v-col>
        </v-row>
      </v-col>
    </v-row>
    <v-row
      v-if="!expandDetails"
      class="image-container"
      style="flex-grow: 2"
      no-gutters
    >
      <ModImageContainer
        v-if="images && images.length > 0"
        :images="images"
      />
    </v-row>
  </v-card>
</template>

<script>
import { mapState, mapGetters } from 'vuex';
import { eq, coerce, valid } from 'semver';
import { markdownAsElement } from '../utils';
import ModImageContainer from './ModImageContainer';

export default {
  components: { ModImageContainer },
  data() {
    return {
      expandDetails: false,
      images: [],
      windowHeight: 0,
    };
  },
  computed: {
    ...mapState([
      'profiles',
      'inProgress',
      'expandedModId',
    ]),
    ...mapGetters([
      'canInstallMods',
    ]),
    mod() {
      return this.$store.state.mods.find((mod) => mod.modInfo.mod_reference === this.expandedModId);
    },
    isFavorite() {
      return this.$store.state.favoriteModIds.includes(this.expandedModId);
    },
    icon() {
      return this.mod.modInfo.logo || 'https://ficsit.app/static/assets/images/no_image.png';
    },
    authors() {
      return this.mod.modInfo.authors.map((author) => author.user.username).join(', ');
    },
    descriptionAsElement() {
      return markdownAsElement(this.mod.modInfo.full_description || '');
    },
    modDescription() {
      const el = this.descriptionAsElement;
      const links = el.getElementsByTagName('a');
      for (let i = 0; i < links.length; i += 1) {
        links[i].target = '_blank';
      }
      return el.innerHTML;
    },
  },
  watch: {
    expandedModId() {
      this.generateImages();
    },
    images() {
      this.expandDetails = !this.images || this.images.length === 0;
    },
  },
  created() {
    window.addEventListener('resize', this.onResize);
  },
  destroyed() {
    window.removeEventListener('resize', this.onResize);
  },
  mounted() {
    this.generateImages();
    this.windowHeight = window.innerHeight;
  },
  methods: {
    close() {
      this.$store.dispatch('unexpandMod');
    },
    toggleExpandDetails() {
      this.expandDetails = !this.expandDetails;
    },
    switchClicked(mod) {
      this.$store.dispatch('switchModInstalled', mod.mod_reference);
    },
    favoriteClicked() {
      this.$store.dispatch('toggleModFavorite', this.mod.modInfo.mod_reference);
    },
    installVersion(version) {
      this.$store.dispatch('installModVersion', { modId: this.mod.modInfo.mod_reference, version });
    },
    toggleIsInProfile(profile) {
      if (profile.items.some((item) => item.id === this.mod.modInfo.mod_reference)) {
        this.$store.dispatch('addModToProfile', { mod: this.$store.state.expandedModId, profile: profile.name });
      } else {
        this.$store.dispatch('removeModFromProfile', { mod: this.$store.state.expandedModId, profile: profile.name });
      }
    },
    async generateImages() {
      this.images = []; // gallery is disabled until width issue is fixed
      // const el = this.descriptionAsElement;
      // const imgs = [...el.getElementsByTagName('img')];
      // await Promise.all(imgs.map(async (img) => {
      //   while (!img.complete) {
      //     // eslint-disable-next-line no-await-in-loop
      //     await new Promise((resolve) => setTimeout(resolve, 100));
      //   }
      // }));
      // imgs.sort((a, b) => a.naturalWidth - b.naturalWidth);
      // this.images = imgs.map((img) => img.src);
    },
    validAndEq(v1, v2) {
      const v1Valid = valid(coerce(v1));
      const v2Valid = valid(coerce(v2));
      if (v1Valid && v2Valid) {
        return eq(v1Valid, v2Valid);
      }
      return false;
    },
    onResize() {
      this.windowHeight = window.innerHeight;
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
.v-btn {
  color: var(--v-text-base) !important;
  font-size: 16px !important;
}
.mod-icon {
  width: 256px;
  height: 256px;
  margin: 0 32px 32px 32px;
}
.header {
  font-weight: 1000;
}
.authors {
  color: var(--v-info-base);
  font-weight: 500;
}

.mod-description {
  display: block;
  overflow: hidden;
  height: 169px;
  word-break: break-word;
}
.mod-description.expanded {
  overflow-y: auto;
  box-shadow: none;
}

.expand-details-button {
  position: relative;
  top: -5px;
}
.expand-details-button.expanded {
  top: 5px;
}

.v-divider {
  border-color: var(--v-text-base) !important;
}

.custom.v-input {
  margin-top: 0 !important;
  color: var(--v-background-base);
}

.control-bar {
  background: rgba(0, 0, 0, 0.45);
}
.control-bar .col {
  padding-top: 2px !important;
  padding-bottom: 2px !important;
}

.row {
  margin-left: 0 !important;
  margin-right: 0 !important;
}

.expand-details-text {
  text-align: center;
  line-height: 10px;
}
:not(.expanded) > .expand-details-text {
  padding-bottom: 0;
  padding-top: 17px;
}
.expended-description-shadow {
  position: absolute;
  width: 100%;
  height: 169px;
  box-shadow: inset 0px -35px 13px -16px rgba(0, 0, 0, 0.45);
}
</style>
