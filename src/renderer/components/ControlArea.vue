<template>
  <v-container
    fluid
  >
    <v-row
      no-gutters
      class="mt-4"
    >
      <v-col
        cols="12"
        class="px-2"
      >
        <v-select
          v-model="selectedInstallModel"
          :disabled="!!inProgress.length || isGameRunning"
          :items="satisfactoryInstalls"
          item-text="displayName"
          return-object
          label="SATISFACTORY INSTALL"
          class="custom"
          append-icon="mdi-chevron-down"
        >
          <template v-slot:selection="{ item }">
            <span>{{ item.displayName }}</span>
            <span v-if="item.smlVersion">&nbsp;-&nbsp;SML v{{ item.smlVersion }}</span>
          </template>
        </v-select>
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col cols="12">
        <v-row
          justify="center"
          no-gutters
        >
          <v-col
            cols="6"
            class="px-2"
          >
            <v-select
              v-model="selectedConfigModel"
              :disabled="!!inProgress.length || isGameRunning"
              :items="configs"
              item-text="name"
              return-object
              label="SELECT CONFIG"
              class="custom"
              append-icon="mdi-chevron-down"
            />
          </v-col>
          <v-col
            cols="6"
            class="buttons"
          >
            <v-row>
              <v-col
                cols="6"
              >
                <v-btn
                  text
                  :disabled="!!inProgress.length || isGameRunning"
                  @click="showCreateConfigDialog"
                >
                  New&nbsp;
                  <v-icon
                    color="green"
                    class="icon"
                  >
                    mdi-plus-circle
                  </v-icon>
                </v-btn>
              </v-col>
              <v-col
                cols="6"
              >
                <v-btn
                  text
                  :disabled="!!inProgress.length || isGameRunning || selectedConfigModel.name === 'vanilla' || selectedConfigModel.name === 'modded' || selectedConfigModel.name === 'development'"
                  @click="showDeleteConfigDialog"
                >
                  Delete&nbsp;
                  <v-icon
                    color="red"
                    class="icon"
                  >
                    mdi-delete
                  </v-icon>
                </v-btn>
              </v-col>
            </v-row>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col
        cols="12"
      >
        <v-row
          justify="center"
          no-gutters
        >
          <v-col
            cols="6"
            class="px-2"
          >
            <v-select
              v-model="modFiltersModel"
              :disabled="!!inProgress.length || isGameRunning"
              :items="modFilters"
              item-text="name"
              :return-object="true"
              label="SHOW"
              class="custom"
              append-icon="mdi-chevron-down"
            >
              <template v-slot:selection="{ item }">
                <span>{{ item.name }}</span>
                &nbsp;
                <span class="green--text">({{ item.mods }})</span>
              </template>
              <template
                slot="item"
                slot-scope="data"
              >
                <!-- HTML that describe how select should render items when the select is open -->
                {{ data.item.name }} ({{ data.item.mods }})
              </template>
            </v-select>
          </v-col>
          <v-col
            cols="6"
            class="px-2"
          >
            <v-select
              v-model="sortByModel"
              :disabled="!!inProgress.length || isGameRunning"
              :items="sortBy"
              label="SORT BY"
              class="custom"
              append-icon="mdi-chevron-down"
            />
          </v-col>
        </v-row>
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col
        cols="12"
        class="px-2"
      >
        <v-text-field
          v-model="searchModel"
          class="custom-search"
          :disabled="!!inProgress.length || isGameRunning"
          label="Search"
        />
      </v-col>
    </v-row>
    <v-dialog
      v-model="newConfigDialog"
      persistent
    >
      <v-card>
        <v-card-title class="headline">
          New config
        </v-card-title>

        <v-card-text>
          <v-form
            ref="newConfigForm"
            v-model="newConfigFormValid"
          >
            <v-text-field
              v-model="newConfigName"
              label="Name"
              required
              :rules="[v => !!v || 'Name is required']"
            />
            <v-switch
              v-model="newConfigCopyCurrent"
              label="Copy current config"
            />
            <v-btn
              color="primary"
              text
              @click="createConfig"
            >
              Create
            </v-btn>
            <v-btn
              color="text"
              text
              @click="cancelCreateConfig"
            >
              Cancel
            </v-btn>
          </v-form>
        </v-card-text>
      </v-card>
    </v-dialog>
    <v-dialog
      v-model="deleteConfigDialog"
      persistent
    >
      <v-card>
        <v-card-title class="headline">
          Delete config
        </v-card-title>

        <v-card-text>
          <span>Are you sure you want to delete config {{ selectedConfigModel.name }}</span>
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="text"
            text
            @click="deleteConfig"
          >
            Delete
          </v-btn>
          <v-btn
            color="primary"
            text
            @click="cancelDeleteConfig"
          >
            Cancel
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
import { mapState } from 'vuex';

export default {
  data() {
    return {
      newConfigFormValid: true,
      newConfigDialog: false,
      newConfigName: '',
      newConfigCopyCurrent: false,
      deleteConfigDialog: false,
    };
  },
  computed: {
    ...mapState([
      'satisfactoryInstalls',
      'configs',
      'modFilters',
      'sortBy',
      'inProgress',
      'isGameRunning',
    ]),
    selectedInstallModel: {
      get() { return this.$store.state.selectedInstall; },
      set(value) { this.$store.dispatch('selectInstall', value); },
    },
    selectedConfigModel: {
      get() { return this.$store.state.selectedConfig; },
      set(value) { this.$store.dispatch('selectConfig', value); },
    },
    modFiltersModel: {
      get() { return this.$store.state.filters.modFilters; },
      set(value) {
        const newFilters = this.$store.state.filters;
        newFilters.modFilters = value;
        this.$store.dispatch('setFilters', newFilters);
      },
    },
    searchModel: {
      get() { return this.$store.state.filters.search; },
      set(value) {
        const newFilters = this.$store.state.filters;
        newFilters.search = value;
        this.$store.dispatch('setFilters', newFilters);
      },
    },
    sortByModel: {
      get() { return this.$store.state.filters.sortBy; },
      set(value) {
        const newFilters = this.$store.state.filters;
        newFilters.sortBy = value;
        this.$store.dispatch('setFilters', newFilters);
      },
    },
  },
  methods: {
    showCreateConfigDialog() {
      this.newConfigDialog = true;
    },
    createConfig() {
      if (this.$refs.newConfigForm.validate()) {
        this.$store.dispatch('createConfig', { configName: this.newConfigName, copyCurrent: this.newConfigCopyCurrent });
        this.cancelCreateConfig();
      }
    },
    cancelCreateConfig() {
      this.newConfigName = '';
      this.newConfigCopyCurrent = false;
      this.newConfigDialog = false;
    },
    showDeleteConfigDialog() {
      this.deleteConfigDialog = true;
    },
    deleteConfig() {
      this.$store.dispatch('deleteConfig', { configName: this.$store.state.selectedConfig.name });
      this.cancelDeleteConfig();
    },
    cancelDeleteConfig() {
      this.deleteConfigDialog = false;
    },
  },
};
</script>

<style>
.custom-search .v-label {
  font-size: 12px !important;
}
</style>

<style scoped>
.buttons {
  padding-left: 3px !important;
}

.v-btn {
  color: var(--v-text-base) !important;
}
</style>
