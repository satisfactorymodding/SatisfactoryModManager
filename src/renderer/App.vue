<template>
  <div id="app">
    <ul class="nav nav-tabs header">
      <li class="nav-item">
        <router-link class="nav-link" exact-active-class="active" to="/">Launcher</router-link>
      </li>
      <li class="nav-item">
        <router-link class="nav-link" exact-active-class="active" to="/mods">Download mods</router-link>
      </li>
    </ul>
    <router-view class="content"></router-view>
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
#app {
  display: flex;
  flex-flow: column;
}
#app .header {
  flex: 0 0 auto;
}
#app .content {
  flex: 1 1 auto;
}
.container-fluid {
  height: 100%;
  padding-left: 0px;
  padding-right: 0px;
}
.selection-row {
  overflow: hidden;
}
</style>
