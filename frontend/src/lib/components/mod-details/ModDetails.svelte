<script lang="ts">
  import type { SizeOptions } from '@floating-ui/dom';
  import { mdiCheck, mdiChevronDown, mdiImport, mdiRocketLaunch, mdiTestTube, mdiWeb } from '@mdi/js';
  import { getContextClient, queryStore } from '@urql/svelte';
  import { SemVer, coerce, minVersion, parse, sort, validRange } from 'semver';

  import ModChangelog from '../modals/ModChangelog.svelte';

  import ModDetailsEntry from './ModDetailsEntry.svelte';

  import Markdown from '$lib/components/Markdown.svelte';
  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import Tooltip from '$lib/components/Tooltip.svelte';
  import { CompatibilityState, GetModDetailsDocument } from '$lib/generated';
  import { type PopupSettings, getModalStore, popup } from '$lib/skeletonExtensions';
  import { addQueuedModAction } from '$lib/store/actionQueue';
  import { canModify, lockfileMods, manifestMods } from '$lib/store/ficsitCLIStore';
  import { error, expandedMod, siteURL } from '$lib/store/generalStore';
  import { search } from '$lib/store/modFiltersStore';
  import { offline } from '$lib/store/settingsStore';
  import { bytesToAppropriate } from '$lib/utils/dataFormats';
  import { getAuthor } from '$lib/utils/getModAuthor';
  import { InstallModVersion, OfflineGetMod } from '$wailsjs/go/ficsitcli/ficsitCLI';
  import type { ficsitcli } from '$wailsjs/go/models';
  import { BrowserOpenURL } from '$wailsjs/runtime/runtime';

  export let focusOnEntry: HTMLElement | undefined = undefined;

  const modalStore = getModalStore();

  const client = getContextClient();

  $: modQuery = queryStore(
    {
      query: GetModDetailsDocument,
      client,
      pause: !$expandedMod || !!$offline,
      variables: {
        modReference: $expandedMod ?? '',
      },
    },
  );

  interface OfflineMod {
    offline: true;
    mod_reference: string;
    name: string;
    authors: {
      role: 'creator';
      user: {
        username: string;
      }
    }[];
    logo?: string;
    versions: ficsitcli.ModVersion[];
  }
  
  let offlineMod: OfflineMod = {
    offline: true,
    mod_reference: '',
    name: '',
    versions: [],
    authors: [],
    logo: undefined,
  };

  $: if($offline && $expandedMod) {
    OfflineGetMod($expandedMod).then((mod) => {
      offlineMod = {
        ...mod,
        authors: mod.authors.map((author) => ({
          role: 'creator',
          user: {
            username: author,
          },
        })),
        offline: true,
      };
    });
  }

  $: mod = $offline ? offlineMod : ($modQuery.fetching ? null : $modQuery.data?.mod);

  $: actualLogo = (mod && 'offline' in mod) ? (mod?.logo ? `data:image/png;base64, ${mod?.logo}` : '/images/no_image.webp') : mod?.logo;
  $: renderedLogo = actualLogo || `${$siteURL}/images/no_image.webp`;
  $: author = getAuthor(mod);

  $: size = mod ? bytesToAppropriate(mod.versions[0]?.size ?? 0) : undefined;

  $: latestVersion = mod ? (mod.versions.length ? sort(mod.versions.map((v) => parse(v.version) ?? coerce(v.version)).filter((v) => !!v) as SemVer[]).reverse()[0] : 'N/A') : undefined;
  $: installedVersion = mod ? ($lockfileMods[mod.mod_reference]?.version ?? 'Not installed') : undefined;

  $: ficsitAppLink = `${$siteURL}/mod/${$expandedMod}`;

  function colorForCompatibilityState(state?: CompatibilityState) {
    switch(state) {
      case CompatibilityState.Broken:
        return 'text-error-500';
      case CompatibilityState.Damaged:
        return 'text-warning-500';
      case CompatibilityState.Works:
        return 'text-primary-700';
    }
    return '';
  }

  $: manifestVersion = mod && $manifestMods[mod.mod_reference]?.version;
  async function installVersion(version: string | null) {
    if(!mod) {
      return;
    }
    
    const modReference = mod.mod_reference;
    const action = async () => InstallModVersion(modReference, version ?? '>=0.0.0').catch((e) => $error = e);
    const actionName = 'install-version';
    return addQueuedModAction(
      modReference,
      actionName,
      action,
    );
  }

  function close() {
    $expandedMod = null;
  }

  $: authorClick = () => {
    if (!author) {
      return;
    }
    $search = `author:"${author}"`;
  };

  const compatEAPopupId = 'mod-details-compat-ea';

  const compatEAPopup = {
    event: 'hover',
    target: compatEAPopupId,
    middleware: {
      offset: 4,
    },
    placement: 'bottom-end',
  } satisfies PopupSettings;

  const compatEXPPopupId = 'mod-details-compat-exp';

  const compatEXPPopup = {
    event: 'hover',
    target: compatEXPPopupId,
    middleware: {
      offset: 4,
    },
    placement: 'bottom-start',
  } satisfies PopupSettings;

  const compatUnknownPopupId = 'mod-details-compat-unknown';

  const compatUnknownPopup = {
    event: 'hover',
    target: compatUnknownPopupId,
    middleware: {
      offset: 4,
    },
    placement: 'bottom',
  } satisfies PopupSettings;

  const authorsMenuPopupId = 'mod-details-authors-menu';

  const authorsMenu = {
    event: 'click',
    target: authorsMenuPopupId,
    middleware: {
      offset: 4,
    },
    placement: 'bottom',
  } satisfies PopupSettings;

  const changeVersionMenuPopupId = 'mod-details-change-version-menu';

  const changeVersionMenu = {
    event: 'click',
    target: changeVersionMenuPopupId,
    middleware: {
      offset: 4,
      size: {
        apply({ availableHeight, elements }) {
          Object.assign(elements.floating.style, {
            maxHeight: `calc(${availableHeight}px - 1rem)`,
          });
        },
      } as SizeOptions,
    },
    placement: 'bottom',
  } satisfies PopupSettings;

  const changelogMenuPopupId = 'mod-details-changelog-menu';

  const changelogMenu = {
    event: 'click',
    target: changelogMenuPopupId,
    middleware: {
      offset: 4,
      size: {
        apply({ availableHeight, elements }) {
          Object.assign(elements.floating.style, {
            maxHeight: `calc(${availableHeight}px - 1rem)`,
          });
        },
      } as SizeOptions,
    },
    placement: 'bottom',
  } satisfies PopupSettings;
