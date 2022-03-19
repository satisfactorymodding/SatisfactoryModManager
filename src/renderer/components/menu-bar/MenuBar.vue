<template>
  <v-container
    fluid
    class="pt-2 pb-2"
  >
    <v-row
      class="px-4 pl-5"
      :class="{ 'update-available': hasUpdate }"
    >
      <v-col cols="auto">
        <SettingsMenu />
      </v-col>
      <v-col cols="auto">
        <v-btn
          class="ma-2 ml-1 px-2 d-inline-flex align-center"
          style="height: 28px"
          @click="helpDialog = true"
        >
          <v-icon
            style="margin-right: 4px; font-size: 16px !important"
          >
            mdi-help-circle
          </v-icon>
          <span>Help</span>
        </v-btn>
      </v-col>
      <v-spacer />
      <v-col cols="auto">
        <span
          v-if="!hasUpdate"
          class="d-inline-flex align-center fill-height mx-3"
          style="font-size: 14px"
        >No updates right now</span>
        <v-btn
          v-else-if="filteredModUpdates.length === 0 || !availableSMMUpdate"
          class="my-2 mx-1"
          style="height: 28px"
          @click="filteredModUpdates.length > 0 ? openModUpdatesDialog() : openSMMUpdateDialog()"
        >
          <span class="mx-1">
            {{ filteredModUpdates.length > 0 ? 'Mod updates are available' : 'SMM updates are available' }}
          </span>
        </v-btn>
        <v-menu v-else>
          <template #activator="{ on, attrs }">
            <v-btn
              class="my-2 mx-1"
              style="height: 28px"
              v-bind="attrs"
              v-on="on"
            >
              <span class="mx-1">
                SMM and mod updates are available
              </span>
            </v-btn>
          </template>
          <v-card>
            <v-list class="menu">
              <v-list-item @click="openModUpdatesDialog">
                <v-list-item-action />
                <v-list-item-content>
                  <v-list-item-title>Mod updates ({{ filteredModUpdates.length }})</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-divider
                inset
                class="custom"
              />

              <v-list-item @click="openSMMUpdateDialog">
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
      </v-col>
      <v-col cols="auto">
        <UpdatesMenu
          :available-s-m-m-update="availableSMMUpdate"
          :filtered-mod-updates="filteredModUpdates"
          :show-ignored-updates.sync="showIgnoredUpdates"
          :update-check-mode.sync="updateCheckMode"
          :update-check-in-progress="updateCheckInProgress"
          @addUpdateListener="addUpdateListener"
          @openSMMUpdateDialog="openSMMUpdateDialog"
          @openModUpdatesDialog="openModUpdatesDialog"
          @checkForUpdates="manualCheckForUpdates"
        />
      </v-col>
    </v-row>
    <ModUpdatesDialog
      ref="modUpdatesDialog"
      :filtered-mod-updates="filteredModUpdates"
      :ignored-updates="ignoredUpdates"
      :is-ignored="isIgnored"
      @ignoreUpdate="ignoreUpdate"
      @unignoreUpdate="unignoreUpdate"
      @updateItem="updateItem"
      @updateAll="updateAll"
      @viewChangelog="viewChangelog"
    />
    <ChangelogDialog
      ref="changelogDialog"
      :view-changelog-update="viewChangelogUpdate"
    />
    <SMMUpdateDialog
      ref="smmUpdateDialog"
      :available-s-m-m-update="availableSMMUpdate"
      :smm-update-notes="smmUpdateNotes"
      @updateSMMNow="updateSMMNow"
    />
    <ProfileExportProgressDialog />
    <ProfileImportProgressDialog />
    <v-dialog
      v-model="helpDialog"
    >
      <v-card>
        <v-card-title>
          Help
        </v-card-title>
        <v-card-text>
          <h3>General troubleshooting</h3>
          If something doesn't behave as expected, the first thing to try is
          <a @click="clearCache">clearing the cache</a><br>
          If that doesn't work, <a @click="exportDebugData">generate debug data</a>
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
          folder name, before step 1)<br>
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
  </v-container>
</template>

<script>
import { mapState } from 'vuex';
import gql from 'graphql-tag';
import {
  ignoreUpdate, unignoreUpdate,
} from '@/utils';
import { getSetting, saveSetting } from '~/settings';
import SettingsMenu from './SettingsMenu';
import UpdatesMenu from './UpdatesMenu';
import ProfileImportProgressDialog from './dialogs/ProfileImportProgressDialog';
import ProfileExportProgressDialog from './dialogs/ProfileExportProgressDialog';
import SMMUpdateDialog from './dialogs/SMMUpdateDialog';
import ChangelogDialog from './dialogs/ChangelogDialog';
import ModUpdatesDialog from './dialogs/ModUpdatesDialog';

const UPDATE_CHECK_INTERVAL = 5 * 60 * 1000;

