<template>
  <v-container
    fluid
  >
    <v-row
      no-gutters
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
              v-model="selectedProfileModel"
              :disabled="!!inProgress.length || isGameRunning"
              :items="profiles"
              item-text="name"
              return-object
              label="SELECT PROFILE"
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
                  @click="showCreateProfileDialog"
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
                  :disabled="!!inProgress.length || isGameRunning || selectedProfileModel.name === 'vanilla' || selectedProfileModel.name === 'modded' || selectedProfileModel.name === 'development'"
                  @click="showDeleteProfileDialog"
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
      v-model="newProfileDialog"
      persistent
    >
      <v-card>
        <v-card-title class="headline">
          New profile
        </v-card-title>

        <v-card-text>
          <v-form
            ref="newProfileForm"
            v-model="newProfileFormValid"
          >
            <v-text-field
              v-model="newProfileName"
              label="Name"
              required
              :rules="[v => !!v || 'Name is required']"
            />
            <v-switch
              v-model="newProfileCopyCurrent"
              label="Copy current profile"
            />
            <v-btn
              color="primary"
              text
              @click="createProfile"
            >
              Create
            </v-btn>
            <v-btn
              color="text"
              text
              @click="cancelCreateProfile"
            >
              Cancel
            </v-btn>
          </v-form>
        </v-card-text>
      </v-card>
    </v-dialog>
    <v-dialog
      v-model="deleteProfileDialog"
      persistent
    >
      <v-card>
        <v-card-title class="headline">
          Delete profile
        </v-card-title>

        <v-card-text>
          <span>Are you sure you want to delete profile {{ selectedProfileModel.name }}</span>
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="text"
            text
            @click="deleteProfile"
          >
            Delete
          </v-btn>
          <v-btn
            color="primary"
            text
            @click="cancelDeleteProfile"
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
      newProfileFormValid: true,
      newProfileDialog: false,
      newProfileName: '',
      newProfileCopyCurrent: false,
      deleteProfileDialog: false,
    };
  },
  computed: {
    ...mapState([
      'satisfactoryInstalls',
      'profiles',
      'modFilters',
      'sortBy',
      'inProgress',
      'isGameRunning',
    ]),
    selectedInstallModel: {
      get() { return this.$store.state.selectedInstall; },
      set(value) { this.$store.dispatch('selectInstall', value); },
    },
    selectedProfileModel: {
      get() { return this.$store.state.selectedProfile; },
      set(value) { this.$store.dispatch('selectProfile', value); },
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
    showCreateProfileDialog() {
      this.newProfileDialog = true;
    },
    createProfile() {
      if (this.$refs.newProfileForm.validate()) {
        this.$store.dispatch('createProfile', { profileName: this.newProfileName, copyCurrent: this.newProfileCopyCurrent });
        this.cancelCreateProfile();
      }
    },
    cancelCreateProfile() {
      this.newProfileName = '';
      this.newProfileCopyCurrent = false;
      this.newProfileDialog = false;
    },
    showDeleteProfileDialog() {
      this.deleteProfileDialog = true;
    },
    deleteProfile() {
      this.$store.dispatch('deleteProfile', { profileName: this.$store.state.selectedProfile.name });
      this.cancelDeleteProfile();
    },
    cancelDeleteProfile() {
      this.deleteProfileDialog = false;
    },
  },
};
</script>

<style>
.custom-search .v-label {
  font-size: 14px !important;
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
