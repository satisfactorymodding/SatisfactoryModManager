<template>
  <v-card
    tile
    flat
    class="mod-detail color-2 h-100"
  >
    <v-card
      tile
      flat
      class="color-2 d-flex p-3"
    >
      <v-avatar
        class="ma-3"
        size="200"
        tile
      >
        <v-img
          :src="sidePanelData.logo || noImageURL"
          contain
          aspect-ratio="1"
          class="color-2"
        >
          <template v-slot:placeholder>
            <v-row
              class="fill-height ma-0"
              align="center"
              justify="center"
            >
              <v-progress-circular
                indeterminate
                class="color-2"
              />
            </v-row>
          </template>
        </v-img>
      </v-avatar>
      <div>
        <v-card-title
          class="headline"
        >
          {{ sidePanelData.name || '' }}
        </v-card-title>
        <v-card-subtitle>
          By {{ sidePanelData.authors.map((author) => author.user.username).join(', ') }} | v{{ sidePanelData.versions[0] ? sidePanelData.versions[0].version : 'N/A' }}
        </v-card-subtitle>
        <v-card-subtitle>
          {{ sidePanelData.short_description }}
        </v-card-subtitle>
      </div>
    </v-card>
    <v-card
      tile
    >
      <v-toolbar dense>
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <v-icon
              v-on="on"
              @click="closenModDetail"
            >
              $backIcon
            </v-icon>
          </template>
          <span>Close Mod Detail</span>
        </v-tooltip>
        <v-spacer />
        <v-btn
          text
          x-small
          :disabled="inProgress.length > 0 || configLoadInProgress"
          @click="openDialogInstallOldVersion()"
        >
          Install Previous Version
        </v-btn>
      </v-toolbar>
    </v-card>
    <v-card
      tile
      flat
      class="mod-detail-description color-2 p-3"
    >
      <!-- eslint-disable-next-line vue/no-v-html -->
      <div v-html="compiledMarkdownDescription" />
    </v-card>
  </v-card>
</template>

<script>
import { mapState } from 'vuex';
import marked from 'marked';
import sanitizeHtml from 'sanitize-html';

export default {
  computed: {
    ...mapState([
      'sidePanelData',
      'inProgress',
      'configLoadInProgress',
    ]),
    noImageURL() {
      return 'https://ficsit.app/static/assets/images/no_image.png';
    },
    compiledMarkdownDescription() {
      const html = sanitizeHtml(marked(this.sidePanelData.full_description || ''), {
        allowedTags: sanitizeHtml.defaults.allowedTags.concat(['img', 'video', 'details', 'summary']),
        allowedAttributes: Object.assign(sanitizeHtml.defaults.allowedAttributes, { img: ['src', 'width', 'height'], video: ['src', 'width', 'height', 'controls'] }),
      });
      const el = document.createElement('html');
      el.innerHTML = html;
      const links = el.getElementsByTagName('a');
      for (let i = 0; i < links.length; i += 1) {
        links[i].target = '_blank';
      }
      return el.innerHTML;
    },
  },
  methods: {
    closenModDetail() {
      this.$eventBus.$emit('open-close-side-panel');
    },
    openDialogInstallOldVersion() {
      this.$store.state.modalInstallModVersion = {};
      this.$store.state.selectedMod = this.sidePanelData;
      this.$eventBus.$emit('show-mod-install-uninstall-dialog', true);
    },
  },
};
</script>
