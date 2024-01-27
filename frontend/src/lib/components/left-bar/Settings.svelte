<script lang="ts">
  import Menu, { SelectionGroup, SelectionGroupIcon } from '@smui/menu';
  import List, { Item, PrimaryText, Text, Separator } from '@smui/list';
  import { mdiBug, mdiCheck, mdiChevronRight, mdiClipboard, mdiCog, mdiDownload, mdiFolderEdit, mdiTune } from '@mdi/js';
  import { getContextClient } from '@urql/svelte';
  import Dialog, { Actions, Content, Title } from '@smui/dialog';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { GenerateDebugInfo } from '$wailsjs/go/bindings/DebugInfo';
  import { startView, konami, launchButton, queueAutoStart, offline, updateCheckMode, cacheDir } from '$lib/store/settingsStore';
  import { manifestMods, lockfileMods } from '$lib/store/ficsitCLIStore';
  import { GetModNameDocument } from '$lib/generated';
  import type { LaunchButtonType, ViewType } from '$lib/wailsTypesExtensions';
  import { OfflineGetMod } from '$wailsjs/go/ficsitcli/FicsitCLI';
  import { OpenDirectoryDialog } from '$lib/generated/wailsjs/go/bindings/App';

  let settingsMenu: Menu;
  let startViewMenu: Menu;

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

  let updateCheckModeMenu: Menu;

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

  let queueModeMenu: Menu;

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

  let launchButtonMenu: Menu;

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

  let cacheDialog = false;
  let cacheError: string | null = null;
  let newCacheLocation = $cacheDir;
  
  let fileDialogOpen = false;
  async function pickCacheLocation() {
    if(fileDialogOpen) {
      return;
    }
    fileDialogOpen = true;
    try {
      let result = await OpenDirectoryDialog({
        defaultDirectory: newCacheLocation ?? undefined,
      });
      if (result) {
        newCacheLocation = result;
      }
    } catch (e) {
      if(e instanceof Error) {
        cacheError = e.message;
      } else if (typeof e === 'string') {
        cacheError = e;
      } else {
        cacheError = 'Unknown error';
      }
    } finally {
      fileDialogOpen = false;
    }
  }

  let cacheMoveInProgress = false;

  async function setCacheLocation() {
    try {
      cacheMoveInProgress = true;
      await cacheDir.asyncSet(newCacheLocation ?? '');
      cacheError = null;
    } catch(e) {
      if (e instanceof Error) {
        cacheError = e.message;
      } else if (typeof e === 'string') {
        cacheError = e;
      } else {
        cacheError = 'Unknown error';
      }
    } finally {
      cacheMoveInProgress = false;
    }
  }

  async function resetCacheLocation() {
    try {
      cacheMoveInProgress = true;
      await cacheDir.asyncSet('');
      newCacheLocation = $cacheDir;
      cacheError = null;
    } catch(e) {
      if (e instanceof Error) {
        cacheError = e.message;
      } else if (typeof e === 'string') {
        cacheError = e;
      } else {
        cacheError = 'Unknown error';
      }
    } finally {
      cacheMoveInProgress = false;
    }
  }
</script>

