<template>
  <v-container
    fluid
  >
    <v-row
      no-gutters
      class="mt-4"
    >
      <v-col cols="12">
        <v-row
          justify="center"
          no-gutters
        >
          <v-col
            cols="6"
            class="px-2"
          >
            <v-select
              v-model="value.config"
              :items="configs"
              item-text="name"
              return-object
              label="SELECT CONFIG"
              class="custom"
              append-icon="mdi-chevron-down"
            />
          </v-col>
          <v-col
            cols="6"
            class="buttons"
          >
            <v-row>
              <v-col
                cols="6"
              >
                <v-btn text>
                  New&nbsp;
                  <v-icon
                    color="green"
                    class="icon"
                  >
                    mdi-plus-circle
                  </v-icon>
                </v-btn>
              </v-col>
              <v-col
                cols="6"
              >
                <v-btn text>
                  Delete&nbsp;
                  <v-icon
                    color="red"
                    class="icon"
                  >
                    mdi-delete
                  </v-icon>
                </v-btn>
              </v-col>
            </v-row>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col
        cols="12"
      >
        <v-row
          justify="center"
          no-gutters
        >
          <v-col
            cols="6"
            class="px-2"
          >
            <v-select
              v-model="value.filters.modFilters"
              :items="modFilters"
              item-text="name"
              :return-object="true"
              label="SHOW"
              class="custom"
              append-icon="mdi-chevron-down"
            >
              <template v-slot:selection="{ item }">
                <span>{{ item.name }}</span>
                &nbsp;
                <span class="green--text">({{ item.mods }})</span>
              </template>
              <template
                slot="item"
                slot-scope="data"
              >
                <!-- HTML that describe how select should render items when the select is open -->
                {{ data.item.name }} ({{ data.item.mods }})
              </template>
            </v-select>
          </v-col>
          <v-col
            cols="6"
            class="px-2"
          >
            <v-select
              v-model="value.filters.sortBy"
              :items="sortBy"
              label="SORT BY"
              class="custom"
              append-icon="mdi-chevron-down"
            />
          </v-col>
        </v-row>
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col
        cols="12"
        class="px-2"
      >
        <v-text-field
          v-model="value.filters.search"
          label="Search"
        />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  props: {
    configs: {
      type: Array,
      default() { return []; },
    },
    modFilters: {
      type: Array,
      default() { return []; },
    },
    sortBy: {
      type: Array,
      default() { return []; },
    },
    value: {
      type: Object,
      default() {
        return {
          config: {},
          filters: {
            compatibility: {},
            sortBy: {},
          },
        };
      },
    },
  },
};
</script>

<style scoped>
.buttons {
  padding-left: 3px !important;
}

.v-btn {
  color: var(--v-text-base) !important;
}
</style>
