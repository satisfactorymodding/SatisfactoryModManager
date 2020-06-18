<template>
  <v-dialog
    v-model="modUpdatesDialog"
    max-width="600"
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
          @click="updateAll"
        >
          Update all
        </v-btn>
        <v-progress-linear
          v-if="isMultiModUpdateInProgress"
          :value="Math.round(currentMultiModUpdateProgress.progress * 100)"
          :class="currentMultiModUpdateProgress.fast ? 'fast' : ''"
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
                  :value="Math.round(currentModProgress.progress * 100)"
                  :class="currentModProgress.fast ? 'fast' : ''"
                  color="warning"
                  height="49"
                  reactive
                  :indeterminate="currentModProgress.progress < 0"
                >
                  <strong>{{ currentModProgress.message }}</strong>
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
                    @click="viewChangelog"
                  >
                    Changelog
                  </v-btn>
                </v-list-item-action>
                <v-list-item-action>
                  <v-btn
                    color="primary"
                    text
                    :disabled="inProgress.length > 0"
                    @click="updateItem"
                  >
                    Update
                  </v-btn>
                </v-list-item-action>
                <v-list-item-action style="margin-left: 0">
                  <v-btn
                    color="text"
                    text
                    @click="isIgnored ? unignoreUpdate : ignoreUpdate"
                  >
                    {{ isIgnored ? 'Unignore' : 'Ignore' }}
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
import { lastElement } from '../../utils';

export default {
  name: 'ModUpdatesDialog',
  props: {
    filteredModUpdates: {
      type: Array,
      default: () => [],
    },
    ignoreUpdate: {
      type: Function,
      default: () => {},
    },
    inProgress: {
      type: Array,
      default: () => [],
    },
    isIgnored: {
      type: Function,
      default: () => {},
    },
    unignoreUpdate: {
      type: Function,
      default: () => {},
    },
    updateAll: {
      type: Function,
      default: () => {},
    },
    updateItem: {
      type: Function,
      default: () => {},
    },
    viewChangelog: {
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
    .dragregion > span {
        flex-grow: 1;
        margin-top: -3px;
    }

    .button > span {
        flex-grow: 1;
        user-select: none;
    }
</style>
