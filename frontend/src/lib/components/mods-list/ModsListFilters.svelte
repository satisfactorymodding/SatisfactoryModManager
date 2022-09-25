<script lang="ts">
  import { filter, filterOptions, order, orderByOptions, search } from '$lib/store/modFiltersStore';
  import Textfield, { Input, type InputComponentDev } from '@smui/textfield';
  import LineRipple, { type LineRippleComponentDev } from '@smui/line-ripple';
  import Select, { Option } from '@smui/select';
  
  let inputA: InputComponentDev;
  let lineRippleA: LineRippleComponentDev;

  export let compact: boolean;
</script>

<div class="px-5 py-2">
  <div class="flex">
    <Textfield
      bind:input={inputA}
      bind:lineRipple={lineRippleA}
      style="height: 30px"
      class="flex-1"
    >
      <Input
        bind:this={inputA}
        bind:value={$search}
        id="input-manual-a"
        aria-controls="helper-text-manual-a"
        aria-describedby="helper-text-manual-a"
        placeholder="Search mods"
      />
      <LineRipple bind:this={lineRippleA} slot="ripple" />
    </Textfield>
    <Select
      class="control-area-input w-1/3 max-w-56 pr-0.5"
      variant="filled"
      bind:value={$filter}
      placeholder="Filter"
      ripple={false}
      selectedText$class="color-primary"
    >
      <span slot="leadingIcon" class="pl-3">
        {#if !compact}
          View:&nbsp;
        {/if}
      </span>
      {#each filterOptions as option}
        <Option value={option}>{option.name}</Option>
      {/each}
    </Select>
    <Select
      class="control-area-input w-1/3 max-w-56 pl-0.5"
      variant="filled"
      bind:value={$order}
      placeholder="Order by"
      ripple={false}
      selectedText$class="color-primary"
    >
      <span slot="leadingIcon" class="pl-3">
        {#if !compact}
          Sort:&nbsp;
        {/if}
      </span>
      {#each orderByOptions as option}
        <Option value={option}>{option.name}</Option>
      {/each}
    </Select>
  </div>
</div>

<style>
  * :global(.control-area-input),
  * :global(.control-area-input .mdc-select__anchor) {
    height: 28px;
  }
  * :global(.control-area-input .mdc-notched-outline),
  * :global(.control-area-input .mdc-select__anchor)  {
    background-color: #1c1c1c;
  }
  * :global(.control-area-input),
  * :global(.control-area-input .mdc-select__anchor) {
    border-radius: 4px;
  }
  * :global(
      .control-area-input.mdc-text-field--with-leading-icon:not(.mdc-text-field--label-floating)
      .mdc-floating-label) {
    left: 16px;
  }
</style>