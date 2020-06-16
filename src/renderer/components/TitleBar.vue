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

                  <v-list-item @click="checkForUpdates">
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
                    If that doesn't work, <a @click="enableDebug">enable debug mode</a>, recreate what you tried to do and went wrong, then <a @click="exportDebugData">generate debug data</a>
                    and upload the generated zip to the modding discord, #help-using-mods channel<br>
                    <br><v-divider /><br>
                    <h3>X Satisfactory installs were found error</h3>
                    <h4>Epic Games</h4>
                    First, check that Epic can start the game. If you changed where the game is located, you need to make Epic update its install info. To do so:<br>
                    1. Rename the game folder to something temporary<br>
                    2. Start install from Epic to the directory you want the game to be in (the original folder name, before step 1)<br>
                    3. After it downloads a bit close Epic<br>
                    4. Copy back the files from the temporary folder EXCEPT the .egstore folder<br>
                    5. Start Epic and resume the install so it finds that it is actually already installed<br>
                    <br><v-divider /><br>
                    <h3>No Satisfactory installs found</h3>
                    <h4>Epic Games / Steam</h4>
                    Make sure you have the game installed, and Epic/Steam can find and launch it.<br>
                    <h4>Cracked</h4>
                    We do not support piracy, thus Satisfactory Mod Manager does not work with cracked copies of the game<br>
                    <br><v-divider /><br>
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
    <v-dialog
      persistent
      :value="isProfileExportInProgress"
      width="500"
      height="230"
    >
      <v-card
        color="loadingBackground !important"
      >
        <v-row
          no-gutters
          justify="center"
        >
          <v-img
            class="mt-4"
            src="static/smm_icon.png"
            max-height="82px"
            max-width="87px"
          />
        </v-row>
        <v-card-title class="loading-text-main">
          EXPORTING PROFILE
        </v-card-title>

        <v-card-text
          v-if="isProfileExportInProgress"
          class="text-center"
        >
          <v-progress-linear
            :value="Math.round(currentProfileExportProgress.progress * 100)"
            :class="currentProfileExportProgress.fast ? 'fast' : ''"
            background-color="#000000"
            color="#5bb71d"
            height="2"
            reactive
            :indeterminate="currentProfileExportProgress.progress < 0"
          />
          {{ currentProfileExportProgress.message || '&nbsp;' }}
        </v-card-text>
      </v-card>
    </v-dialog>
    <v-dialog
      persistent
      :value="isProfileImportInProgress"
      width="500"
      height="230"
    >
      <v-card
        color="loadingBackground !important"
      >
        <v-row
          no-gutters
          justify="center"
        >
          <v-img
            class="mt-4"
            src="static/smm_icon.png"
            max-height="82px"
            max-width="87px"
          />
        </v-row>
        <v-card-title class="loading-text-main">
          IMPORTING PROFILE
        </v-card-title>

        <v-card-text
          v-if="isProfileImportInProgress"
          class="text-center"
        >
          <v-progress-linear
            :value="Math.round(currentProfileImportProgress.progress * 100)"
            :class="currentProfileImportProgress.fast ? 'fast' : ''"
            background-color="#000000"
            color="#5bb71d"
            height="2"
            reactive
            :indeterminate="currentProfileImportProgress.progress < 0"
          />
          {{ currentProfileImportProgress.message || '&nbsp;' }}
        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import {
  clearCache, validAndGreater, getLogFilePath, setDebug,
} from 'satisfactory-mod-manager-api';
import { mapState } from 'vuex';
import StreamZip from 'node-stream-zip';
import JSZip from 'jszip';
import fs from 'fs';
import path from 'path';
import { getCacheFolder } from 'platform-folders';
import { saveSetting, getSetting } from '../settings';
import {
  markdownAsElement, ignoreUpdate, unignoreUpdate, lastElement, filenameFriendlyDate,
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
      importProfileDialog: false,
      importProfileFile: null,
      importProfileName: '',
      importProfileVersions: false,
      importProfileFormValid: true,
      importProfileMetadata: null,
      importProfileMessage: '',
      cachedDebugMode: false,
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
        return this.cachedDebugMode;
      },
      set(value) {
        setDebug(value);
        saveSetting('debugMode', value);
        this.cachedDebugMode = value;
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
    isProfileExportInProgress() {
      return this.inProgress.some((prog) => prog.id === '__exportProfile__');
    },
    profileExportProgress() {
      return this.inProgress.find((prog) => prog.id === '__exportProfile__');
    },
    currentProfileExportProgress() {
      return lastElement(this.profileExportProgress.progresses);
    },
    isProfileImportInProgress() {
      return this.inProgress.some((prog) => prog.id === '__importProfile__');
    },
    profileImportProgress() {
      return this.inProgress.find((prog) => prog.id === '__importProfile__');
    },
    currentProfileImportProgress() {
      return lastElement(this.profileImportProgress.progresses);
    },
  },
  watch: {
    async selectedInstall() {
      await this.checkForUpdates();
    },
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
      if (this.inProgress.length === 0) {
        this.$electron.remote.getCurrentWindow().close();
      }
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
    async exportDebugData() {
      const debugDataZip = new JSZip();
      const metadata = {
        installsFound: this.$store.state.satisfactoryInstalls,
        selectedInstall: this.$store.state.selectedInstall,
        profile: this.$store.state.selectedProfile,
        installedMods: this.$store.state.selectedInstall?.mods,
        smlVersion: this.$store.state.selectedInstall?.smlVersion,
        bootstrapperVersion: this.$store.state.selectedInstall?.bootstrapperVersion,
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
  },
};
</script>

<style scoped>
.app-title {
  font-size: 16px !important;
  padding-top: 5px;
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
  font-size: 18px !important;
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
  max-width: calc(100% - 60px) !important;
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
