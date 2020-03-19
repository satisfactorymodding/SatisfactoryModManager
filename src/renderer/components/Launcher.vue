<template>
  <main>
    <div class="container-fluid my-2 content d-flex flex-column">
      <div class="row flex-grow-0 flex-shrink-0">
        <div class="col-5">
          <select
            v-model="selectedSatisfactoryInstall"
            class="form-control"
            :disabled="inProgress.length > 0 || configLoadInProgress"
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
          <button
            class="btn btn-primary"
            :disabled="inProgress.length > 0 || configLoadInProgress"
            @click="launchSatisfactory"
          >
            Launch Satisfactory
          </button>
        </div>
        <div class="col-auto d-inline-flex align-items-center">
          <strong>SML: {{ selectedSatisfactoryInstall ? (selectedSatisfactoryInstall.smlVersion || 'Install a mod to install SML') : 'Select a Satisfactory Install' }}</strong>
        </div>
        <div class="col-1 d-inline-flex align-items-center">
          <div
            v-if="SMLInProgress"
            class="spinner-border"
            role="status"
          >
            <span class="sr-only">Loading...</span>
          </div>
        </div>
      </div>
      <div class="row flex-grow-0 flex-shrink-0 my-2">
        <div class="col-1 d-inline-flex align-items-center">
          Configs:
        </div>
        <div class="col-4">
          <select
            v-model="selectedConfig"
            class="form-control"
            :disabled="inProgress.length > 0 || configLoadInProgress"
          >
            <option
              v-for="config in availableConfigs"
              :key="config"
              :value="config"
            >
              {{ config }}
            </option>
          </select>
        </div>
        <div class="col-1 d-inline-flex align-items-center">
          <div
            v-if="configLoadInProgress"
            class="spinner-border"
            role="status"
          >
            <span class="sr-only">Loading...</span>
          </div>
        </div>
        <div class="col-auto">
          <button
            class="btn btn-primary"
            :disabled="inProgress.length > 0 || configLoadInProgress"
            @click="deleteSelectedConfig"
          >
            Delete
          </button>
        </div>
        <div class="col-auto">
          <button
            class="btn btn-primary"
            :disabled="inProgress.length > 0 || configLoadInProgress"
            @click="newConfig"
          >
            New
          </button>
        </div>
      </div>
      <div
        class="row my-2 flex-fill container-fluid my-2"
        style="font-size: 14px; width: 100%"
      >
        <div
          class="row selection-row"
          style="height: 50%; width: 100%"
        >
          <div
            class="col-auto d-flex flex-column"
            style="height: 100%; width: 100%"
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
              <template
                slot-scope="{item}"
              >
                <div
                  class="col-1 p-0"
                  style="flex: 0 0 7%; max-width: 7%;"
                  :style="!isModSML20Compatible(item) ? 'background-color: #837971' : ''"
                >
                  <img
                    :src="item.logo || noImageURL"
                    width="50px"
                  >
                </div>
                <div
                  class="d-inline-flex align-items-center text-break"
                  style="flex: 0 0 15%; max-width: 15%;"
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
                  class="col-2 d-inline-flex align-items-center"
                  :style="!isModSML20Compatible(item) ? 'background-color: #837971' : ''"
                >
                  <strong>{{ item.authors.map((author) => author.user.username).join(', ') }}</strong>
                </div>
                <div
                  class="d-inline-flex align-items-center"
                  style="flex: 0 0 10%; max-width: 10%;"
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
                    :disabled="!item.versions[0] || !isModSML20Compatible(item) || inProgress.length > 0 || configLoadInProgress || selectedConfig === 'vanilla'"
                    @click="installUninstallUpdate(item)"
                  >
                    {{ !item.versions[0] ? 'N/A' : (isModSML20Compatible(item) ? (isModVersionInstalled(item.versions[0]) ? "Remove" : (isModInstalled(item) ? "Update" : "Install")) : 'Outdated') }}
                  </button>
                </div>
                <div
                  class="d-inline-flex align-items-center"
                  style="flex: 0 0 15%; max-width: 15%;"
                  :style="!isModSML20Compatible(item) ? 'background-color: #837971' : ''"
                >
                  <button
                    v-if="!isModInstalled(item)"
                    :class="'my-1 btn btn-secondary'"
                    style="font-size: 13px; width: 100%"
                    :disabled="!item.versions[0] || !isModSML20Compatible(item) || inProgress.length > 0 || configLoadInProgress || selectedConfig === 'vanilla'"
                    @click="$bvModal.show('modal-install')"
                  >
                    {{ !item.versions[0] ? 'N/A' : (isModSML20Compatible(item) ? 'Install old version' : 'Outdated') }}
                  </button>
                </div>
                <div
                  v-if="inProgress.includes(item)"
                  class="col-1 d-inline-flex align-items-center"
                >
                  <div
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
        v-if="selectedMod"
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
            v-for="version in selectedMod.versions ? selectedMod.versions.filter((ver) => isVersionSML20Compatible(ver)) : []"
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
        v-if="selectedMod"
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
      title="Filter Mods"
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
        <label for="sortBySelect">Show mods:</label>
        <select
          id="filterInstalledStatus"
          v-model="filters.installedStatus"
          class="form-control"
        >
          <option
            v-for="installedStatusOption in installedStatusOptions"
            :key="installedStatusOption.value"
            :value="installedStatusOption.value"
          >
            {{ installedStatusOption.displayName }}
          </option>
        </select>
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
    <b-modal
      id="modal-new-config"
      title="Add new config"
      @ok="handleModalNewConfigOk"
    >
      <form
        ref="newConfigForm"
        @submit.stop.prevent="handleModalNewConfigSubmit"
      >
        <b-form-group
          :state="newConfigNameState"
          label="Config name"
          label-for="config-name-input"
          invalid-feedback="Config name is required"
        >
          <b-form-input
            id="config-name-input"
            v-model="newConfigName"
            :state="newConfigNameState"
            required
          />
        </b-form-group>
      </form>
    </b-modal>
    <b-modal
      id="modal-update-available"
      :title="`New update available: ${availableUpdate ? availableUpdate.version : ''}`"
      ok-only
    >
      <p>Update will be installed when the app closes.</p>
      <!-- eslint-disable-next-line vue/no-v-html -->
      <div v-html="availableUpdate ? availableUpdate.releaseNotes : ''" />
    </b-modal>
  </main>
