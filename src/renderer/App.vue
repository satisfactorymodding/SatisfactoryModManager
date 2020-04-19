<template>
  <v-app>
    <v-content>
      <router-view />
    </v-content>
  </v-app>
</template>

<script>
export default {
  name: 'SatisfactoryModLauncherGUI',
};
</script>

<style scoped>
div {
  background: var(--v-background-base);
}
</style>

<style>
.theme--dark.v-card {
  background-color: var(--v-background-base) !important;
}

div {
  color: var(--v-text-base) !important;
}

.v-btn {
  letter-spacing: 0.03em !important;
}

html {
  overflow-y: auto !important;
}

.v-btn {
  text-transform: none !important;
  padding: 0px !important;
}

.v-text-field {
  padding-top: 6px !important;
  margin-top: 12px !important;
}

.v-list-item__icon {
  margin-top: 0px !important;
  margin-bottom: 0px !important;
}

.v-divider {
  border-color: var(--v-background-darken3) !important;
}

.v-label {
  font-size: 12px !important;
}

.v-input {
  font-size: 14px !important;
}

.v-input--selection-controls {
  margin-top: 9px !important;
}

.v-input--switch--inset .v-input--switch__track {
  opacity: 0.95 !important;
}

.custom.theme--dark.v-input--switch.v-input--is-disabled:not(.v-input--is-dirty) .v-input--switch__track {
  color: var(--v-error-base) !important;
}

.theme--dark.v-input--switch .v-input--switch__thumb {
  color: #424242!important;
}

.v-text-field__details, .custom>.v-input__control>.v-messages {
  height: 0 !important;
  min-height: 0 !important;
}

.custom.v-text-field>.v-input__control>.v-input__slot:before {
  border-style: none;
}
.custom.v-text-field>.v-input__control>.v-input__slot:after {
  border-style: none;
}
.custom>.v-input__control>.v-input__slot>.v-select__slot>.v-input__append-inner>.v-input__icon>.v-icon:before {
  color: var(--v-text-base);
}

.v-text-field>.v-input__control>.v-input__slot:before {
  border-color: var(--v-primary-base) !important;
}

.icon {
  font-size: 19px !important;
}

::-webkit-scrollbar {
  width: 20px;
}
::-webkit-scrollbar-track {
  background: var(--v-background-base);
  border: solid 3px transparent;
  background-clip: content-box;
  border-radius: 20px;
}
::-webkit-scrollbar-thumb {
  background: var(--v-backgroundSecondary-base);
  border: solid 6px transparent;
  background-clip: content-box;
  border-radius: 20px;
}
::-webkit-scrollbar-thumb:hover {
  background: white;
  border: solid 6px transparent;
  background-clip: content-box;
}
</style>
