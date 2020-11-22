import Vue from 'vue';
import Vuex from 'vuex';
import {
  addDownloadProgressCallback, getProfiles, loadCache, getInstalls, SatisfactoryInstall,
  createProfile,
  deleteProfile,
  renameProfile,
} from 'satisfactory-mod-manager-api';
import { ipcRenderer } from 'electron';
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
    profiles: [],
    selectedProfile: {},
    satisfactoryInstalls: [],
    selectedInstall: null,
    expandedModId: '',
    favoriteModIds: [],
    inProgress: [], // { id: string, progresses: { id: string, progress: number, message: string, fast: boolean }[] }
    currentDownloadProgress: {},
    installSetupError: '',
    error: '',
    errorPersistent: false,
    isGameRunning: false,
    isLaunchingGame: false,
    expandModInfoOnStart: false,
    installedMods: {},
    manifestMods: {},
  },
  mutations: {
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
    refreshInstalledMods(state) {
      state.manifestMods = state.selectedInstall.readManifest().items.reduce((prev, mod) => Object.assign(prev, { [mod.id]: mod.version || null }), {});
      state.installedMods = state.selectedInstall.readLockfile();
    },
    setExpandedMod(state, { modId }) {
      state.expandedModId = modId;
    },
    clearDownloadProgress(state) {
      state.downloadProgress = [];
    },
    showInstallSetupError(state, { e }) {
      state.installSetupError = typeof e === 'string' ? e : e.message;
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
          commit('refreshInstalledMods');
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
          commit('refreshInstalledMods');
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
    async switchModInstalled({
      commit, dispatch, state,
    }, modId) {
      if (state.inProgress.length > 0) {
        dispatch('showError', `Another operation is currently in progress while trying to (un)install a mod: ${state.inProgress.map((progress) => progress.id)}`);
        return;
      }
      commit('clearDownloadProgress');
      const modProgress = { id: modId, progresses: [] };
      state.inProgress.push(modProgress);
      const placeholderProgreess = {
        id: 'placeholder', progress: -1, message: '', fast: false,
      };
      modProgress.progresses.push(placeholderProgreess);
      if (state.installedMods[modId]) {
        placeholderProgreess.message = 'Checking for mods that are no longer needed';
        try {
          await state.selectedInstall.uninstallMod(modId);
          placeholderProgreess.progress = 1;
          commit('refreshInstalledMods');
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
          commit('refreshInstalledMods');
        } catch (e) {
          dispatch('showError', e);
        } finally {
          state.inProgress.remove(modProgress);
        }
      }
    },
    async installModVersion({
      commit, dispatch, state,
    }, { modId, version }) {
      if (state.inProgress.length > 0) {
        dispatch('showError', `Another operation is currently in progress while trying to install a mod version: ${state.inProgress.map((progress) => progress.id)}`);
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
        if (version || !state.installedMods[modId]) {
          await state.selectedInstall.installMod(modId, version);
        } else {
          await state.selectedInstall.updateMod(modId);
        }
        placeholderProgreess.progress = 1;
        commit('refreshInstalledMods');
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
        dispatch('showError', `Another operation is currently in progress while trying to install a SML version: ${state.inProgress.map((progress) => progress.id)}`);
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
          await state.selectedInstall.updateSML(); // this is fine because latest will be reinstalled as a dependency
        }
        placeholderProgreess.progress = 1;
        commit('refreshInstalledMods');
      } catch (e) {
        dispatch('showError', e);
      } finally {
        state.inProgress.remove(modProgress);
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
    renameProfile({ state }, { newProfile: newName }) {
      const oldName = state.selectedProfile.name;
      renameProfile(oldName, newName);
      const selectedProfile = getSetting('selectedProfile', {});
      Object.keys(selectedProfile).forEach((install) => {
        if (selectedProfile[install] === oldName) {
          selectedProfile[install] = newName;
        }
      });
      saveSetting('selectedProfile', selectedProfile);
      const profile = state.profiles.find((p) => p.name === oldName);
      profile.name = newName;
    },
    async initApp({
      commit, dispatch, state,
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
                } else {
                  dispatch('showErrorPersistent', new Error(`${invalidInstalls.length} Satisfactory install was found, but it points to a folder that doesn't exist.\n${invalidInstallsString}`));
                }
              } else {
                dispatch('showErrorPersistent', new Error('No Satisfactory installs found.'));
              }
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
              try {
                await state.selectedInstall.setProfile(savedProfileName);
              } catch (e) {
                commit('setProfile', { newProfile: state.selectedInstall.profile });
                throw e;
              }
            } else {
              state.selectedInstall._profile = savedProfileName;
            }

            dispatch('setupInstalls');
            appLoadProgress.progresses.remove(installValidateProgress);
          })(),
        ]);
      } catch (e) {
        dispatch('showError', e);
      } finally {
        commit('refreshInstalledMods');
        state.inProgress.remove(appLoadProgress);
      }
      setIntervalImmediately(async () => {
        state.isGameRunning = state.isLaunchingGame || await SatisfactoryInstall.isGameRunning();
      }, 5000);
    },
    async setupInstalls({ dispatch, state }) {
      try {
        await Promise.all(state.satisfactoryInstalls.map((install) => (install.setup ? install.setup() : Promise.resolve())));
      } catch (e) {
        dispatch('showInstallSetupError', e);
      }
    },
    showInstallSetupError({ commit }, e) {
      commit('showInstallSetupError', { e });
      // eslint-disable-next-line no-console
      console.error(e);
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
    clearInstallSetupError({ commit }) {
      commit('showInstallSetupError', { e: '' });
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
        id: '', progress: -1, message: `Updating ${update.name} to v${update.version}`, fast: false,
      };
      updateProgress.progresses.push(placeholderProgreess);
      state.inProgress.push(updateProgress);
      try {
        await state.selectedInstall.manifestMutate([], [], [update.item]);
        placeholderProgreess.progress = 1;
        commit('refreshInstalledMods');
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
        commit('refreshInstalledMods');
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
    canInstallMods(state) {
      return state.selectedProfile.name !== 'vanilla' && !state.isGameRunning;
    },
  },
});
