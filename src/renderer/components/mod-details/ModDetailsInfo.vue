<template>
  <v-card
    width="296"
    max-width="296"
    min-width="296"
  >
    <v-container
      fluid
      align
      style="padding: 20px 20px 15px 20px; height: 100%"
      class="mod-stats d-flex flex-column"
    >
      <v-row style="padding-bottom: 15px">
        <img
          :src="icon"
          width="256"
          height="256"
        >
      </v-row>
      <v-row style="padding-bottom: 15px">
        <span class="header">{{ mod.modInfo.name }}</span>
      </v-row>
      <v-row style="padding-bottom: 10px">
        <span>
          A mod by:<br>
          <span
            class="primary--text"
            style="font-weight: 600"
          >{{ mod.modInfo.authors[0].user.username }}</span>
        </span>
      </v-row>
      <v-row
        style="padding-bottom: 10px"
      >
        <v-menu
          open-on-hover
          offset-y
        >
          <template v-slot:activator="{ on, attrs }">
            <v-btn
              color="backgroundModsList-darken1"
              raised
              v-bind="attrs"
              width="100%"
              v-on="on"
            >
              <span>See contributors <span class="primary--text">({{ mod.modInfo.authors.length }})</span></span>
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
              v-for="(author, index) in mod.modInfo.authors"
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
          Size: <span class="header">{{ mod.modInfo.versions[0] && mod.modInfo.versions[0].size ? bytesToAppropriate(mod.modInfo.versions[0].size) : 'N/A' }}</span><br>
          Created: <span class="header">{{ mod.modInfo.created_at.toLocaleDateString() }}</span><br>
          Updated: <span class="header">{{ mod.modInfo.last_version_date ? mod.modInfo.last_version_date.toLocaleString() : 'N/A' }}</span><br>
          Latest version: <span class="header">{{ mod.modInfo.versions[0] ? mod.modInfo.versions[0].version : 'N/A' }}</span><br>
        </span>
      </v-row>
      <v-row
        style="padding-bottom: 10px"
      >
        <v-menu
          offset-y
        >
          <template v-slot:activator="{ on, attrs }">
            <v-btn
              color="backgroundModsList-darken1"
              raised
              v-bind="attrs"
              width="100%"
              :disabled="!!inProgress.length || !mod.isInstalled"
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
          <v-list>
            <v-list-item
              @click="$emit('installVersion', '')"
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
            <v-list-item
              v-for="(version, i) in mod.modInfo.versions"
              :key="i"
              @click="$emit('installVersion', version.version)"
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
          </v-list>
        </v-menu>
      </v-row>
      <v-row style="padding-bottom: 10px">
        <span>
          <!--Multiplayer support: <span
            class="header"
            :class="multiplayerSupportColor + '--text'"
          >{{ multiplayerSupport }}</span><br>-->
          Total downloads: <span class="header">{{ mod.modInfo.downloads.toLocaleString() }}</span><br>
          Views: <span class="header">{{ mod.modInfo.views.toLocaleString() }}</span><br>
          <!--Tags:<br>
          <span class="header">{{ mod.modInfo.tags ? mod.modInfo.tags.map((tag) => `#${tag}`).join(' ') : 'N/A' }}</span><br>-->
        </span>
      </v-row>
      <v-spacer style="flex-basis: 100%" />
      <v-row>
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
import { bytesToAppropriate } from '@/utils';
import { eq, coerce, valid } from 'semver';
import { mapState, mapGetters } from 'vuex';

export default {
  props: {
    mod: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    ...mapState([
      'inProgress',
    ]),
    ...mapGetters([
      'canInstallMods',
    ]),
    icon() {
      return this.mod.modInfo.logo || 'https://ficsit.app/static/assets/images/no_image.png';
    },
    multiplayerSupport() {
      return this.mod.modInfo.multiplayer_support || 'N/A';
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
  },
  methods: {
    validAndEq(v1, v2) {
      const v1Valid = valid(coerce(v1));
      const v2Valid = valid(coerce(v2));
      if (v1Valid && v2Valid) {
        return eq(v1Valid, v2Valid);
      }
      return false;
    },
    bytesToAppropriate,
  },
};
</script>

<style scoped>
.mod-icon {
  width: 256px;
  height: 256px;
  margin: 0 32px 32px 32px;
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
</style>
