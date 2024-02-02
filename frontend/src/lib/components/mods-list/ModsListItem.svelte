<script lang="ts">
  import { mdiDownload, mdiEye, mdiStar, mdiPlay, mdiPause, mdiTrashCan, mdiTrayFull, mdiTrayMinus, mdiSync, mdiLinkLock, mdiArchiveCheck, mdiPauseCircle, mdiPlayCircle, mdiStarMinus, mdiStarPlus, mdiStarOutline, mdiTagMultiple } from '@mdi/js';
  import { createEventDispatcher } from 'svelte';
  import { getContextClient } from '@urql/svelte';
  import { popup, type PopupSettings , ProgressBar } from '@skeletonlabs/skeleton';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import Tooltip from '$lib/components/Tooltip.svelte';
  import { search, type PartialMod } from '$lib/store/modFiltersStore';
  import { favoriteMods, lockfileMods, manifestMods, progress, selectedInstallMetadata } from '$lib/store/ficsitCLIStore';
  import { addQueuedModAction, queuedMods, removeQueuedModAction } from '$lib/store/actionQueue';
  import { error, siteURL } from '$lib/store/generalStore';
  import { DisableMod, EnableMod, InstallMod, RemoveMod } from '$wailsjs/go/ficsitcli/ficsitCLI';
  import { FavoriteMod, UnFavoriteMod } from '$wailsjs/go/settings/settings';
  import { getAuthor } from '$lib/utils/getModAuthor';
  import { getCompatibility, getVersionCompatibility, type CompatibilityWithSource } from '$lib/utils/modCompatibility';
  import { CompatibilityState } from '$lib/generated';
  import Markdown from '$lib/components/Markdown.svelte';
  import type { ButtonDisplay } from '$lib/utils/responsiveButton';
  import ResponsiveButton from '$lib/components/ResponsiveButton.svelte';
  import { installTypeToTargetName } from '$lib/wailsTypesExtensions';

  export let mod: PartialMod;

  const client = getContextClient();

  const dispatch = createEventDispatcher();

  function listingClick() {
    dispatch('click');
  }

  $: authorClick = () => {
    $search = `author:"${author}"`;
  };

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
        iconHoverClass: 'text-red-500',
        tooltip: isInstalled ?
          'This mod is queued to be uninstalled. Click to cancel the operation.' :
            'This mod is queued to be installed. Click to cancel the operation.',
      };
    }

    let display: ButtonDisplay = {
      icon: mdiDownload,
      iconHover: mdiDownload,
      iconHoverClass: 'text-primary-600',
      tooltip: 'Click to install this mod.',
    };
    if (isInstalled) {
      display = {
        icon: mdiArchiveCheck,
        iconClass: 'text-primary-600',
        iconHover: mdiTrashCan,
        iconHoverClass: 'text-red-500',
        tooltip: 'This mod is installed on this profile. Click to uninstall it.',
      };
    } else if (compatibility.state !== CompatibilityState.Works) {
      if (installButtonDisabled) {
        display.tooltip = 'You can\'t install this mod. Reason:';
      } else {
        display.tooltip = 'There are problems reported with this mod, but you can try to install it anyways. Details:';
      }
      if (compatibility.note) {
        display.tooltipHtml = '<br/>' + compatibility.note;
      } else {
        // TODO compatibility.note should always be non-null here, what should our fallback text be if it's not?
        display.tooltip += ' (None specified)';
      }
    }
    if (queued) {
      display.tooltip = 'This mod is already queued for another operation.';
      delete display.tooltipHtml;
    }
    return display;
  })();

  $: enableButtonDisplay = ((): ButtonDisplay => {
    if (queuedEnable) {
      return {
        icon: mdiTrayFull,
        iconHover: mdiTrayMinus,
        iconHoverClass: 'text-red-500',
        tooltip: isEnabled ?
          'This mod is queued to be Disabled. Click to cancel the operation.' :
            'This mod is queued to be Enabled. Click to cancel the operation.',
      };
    }

    let display: ButtonDisplay = {
      icon: mdiPause,
      iconHover: mdiPlayCircle,
      iconHoverClass: 'text-primary-600',
      tooltip: 'This mod is Disabled on this profile. Click to Enable it.',
    };
    if (isEnabled) {
      display = {
        icon: mdiPlay,
        iconHover: mdiPauseCircle,
        tooltip: 'This mod is Enabled on this profile. Click to Disable it, which prevents it from loading when you start the game, but still keeps it a part of this profile.',
      };
    }
    if (queued) {
      display.tooltip = 'This mod is already queued for another operation.';
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
        iconClass: 'text-warning-500',
        tooltip: 'This mod is Favorited. Click to remove it from your favorites.',
      };
    }
    return {
      icon: mdiStarOutline,
      iconHover: mdiStarPlus,
      iconHoverClass: 'text-warning-500',
      tooltip: 'Click to add this mod to your Favorites.',
      tooltipHtml: 'Having a mod Favorited is unrelated to whether or not it\'s installed - it\'s a way to keep track of a mod for later regardless of what Profile you have selected.',
    };
  })();

  let compatibility: CompatibilityWithSource = { state: CompatibilityState.Works, source: 'reported' };
  $: {
    const gameVersion = $selectedInstallMetadata?.version;
    const branch = $selectedInstallMetadata?.branch;
    if(gameVersion && branch) {
      if(!('offline' in mod) && !('missing' in mod)) {
        if(mod.hidden && !isDependency) {
          compatibility = { state: CompatibilityState.Broken, note: 'This mod was hidden by the author.', source: 'reported' };
        } else {
          getCompatibility(mod.mod_reference, branch, gameVersion, installTypeToTargetName($selectedInstallMetadata.type), client).then((result) => {
            if (result.source === 'reported') {
              compatibility = {
                state: result.state,
                note: result.note 
                  ? `This mod has been reported as ${result.state} on this game version.<br>${result.note}` 
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

  $: popupId = `mods-list-item-popup-${mod.mod_reference}`;

  $: popupHover = {
    event: 'hover',
    target: popupId,
    middleware: {
      offset: 4,
    },
    placement: 'bottom-end',
  } satisfies PopupSettings;
</script>

<div
  class="my-1 px-0 @xl/mods-list:h-24 @md/mods-list:h-[5.5rem] h-[4.25rem]"
  class:bg-surface-50-900-token={selected}
  class:rounded-lg={selected}
  role="tab"
  tabindex="0"
  on:click={listingClick}
  on:keypress={listingClick}>
  {#if inProgress}
    <div class="relative h-full">
      <ProgressBar
        class="h-full w-full"
        max={1}
        meter="bg-surface-50-900-token"
        track="bg-surface-200-700-token"
        value={$progress?.progress === -1 ? undefined : $progress?.progress}/>
    </div>
  {/if}
  <div class="flex relative h-full" class:-top-full={inProgress}>
    <img
      class="logo h-full @xl/mods-list:w-24 @md/mods-list:w-[5.5rem] w-[4.25rem]"
      class:grayscale={isInstalled && !isEnabled}
      alt="{mod.name} Logo"
      src={renderedLogo} />
    <div class="ml-2 flex flex-col grow w-0 opacity" class:opacity-30={isInstalled && !isEnabled}>
      <div class="flex items-center" use:popup={popupHover}>
        <div class="shrink min-w-[7rem] truncate">
          <span class="@xl/mods-list:text-xl text-lg font-medium min-w-0 w-full" class:text-error-600={compatibility.state === CompatibilityState.Broken} class:text-warning-500={compatibility.state === CompatibilityState.Damaged}>{mod.name}</span>
        </div>
        <div class="shrink-0 hidden @lg/mods-list:block truncate w-[7rem] grow">
          <span class="pl-1">by</span>
          <!-- We could offer keyboard navigation for clicking this, but it's a waste of the user's time while nagivating via keyboard. If they want to search by author, they could enter the mod description pane -->
          <span
            class="text-primary-600 whitespace-nowrap"
            role="button"
            tabindex="-1"
            on:click|stopPropagation={authorClick}
            on:keypress|stopPropagation={authorClick}>{author}</span>
        </div>
      </div>
      <div class="truncate @md/mods-list:text-base text-sm hidden @md/mods-list:block">{'short_description' in mod ? mod.short_description : ''}</div>
      <div class="flex">
        {#if !inProgress}
          <div class="grow w-0 @xl/mods-list:text-base text-sm">
            <div class="truncate text-base @md/mods-list:text-sm block @md/mods-list:hidden">{'short_description' in mod ? mod.short_description : ''}</div>
            <div class="truncate h-5 @md/mods-list:h-4.5 hidden @md/mods-list:flex items-center space-x-1">
              {#if !('offline' in mod) && !('missing' in mod) && (mod?.tags?.length ?? -1 > 0 )}
                <SvgIcon class="pr-1 py-1 @xl/mods-list:w-7 w-6 shrink-0" icon={mdiTagMultiple}/>
                {#each mod?.tags ?? [] as tag}
                  <span>#{tag.name}</span>
                {/each}
              {/if}
              &nbsp; <!-- keep div height even when no tags are available -->
            </div>
            <div class="flex h-5 @md/mods-list:h-4.5 space-x-2">
              {#if !('offline' in mod) && !('missing' in mod)}
                <div class="w-24 flex items-center space-x-0.5">
                  <SvgIcon class="pr-1 py-1 @xl/mods-list:w-7 w-6" icon={mdiEye}/>
                  <span>{mod.views.toLocaleString()}</span>
                </div>
                <div class="w-24 flex items-center space-x-0.5">
                  <SvgIcon class="pr-1 py-1 @xl/mods-list:w-7 w-6" icon={mdiDownload}/>
                  <span>{mod.downloads.toLocaleString()}</span>
                </div>
              {/if}
            </div>
          </div>
        {:else}
          <span>{$progress?.message}</span>
        {/if}
      </div>
    </div>
    <!-- The purpose of the event handlers here are to prevent navigating to the mod's page when clicking on one of the sub-buttons of the div. Thus, it shouldn't be focusable despite having "interactions" -->
    <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
    <div
      class="pr-2 flex h-full items-center"
      role="separator"
      tabindex="-1"
      on:click|stopPropagation={() => { /* empty */ }}
      on:keypress|stopPropagation={() => { /* empty */ }}>
      <ResponsiveButton
        id="enable-{mod.mod_reference}"
        class="w-8 h-8 @lg/mods-list:mx-1 @xl/mods-list:mx-2"
        buttonClass="w-full h-full"
        disabled={enableButtonDisabled}
        display={enableButtonDisplay}
        onClickAction={toggleModEnabled}
        visible={isInstalled && !isDependency}
      />
      <ResponsiveButton
        id="install-{mod.mod_reference}"
        class="w-8 h-8 @lg/mods-list:mx-1 @xl/mods-list:mx-2"
        buttonClass="w-full h-full"
        disabled={installButtonDisabled}
        display={installButtonDisplay}
        onClickAction={toggleModInstalled}
      />
      <ResponsiveButton
        id="favorite-{mod.mod_reference}"
        class="w-8 h-8 @lg/mods-list:mx-1 @xl/mods-list:mx-2"
        buttonClass="w-full h-full"
        display={favoriteButtonDisplay}
        onClickAction={toggleModFavorite}
      />
    </div>
  </div>
</div>

<Tooltip disabled={!(isInstalled && !isEnabled) && compatibility.state === CompatibilityState.Works} {popupId}>
  {#if isInstalled && !isEnabled}
    This mod is Disabled. Click the pause icon to Enable it. 
  {:else if compatibility.state !== CompatibilityState.Works}
    <Markdown markdown={compatibility.note ?? ''}/>
  {/if}
</Tooltip>

