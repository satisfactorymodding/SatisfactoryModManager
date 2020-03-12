<template>
  <main>
    <div class="container-fluid my-2 content d-flex flex-column">
      <div class="row flex-grow-0 flex-shrink-0">
        <div class="col-5">
          <select
            v-model="selectedSatisfactoryInstall"
            class="form-control"
          >
            <option
              v-for="install in satisfactoryInstalls"
              :key="install.id"
              :value="install"
            >
              {{ install.displayName }}
            </option>
          </select>
        </div>
        <div class="col-auto">
          <div class="column">
            <button
              class="btn btn-primary"
              @click="launchSatisfactory"
            >
              Launch Satisfactory
            </button>
          </div>
        </div>
        <div class="col-auto d-inline-flex align-items-center">
          <strong>SML: {{ selectedSatisfactoryInstall ? (selectedSatisfactoryInstall.smlVersion || 'Not installed') : 'Select a Satisfactory Install' }}</strong>
        </div>
      </div>
      <div
        class="row my-2 flex-fill container-fluid my-2"
        style="font-size: 14px;"
      >
        <div
          class="row selection-row"
          style="height: 50%"
        >
          <div
            class="col-8 d-flex flex-column"
            style="height: 100%"
          >
            <input
              v-model="search"
              class="form-control flex-shrink-0 flex-grow-0"
              type="text"
              placeholder="Search"
              aria-label="Search"
            >
            <br>
            <list
              v-if="searchMods"
              v-model="selectedMod"
              :objects="searchMods"
              :can-select="true"
              class="flex-fill"
            >
              <template slot-scope="{item}">
                <div class="col-1">
                  <img
                    :src="item.logo || noImageURL"
                    width="100%"
                  >
                </div>
                <div
                  class="col-3 d-inline-flex align-items-center text-break"
                >
                  <strong>{{ item.name || '' }}</strong>
                </div>
                <div
                  class="col-1 d-inline-flex align-items-center"
                >
                  <strong>{{ item.versions[0] ? item.versions[0].version : 'N/A' }}</strong>
                </div>
                <div
                  class="col-3 d-inline-flex align-items-center"
                >
                  <strong>{{ item.authors.map((author) => author.user.username).join(', ') }}</strong>
                </div>
                <div
                  class="col-2 d-inline-flex align-items-center"
                >
                  <strong>{{ (item.last_version_date || new Date(0, 0, 0)).toLocaleDateString() }}</strong>
                </div>
                <div
                  class="col-2 d-inline-flex align-items-center"
                >
                  <button
                    :class="'my-1 btn ' + ((!item.versions[0] || isModVersionInstalled(item.versions[0])) ? 'btn-secondary' : 'btn-primary')"
                    style="width: 100%"
                    :disabled="!item.versions[0]"
                    @click="toggleModInstalled(item.versions[0])"
                  >
                    {{ !item.versions[0] ? 'N/A' : (isModVersionInstalled(item.versions[0]) ? "Remove" : "Install") }}
                  </button>
                </div>
              </template>
            </list>
          </div>
          <div
            class="col-4"
            style="height: 100%"
          >
            <list
              v-if="selectedMod && selectedMod.versions"
              :objects="selectedMod.versions"
              :can-select="false"
            >
              <template slot-scope="{item}">
                <div
                  class="col-3"
                  style="min-width: 150px"
                >
                  <button
                    :class="'my-1 btn ' + (isModVersionInstalled(item) ? 'btn-secondary' : 'btn-primary')"
                    style="width: 100%"
                    @click="toggleModInstalled(item)"
                  >
                    {{ isModVersionInstalled(item) ? "Remove" : "Install" }}
                  </button>
                </div>
                <div class="col-2 d-inline-flex align-items-center">
                  <strong>{{ item.version }}</strong>
                </div>
                <div class="col-1">
                  <div
                    v-if="inProgress.includes(item)"
                    class="spinner-border"
                    role="status"
                  >
                    <span class="sr-only">Loading...</span>
                  </div>
                </div>
              </template>
            </list>
          </div>
        </div>
        <div
          v-if="selectedMod != null"
          class="row"
          style="overflow: auto; margin: 10px"
        >
          <!-- eslint-disable-next-line vue/no-v-html -->
          <div v-html="compiledMarkdownDescription" />
        </div>
      </div>
    </div>
    <b-modal
      id="modal-install"
      ref="modal"
      title="Install Mod"
      @ok="handleModalInstallOk"
    >
      <form
        ref="form"
        @submit.stop.prevent="handleModalInstallSubmit"
      >
        <select
          v-model="selectedSatisfactoryInstall"
          class="form-control"
        >
          <option
            v-for="install in satisfactoryInstalls"
            :key="install.id"
            :value="install"
          >
            {{ install.displayName }}
          </option>
        </select>
        <p>Mod: {{ selectedMod.name }}</p>
        <p>Version: {{ modalInstallModVersion.version }}</p>
      </form>
    </b-modal>
    <b-modal
      id="modal-uninstall"
      ref="modal"
      title="Uninstall Mod"
      @ok="handleModalUninstallOk"
    >
      <form
        ref="form"
        @submit.stop.prevent="handleModalUninstallSubmit"
      >
        <select
          v-model="selectedSatisfactoryInstall"
          class="form-control"
        >
          <option
            v-for="install in satisfactoryInstalls"
            :key="install.id"
            :value="install"
          >
            {{ install.displayName }}
          </option>
        </select>
        <p>Mod: {{ selectedMod.name }}</p>
      </form>
    </b-modal>
  </main>
</template>

