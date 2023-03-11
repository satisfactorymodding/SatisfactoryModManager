<script lang="ts">
  import { Icon, Svg } from '@smui/common';

  export let icon: string;
  let clazz = '';

  export { clazz as class };
</script>

<Icon component={Svg} viewBox="0 0 24 24" class="focus:outline-0 {clazz}">
  <path fill="currentColor" d={icon} />
</Icon>