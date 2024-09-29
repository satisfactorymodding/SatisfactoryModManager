<script context="module" lang="ts">
  import type { ComponentProps, ComponentType, SvelteComponent } from 'svelte';
  import type { SvelteHTMLElements } from 'svelte/elements';

  export interface TranslationComponentPart<T extends SvelteComponent> {
    component: ComponentType<T>;
    props: ComponentProps<T>;
  }

  export interface TranslationElementPart<T extends keyof SvelteHTMLElements> {
    element: T;
    props: SvelteHTMLElements[T];
  }

  export function translationComponentPart<T extends SvelteComponent>(component: ComponentType<T>, props: ComponentProps<T>): TranslationComponentPart<T> {
    return { component, props };
  }

  export function translationElementPart<T extends keyof SvelteHTMLElements>(element: T, props: SvelteHTMLElements[T]): TranslationElementPart<T> {
    return { element, props };
  }
</script>

<script lang="ts">
  import { type NsType, type TranslateParams, type TranslationKey, getTranslate } from '@tolgee/svelte';
    
  export let keyName: TranslationKey;
  export let params: TranslateParams | undefined = undefined;
  export let noWrap: boolean | undefined = false;
  export let defaultValue: string | undefined = undefined;
  export let ns: NsType | null | undefined = undefined;
  export let language: string | undefined = undefined;

  export let parts: (TranslationComponentPart<SvelteComponent> | TranslationElementPart<keyof SvelteHTMLElements>)[] = [];

  const { t } = getTranslate();
  
  function split(content: string) {
    if (typeof content !== 'string') {
      return [];
    }
    if (!content?.trim()) {
      return [];
    }
  
    const split = content.split(/(<\d+>.+?<\/\d+>)/g);
    const replaced = split.map(part => {
      const match = part.match(/^<(\d+)>(.+?)<\/(\d+)>$/);
      if (match) {
        if (match[1] !== match[3]) {
          console.warn('Translation part tag index does not match: ' + part, match[1], match[3]);
          return undefined;
        }
        const index = parseInt(match[1], 10);
        if (isNaN(index)) {
          console.warn('Translation part index not a number: ' + part);
          return undefined;
        }
        return {
          component: parts[index - 1],
          content: match[2],
        };
      }
      return part;
    });

    if (replaced.some((p) => p === undefined)) {
      return [ content ];
    }
    return replaced as Array<NonNullable<typeof replaced[number]>>;
  }

  // @tolgee-ignore
  $: translated = $t({
    key: keyName,
    params,
    noWrap,
    defaultValue,
    ns,
    language,
  }); 

  $: contentParts = split(translated);
</script>

{#each contentParts as part}
  {#if typeof part === 'string'}
    <span>{part}</span>
  {:else if 'element' in part.component}
    <svelte:element this={part.component.element} {...part.component.props}>{part.content}</svelte:element>
  {:else}
    <svelte:component this={part.component.component} {...part.component.props}>{part.content}</svelte:component>
  {/if}
{/each}