<script lang="ts">
  import Select from '$lib/components/Select.svelte';
  import { languages } from '$lib/localization';
  import { language } from '$lib/store/settingsStore';

  let languageObject = languages.find((l) => l.languageCode === $language) ?? languages.find((l) => l.languageCode === 'en')!;

  $: $language = languageObject.languageCode;
</script>

<Select
  name="languageCombobox"
  class="w-full"
  buttonClass="bg-surface-200-700-token px-4 text-sm"
  itemActiveClass="!bg-surface-300/20"
  itemClass="bg-surface-50-900-token"
  itemKey="languageCode"
  items={languages}
  menuClass="bg-surface-50-900-token"

  bind:value={languageObject}
>
  <svelte:fragment slot="item" let:item>
    <span>{item.name} ({Math.round(item.completeness * 100)}%)</span>
  </svelte:fragment>
  <!-- TODO: flags
  <svelte:fragment slot="itemTrail" let:item>
    <span slot="lead" class="h-5 w-5 block">
      {item.flag}
    </span>
  </svelte:fragment>
  -->
</Select>