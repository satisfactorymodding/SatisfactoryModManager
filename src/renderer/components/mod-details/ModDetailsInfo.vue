<template>
  <v-card
    class="mod-details-card"
  >
    <v-container
      fluid
      align
      style="padding: 20px 20px 15px 20px; height: 100%"
      class="mod-stats d-flex flex-column"
    >
      <v-row style="padding-bottom: 15px">
        <img
          class="mod-icon"
          :src="icon"
        >
      </v-row>
      <v-row style="padding-bottom: 15px">
        <span class="header">{{ mod.name }}</span>
      </v-row>
      <v-row style="padding-bottom: 10px">
        <span>
          A mod by:<br>
          <span
            class="primary--text"
            style="font-weight: 600"
            @click="searchByAuthor"
          >{{ mod.authors[0].user.username }}</span>
        </span>
      </v-row>
      <v-row
        style="padding-bottom: 10px"
      >
        <v-menu
          open-on-hover
          offset-y
        >
          <template #activator="{ on, attrs }">
            <v-btn
              color="backgroundModsList-darken1"
              raised
              v-bind="attrs"
              width="100%"
              v-on="on"
            >
              <span>See contributors <span class="primary--text">({{ mod.authors.length }})</span></span>
              <v-spacer />
              <v-icon
                right
                dark
              >
                mdi-chevron-down
              </v-icon>
            </v-btn>
          </template>
          <v-list>
            <v-list-item
              v-for="(author, index) in mod.authors"
              :key="index"
            >
              <v-list-item-avatar>
                <v-img :src="author.user.avatar" />
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title>{{ author.user.username }}</v-list-item-title>
                <v-list-item-subtitle>{{ author.role }}</v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-row>
      <v-row style="padding-bottom: 10px">
        <span>
          Mod info:<br>
          Size: <span class="header">{{ mod.versions[0] && mod.versions[0].size ? bytesToAppropriate(mod.versions[0].size) : 'N/A' }}</span><br>
          Created: <span class="header">{{ mod.created_at.toLocaleDateString() }}</span><br>
          Updated: <span class="header">{{ mod.last_version_date ? mod.last_version_date.toLocaleString() : 'N/A' }}</span><br>
          Latest version: <span class="header">{{ mod.versions[0] ? mod.versions[0].version : 'N/A' }}</span><br>
          Installed version: <span class="header">{{ installedVersion ? installedVersion : 'N/A' }}</span><br>
        </span>
      </v-row>
      <v-row
        style="padding-bottom: 10px"
      >
        <v-menu
          offset-y
        >
          <template #activator="{ on, attrs }">
            <v-btn
              color="backgroundModsList-darken1"
              raised
              v-bind="attrs"
              width="100%"
              :disabled="!!inProgress.length || !isInstalled"
              v-on="on"
            >
              <span>Mod version options</span>
              <v-spacer />
              <v-icon
                right
                dark
              >
                mdi-chevron-down
              </v-icon>
            </v-btn>
          </template>
          <v-list
            style="max-height: 90vh"
            class="overflow-y-auto"
          >
            <v-list-item
              @click="$emit('install-version', '')"
            >
              <v-list-item-action>
                <v-icon v-if="!manifestVersion">
                  mdi-check
                </v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Latest</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item
              v-for="(version, i) in mod.versions"
              :key="i"
              @click="$emit('install-version', version.version)"
            >
              <v-list-item-action>
                <v-icon v-if="validAndEq(manifestVersion, version.version)">
                  mdi-check
                </v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>{{ version.version }}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-row>
      <v-row style="padding-bottom: 10px">
        <span>
          <!--Multiplayer support: <span
            class="header"
            :class="multiplayerSupportColor + '--text'"
          >{{ multiplayerSupport }}</span><br>-->
          Total downloads: <span class="header">{{ mod.downloads.toLocaleString() }}</span><br>
          Views: <span class="header">{{ mod.views.toLocaleString() }}</span><br>
          <!--Tags:<br>
          <span class="header">{{ mod.tags ? mod.tags.map((tag) => `#${tag}`).join(' ') : 'N/A' }}</span><br>-->
        </span>
      </v-row>
      <v-row>
        <a
          class="primary--text"
          style="font-weight: 600"
          target="_blank"
          :href="`https://ficsit.app/mod/${mod.id}`"
        >View on ficsit.app</a>
      </v-row>
      <v-row style="flex-basis: 100%" />
      <v-row
        class="pb-2"
      >
        <v-btn
          color="backgroundModsList-darken1"
          raised
          width="100%"
          @click="$emit('close')"
        >
          <v-icon
            left
            dark
          >
            mdi-flip-h mdi-import
          </v-icon>
          <span>Close</span>
        </v-btn>
      </v-row>
    </v-container>
  </v-card>
</template>

<script>
import { mapState } from 'vuex';
import { validAndEq, bytesToAppropriate } from '@/utils';

export default {
  props: {
    mod: {
      type: Object,
      default: () => ({}),
    },
    isInstalled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    ...mapState([
      'inProgress',
    ]),
    icon() {
      return this.mod.logo || 'https://ficsit.app/images/no_image.webp';
    },
    multiplayerSupport() {
      return this.mod.multiplayer_support || 'N/A';
    },
    multiplayerSupportColor() {
      switch (this.multiplayerSupport) {
        case 'Poor':
          return 'red';
        case 'Medium':
          return 'yellow';
        case 'Good':
          return 'green';
        default:
          return 'red';
      }
    },
    manifestVersion() {
      return this.$store.state.manifestItems.find((item) => item.id === this.mod.mod_reference)?.version;
    },
    installedVersion() {
      return this.$store.state.selectedInstall?.mods[this.mod.mod_reference];
    },
  },
  methods: {
    searchByAuthor() {
      this.$root.$emit('updateSearch', `author:"${this.mod.authors[0].user.username}"`);
    },
    validAndEq,
    bytesToAppropriate,
  },
};
</script>

<style scoped>
.mod-icon {
  width: 256px;
  height: 256px;
}
.mod-stats * {
  font-size: 14px;
  line-height: 1.25;
}
.mod-stats .header {
  color: var(--v-text-lighten2);
  font-weight: 600;
}
.mod-stats .v-btn {
  color: var(--v-text-base) !important;
  font-size: 16px !important;
  padding-left: 10px !important;
  padding-right: 15px !important;
}
.row {
  margin-left: 0 !important;
  margin-right: 0 !important;
}
.mod-icon {
  height: 256px;
  width: 256px;
}
.mod-details-card {
  width: 296px;
  min-width: 296px;
  max-width: 296px;
}
@media (max-height: 700px) {
  .mod-icon {
    height: 200px;
    width: 200px;
    margin-left: 10px;
  }
  .mod-details-card {
    width: 260px;
    min-width: 260px;
    max-width: 260px;
  }
}
</style>
