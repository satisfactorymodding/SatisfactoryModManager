<script lang="ts">
  import Button from '@smui/button';
  import Tooltip, { Wrapper } from '@smui/tooltip';

  import SvgIcon from './SVGIcon.svelte';

  import type { ButtonDisplay } from '$lib/utils/responsiveButton';

  export let display: ButtonDisplay;
  export let onClickAction: () => void;
  export let disabled = false;
  export let visible = true;
  let clazz = '';
  export { clazz as class };

</script>

{#if visible}
  <Wrapper>
    <!-- Div required so a tooltip can still be displayed on a disabled button -->
    <div>
      <Button
        on:click={onClickAction}
        disabled={disabled}
        ripple={false}
        variant="text"
        class="min-w-0 w-12 h-12 group {clazz}"
      >
        <SvgIcon icon={display.icon} class="!p-1 !m-0 !w-full !h-full group-hover:!hidden"/>
        <SvgIcon icon={display.iconHover} class="!p-1 !m-0 !w-full !h-full group-hover:!inline-block !hidden"/>
      </Button>
    </div>
    <Tooltip surface$class="max-w-lg text-base">
      {display.tooltip}
      {#if display.tooltipHtml}
        <!-- eslint-disable-next-line svelte/no-at-html-tags -->
        { @html display.tooltipHtml }
      {/if}
    </Tooltip>
  </Wrapper>
{:else}
  <div class="min-w-0 w-12 h-12"/>
{/if}
