<template>
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
                <v-list-item-title>
                  {{ hasUpdate || availableSMMUpdate ? 'Updates ready to install' : 'Update settings' }}
                </v-list-item-title>
              </v-list-item-content>
              <v-list-item-action>
                <v-icon>mdi-chevron-right</v-icon>
              </v-list-item-action>
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

              <v-list-item @click="$emit('checkForUpdates')">
                <v-list-item-action />
                <v-list-item-content>
                  <v-list-item-title>Check for updates</v-list-item-title>
                </v-list-item-content>
                <v-list-item-action />
              </v-list-item>

              <v-divider
                insert
                class="custom"
              />

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
                  <v-switch
                    v-model="showIgnoredUpdatesLocal"
                  />
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
                @click="$emit('openModUpdatesDialog')"
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
                @click="$emit('openSMMUpdateDialog')"
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

        <v-divider class="custom" />

        <v-menu
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
                <v-list-item-title>Profiles</v-list-item-title>
              </v-list-item-content>
              <v-list-item-action>
                <v-icon>mdi-chevron-right</v-icon>
              </v-list-item-action>
            </v-list-item>
          </template>
          <v-card class="app-menu">
            <v-list>
              <v-list-item @click="exportProfile">
                <v-list-item-action>
                  <v-icon color="text">
                    mdi-content-save
                  </v-icon>
                </v-list-item-action>
                <v-list-item-content>
                  <v-list-item-title>Export profile</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-divider
                class="custom"
                inset
              />

              <v-list-item @click="importProfileDialog = true">
                <v-list-item-action>
                  <v-icon color="text">
                    mdi-download
                  </v-icon>
                </v-list-item-action>
                <v-list-item-content>
                  <v-list-item-title>Import profile</v-list-item-title>
                </v-list-item-content>
                <v-dialog
                  v-model="importProfileDialog"
                >
                  <v-card>
                    <v-card-title>
                      Import profile
                    </v-card-title>
                    <v-card-text>
                      <v-form
                        ref="importProfileForm"
                        v-model="importProfileFormValid"
                        @update="importProfileFormValid = $event"
                      >
                        <v-file-input
                          v-model="importProfileFile"
                          label="Profile file"
                          accept=".smmprofile"
                          required
                          :rules="[v => !!v || 'Choose a profile to import']"
                        />
                        <v-text-field
                          v-model="importProfileName"
                          label="Name"
                          required
                          :rules="[v => !!v || 'Profile name is required']"
                        />
                        <v-switch
                          v-model="importProfileVersions"
                          label="Import mod versions"
                        />
                        <span class="warning--text">{{ importProfileMessage }}</span>
                      </v-form>
                    </v-card-text>
                    <v-card-actions>
                      <v-btn
                        color="primary"
                        text
                        @click="importProfile"
                      >
                        Import
                      </v-btn>
                      <v-btn
                        color="primary"
                        text
                        @click="importProfileDialog = false"
                      >
                        Cancel
                      </v-btn>
                    </v-card-actions>
                  </v-card>
                </v-dialog>
              </v-list-item>

              <v-divider
                inset
                class="custom"
              />
            </v-list>
          </v-card>
        </v-menu>

        <v-divider class="custom" />

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
            <v-switch
              v-model="expandModInfoOnStart"
            />
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
            <v-switch
              v-model="darkMode"
            />
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
                If something doesn't behave as expected, the first thing to try is <a
                  @click="clearCache"
                >clearing
                  the cache</a><br>
                If that doesn't work, <a @click="debugMode = true">enable debug mode</a>, recreate what you
                tried
                to do and went wrong, then <a @click="exportDebugData">generate debug data</a>
                and upload the generated zip to the modding discord, #help-using-mods channel<br>
                <br>
                <v-divider />
                <br>
                <h3>X Satisfactory installs were found error</h3>
                <h4>Epic Games</h4>
                First, check that Epic can start the game. If you changed where the game is located, you
                need to make Epic update its install info. To do so:<br>
                1. Rename the game folder to something temporary<br>
                2. Start install from Epic to the directory you want the game to be in (the original
                folder
                name, before step 1)<br>
                3. After it downloads a bit close Epic<br>
                4. Copy back the files from the temporary folder EXCEPT the .egstore folder<br>
                5. Start Epic and resume the install so it finds that it is actually already
                installed<br>
                <br>
                <v-divider />
                <br>
                <h3>No Satisfactory installs found</h3>
                <h4>Epic Games / Steam</h4>
                Make sure you have the game installed, and Epic/Steam can find and launch it.<br>
                <h4>Cracked</h4>
                We do not support piracy, thus Satisfactory Mod Manager does not work with cracked
                copies of
                the game<br>
                <br>
                <v-divider />
                <br>
                <h3>Your issue is not here</h3>
                Ask for help in the modding discord, #help-using-mods channel
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

        <v-menu
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
                <v-list-item-title>Debug</v-list-item-title>
              </v-list-item-content>
              <v-list-item-action>
                <v-icon>mdi-chevron-right</v-icon>
              </v-list-item-action>
            </v-list-item>
          </template>
          <v-card class="app-menu">
            <v-list>
              <v-list-item>
                <v-list-item-action />
                <v-list-item-content>
                  <v-list-item-title>Debug mode</v-list-item-title>
                </v-list-item-content>

                <v-list-item-action>
                  <v-switch
                    v-model="debugMode"
                  />
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

              <v-list-item
                @click="exportDebugData"
              >
                <v-list-item-action />
                <v-list-item-content>
                  <v-list-item-title>Generate debug info</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-divider
                inset
                class="custom"
              />
            </v-list>
          </v-card>
        </v-menu>

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
</template>

