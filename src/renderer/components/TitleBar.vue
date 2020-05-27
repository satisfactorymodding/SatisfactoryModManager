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

            <v-divider class="custom" />

            <v-list-item>
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Expand mod info</v-list-item-title>
              </v-list-item-content>

              <v-list-item-action>
                <v-switch v-model="expandModInfoOnStart" />
              </v-list-item-action>
            </v-list-item>

            <v-divider
              class="custom"
              inset
            />

            <v-list-item>
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Dark mode</v-list-item-title>
              </v-list-item-content>

              <v-list-item-action>
                <v-switch v-model="darkMode" />
              </v-list-item-action>
            </v-list-item>

            <v-divider
              class="custom"
              inset
            />

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

            <v-divider class="custom" />

            <v-list-item @click.stop="creditsDialog = true">
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Credits</v-list-item-title>
              </v-list-item-content>

              <v-list-item-action>
                <v-icon color="text">
                  mdi-information
                </v-icon>
              </v-list-item-action>
              <v-dialog v-model="creditsDialog">
                <v-card>
                  <v-card-title>
                    Credits
                  </v-card-title>
                  <v-card-text>
                    Mircea Roata - programming<br>
                    Deantendo - UI &amp; icons<br>
                    Vilsol - <a
                      href="https://ficsit.app"
                      target="_blank"
                    >ficsit.app</a>
                  </v-card-text>
                  <v-card-actions>
                    <v-btn
                      color="primary"
                      text
                      @click="creditsDialog = false"
                    >
                      Close
                    </v-btn>
                  </v-card-actions>
                </v-card>
              </v-dialog>
            </v-list-item>

            <v-divider
              class="custom"
              inset
            />

            <v-list-item @click.stop="attributionDialog = true">
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Attribution</v-list-item-title>
              </v-list-item-content>

              <v-list-item-action>
                <v-icon color="text">
                  mdi-information
                </v-icon>
              </v-list-item-action>
              <v-dialog v-model="attributionDialog">
                <v-card>
                  <v-card-title>
                    Attribution
                  </v-card-title>
                  <v-card-text>
                    This app uses:<br>
                    <a
                      href="https://github.com/electron/electron"
                      target="_blank"
                    >Electron</a><br>
                    <a
                      href="https://github.com/vuejs/vue"
                      target="_blank"
                    >Vue</a><br>
                    <a
                      href="https://github.com/vuetifyjs/vuetify"
                      target="_blank"
                    >Vuetify</a><br>
                    Licensed under MIT license<br>
                    <a
                      href="https://materialdesignicons.com/"
                      target="_blank"
                    >Material Design Icons</a><br>
                    And many others
                  </v-card-text>
                  <v-card-actions>
                    <v-btn
                      color="primary"
                      text
                      @click="attributionDialog = false"
                    >
                      Close
                    </v-btn>
                  </v-card-actions>
                </v-card>
              </v-dialog>
            </v-list-item>

            <v-divider
              class="custom"
              inset
            />

            <v-list-item @click.stop="helpDialog = true">
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Help</v-list-item-title>
              </v-list-item-content>

              <v-list-item-action>
                <v-icon color="text">
                  mdi-information
                </v-icon>
              </v-list-item-action>
              <v-dialog v-model="helpDialog">
                <v-card>
                  <v-card-title>
                    Attribution
                  </v-card-title>
                  <v-card-text>
                    <h3>General troubleshooting</h3><br>
                    If something doesn't behave as expected, the first thing to try is <a @click="clearCache">clearing the cache</a><br>
                    If that doesn't work, <a @click="enableDebug">enable debug mode</a>, then click on the "Console" tab of the Developer Tool Panel.<br>
                    Recreate what you tried to do and went wrong, then send screenshot of the console<br>
                    <!-- TODO: generate debug info button -->
                    <v-divider />
                    <h3>Why does SML not show mods?</h3><br>
                    Fist, check that Epic can start the game. If you changed where the game is located, you need to make Epic update its install info. To do so:<br>
                    1. Rename the game folder to something temporary<br>
                    2. Start install from Epic to the directory you want the game to be in (the original folder name, before step 1)<br>
                    3. After it downloads a bit close Epic<br>
                    4. Copy back the files from the temporary folder EXCEPT the .egstore folder<br>
                    5. Start Epic and resume the install so it finds that it is actually already installed<br>

                    If you are running a cracked version of the game, ping @Moderator in the <a @click="moddingDiscord">modding discord</a> for help.
                  </v-card-text>
                  <v-card-actions>
                    <v-btn
                      color="primary"
                      text
                      @click="helpDialog = false"
                    >
                      Close
                    </v-btn>
                  </v-card-actions>
                </v-card>
              </v-dialog>
            </v-list-item>

            <v-divider
              class="custom"
              inset
            />

            <v-list-item>
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Debug mode</v-list-item-title>
              </v-list-item-content>

              <v-list-item-action>
                <v-switch v-model="debugMode" />
              </v-list-item-action>
            </v-list-item>

            <v-divider
              class="custom"
              inset
            />

            <v-list-item
              @click="clearCache"
            >
              <v-list-item-action />
              <v-list-item-content>
                <v-list-item-title>Clear cache</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-divider
              class="custom"
              inset
            />

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

            <v-divider class="custom" />

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

            <v-divider
              class="custom"
              inset
            />

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

            <v-divider class="custom" />

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
    <v-dialog v-model="modUpdatesDialog">
      <v-card>
        <v-card-title>
          Mod updates available
        </v-card-title>
        <v-card-text>
          TODO
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="primary"
            text
            @click="modUpdatesDialog = false"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="smmUpdateDialog">
      <v-card>
        <v-card-title>
          SMM update available: {{ availableUpdate ? availableUpdate.version : '' }}
        </v-card-title>
        <v-card-text>
          <!-- eslint-disable-next-line vue/no-v-html -->
          <div v-html="availableUpdate ? availableUpdate.releaseNotes.substr(availableUpdate.releaseNotes.indexOf('<h2>Changelog</h2>')) : ''" />
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="primary"
            text
            @click="installNow"
          >
            Update now
          </v-btn>
          <v-btn
            color="primary"
            text
            @click="installAtExit"
          >
            Update at exit
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { clearCache } from 'satisfactory-mod-manager-api';
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
      creditsDialog: false,
      attributionDialog: false,
      helpDialog: false,
      availableUpdate: null,
      modUpdatesDialog: false,
      smmUpdateDialog: false,
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
    debugMode: {
      get() {
        return this.$store.state.debugMode;
      },
      set(value) {
        this.$store.dispatch('setDebugMode', value);
        if (value) {
          this.$electron.ipcRenderer.send('openDevTools');
        }
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
  mounted() {
    this.$electron.ipcRenderer.on('updateAvailable', (e, updateInfo) => {
      this.availableUpdate = updateInfo;
      this.smmUpdateDialog = true;
    });
    setInterval(() => {
      this.$electron.ipcRenderer.send('checkForUpdates');
    }, 5 * 60 * 1000);
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
    enableDebug() {
      this.debugMode = true;
    },
    clearCache() {
      clearCache();
      if (this.$store.state.selectedInstall) {
        this.$store.state.selectedInstall.clearCache();
      }
    },
    installNow() {
      // TODO
    },
    installAtExit() {
      // TODO
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
.custom.v-divider--inset:not(.v-divider--vertical) {
  margin-left: 30px !important;
  max-width: calc(100% - 30px) !important;
}
.custom.v-divider:not(.v-divider--inset):not(.v-divider--vertical) {
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
