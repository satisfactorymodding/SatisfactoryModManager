<script generics="T" lang="ts">
  import type { SizeOptions } from '@floating-ui/dom';
  import { mdiMenuDown } from '@mdi/js';
  import { ListBox, ListBoxItem, focusTrap } from '@skeletonlabs/skeleton';
  import { createEventDispatcher, tick } from 'svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { type PopupSettings, popup } from '$lib/skeletonExtensions';

  export let name: string;
  // eslint-disable-next-line no-undef
  export let items: T[];
  // eslint-disable-next-line no-undef
  export let value: T;
  export let disabled = false;
  let clazz = '';
  export { clazz as class };
  export let buttonClass = '';
  export let menuClass = '';
  export let itemClass = '';
  export let itemActiveClass = '';
  
  // eslint-disable-next-line no-undef
  export let itemKey: ((item: T) => string) | keyof T = (item) => typeof item === 'string' ? item : JSON.stringify(item);

  // eslint-disable-next-line no-undef
  function getKey(item: T) {
    return typeof itemKey === 'function' ? itemKey(item) : item[itemKey];
  }

  let comboboxOpen = false;
  $: combobox = {
    event: 'click',
    target: name,
    placement: 'bottom-start',
    middleware: {
      offset: 6,
      size: {
        apply({ availableHeight, elements }) {
          Object.assign(elements.floating.style, {
            maxHeight: `calc(${availableHeight}px - 1rem)`,
          });
        },
      } as SizeOptions,
      shift: {
        padding: 0,
      },
    },
    state: ({ state }) => comboboxOpen = state,
    closeQuery: `.${name}-listbox-item`,
  } as PopupSettings;

  let selectMenu: HTMLDivElement;
  $: if (comboboxOpen) {
    const selectedItem = selectMenu.querySelector('[aria-selected=true].listbox-item') as HTMLDivElement | null;
    const firstItem = selectMenu.querySelector('.listbox-item') as HTMLDivElement | null;
    tick().then(() => {
      if (selectedItem) {
        selectedItem.focus();
      } else {
        firstItem?.focus();
      }
    });
  }

  // eslint-disable-next-line no-undef
  const dispatch = createEventDispatcher<{ change: T }>();

  function onChange({ target }: Event) {
    const newValue = (target as HTMLButtonElement).value;
    value = items.find((item) => getKey(item) === newValue)!;
    dispatch('change', value);
  }
</script>

<div class="{clazz} relative">
  <div class="h-full w-full" use:popup={combobox}>
    <!--
    The button scale down on click animation changes its bounds, which the popup uses to position itself.
    Wrap button in a div so that the trigger node location does not change.    
    -->
    <button class="btn w-full h-full {buttonClass}" {disabled}>
      <div class="flex-auto text-start justify-start flex min-w-0 space-x-1 overflow-hidden">
        {#if $$slots.selected}
          <slot name="selected" item={value} />
        {:else}
          <slot name="item" item={value}>
            {value}
          </slot>
        {/if}
      </div>
      <SvgIcon
        class="h-5 w-5 p-0.5 {comboboxOpen ? 'text-primary-600 -scale-y-100' : ''} transition-all shrink-0"
        icon={mdiMenuDown} />
    </button>
  </div>

  <div
    bind:this={selectMenu}
    class="card w-full shadow-xl z-10 duration-0 overflow-y-auto !mt-0 {menuClass}"
    data-popup={name}
    use:focusTrap={comboboxOpen}>
    <!-- 
    Skeleton's popup close function waits for the tranistion duration...
    before actually triggering the transition...
    So we'll just not have a transition...
    -->
    <ListBox
      class="w-full"
      regionDefault="overflow-x-hidden"
      rounded="rounded-none"
      spacing="space-y-0">
      {#each items as item}
        <ListBoxItem
          {name}
          class="{name}-listbox-item {itemClass}"
          active="{itemActiveClass}"
          group={getKey(value)}
          value={getKey(item)}
          on:change={onChange}>
          <slot name="item" {item}>
            {item}
          </slot>
          <slot name="itemTrail" slot="trail" {item} />
        </ListBoxItem>
      {/each}
    </ListBox>
  </div>
</div>
