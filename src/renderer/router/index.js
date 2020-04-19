import Vue from 'vue';
import Router from 'vue-router';

// eslint-disable-next-line import/no-unresolved
const Main = require('@/components/Main');

Vue.use(Router);

export default new Router({
  mode: 'hash',
  routes: [
    {
      path: '/',
      component: Main.default,
    },
    {
      path: '*',
      redirect: '/',
    },
  ],
});
