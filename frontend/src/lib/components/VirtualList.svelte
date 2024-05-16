<script generics="T" lang="ts">
  import { createVirtualizer } from '@tanstack/svelte-virtual';
  import _ from 'lodash';

  // eslint-disable-next-line no-undef
  export let items: T[];
  export let itemHeight: number;
  export let overscan = 5;
  let clazz = '';
  export { clazz as class };
  export let itemClass = '';
  
  let virtualListEl: HTMLDivElement;

  $: virtualizer = createVirtualizer<HTMLDivElement, HTMLDivElement>({
    count: items.length,
    getScrollElement: () => virtualListEl,
    estimateSize: () => itemHeight ?? 100,
    overscan,
  });

  $: vitems = $virtualizer.getVirtualItems();
  
  let virtualItemEls: HTMLDivElement[] = [];
  $: if (virtualItemEls.length) {
    virtualItemEls.forEach((el) => $virtualizer.measureElement(el));
  }
</script>

<div
  bind:this={virtualListEl}
  style="overflow-anchor: none"
  class="relative overflow-y-scroll h-full {clazz}"
>
  <div
    style="height: {$virtualizer.getTotalSize()}px; width: 100%;"
    class="overflow-hidden"
  >
    <div
      style="transform: translateY({vitems[0]
        ? vitems[0].start - $virtualizer.options.scrollMargin
        : 0}px);"
      class="absolute top-0 left-0 w-full overflow-hidden"
    >
      {#each vitems as row, idx (row.index)}
        <div
          bind:this={virtualItemEls[idx]}
          class="overflow-hidden {itemClass}"
          data-index={row.index}>
          <slot item={items[row.index]}>Missing template</slot>
        </div>
      {/each}
    </div>
  </div>
</div>
