<script lang="ts">
  import { T } from '@tolgee/svelte';

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
    <T defaultValue="Welcome to Satisfactory Mod Manager Version 3!" keyName="smm2_migration.title" />
  </header>
  <section class="px-4">
    <p class="text-base text-center">
      <T
        defaultValue="We've made a lot of improvements in this version. Here are some of the highlights:"
        keyName="smm2_migration.intro"
      />
    </p>
  </section>
  <section class="px-4 overflow-y-auto">
    <dl class="list-dl">
      <div>
        <span class="badge bg-primary-800 text-lg">+</span>
        <span class="flex-auto">
          <dt class="text-lg"><T defaultValue="Dedicated server management" keyName="smm2_migration.feature.servers" /></dt>
        </span>
      </div>
      <div>
        <span class="badge bg-primary-800 text-lg">+</span>
        <span class="flex-auto">
          <dt class="text-lg"><T defaultValue="User interface refresh" keyName="smm2_migration.feature.ui" /></dt>
        </span>
      </div>
      <div>
        <span class="badge bg-primary-800 text-lg">+</span>
        <span class="flex-auto">
          <dt class="text-lg"><T defaultValue="Translation support" keyName="smm2_migration.feature.translation" /></dt>
        </span>
      </div>
      <div>
        <span class="badge bg-primary-800 text-lg">+</span>
        <span class="flex-auto">
          <dt class="text-lg"><T defaultValue="Action queueing" keyName="smm2_migration.feature.queue" /></dt>
        </span>
      </div>
      <div>
        <span class="badge bg-primary-800 text-lg">+</span>
        <span class="flex-auto">
          <dt class="text-lg"><T defaultValue="Native Linux + Mac support" keyName="smm2_migration.feature.platform_support" /></dt>
        </span>
      </div>
      <div>
        <span class="badge bg-tertiary-800 text-lg">+</span>
        <span class="flex-auto">
          <dt class="text-lg"><T defaultValue="Improved performance" keyName="smm2_migration.feature.performance" /></dt>
        </span>
      </div>
      <div>
        <span class="badge bg-warning-800 text-lg">⚠</span>
        <span class="flex-auto">
          <dt class="text-lg"><T defaultValue="New profile format" keyName="smm2_migration.feature.profile_format" /></dt>
        </span>
      </div>
    </dl>
  </section>
  <section class="px-4">
    <p class="text-base text-center">
      <button
        class="btn text-primary-600 variant-ringed"
        on:click={OpenWelcomeGuide}
      >
        <T
          defaultValue="Open the SMM2 -> SMM3 Migration Guide"
          keyName="smm2_migration.open_guide"
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
