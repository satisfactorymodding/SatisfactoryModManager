<template>
  <v-dialog
    v-model="dialogVisible"
    width="500"
    height="230"
    persistent
  >
    <v-card color="loadingBackground !important">
      <v-card-title class="loading-text-main">
        Broken installed mods
      </v-card-title>

      <v-card-text v-if="!allLeftover && updateReminder">
        Remember to <a @click="disableOutdatedMods(true)">check for updates</a>, and click "Update all" in the updates menu rather than updating individual mods.
      </v-card-text>

      <v-card-text v-else-if="!allLeftover">
        Some of your mods are broken. Do you want to disable them?
        <ul>
          <li
            v-for="mod in outdatedInstalledMods"
            :key="mod.modReference"
          >
            {{ mod.name }}
          </li>
        </ul>
      </v-card-text>

      <v-card-text v-else>
        Some mods could not be disabled. They are probably dependencies and you should <a @click="disableOutdatedMods(true)">check for updates</a>.
        <ul>
          <li
            v-for="mod in outdatedInstalledMods"
            :key="mod.modReference"
          >
            {{ mod.name }}
          </li>
        </ul>
      </v-card-text>

      <v-card-actions>
        <v-btn
          color="primary"
          text
          @click="disableOutdatedMods(true)"
        >
          <span v-if="!updateReminder">Disable broken mods &amp; update</span>
          <span v-else>Update mods</span>
        </v-btn>
        <v-btn
          color="text"
          text
          @click="inProgress ? null : (dialogVisible = updateReminder = false)"
        >
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { mapGetters, mapState } from 'vuex';
import { getModCompatibilityState, COMPATIBILITY_LEVEL } from '@/utils';

export default {
  data() {
    return {
      outdatedInstalledMods: [],
      previousOutatedInstalledMods: [],
      dialogVisible: false,
      inProgress: false,
      updateReminder: false,
    };
  },
  computed: {
    ...mapState(
      [
        'modsEnabled',
        'manifestMods',
        'installedMods',
        'selectedProfile',
      ],
    ),
    ...mapGetters([
      'branch',
      'gameVersion',
    ]),
    allLeftover() {
      return this.outdatedInstalledMods.length > 0 && this.outdatedInstalledMods.every((mod) => this.previousOutatedInstalledMods.some((prevMod) => prevMod.modReference === mod.modReference));
    },
  },
  watch: {
    modsEnabled() {
      this.checkOutdated();
    },
    installedMods() {
      this.checkOutdated();
    },
    selectedProfile() {
      this.previousOutatedInstalledMods = [];
    },
  },
  methods: {
    async checkOutdated() {
      if (this.modsEnabled) {
        const modStates = await Promise.all(Object.keys(this.installedMods).map((modReference) => getModCompatibilityState(modReference, this.branch, this.gameVersion)));
        this.outdatedInstalledMods = modStates.filter((modState) => modState.compatible === COMPATIBILITY_LEVEL.INCOMPATIBLE);
        this.dialogVisible = this.updateReminder || this.outdatedInstalledMods.length > 0;
      }
    },
    async disableOutdatedMods(update = false) {
      this.previousOutatedInstalledMods = this.outdatedInstalledMods;
      this.inProgress = true;
      this.dialogVisible = false;

      const progress = {
        id: '__loadingApp__',
        progresses: [],
      };
      const placeholderProgreess = {
        id: '', progress: -1, message: `Disabling broken mods${update ? ', updating all' : ''}`, fast: false,
      };
      if (this.outdatedInstalledMods.length === 0) {
        placeholderProgreess.message = 'Updating mods';
      }
      progress.progresses.push(placeholderProgreess);
      this.$store.state.inProgress.push(progress);

      await this.$store.state.selectedInstall.manifestMutate([], [], [], this.outdatedInstalledMods.map((mod) => mod.modReference), update ? Object.keys(this.installedMods) : []);
      this.$store.commit('refreshInstalledMods');

      this.$store.state.inProgress.remove(progress);
      this.inProgress = false;
      this.dialogVisible = !update;
      this.updateReminder = !update;
    },
  },
};
</script>
