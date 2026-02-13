<script lang="ts">
  import type { SizeOptions } from '@floating-ui/dom';
  import { mdiClose, mdiFilter, mdiSort, mdiTagMultiple } from '@mdi/js';
  import { getTranslate } from '@tolgee/svelte';

  import Marquee from '$lib/components/Marquee.svelte';
  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import Select from '$lib/components/Select.svelte';
  import { type PopupSettings, popup } from '$lib/skeletonExtensions';
  import { type FilterField, type OrderByField, type TagOption, filter, filterOptions, order, orderByOptions, search, selectedTags } from '$lib/store/modFiltersStore';
  import { tagSearchMode } from '$lib/store/settingsStore';

  export let availableTags: TagOption[] = [];
  $: selectedTagIds = new Set($selectedTags.map((t) => t.id));
  $: if ($selectedTags.length > 0) {
    $selectedTags = $selectedTags.filter((t) => availableTags.some((a) => a.id === t.id));
  }

  let _tagPopupOpen = false;
  const tagPopupName = 'modsTagFilter';
  const tagPopup: PopupSettings = {
    event: 'click',
    target: tagPopupName,
    placement: 'bottom-start',
    closeQuery: '', // keep open when clicking tags so user can multi-select
    middleware: {
      offset: 6,
      size: {
        apply({ availableHeight, elements }: { availableHeight: number; elements: { floating: HTMLElement } }) {
          Object.assign(elements.floating.style, {
            maxHeight: `calc(${availableHeight}px - 1rem)`,
          });
        },
      } as SizeOptions,
      shift: { padding: 0 },
    },
    state: ({ state }) => (_tagPopupOpen = state),
  };

  function toggleTag(tag: TagOption) {
    if (selectedTagIds.has(tag.id)) {
      $selectedTags = $selectedTags.filter((t) => t.id !== tag.id);
    } else {
      $selectedTags = [...$selectedTags, tag];
    }
  }

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
  <div class="relative !h-full">
    <div class="h-full w-full" use:popup={tagPopup}>
      <button
        class="btn px-2 text-sm space-x-1 !h-full"
        aria-label={$t('mods-list-filter.tag.button-label', 'Filter by tags')}
        type="button"
        on:contextmenu|preventDefault={() => ($selectedTags = [])}
      >
        <SvgIcon class="h-5 w-5 shrink-0" icon={mdiTagMultiple} />
        {#if $selectedTags.length > 0}
          <span class="text-primary-600 font-medium tabular-nums">{$selectedTags.length}</span>
        {/if}
      </button>
    </div>
    <div
      class="card min-w-[24rem] max-h-96 shadow-xl z-10 duration-0 !mt-0 hidden opacity-0 pointer-events-none inert flex flex-col"
      aria-multiselectable="true"
      data-popup={tagPopupName}
      role="listbox"
    >
      <div
        class="flex items-center gap-2 px-3 py-2 border-b border-surface-400-600-token shrink-0"
        aria-label={$t('mods-list-filter.tag.match-mode', 'Match mode')}
        role="group"
      >
        <button
          class="flex-1 px-3 py-1.5 text-sm rounded transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 {$tagSearchMode === 'and' ? 'text-primary-600 font-medium bg-surface-300/20' : 'text-surface-400-700-token hover:bg-surface-300/20'}"
          aria-pressed={$tagSearchMode === 'and'}
          type="button"
          on:click|stopPropagation={() => tagSearchMode.set('and')}
        >
          {$t('mods-list-filter.tag.match-all', 'Match all')}
        </button>
        <button
          class="flex-1 px-3 py-1.5 text-sm rounded transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 {$tagSearchMode === 'any' ? 'text-primary-600 font-medium bg-surface-300/20' : 'text-surface-400-700-token hover:bg-surface-300/20'}"
          aria-pressed={$tagSearchMode === 'any'}
          type="button"
          on:click|stopPropagation={() => tagSearchMode.set('any')}
        >
          {$t('mods-list-filter.tag.match-any', 'Match any')}
        </button>
      </div>
      <div class="overflow-y-auto min-h-0 flex-1">
        {#if availableTags.length > 0}
          <div class="columns-3 [column-gap:0.5rem] min-h-0 p-2">
            {#each availableTags as tag}
              <button
                class="w-full text-left px-3 py-2 text-sm transition-colors rounded-none {selectedTagIds.has(tag.id) ? 'bg-surface-300/20' : 'bg-surface-50-900-token hover:!bg-surface-300/20'} flex items-center gap-2 break-inside-avoid"
                aria-selected={selectedTagIds.has(tag.id)}
                role="option"
                type="button"
                on:click={() => toggleTag(tag)}
              >
                {#if selectedTagIds.has(tag.id)}
                  <span class="text-primary-600 font-medium" aria-hidden="true">✓</span>
                {/if}
                <span class="{selectedTagIds.has(tag.id) ? 'font-medium' : ''}">{tag.name}</span>
              </button>
            {/each}
          </div>
        {:else}
          <div class="px-3 py-2 text-sm text-surface-400-600-token">
            {$t('mods-list-filter.tag.none-available', 'No tags available')}
          </div>
        {/if}
      </div>
    </div>
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
      <Marquee class="text-primary-600 hidden @lg/mod-list-filters:!block w-24 text-left">{filterNames[item.id]}</Marquee>
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
      <Marquee class="text-primary-600 hidden @lg/mod-list-filters:!block w-24 text-left">{orderByNames[item.id]}</Marquee>
    </svelte:fragment>
    <svelte:fragment slot="item" let:item>
      <span>{orderByNames[item.id]}</span>
    </svelte:fragment>
  </Select>
</div>

