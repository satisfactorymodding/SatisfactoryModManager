<script lang="ts">
  import Drawer from '@smui/drawer';
  import { operationStore, query } from '@urql/svelte';
  import { GetModDetailsDocument } from '$lib/generated';
  import { markdown } from '$lib/utils/markdown';
  import Button, { Label } from '@smui/button';
  import MDIIcon from '$lib/components/MDIIcon.svelte';
  import { mdiChevronDown, mdiImport } from '@mdi/js';
  import Menu, { type MenuComponentDev } from '@smui/menu';
  import List, { Item, PrimaryText, SecondaryText, Text } from '@smui/list';
  import { bytesToAppropriate } from '$lib/utils/dataFormats';
  import { createEventDispatcher } from 'svelte';
  import { lockfileMods } from '$lib/ficsitCLIStore';
  import { search } from '$lib/modFiltersStore';

  export let id: string | null = null;

  const modQuery = operationStore(
    GetModDetailsDocument,
    { modReference: id }
  );
  
  $: modQuery.variables = {
    modReference: id
  };

  query(modQuery);

  $: mod = $modQuery.data?.mod;

  $: renderedLogo = mod?.logo ?? 'https://ficsit.app/images/no_image.webp';
  $: descriptionRendered = mod?.full_description ? markdown(mod.full_description) : undefined;

  $: size = mod ? bytesToAppropriate(mod.versions[0]?.size ?? 0) : undefined;

  $: installedVersion = $lockfileMods[mod?.mod_reference]?.version ?? 'Not installed';

  $: ficsitAppLink = `https://ficsit.app/mod/${id}`;

  let authorsMenu: MenuComponentDev;

  const dispatch = createEventDispatcher();

  function close() {
    dispatch('close');
  }
</script>

<div class="h-full flex mods-details">
  <Drawer class="w-64">
    <div class="px-4 py-4 flex flex-col w-full h-full mods-details">
      <img src={renderedLogo} alt="{mod?.name} Logo" class="logo w-full" />
      <span class="pt-4 font-bold text-lg">{mod?.name ?? 'Loading...'}</span>
      <span class="pt-2 font-light">A mod by:</span>
      <span class="font-medium color-primary cursor-pointer" on:click={() => $search = `author:"${mod?.authors[0].user.username}"`}>{mod?.authors[0].user.username ?? 'Loading...'}</span>

      <div class="pt-2" on:mouseenter={() => authorsMenu.setOpen(true)} on:mouseleave={() => authorsMenu.setOpen(false)}>
        <Button variant="unelevated" color="secondary" class="w-full">
          <Label>Contributors <span class="color-primary">({mod?.authors.length ?? 0})</span></Label>
          <MDIIcon icon={mdiChevronDown}/>
        </Button>
        <Menu bind:this={authorsMenu} class="w-full" anchorCorner="BOTTOM_LEFT">
          <List>
            {#each mod?.authors ?? [] as author}
              <Item style="height: 80px" on:SMUI:action={() => $search = `author:"${author.user.username}"`}>
                <img src={author.user.avatar} alt="{author.user.username} Avatar" class="avatar" />
                <Text class="pl-2 -mt-3">
                  <PrimaryText class="text-base">{author.user.username}</PrimaryText>
                  <SecondaryText class="text-base">{author.role}</SecondaryText>
                </Text>
              </Item>
            {/each}
          </List>
        </Menu>
      </div>

      <div class="pt-4">
        <span>Mod info:</span><br>
        <span>Size: </span><span class="font-bold">{size ?? 'Loading...'}</span><br>
        <span>Created: </span><span class="font-bold">{mod ? new Date(mod.created_at).toLocaleDateString() : 'Loading...'}</span><br>
        <span>Updated: </span><span class="font-bold">{mod ? new Date(mod.last_version_date).toLocaleString() : 'Loading...'}</span><br>
        <span>Total downloads: </span><span class="font-bold">{mod?.downloads.toLocaleString() ?? 'Loading...'}</span><br>
        <span>Views: </span><span class="font-bold">{mod?.views.toLocaleString() ?? 'Loading...'}</span><br>
      </div>

      <div class="pt-4">
        <span>Latest version: </span><span class="font-bold">{mod?.versions[0].version ?? 'Loading...'}</span><br>
        <span>Installed version: </span><span class="font-bold">{ installedVersion ?? 'Loading...' }</span><br>
      </div>

      <div class="pt-4">
        <a href={ficsitAppLink} target="_blank" class="color-primary">View on ficsit.app</a>
      </div>

      <div class="grow"></div>

      <Button variant="unelevated" color="secondary" on:click={close}>
        <div class="-scale-x-100">
          <MDIIcon icon={mdiImport}/>
        </div>
        <Label class="pl-4">Close</Label>
      </Button>
    </div>
  </Drawer>  
  <div class="markdown-content break-words flex-1 px-3 overflow-y-scroll">
    {#if descriptionRendered}
      <p>{@html descriptionRendered}</p>
    {:else}
      <p>Loading...</p>
    {/if}
  </div>
</div>

<style>
  .mods-details {
    background-color: #2B2B2B;
  }
  .avatar {
    border-radius: 50%;
    width: 50px;
    height: 50px;
  }
  a:visited {
    color: #249a20;
  }
  ::-webkit-scrollbar {
    width: 25px;
  }
  ::-webkit-scrollbar-track {
    background: black;
    border: solid 10px transparent;
    background-clip: content-box;
    border-radius: 0;
  }
  ::-webkit-scrollbar-thumb {
    background: #fff;
    border: solid 10px transparent;
    border-top-width: 0px;
    border-bottom-width: 0px;
    background-clip: content-box;
    border-radius: 0;
  }
  ::-webkit-scrollbar-thumb:hover {
    border: solid 10px transparent;
    border-top-width: 0px;
    border-bottom-width: 0px;
    background-clip: content-box;
  }
</style>
