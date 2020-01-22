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
        <div class="col-auto justify-content-end my-3 mx-1 flex-grow-0 flex-shrink-0">
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
      <div class="row my-2 flex-fill container-fluid my-2">
        <div
          class="row selection-row"
          style="height: 50%"
        >
          <div
            class="col-7 d-flex flex-column"
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
              v-model="selectedMod"
              :objects="searchMods"
              :can-select="true"
              class="flex-fill"
            >
              <template slot-scope="{item}">
                <div class="col-2">
                  <img
                    :src="item.logo || noImageURL"
                    width="100%"
                  >
                </div>
                <div class="col-auto d-inline-flex align-items-center">
                  <strong>{{ item.name }}</strong>
                </div>
              </template>
            </list>
          </div>
          <div
            class="col-5"
            style="height: 100%"
          >
            <list
              v-if="selectedMod != null"
              :objects="selectedMod.versions"
              :can-select="false"
            >
              <template slot-scope="{item}">
                <div
                  class="col-4"
                  style="min-width: 150px"
                >
                  <button
                    :class="'btn ' + (isModVersionInstalled(item) ? 'btn-secondary' : 'btn-primary')"
                    style="width: 100%"
                    @click="toggleModInstalled(item)"
                  >
                    {{ isModVersionInstalled(item) ? "Remove" : "Install" }}
                  </button>
                </div>
                <div class="col-auto d-inline-flex align-items-center">
                  <strong>{{ item.version }}</strong>
                </div>
                <div class="col-auto">
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
          <div v-html="compiledMarkdownDescription" />
        </div>
      </div>
    </div>
  </main>
</template>

<script>
// TODO: display errors
import semver from 'semver';
import {
  getLatestSMLVersion, getInstalls, getAvailableMods,
} from 'satisfactory-mod-launcher-api';
import marked from 'marked';
import { spawn } from 'child_process';
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
    };
  },
  computed: {
    noImageURL() {
      return 'https://ficsit.app/static/assets/images/no_image.png';
    },
    compiledMarkdownDescription() {
      return marked(this.selectedMod.full_description || '', { sanitize: true });
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
  created() {
    this.refreshSatisfactoryInstalls();
    this.refreshAvailableMods();
    getLatestSMLVersion().then((smlVersion) => {
      this.latestSMLVersion = smlVersion.version;
    });
  },
  methods: {
    refreshSearch() {
      this.searchMods = this.availableMods.filter((mod) => mod.name.toLowerCase().includes(this.search.toLowerCase()));
    },
    refreshAvailableMods() {
      getAvailableMods().then((mods) => {
        this.availableMods = mods;
        this.refreshSearch();
      });
    },
    isModVersionInstalled(modVersion) {
      return this.selectedSatisfactoryInstall.mods[modVersion.mod_id] === modVersion.version;
    },
    refreshCurrentMod() {
      const currentMod = this.selectedMod;
      this.selectedMod = null;
      this.$nextTick().then(() => {
        this.selectedMod = currentMod;
      });
    },
    toggleModInstalled(modVersion) {
      this.inProgress.push(modVersion);
      if (this.isModVersionInstalled(modVersion)) {
        this.selectedSatisfactoryInstall.uninstallMod(modVersion.mod_id).then(() => {
          this.inProgress.splice(this.inProgress.indexOf(modVersion));
          this.refreshCurrentMod();
        });
      } else {
        this.selectedSatisfactoryInstall.installMod(modVersion.mod_id, modVersion.version).then(() => {
          this.inProgress.splice(this.inProgress.indexOf(modVersion));
          this.refreshCurrentMod();
        });
      }
    },
    refreshSatisfactoryInstalls() {
      getInstalls().then((installs) => {
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
      this.selectedSatisfactoryInstall.updateSML().then(() => {
        this.SMLInProgress = false;
      });
    },
    uninstallSML() {
      this.SMLInProgress = true;
      this.selectedSatisfactoryInstall.uninstallSML().then(() => {
        this.SMLInProgress = false;
      });
    },
  },
};
</script>

<style>
</style>
