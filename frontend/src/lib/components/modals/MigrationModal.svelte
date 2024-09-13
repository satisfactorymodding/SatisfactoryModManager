<script lang="ts">
  import { mdiAlert, mdiPlus } from '@mdi/js';
  import { T } from '@tolgee/svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { MarkSmm2MigrationSuccess } from '$wailsjs/go/migration/migration';
  import { BrowserOpenURL } from '$wailsjs/runtime/runtime';

  export let parent: { onClose: () => void };

  const onClose = () => {
    MarkSmm2MigrationSuccess();
    parent.onClose();
  };

  const OpenSmmGuide = () => {
    BrowserOpenURL('https://docs.ficsit.app/satisfactory-modding/latest/ForUsers/SatisfactoryModManager.html');
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
    <ul class="list-dl list space-y-3">
      <li>
        <span class="badge bg-primary-800 text-lg"><SvgIcon class="w-5 h-5 my-1" icon={mdiPlus} /></span>
        <span class="flex-auto text-lg">
          <T defaultValue="Dedicated server management" keyName="smm2_migration.feature.servers" />
          <p class="text-base">
            <T defaultValue="Easily manage servers via filesystem, SFTP, FTP and SMB!" keyName="smm2_migration.feature.servers.description" />
          </p>
        </span>
      </li>
      <li>
        <span class="badge bg-primary-800 text-lg"><SvgIcon class="w-5 h-5 my-1" icon={mdiPlus} /></span>
        <span class="flex-auto text-lg">
          <T defaultValue="User interface refresh" keyName="smm2_migration.feature.ui" />
          <p class="text-base">
            <T defaultValue="A fresh coat of paint to celebrate Satisfactory 1.0!" keyName="smm2_migration.feature.ui.description" />
          </p>
        </span>
      </li>
      <li>
        <span class="badge bg-primary-800 text-lg"><SvgIcon class="w-5 h-5 my-1" icon={mdiPlus} /></span>
        <span class="flex-auto text-lg">
          <T defaultValue="Translation support" keyName="smm2_migration.feature.translation" />
          <p class="text-base">
            <T defaultValue="Join our Discord to translate SMM to your language!" keyName="smm2_migration.feature.translation.description" />
          </p>
        </span>
      </li>
      <li>
        <span class="badge bg-primary-800 text-lg"><SvgIcon class="w-5 h-5 my-1" icon={mdiPlus} /></span>
        <span class="flex-auto text-lg">
          <T defaultValue="Action queueing" keyName="smm2_migration.feature.queue" />
          <p class="text-base">
            <T defaultValue="Queue up multiple mod downloads/uninstalls at once." keyName="smm2_migration.feature.queue.description" />
          </p>
        </span>
      </li>
      <li>
        <span class="badge bg-primary-800 text-lg"><SvgIcon class="w-5 h-5 my-1" icon={mdiPlus} /></span>
        <span class="flex-auto text-lg">
          <T defaultValue="Native Linux Heroic, Steam Snap, Mac support" keyName="smm2_migration.feature.platform_support" />
          <p class="text-base">
            <T defaultValue="Get a build from the GitHub releases if you haven't yet!" keyName="smm2_migration.feature.platform_support.description" />
          </p>
        </span>
      </li>
      <li>
        <span class="badge bg-tertiary-800 text-lg"><SvgIcon class="w-5 h-5 my-1" icon={mdiPlus} /></span>
        <span class="flex-auto text-lg">
          <T defaultValue="Improved performance" keyName="smm2_migration.feature.performance" />
          <p class="text-base">
            <T defaultValue="Keeping you efficient!" keyName="smm2_migration.feature.performance.description" />
          </p>
        </span>
      </li>
      <li>
        <span class="badge bg-warning-800 text-lg"><SvgIcon class="w-5 h-5 my-1" icon={mdiAlert} /></span>
        <span class="flex-auto text-lg">
          <T defaultValue="New profile format" keyName="smm2_migration.feature.profile_format" />
          <p class="text-base">
            <T defaultValue="SMM2 profiles are not compatible with SMM3! Your existing profiles have automatically been migrated, but you can't Import any profiles created with SMM2." keyName="smm2_migration.feature.profile_format.description" />
          </p>
        </span>
      </li>
    </ul>
  </section>
  <section class="px-4">
    <p class="text-base text-center">
      <button
        class="btn text-primary-600 variant-ringed"
        on:click={OpenSmmGuide}
      >
        <T
          defaultValue="Open the SMM3 Documentation"
          keyName="smm2_migration.open_guide"
        />
      </button>
    </p>
  </section>
  <footer class="card-footer">
    <button class="btn" on:click={onClose}>
      <T defaultValue="Close" keyName="common.close" />
    </button>
  </footer>
</div>
