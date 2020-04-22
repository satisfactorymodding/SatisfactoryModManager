<template>
  <div>
    <v-dialog
      v-model="dialog"
      max-width="450"
    >
      <v-alert
        v-if="alert"
        dense
        :type="alert"
        :class="{'mb-0':!alertBody} + ' text-capitalize'"
      >
        {{ alertTitle }}
      </v-alert>
      <v-card v-if="alertBody">
        <v-card-text>
          <p>{{ alertBody }}</p>
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
      alert: null,
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
