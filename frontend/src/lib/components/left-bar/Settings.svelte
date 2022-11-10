<script lang="ts">
  import Button, { Label } from '@smui/button';
  import Menu, { SelectionGroup, SelectionGroupIcon } from '@smui/menu';
  import List, { Item, PrimaryText, Text, Separator } from '@smui/list';
  
  import { mdiBug, mdiCheck, mdiChevronRight, mdiClipboard, mdiCog, mdiDownload, mdiTune } from '@mdi/js';
  
  import SvgIcon from '$lib/components/SVGIcon.svelte';
  
  import { GenerateDebugInfo } from '$wailsjs/go/bindings/DebugInfo';
  
  import { startView, konami, launchButton, queueAutoStart } from '$lib/store/settingsStore';
  import { manifestMods, lockfileMods } from '$lib/store/ficsitCLIStore';
  import { GetModNameDocument } from '$lib/generated';
  
  import { getContextClient } from '@urql/svelte';
  import type { LaunchButtonType, ViewType } from '$lib/wailsTypesExtensions';

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
    }
  ];
  
  const urqlClient = getContextClient();

  async function copyModList() {
    // Generate mod entries
    const modList = await Promise.all(Object.keys($manifestMods).map(async (modReference) => {
      let modName = modReference;
      if(modReference === 'SML') {
        modName = 'Satisfactory Mod Loader';
      } else {
        const result = await urqlClient.query(GetModNameDocument, { modReference }).toPromise();
        if(result?.data?.getModByReference?.name) {
          modName = result.data.getModByReference.name;
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
  <Button variant="unelevated" class="w-full mt-2" on:click={() => settingsMenu.setOpen(true)}>
    <Label>
      SMM settings
    </Label>
    <div class="grow" />
    <SvgIcon icon={mdiTune} class="h-5 w-5" />
  </Button>
  <Menu bind:this={settingsMenu} class="w-full max-h-[32rem] overflow-visible" anchorCorner="TOP_RIGHT">
    <List>
      <Item nonInteractive>
        <SvgIcon icon={mdiBug} class="h-5 w-5" />
        <!-- <div class="w-7"/> -->
        <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
          <PrimaryText class="text-base">Debug</PrimaryText>
        </Text>
        <div class="grow" />
      </Item>
      <Separator insetLeading insetTrailing />
      <Item on:click={() => GenerateDebugInfo()}>
        <div class="w-7"/>
        <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
          <PrimaryText class="text-base">Generate debug info</PrimaryText>
        </Text>
        <div class="grow" />
        <SvgIcon icon={mdiDownload} class="h-5 w-5" />
      </Item>
      <Separator insetLeading insetTrailing insetPadding />
      <Item on:click={() => copyModList()}>
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
        <Item on:click={() => queueModeMenu.setOpen(true)}>
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
      <Separator insetLeading insetTrailing />
      <div>
        <Item on:click={() => startViewMenu.setOpen(true)}>
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
        <Item on:click={() => launchButtonMenu.setOpen(true)}>
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