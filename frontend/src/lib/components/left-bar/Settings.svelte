<script lang="ts">
  import { mdiBug, mdiCheck, mdiChevronRight, mdiClipboard, mdiCog, mdiDownload, mdiFolderEdit, mdiTune } from '@mdi/js';
  import { getContextClient } from '@urql/svelte';
  import { popup, type PopupSettings , ListBox, ListBoxItem } from '@skeletonlabs/skeleton';
  import type { SizeOptions } from '@floating-ui/dom';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { GenerateDebugInfo } from '$wailsjs/go/app/app';
  import { startView, konami, launchButton, queueAutoStart, offline, updateCheckMode } from '$lib/store/settingsStore';
  import { manifestMods, lockfileMods } from '$lib/store/ficsitCLIStore';
  import { GetModNameDocument } from '$lib/generated';
  import type { LaunchButtonType, ViewType } from '$lib/wailsTypesExtensions';
  import { OfflineGetMod } from '$wailsjs/go/ficsitcli/ficsitCLI';
  import { getModalStore } from '$lib/store/skeletonExtensions';

  const modalStore = getModalStore();

  const settingsMenu = {
    event: 'click',
    target: 'settings-menu',
    middleware: {
      offset: 4,
      size: {
        apply({ availableHeight, elements }) {
          Object.assign(elements.floating.style, {
            maxHeight: `${availableHeight * 0.8}px`,
          });
        },
      } as SizeOptions,
    },
    placement: 'right-start',
    closeQuery: '[data-popup="settings-menu"] li:not([data-noclose]):not(.section-header)',
  } satisfies PopupSettings;

  const startViewMenu = {
    event: 'click',
    target: 'start-view-menu',
    middleware: {
      offset: 4,
    },
    placement: 'right-start',
    closeQuery: '[data-popup="start-view-menu"] .listbox-item',
  } satisfies PopupSettings;

  let views: {id: ViewType, name: string}[] = [
    {
      id: 'compact',
      name: 'Compact',
    },
    {
      id: 'expanded',
      name: 'Expanded',
    },
  ];

  const updateCheckModeMenu = {
    event: 'click',
    target: 'update-check-mode-menu',
    middleware: {
      offset: 4,
    },
    placement: 'right-start',
    closeQuery: '[data-popup="update-check-mode-menu"] .listbox-item',
  } satisfies PopupSettings;

  let updateCheckModes: {id: 'launch'|'exit'|'ask', name: string}[] = [
    {
      id: 'launch',
      name: 'On start',
    },
    {
      id: 'exit',
      name: 'On exit',
    },
    {
      id: 'ask',
      name: 'Ask when found',
    },
  ];

  const queueModeMenu = {
    event: 'click',
    target: 'queue-mode-menu',
    middleware: {
      offset: 4,
    },
    placement: 'right-start',
    closeQuery: '[data-popup="queue-mode-menu"] .listbox-item',
  } satisfies PopupSettings;

  let queueModes: {id: boolean, name: string}[] = [
    {
      id: true,
      name: 'Start immediately',
    },
    {
      id: false,
      name: 'Start manually',
    },
  ];

  const launchButtonMenu = {
    event: 'click',
    target: 'launch-button-menu',
    middleware: {
      offset: 4,
    },
    placement: 'right-start',
    closeQuery: '[data-popup="launch-button-menu"] .listbox-item',
  } satisfies PopupSettings;

  let launchButtons: {id: LaunchButtonType, name: string}[] = [
    {
      id: 'normal',
      name: 'Normal',
    },
    {
      id: 'cat',
      name: 'Nyan',
    },
    {
      id: 'button',
      name: 'Launch Button',
    },
  ];
  
  const urqlClient = getContextClient();

  async function copyModList() {
    // Generate mod entries
    const modList = await Promise.all(Object.keys($manifestMods).map(async (modReference) => {
      let modName = modReference;
      if(modReference === 'SML') {
        modName = 'Satisfactory Mod Loader';
      } else {
        if($offline) {
          modName = (await OfflineGetMod(modReference)).name;
        } else {
          const result = await urqlClient.query(GetModNameDocument, { modReference }).toPromise();
          if(result?.data?.getModByReference?.name) {
            modName = result.data.getModByReference.name;
          }
        }
      }
      return {
        friendlyName: modName,
        modReference,
        version: $lockfileMods[modReference].version,
      };
    }));
    // Sort by Friendly Name
    modList.sort((a, b) => {
      const x = a.friendlyName.toLowerCase();
      const y = b.friendlyName.toLowerCase();
      return x.localeCompare(y);
    });
    // Get max lengths to use for padding
    const maxFriendlyNameLen = Math.max(...modList.map((mod) => mod.friendlyName.length));
    const maxModReferenceLen = Math.max(...modList.map((mod) => mod.modReference.length));
    // Create header and add all mods to string
    let modListString = `${'Mod Name'.padEnd(maxFriendlyNameLen + 1) + 'Mod Reference'.padEnd(maxModReferenceLen + 1)}Version\n`;
    modList.forEach((mod) => {
      mod.friendlyName = mod.friendlyName.padEnd(maxFriendlyNameLen, ' ');
      mod.modReference = mod.modReference.padEnd(maxModReferenceLen, ' ');
      modListString += `${mod.friendlyName} ${mod.modReference} ${mod.version}\n`;
    });
    navigator.clipboard.writeText(modListString.trim());
  }