</template>

<script>
import semver from 'semver';
import {
  getInstalls,
  getAvailableMods,
  toggleDebug,
  clearCache,
  getConfigs,
  deleteConfig,
} from 'satisfactory-mod-manager-api';
import marked from 'marked';
import { exec } from 'child_process';
import sanitizeHtml from 'sanitize-html';
import List from './List';
import { getSetting, saveSetting } from '../settings';

export default {
  name: 'Launcher',
  components: {
    List,
  },
  data() {
    return {
      selectedSatisfactoryInstall: null,
      satisfactoryInstalls: [],
      selectedMod: {},
      availableMods: [],
      selectedConfig: '',
      availableConfigs: [],
      newConfigName: '',
      newConfigNameState: null,
      configLoadInProgress: false,
      SMLInProgress: false,
      devSML: false,
      searchMods: [],
      search: '',
      filters: {
        compatibleOnly: true,
        installedStatus: 'any',
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
      installedStatusOptions: [
        {
          value: 'installed',
          displayName: 'Installed',
        },
        {
          value: 'notInstalled',
          displayName: 'Not installed',
        },
        {
          value: 'any',
          displayName: 'All',
        },
      ],
      availableUpdate: null,
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
  },
  watch: {
    search() {
      this.refreshSearch();
    },
    filters: {
      handler() {
        this.refreshSearch();
        saveSetting('filters', this.filters);
      },
      deep: true,
    },
    selectedSatisfactoryInstall() {
      saveSetting('selectedSFInstall', this.selectedSatisfactoryInstall.installLocation);
      this.checkDevSML();
    },
    selectedConfig() {
      saveSetting('selectedConfig', this.selectedConfig);
      this.loadSelectedConfig();
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
        this.modalInstallModVersion = this.selectedMod.versions.filter((ver) => this.isVersionSML20Compatible(ver)).find((ver) => ver.version === version) || this.selectedMod.versions[0];
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
    this.$electron.ipcRenderer.on('toggleDevSML', () => {
      this.toggleDevSML();
    });
    this.$electron.ipcRenderer.on('update-available', (e, updateInfo) => {
      this.availableUpdate = updateInfo;
      this.$bvModal.show('modal-update-available');
    });
  },
  created() {
    const savedSelectedSFInstall = getSetting('selectedSFInstall', undefined);
    this.selectedConfig = getSetting('selectedConfig', 'modded') || 'vanilla';
    this.devSML = getSetting('devSML', false);
    Promise.all(
      [
        this.refreshSatisfactoryInstalls(savedSelectedSFInstall),
        this.refreshAvailableMods(),
        this.refreshAvailableConfigs(),
      ],
    ).then(() => {
      this.$electron.ipcRenderer.send('vue-ready');
      const savedFilters = getSetting('filters', this.filters);
      Object.keys(this.filters).forEach((filter) => {
        if (savedFilters[filter]) {
          this.filters[filter] = savedFilters[filter];
        }
      });
    });
  },
  methods: {
    handleModalInstallOk(bvModalEvt) {
      bvModalEvt.preventDefault();
      this.handleModalInstallSubmit();
    },
    handleModalInstallSubmit() {
      this.installOldVersion(this.selectedMod, this.modalInstallModVersion);
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
        && (!this.filters.compatibleOnly || this.isModSML20Compatible(mod))
        && (this.filters.installedStatus === 'any'
          || (this.isModInstalled(mod) && this.filters.installedStatus === 'installed')
          || (!this.isModInstalled(mod) && this.filters.installedStatus === 'notInstalled')));
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
            if (modB.last_version_date) {
              return 1;
            }
            if (modA.last_version_date) {
              return -1;
            }
            return 0;
        }
      });
      if (this.filters.sortOrder === 'ascending') {
        this.searchMods.reverse();
      }
    },
    saveSelectedConfig() {
      if (this.selectedSatisfactoryInstall) {
        return this.selectedSatisfactoryInstall.saveConfig(this.selectedConfig).catch((err) => {
          this.$bvModal.msgBoxOk(err.toString());
        });
      }
      return Promise.resolve();
    },
    loadSelectedConfig() {
      if (this.selectedSatisfactoryInstall) {
        this.configLoadInProgress = true;
        this.selectedSatisfactoryInstall.loadConfig(this.selectedConfig).then(() => {
          this.refreshAvailableMods();
          this.checkDevSML();
          this.configLoadInProgress = false;
        }).catch((err) => {
          this.$bvModal.msgBoxOk(err.toString());
          this.configLoadInProgress = false;
        });
      }
    },
    deleteSelectedConfig() {
      try {
        deleteConfig(this.selectedConfig);
      } catch (err) {
        this.$bvModal.msgBoxOk(err.toString());
      }
      this.refreshAvailableConfigs();
    },
    newConfig() {
      this.$bvModal.show('modal-new-config');
    },
    checkNewConfigFormValidity() {
      const valid = this.$refs.newConfigForm.checkValidity();
      this.newConfigNameState = valid;
      return valid;
    },
    handleModalNewConfigOk(bvModalEvt) {
      bvModalEvt.preventDefault();
      this.handleModalNewConfigSubmit();
    },
    handleModalNewConfigSubmit() {
      if (!this.checkNewConfigFormValidity()) {
        return;
      }
      this.selectedConfig = this.newConfigName;
      this.newConfigName = '';
      this.newConfigNameState = null;
      this.saveSelectedConfig();
      this.refreshAvailableConfigs();
      this.$nextTick(() => {
        this.$bvModal.hide('modal-new-config');
      });
    },
    refreshAvailableMods() {
      const currentlySelectedModID = this.selectedMod.id;
      return getAvailableMods().then((mods) => {
        this.availableMods = mods;
        this.refreshSearch();
        this.selectedMod = this.searchMods.find((mod) => mod.id === currentlySelectedModID) || this.searchMods[0] || null;
      });
    },
    refreshAvailableConfigs() {
      const currentlySelectedIdx = this.availableConfigs.indexOf(this.selectedConfig);
      this.availableConfigs = getConfigs();
      this.selectedConfig = this.availableConfigs.includes(this.selectedConfig) ? this.selectedConfig : this.availableConfigs[Math.min(currentlySelectedIdx, this.availableConfigs.length - 1)];
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
    installOldVersion(mod, version) {
      this.inProgress.push(mod);
      return this.selectedSatisfactoryInstall
        .installMod(mod.id, version.version)
        .then(() => {
          this.saveSelectedConfig().then(() => {
            this.inProgress.splice(this.inProgress.indexOf(mod));
            this.refreshCurrentMod();
          });
        }).catch((err) => {
          this.$bvModal.msgBoxOk(err.toString());
          this.inProgress.splice(this.inProgress.indexOf(mod));
        });
    },
    installMod(mod) {
      return this.selectedSatisfactoryInstall
        .installMod(mod.id)
        .then(() => {
          this.saveSelectedConfig().then(() => {
            this.inProgress.splice(this.inProgress.indexOf(mod));
            this.refreshCurrentMod();
          });
        }).catch((err) => {
          this.$bvModal.msgBoxOk(err.toString());
          this.inProgress.splice(this.inProgress.indexOf(mod));
        });
    },
    updateMod(mod) {
      return this.selectedSatisfactoryInstall
        .updateMod(mod.id)
        .then(() => {
          this.saveSelectedConfig().then(() => {
            this.inProgress.splice(this.inProgress.indexOf(mod));
            this.refreshCurrentMod();
          });
        }).catch((err) => {
          this.$bvModal.msgBoxOk(err.toString());
          this.inProgress.splice(this.inProgress.indexOf(mod));
        });
    },
    uninstallMod(mod) {
      return this.selectedSatisfactoryInstall
        .uninstallMod(mod.id)
        .then(() => {
          this.saveSelectedConfig().then(() => {
            this.inProgress.splice(this.inProgress.indexOf(mod));
            this.refreshCurrentMod();
          });
        }).catch((err) => {
          this.$bvModal.msgBoxOk(err.toString());
          this.inProgress.splice(this.inProgress.indexOf(mod));
        });
    },
    installUninstallUpdate(mod) {
      if (this.selectedConfig !== 'vanilla') {
        if (this.inProgress.length === 0) {
          this.inProgress.push(mod);
          if (this.isModInstalled(mod)) {
            if (this.isModVersionInstalled(mod.versions[0])) {
              this.uninstallMod(mod).then(() => {
                this.checkDevSML();
              });
            } else {
              this.updateMod(mod);
            }
          } else {
            this.installMod(mod);
          }
        } else {
          this.$bvModal.msgBoxOk('Another operation is currently in progress. Wait for it to finish.');
        }
      } else {
        const defaultModdedExists = this.availableConfigs.includes('modded');
        const hasOtherConfigs = this.availableConfigs.length > (defaultModdedExists ? 2 : 1);
        if (defaultModdedExists || hasOtherConfigs) {
          this.$bvModal.msgBoxOk(`Cannot modify the vanilla config. Choose ${defaultModdedExists ? 'the modded config' : ''}${defaultModdedExists && hasOtherConfigs ? ' or ' : ''}${hasOtherConfigs ? 'one of your custom configs' : ''}`);
        } else {
          this.$bvModal.msgBoxOk('Cannot modify the vanilla config. Create a new config to be able to install mods.');
        }
      }
    },
    refreshSatisfactoryInstalls(savedSelectedInstall) {
      return getInstalls().then((installs) => {
        this.satisfactoryInstalls = installs;
        if (this.satisfactoryInstalls.length > 0) {
          if (savedSelectedInstall) {
            this.selectedSatisfactoryInstall = this.satisfactoryInstalls.find((install) => install.installLocation === savedSelectedInstall) || this.satisfactoryInstalls[0];
          } else {
            const defaultInstall = this.satisfactoryInstalls[0];
            this.selectedSatisfactoryInstall = defaultInstall;
          }
          this.checkDevSML();
        }
      });
    },
    launchSatisfactory() {
      if (this.selectedSatisfactoryInstall) {
        exec(`start "" "${this.selectedSatisfactoryInstall.launchPath}"`).unref();
      }
    },
    checkDevSML() {
      if (this.selectedSatisfactoryInstall) {
        if (this.devSML) {
          if (!this.selectedSatisfactoryInstall.smlVersion) {
            this.SMLInProgress = true;
            return this.selectedSatisfactoryInstall.installSML().then(() => {
              this.SMLInProgress = false;
            });
          }
        }
      }
      return Promise.resolve();
    },
    toggleDevSML() {
      this.devSML = !this.devSML;
      if (!this.devSML) {
        this.SMLInProgress = true;
        return this.selectedSatisfactoryInstall.uninstallSML().then(() => {
          this.SMLInProgress = false;
        });
      }
      return this.checkDevSML();
    },
  },
};
</script>

<style>
</style>
