<template>
  <div>
    <div
      :class="imagePage > 0 ? '' : 'hidden'"
      class="images-button left d-inline-flex align-center"
      @click="imagePage = Math.min(imagePage, Math.ceil(images.length / imagesPerColumn) - 1); imagePage -= 1;"
    >
      <v-icon>mdi-chevron-left</v-icon>
    </div>
    <div
      ref="images"
      class="scrollable-images"
      style="height: calc(100% - 351px);"
    >
      <template v-for="(image, i) in images">
        <img
          v-if="image"
          :key="i"
          ref="image"
          class="mod-gallery-image"
          :src="image"
          @click="bigImageSrc = image"
        >
      </template>
    </div>
    <div
      :class="canScrollImagesRight ? '' : 'hidden'"
      class="images-button right d-inline-flex align-center"
      style="right: 0"
      @click="imagePage += 1"
    >
      <v-icon>mdi-chevron-right</v-icon>
    </div>
    <v-dialog
      v-model="showBigImage"
      width="unset"
    >
      <v-card>
        <img
          :src="bigImageSrc"
          style="display: block;"
          :style="`max-height: ${90/100*windowHeight}px`"
        >
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
export default {
  props: {
    images: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      imagePage: 0,
      imagesPerColumn: 1,
      canScrollImagesRight: false,
      bigImageSrc: '',
      windowHeight: 0,
    };
  },
  computed: {
    showBigImage: {
      get() {
        return !!this.bigImageSrc;
      },
      set(value) {
        if (!value) {
          this.bigImageSrc = '';
        }
      },
    },
  },
  watch: {
    imagePage() {
      this.calculatePageLocation();
    },
    images() {
      this.imagePage = 0;
      setTimeout(() => this.calculatePageLocation(), 500);
    },
  },
  created() {
    window.addEventListener('resize', this.onResize);
  },
  mounted() {
    this.calculatePageLocation();
    this.windowHeight = window.innerHeight;
  },
  destroyed() {
    window.removeEventListener('resize', this.onResize);
  },
  methods: {
    calculatePageLocation() {
      if (this.$refs.image && this.$refs.image[0]) {
        let currentWidth = 0;
        this.imagesPerColumn = Math.round(this.$refs.images.clientHeight / this.$refs.image[0].height);
        for (let i = 0; i < this.imagePage && this.$refs.image[i * this.imagesPerColumn]; i += 1) {
          currentWidth += this.$refs.image[i * this.imagesPerColumn].width;
        }
        this.$refs.images.scrollLeft = currentWidth;
        this.canScrollImagesRight = this.$refs.images.scrollWidth - this.$refs.images.clientWidth > currentWidth;
      } else {
        this.canScrollImagesRight = false;
      }
    },
    onResize() {
      this.calculatePageLocation();
      this.windowHeight = window.innerHeight;
    },
  },
};
</script>

<style scoped>
.image-container {
  overflow-x: hidden;
  width: 100%;
}

.image-container img {
  height: 100%;
  display: block;
}

@media (min-height: 850px) {
  .image-container img {
    height: 50%;
  }
}

@media (min-height: 1000px) {
  .image-container img {
    height: 33.33%;
  }
}

@media (min-height: 1500px) {
  .image-container img {
    height: 25%;
  }
}

.images-button {
  text-align: center;
  line-height: 10px;
  position: absolute;
  bottom: 0;
  top: 351px;
  z-index: 1;
}
.images-button.left {
  box-shadow: inset 50px 0px 50px -50px rgba(0, 0, 0, 1);
}

.images-button.right {
  box-shadow: inset -50px 0px 50px -50px rgba(0, 0, 0, 1);
}
.images-button.hidden {
  visibility: hidden;;
}

.scrollable-images {
  width: 100%;
  transition: all ease-in-out 0.5s;
  position: absolute;
  left: 0;
  right: 0;
  display: flex;
  flex-direction: column;
  flex-wrap: wrap;
  overflow-x: hidden;
  scroll-behavior: smooth;
}

.mod-gallery-image {
  cursor: pointer;
}
</style>
