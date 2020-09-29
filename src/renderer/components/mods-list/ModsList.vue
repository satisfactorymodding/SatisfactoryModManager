<template>
  <div
    class="py-4"
    style="position: relative;"
  >
    <div
      v-if="topShadow"
      class="list-shadow-top"
    />
    <div
      v-if="bottomShadow"
      class="list-shadow-bottom"
    />
    <div
      ref="modsList"
      style="overflow-y: scroll; height: 100%;"
      class="ml-4"
      @scroll="modsListScrolled"
    >
      <v-list
        class="pt-3 mt-n4 custom"
      >
        <template
          v-for="(mod, index) in mods"
        >
          <ModsListItem
            :key="index"
            :mod="mod"
          />
        </template>
      </v-list>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex';
import { lastElement } from '@/utils';
import Fuse from 'fuse.js';
import debounce from 'debounce';
import ModsListItem from './ModsListItem';

export default {
  components: {
    ModsListItem,
  },
  data() {
    return {
      topShadow: false,
      bottomShadow: true,
      mods: [],
    };
  },
  computed: {
    ...mapState([
      'favoriteModIds',
      'expandedModId',
      'inProgress',
      'filters',
    ]),
    ...mapGetters({
      filteredMods: 'filteredMods',
      canInstallMods: 'canInstallMods',
    }),
    search() {
      return this.filters.search;
    },
  },
  watch: {
    mods() {
      setTimeout(() => this.modsListScrolled(), 1);
    },
    search() {
      this.updateSearch();
    },
    filteredMods() {
      this.updateSearch();
    },
  },
  mounted() {
    this.updateSearch('');
  },
  methods: {
    modsListScrolled() {
      this.topShadow = this.$refs.modsList.scrollTop > 0;
      this.bottomShadow = this.$refs.modsList.scrollTop + this.$refs.modsList.offsetHeight < this.$refs.modsList.scrollHeight;
    },
    updateSearch: debounce(function updateSearch() {
      let searchString = this.search;
      if (searchString === '') {
        this.mods = this.filteredMods;
        return;
      }
      searchString = searchString.replace(/(?:author:"(.+?)"|author:([^\s"]+))/g, '="$1$2"');

      const fuse = new Fuse(this.filteredMods, {
        keys: [
          {
            name: 'modInfo.name',
            weight: 2,
          },
          {
            name: 'modInfo.short_description',
            weight: 1,
          },
          {
            name: 'modInfo.full_description',
            weight: 0.75,
          },
          {
            name: 'modInfo.authors.user.username',
            weight: 0.4,
          },
        ],
        useExtendedSearch: true,
        threshold: 0.15,
        minMatchCharLength: searchString.length * 0.8,
      });
      this.mods = fuse.search(searchString).map((result) => result.item);
    }),
    lastElement,
  },
};
</script>

<style scoped>
div {
  background: var(--v-backgroundModsList-base) !important;
}

.list-shadow-top, .list-shadow-bottom {
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  z-index: 1;
  background: transparent !important;
  pointer-events: none;
}

.list-shadow-top {
  box-shadow: inset 0px 45px 20px -20px rgba(0,0,0,0.3);
}
.list-shadow-bottom {
  box-shadow: inset 0px -45px 20px -20px rgba(0,0,0,0.3);
}

::-webkit-scrollbar-track {
  background: var(--v-background-base);
  border: solid 9px transparent;
  background-clip: content-box;
  border-radius: 100px;
}
::-webkit-scrollbar-thumb {
  background: var(--v-background-base);
  border: solid 3.5px transparent;
  background-clip: content-box;
  border-radius: 20px;
}
::-webkit-scrollbar-thumb:hover {
  border: solid 3.5px transparent;
  background-clip: content-box;
}
</style>
