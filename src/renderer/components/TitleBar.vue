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
            <v-menu
              v-model="updatesMenuOpen"
              :close-on-content-click="false"
              offset-x
              :nudge-right="20"
            >
              <template v-slot:activator="{ on }">
                <v-list-item
                  v-on="on"
                >
                  <v-list-item-action />
                  <v-list-item-content>
                    <v-list-item-title>{{ hasUpdate || availableSMMUpdate ? 'Updates ready to install' : 'Update settings' }}</v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </template>
              <v-card class="app-menu">
                <v-list>
                  <v-list-item>
                    <v-list-item-action>
                      <v-icon color="text">
                        mdi-cog
                      </v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                      <v-list-item-title>Update options</v-list-item-title>
                    </v-list-item-content>
                  </v-list-item>

                  <v-divider class="custom" />

                  <v-list-item>
                    <v-list-item-action />
                    <v-list-item-content>
                      <v-list-item-title>Update SMM at</v-list-item-title>
                    </v-list-item-content>
                    <v-list-item-action>
                      <v-select
                        v-model="updateCheckMode"
                        :items="['launch', 'exit', 'ask']"
                        style="width: 108px"
                      />
                    </v-list-item-action>
                  </v-list-item>

                  <v-divider
                    inset
                    class="custom"
                  />

                  <v-list-item>
                    <v-list-item-action />
                    <v-list-item-content>
                      <v-list-item-title>Show ignored updates</v-list-item-title>
                    </v-list-item-content>
                    <v-list-item-action>
                      <v-switch v-model="showIgnoredUpdates" />
                    </v-list-item-action>
                  </v-list-item>

                  <v-divider
                    inset
                    class="custom"
                  />

                  <v-list-item>
                    <v-list-item-action>
                      <v-icon color="text">
                        mdi-update
                      </v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                      <v-list-item-title>Updates</v-list-item-title>
                    </v-list-item-content>
                  </v-list-item>

                  <v-divider class="custom" />

                  <v-list-item
                    :disabled="filteredModUpdates.length === 0"
                    @click="modUpdatesDialog = true"
                  >
                    <v-list-item-action />
                    <v-list-item-content>
                      <v-list-item-title>Mod updates ({{ filteredModUpdates.length }})</v-list-item-title>
                    </v-list-item-content>
                  </v-list-item>

                  <v-divider
                    inset
                    class="custom"
                  />

                  <v-list-item
                    :disabled="!availableSMMUpdate"
                    @click="smmUpdateDialog = true"
                  >
                    <v-list-item-action />
                    <v-list-item-content>
                      <v-list-item-title>SMM updates ({{ smmUpdateCount }})</v-list-item-title>
                    </v-list-item-content>
                  </v-list-item>

                  <v-divider
                    inset
                    class="custom"
                  />
                </v-list>
              </v-card>
            </v-menu>

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
                    Help
                  </v-card-title>
                  <v-card-text>
                    <h3>General troubleshooting</h3>
                    If something doesn't behave as expected, the first thing to try is <a @click="clearCache">clearing the cache</a><br>
                    If that doesn't work, <a @click="enableDebug">enable debug mode</a>, then click on the "Console" tab of the Developer Tool Panel.<br>
                    Recreate what you tried to do and went wrong, then send screenshot of the console<br>
                    <!-- TODO: generate debug info button -->
                    <br><v-divider /><br>
                    <h3>Why does SML not show mods?</h3>
                    <h4>Epic</h4>
                    Fist, check that Epic can start the game. If you changed where the game is located, you need to make Epic update its install info. To do so:<br>
                    1. Rename the game folder to something temporary<br>
                    2. Start install from Epic to the directory you want the game to be in (the original folder name, before step 1)<br>
                    3. After it downloads a bit close Epic<br>
                    4. Copy back the files from the temporary folder EXCEPT the .egstore folder<br>
                    5. Start Epic and resume the install so it finds that it is actually already installed<br>
                    <br>
                    <h4>Cracked</h4>
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
        @click="close"
      >
        <span>&#10005;</span>
      </div>
    </div>
    <v-dialog
      v-model="modUpdatesDialog"
      max-width="600"
    >
      <v-card>
        <v-card-title>
          Mod updates available
        </v-card-title>
        <v-card-text>
          <v-btn
            color="primary"
            text
            :disabled="inProgress.length > 0"
            @click="updateAll"
          >
            Update all
          </v-btn>
          <v-progress-linear
            v-if="isMultiModUpdateInProgress"
            :value="Math.round(currentMultiModUpdateProgress.progress * 100)"
            :class="currentMultiModUpdateProgress.fast ? 'fast' : ''"
            color="warning"
            height="49"
            reactive
            :indeterminate="currentMultiModUpdateProgress.progress < 0"
          >
            <strong>{{ currentMultiModUpdateProgress.message }}</strong>
          </v-progress-linear>
          <v-list
            v-else
            class="custom"
            dense
          >
            <template v-for="(update, index) in filteredModUpdates">
              <div :key="index">
                <v-list-item
                  v-if="inProgress.some((prog) => prog.id === update.item)"
                >
                  <v-progress-linear
                    :value="Math.round(currentModProgress(update.item).progress * 100)"
                    :class="currentModProgress(update.item).fast ? 'fast' : ''"
                    color="warning"
                    height="49"
                    reactive
                    :indeterminate="currentModProgress(update.item).progress < 0"
                  >
                    <strong>{{ currentModProgress(update.item).message }}</strong>
                  </v-progress-linear>
                </v-list-item>
                <v-list-item v-else>
                  <v-list-item-content>
                    <v-list-item-title>{{ update.name }}</v-list-item-title>
                    <v-list-item-subtitle>current: v{{ update.currentVersion }}, available: v{{ update.version }}</v-list-item-subtitle>
                  </v-list-item-content>
                  <v-list-item-action style="margin-left: 10px">
                    <v-btn
                      color="text"
                      text
                      @click="viewChangelog(update)"
                    >
                      Changelog
                    </v-btn>
                  </v-list-item-action>
                  <v-list-item-action>
                    <v-btn
                      color="primary"
                      text
                      :disabled="inProgress.length > 0"
                      @click="updateItem(update)"
                    >
                      Update
                    </v-btn>
                  </v-list-item-action>
                  <v-list-item-action style="margin-left: 0">
                    <v-btn
                      color="text"
                      text
                      @click="isIgnored(update) ? unignoreUpdate(update) : ignoreUpdate(update)"
                    >
                      {{ isIgnored(update) ? 'Unignore' : 'Ignore' }}
                    </v-btn>
                  </v-list-item-action>
                </v-list-item>
              </div>
            </template>
          </v-list>
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
    <v-dialog v-model="changelogDialog">
      <v-card v-if="viewChangelogUpdate">
        <v-card-title>
          {{ viewChangelogUpdate.name }} v{{ viewChangelogUpdate.version }} changelog
        </v-card-title>
        <v-card-text>
          <template v-for="(release, index) in viewChangelogUpdate.releases">
            <div :key="index">
              <h3>v{{ release.version }}</h3>
              <!-- eslint-disable vue/no-v-html -->
              <span v-html="changelogHTML(release)" />
              <v-divider v-if="index != viewChangelogUpdate.releases.length - 1" />
            </div>
          </template>
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="primary"
            text
            @click="changelogDialog = false"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog
      v-model="smmUpdateDialog"
      persistent
    >
      <v-card>
        <v-card-title>
          SMM update available: {{ availableSMMUpdate ? availableSMMUpdate.version : '' }}
        </v-card-title>
        <v-card-text>
          <!-- eslint-disable-next-line vue/no-v-html -->
          <div v-html="smmUpdateNotes" />
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="primary"
            text
            @click="updateSMMNow"
          >
            Update now
          </v-btn>
          <v-btn
            color="primary"
            text
            @click="smmUpdateDialog = false"
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
import { saveSetting, getSetting } from '../settings';
import {
  markdownAsElement, ignoreUpdate, unignoreUpdate, lastElement,
} from '../utils';

