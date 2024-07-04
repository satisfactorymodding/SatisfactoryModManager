<script lang="ts">
  import { mdiClose, mdiFilter, mdiSort } from '@mdi/js';
  import { getTranslate } from '@tolgee/svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import Select from '$lib/components/Select.svelte';
  import { type FilterField, type OrderByField, filter, filterOptions, order, orderByOptions, search } from '$lib/store/modFiltersStore';

  const { t } = getTranslate();

  $: orderByNames = {
    name: $t('mods-list-filter.order-by.name', 'Name'),
    'last-updated': $t('mods-list-filter.order-by.last-updated', 'Last updated'),
    popularity: $t('mods-list-filter.order-by.popularity', 'Popularity'),
    hotness: $t('mods-list-filter.order-by.hotness', 'Hotness'),
    views: $t('mods-list-filter.order-by.views', 'Views'),
    downloads: $t('mods-list-filter.order-by.downloads', 'Downloads'),
  } as Record<OrderByField, string>;

  $: filterNames = {
    all: $t('mods-list-filter.filter.all', 'All mods'),
    compatible: $t('mods-list-filter.filter.compatible', 'Compatible'),
    favorite: $t('mods-list-filter.filter.favorite', 'Favorite'),
    queued: $t('mods-list-filter.filter.queued', 'Queued'),
    installed: $t('mods-list-filter.filter.installed', 'Installed'),
    dependency: $t('mods-list-filter.filter.dependency', 'Dependency'),
    'not-installed': $t('mods-list-filter.filter.not-installed', 'Not installed'),
    enabled: $t('mods-list-filter.filter.enabled', 'Enabled'),
    disabled: $t('mods-list-filter.filter.disabled', 'Disabled'),
  } as Record<FilterField, string>;
</script>

<div class="px-5 py-1 flex @container/mod-list-filters items-center">
  <div class="w-0 grow mr-2 inline-flex items-center">
    <input
      class="w-full bg-transparent border-b-[1px] pr-5 border-secondary-500 hover:border-secondary-50 focus:border-primary-600 duration-500 placeholder-secondary-400 placeholder:font-medium outline-none !ring-0 h-8 transition-colors"
      placeholder={$t('mods-list-filter.search', 'Search mods')}
      bind:value={$search}/>
    <button
      class="-ml-5 my-2 pb-0.5 opacity-0 transition-opacity"
      class:!opacity-100={$search}
      class:pointer-events-none={!$search}
      on:click={() => $search = ''}>
      <SvgIcon class="h-5 w-5 text-error-500/80" icon={mdiClose} />
    </button>
  </div>
  <Select
    name="modsFilter"
    class="!h-full"
    buttonClass="px-2 text-sm space-x-1"
    itemActiveClass="!bg-surface-300/20"
    itemClass="bg-surface-50-900-token"
    itemKey="id"
    items={filterOptions}
    menuClass="min-w-[10rem]"
    bind:value={$filter}
  >
    <svelte:fragment slot="selected" let:item>
      <div class="h-5 w-5">
        <SvgIcon icon={mdiFilter} />
      </div>
      <span class="text-primary-600 hidden @lg/mod-list-filters:!block w-24 text-left">{filterNames[item.id]}</span>
    </svelte:fragment>
    <svelte:fragment slot="item" let:item>
      <span>{filterNames[item.id]}</span>
    </svelte:fragment>
  </Select>
  <Select
    name="modsOrderBy"
    class="!h-full"
    buttonClass="px-2 text-sm space-x-1"
    itemActiveClass="!bg-surface-300/20"
    itemClass="bg-surface-50-900-token"
    itemKey="id"
    items={orderByOptions}
    menuClass="min-w-[10rem]"
    bind:value={$order}
  >
    <svelte:fragment slot="selected" let:item>
      <div class="h-5 w-5">
        <SvgIcon icon={mdiSort} />
      </div>
      <span class="text-primary-600 hidden @lg/mod-list-filters:!block w-24 text-left">{orderByNames[item.id]}</span>
    </svelte:fragment>
    <svelte:fragment slot="item" let:item>
      <span>{orderByNames[item.id]}</span>
    </svelte:fragment>
  </Select>
</div>

