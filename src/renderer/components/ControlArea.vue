<template>
  <v-container
    fluid
  >
    <v-row
      no-gutters
      class="px-2"
    >
      <v-col
        cols="12"
      >
        <v-select
          v-model="selectedInstallModel"
          :disabled="!!inProgress.length || isGameRunning"
          :items="satisfactoryInstalls"
          item-text="displayName"
          return-object
          dense
          filled
          class="custom"
          solo
          append-icon="mdi-chevron-down"
        >
          <template #selection="{ item }">
            <span>Install:&nbsp;</span>
            <span class="green--text">{{ item.displayName }}</span>
            <span
              v-if="item.smlVersion"
              class="green--text"
            >&nbsp;-&nbsp;SML v{{ item.smlVersion }}</span>
          </template>
        </v-select>
      </v-col>
    </v-row>
    <v-row
      no-gutters
      class="px-2"
    >
      <v-col
        cols="6"
      >
        <v-btn-toggle
          v-model="modsEnabled"
          mandatory
          style="width: 100%; height: 28px;"
          borderless
          dense
        >
          <v-btn
            disabled
            class="flex-grow-1 custom"
            style="justify-content: left; height: 100%; color: unset !important; opacity: 1"
          >
            <span style="padding-left: 10.5px; padding-top: 4px">Mods</span>
          </v-btn>

          <v-btn
            :value="false"
            class="custom-btn"
            style="height: 100%; opacity: 1"
            :class="!modsEnabled ? 'red--text' : ''"
          >
            <span style="padding-top: 2px">
              OFF&nbsp;
            </span>
            <v-icon style="color: unset">
              mdi-window-close
            </v-icon>
          </v-btn>

          <v-btn
            :value="true"
            class="custom-btn"
            style="height: 100%; opacity: 1;"
            :class="modsEnabled ? 'green--text' : ''"
          >
            <span style="padding-top: 2px">
              ON&nbsp;
            </span>
            <v-icon style="color: unset">
              mdi-check
            </v-icon>
          </v-btn>
        </v-btn-toggle>
      </v-col>
      <v-col
        class="pl-1 flex-grow-1"
      >
        <v-select
          v-model="selectedProfileModel"
          :disabled="!!inProgress.length || isGameRunning"
          :items="profiles"
          return-object
          item-text="name"
          dense
          filled
          class="custom"
          solo
          append-icon="mdi-chevron-down"
        >
          <template #selection="{ item }">
            <span>Profile:&nbsp;</span>
            <span class="green--text">({{ item.name }})</span>
          </template>
        </v-select>
      </v-col>
      <v-col
        cols="auto"
        class="buttons"
      >
        <v-menu offset-y>
          <template
            #activator="{ on, attrs }"
          >
            <v-btn
              v-bind="attrs"
              class="custom"
              style="min-width: 28px; height: 28px"
              v-on="on"
            >
              <v-icon
                style="color: unset; padding-top: 2px"
              >
                mdi-cog
              </v-icon>
            </v-btn>
          </template>
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
          <v-btn
            text
            :disabled="!!inProgress.length || isGameRunning || selectedProfileModel.name === 'vanilla' || selectedProfileModel.name === 'modded' || selectedProfileModel.name === 'development'"
            @click="showRenameProfileDialog"
          >
            Rename&nbsp;
            <v-icon
              color="yellow"
              class="icon"
            >
              mdi-pencil
            </v-icon>
          </v-btn>
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
        </v-menu>
      </v-col>
    </v-row>
    <v-row
      no-gutters
      class="px-2"
    >
      <v-col
        cols="6"
      >
        <v-select
          v-model="modFiltersModel"
          :items="availableFilters"
          return-object
          item-text="name"
          dense
          filled
          class="custom"
          solo
          append-icon="mdi-chevron-down"
        >
          <template #selection="{ item }">
            <span>Show:&nbsp;</span>
            <span class="green--text">{{ item.name }} ({{ item.mods }})</span>
          </template>
          <template #item="{ item }">
            {{ item.name }} ({{ item.mods }})
          </template>
        </v-select>
      </v-col>
      <v-col
        cols="6"
        class="pl-1"
      >
        <v-select
          v-model="sortByModel"
          :items="availableSorting"
          return-object
          item-text="name"
          dense
          filled
          class="custom"
          solo
          append-icon="mdi-chevron-down"
        >
          <template #selection="{ item }">
            <span>Sort:&nbsp;</span>
            <span class="green--text">{{ item.name }}</span>
          </template>
        </v-select>
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col
        cols="12"
        class="px-2 mb-n2"
      >
        <v-text-field
          v-model="searchModel"
          dense
          filled
          solo
          class="custom-search"
          placeholder="Type to search"
        >
          <template #append>
            <div
              class="d-inline-flex align-center fill-height"
              style="background-color: unset"
              @click="searchModel = ''"
            >
              <span class="pr-1 icon--text">Clear</span>
              <v-icon class="red--text">
                mdi-window-close
              </v-icon>
            </div>
          </template>
        </v-text-field>
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
            @submit.stop.prevent="createProfile"
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
            <span class="warning--text">{{ newProfileMessage }}</span>
          </v-form>
        </v-card-text>
        <v-card-actions>
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
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog
      v-model="renameProfileDialog"
      persistent
    >
      <v-card>
        <v-card-title class="headline">
          Rename profile
        </v-card-title>

        <v-card-text>
          <span>Current profile name: {{ selectedProfileModel.name }}</span>
          <v-form
            ref="renameProfileForm"
            v-model="newProfileFormValid"
            @submit.stop.prevent="renameProfile"
          >
            <v-text-field
              v-model="newProfileName"
              label="Name"
              required
              :rules="[v => !!v || 'Name is required']"
            />
            <span class="warning--text">{{ newProfileMessage }}</span>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="primary"
            text
            @click="renameProfile"
          >
            Rename
          </v-btn>
          <v-btn
            color="text"
            text
            @click="cancelRenameProfile"
          >
            Cancel
          </v-btn>
        </v-card-actions>
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
import { filenamify } from '@/utils';

