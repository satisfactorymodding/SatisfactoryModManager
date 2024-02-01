<script lang="ts">
  // Skeleton passes the parent prop to the modal component, and we would get a warning if the prop is not present here
  export let parent: {onClose: () => void};

  // Just so that it's not unused
  $: parent;

  export let imageSrc: string;
</script>

<div class="card flex flex-col gap-2">
  <img style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" alt="Mod" src={imageSrc}/>
</div>
