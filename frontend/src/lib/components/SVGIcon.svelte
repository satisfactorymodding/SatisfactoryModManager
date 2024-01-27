<script lang="ts">
  export let icon: string;
  let clazz = '';

  export { clazz as class };
</script>

<svg viewBox="0 0 24 24" class="{clazz}">
  <path fill="currentColor" d={icon} />
</svg>
