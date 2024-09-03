<script lang="ts">
  import { mdiTimerSandFull, mdiWeb } from '@mdi/js';
  import { SlideToggle } from '@skeletonlabs/skeleton';
  import { T } from '@tolgee/svelte';

  import LanguageSelector from './LanguageSelector.svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { queueAutoStart } from '$lib/store/settingsStore';
  import { BrowserOpenURL } from '$wailsjs/runtime/runtime';

  export let parent: { onClose: () => void };

  const OpenWelcomeGuide = () => {
    BrowserOpenURL('https://docs.ficsit.app/satisfactory-modding/latest/ForUsers/Welcome.html');
  };
</script>

<div
  style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);"
  class="w-[48rem] card flex flex-col gap-6"
>
  <header class="card-header font-bold text-2xl text-center">
    <T defaultValue="Welcome to the Satisfactory Mod Manager!" keyName="first_time_setup.title" />
  </header>
  <section class="px-4">
    <p class="text-base text-center">
      <T
        defaultValue="Select your preferences to get started."
        keyName="first_time_setup.intro"
      />
    </p>
  </section>
  <section class="px-4 overflow-y-auto">
    <ul class="list">
      <li>
        <span class="badge bg-tertiary-500">
          <SvgIcon
            class="h-5 w-5"
            icon={mdiTimerSandFull}
          />
        </span>
        <span class="flex-auto">
          <p class="text-lg"><T defaultValue="When I add or remove a mod, or switch profiles..." keyName="first_time_setup.option.queue-auto-start.title" /></p>
          <p class="text-base">
            <span class="flex-auto">
              <SlideToggle
                name="slider-queue" 
                active="bg-primary-600"
                bind:checked={$queueAutoStart}>
                <span>
                  {#if $queueAutoStart}
                    <T
                      defaultValue="Apply the change immediately."
                      keyName="first_time_setup.option.queue-auto-start.enabled"
                    />
                  {:else}
                    <T
                      defaultValue='Queue the change and wait for me to press "Apply" to enact queued changes.'
                      keyName="first_time_setup.option.queue-auto-start.enabled"
                    />
                  {/if}
                </span>
              </SlideToggle>
            </span>
        </span>
      </li>
      <li>
        <span class="badge bg-tertiary-500 text-lg">
          <SvgIcon
            class="h-5 w-5"
            icon={mdiWeb}
          />
        </span>
        <span class="w-full">
          <p class="text-lg"><T defaultValue="Use this language where available:" keyName="first_time_setup.option.language.title" /></p>
          <p class="text-base">
            <LanguageSelector />
        </span>
      </li>
    </ul>
  </section>
  <section class="px-4">
    <p class="text-base text-center">
      <T
        defaultValue='Change these settings at any time in the "Mod Manager Settings" menu.'
        keyName="first_time_setup.change_later_hint"
      />
    </p>
  </section>
  <section class="px-4">
    <p class="text-base text-center">
      <button
        class="btn text-primary-600 variant-ringed"
        on:click={OpenWelcomeGuide}
      >
        <T
          defaultValue="Open the Welcome Guide"
          keyName="first_time_setup.open_welcome_guide"
        />
      </button>
    </p>
  </section>
  <footer class="card-footer">
    <button class="btn" on:click={parent.onClose}>
      <T defaultValue="Close" keyName="common.close" />
    </button>
  </footer>
</div>