<script>
import {
  clearCache, getLogFilePath, setDebug, validAndGreater,
} from 'satisfactory-mod-manager-api';
import JSZip from 'jszip';
import fs from 'fs';
import path from 'path';
import { getCacheFolder } from 'platform-folders';
import StreamZip from 'node-stream-zip';
import { filenameFriendlyDate } from '../../utils';
import { getSetting, saveSetting } from '../../../settings';

export default {
  props: {
    availableSMMUpdate: {
      type: Object,
      default: () => ({}),
    },
    filteredModUpdates: {
      type: Array,
      default: () => [],
    },
    showIgnoredUpdates: {
      type: Boolean,
      default: false,
    },
    updateCheckMode: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      attributionDialog: false,
      creditsDialog: false,
      importProfileFormValid: true,
      importProfileDialog: false,
      importProfileFile: null,
      importProfileMetadata: null,
      importProfileName: '',
      importProfileVersions: false,
      importProfileMessage: '',
      helpDialog: false,
      cachedDebugMode: false,
      menuOpen: false,
    };
  },
  computed: {
    version() {
      return this.$electron.remote.app.getVersion();
    },
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
    hasUpdate() {
      return !!this.availableSMMUpdate || this.filteredModUpdates.length > 0;
    },
    smmUpdateCount() {
      if (!this.availableSMMUpdate) {
        return 0;
      }
      return this.availableSMMUpdate.releaseNotes.length;
    },
    expandModInfoOnStart: {
      get() {
        return this.$store.state.expandModInfoOnStart;
      },
      set(value) {
        this.$store.dispatch('setExpandModInfoOnStart', value);
      },
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
        return this.cachedDebugMode;
      },
      set(value) {
        setDebug(value);
        saveSetting('debugMode', value);
        this.cachedDebugMode = value;
      },
    },
    showIgnoredUpdatesLocal: {
      get() {
        return this.showIgnoredUpdates;
      },
      set(value) {
        this.$emit('update:showIgnoredUpdates', value);
      },
    },
  },
  watch: {
    async importProfileFile(file) {
      if (file) {
        const zipData = new StreamZip({ file: file.path });
        try {
          await new Promise((resolve) => zipData.on('ready', resolve));
          const metadata = JSON.parse(zipData.entryDataSync('metadata.json').toString('utf8'));
          zipData.close();
          this.importProfileMetadata = metadata;
        } catch (e) {
          zipData.close();
          this.importProfileMetadata = null;
          this.$store.dispatch('showError', e);
        }
      } else {
        this.importProfileMetadata = null;
      }
    },
    importProfileMetadata(metadata) {
      if (metadata) {
        if (validAndGreater(metadata.gameVersion, this.$store.state.selectedInstall.version)) {
          this.importProfileMessage = `This profile is made for game version ${metadata.gameVersion}, but you're using an older version: ${this.$store.state.selectedInstall.version}. Things might not work as expected.`;
        } else {
          this.importProfileMessage = '';
        }
      } else {
        this.importProfileMessage = '';
      }
    },
  },
  mounted() {
    this.cachedDebugMode = getSetting('debugMode', false);

    if (this.debugMode) {
      setDebug(this.debugMode);
    }
  },
  methods: {
    moddingDiscord() {
      this.$electron.shell.openExternal('https://discord.gg/TShj39G');
    },
    officialDiscord() {
      this.$electron.shell.openExternal('https://discord.gg/Satisfactory');
    },
    clearCache() {
      clearCache();
      if (this.$store.state.selectedInstall) {
        this.$store.state.selectedInstall.clearCache();
      }
    },
    async exportDebugData() {
      const debugDataZip = new JSZip();
      const metadata = {
        installsFound: this.$store.state.satisfactoryInstalls,
        selectedInstall: this.$store.state.selectedInstall,
        profiles: this.$store.state.profiles,
        selectedProfile: this.$store.state.selectedProfile,
        installedMods: this.$store.state.selectedInstall?.mods,
        smlVersion: this.$store.state.selectedInstall?.smlVersion,
        bootstrapperVersion: this.$store.state.selectedInstall?.bootstrapperVersion,
        smmVersion: this.version,
      };
      debugDataZip.file('SatisfactoryModManager.log', fs.createReadStream(getLogFilePath()));
      if (this.$store.state.selectedInstall) {
        debugDataZip.file('pre-launch-debug.log', fs.createReadStream(path.join(this.$store.state.selectedInstall.installLocation, 'pre-launch-debug.log')));
        debugDataZip.file('SatisfactoryModLoader.log', fs.createReadStream(path.join(this.$store.state.selectedInstall.installLocation, 'SatisfactoryModLoader.log')));
        debugDataZip.file('FactoryGame.log', fs.createReadStream(path.join(getCacheFolder(), 'FactoryGame', 'Saved', 'Logs', 'FactoryGame.log')));
      }
      debugDataZip.file('metadata.json', JSON.stringify(metadata, null, 4));

      const result = this.$electron.remote.dialog.showSaveDialogSync(this.$electron.remote.getCurrentWindow(), {
        title: 'Save debug data as',
        filters: [
          { name: 'SMM Debug Data', extensions: ['zip'] },
        ],
        defaultPath: `SMMDebug_${filenameFriendlyDate(new Date())}.zip`,
      });

      if (result) {
        await new Promise((resolve, reject) => {
          debugDataZip.generateNodeStream().pipe(fs.createWriteStream(result)).on('finish', resolve).on('error', reject);
        });
      }
    },
    async importProfile() {
      if (this.$refs.importProfileForm.validate()) {
        const importProfileProgress = {
          id: '__importProfile__',
          progresses: [{
            id: '', progress: -1, message: `Importing profile as ${this.importProfileName}`, fast: false,
          }],
        };
        this.$store.state.inProgress.push(importProfileProgress);
        try {
          await this.$store.state.selectedInstall.importProfile(this.importProfileFile.path, this.importProfileName, this.importProfileVersions);
          this.$store.state.inProgress.remove(importProfileProgress);
          const newProfile = { name: this.importProfileName, items: [] }; // TODO: Items
          this.$store.state.profiles.push(newProfile);
          await this.$store.dispatch('selectProfile', newProfile);
        } catch (e) {
          this.$store.dispatch('showError', e);
          this.$store.state.inProgress.remove(importProfileProgress);
        }
        this.importProfileFile = null;
        this.importProfileName = '';
        this.importProfileVersions = false;
        this.importProfileDialog = false;
      }
    },
    async exportProfile() {
      const result = this.$electron.remote.dialog.showSaveDialogSync(this.$electron.remote.getCurrentWindow(), {
        title: 'Export profile as',
        filters: [
          { name: 'SMM Profile', extensions: ['smmprofile'] },
        ],
        defaultPath: `${this.$store.state.selectedProfile.name}.smmprofile`,
      });
      if (result) {
        const exportProfileProgress = {
          id: '__exportProfile__',
          progresses: [{
            id: '', progress: -1, message: `Exporting profile ${this.$store.state.selectedProfile.name}`, fast: false,
          }],
        };
        this.$store.state.inProgress.push(exportProfileProgress);
        try {
          await this.$store.state.selectedInstall.exportProfile(result);
        } catch (e) {
          this.$store.dispatch('showError', e);
        }
        this.$store.state.inProgress.remove(exportProfileProgress);
      }
    },
  },
};
</script>

<style scoped>
.app-icon {
  padding: 3px 0px 0px 3px;
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
  font-size: 18px !important;
}
.v-list-item {
  padding-left: 10px !important;
}
.v-list-item__action:first-child {
  margin-right: 0px !important;
}
.custom.v-divider--inset:not(.v-divider--vertical) {
  margin-left: 30px !important;
  max-width: calc(100% - 60px) !important;
}
.custom.v-divider:not(.v-divider--inset):not(.v-divider--vertical) {
  margin-left: 10px !important;
  max-width: calc(100% - 40px) !important;
}
</style>
