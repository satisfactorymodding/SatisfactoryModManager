<template>
  <v-container
    v-if="announcements"
    fluid
    class="mt-2 py-2"
    style="user-select: none;"
    :style="`background: ${announcementsBgColor};`"
  >
    <v-row
      align="center"
      no-gutters
    >
      <v-col
        cols="auto"
        class="pl-2 pr-4"
      >
        <v-icon style="font-size: 32px !important">
          {{ announcementsIcon }}
        </v-icon>
      </v-col>
      <v-col
        cols="10"
        :style="`color: ${announcementsFgColor} !important;`"
      >
        <!-- eslint-disable-next-line vue/no-v-html -->
        <span v-html="announcementMessageHTML" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
// import gql from 'graphql-tag';
import { markdownAsElement } from '@/utils';

export default {
  props: {
  },
  data() {
    return {
      announcements: null, /* {
        message: 'Mods do not work on experimental right now. [Why?](https://ficsit.app)',
        type: 'alert',
      }, */
    };
  },
  apollo: {
    /* announcements: {
      query: gql`
        query getAnnouncements() {
          announcements: getAnnouncements()
          {
            id,
            message,
            type,
          }
        }
      `,
      pollInterval: 5 * 60 * 1000,
    }, */
  },
  computed: {
    announcementsIcon() {
      switch (this.announcements.type) {
        case 'alert':
          return 'mdi-alert-outline';
        default:
          return 'mdi-information-outline';
      }
    },
    announcementsBgColor() {
      switch (this.announcements.type) {
        case 'alert':
          return '#bb2828';
        default:
          return '#5daec5';
      }
    },
    announcementsFgColor() {
      switch (this.announcements.type) {
        default:
          return 'white';
      }
    },
    announcementMessageElement() {
      return markdownAsElement(this.announcements.message);
    },
    announcementMessageHTML() {
      const el = this.announcementMessageElement;
      const links = el.getElementsByTagName('a');
      for (let i = 0; i < links.length; i += 1) {
        links[i].target = '_blank';
        links[i].style += ';color: unset;';
      }
      // get content from inside paragraph
      return el.children[1].children[0].innerHTML;
    },
  },
  watch: {
  },
  methods: {
  },
};
</script>
