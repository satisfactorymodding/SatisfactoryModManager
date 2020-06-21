<template>
  <div class="titlebar">
    <div
      class="d-inline-flex align-items-center"
    >
      <AppMenu
        ref="appMenu"
        :available-s-m-m-update="availableSMMUpdate"
        :filtered-mod-updates="filteredModUpdates"
        :show-ignored-updates.sync="showIgnoredUpdates"
        :update-check-mode.sync="updateCheckMode"
        @checkForUpdates="checkForUpdates"
        @addUpdateListener="addUpdateListener"
        @openSMMUpdateDialog="openSMMUpdateDialog"
        @openModUpdatesDialog="openModUpdatesDialog"
      />
    </div>
    <div class="bar">
      <div class="dragregion">
        <span class="app-title">{{ title }}</span>
      </div>
      <div
        class="button minimize"
        @click="minimize"
      >
        <span>&#128469;</span>
      </div>
      <div
        class="button maximize"
        @click="maximize"
      >
        <span>{{ isMaximized ? '&#128471;' : '&#128470;' }}</span>
      </div>
      <div
        class="button close"
        @click="close"
      >
        <span>&#128473;</span>
      </div>
    </div>
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
  </div>
</template>

<script>
import { mapState } from 'vuex';
import { getSetting, saveSetting } from '../../../settings';
import {
  ignoreUpdate, unignoreUpdate,
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
      cachedUpdateCheckMode: 'launch',
      isMaximized: false,
    };
  },
  computed: {
    ...mapState([
      'inProgress',
      'selectedInstall',
      'selectedProfile',
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
    updateCheckMode: {
      get() {
        return this.cachedUpdateCheckMode;
      },
      set(value) {
        saveSetting('updateCheckMode', value);
        this.cachedUpdateCheckMode = value;
      },
    },
  },
  watch: {
    async selectedInstall() {
      await this.checkForUpdates();
      if (this.filteredModUpdates.length > 0) {
        this.openModUpdatesDialog();
      }
    },
    async selectedProfile() {
      await this.checkForUpdates();
      if (this.filteredModUpdates.length > 0) {
        this.openModUpdatesDialog();
      }
    },
  },
  created() {
    this.$electron.remote.getCurrentWindow().on('maximize', this.onMaximize);
    this.$electron.remote.getCurrentWindow().on('unmaximize', this.onUnmaximize);
  },
  mounted() {
    this.ignoredUpdates = getSetting('ignoredUpdates', []);
    this.cachedUpdateCheckMode = getSetting('updateCheckMode', 'launch');
    this.isMaximized = this.$electron.remote.getCurrentWindow().isMaximized();

    if (this.updateCheckmode === 'launch') {
      this.$root.$once('doneLaunchUpdateCheck', () => {
        this.addUpdateListener();
      });
    } else {
      this.addUpdateListener();
    }
    this.nextCheckForUpdates = setTimeout(() => this.checkForUpdates(), UPDATE_CHECK_INTERVAL);
  },
  destroyed() {
    this.$electron.remote.getCurrentWindow().off('maximize', this.onMaximize);
    this.$electron.remote.getCurrentWindow().off('unmaximize', this.onUnmaximize);
  },
  methods: {
    openSMMUpdateDialog() {
      this.$refs.smmUpdateDialog.smmUpdateDialog = true;
    },
    openModUpdatesDialog() {
      this.$refs.modUpdatesDialog.modUpdatesDialog = true;
    },
    minimize() {
      this.$electron.remote.getCurrentWindow().minimize();
    },
    maximize() {
      if (!this.isMaximized) {
        this.$electron.remote.getCurrentWindow().maximize();
      } else {
        this.$electron.remote.getCurrentWindow().unmaximize();
      }
    },
    close() {
      if (this.inProgress.length === 0) {
        this.$electron.remote.getCurrentWindow().close();
      }
    },
    onMaximize() {
      this.isMaximized = true;
    },
    onUnmaximize() {
      this.isMaximized = false;
    },
    addUpdateListener() {
      this.$electron.ipcRenderer.on('updateAvailable', (e, updateInfo) => {
        this.availableSMMUpdate = updateInfo;
        if (this.updateCheckMode === 'ask' || this.updateCheckmode === 'launch') {
          this.$refs.smmUpdateDialog.smmUpdateDialog = true;
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
  background-color: gray;
  color: white !important;
}
.close:hover {
  background-color: red;
  color: white !important;
}
.button>span.dash {
  vertical-align: sub;
  margin-top: 0px;
}
</style>
