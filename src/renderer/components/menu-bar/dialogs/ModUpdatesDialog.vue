<template>
  <v-dialog
    v-model="modUpdatesDialog"
    max-width="500"
  >
    <v-card>
      <v-card-title>
        Mod updates available
      </v-card-title>
      <v-card-text>
        <v-btn
          color="primary"
          text
          :disabled="inProgress.length > 0"
          @click="$emit('updateAll')"
        >
          Update all
        </v-btn>
        <v-progress-linear
          v-if="isMultiModUpdateInProgress"
          :value="Math.round(currentMultiModUpdateProgress.progress * 100)"
          :class="{ 'fast': currentMultiModUpdateProgress.fast }"
          color="warning"
          height="49"
          reactive
          :indeterminate="currentMultiModUpdateProgress.progress < 0"
        >
          <strong>{{ currentMultiModUpdateProgress.message }}</strong>
        </v-progress-linear>
        <v-list
          v-else
          class="custom"
          dense
        >
          <template v-for="(update, index) in filteredModUpdates">
            <div :key="index">
              <v-list-item
                v-if="inProgress.some((prog) => prog.id === update.item)"
              >
                <v-progress-linear
                  :value="Math.round(currentModProgress(update.item).progress * 100)"
                  :class="{ 'fast': currentModProgress(update.item).fast }"
                  color="warning"
                  height="49"
                  reactive
                  :indeterminate="currentModProgress(update.item).progress < 0"
                >
                  <strong style="font-size: 12px">{{ currentModProgress(update.item).message }}</strong>
                </v-progress-linear>
              </v-list-item>
              <v-list-item v-else>
                <v-list-item-content>
                  <v-list-item-title>{{ update.name }}</v-list-item-title>
                  <v-list-item-subtitle>
                    current: v{{ update.currentVersion }}, available: v{{
                      update.version }}
                  </v-list-item-subtitle>
                </v-list-item-content>
                <v-list-item-action style="margin-left: 10px">
                  <v-btn
                    color="text"
                    text
                    @click="$emit('viewChangelog', update)"
                  >
                    Changelog
                  </v-btn>
                </v-list-item-action>
                <v-list-item-action>
                  <v-btn
                    color="primary"
                    text
                    :disabled="inProgress.length > 0"
                    @click="$emit('updateItem', update)"
                  >
                    Update
                  </v-btn>
                </v-list-item-action>
                <v-list-item-action style="margin-left: 0">
                  <v-btn
                    color="text"
                    text
                    @click="isIgnored(update) ? $emit('unignoreUpdate', update) : $emit('ignoreUpdate', update)"
                  >
                    {{ isIgnored(update) ? 'Unignore' : 'Ignore' }}
                  </v-btn>
                </v-list-item-action>
              </v-list-item>
            </div>
          </template>
        </v-list>
      </v-card-text>
      <v-card-actions>
        <v-btn
          color="primary"
          text
          @click="modUpdatesDialog = false"
        >
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { mapState } from 'vuex';
import { lastElement } from '@/utils';

export default {
  props: {
    filteredModUpdates: {
      type: Array,
      default: () => [],
    },
    ignoredUpdates: {
      type: Array,
      default: () => [],
    },
    isIgnored: {
      type: Function,
      default: () => {},
    },
  },
  data() {
    return {
      modUpdatesDialog: false,
    };
  },
  computed: {
    ...mapState([
      'inProgress',
    ]),
    isMultiModUpdateInProgress() {
      return this.inProgress.some((prog) => prog.id === '__updateMods__');
    },
    multiModUpdateProgress() {
      return this.inProgress.find((prog) => prog.id === '__updateMods__');
    },
    currentMultiModUpdateProgress() {
      return lastElement(this.multiModUpdateProgress.progresses);
    },
  },
  methods: {
    modProgress(mod) {
      return this.inProgress.find((prog) => prog.id === mod);
    },
    currentModProgress(mod) {
      return lastElement(this.modProgress(mod).progresses);
    },
  },
};
</script>

<style scoped>
.custom.v-list {
  background-color: var(--v-background-base);
}
.custom.v-list .v-list-item__action {
  margin: 0;
}
.v-list-item {
  padding-left: 10px !important;
}
.v-list-item__action:first-child {
  margin-right: 0px !important;
}
</style>
