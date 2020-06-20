<template>
  <v-dialog
    v-model="smmUpdateDialog"
    persistent
  >
    <v-card>
      <v-card-title>
        SMM update available: {{ availableSMMUpdate ? availableSMMUpdate.version : '' }}
      </v-card-title>
      <v-card-text>
        <!-- eslint-disable-next-line vue/no-v-html -->
        <div v-html="smmUpdateNotes" />
      </v-card-text>
      <v-card-actions>
        <v-btn
          color="primary"
          text
          @click="$emit('updateSMMNow')"
        >
          Update now
        </v-btn>
        <v-btn
          color="primary"
          text
          @click="smmUpdateDialog = false"
        >
          Update at exit
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  props: {
    availableSMMUpdate: {
      type: Object,
      default: () => undefined,
    },
    smmUpdateNotes: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      smmUpdateDialog: false,
    };
  },
};
</script>
