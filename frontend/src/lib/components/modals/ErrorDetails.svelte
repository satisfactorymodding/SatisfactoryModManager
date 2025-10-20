<script lang="ts">
  import { getTranslate } from '@tolgee/svelte';

  import T from '$lib/components/T.svelte';
  import Tooltip from '$lib/components/Tooltip.svelte';
  import { type PopupSettings, popup } from '$lib/skeletonExtensions';
  import { GenerateDebugInfo } from '$wailsjs/go/app/app';
  import { BrowserOpenURL } from '$wailsjs/runtime/runtime';

  export let onClose: (() => void) | null = null;

  export let fullPageMode: boolean = false;
  $: sectionClass = fullPageMode ? 'p-4' : 'px-4';

  const { t } = getTranslate();
  $: openDiscordTooltipText = $t(
    'error.open_modding_discord.must_generate_debug_first',
    'You must generate debug info first',
  );

  let allowOpeningDiscord: boolean = false;
  let debugFileGenerationError: boolean = false;

  let onClickGenerateDebugInfo = async () => {
    try {
      let didUserSaveFile = await GenerateDebugInfo();
      if (didUserSaveFile) {
        // Explicitly set to true -> if people click to save a second time but cancel, don't lock them out
        allowOpeningDiscord = true;
      }
    } catch (error) {
      debugFileGenerationError = true;
      // Enable the Discord button so they can report the error
      allowOpeningDiscord = true;
    }
  };

  let OpenDiscord = () => {
    BrowserOpenURL('https://discord.ficsit.app/');
  };

  let OpenLogDocs = () => {
    BrowserOpenURL('https://docs.ficsit.app/satisfactory-modding/latest/faq.html#Files_Logs');
  };

  const popupId = 'error-details-open-discord-popup';
  const openDiscordPopup = {
    event: 'hover',
    target: popupId,
    middleware: {
      offset: 4,
    },
    placement: 'bottom-end',
  } satisfies PopupSettings;

  export let error: string;
</script>

<header class="card-header font-bold text-2xl text-center">
  <slot name="title" />
</header>
<section class={sectionClass}>
  <p class={fullPageMode ? 'text-base text-center' : ''}>
    <T
      defaultValue="Seems wrong? Click the button below to gather logs, then send the generated zip file on the modding Discord in #help-using-mods."
      keyName="error.reporting_directions"
    />
  </p>
</section>
<section class={`${sectionClass} overflow-y-auto`}>
  <p class="font-mono">{error}</p>
</section>
{#if !allowOpeningDiscord}
  <section class={`${sectionClass} text-center`}>
    <!-- This string intentionally not translated so Fred can always catch it! -->
    <p class="font-mono">Debug info not yet generated</p>
  </section>
{/if}
<section class={sectionClass}>
  <p class={`text-base ${fullPageMode ? 'text-center' : ''}`}>
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
      disabled={!allowOpeningDiscord}
      on:click={OpenDiscord}
      use:popup={openDiscordPopup}
    >
      <T
        defaultValue="Open the Modding Discord"
        keyName="error.open_modding_discord"
      />
    </button>
    <Tooltip disabled={allowOpeningDiscord} {popupId}>
      {openDiscordTooltipText}
    </Tooltip>
  </p>
</section>
{#if debugFileGenerationError}
  <section class={`${sectionClass} ${fullPageMode ? 'text-center' : ''}`}>
    <p class="text-base text-red-500">
      <T
        defaultValue="An error occurred while generating the debug file. Please manually check your Satisfactory Mod Manager log files for more information and report this on the Discord. Use the button below to open the documentation and learn how."
        keyName="error.failed_to_generate_debug"
      />
    </p>
  </section>
  <section class={`${sectionClass} ${fullPageMode ? 'text-center' : ''}`}>
    <button class="btn text-primary-600 variant-ringed" on:click={OpenLogDocs}>
      <T
        defaultValue="Open the Logging Documentation"
        keyName="error.open_log_docs"
      />
    </button>
  </section>
{/if}
{#if !fullPageMode}
  <footer class="card-footer">
    <button class="btn" on:click={onClose}>
      <T defaultValue="Close" keyName="common.close" />
    </button>
  </footer>
{/if}
