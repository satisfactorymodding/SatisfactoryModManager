<template>
  <v-app>
    <v-main>
      <Main />
    </v-main>
  </v-app>
</template>

<script>
import Main from './components/Main';

export default {
  name: 'SatisfactoryModManager',
  components: {
    Main,
  },
};
</script>

<style scoped>
.container {
  background: var(--v-background-base);
}
</style>

<style>
@import url('https://fonts.googleapis.com/css?family=Open+Sans');
.loading-text-main {
  color: #e6e6e4;
  text-align: center;
  font-family: 'Open Sans';
  font-size: 16px !important;
  font-weight: 600;
  display: block !important;
}

.v-card {
  background-color: var(--v-background-base) !important;
}

.container {
  color: var(--v-text-base) !important;
}

html {
  overflow-y: hidden !important;
}

.v-btn {
  letter-spacing: 0.03em !important;
  text-transform: none !important;
  padding: 0px !important;
}

.v-text-field {
  padding-top: 6px !important;
  margin-top: 6px !important;
}

.v-list-item__icon {
  margin-top: 0px !important;
  margin-bottom: 0px !important;
}

.custom .v-label {
  font-size: 14px !important;
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

.custom.v-input--switch.v-input--is-dirty.v-input--is-disabled {
  opacity: 1;
}

.custom.incompatible.v-input.v-input--switch.v-input--is-disabled:not(.v-input--is-dirty) .v-input--switch__track,
.custom.incompatible.v-input.v-input--switch.v-input--is-disabled .v-input--switch__track {
  color: var(--v-error-base) !important;
}

.custom.v-input--switch .v-input--switch__thumb {
  color: #424242 !important;
}

.custom.v-input.v-input--switch.v-input--is-disabled:not(.v-input--is-dirty) .v-input--switch__track {
  color: hsla(0,0%,100%,.25) !important;
}

.custom.v-input--switch.v-input--dense.v-input--switch--inset .v-input--switch__track {
  height: 15px;
  width: 32px;
  top: calc(50% - 7.5px);
}

.custom.v-input--switch.v-input--dense.v-input--switch--inset .v-input--switch__thumb {
  width: 10px;
  height: 10px;
  top: calc(50% - 4.75px);
}

.custom.v-input--switch.v-input--dense.v-input--switch--inset .v-input--selection-controls__input {
  width: 20px;
}

.custom.v-input--switch.v-input--dense.v-input--switch--inset .v-input--selection-controls__ripple {
  width: 20px;
  height: 20px;
  top: calc(50% - 17px);
}

.custom.v-input--switch.v-input--dense.v-input--is-dirty .v-input--switch__thumb {
  transform: translate(15px) !important;
}

.custom.v-input--switch.v-input--dense.v-input--is-dirty .v-input--selection-controls__ripple {
  transform: translate(15px) !important;
}

.custom.v-list .v-list-item, .custom.v-virtual-scroll .v-list-item {
  padding: 0;
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

.v-progress-linear.fast {
  transition-duration: 0.05s;
}

.v-text-field>.v-input__control>.v-input__slot:before {
  border-color: var(--v-primary-base) !important;
}

.v-icon {
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
  background: var(--v-backgroundModsList-base);
  border: solid 6px transparent;
  background-clip: content-box;
  border-radius: 20px;
}
::-webkit-scrollbar-thumb:hover {
  background: var(--v-text-base);
  border: solid 6px transparent;
  background-clip: content-box;
}
.mod-description pre {
  overflow: auto;
}
pre::-webkit-scrollbar-thumb {
  background: #303136 !important;
  border: solid 6px transparent;
  background-clip: content-box;
  border-radius: 20px;
}
</style>
