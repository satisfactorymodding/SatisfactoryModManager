<script context="module" lang="ts">
  export interface ButtonDisplay {
    icon: string;
    iconHover: string;
    iconClass?: string;
    iconHoverClass?: string;
    tooltip: string;
    tooltipMarkdown?: string;
  }
</script>

<script lang="ts">
  import Markdown from '$lib/components/Markdown.svelte';
  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import Tooltip from '$lib/components/Tooltip.svelte';
  import { type PopupSettings, popup } from '$lib/skeletonExtensions';

  export let id: string;
  export let display: ButtonDisplay;
  export let onClickAction: () => void;
  export let disabled = false;
  export let visible = true;
  let clazz = '';
  export { clazz as class };
  export let buttonClass = '';

  $: popupId = `responsiveButtonHover-${id}`;

  $: popupHover = {
    event: 'hover',
    target: popupId,
    middleware: {
      offset: 4,
    },
    placement: 'bottom-end',
  } satisfies PopupSettings;
</script>

{#if visible}
  <div class={clazz} use:popup={popupHover}>
    <button
      class="btn-icon min-w-0 group {buttonClass}"
      {disabled}
      on:click={onClickAction}>
      <SvgIcon
        class="!p-1 !m-0 !w-full !h-full group-hover:!hidden {display.iconClass}"
        icon={display.icon}/>
      <SvgIcon
        class="!p-1 !m-0 !w-full !h-full group-hover:!inline-block !hidden {display.iconHoverClass}"
        icon={display.iconHover}/>
    </button>
  </div>
  <Tooltip fixed {popupId}>
    <span>
      {display.tooltip}
      {#if display.tooltipMarkdown}
        <Markdown inline markdown={display.tooltipMarkdown}/>
      {/if}
    </span>
  </Tooltip>
{:else}
  <div class="min-w-0 w-12 h-12"/>
{/if}
