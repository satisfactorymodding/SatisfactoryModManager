<template>
  <div class="container-fluid" style="overflow: auto">
    <div
      v-for="item in objects"
      :key="item.id"
      v-on:click="clicked(item)"
      :class="'row ' + (item == selected ? 'selected' : '')"
    >
      <slot :item="item"></slot>
    </div>
  </div>
</template>

<script>
export default {
  name: 'list',
  data () {
    return {
      selected: null
    }
  },
  components: {
  },
  methods: {
    clicked (item) {
      if (this.canSelect) {
        this.selected = item
        this.$emit('selectionChanged', this.selected)
      }
    }
  },
  props: ['objects', 'canSelect'],
  created () {
  }
}
</script>

<style>
.selected {
  background-color: lightgray;
}
</style>
