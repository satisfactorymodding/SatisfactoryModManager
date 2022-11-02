<script lang="ts">
  import { mdiDownload, mdiEye, mdiStar, mdiCheckCircle, mdiPlay, mdiPause, mdiTrashCan, mdiTrayFull, mdiTrayMinus } from '@mdi/js';
  import MDIIcon from '$lib/components/MDIIcon.svelte';
  import { createEventDispatcher } from 'svelte';
  import { search, type PartialMod } from '$lib/store/modFiltersStore';
  import Button, { Group, GroupItem, Label } from '@smui/button';
  import LinearProgress from '@smui/linear-progress';
  import Tooltip, { Wrapper } from '@smui/tooltip';
  import { addQueuedModAction, queuedMods, favouriteMods, lockfileMods, manifestMods, progress, selectedInstall, removeQueuedModAction } from '$lib/store/ficsitCLIStore';
  import { error } from '$lib/store/generalStore';
  import { DisableMod, EnableMod, InstallMod, RemoveMod } from '$wailsjs/go/bindings/FicsitCLI';
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
  
  export let compact: boolean;
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

  $: installButtonLabel = (() => {
    if (isDependency) {
      return 'Dependency';
    }
    let prefix = '';
    if (queuedInstall) {
      prefix = 'Queued ';
    }
    if (isInstalled) {
      return `${prefix}Remove`;
    }
    return `${prefix}Install`;
  })();
  
  $: installButtonLabelHover = (() => {
    if (queuedInstall) {
      return 'Unqueue';
    }
    return installButtonLabel;
  })();
  
  $: installButtonIcon = (() => {
    if (isDependency) {
      return mdiCheckCircle;
    }
    if (queuedInstall) {
      return mdiTrayFull;
    }
    if (isInstalled) {
      return mdiTrashCan;
    }
    return mdiDownload;
  })();
  
  $: installButtonIconHover = (() => {
    if (isDependency) {
      return mdiCheckCircle;
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

<div class="my-1 px-0 w-full" class:rounded-lg={selected} class:selected on:click={click}>
  {#if inProgress}
    <div class="relative h-0">
      <LinearProgress progress={$progress?.progress} class="mod-progress-bar h-24 w-full rounded-lg"/>
    </div>
  {/if}
  <div class="flex relative" readonly class:disabled={isInstalled && !isEnabled}>
    <img src={renderedLogo} alt="{mod.name} Logo" class="logo h-24 w-24" />
    <div class="ml-2 flex flex-col grow w-0">
      <Wrapper>
        <div>
          <span class="text-xl font-medium" class:error={compatibility.state === CompatibilityState.Broken} class:warning={compatibility.state === CompatibilityState.Damaged}>{mod.name}</span>
          {#if !compact}
            <span class="pl-1">by</span>
            <span class="color-primary" on:click|stopPropagation={() => $search = `author:"${author}"`}>{author}</span>
          {/if}
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
      <span class="truncate w-full">{mod.short_description}</span>
      <div class="flex grow">
        {#if !inProgress}
          <div class="grow w-0">
            <div class="truncate">
              {#each mod?.tags ?? [] as tag}
                <span class="pr-1">#{tag.name}</span>
              {/each}
              &nbsp; <!-- keep div height even when no tags are available -->
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
          <div class="pr-2 flex" on:click|stopPropagation={() => { /* empty */ }}>
            <Group variant="outlined" class="mr-1">
              <Button
                on:click={toggleModInstalled}
                on:mouseover={() => isInstallButtonHovered = true}
                on:mouseleave={() => isInstallButtonHovered = false}
                variant="unelevated"
                disabled={buttonDisabled || queuedEnable}
                class="{compact ? ((isInstalled ? 'min-w-0' : 'w-[85px]') + ' pl-5 pr-1.5') : (isInstalled ? 'w-24' : 'w-36')} mod-install-button {isInstalled ? 'installed' : ''}">
                {#if !compact}
                  <Label>{ isInstallButtonHovered ? installButtonLabelHover : installButtonLabel }</Label>
                {:else}
                  <MDIIcon icon={ isInstallButtonHovered ? installButtonIconHover : installButtonIcon }/>
                {/if}
                <div class="grow" />
              </Button>
              {#if isInstalled && !isDependency}
                <div use:GroupItem>
                  <Button
                    on:click={ toggleModEnabled }
                    on:mouseover={() => isEnableButtonHovered = true}
                    on:mouseleave={() => isEnableButtonHovered = false}
                    disabled={buttonDisabled || queuedInstall}
                    variant="unelevated"
                    class="min-w-0 pl-5 pr-1.5 mod-enable-button {isEnabled ? 'enabled' : ''}"
                  >
                    <MDIIcon icon={ isEnableButtonHovered ? enableButtonIconHover : enableButtonIcon }/>
                  </Button>
                </div>
              {/if}
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
  .disabled {
    opacity: 0.3;
  }
  .disabled img {
    filter: grayscale(1);
    animation-play-state: paused !important;
  }
</style>