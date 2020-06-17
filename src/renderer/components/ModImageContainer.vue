<template>
  <v-row
    v-if="!expandDetails"
    class="image-container"
    style="flex-grow: 2"
    no-gutters
  >
    <div
      :class="imagePage > 0 ? '' : 'hidden'"
      class="images-button left d-inline-flex align-center"
      @click="imagePage -= 1"
    >
      <v-icon>mdi-chevron-left</v-icon>
    </div>
    <div
      ref="images"
      class="scrollable-images"
      style="height: calc(100% - 351px);"
    >
      <template v-for="n in images.length">
        <img
          v-if="images[n - 1]"
          :key="n"
          ref="image"
          class="mod-gallery-image"
          :src="images[n - 1]"
          @click="bigImage(n - 1)"
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
  </v-row>
</template>
<script>
export default {
  name: 'ModImageContainer',
  props: {
    bigImage: {
      type: Function,
      default: () => {},
    },
    canScrollImagesRight: {
      type: Boolean,
      default: false,
    },
    expandDetails: {
      type: Boolean,
      default: false,
    },
    imagePage: {
      type: Number,
      default: 0,
    },
    images: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
<style>
    .mod-description img {
        max-width: 100%;
    }

    .v-application .mod-description p {
        margin-bottom: 5px;
    }
</style>
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
        top: 370px;
        z-index: 1;
    }

    .images-button.left {
        box-shadow: inset 50px 0px 50px -50px rgba(0, 0, 0, 1);
    }

    .images-button.right {
        box-shadow: inset -50px 0px 50px -50px rgba(0, 0, 0, 1);
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
