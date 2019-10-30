<template>
  <div id="wrapper">
    <main>
      <div class="container-fluid">
        <div class="row my-2">
          <div class="col">
            <select class="form-control" v-model="selectedSatisfactoryInstall">
              <option
                v-for="install in satisfactoryInstalls.models"
                v-bind:key="install.id"
                v-bind:value="install"
              >{{ install.displayName }}</option>
            </select>
          </div>
        </div>
        <div class="row my-2">
          <div class="col">
            <select size="15" class="form-control w-100" v-model="selectedDownloadedMod">
              <option
                v-for="mod in downloadedMods.models.filter(mod => !installedMods.has({mod_id: mod.mod_id, version: mod.version}))"
                v-bind:key="mod.mod_id"
                v-bind:value="mod"
              >{{ mod.displayName }}</option>
            </select>
          </div>
          <div class="col-auto align-self-center">
            <button class="btn btn-primary m-2 w-100" @click="installSelectedMod">Install &gt;&gt;</button>
            <br />
            <button
              class="btn btn-primary m-2 w-100"
              @click="uninstallSelectedMod"
            >&lt;&lt; Uninstall</button>
          </div>
          <div class="col">
            <select size="15" class="form-control w-100" v-model="selectedInstalledMod">
              <option
                v-for="mod in installedMods.models"
                v-bind:key="mod.mod_id"
                v-bind:value="mod"
              >{{ mod.displayName }}</option>
            </select>
          </div>
        </div>
        <div class="row justify-content-end m-2">
          <div class="column">
            <button @click="launchSatisfactory" class="btn btn-primary">Launch Satisfactory</button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import { Collection } from 'vue-mc'
import SatisfactoryInstall from '../model/satisfactoryInstall'
import ModHandler from '../model/modHandler'

export default {
  name: 'launcher',
  components: {
  },
  data () {
    return {
      selectedSatisfactoryInstall: null,
      selectedDownloadedMod: null,
      selectedInstalledMod: null,
      satisfactoryInstalls: new Collection(),
      downloadedMods: new Collection(),
      installedMods: new Collection()
    }
  },
  methods: {
    installSelectedMod () {
      if (this.selectedSatisfactoryInstall && this.selectedDownloadedMod) {
        ModHandler.installModVersion(this.selectedDownloadedMod, this.selectedSatisfactoryInstall).then(() => {
          this.refreshInstalledMods()
        })
      }
    },
    uninstallSelectedMod () {
      if (this.selectedSatisfactoryInstall && this.selectedInstalledMod) {
        ModHandler.uninstallModVersion(this.selectedInstalledMod, this.selectedSatisfactoryInstall).then(() => {
          this.refreshInstalledMods()
        })
      }
    },
    refreshSatisfactoryInstalls () {
      SatisfactoryInstall.getInstalls().then((installs) => {
        this.satisfactoryInstalls = installs
        this.selectedSatisfactoryInstall = this.satisfactoryInstalls.models[0]
      })
    },
    refreshDownloadedMods () {
      ModHandler.getDownloadedMods().then((mods) => {
        this.downloadedMods = mods
        this.selectedDownloadedMod = this.downloadedMods[0]
      })
    },
    refreshInstalledMods () {
      ModHandler.getInstalledMods(this.selectedSatisfactoryInstall).then((mods) => {
        this.installedMods = mods
        this.selectedInstalledMod = this.installedMods.models[0]
      })
    },
    launchSatisfactory () {
      if (this.selectedSatisfactoryInstall) {
        this.selectedSatisfactoryInstall.launch()
      }
    }
  },
  computed: {
  },
  watch: {
    selectedSatisfactoryInstall: {
      handler: function (newVal, oldVal) {
        this.refreshInstalledMods()
      },
      deep: true
    }
  },
  created () {
    this.refreshSatisfactoryInstalls()
    this.refreshDownloadedMods()
  }
}
</script>

<style>
</style>
