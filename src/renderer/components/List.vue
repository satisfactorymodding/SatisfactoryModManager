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
  watch: {
    objects: function () {
      if (this.objects.length > 0) {
        this.clicked(this.objects[0])
      } else {
        this.clicked(null)
      }
    }
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
