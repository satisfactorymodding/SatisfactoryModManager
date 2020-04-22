<template>
  <v-container
    fluid
    class="mod-config color-1 pt-0"
  >
    <v-row align="center">
      <v-col
        class="option d-flex flex-column"
        cols="12"
      >
        <v-select
          v-model="selectedSatisfactoryInstall"
          hide-details
          prepend-inner-icon="$gamepadIcon"
          item-text="displayName"
          :items="satisfactoryInstalls"
          :disabled="inProgress.length > 0 || configLoadInProgress"
          return-object
        >
          <template
            slot="selection"
            slot-scope="data"
          >
            {{ getGameVersion(data.item) }}
          </template>
          <template
            slot="item"
            slot-scope="data"
          >
            {{ getGameVersion(data.item) }}
          </template>
          <template
            v-if="selectedSatisfactoryInstall && selectedSatisfactoryInstall.smlVersion && cachedSMLHasUpdate"
            slot="append-outer"
          >
            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-icon
                  v-on="on"
                  @click="updateSML"
                >
                  $updateIcon
                </v-icon>
              </template>
              <span>Update SML</span>
            </v-tooltip>
          </template>
        </v-select>
      </v-col>
      <v-col
        class="option d-flex flex-column"
        cols="12"
      >
        <v-select
          v-model="selectedConfig"
          hide-details
          prepend-inner-icon="$configIcon"
          item-value="id"
          item-text="name"
          :items="availableConfigs"
          :disabled="inProgress.length > 0 || configLoadInProgress"
        >
          <template slot="append-outer">
            <NewConfigDialog :dialog="newConfigDialog" />
            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-icon
                  :disabled="inProgress.length > 0 || configLoadInProgress"
                  v-on="on"
                  @click="newConfigDialog = true"
                >
                  $addIcon
                </v-icon>
              </template>
              <span>Add Config</span>
            </v-tooltip>
            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-icon
                  :disabled="inProgress.length > 0 || configLoadInProgress"
                  v-on="on"
                  @click="deleteSelectedConfig"
                >
                  $deleteIcon
                </v-icon>
              </template>
              <span>Delete Config'</span>
            </v-tooltip>
          </template>
        </v-select>
      </v-col>
      <v-col
        class="option d-flex flex-column"
        cols="6"
      >
        <v-select
          v-model="filters.installedStatus"
          hide-details
          prepend-inner-icon="$filterIcon"
          item-value="value"
          item-text="displayName"
          :items="installedStatusOptions"
          :disabled="inProgress.length > 0 || configLoadInProgress"
        >
          <template slot="append-outer">
            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-icon
                  :disabled="inProgress.length > 0 || configLoadInProgress"
                  v-on="on"
                  @click="filters.compatibleOnly = !filters.compatibleOnly"
                >
                  {{ filters.compatibleOnly ? '$compatibleIcon' : '$notCompatibleIcon' }}
                </v-icon>
              </template>
              <span>{{ filters.compatibleOnly ? 'Compatible Only' : 'All' }}</span>
            </v-tooltip>
          </template>
        </v-select>
      </v-col>
      <v-col
        class="option d-flex flex-column"
        cols="6"
      >
        <v-select
          v-model="filters.sortBy"
          hide-details
          prepend-inner-icon="$sortIcon"
          item-value="value"
          item-text="displayName"
          :items="sortByOptions"
          :disabled="inProgress.length > 0 || configLoadInProgress"
        >
          <template slot="append-outer">
            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-icon
                  :disabled="inProgress.length > 0 || configLoadInProgress"
                  v-on="on"
                  @click="filters.sortOrder = filters.sortOrder === 'ascending' ? 'descending' : 'ascending'"
                >
                  {{ filters.sortOrder === 'ascending' ? '$sortAscIcon' : '$sortDescIcon' }}
                </v-icon>
              </template>
              <span>{{ filters.sortOrder === 'ascending' ? 'Ascending' : 'Descending' }}</span>
            </v-tooltip>
          </template>
        </v-select>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapActions } from 'vuex';
import { saveSetting } from '../settings';
import NewConfigDialog from './dialogs/NewConfig';

export default {
  name: 'ModConfig',
  components: {
    NewConfigDialog,
  },
  data() {
    return {
      newConfigDialog: false,
    };
  },
  computed: {
    selectedSatisfactoryInstall: {
      get() {
        return this.$store.state.selectedSatisfactoryInstall;
      },
      set(value) {
        this.$store.state.selectedSatisfactoryInstall = value;
      },
    },
    satisfactoryInstalls() {
      return this.$store.state.satisfactoryInstalls;
    },
    installedStatusOptions() {
      return this.$store.state.installedStatusOptions;
    },
    sortByOptions() {
      return this.$store.state.sortByOptions;
    },
    filters: {
      get() {
        return this.$store.state.filters;
      },
      set(value) {
        this.$store.state.filters = value;
      },
    },
    availableConfigs() {
      return this.$store.state.availableConfigs;
    },
    selectedConfig: {
      get() {
        return this.$store.state.selectedConfig;
      },
      set(value) {
        this.$store.state.selectedConfig = value;
      },
    },
    configLoadInProgress() {
      return this.$store.state.configLoadInProgress;
    },
    inProgress() {
      return this.$store.state.inProgress;
    },
    cachedSMLHasUpdate() {
      return this.$store.state.cachedSMLHasUpdate;
    },
  },
  watch: {
    selectedSatisfactoryInstall() {
      this.$store.commit('refreshSearch');
      saveSetting('selectedSFInstall', this.selectedSatisfactoryInstall.installLocation);
      this.$store.dispatch('checkForUpdates');
    },
    filters: {
      handler() {
        this.$store.commit('refreshSearch');
        saveSetting('filters', this.filters);
      },
      deep: true,
    },
    selectedConfig() {
      saveSetting('selectedConfig', this.selectedConfig);
      this.$store.commit('loadSelectedConfig');
    },
  },
  methods: {
    ...mapActions([
      'updateSML',
      'deleteSelectedConfig',
    ]),
    getGameVersion(item) {
      const smlVersion = item.smlVersion ? `- SML v${item.smlVersion}` : '';
      return `${item.displayName} ${smlVersion}`;
    },
  },
};
</script>
