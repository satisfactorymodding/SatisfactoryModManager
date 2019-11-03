<template>
  <div class="container-fluid" style="overflow: auto">
    <div
      v-for="item in objects"
      :key="item.id"
      v-on:click="clicked(item)"
      :class="'row ' + (item == value ? 'selected' : '')"
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
    }
  },
  components: {
  },
  methods: {
    clicked (item) {
      if (this.canSelect) {
        this.$emit('input', item)
      }
    }
  },
  props: ['objects', 'value', 'canSelect'],
  created () {
  }
}
</script>

<style>
.selected {
  background-color: lightgray;
}
</style>
