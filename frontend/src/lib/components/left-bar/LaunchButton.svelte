<script lang="ts">
  import Tooltip, { Wrapper } from '@smui/tooltip';
  import { getContextClient } from '@urql/svelte';
  import { mdiOpenInNew, mdiTrayFull } from '@mdi/js';

  import { selectedInstallMetadata, isGameRunning, lockfileMods, progress } from '$lib/store/ficsitCLIStore';
  import { queuedMods, startQueue } from '$lib/store/actionQueue';
  import { launchButton, queueAutoStart } from '$lib/store/settingsStore';
  import { error, isLaunchingGame } from '$lib/store/generalStore';
  import { LaunchGame } from '$wailsjs/go/ficsitcli/FicsitCLI';
  import { CompatibilityState, type Compatibility } from '$lib/generated';
  import { getCompatibility, type CompatibilityWithSource } from '$lib/utils/modCompatibility';
  import SvgIcon from '$lib/components/SVGIcon.svelte';

  $: isInstallLaunchable = !!$selectedInstallMetadata?.launchPath;

  const client = getContextClient();

  let compatibilities: Record<string, CompatibilityWithSource> = {};
  $: {
    const gameVersion = $selectedInstallMetadata?.version;
    const branch = $selectedInstallMetadata?.branch;
    if (gameVersion && branch) {
      compatibilities = {};
      Object.keys($lockfileMods).map(async (modReference) => {
        if (modReference !== 'SML') {
          compatibilities[modReference] = await getCompatibility(modReference, branch, gameVersion, client);
        }
      });
    }
  }

  $: reportedCompatibilities = Object.entries(compatibilities).reduce((acc, [modReference, compatibility]) => {
    if (compatibility?.source === 'reported') {
      acc[modReference] = compatibility;
    }
    return acc;
  }, {} as Record<string, Compatibility>);

  $: versionCompatibilities = Object.entries(compatibilities).reduce((acc, [modReference, compatibility]) => {
    if (compatibility?.source === 'version') {
      acc[modReference] = compatibility;
    }
    return acc;
  }, {} as Record<string, Compatibility>);

  $: versionIncompatible = Object.keys($lockfileMods).filter((modReference) => versionCompatibilities[modReference]?.state === CompatibilityState.Broken);
  $: versionPossiblyCompatible = Object.keys($lockfileMods).filter((modReference) => versionCompatibilities[modReference]?.state === CompatibilityState.Damaged);
  $: reportedIncompatible = Object.keys($lockfileMods).filter((modReference) => reportedCompatibilities[modReference]?.state === CompatibilityState.Broken);
  $: reportedPossiblyCompatible = Object.keys($lockfileMods).filter((modReference) => reportedCompatibilities[modReference]?.state === CompatibilityState.Damaged);

  $: launchButtonColor = (versionIncompatible.length > 0 || reportedIncompatible.length > 0) ? 'error' : ((versionPossiblyCompatible.length > 0 || reportedPossiblyCompatible.length > 0) ? 'warning' : '');
  $: areOperationsQueued = !$queueAutoStart && $queuedMods.length > 0;

  function launchGame() {
    $isLaunchingGame = true;
    LaunchGame().catch((e) => $error = e);
    setTimeout(() => $isLaunchingGame = false, 10000);
  }

  let launchButtonState: 'press' | 'over' | 'normal' = 'normal';
  let launchButtonPressCount = 0;
  function launchButtonPressed() {
    launchButtonPressCount++;
    setTimeout(() => {
      if(launchButtonPressCount > 0) {
        launchButtonPressCount--;
      }
    }, 3000);
    if(launchButtonPressCount >= 15) {
      launchButtonPressCount = 0;
      launchGame();
    }
  }

  let catPressed = false;
  let catPosition = 0;
  let mouseDownX = 0;
  let mouseDownCatPosition = 0;
  
  function catMouseDown(e: MouseEvent) {
    mouseDownX = e.clientX;
    mouseDownCatPosition = catPosition;
    catPressed = true;
  }

  function catMouseMove(e: MouseEvent) {
    if (catPressed) {
      catPosition = (e.clientX - mouseDownX) / 270 + mouseDownCatPosition;
      catPosition = Math.min(1, Math.max(-0.05, catPosition));
      if (catPosition === 1) {
        catPressed = false;
        setTimeout(() => {
          launchGame();
        }, 1000);
      }
    }
  }
</script>

