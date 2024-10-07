<script lang="ts">
  export let popupId: string;
  let clazz = '';
  export { clazz as class };
  export let disabled = false;
  export let fixed = false;
</script>

<div
  class="rounded-token z-10 {fixed ? 'fixed' : ''} py-2 px-4 bg-[#0d0d0d] block opacity-0 max-w-[30rem] text-base duration-500 {clazz}"
  class:!hidden={disabled}
  data-popup={popupId}
  inert>
  <!--
    fixed allows the popup to be displayed outside the bounds of the parent
    block opacity-0 ensure that the popup gets a width on first render, otherwise it will be 0px wide and floating-ui will "animate" widening it
    inert is required when popup is invisible to ignore mouse inputs
  -->
  <slot />
</div>