export default {
  components: {
    SettingsMenu,
    UpdatesMenu,
    ModUpdatesDialog,
    ChangelogDialog,
    SMMUpdateDialog,
    ProfileExportProgressDialog,
    ProfileImportProgressDialog,
  },
  data() {
    return {
      availableSMMUpdate: null,
      modUpdates: [],
      nextCheckForUpdates: -1,
      updateCheckInProgress: false,
      viewChangelogUpdate: null,
      showIgnoredUpdates: false,
      ignoredUpdates: [],
      cachedUpdateCheckMode: 'launch',
      helpDialog: false,
    };
  },
  computed: {
    ...mapState([
      'inProgress',
      'selectedInstall',
      'selectedProfile',
    ]),
    smmUpdateNotes() {
      if (!this.availableSMMUpdate) {
        return '';
      }
      return this.availableSMMUpdate.releaseNotes.map((release) => release.note.substr(release.note.indexOf('<h2>Changelog</h2>')).replace('<h2>Changelog</h2>', `<h2>v${release.version} changelog</h2>`)).join('\n');
    },
    hasUpdate() {
      return !!this.availableSMMUpdate || this.filteredModUpdates.length > 0;
    },
    filteredModUpdates() {
      return (this.showIgnoredUpdates ? this.modUpdates : this.modUpdates.filter((update) => !this.isIgnored(update)));
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
    smmUpdateCount() {
      if (!this.availableSMMUpdate) {
        return 0;
      }
      return this.availableSMMUpdate.releaseNotes.length;
    },
  },
  watch: {
    async selectedInstall() {
      await this.checkForUpdates();
    },
    async selectedProfile() {
      await this.checkForUpdates();
    },
  },
  mounted() {
    this.ignoredUpdates = getSetting('ignoredUpdates', []);
    this.cachedUpdateCheckMode = getSetting('updateCheckMode', 'launch');

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
    openSMMUpdateDialog() {
      this.$refs.smmUpdateDialog.smmUpdateDialog = true;
    },
    openModUpdatesDialog() {
      this.$refs.modUpdatesDialog.modUpdatesDialog = true;
    },
    addUpdateListener() {
      this.$electron.ipcRenderer.on('updateAvailable', (e, updateInfo) => {
        this.availableSMMUpdate = updateInfo;
        if (this.updateCheckMode === 'ask' || this.updateCheckmode === 'launch') {
          this.$refs.smmUpdateDialog.smmUpdateDialog = true;
        }
      });
    },
    async manualCheckForUpdates() {
      await this.checkForUpdates();
      if (this.filteredModUpdates.length > 0) {
        this.openModUpdatesDialog();
      }
    },
    async checkForUpdates() {
      if (this.updateCheckInProgress) {
        return;
      }
      this.updateCheckInProgress = true;
      clearTimeout(this.nextCheckForUpdates);
      // don't check for updates while something is in progress
      while (this.inProgress.length > 0) {
        // eslint-disable-next-line no-await-in-loop
        await new Promise((res) => setTimeout(() => res(), 500));
      }
      this.$electron.ipcRenderer.send('checkForUpdates');
      const modUpdates = await Promise.all((await this.$store.state.selectedInstall.checkForUpdates()).map(async (update) => Object.assign(update, {
        name: (await this.$apollo.query({
          query: gql`
            query getModName($modReference: ModReference!) {
              mod: getModByReference(modReference: $modReference) {
                id,
                name,
              }
            }
          `,
          variables: {
            modReference: update.item,
          },
        })).data.mod?.name || update.item,
      })));
      this.updateCheckInProgress = false;
      this.modUpdates = modUpdates;
      this.nextCheckForUpdates = setTimeout(() => this.checkForUpdates(), UPDATE_CHECK_INTERVAL);
    },
    updateSMMNow() {
      this.$root.$emit('downloadUpdate');
      this.refs.smmUpdateDialog.smmUpdateDialog = false;
    },
    async updateAll() {
      await this.$store.dispatch('updateMulti', this.filteredModUpdates);
      const currentUpdates = this.filteredModUpdates;
      this.modUpdates.removeWhere((update) => currentUpdates.includes(update));
      if (this.filteredModUpdates.length === 0) {
        this.$refs.modUpdatesDialog.modUpdatesDialog = false;
      }
    },
    async updateItem(update) {
      await this.$store.dispatch('updateSingle', update);
      this.modUpdates.remove(update);
      if (this.filteredModUpdates.length === 0) {
        this.$refs.modUpdatesDialog.modUpdatesDialog = false;
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
      this.$refs.changelogDialog.changelogDialog = true;
    },
    exportDebugData() {
      this.$root.$emit('exportDebugData');
    },
    clearCache() {
      this.$root.$emit('clearCache');
    },
  },
};
</script>

<style>
.menu.v-list {
  background-color: var(--v-backgroundMenu-base) !important;
}
.menu .custom.v-list .v-list-item__action {
  margin: 0;
}
.menu .v-icon {
  font-size: 18px !important;
}
.menu .v-list-item {
  padding-left: 10px !important;
}
.menu .v-list-item__action:first-child {
  margin-right: 0px !important;
}
.menu .custom.v-divider--inset:not(.v-divider--vertical) {
  margin-left: 30px !important;
  max-width: calc(100% - 60px) !important;
}
.menu .custom.v-divider:not(.v-divider--inset):not(.v-divider--vertical) {
  margin-left: 10px !important;
  max-width: calc(100% - 40px) !important;
}
</style>

<style scoped>
.row.update-available, .row.update-available>* {
  -webkit-animation: update-availabe-anim 2s linear infinite alternate both;
          animation: update-availabe-anim 2s linear infinite alternate both;
  color: black !important;
}
.container {
  background-color: var(--v-backgroundMenuBar-base);
}
.col {
  padding: 0;
  padding-top: 8px;
  padding-bottom: 8px;
}

@-webkit-keyframes update-availabe-anim {
  0% {
    background-color: var(--v-backgroundMenuBar-base);
  }
  100% {
    background-color: var(--v-primary-base);
  }
}
@keyframes update-availabe-anim {
  0% {
    background-color: var(--v-backgroundMenuBar-base);
  }
  100% {
    background-color: var(--v-primary-base);
  }
}

</style>