<Wrapper>
  <center>
    {#if areOperationsQueued}
      <button
        class="btn h-8 w-full text-sm bg-error-500"
        on:click={() => startQueue()}
      >
        <span>Apply {$queuedMods.length} change{$queuedMods.length !== 1 ? 's' : ''}</span>
        <div class="grow" />
        <SvgIcon
          class="h-5 w-5"
          icon={mdiTrayFull}/>
      </button>
    {:else if !isInstallLaunchable}
      <button
        class="btn h-8 w-full text-sm bg-surface-200-700-token"
        disabled
      >
        <span>SMM can't launch this install</span>
        <div class="grow" />
      </button>
    {:else if $launchButton === 'normal' || $isGameRunning || $isLaunchingGame }
      <button
        class="btn h-8 w-full text-sm bg-primary-900"
        disabled={!!$progress || $isGameRunning || $isLaunchingGame}
        on:click={() => launchGame()}
      >
        <span>Play Satisfactory</span>
        <div class="grow" />
        <SvgIcon
          class="h-5 w-5"
          icon={mdiOpenInNew}/>
      </button>
    {:else if $launchButton === 'cat' }
      <!-- fixme SMMv3 seems to have broken this button -->
      <!-- svelte-ignore a11y-no-static-element-interactions -->
      <div
        style="height: 32px"
        class="overflow-hidden"
        on:mouseup={() => catPressed = false}
        on:mousemove={catMouseMove}
      >
        <img
          src="/images/launch/cat/bg.png"
          draggable="false"
          alt="Space background"
        >
        <div
          on:mousedown={catMouseDown}
          role="button"
          tabindex="0"
        >
          <img
            src="/images/launch/cat/cat_full.png"
            style="position: relative; top: -57px; zoom: 0.55"
            style:left={`calc(-480px + ${catPosition * 87}%)`}
            draggable="false"
            alt="Cat"
          >
        </div>
        
      </div>
    {:else if $launchButton === 'button' }
      <!-- FIXME: keyboard navigation isn't allowing pressing this button with enter/space -->
      <div
        style="height: 50px"
        role="button"
        tabindex="0"
        on:keydown={launchButtonPressed}
      >
        <img
          src="/images/launch/fun/launch_fun.png"
          draggable="false"
          alt="Launch Button Background"
        >
        <!-- Keyboard interactions for the button are defined in the overall div -->
        <!-- svelte-ignore a11y-no-static-element-interactions -->
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div
          on:click={launchButtonPressed}
          on:mouseup={() => launchButtonState = 'over'}
          on:mousedown={() => launchButtonState = 'press'}
          on:mouseenter={() => launchButtonState = 'over'}
          on:mouseleave={() => launchButtonState = 'normal'}
        >
          <img
            src={`/images/launch/fun/launch_fun_button_${launchButtonState}.png`}
            style="position: relative; zoom: 0.56"
            style:top={launchButtonState === 'press' ? '-97.5px' : '-98px'}
            draggable="false"
            alt="Launch Button"
          >
        </div>
      </div>
    {/if}
  </center>
  {#if versionIncompatible.length > 0 || versionPossiblyCompatible.length > 0 || reportedIncompatible.length > 0 || reportedPossiblyCompatible.length > 0}
    <Tooltip surface$class="max-w-lg text-base">
      You have:
      <ul class="list-disc pl-5">
        {#if versionIncompatible.length > 0}
          <li>
            { versionIncompatible.length } incompatible mod{ versionIncompatible.length > 1 ? 's' : '' } which will either not load or crash your game
          </li>
        {/if}
        {#if reportedIncompatible.length > 0}
          <li>
            { reportedIncompatible.length } mod{ reportedIncompatible.length > 1 ? 's' : '' } that { reportedIncompatible.length > 1 ? 'are' : 'is' } reported as Broken on this game version.
            Read the mod{ reportedIncompatible.length > 1 ? 's\'' : '\'s' } description or compatibility notes for more information.
          </li>
        {/if}
        {#if versionPossiblyCompatible.length > 0}
          <li>
            { versionPossiblyCompatible.length } mod{ versionPossiblyCompatible.length > 1 ? 's' : '' }
            that { versionPossiblyCompatible.length > 1 ? 'are' : 'is' } likely incompatible with your game
          </li>
        {/if}
        {#if reportedPossiblyCompatible.length > 0}
          <li>
            { reportedPossiblyCompatible.length } mod{ reportedPossiblyCompatible.length > 1 ? 's' : '' }
            that { reportedPossiblyCompatible.length > 1 ? 'are' : 'is' } reported as Damaged on this game version.
            Read the mod{ reportedPossiblyCompatible.length > 1 ? 's\'' : '\'s' } description or compatibility notes for more information.
          </li>
        {/if}
      </ul>
      Are you sure you want to launch?
    </Tooltip>
  {:else if areOperationsQueued}
    <Tooltip surface$class="max-w-lg text-base">
      Changes have not yet been made to your mod files. Click the button above to apply the changes you have queued.<br/><br/>(You're in Queue "Start manually" mode)
    </Tooltip>
  {:else if $isGameRunning}
    <Tooltip surface$class="max-w-lg text-base">
      Your game launcher is reporting that the game is already running (or still in the process of closing).
    </Tooltip>
  {:else if $isLaunchingGame}
    <Tooltip surface$class="max-w-lg text-base">
      Launch in progress...
    </Tooltip>
  {:else if !!$progress}
    <Tooltip surface$class="max-w-lg text-base">
      An operation is already in progress.
    </Tooltip>
  {:else if !isInstallLaunchable}
    <Tooltip surface$class="max-w-lg text-base">
      The Mod Manager is not capable of launching this install type, but it will still manage the mod files for you. Launch Satisfactory using your usual game launcher.
    </Tooltip>
  {:else}
    <Tooltip surface$class="max-w-lg text-base">
      You're ready to rumble!<br/><br/>Note: The Mod Manager has already finished installing the mod files for you. You could launch the game using your usual game launcher and mods would still be loaded.
    </Tooltip>
  {/if}
</Wrapper>
