<template>
  <v-dialog
    v-model="dialog"
    max-width="400px"
  >
    <v-card>
      <v-card-title>
        <span class="headline">New Config</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-row>
            <v-form
              ref="form"
              lazy-validation
              class="w-100"
            >
              <v-col cols="12">
                <v-text-field
                  v-model="configName"
                  label="Config Name"
                  :rules="inputRules"
                />
              </v-col>
            </v-form>
          </v-row>
        </v-container>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="blue darken-1"
          text
          @click="submit"
        >
          Save
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  props: {
    dialog: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      configName: '',
      isLoading: false,
      inputRules: [
        (v) => !!v || 'Config Name is required',
      ],
    };
  },
  computed: {
    selectedSatisfactoryInstall: {
      get() {
        return this.$store.state.selectedSatisfactoryInstall;
      },
      set(value) {
        this.$store.state.selectedSatisfactoryInstall = value;
      },
    },
    configLoadInProgress() {
      return this.$store.state.configLoadInProgress;
    },
    inProgress() {
      return this.$store.state.inProgress;
    },
  },
  methods: {
    submit() {
      this.$store.dispatch('handleModalNewConfigSubmit', this.configName).then(() => {
        this.dialog = false;
      });
    },
  },
};
</script>
