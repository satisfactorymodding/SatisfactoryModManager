<template>
  <div class="container-fluid" style="overflow: auto">
    <div
      v-for="item in objects"
      :key="item.id"
      v-on:click="clicked(item)"
      :class="'row ' + (item == value && canSelect ? 'selected' : '')"
    >
      <slot :item="item"></slot>
    </div>
  </div>
</template>

<script>
import arrayEqual from 'array-equal'
export default {
  name: 'list',
  data () {
    return {
      selectedIndex: 0
    }
  },
  watch: {
    objects: function (newObjects, oldObjects) {
      if (this.canSelect && !arrayEqual(newObjects, oldObjects)) {
        if (this.objects.length > 0) {
          this.clicked(this.objects[Math.min(Math.max(this.selectedIndex, 0), this.objects.length - 1)])
        } else {
          this.clicked(null)
        }
      }
    }
  },
  methods: {
    clicked (item) {
      if (this.canSelect) {
        this.selectedIndex = this.objects.indexOf(item)
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
