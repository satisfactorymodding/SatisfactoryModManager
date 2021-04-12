<template>
  <v-menu
    :close-on-content-click="false"
    offset-x
    left
  >
    <template #activator="{ on: menuOn, attrs: menuAttrs }">
      <v-tooltip bottom>
        <template #activator="{ on: tooltipOn, attrs: tooltipAttrs }">
          <v-btn
            style="min-width: 28px; height: 28px"
            class="my-2 mx-1"
            v-bind="{ ...menuAttrs, ...tooltipAttrs}"
            v-on="{ ...menuOn, ...tooltipOn }"
          >
            <v-icon style="font-size: 16px !important">
              mdi-flip-h mdi-sync
            </v-icon>
          </v-btn>
        </template>
        <span>Update settings</span>
      </v-tooltip>
    </template>
    <v-card>
      <v-list class="menu">
        <v-list-item>
          <v-list-item-action>
            <v-icon color="text">
              mdi-cog
            </v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Update options</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-divider
          insert
          class="custom"
        />

        <v-list-item>
          <v-list-item-action />
          <v-list-item-content>
            <v-list-item-title>Update SMM at</v-list-item-title>
          </v-list-item-content>
          <v-list-item-action>
            <v-select
              v-model="updateCheckModeLocal"
              :items="['launch', 'exit', 'ask']"
              style="width: 108px"
            />
          </v-list-item-action>
        </v-list-item>

        <v-divider
          inset
          class="custom"
        />

        <v-list-item>
          <v-list-item-action />
          <v-list-item-content>
            <v-list-item-title>Show ignored updates</v-list-item-title>
          </v-list-item-content>
          <v-list-item-action>
            <v-switch
              v-model="showIgnoredUpdatesLocal"
            />
          </v-list-item-action>
        </v-list-item>

        <v-divider
          inset
          class="custom"
        />

        <v-list-item>
          <v-list-item-action>
            <v-icon color="text">
              mdi-update
            </v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Updates</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-divider class="custom" />

        <v-list-item
          :disabled="filteredModUpdates.length === 0"
          @click="$emit('openModUpdatesDialog')"
        >
          <v-list-item-action />
          <v-list-item-content>
            <v-list-item-title>Mod updates ({{ filteredModUpdates.length }})</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-divider
          inset
          class="custom"
        />

        <v-list-item
          :disabled="!availableSMMUpdate"
          @click="$emit('openSMMUpdateDialog')"
        >
          <v-list-item-action />
          <v-list-item-content>
            <v-list-item-title>SMM updates ({{ smmUpdateCount }})</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-divider
          inset
          class="custom"
        />
      </v-list>
    </v-card>
  </v-menu>
</template>

<script>
export default {
  props: {
    availableSMMUpdate: {
      type: Object,
      default: () => ({}),
    },
    filteredModUpdates: {
      type: Array,
      default: () => [],
    },
    showIgnoredUpdates: {
      type: Boolean,
      default: false,
    },
    updateCheckMode: {
      type: String,
      default: '',
    },
  },
  computed: {
    hasUpdate() {
      return !!this.availableSMMUpdate || this.filteredModUpdates.length > 0;
    },
    smmUpdateCount() {
      if (!this.availableSMMUpdate) {
        return 0;
      }
      return this.availableSMMUpdate.releaseNotes.length;
    },
    updateCheckModeLocal: {
      get() {
        return this.updateCheckMode;
      },
      set(value) {
        this.$emit('update:updateCheckMode', value);
      },
    },
    showIgnoredUpdatesLocal: {
      get() {
        return this.showIgnoredUpdates;
      },
      set(value) {
        this.$emit('update:showIgnoredUpdates', value);
      },
    },
  },
};
</script>
