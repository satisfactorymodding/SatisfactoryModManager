<script lang="ts">
  import { onDestroy } from 'svelte';
  import { cubicInOut } from 'svelte/easing';
  import { tweened } from 'svelte/motion';

  let running = false;

  export let animationSpeed = 8;
  let clazz = '';
  export { clazz as class };

  let element: HTMLElement | null = null;
  $: totalWidth = element?.scrollWidth ?? 0;
  $: visibleWidth = element?.clientWidth || totalWidth;
  $: scrollableWidth = totalWidth - visibleWidth;

  $: animationDuration = scrollableWidth / animationSpeed * 1000;

  $: hoverTranslation = tweened(0, { duration: animationDuration, easing: cubicInOut });
  let interval: ReturnType<typeof setInterval> | null = null;

  function stop() {
    if (interval) {
      clearInterval(interval);
      interval = null;
    }
    $hoverTranslation = 0;
  }

  function start() {
    stop();
    if (animationDuration > 0) {
      $hoverTranslation = 1;
      interval = setInterval(() => {
        $hoverTranslation = 1 - $hoverTranslation;
      }, animationDuration);
    }
  }

  onDestroy(stop);

  $: {
    animationDuration;
    if (running) {
      start();
    } else {
      stop();
    }
  }
</script>

<span
  bind:this={element}
  class="max-w-full overflow-hidden {clazz}"
  role="marquee"
  on:mouseover={() => running = true}
  on:mouseout={() => running = false}
  on:focus={() => running = true}
  on:blur={() => running = false}
>
  <span style="transform: translateX(-{$hoverTranslation * scrollableWidth}px);" class="relative inline-block">
    <slot />
  </span>
</span>