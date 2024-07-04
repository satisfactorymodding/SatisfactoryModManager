<script lang="ts">
  import { mdiOpenInNew, mdiTrayFull } from '@mdi/js';
  import { getContextClient } from '@urql/svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import T from '$lib/components/T.svelte';
  import Tooltip from '$lib/components/Tooltip.svelte';
  import { type Compatibility, CompatibilityState } from '$lib/generated';
  import { type PopupSettings, popup } from '$lib/skeletonExtensions';
  import { queuedMods, startQueue } from '$lib/store/actionQueue';
  import { isGameRunning, lockfileMods, progress, selectedInstallMetadata } from '$lib/store/ficsitCLIStore';
  import { error, hasFetchedMods, isLaunchingGame } from '$lib/store/generalStore';
  import { launchButton, queueAutoStart } from '$lib/store/settingsStore';
  import { type CompatibilityWithSource, getCompatibility } from '$lib/utils/modCompatibility';
  import { installTypeToTargetName } from '$lib/wailsTypesExtensions';
  import { LaunchGame } from '$wailsjs/go/ficsitcli/ficsitCLI';

  $: isInstallLaunchable = !!$selectedInstallMetadata?.info?.launchPath;

  const client = getContextClient();

  let compatibilities: Record<string, CompatibilityWithSource> = {};
  $: {
    const info = $selectedInstallMetadata?.info;
    if (info && $hasFetchedMods) {
      const newCompatibilities: typeof compatibilities = {};
      Object.keys($lockfileMods).map(async (modReference) => {
        if (modReference !== 'SML') {
          newCompatibilities[modReference] = await getCompatibility(modReference, info.branch, info.version, installTypeToTargetName(info.type), client);
        }
      });
      compatibilities = newCompatibilities;
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

  $: launchButtonError = versionIncompatible.length > 0 || reportedIncompatible.length > 0;
  $: launchButtonWarning = !launchButtonError && (versionPossiblyCompatible.length > 0 || reportedPossiblyCompatible.length > 0);
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

  const popupId = 'launch-button';

  const popupHover = {
    event: 'hover',
    target: popupId,
    middleware: {
      offset: 4,
    },
    placement: 'top-start',
  } satisfies PopupSettings;
</script>

<center use:popup={popupHover}>
  {#if areOperationsQueued}
    <button
      class="btn h-8 w-full text-sm bg-error-500"
      on:click={() => startQueue()}
    >
      <span>
        <T defaultValue={'Apply {queued, plural, one {one change} other {# changes}}'} keyName="launch-button.apply-queued" params={{ queued: $queuedMods.length }}/>
      </span>
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
      <span>
        <T defaultValue="SMM can't launch this install" keyName="launch-button.cant-launch"/>
      </span>
      <div class="grow" />
    </button>
  {:else if $launchButton === 'normal' || $isGameRunning || $isLaunchingGame}
    <button
      class="btn h-8 w-full text-sm"
      class:bg-error-500={launchButtonError}
      class:bg-primary-900={!launchButtonError && !launchButtonWarning}
      class:bg-warning-500={launchButtonWarning}
      disabled={!!$progress || $isGameRunning || $isLaunchingGame}
      on:click={() => launchGame()}
    >
      <span>
        <T defaultValue="Play Satisfactory" keyName="launch-button.play"/>
      </span>
      <div class="grow" />
      <SvgIcon
        class="h-5 w-5"
        icon={mdiOpenInNew}/>
    </button>
  {:else if $launchButton === 'cat'}
    <!-- fixme SMMv3 seems to have broken this button -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div
      style="height: 32px"
      class="overflow-hidden"
      on:mouseup={() => catPressed = false}
      on:mousemove={catMouseMove}
    >
      <img
        alt="Space background"
        draggable="false"
        src="/images/launch/cat/bg.png"
      />
      <div
        role="button"
        tabindex="0"
        on:mousedown={catMouseDown}
      >
        <img
          style="position: relative; top: -57px; zoom: 0.55"
          style:left={`calc(-480px + ${catPosition * 87}%)`}
          alt="Cat"
          draggable="false"
          src="/images/launch/cat/cat_full.png"
        />
      </div>
        
    </div>
  {:else if $launchButton === 'button'}
    <!-- FIXME: keyboard navigation isn't allowing pressing this button with enter/space -->
    <div
      style="height: 50px"
      role="button"
      tabindex="0"
      on:keydown={launchButtonPressed}
    >
      <img
        alt="Launch Button Background"
        draggable="false"
        src="/images/launch/fun/launch_fun.png"
      />
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
          style="position: relative; zoom: 0.56"
          style:top={launchButtonState === 'press' ? '-97.5px' : '-98px'}
          alt="Launch Button"
          draggable="false"
          src={`/images/launch/fun/launch_fun_button_${launchButtonState}.png`}
        />
      </div>
    </div>
  {/if}
</center>
<Tooltip class="!mt-0" {popupId}>
  {#if versionIncompatible.length > 0 || versionPossiblyCompatible.length > 0 || reportedIncompatible.length > 0 || reportedPossiblyCompatible.length > 0}
    <span>
      <T defaultValue="You have:" keyName="launch-button.you-have-warning-mods"/>
    </span>
    <ul class="list-disc pl-5">
      {#if versionIncompatible.length > 0}
        <li>
          <span>
            <T defaultValue={'{versionIncompatible} incompatible {versionIncompatible, plural, one {mod} other {mods}} which will not load or crash our game'} keyName="launch-button.incompatible-mods" params={{ versionIncompatible: versionIncompatible.length }}/>
          </span>
        </li>
      {/if}
      {#if reportedIncompatible.length > 0}
        <li>
          <span>
            <T
              defaultValue={'{reportedIncompatible, plural, one {One mod} other {# mods}} that {reportedIncompatible, plural, one {is} other {are}} reported as Broken on this game version. Read the {versionIncompatible, plural, one {mod\'s} other {mods\'}} description or compatibility notes for more information'}
              keyName="launch-button.reported-incompatible-mods"
              params={{ reportedIncompatible: reportedIncompatible.length }}/>
          </span>
        </li>
      {/if}
      {#if versionPossiblyCompatible.length > 0}
        <li>
          <span>
            <T
              defaultValue={'{versionPossiblyCompatible, plural, one {One mod} other {# mods}} that {versionPossiblyCompatible, plural, one {is} other {are}} likely incompatible with your game'} 
              keyName="launch-button.possibly-incompatible-mods"
              params={{ versionPossiblyCompatible: versionPossiblyCompatible.length }}/>
          </span>
        </li>
      {/if}
      {#if reportedPossiblyCompatible.length > 0}
        <li>
          <span>
            <T
              defaultValue={'{reportedPossiblyCompatible, plural, one {One mod} other {# mods}} that {reportedPossiblyCompatible, plural, one {is} other {are}} reported as Damaged on this game version. Read the {versionIncompatible, plural, one {mod\'s} other {mods\'}} description or compatibility notes for more information'}
              keyName="launch-button.reported-possibly-compatible-mods"
              params={{ reportedPossiblyCompatible: reportedPossiblyCompatible.length }}/>
          </span>
        </li>
      {/if}
    </ul>
    <span>
      <T defaultValue="Are you sure you want to launch?" keyName="launch-button.are-you-sure-warning"/>
    </span>
  {:else if areOperationsQueued}
    <span>
      <T defaultValue={'Changes have not yet been made to your mod files. Click the button above to apply the changes you have queued.\n\n(You\'re in Queue "Start manually" mode)'} keyName="launch-button.changes-queued"/>
    </span>
  {:else if $isGameRunning}
    <span>
      <T defaultValue="Your game launcher is reporting that the game is already running (or still in the process of closing)." keyName="launch-button.game-running"/>
    </span>
  {:else if $isLaunchingGame}
    <span>
      <T defaultValue="Launch in progress..." keyName="launch-button.launch-in-progress"/>
    </span>
  {:else if !!$progress}
    <span>
      <T defaultValue="An operation is already in progress." keyName="launch-button.operation-in-progress"/>
    </span>
  {:else if !isInstallLaunchable}
    <span>
      <T defaultValue="The Mod Manager is not capable of launching this install type, but it will still manage the mod files for you. Launch Satisfactory using your usual game launcher." keyName="launch-button.cant-launch-tooltip"/>
    </span>
  {:else}
    <span>
      <T defaultValue="You're ready to rumble!\n\nNote: The Mod Manager has already finished installing the mod files for you. You could launch the game using your usual game launcher and mods would still be loaded." keyName="launch-button.ready"/>
    </span>
  {/if}
</Tooltip>
