<script lang="ts">
  import { mdiDownload, mdiEye, mdiStar, mdiPlay, mdiPause, mdiTrashCan, mdiTrayFull, mdiTrayMinus, mdiSync, mdiLinkLock, mdiArchiveCheck, mdiPauseCircle, mdiPlayCircle, mdiStarMinus, mdiStarPlus, mdiStarOutline, mdiTagMultiple } from '@mdi/js';
  import { createEventDispatcher } from 'svelte';
  import Button from '@smui/button';
  import LinearProgress from '@smui/linear-progress';
  import Tooltip, { Wrapper } from '@smui/tooltip';
  import { getContextClient } from '@urql/svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { search, type PartialMod } from '$lib/store/modFiltersStore';
  import { favoriteMods, lockfileMods, manifestMods, progress, selectedInstall } from '$lib/store/ficsitCLIStore';
  import { addQueuedModAction, queuedMods, removeQueuedModAction } from '$lib/store/actionQueue';
  import { error, siteURL } from '$lib/store/generalStore';
  import { DisableMod, EnableMod, InstallMod, RemoveMod } from '$wailsjs/go/ficsitcli/FicsitCLI';
  import { FavoriteMod, UnFavoriteMod } from '$wailsjs/go/bindings/Settings';
  import { getAuthor } from '$lib/utils/getModAuthor';
  import { getCompatibility, getVersionCompatibility, type CompatibilityWithSource } from '$lib/utils/modCompatibility';
  import type { GameBranch } from '$lib/wailsTypesExtensions';
  import { CompatibilityState } from '$lib/generated';
  import { markdown } from '$lib/utils/markdown';
  import type { ButtonDisplay } from '$lib/utils/responsiveButton';
  import ResponsiveButton from '../ResponsiveButton.svelte';

  export let mod: PartialMod;

  const client = getContextClient();

  const dispatch = createEventDispatcher();

  function listingClick() {
    dispatch('click');
  }

  $: authorClick = () => {
    $search = `author:"${author}"`
  }

  export let selected: boolean;

  $: actualLogo = ('offline' in mod || 'missing' in mod) ? (mod.logo ? `data:image/png;base64, ${mod.logo}` : '/images/no_image.webp') : mod.logo;
  $: renderedLogo = actualLogo || `${$siteURL}/images/no_image.webp`;
  $: author = ('offline' in mod || 'missing' in mod) ? mod.authors[0] : getAuthor(mod);

  $: isInstalled = mod.mod_reference in $manifestMods;
  $: isEnabled = $manifestMods[mod.mod_reference]?.enabled ?? mod.mod_reference in $lockfileMods;
  $: isDependency = !isInstalled && isEnabled;
  $: inProgress = $progress?.item === mod.mod_reference;
  $: queued = $queuedMods.some((q) => q.mod === mod.mod_reference);
  $: queuedInstall = $queuedMods.some((q) => q.mod === mod.mod_reference && (q.action === 'install' || q.action === 'remove'));
  $: queuedEnable = $queuedMods.some((q) => q.mod === mod.mod_reference && (q.action === 'enable' || q.action === 'disable'));
  
  $: installButtonDisplay = ((): ButtonDisplay => {
    if (isDependency) {
      return {
        icon: mdiLinkLock,
        iconHover: mdiLinkLock,
        tooltip: 'This mod is installed a dependency of another mod. It cannot be installed or removed on its own.',
      };
    }
    if (inProgress) {
      return {
        icon: mdiSync,
        iconHover: mdiSync,
        tooltip: 'Wait for the current operation to complete.',
      };
    }
    if (queuedInstall) {
      return {
        icon: mdiTrayFull,
        iconHover: mdiTrayMinus,
        tooltip: isInstalled ?
          'This mod is queued to be uninstalled. Click to cancel the operation.' :
          'This mod is queued to be installed. Click to cancel the operation.',
      };
    }

    let display: ButtonDisplay = {
      icon: mdiDownload,
      iconHover: mdiDownload,
      tooltip: 'Click to install this mod.',
    }
    if (isInstalled) {
      display = {
        icon: mdiArchiveCheck,
        iconHover: mdiTrashCan,
        tooltip: 'This mod is installed on this profile. Click to uninstall it.',
      };
    } else if (compatibility.state !== CompatibilityState.Works) {
      if (installButtonDisabled) {
        display.tooltip = `You can't install this mod. Reason:`;
      } else {
        display.tooltip = `There are problems reported with this mod, but you can try to install it anyways. Details:`;
      }
      if (compatibility.note) {
        display.tooltipHtml = '<br/>' + compatibility.note;
      } else {
        // TODO compatibility.note should always be non-null here, what should our fallback text be if it's not?
        display.tooltip += ' (None specified)';
      }
    }
    if (queued) {
      display.tooltip = "This mod is already queued for another operation.";
      delete display.tooltipHtml;
    }
    return display;
  })();

  $: enableButtonDisplay = ((): ButtonDisplay => {
    if (queuedEnable) {
      return {
        icon: mdiTrayFull,
        iconHover: mdiTrayMinus,
        tooltip: isEnabled ?
          'This mod is queued to be Paused. Click to cancel the operation.' :
          'This mod is queued to be Resumed. Click to cancel the operation.',
      };
    }

    let display: ButtonDisplay = {
      icon: mdiPause,
      iconHover: mdiPlayCircle,
      tooltip: 'This mod is Paused on this profile. Click to Resume it.',
    };
    if (isEnabled) {
      display = {
        icon: mdiPlay,
        iconHover: mdiPauseCircle,
        tooltip: 'This mod is Enabled on this profile. Click to Pause it, which prevents it from loading when you start the game, but still keeps it a part of this profile.',
      };
    }
    if (queued) {
      display.tooltip = "This mod is already queued for another operation.";
    }
    return display;
  })();

  $: controlButtonsDisabled = isDependency || (compatibility.state === CompatibilityState.Broken && compatibility.source === 'version' && !isInstalled);
  $: installButtonDisabled = controlButtonsDisabled || queuedEnable || inProgress;
  $: enableButtonDisabled = controlButtonsDisabled || queuedInstall || inProgress;

  $: isFavorite = $favoriteMods.includes(mod.mod_reference);

  $: favoriteButtonDisplay = ((): ButtonDisplay => {
    if (isFavorite) {
      return {
        icon: mdiStar,
        iconHover: mdiStarMinus,
        tooltip: 'This mod is Favorited. Click to remove it from your favorites.',
      };
    }
    return {
      icon: mdiStarOutline,
      iconHover: mdiStarPlus,
      tooltip: 'Click to add this mod to your Favorites.',
      tooltipHtml: "Having a mod Favorited is unrelated to whether or not it's installed - it's a way to keep track of a mod for later regardless of what Profile you have selected."
    };
  })();

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