<div class="settings-menu">
  <button
    class="btn px-4 h-8 w-full text-sm bg-surface-200-700-token"
    on:click={() => settingsMenu.setOpen(true)}>
    <span>Mod Manager Settings</span>
    <div class="grow" />
    <SvgIcon
      class="h-5 w-5"
      icon={mdiTune} />
  </button>
  <Menu bind:this={settingsMenu} class="w-full max-h-[32rem] overflow-visible" anchorCorner="TOP_RIGHT">
    <List on:SMUIList:action={(e) => { e.stopPropagation(); }}>
      <Item nonInteractive>
        <SvgIcon icon={mdiBug} class="h-5 w-5" />
        <!-- <div class="w-7"/> -->
        <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
          <PrimaryText class="text-base">Debug</PrimaryText>
        </Text>
        <div class="grow" />
      </Item>
      <Separator insetLeading insetTrailing />
      <Item on:SMUI:action={() => GenerateDebugInfo()}>
        <div class="w-7"/>
        <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
          <PrimaryText class="text-base">Generate debug info</PrimaryText>
        </Text>
        <div class="grow" />
        <SvgIcon icon={mdiDownload} class="h-5 w-5" />
      </Item>
      <Separator insetLeading insetTrailing insetPadding />
      <Item on:SMUI:action={() => copyModList()}>
        <div class="w-7"/>
        <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
          <PrimaryText class="text-base">Copy mods list</PrimaryText>
        </Text>
        <div class="grow" />
        <SvgIcon icon={mdiClipboard} class="h-5 w-5" />
      </Item>
      <Separator insetLeading insetTrailing insetPadding />
      <Item nonInteractive>
        <SvgIcon icon={mdiCog} class="h-5 w-5" />
        <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
          <PrimaryText class="text-base">Settings</PrimaryText>
        </Text>
        <div class="grow" />
      </Item>
      <Separator insetLeading insetTrailing />
      <div>
        <Item on:SMUI:action={() => updateCheckModeMenu.setOpen(true)}>
          <div class="w-7"/>
          <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
            <PrimaryText class="text-base">Update check</PrimaryText>
          </Text>
          <div class="grow" />
          <Text class="pr-2 h-full flex flex-col content-center mb-1.5">
            <PrimaryText class="text-base">{updateCheckModes.find((v) => v.id === $updateCheckMode)?.name ?? ''}</PrimaryText>
          </Text>
          <SvgIcon icon={mdiChevronRight} class="h-5 w-5" />
        </Item>
        <Menu bind:this={updateCheckModeMenu} class="w-full max-h-[32rem] overflow-visible" anchorCorner="TOP_RIGHT">
          <List>
            <SelectionGroup>
              {#each updateCheckModes as mode}
                <Item
                  on:SMUI:action={() => ($updateCheckMode = mode.id)}
                  selected={$updateCheckMode === mode.id}
                >
                  <SelectionGroupIcon>
                    <SvgIcon icon={mdiCheck} class="h-5 w-5" />
                  </SelectionGroupIcon>
                  <Text>{mode.name}</Text>
                </Item>
              {/each}
            </SelectionGroup>
          </List>
        </Menu>
      </div>
      <Separator insetLeading insetTrailing insetPadding />
      <div>
        <Item on:SMUI:action={() => queueModeMenu.setOpen(true)}>
          <div class="w-7"/>
          <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
            <PrimaryText class="text-base">Queue</PrimaryText>
          </Text>
          <div class="grow" />
          <Text class="pr-2 h-full flex flex-col content-center mb-1.5">
            <PrimaryText class="text-base">{queueModes.find((v) => v.id === $queueAutoStart)?.name ?? ''}</PrimaryText>
          </Text>
          <SvgIcon icon={mdiChevronRight} class="h-5 w-5" />
        </Item>
        <Menu bind:this={queueModeMenu} class="w-full max-h-[32rem] overflow-visible" anchorCorner="TOP_RIGHT">
          <List>
            <SelectionGroup>
              {#each queueModes as queueMode}
                <Item
                  on:SMUI:action={() => ($queueAutoStart = queueMode.id)}
                  selected={$queueAutoStart === queueMode.id}
                >
                  <SelectionGroupIcon>
                    <SvgIcon icon={mdiCheck} class="h-5 w-5" />
                  </SelectionGroupIcon>
                  <Text>{queueMode.name}</Text>
                </Item>
              {/each}
            </SelectionGroup>
          </List>
        </Menu>
      </div>
      <Separator insetLeading insetTrailing insetPadding />
      <div>
        <Item on:SMUI:action={() => startViewMenu.setOpen(true)}>
          <div class="w-7"/>
          <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
            <PrimaryText class="text-base">Start view</PrimaryText>
          </Text>
          <div class="grow" />
          <Text class="pr-2 h-full flex flex-col content-center mb-1.5">
            <PrimaryText class="text-base">{views.find((v) => v.id === $startView)?.name ?? ''}</PrimaryText>
          </Text>
          <SvgIcon icon={mdiChevronRight} class="h-5 w-5" />
        </Item>
        <Menu bind:this={startViewMenu} class="w-full max-h-[32rem] overflow-visible" anchorCorner="TOP_RIGHT">
          <List>
            <SelectionGroup>
              {#each views as view}
                <Item
                  on:SMUI:action={() => ($startView = view.id)}
                  selected={$startView === view.id}
                >
                  <SelectionGroupIcon>
                    <SvgIcon icon={mdiCheck} class="h-5 w-5" />
                  </SelectionGroupIcon>
                  <Text>{view.name}</Text>
                </Item>
              {/each}
            </SelectionGroup>
          </List>
        </Menu>
      </div>
      <Separator insetLeading insetTrailing insetPadding />
      <Item on:SMUI:action={() => { cacheDialog = true; settingsMenu.setOpen(false); } }>
        <div class="w-7"/>
        <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
          <PrimaryText class="text-base">Change cache location</PrimaryText>
        </Text>
        <div class="grow" />
        <SvgIcon icon={mdiFolderEdit} class="h-5 w-5" />
      </Item>
      <Separator insetLeading insetTrailing insetPadding />
      <Item on:SMUI:action={() => { $offline = !$offline; settingsMenu.setOpen(false); }}>
        <div class="w-7"/>
        <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
          <PrimaryText class="text-base">Go { $offline ? 'Online' : 'Offline' }</PrimaryText>
        </Text>
      </Item>
      {#if $konami}
        <Separator insetLeading insetTrailing insetPadding />
        <Item nonInteractive>
          <SvgIcon icon={mdiCog} class="h-5 w-5" />
          <!-- <div class="w-7"/> -->
          <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
            <PrimaryText class="text-base">Secret settings</PrimaryText>
          </Text>
          <div class="grow" />
        </Item>
        <Separator insetLeading insetTrailing />
        <div>
          <Item on:SMUI:action={() => launchButtonMenu.setOpen(true)}>
            <div class="w-7"/>
            <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
              <PrimaryText class="text-base">Launch button</PrimaryText>
            </Text>
            <div class="grow" />
            <Text class="pr-2 h-full flex flex-col content-center mb-1.5">
              <PrimaryText class="text-base">{launchButtons.find((l) => l.id === $launchButton)?.name ?? ''}</PrimaryText>
            </Text>
            <SvgIcon icon={mdiChevronRight} class="h-5 w-5" />
          </Item>
          <Menu bind:this={launchButtonMenu} class="w-full max-h-[32rem] overflow-visible" anchorCorner="TOP_RIGHT">
            <List>
              <SelectionGroup>
                {#each launchButtons as launch}
                  <Item
                    on:SMUI:action={() => ($launchButton = launch.id)}
                    selected={$launchButton === launch.id}
                  >
                    <SelectionGroupIcon>
                      <SvgIcon icon={mdiCheck} class="h-5 w-5" />
                    </SelectionGroupIcon>
                    <Text>{launch.name}</Text>
                  </Item>
                {/each}
              </SelectionGroup>
            </List>
          </Menu>
        </div>
      {/if}
    </List>
  </Menu>
</div>


<Dialog
  bind:open={cacheDialog}
  scrimClickAction="" escapeKeyAction=""
  surface$style="max-height: calc(100vh - 128px); max-width: calc(100vw - 128px);"
  surface$class="!min-w-[800px] min-h-[400px]"
>
  <Title>Change download cache location</Title>
  <Content>
    <label class="label">
      <span>Cache location</span>
      <div class="flex items-baseline">
        <div class="grow">
          <input type="text"
            class="input px-4 py-2 hover:!cursor-pointer"
            class:input-error={cacheError}
            value={newCacheLocation}
            readonly
            on:click={() => pickCacheLocation()}
          />
          <p>
            {#if cacheError }
              { cacheError }
            {/if}
          </p>
        </div>
        <button
          class="btn mr-4 shrink-0 text-primary-600"
          disabled={cacheMoveInProgress}
          on:click={() => resetCacheLocation()}>
          <span>Reset to default</span>
          <div class="grow" />
        </button>
        <button
          class="btn shrink-0 text-primary-600"
          disabled={cacheMoveInProgress}
          on:click={() => setCacheLocation()}>
          <span>Save and move</span>
          <div class="grow" />
        </button>
      </div>
    </label>
  </Content>
  <Actions>
    <button
      class="btn text-primary-600"
      disabled={cacheMoveInProgress}
      on:click={() => { cacheDialog = false; }}>
      <span>Close</span>
      <div class="grow" />
    </button>
  </Actions>
</Dialog>
