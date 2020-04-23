<template>
  <div>
    <v-dialog
      v-model="dialog"
      max-width="450"
    >
      <v-alert
        v-if="alert"
        dense
        tile
        flat
        :type="alert"
        class="text-capitalize mb-0"
      >
        {{ alertTitle }}
      </v-alert>
      <v-card
        v-if="alertBody"
        tile
        flat
      >
        <v-card-text class="pt-4">
          <v-label>{{ alertBody }}</v-label>
          <!-- eslint-disable-next-line vue/no-v-html -->
          <div v-html="alertHtmlBody" />
        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
export default {
  data() {
    return {
      dialog: false,
      alert: null, // error, info, warning, success
      alertTitle: null,
      alertBody: null,
      alertHtmlBody: null,
    };
  },
  created() {
    this.$eventBus.$on('popup-alert', (alert, alertTitle, alertBody, alertHtmlBody) => {
      this.alert = alert;
      this.alertTitle = alertTitle ?? alert;
      this.alertBody = alertBody;
      this.alertHtmlBody = alertHtmlBody;
      this.dialog = true;
    });
  },
};
</script>
