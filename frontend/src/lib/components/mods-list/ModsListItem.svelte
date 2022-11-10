<script lang="ts">
  import { mdiDownload, mdiEye, mdiStar, mdiCheckCircle, mdiPlay, mdiPause, mdiTrashCan, mdiTrayFull, mdiTrayMinus, mdiSync } from '@mdi/js';
  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { createEventDispatcher } from 'svelte';
  import { search, type PartialMod } from '$lib/store/modFiltersStore';
  import Button from '@smui/button';
  import LinearProgress from '@smui/linear-progress';
  import Tooltip, { Wrapper } from '@smui/tooltip';
  import { addQueuedModAction, queuedMods, favouriteMods, lockfileMods, manifestMods, progress, selectedInstall, removeQueuedModAction } from '$lib/store/ficsitCLIStore';
  import { error } from '$lib/store/generalStore';
  import { DisableMod, EnableMod, InstallMod, RemoveMod } from '$wailsjs/go/ficsitcli_bindings/FicsitCLI';
  import { FavouriteMod, UnFavouriteMod } from '$wailsjs/go/bindings/Settings';
  import { getAuthor } from '$lib/utils/getModAuthor';
  import { getReportedCompatibility, getVersionCompatibility } from '$lib/utils/modCompatibility';
  import type { GameBranch } from '$lib/wailsTypesExtensions';
  import { CompatibilityState, type Compatibility } from '$lib/generated';
  import { markdown } from '$lib/utils/markdown';
  
  export let mod: PartialMod;

  const dispatch = createEventDispatcher();

  function click() {
    dispatch('click');
  }
  
  export let selected: boolean;

  $: renderedLogo = mod.logo || 'https://ficsit.app/images/no_image.webp';
  $: author = getAuthor(mod);

  $: isInstalled = mod.mod_reference in $manifestMods;
  $: isEnabled = mod.mod_reference in $lockfileMods;
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

  $: buttonDisabled = isDependency || (versionCompatibility.state === CompatibilityState.Broken && !isInstalled);

  let isInstallButtonHovered = false;
  let isEnableButtonHovered = false;

  $: isFavourite = $favouriteMods.includes(mod.mod_reference);

  let reportedCompatibility: Compatibility | undefined = { state: CompatibilityState.Works };
  let versionCompatibility: Compatibility = { state: CompatibilityState.Works };
  $: {
    if($selectedInstall && $selectedInstall.info) {
      reportedCompatibility = getReportedCompatibility(mod, $selectedInstall.info.branch as GameBranch);
      if(reportedCompatibility) {
        reportedCompatibility = {
          state: reportedCompatibility.state,
          note: reportedCompatibility.note 
            ? `This mod has been reported as ${reportedCompatibility.state} on this game version.<br>${markdown(reportedCompatibility.note)}` 
            : `This mod has been reported as ${reportedCompatibility.state} on this game version.`
        };
      }
      if (mod.hidden && !isDependency) {
        versionCompatibility = { state: CompatibilityState.Broken, note: 'This mod was hidden by the author.' };
      }
      else {
        getVersionCompatibility(mod.mod_reference, $selectedInstall.info.version).then((compatibility) => {
          versionCompatibility = compatibility;
        });
      }
    }
  }

  $: compatibility = reportedCompatibility ?? versionCompatibility;

  async function toggleModInstalled() {
    // Svelte does not recreate the component, but reuse it, so the associated mod reference might change
    const modReference = mod.mod_reference;
    const action = isInstalled ? async () => RemoveMod(modReference) : async () => InstallMod(modReference);
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
    const action = isEnabled ? async () => DisableMod(modReference) : async () => EnableMod(modReference);
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

  async function toggleModFavourite() {
    try {
      if(!isFavourite) {
        await FavouriteMod(mod.mod_reference);
      } else {
        await UnFavouriteMod(mod.mod_reference);
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
    <div class="relative h-0">
      <LinearProgress progress={$progress?.progress} class="mod-progress-bar h-full rounded-lg"/>
    </div>
  {/if}
  <div class="flex relative h-full" readonly class:disabled={isInstalled && !isEnabled}>
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
            { @html compatibility.note }
          </Tooltip>
        {/if}
      </Wrapper>
      <div class="truncate md:text-base text-sm hidden md:block">{mod.short_description}</div>
      <div class="flex">
        {#if !inProgress}
          <div class="grow w-0 lg:text-base text-sm">
            <div class="truncate text-base md:text-sm block md:hidden">{mod.short_description}</div>
            <div class="truncate hidden md:block">
              {#each mod?.tags ?? [] as tag}
                <span class="pr-1">#{tag.name}</span>
              {/each}
              &nbsp; <!-- keep div height even when no tags are available -->
            </div>
            <div class="flex h-5 md:h-4.5">
              <div class="w-24 flex items-center">
                <div class="pr-1 inline-flex items-center justify-items-center lg:w-7 w-6"><SvgIcon icon={mdiEye} class=""/></div><span>{mod.views.toLocaleString()}</span>
              </div>
              <div class="w-24 flex items-center">
                <div class="pr-1 inline-flex items-center justify-items-center lg:w-7 w-6"><SvgIcon icon={mdiDownload} class=""/></div><span>{mod.downloads.toLocaleString()}</span>
              </div>
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
          on:mouseover={() => isEnableButtonHovered = true}
          on:mouseleave={() => isEnableButtonHovered = false}
          disabled={buttonDisabled || queuedInstall || inProgress}
          ripple={false}
          variant="text"
          class="min-w-0 w-12 h-12 mod-enable-button {isEnabled || queued ? 'enabled' : ''} "
        >
          <SvgIcon icon={ isEnableButtonHovered ? enableButtonIconHover : enableButtonIcon } class="!p-1 !m-0 !w-full !h-full"/>
        </Button>
      {:else}
        <div class="min-w-0 w-12 h-12"/>
      {/if}
      <Button
        on:click={toggleModInstalled}
        on:mouseover={() => isInstallButtonHovered = true}
        on:mouseleave={() => isInstallButtonHovered = false}
        variant="text"
        disabled={buttonDisabled || queuedEnable || inProgress}
        ripple={false}
        class="min-w-0 w-12 h-12 mod-install-button {isInstalled || queued ? 'installed' : ''}"
      >
        <SvgIcon icon={ isInstallButtonHovered ? installButtonIconHover : installButtonIcon } class="!p-1 !m-0 !w-full !h-full"/>
      </Button>
      <Button
        on:click={toggleModFavourite}
        variant="text"
        ripple={false}
        class="min-w-0 w-12 h-12 mod-favourite-button {isFavourite ? 'favourite' : ''}"
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