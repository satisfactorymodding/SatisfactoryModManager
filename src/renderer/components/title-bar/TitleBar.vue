<template>
  <div class="titlebar">
    <div
      class="d-inline-flex align-items-center"
      @click="$refs.appMenu.menuOpen = !$refs.appMenu.menuOpen"
    >
      <AppMenu
        ref="appMenu"
        :available-s-m-m-update="availableSMMUpdate"
        :check-for-updates="checkForUpdates"
        :filtered-mod-updates="filteredModUpdates"
        :show-ignored-updates="showIgnoredUpdates"
        :add-update-listener="addUpdateListener"
        :set-show-ignored-updates="setShowIgnoredUpdates"
        :open-smm-update-dialog="openSmmUpdateDialog"
        :open-mod-updates-dialog="openModUpdatesDialog"
      />
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
    <ModUpdatesDialog
      :filtered-mod-updates="filteredModUpdates"
      :ignore-update="ignoreUpdate"
      :in-progress="inProgress"
      :is-ignored="isIgnored"
      :unignore-update="unignoreUpdate"
      :update-all="updateAll"
      :update-item="updateItem"
      :view-changelog="viewChangelog"
    />
    <ChangelogDialog
      :changelog-h-t-m-l="changelogHTML"
      :view-changelog-update="viewChangelogUpdate"
    />
    <SMMUpdateDialog
      ref="smmUpdateDialog"
      :available-s-m-m-update="availableSMMUpdate"
      :smm-update-notes="smmUpdateNotes"
      :update-s-m-m-now="updateSMMNow"
    />
    <ProfileExportProgressDialog :in-progress="inProgress" />
    <ProfileImportProgressDialog :in-progress="inProgress" />
  </div>
</template>

<script>
import {
  validAndGreater,
} from 'satisfactory-mod-manager-api';
import { mapState } from 'vuex';
import { getSetting } from '../../settings';
import {
  ignoreUpdate, markdownAsElement, unignoreUpdate,
} from '../../utils';
import ProfileImportProgressDialog from './ProfileImportProgressDialog';
import ProfileExportProgressDialog from './ProfileExportProgressDialog';
import SMMUpdateDialog from './SMMUpdateDialog';
import ChangelogDialog from './ChangelogDialog';
import ModUpdatesDialog from './ModUpdatesDialog';
import AppMenu from './AppMenu';

const UPDATE_CHECK_INTERVAL = 5 * 60 * 1000;

export default {
  components: {
    AppMenu,
    ModUpdatesDialog,
    ChangelogDialog,
    SMMUpdateDialog,
    ProfileExportProgressDialog,
    ProfileImportProgressDialog,
  },
  props: {
    title: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      availableSMMUpdate: null,
      modUpdates: [],
      nextCheckForUpdates: -1,
      viewChangelogUpdate: null,
      showIgnoredUpdates: false,
      ignoredUpdates: [],
      importProfileMetadata: null,
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
  },
  watch: {
    async selectedInstall() {
      await this.checkForUpdates();
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
    this.ignoredUpdates = getSetting('ignoredUpdates', []);
    this.nextCheckForUpdates = setTimeout(() => this.checkForUpdates(), UPDATE_CHECK_INTERVAL);
  },
  methods: {
    setShowIgnoredUpdates(value) {
      this.showIgnoredUpdates = value;
    },
    openSmmUpdateDialog() {
      this.refs.smmUpdateDialog.smmUpdateDialog = true;
    },
    openModUpdatesDialog() {
      this.$refs.modUpdatesDialog.modUpdatesDialog = true;
    },
    close() {
      if (this.inProgress.length === 0) {
        this.$electron.remote.getCurrentWindow().close();
      }
    },
    addUpdateListener() {
      this.$electron.ipcRenderer.on('updateAvailable', (e, updateInfo) => {
        this.availableSMMUpdate = updateInfo;
        if (this.updateCheckMode === 'ask') {
          this.$refs.smmUpdateDialog.smmUpdateDialogdateDialog = true;
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
    changelogHTML(release) {
      return markdownAsElement(release.changelog).innerHTML;
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
  padding-bottom: 4px;
}
.titlebar, .titlebar > * {
  color: var(--v-text2-base) !important;
  background-color: var(--v-background-base);
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
</style>
