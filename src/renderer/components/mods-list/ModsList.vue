<template>
  <div
    class="py-4 mods-list"
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
    <v-virtual-scroll
      ref="modsList"
      class="pt-3 ml-5 mr-3 custom"
      :items="mods"
      bench="20"
      item-height="51"
      style="overflow-y: scroll;"
    >
      <template #default="{ item: mod, index }">
        <ModsListItem
          :key="index"
          :mod="mod"
          class="mod-list-item"
        />
      </template>
    </v-virtual-scroll>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import { lastElement, setIntervalImmediately, isCompatibleFast } from '@/utils';
import Fuse from 'fuse.js';
import debounce from 'debounce';
import gql from 'graphql-tag';
import ModsListItem from './ModsListItem';
import { getSetting, saveSetting } from '~/settings';

const MODS_PER_QUERY = 50;

export default {
  components: {
    ModsListItem,
  },
  props: {
    filters: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      isMounted: false,
      mods: [],
      availableMods: [],
      hiddenInstalledMods: [],
      compatibleMods: [],
      availableFilters: [
        {
          name: 'All mods',
          filter(mods) {
            return mods;
          },
          mods: 0,
        },
        {
          name: 'Compatible',
          filter(mods, vm) {
            return vm.compatibleMods;
          },
          mods: 0,
        },
        {
          name: 'Favourite',
          filter(mods, vm) {
            return mods.filter((mod) => vm.$store.state.favoriteModIds.includes(mod.mod_reference));
          },
          mods: 0,
        },
        {
          name: 'Installed',
          filter(mods, vm) {
            return mods.filter((mod) => vm.$store.state.manifestItems.some((item) => item.id === mod.mod_reference));
          },
          mods: 0,
        },
        {
          name: 'Not installed',
          filter(mods, vm) {
            return mods.filter((mod) => !vm.$store.state.manifestItems.some((item) => item.id === mod.mod_reference));
          },
          mods: 0,
        },
        {
          name: 'Enabled',
          filter(mods, vm) {
            const installedMods = Object.keys(vm.$store.state.installedMods);
            return mods.filter((mod) => installedMods.includes(mod.mod_reference));
          },
          mods: 0,
        },
        {
          name: 'Disabled',
          filter(mods, vm) {
            return mods.filter((mod) => {
              const manifestItem = vm.$store.state.manifestItems.find((item) => item.id === mod.mod_reference);
              if (!manifestItem) return false;
              return !manifestItem.enabled;
            });
          },
          mods: 0,
        },
      ],
      availableSorting: [
        {
          name: 'Last updated',
          sort(mods) {
            return mods.sort((a, b) => b.last_version_date - a.last_version_date);
          },
        },
        {
          name: 'Name',
          sort(mods) {
            return mods.sort((a, b) => a.name.localeCompare(b.name));
          },
        },
        {
          name: 'Popularity',
          sort(mods) {
            return mods.sort((a, b) => b.popularity - a.popularity);
          },
        },
        {
          name: 'Hotness',
          sort(mods) {
            return mods.sort((a, b) => b.hotness - a.hotness);
          },
        },
        {
          name: 'Views',
          sort(mods) {
            return mods.sort((a, b) => b.views - a.views);
          },
        },
        {
          name: 'Downloads',
          sort(mods) {
            return mods.sort((a, b) => b.downloads - a.downloads);
          },
        },
      ],
    };
  },
  computed: {
    ...mapState([
      'favoriteModIds',
      'expandedModId',
      'inProgress',
      'installedMods',
      'selectedInstall',
    ]),
    search() {
      return this.filters.search;
    },
    hiddenInstalledModReferences() {
      return Object.keys(this.$store.state.installedMods)
        .filter((modReference) => !this.availableMods.some((mod) => mod.mod_reference === modReference)); // not in the available mods list
    },
    allMods() {
      return [...this.availableMods, ...this.hiddenInstalledMods].filter((mod, idx, arr) => arr.findIndex((other) => other.mod_reference === mod.mod_reference) === idx);
    },
    filteredMods() {
      if (!this.filters.modFilters?.filter || !this.filters.sortBy?.sort) { return this.allMods; }
      // eslint-disable-next-line vue/no-side-effects-in-computed-properties, vue/no-mutating-props
      return this.filters.sortBy.sort(this.filters.modFilters.filter(this.allMods, this));
    },
    topShadow() {
      if (!this.isMounted) return false;
      return (this.$refs.modsList?.first || 0) > 0;
    },
    bottomShadow() {
      if (!this.isMounted) return false;
      return (this.$refs.modsList?.last || 0) < (this.$refs.modsList?.items.length || 0) - 1;
    },
  },
  apollo: {
    hiddenInstalledMods: {
      query: gql`
        query getHiddenInstalledMods($references: [String!]) {
          getMods(filter: { references: $references, hidden: true }) {
            mods {
              id,
              name,
              mod_reference,
              short_description,
              full_description,
              authors {
                user {
                  username,
                }
              },
              downloads,
              views,
              popularity,
              hotness,
              logo,
              hidden,
              last_version_date,
              versions(filter: { limit: 100 }) {
                id,
                sml_version,
              }
            }
          }
        }
      `,
      variables() {
        return {
          references: this.hiddenInstalledModReferences.length > 0 ? this.hiddenInstalledModReferences : [''], // requesting getMods with empty references will return all (hidden) mods
        };
      },
      update: (data) => data.getMods.mods,
    },
  },
  watch: {
    search() {
      this.updateSearch();
    },
    filteredMods() {
      this.updateSearch();
    },
    filters(newValue, oldValue) {
      if (newValue.modFilters !== oldValue.modFilters || newValue.sortBy !== oldValue.sortBy) {
        saveSetting('filters', { modFilters: this.filters.modFilters.name, sortBy: this.filters.sortBy.name });
      }
    },
    selectedInstall() {
      this.updateFilters();
    },
    allMods() {
      this.updateFilters();
    },
    installedMods() {
      this.availableFilters.forEach(async (filter) => { filter.mods = filter.filter(this.allMods, this).length; });
    },
    favoriteModIds() {
      this.availableFilters.forEach(async (filter) => { filter.mods = filter.filter(this.allMods, this).length; });
    },
  },
  mounted() {
    this.$emit('set-available-filters', this.availableFilters);
    this.$emit('set-available-sorting', this.availableSorting);

    this.$nextTick(() => {
      setIntervalImmediately(async () => {
        const availableModsCount = (await this.$apollo.query({
          query: gql`
            query getMods {
              availableMods: getMods {
                count
              }
            }
          `,
        })).data.availableMods.count;
        if (availableModsCount !== this.availableMods.length) {
          const currentLength = this.availableMods.length;
          this.availableMods = (await Promise.all(Array.from({ length: Math.ceil(availableModsCount / MODS_PER_QUERY) }).map(async (_, page) => (await this.$apollo.query({
            query: gql`
              query getMods($limit: Int!, $offset: Int!) {
                availableMods: getMods(filter: { limit: $limit, offset: $offset }) {
                  mods {
                    id,
                    name,
                    mod_reference,
                    short_description,
                    full_description,
                    authors {
                      user {
                        username,
                      }
                    },
                    downloads,
                    views,
                    popularity,
                    hotness,
                    logo,
                    hidden,
                    last_version_date,
                    versions(filter: { limit: 100 }) {
                      id,
                      sml_version,
                    }
                  }
                }
              }
            `,
            variables: {
              limit: MODS_PER_QUERY,
              offset: page * MODS_PER_QUERY,
            },
          })).data.availableMods.mods))).flat(1);
          if (this.$store.state.expandModInfoOnStart && currentLength === 0) {
            this.$store.dispatch('expandMod', this.filteredMods[0].mod_reference);
          }
        }
      }, 5 * 60 * 1000);
    });

    const savedFilters = getSetting('filters', { modFilters: this.availableFilters[1].name, sortBy: this.availableSorting[0].name }); // default Compatible, Last Updated
    this.$emit('update:filters', {
      modFilters: this.availableFilters.find((modFilter) => modFilter.name === savedFilters.modFilters) || this.availableFilters[1], // default Compatible
      sortBy: this.availableSorting.find((item) => item === savedFilters.sortBy) || this.availableSorting[0], // default Last Updated
      search: '',
    });
    this.$root.$on('updateSearch', (search) => {
      this.$emit('update:filters', {
        ...this.filters,
        search,
      });
      this.updateSearch();
    });
    this.updateSearch();
    this.isMounted = true;
  },
  methods: {
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
            name: 'name',
            weight: 2,
          },
          {
            name: 'short_description',
            weight: 1,
          },
          {
            name: 'full_description',
            weight: 0.75,
          },
          {
            name: 'authors.user.username',
            weight: 0.4,
          },
        ],
        useExtendedSearch: true,
        threshold: 0.2,
        ignoreLocation: true,
      });
      this.mods = fuse.search(searchString).map((result) => result.item);
    }),
    async updateFilters() {
      if (!this.$store.state.selectedInstall) {
        this.compatibleMods = [];
        return;
      }
      const { version } = this.$store.state.selectedInstall;
      this.compatibleMods = (await Promise.all(this.allMods.map(async (mod) => ((await isCompatibleFast(mod, version)) ? mod : null)))).filter((mod) => !!mod);
      this.availableFilters.forEach(async (filter) => { filter.mods = filter.filter(this.allMods, this).length; });
    },
    async isCompatible(mod) {
      if (!this.$store.state.selectedInstall) return false;
      return isCompatibleFast(mod, this.$store.state.selectedInstall.version);
    },
    lastElement,
  },
};
</script>

<style scoped>
.container {
  background: var(--v-backgroundModsList-base) !important;
}
.mods-list {
  background: var(--v-backgroundModsList-base) !important;
}

.list-shadow-top, .list-shadow-bottom {
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  z-index: 2;
  background: transparent !important;
  pointer-events: none;
}

.list-shadow-top {
  box-shadow: inset 0px 45px 20px -20px rgba(0,0,0,0.3);
}
.list-shadow-bottom {
  box-shadow: inset 0px -45px 20px -20px rgba(0,0,0,0.3);
}

::-webkit-scrollbar {
  width: 25px;
}
::-webkit-scrollbar-track {
  background: black;
  border: solid 10px transparent;
  background-clip: content-box;
  border-radius: 0;
}
::-webkit-scrollbar-thumb {
  background: var(--v-background-base);
  border: solid 10px transparent;
  background-clip: content-box;
  border-radius: 0;
}
::-webkit-scrollbar-thumb:hover {
  border: solid 10px transparent;
  background-clip: content-box;
}
</style>
