<script lang="ts">
  export let loading: boolean;
  export let label: string;
</script>

<div class="inline-flex w-full">
  <span class="mr-2">{label}:</span>
  <span class="font-bold grow" class:animate-pulse={loading} class:placeholder={loading}><slot/></span>
</div>
