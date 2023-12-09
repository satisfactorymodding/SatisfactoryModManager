<script lang="ts">
  import { mdiDownload, mdiEye, mdiStar, mdiCheckCircle, mdiPlay, mdiPause, mdiTrashCan, mdiTrayFull, mdiTrayMinus, mdiSync } from '@mdi/js';
  import { createEventDispatcher } from 'svelte';
  import Button from '@smui/button';
  import LinearProgress from '@smui/linear-progress';
  import Tooltip, { Wrapper } from '@smui/tooltip';
  import { getContextClient } from '@urql/svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { search, type PartialMod } from '$lib/store/modFiltersStore';
  import { favoriteMods, lockfileMods, manifestMods, progress, selectedInstall } from '$lib/store/ficsitCLIStore';
  import { addQueuedModAction, queuedMods, removeQueuedModAction } from '$lib/store/actionQueue';
  import { error } from '$lib/store/generalStore';
  import { DisableMod, EnableMod, InstallMod, RemoveMod } from '$wailsjs/go/ficsitcli/FicsitCLI';
  import { FavoriteMod, UnFavoriteMod } from '$wailsjs/go/bindings/Settings';
  import { getAuthor } from '$lib/utils/getModAuthor';
  import { getCompatibility, getVersionCompatibility, type CompatibilityWithSource } from '$lib/utils/modCompatibility';
  import type { GameBranch } from '$lib/wailsTypesExtensions';
  import { CompatibilityState } from '$lib/generated';
  import { markdown } from '$lib/utils/markdown';
  
  export let mod: PartialMod;

  const client = getContextClient();

  const dispatch = createEventDispatcher();

  function click() {
    dispatch('click');
  }
  
  export let selected: boolean;

  $: actualLogo = ('offline' in mod || 'missing' in mod) ? (mod.logo ? `data:image/png;base64, ${mod.logo}` : '/images/no_image.webp') : mod.logo;
  $: renderedLogo = actualLogo || 'https://ficsit.app/images/no_image.webp';
  $: author = ('offline' in mod || 'missing' in mod) ? mod.authors[0] : getAuthor(mod);

  $: isInstalled = mod.mod_reference in $manifestMods;
  $: isEnabled = $manifestMods[mod.mod_reference]?.enabled ?? mod.mod_reference in $lockfileMods;
  $: isDependency = !isInstalled && isEnabled;
  $: inProgress = $progress?.item === mod.mod_reference;
  $: queued = $queuedMods.some((q) => q.mod === mod.mod_reference);
  $: queuedInstall = $queuedMods.some((q) => q.mod === mod.mod_reference && (q.action === 'install' || q.action === 'remove'));
  $: queuedEnable = $queuedMods.some((q) => q.mod === mod.mod_reference && (q.action === 'enable' || q.action === 'disable'));
  
  $: installButtonIcon = (() => {
    if (isDependency) {
      return mdiCheckCircle;
    }
    if (inProgress) {
      return mdiSync;
    }
    if (queuedInstall) {
      return mdiTrayFull;
    }
    if (isInstalled) {
      return mdiCheckCircle;
    }
    return mdiDownload;
  })();
  
  $: installButtonIconHover = (() => {
    if (isDependency) {
      return mdiCheckCircle;
    }
    if (inProgress) {
      return mdiSync;
    }
    if (queuedInstall) {
      return mdiTrayMinus;
    }
    if (isInstalled) {
      return mdiTrashCan;
    }
    return mdiDownload;
  })();
  
  $: enableButtonIcon = (() => {
    if (queuedEnable) {
      return mdiTrayFull;
    }
    if (isEnabled) {
      return mdiPlay;
    }
    return mdiPause;
  })();
  
  $: enableButtonIconHover = (() => {
    if (queuedEnable) {
      return mdiTrayMinus;
    }
    if (isEnabled) {
      return mdiPause;
    }
    return mdiPlay;
  })();

  $: buttonDisabled = isDependency || (compatibility.state === CompatibilityState.Broken && compatibility.source === 'version' && !isInstalled);

  $: isFavorite = $favoriteMods.includes(mod.mod_reference);

  let compatibility: CompatibilityWithSource = { state: CompatibilityState.Works, source: 'reported' };
  $: {
    const gameVersion = $selectedInstall?.version;
    const branch = $selectedInstall?.branch as GameBranch;
    if(gameVersion && branch) {
      if(!('offline' in mod) && !('missing' in mod)) {
        if(mod.hidden && !isDependency) {
          compatibility = { state: CompatibilityState.Broken, note: 'This mod was hidden by the author.', source: 'reported' };
        } else {
          getCompatibility(mod.mod_reference, branch, gameVersion, client).then((result) => {
            if (result.source === 'reported') {
              compatibility = {
                state: result.state,
                note: result.note 
                  ? `This mod has been reported as ${result.state} on this game version.<br>${markdown(result.note)}` 
                  : `This mod has been reported as ${result.state} on this game version. (No further notes provided)`,
                source: 'reported',
              };
            } else {
              compatibility = result;
            }
          });
        }
      } else if('missing' in mod) {
        compatibility = { state: CompatibilityState.Broken, note: 'This mod is no longer available on ficsit.app. You may want to remove it.', source: 'version' };
      } else {
        getVersionCompatibility(mod.mod_reference, gameVersion, client).then((result) => {
          compatibility = {
            ...result,
            source: 'version',
          };
        });
      }
    }
  }

  async function toggleModInstalled() {
    // Svelte does not recreate the component, but reuse it, so the associated mod reference might change
    const modReference = mod.mod_reference;
    const action = isInstalled ? async () => RemoveMod(modReference).catch((e) => $error = e) : async () => InstallMod(modReference).catch((e) => $error = e);
    const actionName = isInstalled ? 'remove' : 'install';
    if(queued) {
      removeQueuedModAction(modReference);
      return;
    }
    return addQueuedModAction(
      modReference,
      actionName,
      action,
    );
  }

  async function toggleModEnabled() {
    // Svelte does not recreate the component, but reuse it, so the associated mod reference might change
    const modReference = mod.mod_reference;
    const action = isEnabled ? async () => DisableMod(modReference).catch((e) => $error = e) : async () => EnableMod(modReference).catch((e) => $error = e);
    const actionName = isEnabled ? 'disable' : 'enable';
    if(queued) {
      removeQueuedModAction(modReference);
      return;
    }
    return addQueuedModAction(
      modReference,
      actionName,
      action,
    );
  }

  async function toggleModFavorite() {
    try {
      if(!isFavorite) {
        await FavoriteMod(mod.mod_reference);
      } else {
        await UnFavoriteMod(mod.mod_reference);
      }
    } catch(e) {
      if (e instanceof Error) {
        $error = e.message;
      } else if (typeof e === 'string') {
        $error = e;
      } else {
        $error = 'Unknown error';
      }
    }
  }