<div class="my-1 px-0 lg:h-24 md:h-[5.5rem] h-[4.25rem]" class:rounded-lg={selected} class:selected on:click={listingClick} on:keypress={listingClick} role="tab" tabindex="0">
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
            <!-- We could offer keyboard navigation for clicking this, but it's a waste of the user's time while nagivating via keyboard. If they want to search by author, they could enter the mod description pane -->
            <span class="color-primary whitespace-nowrap" on:click|stopPropagation={authorClick} on:keypress|stopPropagation={authorClick} role="button" tabindex="-1">{author}</span>
          </div>
        </div>
        {#if isInstalled && !isEnabled}
          <Tooltip surface$class="max-w-lg text-base">
            This mod is Paused. Press the pause icon to enable it. 
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
              {#if !('offline' in mod) && !('missing' in mod) && (mod?.tags?.length ?? -1 > 0 )}
                <div class="flex inline-flex items-center justify-items-center lg:w-4 w-3">
                  <!-- TODO this icon is not vertically centered -->
                  <SvgIcon icon={mdiTagMultiple} class=""/>
                </div>
                {#each mod?.tags ?? [] as tag}
                  <span class="pr-1">#{tag.name}</span>
                {/each}
              {/if}
              &nbsp; <!-- keep div height even when no tags are available -->
            </div>
            <div class="flex h-5 md:h-4.5">
              {#if !('offline' in mod) && !('missing' in mod)}
                <div class="w-24 flex items-center">
                  <div class="pr-1 inline-flex items-center justify-items-center lg:w-7 w-6">
                    <SvgIcon icon={mdiEye} class=""/>
                  </div>
                  <span>{mod.views.toLocaleString()}</span>
                </div>
                <div class="w-24 flex items-center">
                  <div class="pr-1 inline-flex items-center justify-items-center lg:w-7 w-6">
                    <SvgIcon icon={mdiDownload} class=""/>
                  </div>
                  <span>{mod.downloads.toLocaleString()}</span>
                </div>
              {/if}
            </div>
          </div>
        {:else}
          <span>{ $progress?.message }</span>
        {/if}
      </div>
    </div>
    <!-- The purpose of the event handlers here are to prevent navigating to the mod's page when clicking on one of the sub-buttons of the div. Thus, it shouldn't be focusable despite having "interactions" -->
    <!-- svelte-ignore a11y-no-noninteractive-element-interactions a11y-click-events-have-key-events -->
    <div class="pr-2 flex h-full items-center" role="separator" tabindex="-1" on:click|stopPropagation={() => { /* empty */ }} on:keypress|stopPropagation={() => { /* empty */ }}>
      <ResponsiveButton
        display={enableButtonDisplay}
        disabled={enableButtonDisabled}
        onClickAction={toggleModEnabled}
        class={'mod-enable-button ' + (isEnabled || queued ? 'enabled' : '')}
        visible={isInstalled && !isDependency}
      />
      <ResponsiveButton
        display={installButtonDisplay}
        disabled={installButtonDisabled}
        onClickAction={toggleModInstalled}
        class={'mod-install-button ' + (isInstalled || queued ? 'installed' : '')}
      />
      <ResponsiveButton
        display={favoriteButtonDisplay}
        onClickAction={toggleModFavorite}
        class={'mod-favorite-button ' + (isFavorite ? 'favorite' : '')}
      />
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
