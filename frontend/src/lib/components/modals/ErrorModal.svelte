<script lang="ts">
  import { T } from "@tolgee/svelte";
  import ErrorDetails from "./ErrorDetails.svelte";

  export let parent: { onClose: () => void };

  export let error: string;
</script>

<div
  style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);"
  class="w-[48rem] card flex flex-col gap-6"
>
  <ErrorDetails {error} {parent}>
    <T slot="title" defaultValue="Something went wrong" keyName="error.title" />
  </ErrorDetails>
</div>