</script>

<div class="@container/mod-details h-full flex w-full bg-surface-200-700-token @3xl/mod-details:text-base text-sm">
  <div style="border-right-color: rgba(239, 239, 239, 0.12);" class="p-4 space-y-4 flex flex-col h-full @3xl/mod-details:w-64 w-52 overflow-y-auto">
    <div class="flex flex-col flex-auto overflow-y-auto space-y-4">
      <div class="flex flex-col">
        <img class="logo w-full" alt="{mod?.name} Logo" src={renderedLogo} />
        <span 
          class="font-bold @3xl/mod-details:text-lg mt-2 text-base"
          class:animate-pulse={!mod}
          class:placeholder={!mod}
        >
          {mod?.name ?? ''}
        </span>
        <span class="font-light">A mod by:</span>
        <span
          bind:this={focusOnEntry}
          class="font-medium text-primary-600 cursor-pointer"
          class:animate-pulse={!mod}
          class:placeholder={!mod}
          role="button"
          tabindex="0"
          on:click={authorClick}
          on:keypress={authorClick} 
        >
          {author ?? ''}
        </span>
      </div>
    
      <div use:popup={authorsMenu}>
        <button class="btn px-4 h-10 text-sm w-full bg-secondary-600">
          <span class="whitespace-break-spaces">Contributors <span class="text-primary-600">
            (<span class:animate-pulse={!mod} class:placeholder={!mod}>{mod ? (mod.authors.length ?? 0) : '    '}</span>)</span></span>
          <SvgIcon
            class="h-5 w-5"
            icon={mdiChevronDown}/>
        </button>
      </div>
      <div class="card shadow-xl min-w-[11rem] z-10 duration-0 overflow-y-auto max-h-[95vh] !mt-0" data-popup={authorsMenuPopupId}>
        <!-- 
        Skeleton's popup close function waits for the tranistion duration...
        before actually triggering the transition...
        So we'll just not have a transition...
        -->
      
        <ul>
          {#each mod?.authors ?? [] as author}
            <li>
              <button class="btn w-full h-full space-x-4" on:click={() => $search = `author:"${author.user.username}"`}>
                <div class="h-12 w-12">
                  {#if 'avatar' in author.user}
                    <img class="rounded-full w-ful h-full" alt="{author.user.username} Avatar" src={author.user.avatar} />
                  {/if}
                </div>
                <div class="flex-auto flex flex-col text-left">
                  <span>{author.user.username}</span>
                  <span>{author.role}</span>
                </div>
              </button>
            </li>
          {/each}
        </ul>
      </div>

      <div>
        <span>Mod info:</span><br/>
        <ModDetailsEntry label="Size" loading={!mod}>{size ?? ''}</ModDetailsEntry>
        {#if (!mod || !('offline' in mod)) && !$offline}
          <ModDetailsEntry label="Created" loading={!mod}>{mod ? new Date(mod.created_at).toLocaleDateString() : ''}</ModDetailsEntry>
          <ModDetailsEntry label="Updated" loading={!mod}>{mod ? new Date(mod.last_version_date).toLocaleString() : ''}</ModDetailsEntry>
          <ModDetailsEntry label="Total downloads" loading={!mod}>{mod?.downloads.toLocaleString() ?? ''}</ModDetailsEntry>
          <ModDetailsEntry label="Views" loading={!mod}>{mod?.views.toLocaleString() ?? ''}</ModDetailsEntry>
          <ModDetailsEntry label="Compatibility" loading={!mod}>
            {#if mod?.compatibility}
              <div class="flex pl-1">
                <div use:popup={compatEAPopup}>
                  <SvgIcon class="{colorForCompatibilityState(mod.compatibility.EA.state)} w-5" icon={mdiRocketLaunch} />
                </div>
              
                <Tooltip popupId={compatEAPopupId}>
                  <span class="text-base">
                    This mod has been reported as {mod.compatibility.EA.state} on Early Access.
                  </span>
                  {#if mod.compatibility.EA.note}
                    <Markdown class="[&>p]:my-0" markdown={mod.compatibility.EA.note} />
                  {:else}
                    (No further notes provided)
                  {/if}
                </Tooltip>
                <div use:popup={compatEXPPopup}>
                  <SvgIcon class="{colorForCompatibilityState(mod.compatibility.EXP.state)} w-5" icon={mdiTestTube} />
                </div>
                <Tooltip popupId={compatEXPPopupId}>
                  <span class="text-base">
                    This mod has been reported as {mod.compatibility.EXP.state} on Experimental.
                  </span>
                  {#if mod.compatibility.EXP.note}
                    <Markdown class="[&>p]:my-0" markdown={mod.compatibility.EXP.note} />
                  {:else}
                    (No further notes provided)
                  {/if}
                </Tooltip>
              </div>
            {:else if mod}
              <span class="font-bold" use:popup={compatUnknownPopup}>&nbsp;Unknown</span>
              <Tooltip popupId={compatUnknownPopupId}>
                <span class="text-base">No compatibility information has been reported for this mod yet. Try it out and contact us on the Discord so it can be updated!</span>
              </Tooltip>
            {/if}
          </ModDetailsEntry>
        {/if}
      </div>

      <div>
        <ModDetailsEntry label="Latest version" loading={!mod}>{latestVersion ?? ''}</ModDetailsEntry>
        <ModDetailsEntry label="Installed version" loading={!mod}>{installedVersion ?? ''}</ModDetailsEntry>
        <div class="pt-2" use:popup={changeVersionMenu}>
          <button
            class="btn px-4 h-10 text-sm w-full bg-secondary-600"
            disabled={!$canModify}
          >
            <span>Change version</span>
            <SvgIcon
              class="h-5 w-5"
              icon={mdiChevronDown}/>
          </button>
        </div>
        <div class="card shadow-xl min-w-[11rem] z-10 duration-0 overflow-y-auto !mt-0" data-popup={changeVersionMenuPopupId}>
          <!-- 
          Skeleton's popup close function waits for the tranistion duration...
          before actually triggering the transition...
          So we'll just not have a transition...
          -->
        
          <ul>
            <li>
              <button class="btn w-full h-full text-left" on:click={() => installVersion(null)}>
                <div class="w-7 h-7 p-1">
                  {#if manifestVersion === '>=0.0.0'}
                    <SvgIcon class="h-full w-full" icon={mdiCheck} />
                  {/if}
                </div>
                <span class="flex-auto">Any</span>
              </button>
            </li>
            {#each mod?.versions ?? [] as version}
              <li class="flex">
                <button class="btn w-full h-full text-left" on:click={() => installVersion(version.version)}>
                  <div class="w-7 h-7 p-1">
                    {#if manifestVersion === version.version}
                      <SvgIcon class="h-full w-full" icon={mdiCheck} />
                    {/if}
                  </div>
                  <span class="flex-auto">{version.version}</span>
                </button>
                <button class="btn w-full h-full text-left" on:click={() => installVersion(`>=${version.version}`)}>
                  <span class="flex-auto">or newer</span>
                  <div class="w-7 h-7 p-1">
                    {#if manifestVersion && manifestVersion !== version.version && validRange(manifestVersion) && minVersion(manifestVersion)?.format() === version.version}
                      <SvgIcon class="h-full w-full" icon={mdiCheck} />
                    {/if}
                  </div>
                </button>
              </li>
            {/each}
          </ul>
        </div>
        {#if (!mod || !('offline' in mod)) && !$offline}
          <div class="pt-2" use:popup={changelogMenu}>
            <button class="btn px-4 h-10 text-sm w-full bg-secondary-600">
              <span>Changelogs</span>
              <SvgIcon
                class="h-5 w-5"
                icon={mdiChevronDown}/>
            </button>
          </div>
          <div class="card shadow-xl min-w-[11rem] z-10 duration-0 overflow-y-auto !mt-0" data-popup={changelogMenuPopupId}>
            <!-- 
            Skeleton's popup close function waits for the tranistion duration...
            before actually triggering the transition...
            So we'll just not have a transition...
            -->
          
            <ul>
              {#each mod?.versions ?? [] as version}
                <li>
                  <button class="btn w-full h-full" on:click={() => modalStore.trigger({ type:'component', component: { ref: ModChangelog, props:{ mod: mod?.mod_reference, versionRange: version.version } } })}>
                    <span>{version.version}</span>
                  </button>
                </li>
              {/each}
            </ul>
          </div>
        {/if}
        <div class="pt-2">
          <button
            class="btn px-4 h-10 text-sm w-full bg-primary-900"
            on:click={() => BrowserOpenURL(ficsitAppLink)}>
            <span class="whitespace-break-spaces">View on ficsit.app</span>
            <SvgIcon
              class="h-5 w-5"
              icon={mdiWeb}/>
          </button>
        </div>
      </div>
    </div>

    <button
      class="btn px-4 h-8 w-full bg-secondary-600 text-sm"
      on:click={close}>
      <SvgIcon
        class="h-5 w-5 -scale-x-100"
        icon={mdiImport}/>
      <span>Close</span>
    </button>
  </div>
  <div class="break-words overflow-wrap-anywhere flex-1 px-3 mr-3 my-4 overflow-y-scroll overflow-x-hidden w-0">
    {#if $offline}
      <div class="flex items-center justify-center h-full text-center font-bold">Offline mode is enabled. Changelogs and descriptions are not available.</div>
    {:else if mod && 'full_description' in mod && mod.full_description}
      <Markdown markdown={mod.full_description} />
    {:else}
      <div class="p-4 space-y-4 h-full flex flex-col">
        <div class="placeholder" />
        <div class="grid grid-cols-3 gap-4 md:gap-8">
          <div class="placeholder" />
          <div class="placeholder" />
          <div class="placeholder" />
        </div>
        <div class="grid grid-cols-4 gap-2 md:gap-4">
          <div class="placeholder" />
          <div class="placeholder" />
          <div class="placeholder" />
          <div class="placeholder" />
        </div>
        <div class="placeholder grow" />
      </div>
    {/if}
  </div>
</div>

<style>
  .overflow-wrap-anywhere {
    overflow-wrap: anywhere;
  }
</style>
