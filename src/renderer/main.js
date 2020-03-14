import Vue from 'vue';
import axios from 'axios';

import BootstrapVue from 'bootstrap-vue';
import vueElectron from 'vue-electron';
import App from './App';
import router from './router';

import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-vue/dist/bootstrap-vue.css';

if (!process.env.IS_WEB) Vue.use(vueElectron);

// eslint-disable-next-line no-multi-assign
Vue.http = Vue.prototype.$http = axios;
Vue.config.productionTip = false;
Vue.use(BootstrapVue);

/* eslint-disable no-new */
new Vue({
  components: { App },
  router,
  template: '<App/>',
}).$mount('#app');