</script>

<div class="settings-menu">
  <div class="w-full h-8" use:popup={settingsMenu}>
    <button class="btn px-4 h-full w-full text-sm bg-surface-200-700-token"
    >
      <span>Mod Manager Settings</span>
      <div class="grow" />
      <SvgIcon
        class="h-5 w-5"
        icon={mdiTune} />
    </button>
  </div>
  <!-- #if gets executed before lower elements are added to the dom, so we have the submenus before to ensure they exist when use:popup is called-->
  <div class="card shadow-xl w-48 z-10 duration-0 overflow-y-auto" data-popup="update-check-mode-menu">
    <!-- 
    Skeleton's popup close function waits for the tranistion duration...
    before actually triggering the transition...
    So we'll just not have a transition...
    -->
    <ul class="menu">
      <ListBox class="w-full" rounded="rounded-none" spacing="space-y-0">
        {#each updateCheckModes as item}
          <ListBoxItem
            name="update-check-mode"
            class="bg-surface-50-900-token"
            active=""
            value={item.id}
            bind:group={$updateCheckMode}>
            {item.name}
            <span slot="trail" class="h-5 w-5 block">
              {#if $updateCheckMode === item.id}
                <SvgIcon class="h-full w-full" icon={mdiCheck}/>
              {/if}
            </span>
          </ListBoxItem>
        {/each}
      </ListBox>
    </ul>
  </div>
  <div class="card shadow-xl w-56 z-10 duration-0 overflow-y-auto" data-popup="queue-mode-menu">
    <!-- 
    Skeleton's popup close function waits for the tranistion duration...
    before actually triggering the transition...
    So we'll just not have a transition...
    -->
    <ul class="menu">
      <ListBox class="w-full" rounded="rounded-none" spacing="space-y-0">
        {#each queueModes as item}
          <ListBoxItem
            name="update-check-mode"
            class="bg-surface-50-900-token"
            active=""
            value={item.id}
            bind:group={$queueAutoStart}>
            {item.name}
            <span slot="trail" class="h-5 w-5 block">
              {#if $queueAutoStart === item.id}
                <SvgIcon class="h-full w-full" icon={mdiCheck}/>
              {/if}
            </span>
          </ListBoxItem>
        {/each}
      </ListBox>
    </ul>
  </div>
  <div class="card shadow-xl w-48 z-10 duration-0 overflow-y-auto" data-popup="start-view-menu">
    <!-- 
    Skeleton's popup close function waits for the tranistion duration...
    before actually triggering the transition...
    So we'll just not have a transition...
    -->
    <ul class="menu">
      <ListBox class="w-full" rounded="rounded-none" spacing="space-y-0">
        {#each views as item}
          <ListBoxItem
            name="update-check-mode"
            class="bg-surface-50-900-token"
            active=""
            value={item.id}
            bind:group={$startView}>
            {item.name}
            <span slot="trail" class="h-5 w-5 block">
              {#if $startView === item.id}
                <SvgIcon class="h-full w-full" icon={mdiCheck}/>
              {/if}
            </span>
          </ListBoxItem>
        {/each}
      </ListBox>
    </ul>
  </div>
  <div class="card shadow-xl w-48 z-10 duration-0 overflow-y-auto" data-popup="launch-button-menu">
    <!-- 
    Skeleton's popup close function waits for the tranistion duration...
    before actually triggering the transition...
    So we'll just not have a transition...
    -->
    <ul class="menu">
      <ListBox class="w-full" rounded="rounded-none" spacing="space-y-0">
        {#each launchButtons as item}
          <ListBoxItem
            name="update-check-mode"
            class="bg-surface-50-900-token"
            active=""
            value={item.id}
            bind:group={$launchButton}>
            {item.name}
            <span slot="trail" class="h-5 w-5 block">
              {#if $launchButton === item.id}
                <SvgIcon class="h-full w-full" icon={mdiCheck}/>
              {/if}
            </span>
          </ListBoxItem>
        {/each}
      </ListBox>
    </ul>
  </div>

  <!-- main settings menu starts here -->
  <div class="card shadow-xl z-10 duration-0 overflow-y-auto py-2" data-popup="settings-menu">
    <!-- 
    Skeleton's popup close function waits for the tranistion duration...
    before actually triggering the transition...
    So we'll just not have a transition...
    -->
    <ul class="menu">
      <li class="section-header">
        <span class="h-5 w-5"><SvgIcon class="h-full w-full" icon={mdiBug}/></span>
        <span class="flex-auto">Debug</span>
      </li>
      <hr class="divider" />
      <li>
        <button on:click={() => GenerateDebugInfo()}>
          <span class="h-5 w-5"/>
          <span class="flex-auto">Generate debug info</span>
          <span class="h-5 w-5"><SvgIcon class="h-full w-full" icon={mdiDownload}/></span>
        </button>
      </li>
      <hr class="divider" />
      <li>
        <button on:click={() => copyModList()}>
          <span class="h-5 w-5"/>
          <span class="flex-auto">Copy mod list</span>
          <span class="h-5 w-5"><SvgIcon class="h-full w-full" icon={mdiClipboard}/></span>
        </button>
      </li>
      <hr class="divider" />
      <li class="section-header">
        <span class="h-5 w-5"><SvgIcon class="h-full w-full" icon={mdiCog}/></span>
        <span class="flex-auto">Settings</span>
      </li>
      <hr class="divider" />
      <li data-noclose use:popup={updateCheckModeMenu}>
        <button>
          <span class="h-5 w-5"/>
          <span class="flex-auto">Update check</span>
          <span>{updateCheckModes.find((m) => m.id === $updateCheckMode)?.name}</span>
          <span class="h-5 w-5">
            <SvgIcon class="h-full w-full" icon={mdiChevronRight}/>
          </span>
        </button>
      </li>
      <hr class="divider" />
      <li data-noclose use:popup={queueModeMenu}>
        <button>
          <span class="h-5 w-5"/>
          <span class="flex-auto">Queue</span>
          <span>{queueModes.find((m) => m.id === $queueAutoStart)?.name}</span>
          <span class="h-5 w-5">
            <SvgIcon class="h-full w-full" icon={mdiChevronRight}/>
          </span>
        </button>
      </li>
      <hr class="divider" />
      <li data-noclose use:popup={startViewMenu}>
        <button>
          <span class="h-5 w-5"/>
          <span class="flex-auto">Start view</span>
          <span>{views.find((m) => m.id === $startView)?.name}</span>
          <span class="h-5 w-5">
            <SvgIcon class="h-full w-full" icon={mdiChevronRight}/>
          </span>
        </button>
      </li>
      <hr class="divider" />
      <li>
        <button on:click={() => modalStore.trigger({ type: 'component', component: 'cacheLocationPicker' })}>
          <span class="h-5 w-5"/>
          <span class="flex-auto">Change cache location</span>
          <span class="h-5 w-5"><SvgIcon class="h-full w-full" icon={mdiFolderEdit}/></span>
        </button>
      </li>
      <hr class="divider" />
      <li>
        <button on:click={() => { $offline = !$offline; }}>
          <span class="h-5 w-5"/>
          <span class="flex-auto">Go {$offline ? 'online' : 'offline'}</span>
          <span class="h-5 w-5"/>
        </button>
      </li>
      {#if $konami}
        <hr class="divider" />
        <li class="section-header">
          <span class="h-5 w-5"><SvgIcon class="h-full w-full" icon={mdiCog}/></span>
          <span class="flex-auto">Secret settings</span>
        </li>
        <hr class="divider" />
        <li data-noclose use:popup={launchButtonMenu}>
          <button>
            <span class="h-5 w-5"/>
            <span class="flex-auto">Launch button</span>
            <span>{launchButtons.find((l) => l.id === $launchButton)?.name ?? ''}</span>
            <span class="h-5 w-5">
              <SvgIcon class="h-full w-full" icon={mdiChevronRight}/>
            </span>
          </button>
        </li>
      {/if}
    </ul>
  </div>
</div>

<style lang="postcss">
  .menu {
    @apply list;
    > li {
      @apply h-10
    }
    > hr.divider {
      @apply border-surface-50 border-opacity-20 mr-4;
    }
    > .section-header {
      @apply pl-4 tracking-wider;
      + hr.divider {
        @apply ml-4;
      }
    }
    > :not(.section-header) + hr.divider {
      @apply ml-[3.25rem];
    }
    > :not(.section-header) > button {
      @apply btn px-4 w-full space-x-4 text-left;
      &:active {
        @apply bg-surface-400;
      }
    }
  }
</style>
