<template>
  <main>
    <div class="container-fluid my-2 content d-flex flex-column">
      <div class="row flex-grow-0 flex-shrink-0">
        <div class="col-5">
          <select class="form-control" v-model="selectedSatisfactoryInstall">
            <option
              v-for="install in satisfactoryInstalls.models"
              v-bind:key="install.id"
              v-bind:value="install"
            >{{ install.displayName }}</option>
          </select>
        </div>
        <div class="col-auto d-inline-flex align-items-center">
          <strong>SML: {{ installedSMLVersion }}</strong>
        </div>
        <div class="col-auto">
          <button
            class="btn btn-primary"
            v-if="hasSMLUpdates"
            @click="updateSML"
          >{{ latestSMLVersion.tag_name }} available</button>
          <button class="btn btn-secondary" v-if="isSMLInstalled" @click="uninstallSML">Uninstall</button>
        </div>
        <div class="col-auto">
          <div class="spinner-border" role="status" v-if="SMLInProgress">
            <span class="sr-only">Loading...</span>
          </div>
        </div>
      </div>
      <div class="row my-2 flex-fill">
        <div class="col">
          <list
            :objects="downloadedMods.models.filter(mod => !installedMods.has({mod_id: mod.mod_id, version: mod.version}))"
            :canSelect="true"
            v-model="selectedDownloadedMod"
          >
            <template slot-scope="{item}">
              <div
                class="col-auto d-inline-flex align-items-center"
                @contextmenu.prevent="$refs.modContextMenu.open($event, {mod: item, hasModUpdate: false})"
              >
                <strong>{{item.displayName}}</strong>
              </div>
            </template>
          </list>
        </div>
        <div class="col-3 align-self-center align-items-center">
          <button
            class="btn btn-primary m-2 w-100 d-inline-flex align-items-center justify-content-center"
            style="min-height: 48px"
            @click="installSelectedMod"
          >
            Install &gt;&gt;&nbsp;
            <div class="spinner-border" role="status" v-if="installInProgress">
              <span class="sr-only">Loading...</span>
            </div>
          </button>
          <br />
          <button
            class="btn btn-primary m-2 w-100 d-inline-flex align-items-center justify-content-center"
            style="min-height: 48px"
            @click="uninstallSelectedMod"
          >
            &lt;&lt; Uninstall&nbsp;
            <div class="spinner-border" role="status" v-if="uninstallInProgress">
              <span class="sr-only">Loading...</span>
            </div>
          </button>
        </div>
        <div class="col">
          <list :objects="installedMods.models" :canSelect="true" v-model="selectedInstalledMod">
            <template slot-scope="{item}">
              <div class="col-auto d-inline-flex align-items-center">
                <strong>{{item.displayName}}</strong>
              </div>
            </template>
          </list>
        </div>
      </div>
      <div class="row justify-content-end my-3 mx-1 flex-grow-0 flex-shrink-0">
        <div class="column">
          <button @click="launchSatisfactory" class="btn btn-primary">Launch Satisfactory</button>
        </div>
      </div>
    </div>
    <context-menu
      id="context-menu"
      ref="modContextMenu"
      @ctx-open="onModMenuOpen"
      @ctx-cancel="onModMenuCancel"
      @ctx-close="onModMenuClose"
    >
      <li
        v-if="modContextMenuData.hasModUpdate"
        @click="updateMod(modContextMenuData.mod)"
      >Update mod</li>
      <li @click="removeMod(modContextMenuData.mod)">Remove</li>
    </context-menu>
  </main>
</template>

<script>
import { Collection } from 'vue-mc'
import contextMenu from 'vue-context-menu'
import List from './List'
import SatisfactoryInstall from '../model/satisfactoryInstall'
import ModHandler from '../model/modHandler'
// import FicsitApp from '../model/ficsitApp'
import SMLHandler from '../model/smlHandler'
import semver from 'semver'

export default {
  name: 'launcher',
  components: {
    List,
    contextMenu
  },
  data () {
    return {
      selectedSatisfactoryInstall: null,
      selectedDownloadedMod: null,
      selectedInstalledMod: null,
      satisfactoryInstalls: new Collection(),
      downloadedMods: new Collection(),
      installedMods: new Collection(),
      installedSMLVersion: '',
      latestSMLVersion: new SMLHandler.SMLRelease(),
      SMLInProgress: false,
      installInProgress: false,
      uninstallInProgress: false,
      modContextMenuData: {
        mod: null,
        hasModUpdate: false
      }
    }
  },
  methods: {
    installSelectedMod () {
      if (this.selectedSatisfactoryInstall && this.selectedDownloadedMod) {
        this.installInProgress = true
        ModHandler.installModVersion(this.selectedDownloadedMod, this.selectedSatisfactoryInstall).then(() => {
          this.refreshDownloadedMods()
          this.refreshInstalledMods()
          this.installInProgress = false
        })
      }
    },
    uninstallSelectedMod () {
      if (this.selectedSatisfactoryInstall && this.selectedInstalledMod) {
        this.uninstallInProgress = true
        ModHandler.uninstallModVersion(this.selectedInstalledMod, this.selectedSatisfactoryInstall).then(() => {
          this.refreshDownloadedMods()
          this.refreshInstalledMods()
          this.uninstallInProgress = false
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
    },
    refreshSMLVersion () {
      SMLHandler.getInstalledSMLVersion(this.selectedSatisfactoryInstall).then((version) => {
        this.installedSMLVersion = version
      })
    },
    updateSML () {
      this.SMLInProgress = true
      SMLHandler.updateSML(this.selectedSatisfactoryInstall).then(() => {
        this.SMLInProgress = false
        this.refreshSMLVersion()
      })
    },
    uninstallSML () {
      this.SMLInProgress = true
      SMLHandler.uninstallSML(this.selectedSatisfactoryInstall).then(() => {
        this.SMLInProgress = false
        this.refreshSMLVersion()
      })
    },
    onModMenuOpen (locals) {
      this.modContextMenuData = locals
      ModHandler.hasModUpdate(this.modContextMenuData.mod.mod_id).then((hasUpdate) => {
        this.modContextMenuData.hasModUpdate = hasUpdate
      })
    },
    onModMenuClose (locals) {
    },
    onModMenuCancel () {
      this.modContextMenuData = {
        mod: null,
        hasModUpdate: false
      }
    },
    updateMod (mod) {
      ModHandler.updateMod(mod).then(() => {
        this.refreshDownloadedMods()
      })
    },
    removeMod (mod) {
      ModHandler.removeModVersion(mod).then(() => {
        this.refreshDownloadedMods()
      })
    }
  },
  computed: {
    hasSMLUpdates () {
      return !semver.valid(this.installedSMLVersion) || (semver.valid(this.latestSMLVersion.tag_name) && semver.lt(this.installedSMLVersion, this.latestSMLVersion.tag_name))
    },
    isSMLInstalled () {
      return !!semver.valid(this.installedSMLVersion)
    }
  },
  watch: {
    selectedSatisfactoryInstall: {
      handler: function (newVal, oldVal) {
        this.refreshInstalledMods()
        this.refreshSMLVersion()
      },
      deep: true
    }
  },
  created () {
    this.refreshSatisfactoryInstalls()
    this.refreshDownloadedMods()
    SMLHandler.getLatestSMLVersion().then((version) => {
      this.latestSMLVersion = version
    })
  }
}
</script>

<style>
</style>
