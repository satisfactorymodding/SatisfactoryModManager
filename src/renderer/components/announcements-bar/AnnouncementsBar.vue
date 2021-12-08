<template>
  <v-container
    v-if="announcements"
    fluid
    class="ma-0 pa-0"
    style="user-select: none;"
  >
    <v-carousel
      hide-delimiters
      height="auto"
      :show-arrows="false"
      cycle
      continuous
      interval="10000"
    >
      <v-carousel-item
        v-for="announcement in announcements"
        :key="announcement.id"
      >
        <Announcement
          :announcement="announcement"
        />
      </v-carousel-item>
    </v-carousel>
  </v-container>
</template>

<script>
import gql from 'graphql-tag';
import Announcement from './Announcement';

export default {
  components: { Announcement },
  apollo: {
    announcements: {
      query: gql`
        query getAnnouncements {
          announcements: getAnnouncements
          {
            id,
            message,
            importance,
          }
        }
      `,
      pollInterval: 60 * 1000,
    },
  },
  watch: {
  },
  methods: {
  },
};
</script>
