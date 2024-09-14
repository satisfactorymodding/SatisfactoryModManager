<script lang="ts">
  import { mdiArchiveCheck, mdiDownload, mdiEye, mdiLinkLock, mdiPause, mdiPauseCircle, mdiPlay, mdiPlayCircle, mdiStar, mdiStarMinus, mdiStarOutline, mdiStarPlus, mdiSync, mdiTagMultiple, mdiTrashCan, mdiTrayFull, mdiTrayMinus } from '@mdi/js';
  import { ProgressBar } from '@skeletonlabs/skeleton';
  import { getTranslate } from '@tolgee/svelte';
  import { getContextClient } from '@urql/svelte';
  import { createEventDispatcher } from 'svelte';

  import Markdown from '$lib/components/Markdown.svelte';
  import ResponsiveButton from '$lib/components/ResponsiveButton.svelte';
  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import T from '$lib/components/T.svelte';
  import Tooltip from '$lib/components/Tooltip.svelte';
  import { CompatibilityState } from '$lib/generated';
  import { type PopupSettings, popup } from '$lib/skeletonExtensions';
  import { addQueuedModAction, queuedMods, removeQueuedModAction } from '$lib/store/actionQueue';
  import { canInstallMods, favoriteMods, lockfileMods, manifestMods, progress, progressMessage, progressPercent, selectedInstallMetadata } from '$lib/store/ficsitCLIStore';
  import { error, siteURL } from '$lib/store/generalStore';
  import { type PartialMod, search } from '$lib/store/modFiltersStore';
  import { largeNumberFormat } from '$lib/utils/dataFormats';
  import { getAuthor } from '$lib/utils/getModAuthor';
  import { type CompatibilityWithSource, getCompatibility, getVersionCompatibility } from '$lib/utils/modCompatibility';
  import type { ButtonDisplay } from '$lib/utils/responsiveButton';
  import { installTypeToTargetName } from '$lib/wailsTypesExtensions';
  import { DisableMod, EnableMod, InstallMod, RemoveMod } from '$wailsjs/go/ficsitcli/ficsitCLI';
  import { ficsitcli } from '$wailsjs/go/models';
  import { FavoriteMod, UnFavoriteMod } from '$wailsjs/go/settings/settings';

  export let mod: PartialMod;

  const { t } = getTranslate();

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

  const modActions = [
    ficsitcli.Action.INSTALL,
    ficsitcli.Action.UNINSTALL,
    ficsitcli.Action.ENABLE,
    ficsitcli.Action.DISABLE,
  ];

  $: isInstalled = mod.mod_reference in $manifestMods;
  $: isEnabled = $manifestMods[mod.mod_reference]?.enabled ?? mod.mod_reference in $lockfileMods;
  $: isDependency = !isInstalled && isEnabled;
  $: inProgress = !!$progress && modActions.includes($progress.action) && $progress.item.name === mod.mod_reference;
  $: queued = $queuedMods.some((q) => q.mod === mod.mod_reference);
  $: queuedInstall = $queuedMods.some((q) => q.mod === mod.mod_reference && (q.action === 'install' || q.action === 'remove'));
  $: queuedEnable = $queuedMods.some((q) => q.mod === mod.mod_reference && (q.action === 'enable' || q.action === 'disable'));
  
  $: installButtonDisplay = ((): ButtonDisplay => {
    if (inProgress) {
      return {
        icon: mdiSync,
        iconHover: mdiSync,
        tooltip: $t('mod-list-item.wait', 'Wait for the current operation to complete.'),
      };
    }
    if (isDependency) {
      return {
        icon: mdiLinkLock,
        iconHover: mdiLinkLock,
        tooltip: $t('mod-list-item.dependency', 'This mod is installed as a dependency of another mod. It cannot be installed or removed on its own.'),
      };
    }
    if (queuedInstall) {
      return {
        icon: mdiTrayFull,
        iconHover: mdiTrayMinus,
        iconHoverClass: 'text-red-500',
        tooltip: isInstalled ?
          $t('mod-list-item.queued-uninstall', 'This mod is queued to be uninstalled. Click to cancel the operation.') :
            $t('mod-list-item.queued-install', 'This mod is queued to be installed. Click to cancel the operation.'),
      };
    }

    let display: ButtonDisplay = {
      icon: mdiDownload,
      iconHover: mdiDownload,
      iconHoverClass: 'text-primary-600',
      tooltip: $t('mod-list-item.install', 'Click to install this mod.'),
    };
    if (isInstalled) {
      display = {
        icon: mdiArchiveCheck,
        iconClass: 'text-primary-600',
        iconHover: mdiTrashCan,
        iconHoverClass: 'text-red-500',
        tooltip: $t('mod-list-item.uninstall', 'This mod is installed on this profile. Click to uninstall this mod.'),
      };
    } else if (compatibility.state !== CompatibilityState.Works) {
      if (compatibility.state === CompatibilityState.Broken && compatibility.source === 'version') {
        display.tooltip = $t('mod-list-item.not-installable', 'You can\'t install this mod. Reason:');
      } else {
        display.tooltip = $t('mod-list-item.compatibility-warning', 'There are problems reported with this mod, but you can try to install it anyways. Details:');
      }
      if (compatibility.note) {
        display.tooltipHtml = '<br/>' + compatibility.note;
      } else {
        // TODO compatibility.note should always be non-null here, what should our fallback text be if it's not?
        display.tooltip += $t('mod-list-item.compatibility-note-none', ' (None specified)');
      }
    }
    if (queued) {
      display.tooltip = $t('mod-list-item.queued', 'This mod is already queued for another operation.');
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
          $t('mod-list-item.disable-queued', 'This mod is queued to be Disabled. Click to cancel the operation.') :
            $t('mod-list-item.enable-queued', 'This mod is queued to be Enabled. Click to cancel the operation.'),
      };
    }

    let display: ButtonDisplay = {
      icon: mdiPause,
      iconHover: mdiPlayCircle,
      iconHoverClass: 'text-primary-600',
      tooltip: $t('mod-list-item.enable', 'Click to enable this mod.'),
    };
    if (isEnabled) {
      display = {
        icon: mdiPlay,
        iconHover: mdiPauseCircle,
        tooltip: $t('mod-list-item.disable', 'This mod is Enabled on this profile. Click to Disable it, which prevents it from loading when you start the game, but still keeps it a part of this profile.'),
      };
    }
    if (queued) {
      display.tooltip = $t('mod-list-item.queued', 'This mod is already queued for another operation.');
    }
    return display;
  })();

  $: controlButtonsDisabled = !$canInstallMods || isDependency || (compatibility.state === CompatibilityState.Broken && compatibility.source === 'version' && !isInstalled);
  $: installButtonDisabled = controlButtonsDisabled || queuedEnable || inProgress;
  $: enableButtonDisabled = controlButtonsDisabled || queuedInstall || inProgress;

  $: isFavorite = $favoriteMods.includes(mod.mod_reference);

  $: favoriteButtonDisplay = ((): ButtonDisplay => {
    if (isFavorite) {
      return {
        icon: mdiStar,
        iconHover: mdiStarMinus,
        iconClass: 'text-warning-500',
        tooltip: $t('mod-list-item.unfavorite', 'Click to remove this mod from your Favorites.'),
      };
    }
    return {
      icon: mdiStarOutline,
      iconHover: mdiStarPlus,
      iconHoverClass: 'text-warning-500',
      tooltip: $t('mod-list-item.favorite', 'Click to add this mod to your Favorites. Having a mod Favorited is unrelated to whether or not it\'s installed - it\'s a way to keep track of a mod for later regardless of what Profile you have selected.'),
    };
  })();

  let compatibility: CompatibilityWithSource = { state: CompatibilityState.Works, source: 'reported' };
  $: {
    const info = $selectedInstallMetadata?.info;
    if(info) {
      if(!('offline' in mod) && !('missing' in mod)) {
        if(mod.hidden && !isDependency) {
          compatibility = { state: CompatibilityState.Broken, note: $t('mod-list-item.hidden', 'This mod was hidden by the author.'), source: 'reported' };
        } else {
          getCompatibility(mod.mod_reference, info.branch, info.version, installTypeToTargetName(info.type), client).then((result) => {
            if (result.source === 'reported') {
              compatibility = {
                state: result.state,
                note: $t('mod-list-item.compatibility-note', 'This mod has been reported as {state} on this game version.', { state: result.state }) 
                  + (result.note ? ('<br>' + result.note) : (' ' + $t('mod.compatibility-no-notes', '(No further notes provided)'))),
                source: 'reported',
              };
            } else {
              compatibility = result;
            }
          });
        }
      } else if('missing' in mod) {
        compatibility = { state: CompatibilityState.Broken, note: $t('mod-list-item.unavailable', 'This mod is no longer available on ficsit.app. You may want to remove it.'), source: 'version' };
      } else {
        getVersionCompatibility(mod.mod_reference, info.version, client).then((result) => {
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
    placement: 'bottom-start',
  } satisfies PopupSettings;
</script>

<div
  class="my-1 px-0 @md/mods-list:h-[5.5rem] h-[4.25rem] overflow-hidden"
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
        value={$progressPercent}/>
    </div>
  {/if}
  <div class="flex relative h-full space-x-2" class:-top-full={inProgress}>
    <img
      class="logo h-full @md/mods-list:w-[5.5rem] w-[4.25rem]"
      class:grayscale={isInstalled && !isEnabled}
      alt="{mod.name} Logo"
      src={renderedLogo} />
    <div class="flex flex-col h-full grow w-0" class:opacity-30={isInstalled && !isEnabled}>
      <div class="flex items-center" use:popup={popupHover}>
        <span 
          class="shrink min-w-[7rem] truncate text-lg font-medium !leading-6" 
          class:text-error-600={compatibility.state === CompatibilityState.Broken}
          class:text-warning-500={compatibility.state === CompatibilityState.Damaged}>
          {mod.name}
        </span>
        <div class="shrink-0 hidden @lg/mods-list:block truncate w-[7rem] grow">
          <span class="pl-1"><T defaultValue="by" keyName="mod-list-item.by-author" /></span>
          <!-- We could offer keyboard navigation for clicking this, but it's a waste of the user's time while nagivating via keyboard. If they want to search by author, they could enter the mod description pane -->
          <span
            class="text-primary-600 whitespace-nowrap"
            role="button"
            tabindex="-1"
            on:click|stopPropagation={authorClick}
            on:keypress|stopPropagation={authorClick}>{author}</span>
        </div>
      </div>
      <div class="flex flex-col grow h-0 overflow-hidden justify-around flex-wrap">
        <div class="truncate w-full @md/mods-list:text-base text-sm">{'short_description' in mod ? mod.short_description : ''}</div>
        <div class="truncate w-full h-5 hidden @md/mods-list:flex items-center space-x-1">
          <SvgIcon class="w-5 shrink-0" icon={mdiTagMultiple}/>
          {#if !('offline' in mod) && !('missing' in mod) && (mod?.tags?.length ?? -1 > 0 )}
            {#each mod?.tags ?? [] as tag}
              <span>#{tag.name}</span>
            {/each}
          {:else}
            <span><T defaultValue="(none available)" keyName="mod-list-item.no-tags" /></span>
          {/if}
        </div>
        <div class="text-sm w-full">
          <div class="flex h-5 space-x-2">
            {#if !('offline' in mod) && !('missing' in mod)}
              <div class="w-16 flex items-center space-x-1">
                <SvgIcon class="w-4 @md/mods-list:w-5" icon={mdiEye}/>
                <span>{largeNumberFormat(mod.views)}</span>
              </div>
              <div class="w-16 flex items-center space-x-1">
                <SvgIcon class="w-4 @md/mods-list:w-5" icon={mdiDownload}/>
                <span>{largeNumberFormat(mod.downloads)}</span>
              </div>
            {/if}
          </div>
        </div>
      </div>
      {#if inProgress}
        <span class="shrink-0 text-sm">{$progressMessage}</span>
      {/if}
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
        class="w-8 h-8 @lg/mods-list:mx-1"
        buttonClass="w-full h-full"
        disabled={enableButtonDisabled}
        display={enableButtonDisplay}
        onClickAction={toggleModEnabled}
        visible={isInstalled && !isDependency}
      />
      <ResponsiveButton
        id="install-{mod.mod_reference}"
        class="w-8 h-8 @lg/mods-list:mx-1"
        buttonClass="w-full h-full"
        disabled={installButtonDisabled}
        display={installButtonDisplay}
        onClickAction={toggleModInstalled}
      />
      <ResponsiveButton
        id="favorite-{mod.mod_reference}"
        class="w-8 h-8 @lg/mods-list:mx-1"
        buttonClass="w-full h-full"
        display={favoriteButtonDisplay}
        onClickAction={toggleModFavorite}
      />
    </div>
  </div>
</div>

<Tooltip disabled={!(isInstalled && !isEnabled) && compatibility.state === CompatibilityState.Works} {popupId}>
  {#if isInstalled && !isEnabled}
    <T defaultValue="This mod is Disabled. Click the pause icon to Enable it." keyName="mod-list-item.disabled-tooltip" />
  {:else if compatibility.state !== CompatibilityState.Works}
    <Markdown class="[&>p]:my-0" markdown={compatibility.note ?? ''}/>
  {/if}
</Tooltip>

