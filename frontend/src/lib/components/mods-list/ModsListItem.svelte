<script lang="ts">
  import { mdiDownload, mdiEye, mdiStar, mdiCheckCircle, mdiPlusBoxMultiple, mdiPlay, mdiPause, mdiTrashCan } from '@mdi/js';
  import MDIIcon from '$lib/components/MDIIcon.svelte';
  import { createEventDispatcher } from 'svelte';
  import { search, type PartialMod } from '$lib/store/modFiltersStore';
  import Button, { Group, GroupItem, Label } from '@smui/button';
  import LinearProgress from '@smui/linear-progress';
  import Tooltip, { Wrapper } from '@smui/tooltip';
  import { favouriteMods, lockfileMods, manifestMods, progress, selectedInstall } from '$lib/store/ficsitCLIStore';
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

  $: installButtonLabel = isDependency ? 'Dependency' : (isInstalled ? 'Remove' : 'Install');
  $: installButtonIcon = isDependency ? mdiCheckCircle : (isInstalled ? mdiCheckCircle : mdiDownload);
  $: installButtonIconHover = isDependency ? mdiCheckCircle : (isInstalled ? mdiTrashCan : mdiDownload);
  $: enableButtonIcon = isInstalled ? (isEnabled ? mdiPlay : mdiPause) : '';
  $: enableButtonIconHover = isInstalled ? (isEnabled ? mdiPause : mdiPlay) : '';
  $: buttonDisabled = isDependency || (!!$progress) || (versionCompatibility.state === CompatibilityState.Broken && !isInstalled);

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

  function toggleModInstalled() {
    if(isInstalled) {
      RemoveMod(mod.mod_reference);
    } else {
      InstallMod(mod.mod_reference);
    }
  }

  function toggleModEnabled() {
    if(isEnabled) {
      DisableMod(mod.mod_reference);
    } else {
      EnableMod(mod.mod_reference);
    }
  }

  function toggleModFavourite() {
    if(!isFavourite) {
      FavouriteMod(mod.mod_reference);
    } else {
      UnFavouriteMod(mod.mod_reference);
    }
  }

  function toggleModQueueInstall() {
    console.log('test');
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
          <Tooltip class="mods-list" surface$class="max-w-lg text-base">
            This mod is disabled. Press the pause icon to enable it. 
          </Tooltip>
        {:else if compatibility.state !== CompatibilityState.Works}
          <Tooltip class="mods-list" surface$class="max-w-lg text-base">
            { @html compatibility.note }
          </Tooltip>
        {/if}
      </Wrapper>
      <span class="truncate w-full">{mod.short_description}</span>
      <div class="flex grow">
        {#if !inProgress}
          <div class="grow w-0">
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
          <div class="pr-2 flex" on:click|stopPropagation={null}>
              <Group variant="outlined" class="mr-1">
                <Button
                  on:click={toggleModInstalled}
                  on:mouseover={() => isInstallButtonHovered = true}
                  on:mouseleave={() => isInstallButtonHovered = false}
                  variant="unelevated"
                  disabled={buttonDisabled}
                  class="{compact ? 'min-w-0 pl-5 pr-1.5' : 'w-28'} mod-install-button {isInstalled ? 'installed' : ''}">
                  {#if !compact}
                    <Label>{ installButtonLabel }</Label>
                  {:else}
                    <MDIIcon icon={ isInstallButtonHovered ? installButtonIconHover : installButtonIcon }/>
                  {/if}
                </Button>
                <div use:GroupItem>
                  <Button
                    on:click={ isInstalled ? toggleModEnabled : toggleModQueueInstall }
                    on:mouseover={() => isEnableButtonHovered = true}
                    on:mouseleave={() => isEnableButtonHovered = false}
                    disabled={buttonDisabled}
                    variant="unelevated"
                    class="min-w-0 pl-5 pr-1.5 {isInstalled ? ('mod-enable-button ' + (isEnabled ? 'enabled' : '')) : 'mod-install-button'}"
                  >
                    <MDIIcon icon={ isInstalled ? (isEnableButtonHovered ? enableButtonIconHover : enableButtonIcon) : mdiPlusBoxMultiple }/>
                  </Button>
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
  .disabled {
    opacity: 0.3;
  }
  .disabled img {
    filter: grayscale(1);
    animation-play-state: paused !important;
  }
</style>