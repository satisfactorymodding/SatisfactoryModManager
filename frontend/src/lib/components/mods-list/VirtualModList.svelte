<script lang="ts">
  import { tick } from 'svelte';

  import type { PartialMod } from '$lib/store/modFiltersStore';

  export let items: PartialMod[];
  export let itemHeight: number | undefined = undefined;
  export let bench = 10;
  let clazz = '';
  export { clazz as class };

  let start = 0;
  let end = 10;

  let viewport: HTMLElement;
  let container: HTMLElement;
  let heightMap: number[] = [];
  
  $: {
    heightMap = Array.from({length: items.length});
    if(viewport && container && items.length > 0) {
      tick().then(onScroll);
    }
  }

  $: top = heightMap.slice(0, start).reduce((acc, curr) => acc + curr, 0);
  $: bottom = heightMap.slice(end).reduce((acc, curr) => acc + curr, 0);
  $: visibleItems = items.map((data, index) => ({ index, data })).slice(start, end);

  function updateHeightMap() {
    const itemElements = Array.from(container?.children ?? []);
    itemElements.forEach((elem, idx) => {
      heightMap[start + idx] = elem.clientHeight;
    });
    
    const setHeights = heightMap.filter((x) => !!x);
    const averageHeight = setHeights.reduce((acc, curr) => acc + curr, 0) / setHeights.length;
    
    for(let i = 0; i < heightMap.length; i++) {
      heightMap[i] = itemHeight ?? (heightMap[i] ?? averageHeight);
    }
  }

  function onScroll() {
    updateHeightMap();

    
    let height = 0;
    let newStart = 0;
    while(newStart < items.length && height + heightMap[newStart] < viewport.scrollTop) {
      height += heightMap[newStart];
      newStart++;
    }

    let newEnd = newStart;
    while(newEnd < items.length && height < viewport.scrollTop + viewport.clientHeight) {
      height += heightMap[newStart];
      newEnd++;
    }

    newStart = Math.max(newStart - bench, 0);
    newEnd = Math.min(newEnd + bench, items.length);

    start = newStart;
    end = newEnd;
  }
</script>

<style>
  ::-webkit-scrollbar {
    width: 25px;
  }
  ::-webkit-scrollbar-track {
    background: black;
    border: solid 10px transparent;
    border-top-width: 0px;
    border-bottom-width: 0px;
    background-clip: content-box;
    border-radius: 0;
  }
  ::-webkit-scrollbar-thumb {
    background: #fff;
    border: solid 10px transparent;
    border-top-width: 0px;
    border-bottom-width: 0px;
    background-clip: content-box;
    border-radius: 0;
  }
  ::-webkit-scrollbar-thumb:hover {
    border: solid 10px transparent;
    border-top-width: 0px;
    border-bottom-width: 0px;
    background-clip: content-box;
  }
</style>

<div
  bind:this={viewport}
  on:scroll={onScroll}
  class="relative overflow-y-scroll h-full {clazz}"
>
  <div
    bind:this={container}
    style="padding-top: {top}px; padding-bottom: {bottom}px;"
  >
    {#each visibleItems as item (item.index)}
      <div class="overflow-hidden">
        <slot item={item.data}>Missing template</slot>
      </div>
    {/each}
  </div>
</div>