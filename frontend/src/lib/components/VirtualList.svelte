<script generics="T" lang="ts">
  import { createVirtualizer } from '@tanstack/svelte-virtual';

  // eslint-disable-next-line no-undef
  export let items: T[];
  export let itemHeight: number;
  export let overscan = 5;
  let clazz = '';
  export { clazz as class };
  export let itemClass = '';

  let virtualListEl: HTMLDivElement;

  let virtualizer = createVirtualizer<HTMLDivElement, HTMLDivElement>({
    count: items.length,
    getScrollElement: () => virtualListEl,
    estimateSize: () => itemHeight ?? 100,
    overscan,
  });

  let virtualizerUpdateWrapper = virtualizer;

  function updateVirtualizer() {
    $virtualizer.setOptions({
      count: items.length,
      getScrollElement: () => virtualListEl,
      estimateSize: () => itemHeight ?? 100,
      overscan,
    });
    virtualizerUpdateWrapper = virtualizer;
  }

  $: {
    // Don't subscribe to the store here, so we don't call setOptions on scroll (the store pushes a new value on scroll)
    items;
    itemHeight;
    overscan;
    virtualListEl;
    updateVirtualizer();
  }

  $: vitems = $virtualizerUpdateWrapper.getVirtualItems();

  function measure(el: HTMLDivElement) {
    $virtualizerUpdateWrapper.measureElement(el);
  }

  let virtualItemEls: HTMLDivElement[] = [];
  $: if (virtualItemEls.length) {
    virtualItemEls.forEach(measure);
  }
</script>

<div
  bind:this={virtualListEl}
  style="overflow-anchor: none"
  class="relative overflow-y-scroll h-full {clazz}"
>
  <div
    style="height: {$virtualizerUpdateWrapper.getTotalSize()}px; width: 100%;"
    class="overflow-hidden"
  >
    <div
      style="transform: translateY({vitems[0]
        ? vitems[0].start - $virtualizerUpdateWrapper.options.scrollMargin
        : 0}px);"
      class="absolute top-0 left-0 w-full overflow-hidden"
    >
      {#each vitems as row, idx (row.index)}
        <div
          bind:this={virtualItemEls[idx]}
          class="overflow-hidden {itemClass}"
          data-index={row.index}>
          <slot index={row.index} item={items[row.index]}>Missing template</slot>
        </div>
      {/each}
    </div>
  </div>
</div>
