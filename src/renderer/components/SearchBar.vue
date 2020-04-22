<template>
  <v-card
    tile
    flat
    class="mod-search pt-0"
  >
    <v-row>
      <v-col
        class="option d-flex pt-0"
        cols="12"
      >
        <v-text-field
          v-model="search"
          label="Search"
          single-line
          hide-details
          :disabled="inProgress.length > 0 || configLoadInProgress"
        />
      </v-col>
    </v-row>
  </v-card>
</template>

<script>
import { mapState } from 'vuex';

export default {
  name: 'SearchBar',
  computed: {
    ...mapState([
      'inProgress',
      'configLoadInProgress',
    ]),
    search: {
      get() {
        return this.$store.state.search;
      },
      set(value) {
        this.$store.commit('setSearch', value);
      },
    },
  },
};
</script>
