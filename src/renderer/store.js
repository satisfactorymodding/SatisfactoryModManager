import Vue from 'vue';
import Vuex from 'vuex';
import {
  addDownloadProgressCallback, getProfiles, loadCache, getInstalls, SatisfactoryInstall,
  getProfileFolderPath, readManifest, readLockfile,
  createProfile,
  deleteProfile,
  renameProfile,
  clearOutdatedCache,
  LogLevel,
} from 'satisfactory-mod-manager-api';
import path from 'path';
import { ipcRenderer } from 'electron';
import { saveSetting, getSetting } from '~/settings';
import { bytesToAppropriate, secondsToAppropriate, setIntervalImmediately } from './utils';
import { write as writeLog } from './logging';

Vue.use(Vuex);

const DOWNLOAD_SPEED_AVG_TIME = 10 * 1000;
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
    modsEnabled: true,
    satisfactoryInstalls: [],
    selectedInstall: null,
    expandedModId: '',
    favoriteModIds: [],
    inProgress: [], // { id: string, progresses: { id: string, progress: number, message: string, fast: boolean }[] }
    downloadProgress: {},
    allDownloadProgress: null,
    installSetupError: '',
    error: '',
    errorPersistent: false,
    isGameRunning: false,
    isLaunchingGame: false,
    expandModInfoOnStart: false,
    installedMods: {},
    manifestItems: [],
    konami: false,
    launchButton: false,
    launchCat: false,
  },
  mutations: {
    setInstall(state, { newInstall }) {
      state.selectedInstall = newInstall;
    },
    setProfile(state, { newProfile }) {
      state.selectedProfile = newProfile;
    },
    setModsEnabled(state, modsEnabled) {
      state.modsEnabled = modsEnabled;
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
      state.manifestItems = readManifest(path.join(getProfileFolderPath(state.selectedProfile.name), 'manifest.json')).items;
      state.installedMods = readLockfile(path.join(getProfileFolderPath(state.selectedProfile.name), state.selectedInstall.lockfileName));
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
      writeLog(LogLevel.ERROR, e);
    },
    showErrorPersistent(state, { e }) {
      state.error = typeof e === 'string' ? e : e.message;
      state.errorPersistent = true;
      writeLog(LogLevel.ERROR, e);
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
    konami(state) {
      state.konami = true;
    },
    launchButton(state, value) {
      state.launchButton = value;
    },
    launchCat(state, value) {
      state.launchCat = value;
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
        let savedProfileName = getSetting('selectedProfile', {})[state.selectedInstall.installLocation] || 'modded';
        if (!state.profiles.some((profile) => profile.name.toLowerCase() === savedProfileName.toLowerCase())) {
          savedProfileName = 'modded'; // If profile is missing, default to modded
        }
        commit('setProfile', { newProfile: state.profiles.find((profile) => profile.name.toLowerCase() === savedProfileName.toLowerCase()) });
        let savedModsEnabled = getSetting('modsEnabled', {})[state.selectedInstall.installLocation];
        if (savedModsEnabled === undefined) {
          savedModsEnabled = true;
        }
        commit('setModsEnabled', savedModsEnabled);
        try {
          if (state.modsEnabled) {
            await newInstall.setProfile(savedProfileName);
          } else {
            await newInstall.setProfile('vanilla');
          }
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
        if (state.modsEnabled) {
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
        } else {
          commit('refreshInstalledMods');
        }
      }
    },
    async setModsEnabled({ commit, dispatch, state }, modsEnabled) {
      commit('setModsEnabled', modsEnabled);
      if (!state.inProgress.some((prog) => prog.id === '__loadingApp__')) {
        const loadProgress = {
          id: '__loadingApp__',
          progresses: [{
            id: '', progress: -1, message: 'Validating mod install', fast: false,
          }],
        };
        state.inProgress.push(loadProgress);
        try {
          if (state.modsEnabled) {
            await state.selectedInstall.setProfile(state.selectedProfile.name);
          } else {
            await state.selectedInstall.setProfile('vanilla');
          }
          commit('refreshInstalledMods');
          const current = getSetting('modsEnabled', {});
          current[state.selectedInstall.installLocation] = state.modsEnabled;
          saveSetting('modsEnabled', current);
        } catch (e) {
          dispatch('showError', e);
        } finally {
          state.inProgress.remove(loadProgress);
        }
      }
    },
    async installMod({
      commit, dispatch, state,
    }, modId) {
      if (!state.manifestItems.some((item) => item.id === modId)) {
        if (state.inProgress.length > 0) {
          dispatch('showError', `Another operation is currently in progress while trying to install a mod: ${state.inProgress.map((progress) => progress.id)}`);
          return;
        }
        if (!state.modsEnabled) {
          dispatch('showError', 'Enable mods to be able to make changes');
          return;
        }
        commit('clearDownloadProgress');
        const modProgress = { id: modId, progresses: [] };
        state.inProgress.push(modProgress);
        const placeholderProgreess = {
          id: 'placeholder', progress: -1, message: '', fast: false,
        };
        modProgress.progresses.push(placeholderProgreess);
        placeholderProgreess.message = 'Finding the best version to install';
        try {
          await state.selectedInstall.installMod(modId);
          placeholderProgreess.progress = 1;
          commit('refreshInstalledMods');
        } catch (e) {
          dispatch('showError', e);
        } finally {
          // Allow the UI to update properly
          setTimeout(() => {
            state.inProgress.remove(modProgress);
          }, 500);
        }
      }
    },
    async uninstallMod({
      commit, dispatch, state,
    }, modId) {
      if (state.manifestItems.some((item) => item.id === modId)) {
        if (state.inProgress.length > 0) {
          dispatch('showError', `Another operation is currently in progress while trying to uninstall a mod: ${state.inProgress.map((progress) => progress.id)}`);
          return;
        }
        if (!state.modsEnabled) {
          dispatch('showError', 'Enable mods to be able to make changes');
          return;
        }
        commit('clearDownloadProgress');
        const modProgress = { id: modId, progresses: [] };
        state.inProgress.push(modProgress);
        const placeholderProgreess = {
          id: 'placeholder', progress: -1, message: '', fast: false,
        };
        modProgress.progresses.push(placeholderProgreess);
        placeholderProgreess.message = 'Checking for mods that are no longer needed';
        try {
          await state.selectedInstall.uninstallMod(modId);
          placeholderProgreess.progress = 1;
          commit('refreshInstalledMods');
        } catch (e) {
          dispatch('showError', e);
        } finally {
          // Allow the UI to update properly
          setTimeout(() => {
            state.inProgress.remove(modProgress);
          }, 500);
        }
      }
    },
    async enableMod({
      commit, dispatch, state,
    }, modId) {
      if (state.manifestItems.some((item) => item.id === modId)) {
        if (state.inProgress.length > 0) {
          dispatch('showError', `Another operation is currently in progress while trying to enable a mod: ${state.inProgress.map((progress) => progress.id)}`);
          return;
        }
        if (!state.modsEnabled) {
          dispatch('showError', 'Enable mods to be able to make changes');
          return;
        }
        commit('clearDownloadProgress');
        const modProgress = { id: modId, progresses: [] };
        state.inProgress.push(modProgress);
        const placeholderProgreess = {
          id: 'placeholder', progress: -1, message: '', fast: false,
        };
        modProgress.progresses.push(placeholderProgreess);
        placeholderProgreess.message = 'Finding the best version to install';
        try {
          await state.selectedInstall.enableMod(modId);
          placeholderProgreess.progress = 1;
          commit('refreshInstalledMods');
        } catch (e) {
          dispatch('showError', e);
        } finally {
          // Allow the UI to update properly
          setTimeout(() => {
            state.inProgress.remove(modProgress);
          }, 500);
        }
      }
    },
    async disableMod({
      commit, dispatch, state,
    }, modId) {
      if (state.installedMods[modId]) {
        if (state.inProgress.length > 0) {
          dispatch('showError', `Another operation is currently in progress while trying to disable a mod: ${state.inProgress.map((progress) => progress.id)}`);
          return;
        }
        if (!state.modsEnabled) {
          dispatch('showError', 'Enable mods to be able to make changes');
          return;
        }
        commit('clearDownloadProgress');
        const modProgress = { id: modId, progresses: [] };
        state.inProgress.push(modProgress);
        const placeholderProgreess = {
          id: 'placeholder', progress: -1, message: '', fast: false,
        };
        modProgress.progresses.push(placeholderProgreess);
        placeholderProgreess.message = 'Checking for mods that are no longer needed';
        try {
          await state.selectedInstall.disableMod(modId);
          placeholderProgreess.progress = 1;
          commit('refreshInstalledMods');
        } catch (e) {
          dispatch('showError', e);
        } finally {
          // Allow the UI to update properly
          setTimeout(() => {
            state.inProgress.remove(modProgress);
          }, 500);
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
      if (!state.modsEnabled) {
        dispatch('showError', 'Enable mods to be able to make changes');
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
        // Allow the UI to update properly
        setTimeout(() => {
          state.inProgress.remove(modProgress);
        }, 500);
      }
    },
    async installSMLVersion({
      commit, dispatch, state,
    }, version) {
      if (state.inProgress.length > 0) {
        dispatch('showError', `Another operation is currently in progress while trying to install a SML version: ${state.inProgress.map((progress) => progress.id)}`);
        return;
      }
      if (!state.modsEnabled) {
        dispatch('showError', 'Enable mods to be able to make changes');
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
        commit('refreshInstalledMods');
      } catch (e) {
        dispatch('showError', e);
      } finally {
        // Allow the UI to update properly
        setTimeout(() => {
          state.inProgress.remove(modProgress);
        }, 500);
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
    renameProfile({ state }, { profile, newName }) {
      const oldName = profile.name;
      renameProfile(oldName, newName);
      const selectedProfile = getSetting('selectedProfile', {});
      Object.keys(selectedProfile).forEach((install) => {
        if (selectedProfile[install] === oldName) {
          selectedProfile[install] = newName;
        }
      });
      saveSetting('selectedProfile', selectedProfile);
      state.profiles.find((p) => p.name === profile.name).name = newName;
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
        const downloadPastProgresses = new Map();
        if (state.downloadProgress[url]) {
          state.downloadProgress[url].downloadPastProgresses.forEach((prog, time) => {
            if (elapsedTime - time <= DOWNLOAD_SPEED_AVG_TIME) {
              downloadPastProgresses.set(time, prog);
            }
          });
        }
        downloadPastProgresses.set(elapsedTime, progress);
        state.downloadProgress[url] = {
          progress,
          name,
          version,
          elapsedTime,
          downloadPastProgresses,
          speed() {
            const entries = Array.from(this.downloadPastProgresses.entries());
            const first = entries[0];
            const last = entries[entries.length - 1];
            return (this.progress.transferred !== this.progress.total && first[0] !== last[0])
              ? (last[1].transferred - first[1].transferred) / ((last[0] - first[0]) / 1000) : 0;
          },
          ETA() {
            return this.speed() !== 0 ? (this.progress.total - this.progress.transferred) / this.speed() : 0;
          },
        };
        if (!state.allDownloadProgress) {
          state.allDownloadProgress = {
            id: '__downloadingFiles__', progress: 0, message: '', fast: true,
          };
          state.inProgress[0].progresses.push(state.allDownloadProgress);
        }
        const speed = Object.values(state.downloadProgress)
          .map((downloadProgress) => downloadProgress.speed())
          .reduce((a, b) => a + b, 0);
        const total = Object.values(state.downloadProgress)
          .map((downloadProgress) => downloadProgress.progress.total)
          .reduce((a, b) => a + b, 0);
        const transferred = Object.values(state.downloadProgress)
          .map((downloadProgress) => downloadProgress.progress.transferred)
          .reduce((a, b) => a + b, 0);
        const percent = transferred / total;
        const ETA = (total - transferred) / speed;
        state.allDownloadProgress.message = `${Object.keys(state.downloadProgress).length > 1 ? `${Object.keys(state.downloadProgress).length} files` : `${limitDownloadNameLength(name)} v${version}`}: ${Math.round(percent * 100)}% | ${bytesToAppropriate(transferred)} / ${bytesToAppropriate(total)} | ${bytesToAppropriate(speed)}/s | ${secondsToAppropriate(ETA)}`;
        state.allDownloadProgress.progress = percent;
        if (progress.percent === 1) {
          setTimeout(() => {
            if (total === transferred) {
              if (state.inProgress.length > 0) {
                state.inProgress[0].progresses.remove(state.allDownloadProgress);
              }
              state.allDownloadProgress = null;
              state.downloadProgress = {};
            }
          }, 100);
        }
      });
      commit('setFavoriteModIds', { favoriteModIds: getSetting('favoriteMods', []) });
      commit('setProfiles', { profiles: getProfiles().filter((profile) => profile.name !== 'vanilla') });
      commit('setExpandModInfoOnStart', getSetting('expandModInfoOnStart', false));
      if (getSetting('konami', false)) {
        commit('konami');
      }
      commit('launchButton', getSetting('launchButton', false));
      commit('launchCat', getSetting('launchCat', false));

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
              state.inProgress.remove(appLoadProgress);
              return;
            }
            commit('setInstalls', { installs });
            const installValidateProgress = { id: 'validatingInstall', progress: -1, message: 'Validating mod install' };
            appLoadProgress.progresses.push(installValidateProgress);
            const savedLocation = getSetting('selectedInstall');
            commit('setInstall', { newInstall: state.satisfactoryInstalls.find((install) => install.installLocation === savedLocation) || state.satisfactoryInstalls[0] });
            let savedProfileName = getSetting('selectedProfile', {})[state.selectedInstall.installLocation] || 'modded';
            if (savedProfileName === 'vanilla') {
              savedProfileName = 'modded'; // Removed vanilla from profiles list
            }
            if (!state.profiles.some((profile) => profile.name.toLowerCase() === savedProfileName.toLowerCase())) {
              savedProfileName = 'modded'; // If profile is missing, default to modded
            }
            commit('setProfile', { newProfile: state.profiles.find((profile) => profile.name.toLowerCase() === savedProfileName.toLowerCase()) });
            let savedModsEnabled = getSetting('modsEnabled', {})[state.selectedInstall.installLocation];
            if (savedModsEnabled === undefined) {
              savedModsEnabled = true;
            }
            commit('setModsEnabled', savedModsEnabled);

            if (!await SatisfactoryInstall.isGameRunning()) {
              try {
                if (savedModsEnabled) {
                  await state.selectedInstall.setProfile(savedProfileName);
                } else {
                  await state.selectedInstall.setProfile('vanilla');
                }
              } catch (e) {
                state.selectedInstall._profile = savedProfileName;
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

      clearOutdatedCache(); // Clear outdated cached mods, SML versions, etc.

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
    konami({ commit }) {
      saveSetting('konami', true);
      commit('konami');
    },
    launchButton({ commit }, value) {
      saveSetting('launchButton', value);
      saveSetting('launchCat', false);
      commit('launchButton', value);
      commit('launchCat', false);
    },
    launchCat({ commit }, value) {
      saveSetting('launchCat', value);
      saveSetting('launchButton', false);
      commit('launchCat', value);
      commit('launchButton', false);
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
        await state.selectedInstall.manifestMutate([], [], [], [], [update.item]);
        placeholderProgreess.progress = 1;
        commit('refreshInstalledMods');
      } catch (e) {
        dispatch('showError', e);
      } finally {
        // Allow the UI to update properly
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
        id: '', progress: -1, message: `Updating ${updates.length} mod${updates.length > 1 ? 's' : ''}`, fast: false,
      };
      updateProgress.progresses.push(placeholderProgreess);
      state.inProgress.push(updateProgress);
      try {
        await state.selectedInstall.manifestMutate([], [], [], [], updates.map((update) => update.item));
        placeholderProgreess.progress = 1;
        commit('refreshInstalledMods');
      } catch (e) {
        dispatch('showError', e);
      } finally {
        // Allow the UI to update properly
        setTimeout(() => {
          state.inProgress.remove(updateProgress);
        }, 500);
      }
    },
  },
});
