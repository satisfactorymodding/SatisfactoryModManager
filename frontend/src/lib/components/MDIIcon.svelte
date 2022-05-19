<script lang="ts">
  import { Icon } from '@smui/button';
  import { Svg } from '@smui/common/elements';

  export let icon: string;
</script>

<Icon component={Svg} style="width: 1.1em; height: auto;" viewBox="0 0 24 24">
  <path fill="currentColor" d={icon} />
</Icon>
