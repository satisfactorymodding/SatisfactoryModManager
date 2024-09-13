<script lang="ts">
  import Select from '$lib/components/Select.svelte';
  import { i18n } from '$lib/generated';
  import { language } from '$lib/store/settingsStore';

  let languages: string[] = Object.keys(i18n);

  function localeName(locale: string) {
    if (!locale) return 'N/A';
    return new Intl.DisplayNames([locale], { type: 'language' }).of(locale) ?? `Error:${locale}`;
  }
</script>

<Select
  name="languageCombobox"
  class="w-full"
  buttonClass="bg-surface-200-700-token px-4 text-sm"
  itemActiveClass="!bg-surface-300/20"
  itemClass="bg-surface-50-900-token"
  items={languages}
  menuClass="bg-surface-50-900-token"

  bind:value={$language}
>
  <svelte:fragment slot="item" let:item>
    <span>{localeName(item)}</span>
  </svelte:fragment>
  <!-- TODO: dynamic flags 
  <svelte:fragment slot="itemTrail" let:item>
    <span slot="lead" class="h-5 w-5 block">
      {localeFlag(item.value)}
    </span>
  </svelte:fragment>
  -->
</Select>