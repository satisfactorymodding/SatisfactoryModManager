<template>
  <v-dialog
    v-model="dialogVisible"
    width="500"
    height="230"
    persistent
  >
    <v-card color="loadingBackground !important">
      <v-card-title class="loading-text-main">
        Outdated installed mods
      </v-card-title>

      <v-card-text v-if="!allLeftover">
        Some of your mods are outdated. Do you want to disable them?
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
        Some mods could not be disabled. They are probably dependencies and you should check for updates of their parent mod.
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
          v-if="!allLeftover"
          color="primary"
          text
          @click="disableOutdatedMods"
        >
          Disable outdated mods
        </v-btn>
        <v-btn
          color="text"
          text
          @click="inProgress ? null : (dialogVisible = false)"
        >
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import gql from 'graphql-tag';
import { mapState } from 'vuex';
import { isCompatibleFast } from '@/utils';

export default {
  data() {
    return {
      outdatedInstalledMods: [],
      previousOutatedInstalledMods: [],
      dialogVisible: false,
      inProgress: false,
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
    allLeftover() {
      return this.outdatedInstalledMods.every((mod) => this.previousOutatedInstalledMods.some((prevMod) => prevMod.modReference === mod.modReference));
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
        const modStates = await Promise.all(Object.keys(this.installedMods).map(async (modReference) => {
          if (modReference === 'SML' || modReference === 'bootstrapper') {
            return { modReference, name: modReference, compatible: true };
          }
          const { mod } = (await this.$apollo.query({
            query: gql`            
              query checkOutdatedMod($modReference: ModReference!) {
                mod: getModByReference(modReference: $modReference) {
                  id,
                  mod_reference,
                  name,
                  versions(filter: { limit: 100 }) {
                    id,
                    sml_version,
                  }
                }
              }
            `,
            variables: {
              modReference,
            },
          })).data;
          if (!mod) {
            return { modReference, name: modReference, compatible: false };
          }
          return { modReference, name: mod.name, compatible: await isCompatibleFast(mod, this.$store.state.selectedInstall.version) };
        }));
        this.previousOutatedInstalledMods = this.outdatedInstalledMods;
        this.outdatedInstalledMods = modStates.filter((modState) => !modState.compatible);
        this.dialogVisible = this.outdatedInstalledMods.length > 0;
      }
    },
    async disableOutdatedMods() {
      this.inProgress = true;
      await this.$store.state.selectedInstall.manifestMutate([], [], [], this.outdatedInstalledMods.map((mod) => mod.modReference), []);
      this.$store.commit('refreshInstalledMods');
      this.inProgress = false;
      this.dialogVisible = false;
    },
  },
};
</script>
