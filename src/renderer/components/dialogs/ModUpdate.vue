<template>
  <v-dialog
    v-model="dialog"
    persistent
    max-width="500px"
  >
    <v-card class="position-relative">
      <v-card
        v-show="inProgress.length > 0 || configLoadInProgress"
        class="position-absolute w-100"
      >
        <v-progress-linear
          color="success"
          indeterminate
          height="2"
        ></v-progress-linear>
      </v-card>
      <v-card-title>
        <span class="headline">Mod Update Available</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-list>
            <v-list-item
              v-for="item in updates"
              :key="item.id"
            >
              <v-list-item-content>
                {{ (availableMods.find((mod) => mod.id === item.id) || { name: item.id }).name }} v{{ item.version }}
              </v-list-item-content>
              <v-list-item-action
                class="d-flex flex-row"
              >
                <v-btn
                  icon
                  :disabled="inProgress.length > 0 || configLoadInProgress"
                  @click="updateById(item.id)"
                >
                  <v-icon color="green lighten-1">
                    $updateIcon
                  </v-icon>
                </v-btn>
                <v-btn
                  icon
                  :disabled="inProgress.length > 0 || configLoadInProgress"
                  @click="ignoreVersion(item)"
                >
                  <v-icon color="red lighten-1">
                    $ignoreIcon
                  </v-icon>
                </v-btn>
              </v-list-item-action>
            </v-list-item>
          </v-list>
        </v-container>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="blue darken-1"
          text
          @click="dialog = false"
        >
          Close
        </v-btn>
        <v-btn
          color="blue darken-1"
          text
          :disabled="inProgress.length > 0 || configLoadInProgress"
          @click="updateAll"
        >
          Update All
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { mapState } from 'vuex';

export default {
  data() {
    return {
      dialog: false,
    };
  },
  computed: {
    ...mapState([
      'inProgress',
      'configLoadInProgress',
    ]),
    updates() {
      return this.$store.state.updates;
    },
    availableMods() {
      return this.$store.state.availableMods;
    },
  },
  created() {
    this.$eventBus.$on('show-mod-update-dialog', () => {
      this.dialog = true;
    });
    this.$eventBus.$on('hide-mod-update-dialog', () => {
      this.dialog = false;
    });
  },
  methods: {
    updateById(id) {
      this.$store.dispatch('updateById', id);
    },
    ignoreVersion(item) {
      this.$store.dispatch('ignoreVersion', item);
    },
    updateAll() {
      this.$store.dispatch('updateAll');
    },
  },
};
</script>
