import Vue from 'vue';

import vueElectron from 'vue-electron';
import Vuetify from 'vuetify';
import VueApollo from 'vue-apollo';
import AsyncComputed from 'vue-async-computed';
import App from './App';
import 'vuetify/dist/vuetify.min.css';
import '@mdi/font/css/materialdesignicons.css';
import store from './store';
import { apolloClient } from './graphql';

if (!process.env.IS_WEB) Vue.use(vueElectron);

Vue.config.productionTip = false;
Vue.use(Vuetify);

Vue.use(VueApollo);
Vue.use(AsyncComputed);

/* eslint-disable no-new */
new Vue({
  components: { App },
  apolloProvider: new VueApollo({
    defaultClient: apolloClient,
    defaultOptions: {
      $query: {
        fetchPolicy: 'cache-and-network',
      },
    },
  }),
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
          backgroundModsList: '#2b2b2b',
          backgroundMenu: '#000000',
          backgroundMenuBar: '#0c0c0c',
          text: '#d5d5d5',
          icon: '#9e9e9e',
          text2: '#b7b8bc',
          loadingBackground: '#2f3136',
          ficsitOrange: '#e59445',
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
          ficsitOrange: '#e59445',
        },
      },
      dark: true,
    },
  }),
  store,
  template: '<App/>',
}).$mount('#app');
