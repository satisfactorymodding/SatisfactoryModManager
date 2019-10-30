import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
  mode: 'hash',
  routes: [
    {
      path: '/',
      component: require('@/components/Launcher').default
    },
    {
      path: '/mods',
      component: require('@/components/Mods').default
    },
    {
      path: '*',
      redirect: '/'
    }
  ]
})
