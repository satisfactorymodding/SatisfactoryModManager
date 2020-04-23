import Vue from 'vue';
import axios from 'axios';

import vueElectron from 'vue-electron';
import Vuetify from 'vuetify';
import App from './App';
import router from './router';
import 'vuetify/dist/vuetify.min.css';
import '@mdi/font/css/materialdesignicons.css';

if (!process.env.IS_WEB) Vue.use(vueElectron);

// eslint-disable-next-line no-multi-assign
Vue.http = Vue.prototype.$http = axios;
Vue.config.productionTip = false;
Vue.use(Vuetify);

/* eslint-disable no-new */
new Vue({
  components: { App },
  router,
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
          text: '#9e9e9e',
          text2: '#b7b8bc',
        },
        light: {
          primary: '#249a20',
          warning: '#ffc107',
          error: '#e51c22',
          info: '#039ae5',
          background: '#fcfcfc',
          backgroundSecondary: '#d4d4d4',
          text: '#000',
          text2: '#000',
        },
      },
      dark: true,
    },
  }),
  template: '<App/>',
}).$mount('#app');
