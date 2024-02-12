<script lang="ts">
  import { mdiFilter, mdiSort } from '@mdi/js';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import Select from '$lib/components/Select.svelte';
  import { filter, filterOptions, order, orderByOptions, search } from '$lib/store/modFiltersStore';
</script>

<div class="px-5 py-2 flex @container/mod-list-filters">
  <div class="grow mr-2">
    <input class="w-full bg-transparent border-b-[1px] border-secondary-500 hover:border-secondary-50 focus:border-primary-600 duration-500 placeholder-secondary-400 placeholder:font-medium outline-none !ring-0 h-8 transition-colors" placeholder="Search mods" bind:value={$search}/>
  </div>
  <div class="flex grow shrink-0 min-w-[10rem] w-0 @lg/mod-list-filters:min-w-[21rem]">
    <Select
      name="modsFilter"
      class="!h-full w-1/2 shrink-0"
      buttonClass="px-4 text-sm space-x-0.5"
      itemActiveClass="!bg-surface-300/20"
      itemClass="bg-surface-50-900-token"
      itemKey="name"
      items={filterOptions}
      menuClass="min-w-[10rem]"
      bind:value={$filter}
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
      name="modsOrderBy"
      class="!h-full w-1/2 shrink-0"
      buttonClass="px-4 text-sm space-x-0.5"
      itemActiveClass="!bg-surface-300/20"
      itemClass="bg-surface-50-900-token"
      itemKey="name"
      items={orderByOptions}
      menuClass="min-w-[10rem]"
      bind:value={$order}
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

