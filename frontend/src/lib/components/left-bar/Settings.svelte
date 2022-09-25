<script lang="ts">
  import Button, { Label } from '@smui/button';
  import Menu, { type MenuComponentDev, SelectionGroup, SelectionGroupIcon } from '@smui/menu';
  import List, { Item, PrimaryText, Text, Separator } from '@smui/list';
  
  import { mdiBug, mdiCheck, mdiChevronRight, mdiClipboard, mdiCog, mdiDownload, mdiTune } from '@mdi/js';
  
  import MdiIcon from '$lib/components/MDIIcon.svelte';
  
  import { GenerateDebugInfo } from '$wailsjs/go/bindings/DebugInfo';
  
  import { startView, type View } from '$lib/store/settingsStore';
  import { manifestMods, lockfileMods } from '$lib/store/ficsitCLIStore';
  import { GetModNameDocument } from '$lib/generated';
  
  import { getClient } from '@urql/svelte';

  let settingsMenu: MenuComponentDev;
  let startViewMenu: MenuComponentDev;

  let views: {id: View, name: string}[] = [
    {
      id: 'compact',
      name: 'Compact',
    },
    {
      id: 'expanded',
      name: 'Expanded',
    },
  ];
  
  const urqlClient = getClient();

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
    <MdiIcon icon={mdiTune} class="h-5" />
  </Button>
  <Menu bind:this={settingsMenu} class="w-full max-h-[32rem] overflow-visible" anchorCorner="TOP_RIGHT">
    <List>
      <Item nonInteractive>
        <MdiIcon icon={mdiBug} class="h-5" />
        <!-- <div class="w-7"/> -->
        <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
          <PrimaryText class="text-base">Debug</PrimaryText>
        </Text>
        <div class="grow" />
      </Item>
      <Separator insetLeading insetTrailing />
      <Item on:click={() => GenerateDebugInfo()}>
        <MdiIcon icon={mdiDownload} class="h-5" />
        <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
          <PrimaryText class="text-base">Generate debug info</PrimaryText>
        </Text>
        <div class="grow" />
      </Item>
      <Separator insetLeading insetTrailing insetPadding />
      <Item on:click={() => copyModList()}>
        <MdiIcon icon={mdiClipboard} class="h-5" />
        <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
          <PrimaryText class="text-base">Copy mods list</PrimaryText>
        </Text>
        <div class="grow" />
      </Item>
      <Separator insetLeading insetTrailing insetPadding />
      <Item nonInteractive>
        <MdiIcon icon={mdiCog} class="h-5" />
        <!-- <div class="w-7"/> -->
        <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
          <PrimaryText class="text-base">Settings</PrimaryText>
        </Text>
        <div class="grow" />
      </Item>
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
          <MdiIcon icon={mdiChevronRight} class="h-5" />
        </Item>
        <Menu bind:this={startViewMenu} class="w-full max-h-[32rem] overflow-visible" anchorCorner="TOP_RIGHT">
          <SelectionGroup>
            {#each views as view}
              <Item
                on:SMUI:action={() => ($startView = view.id)}
                selected={$startView === view.id}
              >
                <SelectionGroupIcon>
                  <MdiIcon icon={mdiCheck} class="h-5" />
                </SelectionGroupIcon>
                <Text>{view.name}</Text>
              </Item>
            {/each}
          </SelectionGroup>
        </Menu>
      </div>
    </List>
  </Menu>
</div>