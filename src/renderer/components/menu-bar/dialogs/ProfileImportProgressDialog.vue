<template>
  <v-dialog
    persistent
    :value="isProfileImportInProgress"
    width="500"
    height="230"
  >
    <v-card
      color="loadingBackground !important"
    >
      <v-row
        no-gutters
        justify="center"
      >
        <v-img
          class="mt-4"
          src="static/smm_icon.png"
          max-height="82px"
          max-width="87px"
        />
      </v-row>
      <v-card-title class="loading-text-main">
        IMPORTING PROFILE
      </v-card-title>

      <v-card-text
        v-if="isProfileImportInProgress"
        class="text-center"
      >
        <v-progress-linear
          :value="Math.round(currentProfileImportProgress.progress * 100)"
          :class="{ 'fast': currentProfileImportProgress.fast }"
          background-color="#000000"
          color="#5bb71d"
          height="2"
          reactive
          :indeterminate="currentProfileImportProgress.progress < 0"
        />
        {{ currentProfileImportProgress.message || '&nbsp;' }}
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script>
import { mapState } from 'vuex';
import { lastElement } from '@/utils';

export default {
  computed: {
    ...mapState([
      'inProgress',
    ]),
    isProfileImportInProgress() {
      return this.inProgress.some((prog) => prog.id === '__importProfile__');
    },
    profileImportProgress() {
      return this.inProgress.find((prog) => prog.id === '__importProfile__');
    },
    currentProfileImportProgress() {
      return lastElement(this.profileImportProgress.progresses);
    },
  },
};
</script>
