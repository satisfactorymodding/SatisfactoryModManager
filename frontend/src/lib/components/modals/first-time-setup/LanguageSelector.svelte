<script lang="ts">
  import Select from '$lib/components/Select.svelte';
  import { i18n } from '$lib/generated';
  import { language } from '$lib/store/settingsStore';

  let languages: string[] = Object.keys(i18n);

  type DropdownItem = {
    value: string;
    label: string;
  };

  function prepLocale(locale: string) {
    return {
      value: locale,
      label: localeName(locale),
    };
  }

  let dropdownItemData = languages.map(prepLocale);
  
  function localeName(locale: string) {
    if (!locale) return 'N/A';
    return new Intl.DisplayNames([locale], { type: 'language' }).of(locale) ?? `Error:${locale}`;
  }

  async function languageSelectChanged(selection: CustomEvent<DropdownItem>) {
    $language = selection.detail.value;
  }
</script>

<!-- TODO the formatting on this is totally botched, isn't going vertical, is starting open? -->
<Select
  name="languageCombobox"
  class="w-full h-15"
  buttonClass="bg-surface-200-700-token px-4 text-sm"
  itemActiveClass="!bg-surface-300/20"
  itemClass="bg-surface-50-900-token"
  items={dropdownItemData}
  menuClass="bg-surface-50-900-token"

  value={prepLocale($language) ?? 'UNSET'}
  on:change={languageSelectChanged}
>
  <svelte:fragment slot="item" let:item>
    <span>{item.label}</span>
  </svelte:fragment>
  <!-- TODO: dynamic flags 
    <svelte:fragment slot="itemTrail" let:item>
      <span slot="lead" class="h-5 w-5 block">
        {localeFlag(item.value)}
      </span>
    </svelte:fragment>
    -->
</Select>