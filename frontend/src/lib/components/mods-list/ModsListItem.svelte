<script lang="ts">
  import { mdiDownload, mdiEye, mdiChevronDown, mdiStar, mdiCheckCircle } from '@mdi/js';
  import MDIIcon from '$lib/components/MDIIcon.svelte';
  import { createEventDispatcher } from 'svelte';
  import type { PartialMod } from './modFilters';
  import Button, { Group, GroupItem, Label } from '@smui/button';
  import Menu, { type MenuComponentDev } from '@smui/menu';
  import List, { Item, Text } from '@smui/list';
  import LinearProgress from '@smui/linear-progress';
  import { favouriteMods, lockfileMods, manifestMods, progress } from '$lib/ficsitCLIStore';
  import { InstallMod, RemoveMod } from '$wailsjs/go/bindings/FicsitCLI';
  import { FavouriteMod, UnFavouriteMod } from '$wailsjs/go/bindings/Settings';
  
  export let mod: PartialMod;

  const dispatch = createEventDispatcher();

  function click() {
    dispatch('click');
  }

  let installOptionsMenu: MenuComponentDev;
  
  export let compact: boolean;
  export let selected: boolean;

  $: renderedLogo = mod.logo || 'https://ficsit.app/images/no_image.webp';
  $: author = mod.authors[0].user.username;

  $: isInstalled = mod.mod_reference in $manifestMods;
  $: isEnabled = mod.mod_reference in $lockfileMods;
  $: isDependency = !isInstalled && isEnabled;
  $: inProgress = $progress?.item === mod.mod_reference;

  $: installButtonLabel = isDependency ? 'Dependency' : (isInstalled ? 'Remove' : 'Install');
  $: installButtonIcon = isDependency ? mdiCheckCircle : (isInstalled ? mdiCheckCircle : mdiDownload);
  $: buttonDisabled = isDependency || (!!$progress);

  $: isFavourite = $favouriteMods.includes(mod.mod_reference);

  function toggleModInstalled() {
    if(isInstalled) {
      RemoveMod(mod.mod_reference);
    } else {
      InstallMod(mod.mod_reference);
    }
  }

  function toggleModFavourite() {
    if(!isFavourite) {
      FavouriteMod(mod.mod_reference);
    } else {
      UnFavouriteMod(mod.mod_reference);
    }
  }

  function toggleAddToQueue() {
    console.log('test');
  }
</script>

<div class="my-1 px-0 w-full" class:rounded-lg={selected} class:selected>
  {#if inProgress}
    <div class="relative h-0">
      <LinearProgress progress={$progress?.progress} class="mod-progress-bar h-24 w-full rounded-lg"/>
    </div>
  {/if}
  <div class="flex relative" readonly>
    <img src={renderedLogo} alt="{mod.name} Logo" class="logo h-24 w-24" on:click={click} />
    <div class="ml-2 flex flex-col grow w-0">
      <div on:click={click}>
        <span class="text-xl font-medium">{mod.name}</span>
        {#if !compact}
          <span class="pl-1">by</span>
          <span class="color-primary">{author}</span>
        {/if}
      </div>
      <span class="truncate w-full" on:click={click}>{mod.short_description}</span>
      <div class="flex grow">
        {#if !inProgress}
          <div class="grow w-0" on:click={click}>
            <div>
              #tags
            </div>
            <div class="flex h-5">
              <div class="w-24 flex">
                <div class="pr-1 inline-flex items-center justify-items-center"><MDIIcon icon={mdiEye}/></div>{mod.views.toLocaleString()}
              </div>
              <div class="w-24 flex">
                <div class="pr-1 inline-flex items-center justify-items-center"><MDIIcon icon={mdiDownload}/></div>{mod.downloads.toLocaleString()}
              </div>
            </div>
          </div>
          <div class="pr-2 flex">
            <Group variant="outlined" class="mr-1">
              <Button on:click={toggleModInstalled} variant="unelevated" disabled={buttonDisabled} class="{compact ? 'min-w-0 pl-5 pr-1.5' : 'w-32'} mod-install-button {isInstalled ? 'installed' : ''}">
                {#if !compact}
                  <Label>{ installButtonLabel }</Label>
                {:else}
                  <MDIIcon icon={ installButtonIcon }/>
                {/if}
              </Button>
              <div use:GroupItem>
                <Button
                  on:click={() => installOptionsMenu.setOpen(true)}
                  disabled={buttonDisabled}
                  variant="unelevated"
                  style="padding: 0; min-width: 36px;"
                  class="mod-install-button {isInstalled ? 'installed' : ''}"
                >
                  <MDIIcon icon={ mdiChevronDown }/>
                </Button>
                <Menu bind:this={installOptionsMenu} anchorCorner="TOP_LEFT">
                  <List>
                    <Item on:SMUI:action={toggleAddToQueue}>
                      <Text>Queue { installButtonLabel }</Text>
                    </Item>
                  </List>
                </Menu>
              </div>
            </Group>
            <Button on:click={toggleModFavourite} variant="unelevated" class="min-w-0 pl-5 pr-1.5 mod-favourite-button {isFavourite ? 'favourite' : ''}">
              <MDIIcon icon={ mdiStar }/>
            </Button>
          </div>
        {:else}
          <span>{ $progress?.message }</span>
        {/if}
      </div>
    </div>
  </div>
</div>

<style>
  :global(.mdc-linear-progress__bar-inner) {
    border-top-width: 100px;
  }
  .selected {
    background-color: #1c1c1c;
  }
</style>