<script lang="ts">
  import Tooltip, { Wrapper } from '@smui/tooltip';
  import { selectedInstall, isGameRunning, lockfileMods, progress, queuedMods, startQueue } from '$lib/store/ficsitCLIStore';
  import { launchButton, queueAutoStart } from '$lib/store/settingsStore';
  import { isLaunchingGame } from '$lib/store/generalStore';
  import { LaunchGame } from '$wailsjs/go/bindings/FicsitCLI';
  import Button, { Label } from '@smui/button';
  import { getClient } from '@urql/svelte';
  import { CompatibilityState, ModReportedCompatibilityDocument, type Compatibility } from '$lib/generated';
  import { getReportedCompatibility, getVersionCompatibility } from '$lib/utils/modCompatibility';
  import type { GameBranch } from '$lib/wailsTypesExtensions';
  import { mdiCheckCircleOutline, mdiOpenInNew } from '@mdi/js';
  import MDIIcon from '$lib/components/MDIIcon.svelte';

  const client = getClient();

  let reportedCompatibilities: Record<string, Compatibility | undefined> = {};
  $: {
    const branch = $selectedInstall?.info?.branch as GameBranch;
    if(branch) {
      reportedCompatibilities = {};
      Object.keys($lockfileMods).map(async (modReference) => {
        const result = await client.query(ModReportedCompatibilityDocument, { modReference }).toPromise();
        if(!result.data?.getModByReference) {
          return;
        }
        reportedCompatibilities[modReference] = getReportedCompatibility(result.data.getModByReference, branch);
      });
    }
  }

  let versionCompatibilities: Record<string, Compatibility> = {};
  $: {
    const gameVersion = $selectedInstall?.info?.version;
    if(gameVersion) {
      versionCompatibilities = {};
      Object.keys($lockfileMods).map(async (modReference) => {
        versionCompatibilities[modReference] = await getVersionCompatibility(modReference, gameVersion);
      });
    }
  }

  $: versionIncompatible = Object.keys($lockfileMods).filter((modReference) => versionCompatibilities[modReference]?.state === CompatibilityState.Broken);
  $: versionPossiblyCompatible = Object.keys($lockfileMods).filter((modReference) => versionCompatibilities[modReference]?.state === CompatibilityState.Damaged);
  $: reportedIncompatible = Object.keys($lockfileMods).filter((modReference) => reportedCompatibilities[modReference]?.state === CompatibilityState.Broken);
  $: reportedPossiblyCompatible = Object.keys($lockfileMods).filter((modReference) => reportedCompatibilities[modReference]?.state === CompatibilityState.Damaged);

  $: launchButtonColor = (versionIncompatible.length > 0 || reportedIncompatible.length > 0) ? 'error' : ((versionPossiblyCompatible.length > 0 || reportedPossiblyCompatible.length > 0) ? 'warning' : '');

  function launchGame() {
    $isLaunchingGame = true;
    LaunchGame();
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
    {#if !$queueAutoStart && $queuedMods.length > 0}
      <Button variant="unelevated" class="h-12 w-full launch-game error" on:click={() => startQueue()}>
        <Label>Apply {$queuedMods.length} changes</Label>
        <div class="grow" />
        <MDIIcon icon={ mdiCheckCircleOutline }/>
      </Button>
    {:else if $launchButton === 'normal' || $isGameRunning || $isLaunchingGame }
      <Button variant="unelevated" class="h-12 w-full launch-game {launchButtonColor}" disabled={$progress || $isGameRunning || $isLaunchingGame} on:click={() => launchGame()}>
        <Label>Play Satisfactory</Label>
        <div class="grow" />
        <MDIIcon icon={ mdiOpenInNew }/>
      </Button>
    {:else if $launchButton === 'cat' }
      <div
        style="height: 32px"
        class="overflow-hidden"
        on:mouseup={() => catPressed = false}
        on:mousemove={catMouseMove}
      >
        <img
          src="/images/launch/cat/bg.png"
          draggable="false"
          alt="Nyan background"
        >
        <img
          src="/images/launch/cat/cat_full.png"
          style="position: relative; top: -57px; zoom: 0.55"
          style:left={`calc(-480px + ${catPosition * 87}%)`}
          draggable="false"
          on:mousedown={catMouseDown}
          alt="Nyan"
        >
      </div>
    {:else if $launchButton === 'button' }
      <div style="height: 50px">
        <img
          src="/images/launch/fun/launch_fun.png"
          draggable="false"
          alt="Launch Button Background"
        >
        <img
          src={`/images/launch/fun/launch_fun_button_${launchButtonState}.png`}
          style="position: relative; zoom: 0.56"
          style:top={launchButtonState === 'press' ? '-97.5px' : '-98px'}
          draggable="false"
          on:click={() => launchButtonPressed()}
          on:mousedown={() => launchButtonState = 'press'}
          on:mouseup={() => launchButtonState = 'over'}
          on:mouseenter={() => launchButtonState = 'over'}
          on:mouseleave={() => launchButtonState = 'normal'}
          alt="Launch Button"
        >
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
    {/if}
</Wrapper>