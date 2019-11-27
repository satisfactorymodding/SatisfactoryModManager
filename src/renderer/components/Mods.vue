<template>
  <div class="container-fluid my-2">
    <div class="row selection-row" style="height: 50%">
      <div class="col-7 d-flex flex-column">
        <input
          class="form-control flex-shrink-0 flex-grow-0"
          v-model="search"
          type="text"
          placeholder="Search"
          aria-label="Search"
        />
        <br />
        <list
          :objects="searchMods.models"
          :canSelect="true"
          v-model="selectedMod"
          class="flex-fill"
        >
          <template slot-scope="{item}">
            <div class="col-2">
              <img :src="item.logo || noImageURL" width="100%" />
            </div>
            <div class="col-auto d-inline-flex align-items-center">
              <strong>{{item.name}}</strong>
            </div>
          </template>
        </list>
      </div>
      <div class="col-5">
        <list v-if="selectedMod != null" :objects="selectedMod.versions.models" :canSelect="false">
          <template slot-scope="{item}">
            <div class="col-4" style="min-width: 150px">
              <button
                :class="'btn ' + (item.isDownloaded ? 'btn-secondary' : 'btn-primary')"
                v-on:click="toggleModDownloaded(item)"
                style="width: 100%"
              >{{item.isDownloaded ? "Remove" : "Download"}}</button>
            </div>
            <div class="col-auto d-inline-flex align-items-center">
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
    <div class="row" style="overflow: auto; margin: 10px" v-if="selectedMod != null">
      <div v-html="compiledMarkdownDescription"></div>
    </div>
  </div>
</template>

<script>
import FicsitApp from '../model/ficsitApp'
import ModHandler from '../model/modHandler'
import List from './List'
import { Collection } from 'vue-mc'
import marked from 'marked'
export default {
  name: 'download-mods',
  components: {
    List
  },
  computed: {
    noImageURL () {
      return FicsitApp.noImageURL
    },
    compiledMarkdownDescription: function () {
      return marked(this.selectedMod.full_description, { sanitize: true })
    }
  },
  watch: {
    search: function () {
      this.refreshSearch()
    },
    selectedMod: function () {
      this.refreshDownloaded()
    }
  },
  methods: {
    toggleModDownloaded (version) {
      this.isModVersionDownloaded(version).then((downloaded) => {
        version.inProgress = true
        if (downloaded) {
          ModHandler.removeModVersion(version).then(() => {
            version.inProgress = false
            this.refreshDownloaded()
          })
        } else {
          ModHandler.downloadModVersion(version).then(() => {
            version.inProgress = false
            this.refreshDownloaded()
          })
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
      if (this.selectedMod != null && this.selectedMod.versions != null) {
        this.selectedMod.versions.models.forEach((version) => {
          this.isModVersionDownloaded(version).then((downloaded) => {
            version.isDownloaded = downloaded
          })
        })
      }
    },
    refreshSearch () {
      this.searchMods = this.mods.filter((mod) => mod.name.toLowerCase().includes(this.search.toLowerCase()))
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
</style>
