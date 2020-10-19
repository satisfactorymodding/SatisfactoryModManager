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
import { markdownAsElement } from '@/utils';

export default {
  props: {
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
  methods: {
    changelogElement(release) {
      return markdownAsElement(release.changelog);
    },
    changelogHTML(release) {
      const el = this.changelogElement(release);
      const links = el.getElementsByTagName('a');
      for (let i = 0; i < links.length; i += 1) {
        links[i].target = '_blank';
      }
      return el.innerHTML;
    },
  },
};
</script>
