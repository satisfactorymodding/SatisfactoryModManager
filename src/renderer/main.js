import Vue from 'vue';

import vueElectron from 'vue-electron';
import Vuetify from 'vuetify';
import App from './App';
import 'vuetify/dist/vuetify.min.css';
import '@mdi/font/css/materialdesignicons.css';
import store from './store';

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
          backgroundModsList: '#303136',
          backgroundMenu: '#000000',
          backgroundMenuBar: '#141414',
          text: '#9e9e9e',
          text2: '#b7b8bc',
          loadingBackground: '#2f3136',
        },
        light: {
          primary: '#249a20',
          warning: '#ffc107',
          error: '#e51c22',
          info: '#039ae5',
          background: '#d4d4d4',
          backgroundModsList: '#fcfcfc',
          backgroundMenu: '#ffffff',
          backgroundMenuBar: '#ffffff',
          text: '#000',
          text2: '#000',
          loadingBackground: '#ffffff',
        },
      },
      dark: true,
    },
  }),
  store,
  template: '<App/>',
}).$mount('#app');
