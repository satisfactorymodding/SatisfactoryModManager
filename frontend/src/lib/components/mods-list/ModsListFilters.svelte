<script lang="ts">
  import Textfield, { Input } from '@smui/textfield';
  import LineRipple from '@smui/line-ripple';
  import Select, { Option } from '@smui/select';
  import { mdiFilter, mdiSort } from '@mdi/js';

  import SvgIcon from '../SVGIcon.svelte';

  import { filter, filterOptions, order, orderByOptions, search } from '$lib/store/modFiltersStore';
  
  let inputA: Input;
  let lineRippleA: LineRipple;
</script>

<div class="px-5 py-2 flex @container/mod-list-filters">
  <div class="grow mr-2">
    <Textfield
      bind:input={inputA}
      bind:lineRipple={lineRippleA}
      class="w-full h-[30px]"
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
  </div>
  <div class="flex grow shrink-0 min-w-[140px] w-0 @lg/mod-list-filters:min-w-[24rem]">
    <Select
      class="control-area-input w-20 @lg/mod-list-filters:w-1/2 pr-0.5"
      variant="filled"
      bind:value={$filter}
      placeholder="Filter"
      ripple={false}
      selectedText$class="color-primary"
      selectedTextContainer$class="!hidden @lg/mod-list-filters:!block"
      dropdownIcon$class="ml-0"
      anchor$class="!items-center w-0"
    >
      <div slot="leadingIcon" class="pl-3 pr-1">
        <div class="h-5 w-5">
          <SvgIcon icon={mdiFilter} />
        </div>
      </div>
      {#each filterOptions as option}
        <Option value={option}>{option.name}</Option>
      {/each}
    </Select>
    <Select
      class="control-area-input w-20 @lg/mod-list-filters:w-1/2 pr-0.5"
      variant="filled"
      bind:value={$order}
      placeholder="Order by"
      ripple={false}
      selectedText$class="color-primary"
      selectedTextContainer$class="!hidden @lg/mod-list-filters:!block"
      dropdownIcon$class="ml-0"
      anchor$class="!items-center w-0"
    >
      <div slot="leadingIcon" class="pl-3 pr-1">
        <div class="h-5 w-5">
          <SvgIcon icon={mdiSort} />
        </div>
      </div>
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
