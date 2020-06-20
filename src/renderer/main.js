import Vue from 'vue';

import vueElectron from 'vue-electron';
import Vuetify from 'vuetify';
import App from './App';
import 'vuetify/dist/vuetify.min.css';
import '@mdi/font/css/materialdesignicons.css';
import store from './store';
import { getSetting } from './settings';

if (!process.env.IS_WEB) Vue.use(vueElectron);

Vue.config.productionTip = false;
Vue.use(Vuetify);

/* eslint-disable no-new */
new Vue({
  components: { App },
  vuetify: new Vuetify({
    icons: {
      iconfont: 'mdi',
    },
    theme: {
      options: {
        customProperties: true,
      },
      themes: {
        dark: {
          primary: '#249a20',
          warning: '#ffc107',
          error: '#e51c22',
          info: '#039ae5',
          background: '#212226',
          backgroundSecondary: '#303136',
          menuBackground: '#000000',
          text: '#9e9e9e',
          text2: '#b7b8bc',
          loadingBackground: '#2f3136',
        },
        light: {
          primary: '#249a20',
          warning: '#ffc107',
          error: '#e51c22',
          info: '#039ae5',
          background: '#fcfcfc',
          backgroundSecondary: '#d4d4d4',
          menuBackground: '#ffffff',
          text: '#000',
          text2: '#000',
          loadingBackground: '#ffffff',
        },
      },
      dark: getSetting('darkMode', true),
    },
  }),
  store,
  template: '<App/>',
}).$mount('#app');
