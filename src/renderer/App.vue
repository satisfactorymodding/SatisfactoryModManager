<template>
  <div id="app">
    <ul class="nav nav-tabs header">
      <li class="nav-item">
        <router-link
          class="nav-link"
          exact-active-class="active"
          to="/"
        >
          Launcher
        </router-link>
      </li>
    </ul>
    <router-view class="content" />
  </div>
</template>

<script>
export default {
  name: 'SatisfactoryModManagerGUI',
};
</script>

<style>
html,
body,
#app,
main {
  height: 100%;
}
#app,
main {
  display: flex !important;
  flex-flow: column !important;
}
#app .header,
main .header {
  flex: 0 0 auto;
}
#app .content,
main .content {
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
