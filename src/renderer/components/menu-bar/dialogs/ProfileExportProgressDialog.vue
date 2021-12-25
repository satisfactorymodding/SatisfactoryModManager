<template>
  <v-dialog
    persistent
    :value="isProfileExportInProgress"
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
        EXPORTING PROFILE
      </v-card-title>

      <v-card-text
        v-if="isProfileExportInProgress"
        class="text-center"
      >
        <v-progress-linear
          :value="Math.round(currentProfileExportProgress.progress * 100)"
          :class="{ 'fast': currentProfileExportProgress.fast }"
          background-color="#000000"
          color="#5bb71d"
          height="2"
          reactive
          :indeterminate="currentProfileExportProgress.progress < 0"
        />
        {{ currentProfileExportProgress.message || '&nbsp;' }}
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
    isProfileExportInProgress() {
      return this.inProgress.some((prog) => prog.id === '__exportProfile__');
    },
    profileExportProgress() {
      return this.inProgress.find((prog) => prog.id === '__exportProfile__');
    },
    currentProfileExportProgress() {
      return lastElement(this.profileExportProgress.progresses);
    },
  },
};
</script>
