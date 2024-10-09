<script lang="ts">
  import { fade } from 'svelte/transition';
  import { thumbHashToDataURL } from 'thumbhash';

  export let image: string;
  export let thumbhash: string;
  export let alt: string;

  $: thumbHashData = (() => {
    try {
      return thumbHashToDataURL(
        new Uint8Array(
          atob(thumbhash)
            .split('')
            .map((x) => x.charCodeAt(0)),
        ),
      );
    } catch (e) {
      console.error(e);
    }
  })();

  let imageLoaded = false;
  let thumbnailLoaded = false;

  $: {
    image;
    imageLoaded = false;
  }

  $: {
    thumbhash;
    thumbnailLoaded = false;
  }
</script>
<div class="relative min-h-full min-w-full max-h-full max-w-full">
  <img
    class="absolute max-h-full min-h-full min-w-full max-w-full transition-opacity delay-100 duration-200 ease-linear"
    class:invisible={!imageLoaded}
    class:opacity-0={!imageLoaded}
    {alt}
    {...$$restProps}
    src={image}
    on:load={() => (imageLoaded = true)} />
  {#if !imageLoaded && thumbHashData}
    <img
      class="absolute max-h-full min-h-full min-w-full max-w-full"
      class:invisible={!thumbnailLoaded}
      {alt}
      {...$$restProps}
      src={thumbHashData}
      on:load={() => (thumbnailLoaded = true)}
      in:fade={{ duration: 200 }}
      out:fade={{ duration: 200, delay: 100 }} />
  {/if}
</div>