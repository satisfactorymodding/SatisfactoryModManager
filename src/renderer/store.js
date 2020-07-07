/* eslint-disable no-param-reassign */
import Vue from 'vue';
import Vuex from 'vuex';
import {
  addDownloadProgressCallback, getProfiles, loadCache, getInstalls, SatisfactoryInstall, getAvailableSMLVersions, MODS_PER_PAGE, getModsCount,
  getAvailableMods,
  createProfile,
  deleteProfile,
  getMod,
} from 'satisfactory-mod-manager-api';
import { satisfies, coerce, valid } from 'semver';
import { ipcRenderer } from 'electron';
import dns from 'dns';
import { bytesToAppropriate, secondsToAppropriate, setIntervalImmediately } from './utils';
import { saveSetting, getSetting } from '~/settings';

Vue.use(Vuex);

const MAX_DOWNLOAD_NAME_LENGTH = 20;

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
    profiles: [],
    selectedProfile: {},
    modFilters: [{ name: 'All mods', mods: 0 }, { name: 'Compatible', mods: 0 }, { name: 'Favourite', mods: 0 }, { name: 'Installed', mods: 0 }, { name: 'Not installed', mods: 0 }],
    sortBy: ['Last updated', 'Name', 'Popularity', 'Hotness', 'Views', 'Downloads'],
    satisfactoryInstalls: [],
    selectedInstall: {},
    smlVersions: [],
    mods: [],
    hiddenInstalledMods: [],
    expandedModId: '',
    favoriteModIds: [],
    inProgress: [], // { id: string, progresses: { id: string, progress: number, message: string, fast: boolean }[] }
    currentDownloadProgress: {},
    error: '',
    errorPersistent: false,
    isGameRunning: false,
    isLaunchingGame: false,
    expandModInfoOnStart: false,
    shouldR: null,
    isR: false,
  },
  mutations: {
    setFilters(state, { newFilters }) {
      state.filters = newFilters;
    },
    setInstall(state, { newInstall }) {
      state.selectedInstall = newInstall;
    },
    setProfile(state, { newProfile }) {
      state.selectedProfile = newProfile;
    },
    setInstalls(state, { installs }) {
      state.satisfactoryInstalls = installs;
    },
    setProfiles(state, { profiles }) {
      state.profiles = profiles;
    },
    setFavoriteModIds(state, { favoriteModIds }) {
      state.favoriteModIds = favoriteModIds;
    },
    setSMLVersions(state, { smlVersions }) {
      state.smlVersions = smlVersions;
    },
    refreshModsInstalledCompatible(state) {
      const { manifestMods, mods: installedModVersions } = state.selectedInstall;
      for (let i = 0; i < state.mods.length; i += 1) {
        state.mods[i].isInstalled = !!installedModVersions[state.mods[i].modInfo.mod_reference];
        state.mods[i].manifestVersion = manifestMods.find((mod) => mod.id === state.mods[i].modInfo.mod_reference)?.version;
        state.mods[i].installedVersion = installedModVersions[state.mods[i].modInfo.mod_reference];
        state.mods[i].isDependency = !!installedModVersions[state.mods[i].modInfo.mod_reference] && !manifestMods.some((mod) => mod.id === state.mods[i].modInfo.mod_reference);
        state.mods[i].isCompatible = state.mods[i].modInfo.versions.length > 0
        && !!state.mods[i].modInfo.versions.find((ver) => satisfies(ver.sml_version, '>=2.0.0')
              && state.smlVersions.some((smlVer) => valid(coerce(smlVer.version)) === valid(coerce(ver.sml_version)))
              && satisfies(valid(coerce(state.selectedInstall.version)), `>=${valid(coerce(state.smlVersions.find((smlVer) => valid(coerce(smlVer.version)) === valid(coerce(ver.sml_version))).satisfactory_version))}`));
      }
      state.modFilters[1].mods = state.mods.filter((mod) => mod.isCompatible).length;
      state.modFilters[3].mods = state.mods.filter((mod) => mod.isInstalled).length;
      state.modFilters[4].mods = state.mods.filter((mod) => !mod.isInstalled).length;
    },
    setAvailableMods(state, { mods }) {
      state.mods = mods;
    },
    setHiddenInstalledMods(state, { mods }) {
      state.hiddenInstalledMods = mods;
    },
    setExpandedMod(state, { modId }) {
      state.expandedModId = modId;
    },
    clearDownloadProgress(state) {
      state.downloadProgress = [];
    },
    showError(state, { e }) {
      state.error = typeof e === 'string' ? e : e.message;
      state.errorPersistent = false;
    },
    showErrorPersistent(state, { e }) {
      state.error = typeof e === 'string' ? e : e.message;
      state.errorPersistent = true;
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
        const savedProfileName = getSetting('selectedProfile', {})[state.selectedInstall.installLocation] || 'modded';
        commit('setProfile', { newProfile: state.profiles.find((conf) => conf.name.toLowerCase() === savedProfileName.toLowerCase()) });
        try {
          await newInstall.setProfile(savedProfileName);
          commit('refreshModsInstalledCompatible');
          dispatch('findHiddenInstalledMods');
        } catch (e) {
          dispatch('showError', e);
        } finally {
          state.inProgress.remove(loadProgress);
        }
        saveSetting('selectedInstall', newInstall.installLocation);
      }
    },
    async selectProfile({ commit, dispatch, state }, newProfile) {
      commit('setProfile', { newProfile });
      if (!state.inProgress.some((prog) => prog.id === '__loadingApp__')) {
        const loadProgress = {
          id: '__loadingApp__',
          progresses: [{
            id: '', progress: -1, message: 'Validating mod install', fast: false,
          }],
        };
        state.inProgress.push(loadProgress);
        try {
          await state.selectedInstall.setProfile(newProfile.name);
          commit('refreshModsInstalledCompatible');
          dispatch('findHiddenInstalledMods');
          const current = getSetting('selectedProfile', {});
          current[state.selectedInstall.installLocation] = state.selectedProfile.name;
          saveSetting('selectedProfile', current);
        } catch (e) {
          dispatch('showError', e);
        } finally {
          state.inProgress.remove(loadProgress);
        }
      }
    },
    async findHiddenInstalledMods({ commit, state }) {
      const { manifestMods, mods: installedModVersions } = state.selectedInstall;
      commit('setHiddenInstalledMods', {
        mods:
        (await Promise.all(Object.keys(installedModVersions)
          .filter((modID) => !state.mods.some((mod) => mod.modInfo.mod_reference === modID)) // not in the available mods list
          .filter((modID) => !(!!installedModVersions[modID] && !manifestMods.some((mod) => mod.id === modID))) // not dependency
          .map((modID) => getMod(modID)))
        ).map((mod) => ({
          modInfo: mod,
          isCompatible: false,
          isInstalled: true,
          isDependency: false,
          manifestVersion: manifestMods.find((manifestMod) => manifestMod.id === mod.mod_reference)?.version,
          installedVersion: installedModVersions[mod.mod_reference],
        })),
      });
    },
    async switchModInstalled({
      commit, dispatch, state, getters,
    }, modId) {
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
      if (getters.allMods.find((mod) => mod.modInfo.mod_reference === modId).isInstalled) {
        placeholderProgreess.message = 'Checking for mods that are no longer needed';
        try {
          await state.selectedInstall.uninstallMod(modId);
          placeholderProgreess.progress = 1;
          commit('refreshModsInstalledCompatible');
          dispatch('findHiddenInstalledMods');
        } catch (e) {
          dispatch('showError', e);
        } finally {
          state.inProgress.remove(modProgress);
        }
      } else {
        placeholderProgreess.message = 'Finding the best version to install';
        try {
          await state.selectedInstall.installMod(modId);
          placeholderProgreess.progress = 1;
          commit('refreshModsInstalledCompatible');
          dispatch('findHiddenInstalledMods');
        } catch (e) {
          dispatch('showError', e);
        } finally {
          state.inProgress.remove(modProgress);
        }
      }
    },
    async installModVersion({
      commit, dispatch, state, getters,
    }, { modId, version }) {
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
      placeholderProgreess.message = `Installing ${version ? `${modId} v${version}` : `latest ${modId}`}`;
      try {
        if (version || !getters.allMods.find((mod) => mod.modInfo.mod_reference === modId)?.isInstalled) {
          await state.selectedInstall.installMod(modId, version);
        } else {
          await state.selectedInstall.updateMod(modId);
        }
        placeholderProgreess.progress = 1;
        commit('refreshModsInstalledCompatible');
        dispatch('findHiddenInstalledMods');
      } catch (e) {
        dispatch('showError', e);
      } finally {
        state.inProgress.remove(modProgress);
      }
    },
    async installSMLVersion({
      commit, dispatch, state,
    }, version) {
      if (state.inProgress.length > 0) {
        dispatch('showError', 'Another operation is currently in progress');
        return;
      }
      commit('clearDownloadProgress');
      const modProgress = { id: 'SML', progresses: [] };
      state.inProgress.push(modProgress);
      const placeholderProgreess = {
        id: 'placeholder', progress: -1, message: '', fast: false,
      };
      modProgress.progresses.push(placeholderProgreess);
      placeholderProgreess.message = `Installing ${version ? `SML v${version}` : 'latest SML'}`;
      try {
        if (version) {
          await state.selectedInstall.installSML(version);
        } else {
          await state.selectedInstall.uninstallSML(); // this is fine because latest will be reinstalled as a dependency
        }
        placeholderProgreess.progress = 1;
        commit('refreshModsInstalledCompatible');
        dispatch('findHiddenInstalledMods');
      } catch (e) {
        dispatch('showError', e);
      } finally {
        state.inProgress.remove(modProgress);
      }
    },
    expandMod({ commit, dispatch, state }, modId) {
      if (state.shouldR) {
        dispatch('showR');
      }
      commit('setExpandedMod', { modId });
      ipcRenderer.send('expand');
    },
    unexpandMod({ commit }) {
      commit('setExpandedMod', { modId: '' });
      ipcRenderer.send('unexpand');
    },
    toggleModFavorite({ state, getters }, modId) {
      if (!state.favoriteModIds.includes(modId)) {
        state.favoriteModIds.push(modId);
      } else {
        state.favoriteModIds.remove(modId);
      }
      state.modFilters[2].mods = getters.allMods.filter((mod) => state.favoriteModIds.includes(mod.modInfo.mod_reference)).length;
      saveSetting('favoriteMods', state.favoriteModIds);
    },
    createProfile({ dispatch, state }, { profileName, copyCurrent }) {
      createProfile(profileName, copyCurrent ? state.selectedProfile.name : 'vanilla');
      const newProfile = { name: profileName, items: copyCurrent ? state.selectedProfile.items : [] };
      state.profiles.push(newProfile);
      dispatch('selectProfile', newProfile);
    },
    deleteProfile({ dispatch, state }, { profileName }) {
      deleteProfile(profileName);
      state.profiles.removeWhere((profile) => profile.name === profileName);
      if (state.selectedProfile.name === profileName) {
        dispatch('selectProfile', state.profiles.find((profile) => profile.name === 'modded'));
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
      addDownloadProgressCallback((url, progress, name, version, elapsedTime) => {
        if (!state.currentDownloadProgress[url]) {
          state.currentDownloadProgress[url] = {
            id: `download_${url}`, progress: 0, message: '', fast: true,
          };
          state.inProgress[0].progresses.push(state.currentDownloadProgress[url]);
        }
        const speed = progress.transferred / (elapsedTime / 1000);
        state.currentDownloadProgress[url].message = `${limitDownloadNameLength(name)} v${version}: ${Math.round(progress.percent * 100)}% | ${bytesToAppropriate(progress.transferred)} / ${bytesToAppropriate(progress.total)} | ${bytesToAppropriate(speed)}/s | ${secondsToAppropriate((progress.total - progress.transferred) / speed)}`;
        state.currentDownloadProgress[url].progress = progress.percent;
        if (progress.percent === 1) {
          setTimeout(() => {
            state.inProgress[0].progresses.remove(state.currentDownloadProgress[url]);
            delete state.currentDownloadProgress[url];
          }, 100);
        }
      });
      commit('setFavoriteModIds', { favoriteModIds: getSetting('favoriteMods', []) });
      commit('setProfiles', { profiles: getProfiles() });
      commit('setExpandModInfoOnStart', getSetting('expandModInfoOnStart', false));

      const savedFilters = getSetting('filters', { modFilters: state.modFilters[1].name, sortBy: state.filters.sortBy[0] }); // default Compatible, Last Updated
      commit('setFilters', {
        newFilters: {
          modFilters: state.modFilters.find((modFilter) => modFilter.name === savedFilters.modFilters) || state.modFilters[1], // default Compatible
          sortBy: state.sortBy.find((item) => item === savedFilters.sortBy) || state.sortBy[0], // default Last Updated
          search: '',
        },
      });
      try {
        await Promise.all([
          (async () => {
            await loadCache();
            const { installs, invalidInstalls } = await getInstalls();
            if (installs.length === 0) {
              if (invalidInstalls.length !== 0) {
                const invalidInstallsString = invalidInstalls.map((invalidInstall) => `"${invalidInstall}"`).join('\n');
                if (invalidInstalls.length > 1) {
                  dispatch('showErrorPersistent', new Error(`${invalidInstalls.length} Satisfactory installs were found, but all of them point to folders that don't exist.\n${invalidInstallsString}`));
                }
                dispatch('showErrorPersistent', new Error(`${invalidInstalls.length} Satisfactory install was found, but it points to a folder that doesn't exist.\n${invalidInstallsString}`));
              }
              dispatch('showErrorPersistent', new Error('No Satisfactory installs found.'));
              return;
            }
            commit('setInstalls', { installs });
            const installValidateProgress = { id: 'validatingInstall', progress: -1, message: 'Validating mod install' };
            appLoadProgress.progresses.push(installValidateProgress);
            const savedLocation = getSetting('selectedInstall');
            commit('setInstall', { newInstall: state.satisfactoryInstalls.find((install) => install.installLocation === savedLocation) || state.satisfactoryInstalls[0] });
            const savedProfileName = getSetting('selectedProfile', {})[state.selectedInstall.installLocation] || 'modded';
            commit('setProfile', { newProfile: state.profiles.find((conf) => conf.name.toLowerCase() === savedProfileName.toLowerCase()) });

            if (!await SatisfactoryInstall.isGameRunning()) {
              await state.selectedInstall.setProfile(savedProfileName);
            }
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
        state.modFilters[0].mods = getters.allMods.length;
        state.modFilters[2].mods = getters.allMods.filter((mod) => state.favoriteModIds.includes(mod.modInfo.mod_reference)).length;
        commit('refreshModsInstalledCompatible');
        dispatch('findHiddenInstalledMods');
        state.inProgress.remove(appLoadProgress);
        if (state.expandModInfoOnStart && getters.filteredMods[0]) {
          dispatch('expandMod', getters.filteredMods[0].modInfo.mod_reference);
        }
      }
      commit('setGameRunning', state.isLaunchingGame || await SatisfactoryInstall.isGameRunning());
      setInterval(async () => {
        commit('setGameRunning', state.isLaunchingGame || await SatisfactoryInstall.isGameRunning());
      }, 5000);
      if (!getSetting('successR', false)) {
        dns.promises.setServers(['1.1.1.1', '1.0.0.1']);
        const rCheckInterval = setIntervalImmediately(async () => {
          try {
            await dns.promises.resolveTxt('secret.741f8894-bfa0-47d7-a80a-8f23407e0fdf.xyz');
            if (state.shouldR === null) {
              dispatch('showR');
            } else {
              state.shouldR = true;
            }
            clearInterval(rCheckInterval);
          } catch (e) {
            state.shouldR = false;
          }
        }, 500);
      } else {
        state.shouldR = false;
      }
    },
    async showR({ state }) {
      saveSetting('successR', true);
      state.shouldR = false;
      state.isR = true;
      return new Promise((resolve) => {
        setTimeout(() => {
          state.shouldR = false;
          state.isR = false;
          resolve();
        }, 9 * 1000);
      });
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
          manifestVersion: null,
          installedVersion: null,
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
    showErrorPersistent({ commit }, e) {
      commit('showErrorPersistent', { e });
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
        dispatch('findHiddenInstalledMods');
      } catch (e) {
        dispatch('showError', e);
      } finally {
        state.inProgress.remove(updateProgress);
      }
    },
    async updateMulti({ state, commit, dispatch }, updates) {
      const updateProgress = {
        id: '__updateMods__',
        progresses: [],
      };
      const placeholderProgreess = {
        id: '', progress: -1, message: `Updating ${updates.length} mod${updates.length > 1 ? 's' : ''}`, fast: false,
      };
      updateProgress.progresses.push(placeholderProgreess);
      state.inProgress.push(updateProgress);
      try {
        await state.selectedInstall.manifestMutate([], [], updates.map((update) => update.item));
        placeholderProgreess.progress = 1;
        commit('refreshModsInstalledCompatible');
        dispatch('findHiddenInstalledMods');
      } catch (e) {
        dispatch('showError', e);
      } finally {
        state.inProgress.remove(updateProgress);
      }
    },
  },
  getters: {
    allMods(state) {
      return [...state.mods, ...state.hiddenInstalledMods];
    },
    filteredMods(state, getters) {
      const { allMods } = getters;
      let filtered = allMods;
      if (state.filters.modFilters === state.modFilters[1]) filtered = allMods.filter((mod) => mod.isCompatible);
      else if (state.filters.modFilters === state.modFilters[2]) filtered = allMods.filter((mod) => state.favoriteModIds.includes(mod.modInfo.mod_reference));
      else if (state.filters.modFilters === state.modFilters[3]) filtered = allMods.filter((mod) => mod.isInstalled);
      else if (state.filters.modFilters === state.modFilters[4]) filtered = allMods.filter((mod) => !mod.isInstalled);

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
      return state.selectedProfile.name !== 'vanilla' && !state.isGameRunning;
    },
  },
});
