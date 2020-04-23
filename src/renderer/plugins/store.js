import Vue from 'vue';
import Vuex from 'vuex';
import { ipcRenderer } from 'electron';
import {
  compare, satisfies, coerce, valid,
} from 'semver';
import {
  getInstalls,
  getAvailableMods,
  getConfigs,
  deleteConfig,
  getAvailableSMLVersions,
} from 'satisfactory-mod-manager-api';
import { exec } from 'child_process';
import {
  isModInstalled,
  isModCompatible,
  isModVersionInstalled,
  isVersionSML20Compatible,
  isModHasUpdate,
} from '../utils/helper';
import { getSetting, saveSetting } from '../settings';
import EventBus from './eventBus';

Vue.use(Vuex);

const opt = {
  state: {
    selectedSatisfactoryInstall: null,
    satisfactoryInstalls: [],
    selectedMod: {},
    availableMods: [],
    selectedConfig: '',
    availableConfigs: [],
    configLoadInProgress: false,
    SMLInProgress: false,
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
    updates: [],
    updatingAll: false,
    cachedSMLHasUpdate: false,
    isSidePanelOpen: false,
    sidePanelData: null,
    featureFlag: {
      favorite: false,
  },
  },
  mutations: {
    setSearch(state, search) {
      state.search = search;
      this.commit('refreshSearch');
    },
    refreshSatisfactoryInstalls(state, savedSelectedInstall) {
      return getInstalls().then((installs) => {
        state.satisfactoryInstalls = installs;
        if (state.satisfactoryInstalls.length > 0) {
          if (savedSelectedInstall) {
            state.selectedSatisfactoryInstall = state.satisfactoryInstalls.find((install) => install.installLocation === savedSelectedInstall) || state.satisfactoryInstalls[0];
          } else {
            const defaultInstall = state.satisfactoryInstalls[0];
            state.selectedSatisfactoryInstall = defaultInstall;
          }
        }
      });
    },
    refreshAvailableMods(state) {
      const currentlySelectedModID = state.selectedMod ? state.selectedMod.id : '';
      return getAvailableMods().then((mods) => {
        state.availableMods = mods;
        this.commit('refreshSearch');
        state.selectedMod = state.searchMods.find((mod) => mod.id === currentlySelectedModID) || state.searchMods[0] || null;
      }).catch((err) => {
        EventBus.$emit('popup-alert', 'error', null, err.toString());
      });
    },
    refreshAvailableConfigs(state) {
      const currentlySelectedIdx = state.availableConfigs.indexOf(state.selectedConfig);
      state.availableConfigs = getConfigs();
      state.selectedConfig = state.availableConfigs.includes(state.selectedConfig) ? state.selectedConfig : state.availableConfigs[Math.min(currentlySelectedIdx, state.availableConfigs.length - 1)];
    },
    refreshSearch(state) {
      state.searchMods = state.availableMods.filter((mod) => mod.name.toLowerCase().includes(state.search.toLowerCase())
      && (!state.filters.compatibleOnly || isModCompatible(mod))
      && (state.filters.installedStatus === 'any'
        || (isModInstalled(state, mod) && state.filters.installedStatus === 'installed')
        || (!isModInstalled(state, mod) && state.filters.installedStatus === 'notInstalled')));
      state.searchMods.sort((modA, modB) => {
        switch (state.filters.sortBy) {
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
      if (state.filters.sortOrder === 'ascending') {
        state.searchMods.reverse();
      }
    },
    loadSelectedConfig(state) {
      if (state.selectedSatisfactoryInstall) {
        state.configLoadInProgress = true;
        state.selectedSatisfactoryInstall.loadConfig(state.selectedConfig).then(() => {
          this.commit('refreshAvailableMods');
          state.configLoadInProgress = false;
        }).catch((err) => {
          EventBus.$emit('popup-alert', 'error', null, err.toString());
          state.configLoadInProgress = false;
        });
      }
    },
    refreshCurrentMod(state) {
      const currentModId = state.selectedMod ? state.selectedMod.id : '';
      this.commit('refreshAvailableMods').then(() => {
        state.selectedMod = state.searchMods.find((mod) => mod.id === currentModId) || state.searchMods[0] || null;
      });
    },
  },
  actions: {
    initLauncher({ state, commit, dispatch }) {
      const savedSelectedSFInstall = getSetting('selectedSFInstall', undefined);
      state.selectedConfig = getSetting('selectedConfig', 'modded') || 'vanilla';
      Promise.all(
        [
          commit('refreshSatisfactoryInstalls', savedSelectedSFInstall),
          commit('refreshAvailableMods'),
          commit('refreshAvailableConfigs'),
        ],
      ).then(() => {
        ipcRenderer.send('vue-ready');
        const savedFilters = getSetting('filters', state.filters);
        Object.keys(state.filters).forEach((filter) => {
          if (savedFilters[filter] !== undefined) {
            state.filters[filter] = savedFilters[filter];
          }
        });
        dispatch('checkForUpdates');
      });
    },
    launchSatisfactory({ state }) {
      if (state.selectedSatisfactoryInstall) {
        state.inProgress.push('LAUNCH');
        exec(`start "" "${state.selectedSatisfactoryInstall.launchPath}"`).unref();
        setTimeout(() => {
          state.inProgress.remove('LAUNCH');
        }, 5000);
      }
    },
    handleModalNewConfigSubmit({ state, commit, dispatch }, newConfigName) {
      state.selectedConfig = newConfigName;
      dispatch('saveSelectedConfig').then(() => {
        commit('refreshAvailableConfigs');
      });
    },
    handleModalInstallSubmit(state) {
      this.dispatch.installOldVersion(state.selectedMod, state.modalInstallModVersion);
      this.$nextTick(() => {
        this.$bvModal.hide('modal-install');
      });
    },
    handleModalUninstallSubmit({ state, dispatch }) {
      dispatch('uninstallMod', (state.selectedMod.versions.find((ver) => isModVersionInstalled(ver))));
      this.$nextTick(() => {
        EventBus.$emit('hide-mod-uninstall-confirm');
      });
    },
    saveSelectedConfig({ state }) {
      if (state.selectedSatisfactoryInstall) {
        return state.selectedSatisfactoryInstall.saveConfig(state.selectedConfig).catch((err) => {
          EventBus.$emit('popup-alert', 'error', null, err.toString());
        });
      }
      return Promise.resolve();
    },
    checkForUpdates({ state, dispatch }) {
      state.updates = state.availableMods.filter((mod) => isModHasUpdate(state, mod)).map((mod) => ({ id: mod.id, version: mod.versions[0].version }));
      dispatch('hasSMLUpdate').then((hasUpdate) => {
        getAvailableSMLVersions().then((versions) => versions.sort((a, b) => -compare(a.version, b.version))[0].version).then((latestSMLVersion) => {
          if (hasUpdate) {
            state.updates.push({ id: 'SML', version: latestSMLVersion });
          }
          const ignoredUpdates = getSetting('ignoredUpdates', []);
          state.updates = state.updates.filter((update) => !ignoredUpdates.find((ignored) => ignored.id === update.id && ignored.version === update.version));
          if (state.updates.length > 0) {
            EventBus.$emit('show-mod-update-dialog');
          }
        });
      });
    },
    hasSMLUpdate({ state }) {
      return getAvailableSMLVersions()
        .then((versions) => {
          const compatibleVersions = versions
            .filter((ver) => satisfies(valid(coerce(state.selectedSatisfactoryInstall.version)), `>=${ver.satisfactory_version}`));
          state.cachedSMLHasUpdate = compatibleVersions.length > 0
            && state.selectedSatisfactoryInstall.smlVersion
            && state.selectedSatisfactoryInstall.smlVersion !== compatibleVersions.sort((a, b) => -compare(a.version, b.version))[0].version;
          return state.cachedSMLHasUpdate;
        });
    },
    updateSML({ state, dispatch }) {
      state.inProgress.push('SML');
      return state.selectedSatisfactoryInstall
        .updateSML()
        .then(() => {
          dispatch('saveSelectedConfig').then(() => {
            state.cachedSMLHasUpdate = false;
            state.inProgress.remove('SML');
          });
        }).catch((err) => {
          EventBus.$emit('popup-alert', 'error', null, err.toString());
          state.inProgress.remove('SML');
        });
    },
    installUninstallUpdate({ state, dispatch }, mod) {
      if (state.selectedConfig !== 'vanilla') {
        if (state.inProgress.length === 0) {
          state.inProgress.push(mod);
          if (isModInstalled(state, mod)) {
            if (isModHasUpdate(state, mod)) {
              dispatch('updateMod', mod);
            } else {
              dispatch('uninstallMod', mod);
            }
          } else {
            dispatch('installMod', mod);
          }
        } else {
          EventBus.$emit('popup-alert', 'warning', null, 'Another operation is currently in progress. Wait for it to finish.');
        }
      } else {
        const defaultModdedExists = state.availableConfigs.includes('modded');
        const hasOtherConfigs = state.availableConfigs.length > (defaultModdedExists ? 2 : 1);
        let message = 'Cannot modify the vanilla config. Create a new config to be able to install mods.';
        if (defaultModdedExists || hasOtherConfigs) {
          message = `Cannot modify the vanilla config. Choose ${defaultModdedExists ? 'the modded config' : ''}${defaultModdedExists && hasOtherConfigs ? ' or ' : ''}${hasOtherConfigs ? 'one of your custom configs' : ''}`;
        }
        EventBus.$emit('popup-alert', 'warning', null, message);
      }
    },
    updateAll({ state, dispatch }) {
      if (state.updates.length === 0) {
        state.updatingAll = false;
        this.$nextTick(() => {
          EventBus.$emit('hide-mod-update-dialog');
        });
      } else {
        state.updatingAll = true;
        dispatch('updateById', this.updates[0].id)
          .then(() => {
            dispatch('updateAll');
          });
      }
    },
    updateById({ state, dispatch }, id) {
      state.inProgress.push(id);
      return (id === 'SML' ? dispatch('updateSML') : dispatch('updateMod', (state.availableMods.find((mod) => mod.id === id)))).then(() => {
        state.updates.removeWhere((update) => update.id === id);
        state.inProgress.remove(id);
        if (state.updates.length === 0) {
          this.$nextTick(() => {
            EventBus.$emit('hide-mod-update-dialog');
          });
        }
      }).catch(() => this.inProgress.remove(id));
    },
    ignoreVersion({ state }, update) {
      state.updates.remove(update);
      const ignoredUpdates = getSetting('ignoredUpdates', []);
      ignoredUpdates.push(update);
      saveSetting('ignoredUpdates', ignoredUpdates);
      if (state.updates.length === 0) {
        this.$nextTick(() => {
          EventBus.$emit('hide-mod-update-dialog');
        });
      }
    },
    installOldVersion({ state, commit, dispatch }, mod, version) {
      state.inProgress.push(mod);
      return state.selectedSatisfactoryInstall
        .installMod(mod.id, version.version)
        .then(() => {
          dispatch('saveSelectedConfig').then(() => {
            state.inProgress.remove(mod);
            commit('refreshCurrentMod');
          });
        }).catch((err) => {
          EventBus.$emit('popup-alert', 'error', null, err.toString());
          state.inProgress.remove(mod);
        });
    },
    installMod({ state, commit, dispatch }, mod) {
      return state.selectedSatisfactoryInstall
        .installMod(mod.id)
        .then(() => {
          dispatch('saveSelectedConfig').then(() => {
            state.inProgress.remove(mod);
            commit('refreshCurrentMod');
          });
        }).catch((err) => {
          EventBus.$emit('popup-alert', 'error', null, err.toString());
          state.inProgress.remove(mod);
        });
    },
    updateMod({ state, commit, dispatch }, mod) {
      return state.selectedSatisfactoryInstall
        .updateMod(mod.id)
        .then(() => {
          dispatch('saveSelectedConfig').then(() => {
            state.inProgress.remove(mod);
            if (state.updates.find((update) => mod.id === update.id)) {
              state.updates.remove(state.updates.find((update) => mod.id === update.id));
            }
            commit('refreshCurrentMod');
          });
        }).catch((err) => {
          EventBus.$emit('popup-alert', 'error', null, err.toString());
          state.inProgress.remove(mod);
        });
    },
    uninstallMod({ state, commit, dispatch }, mod) {
      return state.selectedSatisfactoryInstall
        .uninstallMod(mod.id)
        .then(() => {
          dispatch('saveSelectedConfig').then(() => {
            state.inProgress.remove(mod);
            commit('refreshCurrentMod');
          });
        }).catch((err) => {
          EventBus.$emit('popup-alert', 'error', null, err.toString());
          state.inProgress.remove(mod);
        });
    },
    deleteSelectedConfig({ state, commit }) {
      try {
        deleteConfig(state.selectedConfig);
      } catch (err) {
        EventBus.$emit('popup-alert', 'error', null, err.toString());
      }
      commit('refreshAvailableConfigs');
    },
    openedByUrl({ state }, url) {
      const parsed = new URL(url);
      const command = parsed.pathname.replace(/^\/+|\/+$/g, '');
      if (command === 'install') {
        const modID = parsed.searchParams.get('modID');
        const version = parsed.searchParams.get('version');
        state.selectedMod = state.availableMods.find((mod) => mod.id === modID);
        state.modalInstallModVersion = state.selectedMod.versions.filter((ver) => isVersionSML20Compatible(ver)).find((ver) => ver.version === version) || state.selectedMod.versions[0];
        // this.$bvModal.show('modal-install');
      } else if (command === 'uninstall') {
        const modID = parsed.searchParams.get('modID');
        state.selectedMod = state.availableMods.find((mod) => mod.id === modID);
        // this.$bvModal.show('modal-uninstall');
      }
    },
  },
};

export default new Vuex.Store(opt);