</script>

<div class="my-1 px-0 lg:h-24 md:h-[5.5rem] h-[4.25rem]" class:rounded-lg={selected} class:selected on:click={click}>
  {#if inProgress}
    <div class="relative h-full">
      <LinearProgress progress={$progress?.progress} class="mod-progress-bar h-full rounded-lg"/>
    </div>
  {/if}
  <div class="flex relative h-full" class:-top-full={inProgress} class:disabled={isInstalled && !isEnabled}>
    <img src={renderedLogo} alt="{mod.name} Logo" class="logo h-full lg:w-24 md:w-[5.5rem] w-[4.25rem]" />
    <div class="ml-2 flex flex-col grow w-0">
      <Wrapper>
        <div class="flex items-center">
          <div class="shrink min-w-0 truncate">
            <span class="lg:text-xl text-lg font-medium min-w-0 w-full" class:error={compatibility.state === CompatibilityState.Broken} class:warning={compatibility.state === CompatibilityState.Damaged}>{mod.name}</span>
          </div>
          <div class="shrink-0">
            <span class="pl-1">by</span>
            <span class="color-primary whitespace-nowrap" on:click|stopPropagation={() => $search = `author:"${author}"`}>{author}</span>
          </div>
        </div>
        {#if isInstalled && !isEnabled}
          <Tooltip surface$class="max-w-lg text-base">
            This mod is disabled. Press the pause icon to enable it. 
          </Tooltip>
        {:else if compatibility.state !== CompatibilityState.Works}
          <Tooltip surface$class="max-w-lg text-base">
            <!-- eslint-disable-next-line svelte/no-at-html-tags -->
            { @html compatibility.note }
          </Tooltip>
        {/if}
      </Wrapper>
      <div class="truncate md:text-base text-sm hidden md:block">{'short_description' in mod ? mod.short_description : ''}</div>
      <div class="flex">
        {#if !inProgress}
          <div class="grow w-0 lg:text-base text-sm">
            <div class="truncate text-base md:text-sm block md:hidden">{'short_description' in mod ? mod.short_description : ''}</div>
            <div class="truncate hidden md:block">
              {#if !('offline' in mod) && !('missing' in mod)}
                {#each mod?.tags ?? [] as tag}
                  <span class="pr-1">#{tag.name}</span>
                {/each}
              {/if}
              &nbsp; <!-- keep div height even when no tags are available -->
            </div>
            <div class="flex h-5 md:h-4.5">
              {#if !('offline' in mod) && !('missing' in mod)}
                <div class="w-24 flex items-center">
                  <div class="pr-1 inline-flex items-center justify-items-center lg:w-7 w-6"><SvgIcon icon={mdiEye} class=""/></div><span>{mod.views.toLocaleString()}</span>
                </div>
                <div class="w-24 flex items-center">
                  <div class="pr-1 inline-flex items-center justify-items-center lg:w-7 w-6"><SvgIcon icon={mdiDownload} class=""/></div><span>{mod.downloads.toLocaleString()}</span>
                </div>
              {/if}
            </div>
          </div>
        {:else}
          <span>{ $progress?.message }</span>
        {/if}
      </div>
    </div>
    <div class="pr-2 flex h-full items-center" on:click|stopPropagation={() => { /* empty */ }}>
      {#if isInstalled && !isDependency}
        <Button
          on:click={ toggleModEnabled }
          disabled={buttonDisabled || queuedInstall || inProgress}
          ripple={false}
          variant="text"
          class="min-w-0 w-12 h-12 mod-enable-button {isEnabled || queued ? 'enabled' : ''} group"
        >
          <SvgIcon icon={enableButtonIcon} class="!p-1 !m-0 !w-full !h-full group-hover:!hidden"/>
          <SvgIcon icon={enableButtonIconHover} class="!p-1 !m-0 !w-full !h-full group-hover:!inline-block !hidden"/>
        </Button>
      {:else}
        <div class="min-w-0 w-12 h-12"/>
      {/if}
      <Button
        on:click={toggleModInstalled}
        variant="text"
        disabled={buttonDisabled || queuedEnable || inProgress}
        ripple={false}
        class="min-w-0 w-12 h-12 mod-install-button {isInstalled || queued ? 'installed' : ''} group"
      >
        <SvgIcon icon={installButtonIcon} class="!p-1 !m-0 !w-full !h-full group-hover:!hidden"/>
        <SvgIcon icon={installButtonIconHover} class="!p-1 !m-0 !w-full !h-full group-hover:!inline-block !hidden"/>
      </Button>
      <Button
        on:click={toggleModFavorite}
        variant="text"
        ripple={false}
        class="min-w-0 w-12 h-12 mod-favorite-button {isFavorite ? 'favorite' : ''}"
      >
        <SvgIcon icon={ mdiStar } class="!p-1 !m-0 !w-full !h-full"/>
      </Button>
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
  .disabled {
    opacity: 0.3;
  }
  .disabled img {
    filter: grayscale(1);
    animation-play-state: paused !important;
  }
</style>
