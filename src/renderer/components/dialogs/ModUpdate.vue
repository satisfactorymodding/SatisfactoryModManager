<template>
  <v-dialog
    v-model="dialog"
    persistent
    max-width="500px"
  >
    <v-card>
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
                  @click="updateById(item.id)"
                >
                  <v-icon color="green lighten-1">
                    $updateIcon
                  </v-icon>
                </v-btn>
                <v-btn
                  icon
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
          @click="dialog = false"
        >
          Update All
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  data() {
    return {
      dialog: false,
    };
  },
  computed: {
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
  },
};
</script>
