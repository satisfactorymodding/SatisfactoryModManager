<!-- eslint-disable vue/no-v-html -->
<template>
  <v-container
    style="width: 100%; height: 800px; padding: 0; box-shadow: inset 10px 0px 10px -10px rgba(0,0,0,1);"
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
        style="padding-right: 585px"
        @click="close"
      >
        <v-icon>
          mdi-flip-h mdi-import
        </v-icon>
        <span
          class="mx-2"
        >Close</span>
      </v-col>
      <!--<v-col
        cols="auto"
      >
        <span
          class="mx-2"
        >Previous versions</span>
        <v-icon>
          mdi-chevron-down
        </v-icon>
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
              :disabled="!mod.isCompatible || !!inProgress.id"
              @click.stop.prevent="switchClicked(mod.modInfo)"
            />
          </v-col>
        </v-row>
      </v-col>
      <v-col
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
              :disabled="!!inProgress.id"
              v-on="on"
            >
              Add to config&nbsp;
              <v-icon
                class="icon"
              >
                mdi-chevron-down
              </v-icon>
            </v-btn>
          </template>
          <v-list>
            <v-list-item
              v-for="(conf, i) in configs"
              :key="i"
              @click="toggleIsInConfig(conf)"
            >
              <v-list-item-title>
                <v-icon v-if="conf.items.includes(mod.modInfo.id)">
                  mdi-check
                </v-icon><span v-else>&nbsp;&nbsp;&nbsp;&nbsp;</span>&nbsp;{{ conf.name }}
              </v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-col>
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
              :value="isFavorite"
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
      <template v-for="n in Math.min(images.length, 6)">
        <v-col
          :key="n"
          cols="4"
        >
          <img :src="images[n - 1]">
        </v-col>
      </template>
    </v-row>
  </v-container>
</template>

<script>
import marked from 'marked';
import sanitizeHtml from 'sanitize-html';

export default {
  props: {
    mod: {
      type: Object,
      default() { return {}; },
    },
    isFavorite: {
      type: Boolean,
      default: false,
    },
    inProgress: {
      type: Object,
      default() { return {}; },
    },
    configs: {
      type: Array,
      default() { return []; },
    },
  },
  data() {
    return {
      expandDetails: false,
    };
  },
  computed: {
    icon() {
      return this.mod.modInfo.logo || 'https://ficsit.app/static/assets/images/no_image.png';
    },
    authors() {
      return this.mod.modInfo.authors.map((author) => author.user.username).join(', ');
    },
    descriptionAsElement() {
      const html = sanitizeHtml(marked(this.mod.modInfo.full_description || ''), {
        allowedTags: sanitizeHtml.defaults.allowedTags.concat(['img', 'video', 'details', 'summary', 'h1', 'h2']),
        allowedAttributes: Object.assign(sanitizeHtml.defaults.allowedAttributes, { img: ['src', 'width', 'height'], video: ['src', 'width', 'height', 'controls'] }),
      });
      const el = document.createElement('html');
      el.innerHTML = html;
      return el;
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
  },
  methods: {
    close() {
      this.$emit('close');
    },
    toggleExpandDetails() {
      this.expandDetails = !this.expandDetails;
    },
    switchClicked(mod) {
      this.$emit('switchMod', mod.id);
    },
    favoriteClicked() {
      if (!this.isFavorite) {
        this.$emit('favoriteMod', this.mod.modInfo.id);
      } else {
        this.$emit('unfavoriteMod', this.mod.modInfo.id);
      }
    },
    toggleIsInConfig(config) {
      if (config.items.includes(this.mod.modInfo.id)) {
        this.$emit('addToConfig', config.name);
      } else {
        this.$emit('removeFromConfig', config.name);
      }
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
  height: 165px;
  box-shadow: inset 0px -35px 13px -16px rgba(0, 0, 0, 0.45);
}

.mod-description.expanded {
  overflow-y: auto;
  height: 598px;
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
.image-container img {
  width: 384px;
  height: 216px;
  display: block;
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
</style>
