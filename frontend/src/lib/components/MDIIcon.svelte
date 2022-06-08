<script lang="ts">
  import { classMap } from '@smui/common/internal';
  import { getContext } from 'svelte';


  export let icon: string;
  let clazz = '';

  const context = getContext<string | undefined>('SMUI:icon:context');

  $: fullClass = classMap({
    [clazz]: true,
    'mdc-button__icon': context === 'button',
    'mdc-fab__icon': context === 'fab',
    'mdc-icon-button__icon': context === 'icon-button',
    'mdc-tab__icon': context === 'tab',
    'mdc-banner__icon': context === 'banner',
    'mdc-segmented-button__icon': context === 'segmented-button',
  });

  export {clazz as class};
</script>

<svg viewBox="0 0 24 24" width="24" height="24" class="max-w-full max-h-full {fullClass}">
  <path fill="currentColor" d={icon} />
</svg>
