<script lang="ts">
  import T from "$lib/components/T.svelte";
  import { BrowserOpenURL } from "$lib/generated/wailsjs/runtime/runtime";
  import { GenerateDebugInfo } from "$wailsjs/go/app/app";
  import { getTranslate } from "@tolgee/svelte";

  export let parent: { onClose: () => void };

  export let fullPageMode: boolean = false;
  let sectionClass: string = fullPageMode ? "p-4" : "px-4";

  let openDiscordTooltipText: string;
  const { t } = getTranslate();
  t.subscribe((getTranslationText) => {
    openDiscordTooltipText = getTranslationText(
      "error.open_modding_discord.must_generate_debug_first",
      "You must generate debug info first",
    );
  });

  let allowOpeningDiscord: boolean = false;

  let onClickGenerateDebugInfo = async () => {
    let didUserSaveFile = await GenerateDebugInfo();
    if (didUserSaveFile) {
      allowOpeningDiscord = true;
      openDiscordTooltipText = "";
    } else {
      alert(
        "Failed to generate and save the debug info file. Did you click the Cancel button? If not, manually check your Satisfactory Mod Manager log files for more information: https://docs.ficsit.app/satisfactory-modding/latest/faq.html#Files_Logs",
      );
    }
  };

  let OpenDiscord = () => {
    BrowserOpenURL("https://discord.ficsit.app/");
  };
  export let error: string;
</script>

<!-- Replace with svelte snippets once in Svelte 5 -->
<header class="card-header font-bold text-2xl text-center">
  <slot name="title" />
  <!-- <T defaultValue="Something went wrong" keyName="error.title" /> -->
</header>
<section class={`${sectionClass} overflow-y-auto`}>
  <p>{error}</p>
</section>
<section class={sectionClass}>
  <p class={fullPageMode ? "text-base text-center" : ""}>
    <T
      defaultValue="Seems wrong? Click the button below to gather logs, then send the generated zip file on the modding discord in #help-using-mods."
      keyName="error.reporting_directions"
    />
  </p>
</section>
<section class={sectionClass}>
  <p class="text-base text-center">
    <button
      class="btn text-primary-600 variant-ringed"
      on:click={onClickGenerateDebugInfo}
    >
      <T
        defaultValue="Generate debug info"
        keyName="error.generate_debug_info"
      />
    </button>
    <button
      class="btn text-primary-600 variant-ringed"
      on:click={OpenDiscord}
      disabled={!allowOpeningDiscord}
      title={openDiscordTooltipText}
    >
      <T
        defaultValue="Open the Modding Discord"
        keyName="error.open_modding_discord"
      />
    </button>
  </p>
</section>
{#if !fullPageMode}
  <footer class="card-footer">
    <button class="btn" on:click={parent.onClose}>
      <T defaultValue="Close" keyName="common.close" />
    </button>
  </footer>
{/if}
