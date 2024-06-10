<script lang="ts">
  import { getContextClient } from '@urql/svelte';

  import ModImage from '$lib/components/modals/ModImage.svelte';
  import { GetModReferenceDocument } from '$lib/generated';
  import { getModalStore } from '$lib/skeletonExtensions';
  import { expandedMod } from '$lib/store/generalStore';
  import { markdown as renderMarkdown } from '$lib/utils/markdown';
  import { BrowserOpenURL } from '$wailsjs/runtime/runtime';

  export let markdown: string;

  $: rendered = renderMarkdown(markdown);

  const modalStore = getModalStore();

  const client = getContextClient();

  // Does not need offline support, since descriptions are disabled in offline mode
  function handleElementClick(element: HTMLElement) {
    if(element instanceof HTMLAnchorElement) {
      const url = new URL(element.href);
      if(url.hostname === 'ficsit.app' && url.pathname.startsWith('/mod/')) {
        const modIdOrReference = url.pathname.split('/')[2];
        if(modIdOrReference) {
          client.query(GetModReferenceDocument, {
            modIdOrReference,
          }).toPromise()
            .then((result) => {
              if (result.data?.getModByIdOrReference?.mod_reference) {
                $expandedMod = result.data.getModByIdOrReference.mod_reference;
              } else {
                console.error(`Failed to GetModReferenceDocument for modIdOrReference '${modIdOrReference}', so opening the link '${element.href}' in the browser instead.`);
                BrowserOpenURL(element.href);
              }
            });
        }
        return true;
      }
      BrowserOpenURL(element.href);
      return true;
    }
    if(element instanceof HTMLImageElement) {
      modalStore.trigger({
        type: 'component',
        component: {
          ref: ModImage,
          props: {
            imageSrc: element.src,
          },
        },
      });
      return true;
    }
    return false;
  }

  function handleDescriptionClick(event: MouseEvent) {
    let element: HTMLElement | null = event.target as HTMLElement;
    while(element) {
      if(handleElementClick(element)) {
        event.preventDefault();
      }
      element = element.parentElement;
    }
  }
</script>

<!-- Intercepting mouse clicks for the link interrupter also seems to work for pressing Enter on the keyboard without a specific key handler added -->
<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
<div
  {...$$props}
  class="markdown-content {$$props.class}"
  role="article"
  on:click={handleDescriptionClick}>
  <!-- eslint-disable-next-line svelte/no-at-html-tags -->
  {@html rendered}
</div>
