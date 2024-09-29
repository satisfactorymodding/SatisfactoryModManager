import { i18n } from '$lib/generated';

interface StringTable extends Record<string, null | string | StringTable> {}

function allStrings(stringTable: StringTable): (string | null)[] {
  return Object.values(stringTable).map((entry) => {
    if (!entry) {
      return [null];
    }
    if (typeof entry === 'string') {
      return [entry];
    }
    return allStrings(entry);
  }).flat();
}

export const languages = Object.entries(i18n).map(([lang, stringTable]) => {
  const strings = allStrings(stringTable);
  return {
    languageCode: lang,
    stringTable,
    name: new Intl.DisplayNames([lang], { type: 'language' }).of(lang) ?? lang,
    // TODO flag
    completeness: strings.filter((s) => !!s).length / strings.length,
  };
}).filter((lang) => lang.completeness > 0);