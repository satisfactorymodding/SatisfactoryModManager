/* eslint-disable no-param-reassign */
import Vue from 'vue';
import Vuex from 'vuex';
import {} from './utils';
import {
  setDebug, addDownloadProgressCallback, getConfigs, loadCache, getInstalls, SatisfactoryInstall, getAvailableSMLVersions, MODS_PER_PAGE, getModsCount,
  getAvailableMods,
  createConfig,
  deleteConfig,
} from 'satisfactory-mod-manager-api';
import { satisfies, coerce, valid } from 'semver';
import { ipcRenderer } from 'electron';
import { saveSetting, getSetting } from './settings';

Vue.use(Vuex);

const MAX_DOWNLOAD_NAME_LENGTH = 25;

function limitDownloadNameLength(name) {
  if (name.length > MAX_DOWNLOAD_NAME_LENGTH) {
    return `${name.substr(0, MAX_DOWNLOAD_NAME_LENGTH - 3)}...`;
  }
  return name;
}

export default new Vuex.Store({
  state: {
    filters: {
      modFilters: {},
      sortBy: '',
      search: '',
    },
    configs: [],
    selectedConfig: {},
    modFilters: [{ name: 'All mods', mods: 0 }, { name: 'Compatible', mods: 0 }, { name: 'Favourite', mods: 0 }, { name: 'Installed', mods: 0 }, { name: 'Not installed', mods: 0 }],
    sortBy: ['Name', 'Last updated', 'Popularity', 'Hotness', 'Views', 'Downloads'],
    satisfactoryInstalls: [],
    selectedInstall: {},
    smlVersions: [],
    mods: [],
    expandedModId: '',
    favoriteModIds: [],
    manifestMods: [],
    inProgress: [], // { id: string, progresses: { id: string, progress: number, message: string, fast: boolean }[] }
    currentDownloadProgress: {},
    error: '',
    isGameRunning: false,
    isLaunchingGame: false,
    expandModInfoOnStart: false,
    debugMode: false,
  },
  mutations: {
    setFilters(state, { newFilters }) {
      state.filters = newFilters;
    },
    setInstall(state, { newInstall }) {
      state.selectedInstall = newInstall;
    },
    setConfig(state, { newConfig }) {
      state.selectedConfig = newConfig;
    },
    setInstalls(state, { installs }) {
      state.satisfactoryInstalls = installs;
    },
    setConfigs(state, { configs }) {
      state.configs = configs;
    },
    setFavoriteModIds(state, { favoriteModIds }) {
      state.favoriteModIds = favoriteModIds;
    },
    setSMLVersions(state, { smlVersions }) {
      state.smlVersions = smlVersions;
    },
    refreshModsInstalledCompatible(state) {
      const installedMods = Object.keys(state.selectedInstall.mods);
      state.manifestMods = state.selectedInstall.manifestMods;
      for (let i = 0; i < state.mods.length; i += 1) {
        state.mods[i].isInstalled = installedMods.includes(state.mods[i].modInfo.mod_reference);
        state.mods[i].isDependency = installedMods.includes(state.mods[i].modInfo.mod_reference) && !state.manifestMods.includes(state.mods[i].modInfo.mod_reference);
        state.mods[i].isCompatible = state.mods[i].modInfo.versions.length > 0
        && !!state.mods[i].modInfo.versions.find((ver) => satisfies(ver.sml_version, '>=2.0.0')
              && state.smlVersions.some((smlVer) => valid(coerce(smlVer.version)) === valid(coerce(ver.sml_version)))
              && satisfies(valid(coerce(state.selectedInstall.version)), `>=${valid(coerce(state.smlVersions.find((smlVer) => valid(coerce(smlVer.version)) === valid(coerce(ver.sml_version))).satisfactory_version))}`));
      }
      state.modFilters[1].mods = state.mods.filter((mod) => mod.isCompatible).length;
      state.modFilters[3].mods = state.mods.filter((mod) => mod.isInstalled).length;
      state.modFilters[4].mods = state.mods.filter((mod) => !mod.isInstalled).length;
    },
    downloadProgress(state, {
      url, progress, name, version,
    }) {
      if (!state.currentDownloadProgress[url]) {
        state.currentDownloadProgress[url] = {
          id: `download_${url}`, progress: 0, message: '', fast: true,
        };
        state.inProgress[0].progresses.push(state.currentDownloadProgress[url]);
      }
      state.currentDownloadProgress[url].message = `Downloading ${limitDownloadNameLength(name)} v${version} ${Math.round(progress.percent * 100)}%`;
      state.currentDownloadProgress[url].progress = progress.percent;
      if (progress.percent === 1) {
        setTimeout(() => {
          state.inProgress[0].progresses.remove(state.currentDownloadProgress[url]);
          delete state.currentDownloadProgress[url];
        }, 100);
      }
    },
    setAvailableMods(state, { mods }) {
      state.mods = mods;
    },
    setExpandedMod(state, { modId }) {
      state.expandedModId = modId;
    },
    clearDownloadProgress(state) {
      state.downloadProgress = [];
    },
    showError(state, { e }) {
      state.error = typeof e === 'string' ? e : e.message;
    },
    launchGame(state) {
      state.isLaunchingGame = true;
      state.isGameRunning = true;
      setTimeout(() => { state.isLaunchingGame = false; }, 10000);
    },
    setGameRunning(state, isGameRunning) {
      state.isGameRunning = isGameRunning;
    },
    setExpandModInfoOnStart(state, value) {
      state.expandModInfoOnStart = value;
    },
    setDebugMode(state, value) {
      state.debugMode = value;
    },
  },
  actions: {
    setFilters({ commit }, newFilters) {
      commit('setFilters', { newFilters });
      saveSetting('filters', { modFilters: newFilters.modFilters.name, sortBy: newFilters.sortBy });
    },
    async selectInstall({ commit, dispatch, state }, newInstall) {
      commit('setInstall', { newInstall });
      if (!state.inProgress.some((prog) => prog.id === '__loadingApp__')) {
        const loadProgress = {
          id: '__loadingApp__',
          progresses: [{
            id: '', progress: -1, message: 'Validating mod install', fast: false,
          }],
        };
        state.inProgress.push(loadProgress);
        const savedConfigName = getSetting('selectedConfig', {})[state.selectedInstall.installLocation] || 'modded';
        commit('setConfig', { newConfig: state.configs.find((conf) => conf.name === savedConfigName) });
        try {
          await newInstall.setConfig(savedConfigName);
          commit('refreshModsInstalledCompatible');
        } catch (e) {
          dispatch('showError', e);
        } finally {
          state.inProgress.remove(loadProgress);
        }
        saveSetting('selectedInstall', newInstall.installLocation);
      }
    },
    async selectConfig({ commit, dispatch, state }, newConfig) {
      commit('setConfig', { newConfig });
      if (!state.inProgress.some((prog) => prog.id === '__loadingApp__')) {
        const loadProgress = {
          id: '__loadingApp__',
          progresses: [{
            id: '', progress: -1, message: 'Validating mod install', fast: false,
          }],
        };
        state.inProgress.push(loadProgress);
        try {
          await state.selectedInstall.setConfig(newConfig.name);
          commit('refreshModsInstalledCompatible');
          let current = getSetting('selectedConfig', {});
          if (typeof current !== 'object') { current = {}; }
          current[state.selectedInstall.installLocation] = state.selectedConfig.name;
          saveSetting('selectedConfig', current);
        } catch (e) {
          dispatch('showError', e);
        } finally {
          state.inProgress.remove(loadProgress);
        }
      }
    },
    async switchModInstalled({ commit, dispatch, state }, modId) {
      if (state.inProgress.length > 0) {
        dispatch('showError', 'Another operation is currently in progress');
        return;
      }
      commit('clearDownloadProgress');
      const modProgress = { id: modId, progresses: [] };
      state.inProgress.push(modProgress);
      const placeholderProgreess = {
        id: 'placeholder', progress: -1, message: '', fast: false,
      };
      modProgress.progresses.push(placeholderProgreess);
      if (state.mods.find((mod) => mod.modInfo.mod_reference === modId).isInstalled) {
        placeholderProgreess.message = 'Checking for mods that are no longer needed';
        try {
          await state.selectedInstall.uninstallMod(modId);
          placeholderProgreess.progress = 1;
          commit('refreshModsInstalledCompatible');
        } catch (e) {
          dispatch('showError', e);
        } finally {
          setTimeout(() => {
            state.inProgress.remove(modProgress);
          }, 500);
        }
      } else {
        placeholderProgreess.message = 'Finding the best version to install';
        try {
          await state.selectedInstall.installMod(modId);
          placeholderProgreess.progress = 1;
          commit('refreshModsInstalledCompatible');
        } catch (e) {
          dispatch('showError', e);
        } finally {
          setTimeout(() => {
            state.inProgress.remove(modProgress);
          }, 500);
        }
      }
    },
    expandMod({ commit }, modId) {
      commit('setExpandedMod', { modId });
      ipcRenderer.send('expand');
    },
    unexpandMod({ commit }) {
      commit('setExpandedMod', { modId: '' });
      ipcRenderer.send('unexpand');
    },
    toggleModFavorite({ state }, modId) {
      if (!state.favoriteModIds.includes(modId)) {
        state.favoriteModIds.push(modId);
      } else {
        state.favoriteModIds.remove(modId);
      }
      state.modFilters[2].mods = state.mods.filter((mod) => state.favoriteModIds.includes(mod.modInfo.mod_reference)).length;
      saveSetting('favoriteMods', state.favoriteModIds);
    },
    createConfig({ dispatch, state }, { configName, copyCurrent }) {
      createConfig(configName, copyCurrent ? state.selectedConfig.name : 'vanilla');
      const newConfig = { name: configName, items: copyCurrent ? state.selectedConfig.items : [] };
      state.configs.push(newConfig);
      dispatch('selectConfig', newConfig);
    },
    deleteConfig({ dispatch, state }, { configName }) {
      deleteConfig(configName);
      state.configs.removeWhere((config) => config.name === configName);
      if (state.selectedConfig.name === configName) {
        dispatch('selectConfig', state.configs.find((config) => config.name === 'modded'));
      }
    },
    async initApp({
      commit, dispatch, state, getters,
    }) {
      const appLoadProgress = {
        id: '__loadingApp__',
        progresses: [{
          id: '', progress: -1, message: 'Loading', fast: false,
        }],
      };
      state.inProgress.push(appLoadProgress);
      dispatch('setDebugMode', getSetting('debugMode', false));
      addDownloadProgressCallback((url, progress, name, version) => commit('downloadProgress', {
        url, progress, name, version,
      }));
      commit('setFavoriteModIds', { favoriteModIds: getSetting('favoriteMods', []) });
      commit('setConfigs', { configs: getConfigs() });
      commit('setExpandModInfoOnStart', getSetting('expandModInfoOnStart', false));

      const savedFilters = getSetting('filters', { modFilters: state.modFilters[0].name, sortBy: state.filters.sortBy[0] });
      commit('setFilters', {
        newFilters: {
          modFilters: state.modFilters.find((modFilter) => modFilter.name === savedFilters.modFilters) || state.modFilters[0],
          sortBy: state.sortBy.find((item) => item === savedFilters.sortBy) || state.sortBy[0],
          search: '',
        },
      });
      try {
        await Promise.all([
          (async () => {
            await loadCache();
            commit('setInstalls', { installs: getInstalls() });
            const installValidateProgress = { id: 'validatingInstall', progress: -1, message: 'Validating mod install' };
            appLoadProgress.progresses.push(installValidateProgress);
            const savedLocation = getSetting('selectedInstall');
            commit('setInstall', { newInstall: state.satisfactoryInstalls.find((install) => install.installLocation === savedLocation) || state.satisfactoryInstalls[0] });
            const savedConfigName = getSetting('selectedConfig', {})[state.selectedInstall.installLocation] || 'modded';
            commit('setConfig', { newConfig: state.configs.find((conf) => conf.name === savedConfigName) });

            await state.selectedInstall.setConfig(savedConfigName);
            appLoadProgress.progresses.remove(installValidateProgress);
          })(),
          (async () => {
            commit('setSMLVersions', { smlVersions: await getAvailableSMLVersions() });
          })(),
          dispatch('getAllMods', { progress: appLoadProgress }),
        ]);
      } catch (e) {
        dispatch('showError', e);
      } finally {
        state.modFilters[0].mods = state.mods.length;
        state.modFilters[2].mods = state.mods.filter((mod) => state.favoriteModIds.includes(mod.modInfo.mod_reference)).length;
        commit('refreshModsInstalledCompatible');
        state.inProgress.remove(appLoadProgress);
        if (state.expandModInfoOnStart) {
          dispatch('expandMod', getters.filteredMods[0].modInfo.mod_reference);
        }
      }
      setInterval(async () => {
        commit('setGameRunning', state.isLaunchingGame || await SatisfactoryInstall.isGameRunning());
      }, 5000);
    },
    async getAllMods({ commit }, { progress }) {
      const getModsProgress = { id: 'getMods', progress: -1, message: 'Getting available mods' };
      if (progress) {
        progress.progresses.push(getModsProgress);
      }
      let modsGot = 0;
      const modCount = await getModsCount();
      getModsProgress.message = `Getting available mods (${modsGot}/${modCount})`;
      getModsProgress.progress = 0;
      const modPages = Math.ceil(modCount / MODS_PER_PAGE);
      const mods = (await Promise.all(Array.from({ length: modPages }).map(async (_, i) => {
        const page = await getAvailableMods(i);
        modsGot += page.length;
        getModsProgress.progress += 1 / modPages;
        getModsProgress.message = `Getting available mods (${modsGot}/${modCount})`;
        return page;
      }))).flat(1);
      commit('setAvailableMods', {
        mods: mods.map((mod) => ({
          modInfo: mod,
          isInstalled: false,
          isCompatible: true,
          isDependency: false,
        })),
      });
      if (progress) {
        progress.progresses.remove(getModsProgress);
      }
    },
    showError({ commit }, e) {
      commit('showError', { e });
      // eslint-disable-next-line no-console
      console.error(e);
    },
    clearError({ commit }) {
      commit('showError', { e: '' });
    },
    setExpandModInfoOnStart({ commit }, value) {
      commit('setExpandModInfoOnStart', value);
      saveSetting('expandModInfoOnStart', value);
    },
    setDebugMode({ commit }, value) {
      setDebug(value);
      commit('setDebugMode', value);
      saveSetting('debugMode', value);
    },
    async updateSingle({ state, commit, dispatch }, update) {
      const updateProgress = {
        id: update.item,
        progresses: [],
      };
      const placeholderProgreess = {
        id: '', progress: -1, message: `Updating ${update.item} to v${update.version}`, fast: false,
      };
      updateProgress.progresses.push(placeholderProgreess);
      state.inProgress.push(updateProgress);
      try {
        await state.selectedInstall.manifestMutate([], [], [update.item]);
        placeholderProgreess.progress = 1;
        commit('refreshModsInstalledCompatible');
      } catch (e) {
        dispatch('showError', e);
      } finally {
        setTimeout(() => {
          state.inProgress.remove(updateProgress);
        }, 500);
      }
    },
    async updateMulti({ state, commit, dispatch }, updates) {
      const updateProgress = {
        id: '__updateMods__',
        progresses: [],
      };
      const placeholderProgreess = {
        id: '', progress: -1, message: `Updating ${updates.length} mods`, fast: false,
      };
      updateProgress.progresses.push(placeholderProgreess);
      state.inProgress.push(updateProgress);
      try {
        await state.selectedInstall.manifestMutate([], [], updates.map((update) => update.item));
        placeholderProgreess.progress = 1;
        commit('refreshModsInstalledCompatible');
      } catch (e) {
        dispatch('showError', e);
      } finally {
        setTimeout(() => {
          state.inProgress.remove(updateProgress);
        }, 500);
      }
    },
  },
  getters: {
    filteredMods(state) {
      let filtered;
      if (state.filters.modFilters === state.modFilters[1]) filtered = state.mods.filter((mod) => mod.isCompatible);
      else if (state.filters.modFilters === state.modFilters[2]) filtered = state.mods.filter((mod) => state.favoriteModIds.includes(mod.modInfo.mod_reference));
      else if (state.filters.modFilters === state.modFilters[3]) filtered = state.mods.filter((mod) => mod.isInstalled);
      else if (state.filters.modFilters === state.modFilters[4]) filtered = state.mods.filter((mod) => !mod.isInstalled);
      else filtered = [...state.mods];

      if (state.filters.search !== '') {
        filtered = filtered.filter((mod) => mod.modInfo.name.toLowerCase().includes(state.filters.search.toLowerCase())); // TODO: maybe search in description too
      }

      if (state.filters.sortBy === 'Name') filtered = filtered.sort((a, b) => a.modInfo.name.localeCompare(b.modInfo.name));
      if (state.filters.sortBy === 'Last updated') filtered = filtered.sort((a, b) => b.modInfo.last_version_date - a.modInfo.last_version_date);
      if (state.filters.sortBy === 'Popularity') filtered = filtered.sort((a, b) => b.modInfo.popularity - a.modInfo.popularity);
      if (state.filters.sortBy === 'Hotness') filtered = filtered.sort((a, b) => b.modInfo.hotness - a.modInfo.hotness);
      if (state.filters.sortBy === 'Views') filtered = filtered.sort((a, b) => b.modInfo.views - a.modInfo.views);
      if (state.filters.sortBy === 'Downloads') filtered = filtered.sort((a, b) => b.modInfo.downloads - a.modInfo.downloads);

      return filtered;
    },
    canInstallMods(state) {
      return state.selectedConfig.name !== 'vanilla' && !state.isGameRunning;
    },
  },
});
