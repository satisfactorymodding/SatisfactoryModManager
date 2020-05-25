<template>
  <div class="titlebar">
    <div
      class="d-inline-flex align-items-center"
      @click="settingsClicked"
    >
      <v-menu
        v-model="menuOpen"
        :close-on-content-click="false"
        offset-x
      >
        <template v-slot:activator="{ on }">
          <v-icon
            :color="getColorForState"
            class="ma-1 app-icon"
            v-on="on"
          >
            mdi-cog
          </v-icon>
        </template>
        <v-card class="app-menu">
          <v-list>
            <v-list-item
              v-if="hasUpdate"
              @click="installUpdates"
            >
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Updates ready to install</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-list-item>
              <v-list-item-action>
                <v-icon color="text">
                  mdi-cog
                </v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Settings</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-divider />

            <v-list-item>
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Expand mod info</v-list-item-title>
              </v-list-item-content>

              <v-list-item-action>
                <v-switch v-model="expandModInfoOnStart" />
              </v-list-item-action>
            </v-list-item>

            <v-divider inset />

            <v-list-item>
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Dark mode</v-list-item-title>
              </v-list-item-content>

              <v-list-item-action>
                <v-switch v-model="darkMode" />
              </v-list-item-action>
            </v-list-item>

            <v-divider inset />

            <v-list-item>
              <v-list-item-action>
                <v-icon color="text">
                  mdi-information
                </v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>About</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-divider />

            <v-list-item>
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Credits</v-list-item-title>
              </v-list-item-content>

              <v-list-item-action>
                <v-icon color="text">
                  mdi-information
                </v-icon>
              </v-list-item-action>
            </v-list-item>

            <v-divider inset />

            <v-list-item>
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Attribution</v-list-item-title>
              </v-list-item-content>

              <v-list-item-action>
                <v-icon color="text">
                  mdi-information
                </v-icon>
              </v-list-item-action>
            </v-list-item>

            <v-divider inset />

            <v-list-item>
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Help</v-list-item-title>
              </v-list-item-content>

              <v-list-item-action>
                <v-icon color="text">
                  mdi-information
                </v-icon>
              </v-list-item-action>
            </v-list-item>

            <v-divider inset />

            <v-list-item>
              <v-list-item-action>
                <v-icon color="text">
                  mdi-discord
                </v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Discord</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-divider />

            <v-list-item @click="moddingDiscord">
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Satisfactory Modding</v-list-item-title>
              </v-list-item-content>

              <v-list-item-action>
                <v-icon color="text">
                  mdi-open-in-new
                </v-icon>
              </v-list-item-action>
            </v-list-item>

            <v-divider inset />

            <v-list-item @click="officialDiscord">
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Official Satisfactory</v-list-item-title>
              </v-list-item-content>

              <v-list-item-action>
                <v-icon color="text">
                  mdi-open-in-new
                </v-icon>
              </v-list-item-action>
            </v-list-item>

            <v-divider />

            <v-list-item>
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Satisfactory Mod Manager v{{ version }}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-card>
      </v-menu>
    </div>
    <div class="bar">
      <div class="dragregion">
        <span class="app-title">{{ title }}</span>
      </div>
      <div
        class="button close"
        @click="onClose"
      >
        <span>&#10005;</span>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import { saveSetting } from '../settings';

export default {
  props: {
    title: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      menuOpen: false,
    };
  },
  computed: {
    ...mapState([
      'hasUpdate',
    ]),
    state() {
      if (this.menuOpen) {
        return 'on';
      }
      if (this.hasUpdate) {
        return 'notify';
      }
      return 'off';
    },
    getColorForState() {
      if (this.state === 'notify') {
        return '#ffc107';
      }
      if (this.state === 'on') {
        return 'primary';
      }
      return '#9e9e9e';
    },
    darkMode: {
      get() {
        return this.$vuetify.theme.dark;
      },
      set(value) {
        this.$vuetify.theme.dark = value;
        saveSetting('darkMode', value);
      },
    },
    expandModInfoOnStart: {
      get() {
        return this.$store.state.expandModInfoOnStart;
      },
      set(value) {
        this.$store.dispatch('setExpandModInfoOnStart', value);
      },
    },
    version() {
      return this.$electron.remote.app.getVersion();
    },
  },
  methods: {
    onClose() {
      this.$electron.remote.getCurrentWindow().close();
    },
    settingsClicked() {
      this.menuOpen = !this.menuOpen;
    },
    installUpdates() {
      console.log('UPDATE');
    },
    moddingDiscord() {
      this.$electron.shell.openExternal('https://discord.gg/TShj39G');
    },
    officialDiscord() {
      this.$electron.shell.openExternal('https://discord.gg/Satisfactory');
    },
  },
};
</script>

<style scoped>
.app-title {
  font-size: 16px !important;
}
.titlebar {
  display: flex;
  height: var(--titlebar-height);
}
.titlebar, .titlebar > * {
  color: var(--v-text2-base) !important;
  background-color: var(--v-background-base);
}
.app-menu .v-list {
  background-color: var(--v-menuBackground-base);
}
.v-icon {
  font-size: 16px !important;
}
.app-icon {
  padding: 3px 0px 0px 3px;
}
.app-icon.v-icon.v-icon:after {
    background-color: rgba(0,0,0,0);
}
.v-list-item {
  padding-left: 10px !important;
}
.v-list-item__action:first-child {
  margin-right: 0px !important;
}
.v-divider--inset:not(.v-divider--vertical) {
  margin-left: 30px !important;
  max-width: calc(100% - 30px) !important;
}
.v-divider:not(.v-divider--inset):not(.v-divider--vertical) {
  margin-left: 10px !important;
  max-width: calc(100% - 40px) !important;
}
.bar {
  flex-grow: 1;
  display: flex;
}
.dragregion {
  flex-grow: 1;
  text-align: center;
  vertical-align: middle;
  margin: 3px 3px 0px 0px;
  -webkit-app-region: drag;
  display: flex;
  align-items: center;
}
.dragregion>span {
  flex-grow: 1;
  margin-top: -3px;
}
.button {
  width: 44px;
  text-align: center;
  font-size: 12pt;
  display: flex;
  align-items: center;
  cursor: pointer;
  color: var(--v-text-base);
}
.button>span {
  flex-grow: 1;
  margin-top: -3px;
  user-select: none;
}
.button:hover {
  background-color: var(--titlebar-button-hover-color);
}
.close:hover {
  background-color: red;
  color: white;
}
.button>span.dash {
  vertical-align: sub;
  margin-top: 0px;
}
#settingsIcon {
  height: 18px;
  margin: 4.5px;
}
</style>