const UPDATE_CHECK_INTERVAL = 5 * 60 * 1000;

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
      updatesMenuOpen: false,
      creditsDialog: false,
      attributionDialog: false,
      helpDialog: false,
      availableSMMUpdate: null,
      modUpdates: [],
      modUpdatesDialog: false,
      smmUpdateDialog: false,
      cachedUpdateCheckMode: 'launch',
      nextCheckForUpdates: -1,
      viewChangelogUpdate: null,
      changelogDialog: false,
      showIgnoredUpdates: false,
      ignoredUpdates: [],
    };
  },
  computed: {
    ...mapState([
      'inProgress',
      'selectedInstall',
    ]),
    ...mapState({
      allMods: 'mods',
    }),
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
    smmUpdateNotes() {
      if (!this.availableSMMUpdate) {
        return '';
      }
      return this.availableSMMUpdate.releaseNotes.map((release) => release.note.substr(release.note.indexOf('<h2>Changelog</h2>')).replace('<h2>Changelog</h2>', `<h2>v${release.version} changelog</h2>`)).join('\n');
    },
    smmUpdateCount() {
      if (!this.availableSMMUpdate) {
        return 0;
      }
      return this.availableSMMUpdate.releaseNotes.length;
    },
    updateCheckMode: {
      get() {
        return this.cachedUpdateCheckMode;
      },
      set(value) {
        saveSetting('updateCheckMode', value);
        this.cachedUpdateCheckMode = value;
      },
    },
    hasUpdate() {
      return !!this.availableSMMUpdate || this.filteredModUpdates.length > 0;
    },
    filteredModUpdates() {
      return (this.showIgnoredUpdates ? this.modUpdates : this.modUpdates.filter((update) => !this.isIgnored(update)));
    },
    isMultiModUpdateInProgress() {
      return this.inProgress.some((prog) => prog.id === '__updateMods__');
    },
    multiModUpdateProgress() {
      return this.inProgress.find((prog) => prog.id === '__updateMods__');
    },
    currentMultiModUpdateProgress() {
      return lastElement(this.multiModUpdateProgress.progresses);
    },
  },
  watch: {
    async selectedInstall() {
      await this.checkForUpdates();
    },
  },
  mounted() {
    this.cachedUpdateCheckMode = getSetting('updateCheckMode', 'launch');
    this.ignoredUpdates = getSetting('ignoredUpdates', []);
    if (this.updateCheckMode === 'launch') {
      this.$root.$once('doneLaunchUpdateCheck', () => {
        this.addUpdateListener();
      });
    } else {
      this.addUpdateListener();
    }
    this.nextCheckForUpdates = setTimeout(() => this.checkForUpdates(), UPDATE_CHECK_INTERVAL);
  },
  methods: {
    close() {
      this.$electron.remote.getCurrentWindow().close();
    },
    settingsClicked() {
      this.menuOpen = !this.menuOpen;
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
    addUpdateListener() {
      this.$electron.ipcRenderer.on('updateAvailable', (e, updateInfo) => {
        this.availableSMMUpdate = updateInfo;
        if (this.updateCheckMode === 'ask') {
          this.smmUpdateDialog = true;
        }
      });
    },
    async checkForUpdates() {
      clearTimeout(this.nextCheckForUpdates);
      // don't check for updates while something is in progress
      while (this.inProgress.length > 0) {
        // eslint-disable-next-line no-await-in-loop
        await new Promise((resolve) => setTimeout(() => resolve(), 500));
      }
      this.$electron.ipcRenderer.send('checkForUpdates');
      this.modUpdates = (await this.$store.state.selectedInstall.checkForUpdates()).map((update) => Object.assign(update, {
        name: this.allMods.find((mod) => mod.modInfo.mod_reference === update.item)?.modInfo.name || update.item,
      }));
      this.nextCheckForUpdates = setTimeout(() => this.checkForUpdates(), UPDATE_CHECK_INTERVAL);
    },
    updateSMMNow() {
      this.$root.$emit('downloadUpdate');
      this.smmUpdateDialog = false;
    },
    async updateAll() {
      await this.$store.dispatch('updateMulti', this.filteredModUpdates);
      const currentUpdates = this.filteredModUpdates;
      this.modUpdates.removeWhere((update) => currentUpdates.includes(update));
      if (this.filteredModUpdates.length === 0) {
        this.modUpdatesDialog = false;
      }
    },
    async updateItem(update) {
      await this.$store.dispatch('updateSingle', update);
      this.modUpdates.remove(update);
      if (this.filteredModUpdates.length === 0) {
        this.modUpdatesDialog = false;
      }
    },
    ignoreUpdate(update) {
      this.ignoredUpdates = ignoreUpdate(update.item, update.version);
    },
    unignoreUpdate(update) {
      this.ignoredUpdates = unignoreUpdate(update.item, update.version);
    },
    isIgnored(update) {
      return this.ignoredUpdates.some((ignoredUpdate) => ignoredUpdate.item === update.item && ignoredUpdate.version === update.version);
    },
    viewChangelog(update) {
      this.viewChangelogUpdate = update;
      this.changelogDialog = true;
    },
    changelogHTML(release) {
      return markdownAsElement(release.changelog).innerHTML;
    },
    modProgress(mod) {
      return this.inProgress.find((prog) => prog.id === mod);
    },
    currentModProgress(mod) {
      return lastElement(this.modProgress(mod).progresses);
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
.custom.v-list {
  background-color: var(--v-background-base);
}
.custom.v-list .v-list-item__action {
  margin: 0;
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
