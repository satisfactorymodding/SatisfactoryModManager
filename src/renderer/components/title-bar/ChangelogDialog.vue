<template>
  <v-dialog v-model="changelogDialog">
    <v-card v-if="viewChangelogUpdate">
      <v-card-title>
        {{ viewChangelogUpdate.name }} v{{ viewChangelogUpdate.version }} changelog
      </v-card-title>
      <v-card-text>
        <template v-for="(release, index) in viewChangelogUpdate.releases">
          <div :key="index">
            <h3>v{{ release.version }}</h3>
            <!-- eslint-disable vue/no-v-html -->
            <span v-html="changelogHTML(release)" />
            <v-divider v-if="index != viewChangelogUpdate.releases.length - 1" />
          </div>
        </template>
      </v-card-text>
      <v-card-actions>
        <v-btn
          color="primary"
          text
          @click="changelogDialog = false"
        >
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script>
export default {
  name: 'ChangelogDialog',
  props: {
    changelogHTML: {
      type: Function,
      default: () => {},
    },
    viewChangelogUpdate: {
      type: Object,
      default: () => undefined,
    },
  },
  data() {
    return {
      changelogDialog: false,
    };
  },
};
</script>
<style scoped>
</style>
