<script lang="ts">
  export let icon: string;
  let clazz: string | undefined = undefined;
  export {clazz as class};
</script>

<svg viewBox="0 0 24 24" width="24" height="24" class="max-w-full max-h-full {clazz}">
  <path fill="currentColor" d={icon} />
</svg>
