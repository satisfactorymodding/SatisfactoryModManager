<script lang="ts">
  import { mdiClock, mdiPlay, mdiTimerSandFull, mdiTimerSandPaused, mdiWatch, mdiWeb } from '@mdi/js';
  import { SlideToggle } from '@skeletonlabs/skeleton';
  import { T } from '@tolgee/svelte';

  import LanguageSelector from './LanguageSelector.svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { queueAutoStart } from '$lib/store/settingsStore';
  import { SetNewUserSetupComplete } from '$wailsjs/go/settings/settings';
  import { BrowserOpenURL } from '$wailsjs/runtime/runtime';

  export let parent: { onClose: () => void };

  const OpenWelcomeGuide = () => {
    BrowserOpenURL('https://docs.ficsit.app/satisfactory-modding/latest/ForUsers/Welcome.html');
  };

  // Modal should be persistent when used because pressing Escape doesn't trigger this
  function onClose() {
    SetNewUserSetupComplete(true);
    parent.onClose();
  }
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
  <section class="px-4 overflow-y-visible">
    <ul class="list">
      <li>
        <span class="badge bg-tertiary-500">
          <SvgIcon
            class="h-6 w-6 my-1"
            icon={mdiTimerSandFull}
          />
        </span>
        <div class="flex-auto">
          <p class="text-lg"><T defaultValue="When I add or remove a mod, or switch profiles..." keyName="first_time_setup.option.queue-auto-start.title" /></p>
          <SlideToggle
            name="slider-queue"
            class="flex-auto" 
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
                  keyName="first_time_setup.option.queue-auto-start.disabled"
                />
              {/if}
            </span>
          </SlideToggle>
          <!-- TODO testing the button group approach -->
          <div class="btn-group bg-surface-200-700-token">
            <button
              class="!btn-sm !px-4 text-lg"
              class:!bg-primary-900={$queueAutoStart}
              on:click={() => { $queueAutoStart = true; }}
            >
              <T defaultValue="Apply the change immediately." keyName="first_time_setup.option.queue-auto-start.enabled"/>
              <div class="grow"/>
              <SvgIcon
                class="h-5 w-5"
                icon={mdiPlay} />
            </button>
            <button
              class="!btn-sm !px-4 text-lg"
              class:!bg-primary-900={!$queueAutoStart}
              on:click={() => { $queueAutoStart = false; }}
            >
              <T defaultValue='Queue the change and wait for me to press "Apply" to enact queued changes.' keyName="first_time_setup.option.queue-auto-start.disabled"/>
              <div class="grow"/>
              <SvgIcon
                class="h-5 w-5"
                icon={mdiClock} />
            </button>
          </div>
        </div>
      </li>
      <li>
        <span class="badge bg-tertiary-500 text-lg">
          <SvgIcon
            class="h-6 w-6 my-1"
            icon={mdiWeb}
          />
        </span>
        <div class="w-full">
          <span class="text-lg"><T defaultValue="Use this language where available:" keyName="first_time_setup.option.language.title" /></span>
          <LanguageSelector />
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
    <button class="btn variant-ringed" on:click={onClose}>
      <T defaultValue="Get Started!" keyName="first_time_setup.acknowledge" />
    </button>
  </footer>
</div>
