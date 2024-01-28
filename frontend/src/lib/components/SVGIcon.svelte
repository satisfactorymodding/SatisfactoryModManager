<script lang="ts">
  export let icon: string;
  let clazz = '';

  export { clazz as class };
</script>

<svg class="{clazz}" viewBox="0 0 24 24">
  <path d={icon} fill="currentColor" />
</svg>
