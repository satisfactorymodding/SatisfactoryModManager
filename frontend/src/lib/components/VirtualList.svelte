<script generics="T" lang="ts">
  import _ from 'lodash';
  import { tick } from 'svelte';

  // eslint-disable-next-line no-undef
  export let items: T[];
  export let itemHeight: number | undefined = undefined;
  export let bench = 10;
  let clazz = '';
  export { clazz as class };
  export let containerClass = '';

  let start = 0;
  let end = 10;

  let viewport: HTMLElement;
  let container: HTMLElement;
  let heightMap: number[] = [];
  
  $: {
    heightMap = Array.from({ length: items.length });
    if(viewport && container && items.length > 0) {
      tick().then(updateVisible);
    }
  }

  function updateHeightMap() {
    const virtualRows = Array.from(container?.children ?? []);
    virtualRows.forEach((elem, idx) => {
      heightMap[start + idx] = elem.clientHeight;
    });
  }

  $: knownHeights = heightMap.filter((x) => !!x);
  $: averageHeight = knownHeights.reduce((acc, curr) => acc + curr, 0) / knownHeights.length;

  function getHeight(item: number): number {
    return heightMap[item] ?? (itemHeight ?? averageHeight);
  }

  let viewportHeight: number;

  $: if(viewport) {
    // Add or remove elements when the viewport height changes
    viewportHeight;
    updateVisible();
  }

  async function updateVisible() {
    const { scrollTop } = viewport;

    updateHeightMap();
    
    let height = 0;
    let newStart = 0;
    while(newStart < items.length && height + getHeight(newStart) < scrollTop) {
      height += getHeight(newStart);
      newStart++;
    }

    let newEnd = newStart;
    while(newEnd < items.length && height < scrollTop + viewport.clientHeight) {
      height += getHeight(newEnd);
      newEnd++;
    }

    newStart = Math.max(newStart - bench, 0);
    newEnd = Math.min(newEnd + bench, items.length);

    const old_start = start;

    start = newStart;
    end = newEnd;
    
    // prevent jumping if we scrolled up
    // this became an issue when using using skeleton popups as tooltips,
    // specifically when calling window.getComputedStyle(tooltipElement).<anyPropertyHere>
    // I have no idea why that causes this issue
    if (start < old_start) {
      await tick();
      viewport.scrollTo({ top: scrollTop, behavior: 'instant' });
    }
  }

  $: top = _.range(0, start).map(getHeight).reduce((acc, curr) => acc + curr, 0);
  $: bottom = _.range(end, heightMap.length).map(getHeight).reduce((acc, curr) => acc + curr, 0);
  $: visibleItems = items.map((data, index) => ({ index, data })).slice(start, end);

  function itemCreated(_element: HTMLElement) {
    updateHeightMap();
  }
</script>

<div
  bind:this={viewport}
  class="relative overflow-y-scroll h-full {clazz}"
  bind:offsetHeight={viewportHeight}
  on:scroll={updateVisible}
>
  <div
    bind:this={container}
    style="padding-top: {top}px; padding-bottom: {bottom}px;"
    class={containerClass}
  >
    {#each visibleItems as item (item.index)}
      <div class="overflow-hidden" use:itemCreated>
        <slot item={item.data}>Missing template</slot>
      </div>
    {/each}
  </div>
</div>
