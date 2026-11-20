<script context="module" lang="ts">
  import type { MissingMod, OfflineMod, PartialSMRMod } from '$lib/store/modFiltersStore';

  type OnlinePreview = Pick<PartialSMRMod, 'mod_reference' | 'name' | 'logo' | 'logo_thumbhash' | 'short_description' | 'tags' | 'authors' | 'views' | 'downloads'>;
  export type ModPreviewData = OnlinePreview | OfflineMod | MissingMod;
</script>

<script lang="ts">
  import { mdiDownload, mdiEye, mdiTagMultiple } from '@mdi/js';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import T from '$lib/components/T.svelte';
  import Thumbhash from '$lib/components/Thumbhash.svelte';
  import { siteURL } from '$lib/store/generalStore';
  import { largeNumberFormat } from '$lib/utils/dataFormats';
  import { getAuthor } from '$lib/utils/getModAuthor';

  export let mod: ModPreviewData;
  export let disabled = false;

  $: actualLogo = ('offline' in mod || 'missing' in mod) ? (mod.logo ? `data:image/png;base64, ${mod.logo}` : '/images/no_image.webp') : mod.logo;
  $: renderedLogo = actualLogo || `${$siteURL}/images/no_image.webp`;
  $: renderedThumbhash = ('logo_thumbhash' in mod ? mod.logo_thumbhash : undefined) || '2/eFDQIsFmh9h4BreKeAeQqYBxd3d3J4Jw';
  $: author = ('offline' in mod || 'missing' in mod) ? mod.authors[0] : getAuthor(mod);
</script>

<div class="logo h-full @md/mods-list:w-[5.5rem] w-[4.25rem] shrink-0" class:grayscale={disabled}>
  <Thumbhash alt="{mod.name} Logo" image={renderedLogo} thumbhash={renderedThumbhash} />
</div>
<div class="flex flex-col h-full grow w-0" class:opacity-30={disabled}>
  <slot name="heading">
    <div class="flex items-center">
      <span class="shrink min-w-0 truncate text-lg font-medium !leading-6">{mod.name}</span>
      <div class="shrink-0 hidden @lg/mods-list:block truncate w-[7rem] grow">
        <span class="pl-1"><T defaultValue="by" keyName="mod-list-item.by-author" /></span>
        <span class="text-primary-600 whitespace-nowrap">{author}</span>
      </div>
    </div>
  </slot>
  <div class="flex flex-col grow h-0 overflow-hidden justify-around flex-wrap">
    <div class="truncate w-full @md/mods-list:text-base text-sm">{'short_description' in mod ? mod.short_description : ''}</div>
    <div class="truncate w-full h-5 hidden @md/mods-list:flex items-center space-x-1">
      <SvgIcon class="w-5 shrink-0" icon={mdiTagMultiple} />
      {#if !('offline' in mod) && !('missing' in mod) && (mod.tags?.length ?? 0) > 0}
        {#each mod.tags ?? [] as tag}
          <span>#{tag.name}</span>
        {/each}
      {:else}
        <span><T defaultValue="(none available)" keyName="mod-list-item.no-tags" /></span>
      {/if}
    </div>
    <div class="text-sm w-full">
      <div class="flex h-5 space-x-2">
        {#if !('offline' in mod) && !('missing' in mod)}
          <div class="w-16 flex items-center space-x-1">
            <SvgIcon class="w-4 @md/mods-list:w-5" icon={mdiEye} />
            <span>{largeNumberFormat(mod.views)}</span>
          </div>
          <div class="w-16 flex items-center space-x-1">
            <SvgIcon class="w-4 @md/mods-list:w-5" icon={mdiDownload} />
            <span>{largeNumberFormat(mod.downloads)}</span>
          </div>
        {/if}
      </div>
    </div>
  </div>
  <slot />
</div>
