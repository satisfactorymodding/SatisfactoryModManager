<template>
  <div class="container-fluid">
    <div class="row selection-row">
      <div class="col-7">
        <input
          class="form-control"
          v-model="search"
          type="text"
          placeholder="Search"
          aria-label="Search"
        />
        <br />
        <list
          :objects="searchMods.models"
          :canSelect="true"
          v-on:selectionChanged="selectedModChanged"
        >
          <template slot-scope="{item}">
            <div class="col-2">
              <img :src="item.logo || noImageURL" width="100%" />
            </div>
            <div class="col-auto" style="display: inline-flex; align-items: center;">
              <strong>{{item.name}}</strong>
            </div>
          </template>
        </list>
      </div>
      <div class="col-5">
        <list :objects="selectedMod.versions.models" :canSelect="false">
          <template slot-scope="{item}">
            <div class="col-4" style="min-width: 150px">
              <button
                :class="'btn ' + (item.isDownloaded ? 'btn-secondary' : 'btn-primary')"
                v-on:click="toggleModDownloaded(item)"
                style="width: 100%"
              >{{item.isDownloaded ? "Remove" : "Download"}}</button>
            </div>
            <div class="col-auto" style="display: inline-flex; align-items: center;">
              <strong>{{item.version}}</strong>
            </div>
            <div class="col-auto">
              <div class="spinner-border" role="status" v-if="item.inProgress">
                <span class="sr-only">Loading...</span>
              </div>
            </div>
          </template>
        </list>
      </div>
    </div>
    <div class="row" style="overflow: auto; margin: 10px">
      <vue-markdown :source="selectedMod.full_description"></vue-markdown>
    </div>
  </div>
</template>

<script>
import FicsitApp from '../model/ficsitApp'
import ModHandler from '../model/modHandler'
import List from './List'
import { Collection } from 'vue-mc'
import VueMarkdown from 'vue-markdown'
export default {
  name: 'download-mods',
  components: {
    List,
    VueMarkdown
  },
  computed: {
    noImageURL () {
      return FicsitApp.noImageURL
    }
  },
  watch: {
    search: function () {
      this.refreshSearch()
    }
  },
  methods: {
    selectedModChanged (newMod) {
      this.selectedMod = newMod
      this.refreshDownloaded()
    },
    toggleModDownloaded (version) {
      this.isModVersionDownloaded(version).then((downloaded) => {
        version.inProgress = true
        if (downloaded) {
          ModHandler.removeModVersion(version).then(() => { this.refreshDownloaded() })
        } else {
          ModHandler.downloadModVersion(version).then(() => { this.refreshDownloaded() })
        }
      })
    },
    isModVersionDownloaded (version) {
      return new Promise((resolve, reject) => {
        ModHandler.isModVersionDownloaded(version.mod_id, version.version).then((downloaded) => {
          resolve(downloaded)
        })
      })
    },
    refreshMods () {
      FicsitApp.getMods().then((mods) => {
        this.mods = mods
        this.refreshSearch()
      })
    },
    refreshDownloaded () {
      this.selectedMod.versions.models.forEach((version) => {
        this.isModVersionDownloaded(version).then((downloaded) => {
          version.inProgress = false
          version.isDownloaded = downloaded
        })
      })
    },
    refreshSearch () {
      this.searchMods = this.mods.filter((mod) => mod.name.toLowerCase().includes(this.search))
    }
  },
  data () {
    return {
      mods: new Collection(),
      selectedMod: new FicsitApp.FicsitAppMod(),
      searchMods: new Collection(),
      search: ''
    }
  },
  created () {
    this.refreshMods()
  }
}
</script>

<style>
.container-fluid {
  height: 100%;
  padding-left: 0px;
  padding-right: 0px;
}
.selection-row {
  height: 50%;
  overflow: hidden;
}
</style>
