<template>
  <v-hover v-slot="{ hover }">
    <v-tooltip
      v-if="disabled"
      color="background"
      left
    >
      <template #activator="{ on, attrs }">
        <div
          v-bind="attrs"
          class="d-inline-flex align-center justify-center"
          style="height: 30px; width: 30px"
          :class="hover ? backgroundHoverClass : backgroundNormalClass"
          v-on="on"
        >
          <v-icon
            :color="hover ? iconHoverColor : iconNormalColor"
          >
            {{ hover ? hoverIcon : normalIcon }}
          </v-icon>
        </div>
      </template>
      <slot name="tooltip" />
    </v-tooltip>
    <div
      v-else
      class="d-inline-flex align-center justify-center"
      style="height: 30px; width: 30px"
      :class="hover ? backgroundHoverClass : backgroundNormalClass"
      @click="disabled ? null : buttonClicked()"
    >
      <v-icon
        :color="hover ? iconHoverColor : iconNormalColor"
      >
        {{ hover ? hoverIcon : normalIcon }}
      </v-icon>
    </div>
  </v-hover>
</template>

<script>
export default {
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
    backgroundNormalClass: {
      type: String,
      default: '',
    },
    backgroundHoverClass: {
      type: String,
      default: '',
    },
    iconNormalColor: {
      type: String,
      default: '',
    },
    iconHoverColor: {
      type: String,
      default: '',
    },
    normalIcon: {
      type: String,
      default: '',
    },
    hoverIcon: {
      type: String,
      default: '',
    },
  },
  methods: {
    buttonClicked() {
      this.$emit('click');
    },
  },
};
</script>
