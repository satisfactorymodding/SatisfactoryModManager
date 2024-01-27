<script lang="ts">
  import Textfield, { Input } from '@smui/textfield';
  import LineRipple from '@smui/line-ripple';
  import { mdiFilter, mdiSort } from '@mdi/js';

  import { filter, filterOptions, order, orderByOptions, search } from '$lib/store/modFiltersStore';
  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import Select from '$lib/components/Select.svelte';

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
  <div class="flex grow shrink-0 min-w-[140px] w-0 @lg/mod-list-filters:min-w-[21rem]">
    <Select
      bind:value={$filter}
      items={filterOptions}
      itemKey="name"
      name="modsFilter"
      class="!h-full w-1/2 shrink-0"
      buttonClass="px-4 text-sm space-x-0.5"
      menuClass="min-w-[10rem]"
      itemClass="bg-surface-50-900-token"
      itemActiveClass="!bg-surface-300/20"
    >
      <svelte:fragment slot="selected" let:item>
        <div class="h-5 w-5">
          <SvgIcon icon={mdiFilter} />
        </div>
        <span class="text-primary-600 hidden @lg/mod-list-filters:!block">{item.name}</span>
      </svelte:fragment>
      <svelte:fragment slot="item" let:item>
        <span>{item.name}</span>
      </svelte:fragment>
    </Select>
    <Select
      bind:value={$order}
      items={orderByOptions}
      itemKey="name"
      name="modsOrderBy"
      class="!h-full w-1/2 shrink-0"
      menuClass="min-w-[10rem]"
      buttonClass="px-4 text-sm space-x-0.5"
      itemClass="bg-surface-50-900-token"
      itemActiveClass="!bg-surface-300/20"
    >
      <svelte:fragment slot="selected" let:item>
        <div class="h-5 w-5">
          <SvgIcon icon={mdiSort} />
        </div>
        <span class="text-primary-600 hidden @lg/mod-list-filters:!block">{item.name}</span>
      </svelte:fragment>
      <svelte:fragment slot="item" let:item>
        <span>{item.name}</span>
      </svelte:fragment>
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
