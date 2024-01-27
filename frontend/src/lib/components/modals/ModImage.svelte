<script lang="ts">
  // eslint-ignore-next-line unused-export-let
  export let parent: {onClose: () => void};

  export let imageSrc: string;
</script>

<div class="card flex flex-col gap-2" style="max-height: calc(100vh - 128px); max-width: calc(100vw - 128px);">
  <img src={imageSrc} alt="Mod" style="max-height: calc(100vh - 128px); max-width: calc(100vw - 128px);"/>
</div>
