<template>
  <div id="app">
    <ul class="nav nav-tabs" id="myTab" role="tablist">
      <li class="nav-item">
        <router-link class="nav-link" exact-active-class="active" to="/">Launcher</router-link>
      </li>
      <li class="nav-item">
        <router-link class="nav-link" exact-active-class="active" to="/mods">Download mods</router-link>
      </li>
    </ul>
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  name: 'satisfactory-mod-launcher-gui'
}
</script>

<style>
html,
body,
#app {
  height: 100%;
}
</style>
