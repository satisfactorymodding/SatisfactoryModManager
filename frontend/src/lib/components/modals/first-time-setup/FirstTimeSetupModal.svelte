<script lang="ts">
  import { mdiClock, mdiPlay, mdiTimerSandFull, mdiWeb } from '@mdi/js';
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
    <ul class="list space-y-4">
      <li>
        <span class="badge bg-tertiary-500 self-end">
          <SvgIcon
            class="h-6 w-6 my-1"
            icon={mdiTimerSandFull}
          />
        </span>
        <div class="flex-auto">
          <p class="text-lg"><T defaultValue="When I add or remove a mod, or switch profiles..." keyName="first_time_setup.option.queue-auto-start.title" /></p>
          <div class="btn-group bg-surface-200-700-token h-10 w-full">
            <button
              class="!btn-sm !px-4 text-lg w-1/2 !justify-between"
              class:!bg-primary-900={$queueAutoStart}
              on:click={() => { $queueAutoStart = true; }}
            >
              <span><T defaultValue="Apply changes immediately" keyName="first_time_setup.option.queue-auto-start.enabled"/></span>
              <SvgIcon
                class="h-5 w-5"
                icon={mdiPlay} />
            </button>
            <button
              class="!btn-sm !px-4 text-lg w-1/2 !justify-between"
              class:!bg-primary-900={!$queueAutoStart}
              on:click={() => { $queueAutoStart = false; }}
            >
              <span><T defaultValue='Wait for me to press "Apply"' keyName="first_time_setup.option.queue-auto-start.disabled"/></span>
              <SvgIcon
                class="h-5 w-5"
                icon={mdiClock} />
            </button>
          </div>
        </div>
      </li>
      <li>
        <span class="badge bg-tertiary-500 text-lg self-end">
          <SvgIcon
            class="h-6 w-6 my-1"
            icon={mdiWeb}
          />
        </span>
        <div class="w-full">
          <span class="text-lg"><T defaultValue="Use this language where available:" keyName="first_time_setup.option.language.title" /></span>
          <LanguageSelector />
        </div>
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
