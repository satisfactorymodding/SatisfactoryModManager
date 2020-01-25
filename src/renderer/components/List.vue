<template>
  <div
    class="container-fluid"
    style="overflow: auto"
  >
    <div
      v-for="item in objects"
      :key="item.id"
      :class="'row ' + (item == value && canSelect ? 'selected' : '')"
      @click="clicked(item)"
    >
      <slot :item="item" />
    </div>
  </div>
</template>

<script>
import arrayEqual from 'array-equal';

export default {
  name: 'List',
  props: {
    objects: {
      type: Array,
      required: true,
    },
    value: {
      type: Object,
      default() {
        return {};
      },
    },
    canSelect: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      selectedIndex: 0,
    };
  },
  watch: {
    objects(newObjects, oldObjects) {
      if (this.canSelect && !arrayEqual(newObjects, oldObjects)) {
        if (this.objects.length > 0) {
          this.clicked(this.objects[Math.min(Math.max(this.selectedIndex, 0), this.objects.length - 1)]);
        } else {
          this.clicked(null);
        }
      }
    },
  },
  created() {
  },
  methods: {
    clicked(item) {
      if (this.canSelect) {
        this.selectedIndex = this.objects.indexOf(item);
        this.$emit('input', item);
      }
    },
  },
};
</script>

<style>
.selected {
  background-color: lightgray;
}
</style>