<script>
// TODO: display errors
import semver from 'semver';
import {
  getLatestSMLVersion,
  getInstalls,
  getAvailableMods,
} from 'satisfactory-mod-launcher-api';
import marked from 'marked';
import { spawn } from 'child_process';
import sanitizeHtml from 'sanitize-html';
import List from './List';


export default {
  name: 'Launcher',
  components: {
    List,
  },
  data() {
    return {
      selectedSatisfactoryInstall: null,
      satisfactoryInstalls: [],
      availableMods: [],
      installedSMLVersion: '',
      latestSMLVersion: {},
      SMLInProgress: false,
      selectedMod: {},
      searchMods: [],
      search: '',
      inProgress: [],
      modalInstallModVersion: {},
    };
  },
  computed: {
    noImageURL() {
      return 'https://ficsit.app/static/assets/images/no_image.png';
    },
    compiledMarkdownDescription() {
      return sanitizeHtml(marked(this.selectedMod.full_description || ''));
    },
    hasSMLUpdates() {
      return (
        !semver.valid(this.installedSMLVersion)
        || (semver.valid(this.latestSMLVersion.version)
          && semver.lt(this.installedSMLVersion, this.latestSMLVersion.version))
      );
    },
    isSMLInstalled() {
      return !!semver.valid(this.installedSMLVersion);
    },
  },
  watch: {
    search() {
      this.refreshSearch();
    },
  },
  mounted() {
    this.$electron.ipcRenderer.on('openedByUrl', (e, url) => {
      const parsed = new URL(url);
      const command = parsed.pathname.replace(/^\/+|\/+$/g, '');
      if (command === 'install') {
        const modID = parsed.searchParams.get('modID');
        const version = parsed.searchParams.get('version');
        this.selectedMod = this.availableMods.find((mod) => mod.id === modID);
        this.modalInstallModVersion = this.selectedMod.versions.find((ver) => ver.version === version);
        this.$bvModal.show('modal-install');
      } else if (command === 'uninstall') {
        const modID = parsed.searchParams.get('modID');
        this.selectedMod = this.availableMods.find((mod) => mod.id === modID);
        this.$bvModal.show('modal-uninstall');
      }
    });
  },
  created() {
    Promise.all(
      [
        this.refreshSatisfactoryInstalls(),
        this.refreshAvailableMods(),
        getLatestSMLVersion().then((smlVersion) => {
          this.latestSMLVersion = smlVersion.version;
        }),
      ],
    ).then(() => {
      this.$electron.ipcRenderer.send('vue-ready');
    });
  },
  methods: {
    handleModalInstallOk(bvModalEvt) {
      bvModalEvt.preventDefault();
      this.handleModalInstallSubmit();
    },
    handleModalInstallSubmit() {
      this.installMod(this.modalInstallModVersion);
      this.$nextTick(() => {
        this.$bvModal.hide('modal-install');
      });
    },
    handleModalUninstallOk(bvModalEvt) {
      bvModalEvt.preventDefault();
      this.handleModalInstallSubmit();
    },
    handleModalUninstallSubmit() {
      this.uninstallMod(this.selectedMod.versions.find((ver) => this.isModVersionInstalled(ver)));
      this.$nextTick(() => {
        this.$bvModal.hide('modal-uninstall');
      });
    },
    refreshSearch() {
      this.searchMods = this.availableMods.filter((mod) => mod.name.toLowerCase().includes(this.search.toLowerCase()));
    },
    refreshAvailableMods() {
      return getAvailableMods().then((mods) => {
        this.availableMods = mods;
        this.refreshSearch();
      });
    },
    isModVersionInstalled(modVersion) {
      if (modVersion && modVersion.mod_id && modVersion.version) {
        return this.selectedSatisfactoryInstall.mods[modVersion.mod_id] === modVersion.version;
      }
      return false;
    },
    refreshCurrentMod() {
      const currentModId = this.selectedMod.mod_id;
      this.refreshAvailableMods().then(() => {
        this.selectedMod = this.searchMods.find((mod) => mod.mod_id === currentModId);
      });
    },
    installMod(modVersion) {
      return this.selectedSatisfactoryInstall
        .installMod(modVersion.mod_id, modVersion.version)
        .then(() => {
          this.inProgress.splice(this.inProgress.indexOf(modVersion));
          this.refreshCurrentMod();
        });
    },
    uninstallMod(modVersion) {
      return this.selectedSatisfactoryInstall
        .uninstallMod(modVersion.mod_id)
        .then(() => {
          this.inProgress.splice(this.inProgress.indexOf(modVersion));
          this.refreshCurrentMod();
        });
    },
    toggleModInstalled(modVersion) {
      this.inProgress.push(modVersion);
      if (this.isModVersionInstalled(modVersion)) {
        this.uninstallMod(modVersion);
      } else {
        this.installMod(modVersion);
      }
    },
    refreshSatisfactoryInstalls() {
      return getInstalls().then((installs) => {
        this.satisfactoryInstalls = installs;
        if (this.satisfactoryInstalls.length > 0) {
          const defaultInstall = this.satisfactoryInstalls[0];
          this.selectedSatisfactoryInstall = defaultInstall;
        }
      });
    },
    launchSatisfactory() {
      if (this.selectedSatisfactoryInstall) {
        spawn(this.selectedSatisfactoryInstall.launchPath, { detached: true }).unref();
      }
    },
    updateSML() {
      this.SMLInProgress = true;
      return this.selectedSatisfactoryInstall.updateSML().then(() => {
        this.SMLInProgress = false;
      });
    },
    uninstallSML() {
      this.SMLInProgress = true;
      return this.selectedSatisfactoryInstall.uninstallSML().then(() => {
        this.SMLInProgress = false;
      });
    },
  },
};
</script>

<style>
</style>
