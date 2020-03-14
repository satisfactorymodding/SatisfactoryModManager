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
          <strong>SML: {{ selectedSatisfactoryInstall ? (selectedSatisfactoryInstall.smlVersion || 'Install a mod to install SML') : 'Select a Satisfactory Install' }}</strong>
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
            <button
              v-b-modal.filters-modal
              class="btn btn-primary"
            >
              Sort/Filter
            </button>
            <br>
            <list
              v-if="searchMods"
              v-model="selectedMod"
              :objects="searchMods"
              :can-select="true"
              class="flex-fill"
            >
              <template slot-scope="{item}">
                <div
                  class="col-1 p-0"
                >
                  <img
                    :src="item.logo || noImageURL"
                    width="100%"
                    :style="!isModSML20Compatible(item) ? 'background-color: #837971' : ''"
                  >
                </div>
                <div
                  class="col-3 d-inline-flex align-items-center text-break"
                  :style="!isModSML20Compatible(item) ? 'background-color: #837971' : ''"
                >
                  <strong>{{ item.name || '' }}</strong>
                </div>
                <div
                  class="col-1 d-inline-flex align-items-center"
                  :style="!isModSML20Compatible(item) ? 'background-color: #837971' : ''"
                >
                  <strong>{{ item.versions[0] ? item.versions[0].version : 'N/A' }}</strong>
                </div>
                <div
                  class="col-3 d-inline-flex align-items-center"
                  :style="!isModSML20Compatible(item) ? 'background-color: #837971' : ''"
                >
                  <strong>{{ item.authors.map((author) => author.user.username).join(', ') }}</strong>
                </div>
                <div
                  class="col-2 d-inline-flex align-items-center"
                  :style="!isModSML20Compatible(item) ? 'background-color: #837971' : ''"
                >
                  <strong>{{ item.last_version_date ? item.last_version_date.toLocaleDateString() : 'N/A' }}</strong>
                </div>
                <div
                  class="col-2 d-inline-flex align-items-center"
                  :style="!isModSML20Compatible(item) ? 'background-color: #837971' : ''"
                >
                  <button
                    :class="'my-1 btn ' + ((!item.versions[0] || isModVersionInstalled(item.versions[0])) ? 'btn-secondary' : 'btn-primary')"
                    style="font-size: 13px; width: 100%"
                    :disabled="!item.versions[0] || !isModSML20Compatible(item)"
                    @click="toggleModInstalled(item.versions[0])"
                  >
                    {{ !item.versions[0] ? 'N/A' : (isModSML20Compatible(item) ? (isModVersionInstalled(item.versions[0]) ? "Remove" : (isModInstalled(item) ? "Update" : "Install")) : 'Outdated') }}
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
                    :disabled="!isVersionSML20Compatible(item)"
                    @click="toggleModInstalled(item)"
                  >
                    {{ isVersionSML20Compatible(item) ? (isModVersionInstalled(item) ? "Remove" : "Install") : "Outdated" }}
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
      title="Install Mod"
      @ok="handleModalInstallOk"
    >
      <form
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
        <label for="modalInstallVersion">Version:</label>
        <select
          id="modalInstallVersion"
          v-model="modalInstallModVersion"
          class="form-control"
        >
          <option
            v-for="version in selectedMod.versions"
            :key="version.version"
            :value="version"
          >
            {{ version.version }}
          </option>
        </select>
      </form>
    </b-modal>
    <b-modal
      id="modal-uninstall"
      title="Uninstall Mod"
      @ok="handleModalUninstallOk"
    >
      <form
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
    <b-modal
      id="filters-modal"
      title="Uninstall Mod"
      ok-only
    >
      <form>
        <b-form-checkbox
          v-model="filters.compatibleOnly"
          switch
          size="lg"
        >
          Compatible with Update 3
        </b-form-checkbox>
        <label for="sortBySelect">Sort by:</label>
        <select
          id="sortBySelect"
          v-model="filters.sortBy"
          class="form-control"
        >
          <option
            v-for="sortByOption in sortByOptions"
            :key="sortByOption.value"
            :value="sortByOption.value"
          >
            {{ sortByOption.displayName }}
          </option>
        </select>
        <label for="sortOrderSelect">Order:</label>
        <select
          id="sortOrderSelect"
          v-model="filters.sortOrder"
          class="form-control"
        >
          <option
            v-for="sortOrderOption in sortOrderOptions"
            :key="sortOrderOption.value"
            :value="sortOrderOption.value"
          >
            {{ sortOrderOption.displayName }}
          </option>
        </select>
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
  toggleDebug,
  clearCache,
} from 'satisfactory-mod-launcher-api';
import marked from 'marked';
import { exec } from 'child_process';
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
      filters: {
        compatibleOnly: false,
        sortBy: 'lastVersionDate', // lastVersionDate, popularity, hotness, downloads, views
        sortOrder: 'descending', // ascending, descending
      },
      inProgress: [],
      modalInstallModVersion: {},
      sortByOptions: [
        {
          value: 'name',
          displayName: 'Name',
        },
        {
          value: 'lastVersionDate',
          displayName: 'Last Version Date',
        },
        {
          value: 'popularity',
          displayName: 'Popularity (downloads)',
        },
        {
          value: 'hotness',
          displayName: 'Hotness (views)',
        },
        {
          value: 'downloads',
          displayName: 'Downloads',
        },
        {
          value: 'views',
          displayName: 'Views',
        },
      ],
      sortOrderOptions: [
        {
          value: 'ascending',
          displayName: 'Ascending',
        },
        {
          value: 'descending',
          displayName: 'Descending',
        },
      ],
    };
  },
  computed: {
    noImageURL() {
      return 'https://ficsit.app/static/assets/images/no_image.png';
    },
    compiledMarkdownDescription() {
      const html = sanitizeHtml(marked(this.selectedMod.full_description || ''));
      const el = document.createElement('html');
      el.innerHTML = html;
      const links = el.getElementsByTagName('a');
      for (let i = 0; i < links.length; i += 1) {
        links[i].target = '_blank';
      }
      return el.innerHTML;
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
    filters: {
      handler() {
        this.refreshSearch();
      },
      deep: true,
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
        this.modalInstallModVersion = this.selectedMod.versions.find((ver) => ver.version === version) || this.selectedMod.versions[0];
        this.$bvModal.show('modal-install');
      } else if (command === 'uninstall') {
        const modID = parsed.searchParams.get('modID');
        this.selectedMod = this.availableMods.find((mod) => mod.id === modID);
        this.$bvModal.show('modal-uninstall');
      }
    });
    this.$electron.ipcRenderer.on('toggleDebug', () => {
      toggleDebug();
    });
    this.$electron.ipcRenderer.on('clearCache', () => {
      clearCache();
      if (this.selectedSatisfactoryInstall) {
        this.selectedSatisfactoryInstall.clearCache();
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
    isModSML20Compatible(mod) {
      return mod.versions.length !== 0 && semver.satisfies(mod.versions[0].sml_version, '>=2.0.0');
    },
    isVersionSML20Compatible(version) {
      return semver.satisfies(version.sml_version, '>=2.0.0');
    },
    refreshSearch() {
      this.searchMods = this.availableMods.filter((mod) => mod.name.toLowerCase().includes(this.search.toLowerCase())
        && (!this.filters.compatibleOnly || this.isModSML20Compatible(mod)));
      this.searchMods.sort((modA, modB) => {
        switch (this.filters.sortBy) {
          case 'name':
            return modB.name.localeCompare(modA.name);
          case 'popularity':
            return modB.popularity - modA.popularity;
          case 'hotness':
            return modB.hotness - modA.hotness;
          case 'downloads':
            return modB.downloads - modA.downloads;
          case 'views':
            return modB.views - modA.views;
          case 'lastVersionDate':
          default:
            if (modB.last_version_date && modA.last_version_date) {
              return modB.last_version_date.getTime() - modA.last_version_date.getTime();
            }
            return 0;
        }
      });
      if (this.filters.sortOrder === 'ascending') {
        this.searchMods.reverse();
      }
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
    isModInstalled(mod) {
      return mod.versions.some((version) => this.isModVersionInstalled(version));
    },
    refreshCurrentMod() {
      const currentModId = this.selectedMod.id;
      this.refreshAvailableMods().then(() => {
        this.selectedMod = this.searchMods.find((mod) => mod.id === currentModId);
      });
    },
    installMod(modVersion) {
      return this.selectedSatisfactoryInstall
        .installMod(modVersion.mod_id, modVersion.version)
        .then(() => {
          this.inProgress.splice(this.inProgress.indexOf(modVersion));
          this.refreshCurrentMod();
        }).catch((err) => {
          this.$bvModal.msgBoxOk(err.toString());
          this.inProgress.splice(this.inProgress.indexOf(modVersion));
        });
    },
    uninstallMod(modVersion) {
      return this.selectedSatisfactoryInstall
        .uninstallMod(modVersion.mod_id)
        .then(() => {
          this.inProgress.splice(this.inProgress.indexOf(modVersion));
          this.refreshCurrentMod();
        }).catch((err) => {
          this.$bvModal.msgBoxOk(err.toString());
          this.inProgress.splice(this.inProgress.indexOf(modVersion));
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
        exec(`start "" "${this.selectedSatisfactoryInstall.launchPath}"`).unref();
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
