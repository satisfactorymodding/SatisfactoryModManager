<script lang="ts">
  import { onDestroy } from 'svelte';
  import { sineInOut } from 'svelte/easing';
  import { type Tweened, tweened } from 'svelte/motion';

  let running = false;

  export let animationSpeed = 8;
  let clazz = '';
  export { clazz as class };

  let element: HTMLElement | null = null;

  // I don't know why, but the update function doesn't trigger reactivity, so just update everything at once. I'm looking forward to Svelte 5.
  let scrollableWidth = 0;
  let animationDuration = 0;
  let hoverTranslation: Tweened<number>;

  // widths aren't reactive, so we need to recompute them during the animation, in case the content changes
  function update() {
    const totalWidth = element?.scrollWidth ?? 0;
    const visibleWidth = element?.clientWidth || totalWidth;
    scrollableWidth = totalWidth - visibleWidth;
    animationDuration = scrollableWidth / animationSpeed * 1000;
    if (scrollableWidth < 2) {
      animationDuration = 0;
    }
    hoverTranslation = tweened(0, { duration: animationDuration, easing: sineInOut });
  }

  $: {
    element;
    update();
  }

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
    update();
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