export default {
  props: {
    filters: {
      type: Object,
      required: true,
    },
    availableFilters: {
      type: Array,
      required: true,
    },
    availableSorting: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      newProfileFormValid: true,
      newProfileDialog: false,
      newProfileName: '',
      newProfileCopyCurrent: false,
      newProfileMessage: '',
      deleteProfileDialog: false,
      renameProfileDialog: false,
    };
  },
  computed: {
    ...mapState([
      'satisfactoryInstalls',
      'profiles',
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
    modsEnabled: {
      get() { return this.$store.state.modsEnabled; },
      set(value) { this.$store.dispatch('setModsEnabled', value); },
    },
    modFiltersModel: {
      get() { return this.filters.modFilters; },
      set(value) { this.$emit('update:filters', { ...this.filters, modFilters: value }); },
    },
    searchModel: {
      get() { return this.filters.search; },
      set(value) { this.$emit('update:filters', { ...this.filters, search: value }); },
    },
    sortByModel: {
      get() { return this.filters.sortBy; },
      set(value) { this.$emit('update:filters', { ...this.filters, sortBy: value }); },
    },
    filtersModel: {
      get() { return this.filters; },
      set(data) { this.$emit('update:filters', data); },
    },
  },
  watch: {
    newProfileName(name) {
      const validName = filenamify(name);
      if (name !== validName) {
        this.newProfileMessage = `Profile will be saved as ${validName}`;
      } else {
        this.newProfileMessage = '';
      }
    },
  },
  methods: {
    showCreateProfileDialog() {
      this.newProfileDialog = true;
    },
    createProfile() {
      if (this.$refs.newProfileForm.validate()) {
        this.$store.dispatch('createProfile', { profileName: filenamify(this.newProfileName), copyCurrent: this.newProfileCopyCurrent });
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
    showRenameProfileDialog() {
      this.renameProfileDialog = true;
    },
    renameProfile() {
      if (this.$refs.renameProfileForm.validate()) {
        this.$store.dispatch('renameProfile', { newProfile: filenamify(this.newProfileName) });
        this.cancelRenameProfile();
      }
    },
    cancelRenameProfile() {
      this.newProfileName = '';
      this.renameProfileDialog = false;
    },
  },
};
</script>

<style>
.custom-search.v-text-field {
  padding-top: 0px !important;
  margin-top: 0px !important;
}
.custom-search .v-label {
  font-size: 14px !important;
}
.custom-search.v-text-field.v-text-field--solo.v-input--dense>.v-input__control {
  min-height: 28px;
}
.custom.v-text-field--filled.v-input--dense.v-text-field--filled>.v-input__control>.v-input__slot, .custom.v-text-field.v-text-field--solo.v-input--dense>.v-input__control,
.custom-search.v-text-field--filled.v-input--dense.v-text-field--filled>.v-input__control>.v-input__slot {
  min-height: 28px !important;
}
.custom.v-input--dense>.v-input__control>.v-input__slot {
  margin-bottom: 0px;
}
.custom.theme--dark.v-btn.v-btn--disabled.v-btn--has-bg {
  background-color: #2b2b2b !important;
}
.custom.theme--dark.v-text-field--solo>.v-input__control>.v-input__slot,
.custom-search.theme--dark.v-text-field--solo>.v-input__control>.v-input__slot {
  background-color: #2b2b2b !important;
}
.custom.theme--dark.v-btn.v-btn--disabled.v-btn--has-bg:hover {
  background-color: #575757 !important;
}
.custom.theme--dark.v-text-field--solo:hover>.v-input__control>.v-input__slot,
.custom-search.theme--dark.v-text-field--solo:hover>.v-input__control>.v-input__slot {
  background-color: #575757 !important;
}
.v-list--dense .v-list-item .v-list-item__title {
  font-size: 0.9125rem !important;
}
</style>

<style scoped>
.buttons {
  padding-left: 3px !important;
}

.v-btn {
  color: var(--v-text-base) !important;
}

.custom {
  padding-top: 0px !important;
  margin-top: 0px !important;
  padding-bottom: 4px !important;
  height: 32px;
  font-size: 13px !important;
}

.custom-btn {
  font-size: 12px !important;
}

.custom.theme--dark.v-btn.v-btn--has-bg {
  background-color: #2b2b2b;
}

.custom-btn.theme--dark:before {
  background-color: var(--v-backgroundMenuBar-base);
  opacity: 1;
}
.custom-btn.theme--dark.v-btn--active:before {
  background-color: #2b2b2b;
}

div {
  background-color: var(--v-backgroundMenuBar-base);
}
</style>
