<template>
  <v-dialog
    v-show="selectedMod"
    v-model="dialog"
    persistent
    max-width="500px"
  >
    <v-card>
      <v-card-title>
        <span class="headline">{{ isInstall ? 'Install Mod' : 'Uninstall Mod' }}</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-row>
            <v-col cols="12">
              <v-select
                v-model="selectedSatisfactoryInstall"
                hide-details
                prepend-inner-icon="$gamepadIcon"
                item-text="displayName"
                :items="satisfactoryInstalls"
                return-object
              />
            </v-col>
            <v-col cols="12">
              <p>Mod: {{ selectedMod.name }}</p>
            </v-col>
            <v-col cols="12">
              <v-select
                v-model="modalInstallModVersion"
                hide-details
                prepend-inner-icon="$targetVersionIcon"
                item-text="version"
                :items="modVersions"
                return-object
              />
            </v-col>
          </v-row>
        </v-container>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="blue darken-1"
          text
          @click="dialog = false"
        >
          Close
        </v-btn>
        <v-btn
          color="blue darken-1"
          text
          :disabled="inProgress.length > 0 || configLoadInProgress"
          @click="installUninstallMod"
        >
          {{ isInstall ? 'Install' : 'Uninstall' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { mapState } from 'vuex';
import { isVersionSML20Compatible } from '../../utils/helper';

export default {
  data() {
    return {
      dialog: false,
      isInstall: true,
    };
  },
  computed: {
    ...mapState([
      'satisfactoryInstalls',
      'selectedMod',
      'inProgress',
      'configLoadInProgress',
    ]),
    selectedSatisfactoryInstall: {
      get() {
        return this.$store.state.selectedSatisfactoryInstall;
      },
      set(value) {
        this.$store.state.selectedSatisfactoryInstall = value;
      },
    },
    modalInstallModVersion: {
      get() {
        return this.$store.state.modalInstallModVersion;
      },
      set(value) {
        this.$store.state.modalInstallModVersion = value;
      },
    },
    modVersions() {
      return this.$store.state.selectedMod.versions ? this.$store.state.selectedMod.versions.filter((ver) => isVersionSML20Compatible(ver)) : [];
    },
  },
  created() {
    this.$eventBus.$on('show-mod-install-uninstall-dialog', (isInstall) => {
      this.isInstall = isInstall;
      this.dialog = true;
    });
    this.$eventBus.$on('hide-mod-install-uninstall-dialog', () => {
      this.dialog = false;
    });
  },
  methods: {
    installUninstallMod() {
      const submitMethod = this.isInstall ? 'handleModalInstallSubmit' : 'handleModalUninstallSubmit';
      this.$store.dispatch(submitMethod);
    },
  },
};
</script>
