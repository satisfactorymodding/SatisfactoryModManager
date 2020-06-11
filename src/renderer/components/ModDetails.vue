<!-- eslint-disable vue/no-v-html -->
<template>
  <v-container
    style="width: 100%; height: 1075px; padding: 0; box-shadow: inset 10px 0px 10px -10px rgba(0,0,0,1);"
    fluid
  >
    <v-row
      no-gutters
      style="padding-top: 32px;"
    >
      <v-col cols="auto">
        <img
          :src="icon"
          class="mod-icon"
        >
      </v-col>
      <v-col>
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
          <v-col cols="4">
            <span class="header">{{ mod.modInfo.name }}</span>
          </v-col>
        </v-row>
        <v-row
          class="mod-description"
          :class="expandDetails ? 'expanded' : ''"
          v-html="modDescription"
        />
        <v-row
          class="expand-details-button"
          :class="expandDetails ? 'expanded' : ''"
          @click="toggleExpandDetails"
        >
          <v-col
            style="text-align: center; line-height: 10px"
          >
            <span>Show {{ expandDetails ? 'less' : 'more' }}<br><v-icon v-if="!expandDetails">mdi-chevron-down</v-icon><v-icon v-else>mdi-chevron-up</v-icon></span>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
    <v-row
      class="control-bar"
      align="center"
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
      <!--<v-col
        cols="auto"
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
              <span
                class="mx-2"
              >Previous versions</span>
              <v-icon>
                mdi-chevron-down
              </v-icon>
            </v-btn>
          </template>
          <v-list>
            <v-list-item>
              <v-list-item-title>
                <v-icon v-if="installedVersion === ''">
                  mdi-check
                </v-icon><span v-else>&nbsp;&nbsp;&nbsp;&nbsp;</span>&nbsp;Latest
              </v-list-item-title>
            </v-list-item>
            <v-list-item
              v-for="(version, i) in mod.modInfo.versions"
              :key="i"
              @click="installOldVersion(conf)"
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
      no-gutters
    >
      <div
        :class="imagePage > 0 ? '' : 'hidden'"
        class="images-button left d-inline-flex align-center"
        @click="imagePageLeft"
      >
        <v-icon>mdi-chevron-left</v-icon>
      </div>
      <v-row
        class="scrollable-images"
        :style="`left: ${-432 * imagePage}px; width: ${432 * Math.ceil(images.length / 2)}px`"
      >
        <template v-for="n in images.length">
          <img
            v-if="images[n - 1]"
            :key="n"
            :src="images[n - 1]"
            @click="bigImage(n - 1)"
          >
        </template>
      </v-row>
      <div
        :class="imagePage < Math.ceil(images.length / 2) - 2 ? '' : 'hidden'"
        class="images-button right d-inline-flex align-center"
        style="left: 1056px;"
        @click="imagePageRight"
      >
        <v-icon>mdi-chevron-right</v-icon>
      </div>
    </v-row>
    <v-dialog
      v-model="showBigImage"
      width="unset"
    >
      <v-card>
        <img
          :src="bigImageSrc"
          style="display: block; max-height: 750px"
        >
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
import { mapState, mapGetters } from 'vuex';
import { markdownAsElement } from '../utils';

export default {
  data() {
    return {
      expandDetails: false,
      imagePage: 0,
      bigImageSrc: '',
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
    images() {
      // TODO: other image sources
      const el = this.descriptionAsElement;
      const imgs = el.getElementsByTagName('img');
      const imgUrls = [];
      for (let i = 0; i < imgs.length; i += 1) {
        if (imgs[i].src) {
          imgUrls.push(imgs[i].src);
        }
      }
      return imgUrls;
    },
    showBigImage: {
      get() {
        return !!this.bigImageSrc;
      },
      set(value) {
        if (!value) {
          this.bigImageSrc = '';
        }
      },
    },
  },
  watch: {
    expandedModId() {
      this.imagePage = 0;
    },
  },
  methods: {
    close() {
      this.$store.dispatch('unexpandMod');
    },
    toggleExpandDetails() {
      this.expandDetails = !this.expandDetails;
    },
    imagePageLeft() {
      if (this.imagePage > 0) this.imagePage -= 1;
    },
    imagePageRight() {
      if (this.imagePage < Math.ceil(this.images.length / 2) - 2) this.imagePage += 1;
    },
    switchClicked(mod) {
      this.$store.dispatch('switchModInstalled', mod.mod_reference);
    },
    favoriteClicked() {
      this.$store.dispatch('toggleModFavorite', this.mod.modInfo.mod_reference);
    },
    toggleIsInProfile(profile) {
      if (profile.items.some((item) => item.id === this.mod.modInfo.mod_reference)) {
        this.$store.dispatch('addModToProfile', { mod: this.$store.state.expandedModId, profile: profile.name });
      } else {
        this.$store.dispatch('removeModFromProfile', { mod: this.$store.state.expandedModId, profile: profile.name });
      }
    },
    bigImage(idx) {
      this.bigImageSrc = this.images[idx];
    },
  },
};
</script>

<style>
.mod-description img {
  max-width: 600px;
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
  box-shadow: inset 0px -35px 13px -16px rgba(0, 0, 0, 0.45);
}

.mod-description.expanded {
  overflow-y: auto;
  height: 655px;
  box-shadow: none;
}
.control-bar {
  background: rgba(0, 0, 0, 0.45);
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
.control-bar .col {
  padding-top: 2px !important;
  padding-bottom: 2px !important;
}
.row {
  margin-left: 0 !important;
  margin-right: 0 !important;
}
.image-container {
  overflow-x: hidden;
  width: 100%;
}
.image-container img {
  width: 432px;
  height: 243px;
  display: block;
}
.images-button {
  text-align: center;
  line-height: 10px;
  position: absolute;
  bottom: 217px;
  height: 486px;
  z-index: 1;
}
.images-button.left {
  box-shadow: inset 50px 0px 50px -50px rgba(0,0,0,1);
}
.images-button.right {
  box-shadow: inset -50px 0px 50px -50px rgba(0,0,0,1);
}
.images-button.hidden {
  visibility: hidden;
}
.scrollable-images {
  width: 100%;
  position: relative;
  transition: all ease-in-out 0.5s;
  position: absolute;
  left: 0;
  right: 0;
}
</style>
