<script lang="ts">
  import { ProgressBar } from '@skeletonlabs/skeleton';
  import { getTranslate } from '@tolgee/svelte';
  import { getContextClient } from '@urql/svelte';
  import { onDestroy, onMount } from 'svelte';

  import T from '$lib/components/T.svelte';
  import ModPreview, { type ModPreviewData } from '$lib/components/mods-list/ModPreview.svelte';
  import { GetProfileModsDocument } from '$lib/generated';
  import { hasPendingProfileChange } from '$lib/store/actionQueue';
  import { canModify, checkForUpdates, profiles, progressMessage, progressPercent, selectedProfile } from '$lib/store/ficsitCLIStore';
  import { error } from '$lib/store/generalStore';
  import { offline } from '$lib/store/settingsStore';
  import { GetProfileDiff, MergeProfile, OfflineGetModsByReferences } from '$wailsjs/go/ficsitcli/ficsitCLI';
  import type { ficsitcli } from '$wailsjs/go/models';

  export let parent: { onClose: () => void };
  export let profile: string;

  const { t } = getTranslate();
  const client = getContextClient();

  let currentProfile = profile;
  let previousCurrentProfile = profile;
  $: currentProfiles = [...$profiles].sort();
  $: incomingProfiles = $profiles.filter((name) => name !== currentProfile).sort();
  let incomingProfile = $profiles.filter((name) => name !== profile).sort()[0] ?? '';
  let entries: ficsitcli.ProfileModDiff[] = [];
  let enabledMods: string[] = [];
  let selectedVersions: Record<string, string> = {};
  let previews: Record<string, ModPreviewData> = {};
  let loading = true;
  let merging = false;
  let comparisonError = '';
  let requestId = 0;

  $: selections = entries.map((entry): ficsitcli.ProfileModSelection => ({
    ...entry,
    enabled: enabledMods.includes(entry.mod),
    version: selectedVersions[entry.mod],
  })).filter((entry) => entry.current
    ? entry.enabled !== entry.current.enabled || entry.version !== entry.current.version
    : entry.enabled);

  function errorMessage(e: unknown) {
    return e instanceof Error ? e.message : typeof e === 'string' ? e : 'Unknown error';
  }

  $: previewFor = (reference: string): ModPreviewData => previews[reference]
    ?? { missing: true, mod_reference: reference, name: reference, authors: ['N/A'] };

  async function loadPreviews(references: string[], request: number) {
    if (!references.length) return;
    const loaded: Record<string, ModPreviewData> = {};
    if (!$offline) {
      const batches = Array.from({ length: Math.ceil(references.length / 100) }, (_, i) => references.slice(i * 100, (i + 1) * 100));
      const results = await Promise.allSettled(batches.map((modReferences) => client.query(GetProfileModsDocument, {
        modReferences,
        limit: modReferences.length,
      }).toPromise()));
      for (const result of results) {
        if (result.status === 'fulfilled') {
          for (const mod of result.value.data?.getMods.mods ?? []) loaded[mod.mod_reference] = mod;
        }
      }
      if (request !== requestId) return;
      previews = { ...loaded };
    }
    const missing = references.filter((reference) => !loaded[reference]);
    if (missing.length) {
      try {
        for (const mod of await OfflineGetModsByReferences(missing) ?? []) {
          loaded[mod.mod_reference] = { ...mod, logo: mod.logo ?? undefined, offline: true };
        }
      } catch (e) {
        // Metadata is optional: unavailable mods must still appear in the comparison.
        console.warn('Failed to load cached profile mod previews', e);
      }
    }
    if (request === requestId) previews = loaded;
  }

  async function loadComparison() {
    const request = ++requestId;
    entries = [];
    enabledMods = [];
    selectedVersions = {};
    previews = {};
    comparisonError = '';
    loading = true;
    try {
      if (incomingProfile) {
        const result = await GetProfileDiff(currentProfile, incomingProfile);
        if (request !== requestId) return;
        entries = result ?? [];
        enabledMods = entries.filter((entry) => entry.current?.enabled).map((entry) => entry.mod);
        selectedVersions = Object.fromEntries(entries.map((entry) => [entry.mod, entry.current?.version ?? entry.incoming?.version ?? '']));
        void loadPreviews(entries.map((entry) => entry.mod), request);
      }
    } catch (e) {
      if (request === requestId) comparisonError = errorMessage(e);
    } finally {
      if (request === requestId) loading = false;
    }
  }

  function currentProfileChanged() {
    if (currentProfile === incomingProfile) {
      incomingProfile = previousCurrentProfile;
    }
    previousCurrentProfile = currentProfile;
    void loadComparison();
  }

  async function finishMerge() {
    if (merging || loading || !$canModify || !selections.length) return;
    merging = true;
    comparisonError = '';
    try {
      const destination = currentProfile;
      const result = await MergeProfile(destination, incomingProfile, selections);
      if (destination === $selectedProfile) {
        $hasPendingProfileChange = !result.applied;
        checkForUpdates().catch(console.error);
      }
      parent.onClose();
      if (result.applyError) {
        $error = $t('profiles.merge.apply-error', 'The merged profile was saved, but could not be applied: {message}', { message: result.applyError });
      }
    } catch (e) {
      comparisonError = errorMessage(e);
    } finally {
      merging = false;
    }
  }

  onMount(() => { void loadComparison(); });
  onDestroy(() => { requestId++; });
</script>

<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="w-[64rem] card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    <T defaultValue="Compare / Merge profiles" keyName="profiles.merge.title" />
  </header>
  <section class="p-4 min-h-0 overflow-auto space-y-4">
    <p class="text-sm">
      <T defaultValue="Checked means enabled in Current; unchecked means disabled or not present. Incoming is read-only. Choose a version separately when constraints differ." keyName="profiles.merge.enabled-help" />
    </p>
    {#if comparisonError}
      <p class="text-error-500 break-words" role="alert">{comparisonError}</p>
    {/if}
    <div class="flex flex-wrap gap-2">
      <button class="btn btn-sm variant-filled-surface" disabled={loading || merging || !$canModify} on:click={loadComparison}>
        <T defaultValue="Refresh" keyName="profiles.merge.refresh" />
      </button>
      <button class="btn btn-sm variant-filled-surface" disabled={loading || merging || !entries.length || !$canModify} on:click={() => enabledMods = entries.map((entry) => entry.mod)}>
        <T defaultValue="Enable all" keyName="profiles.merge.enable-all" />
      </button>
      <button class="btn btn-sm variant-filled-surface" disabled={loading || merging || !enabledMods.length || !$canModify} on:click={() => enabledMods = []}>
        <T defaultValue="Disable all" keyName="profiles.merge.disable-all" />
      </button>
    </div>
    <div class="@container/mods-list">
      <table class="w-full min-w-[40rem] table-fixed border-separate border-spacing-x-0 border-spacing-y-2" aria-label={$t('profiles.merge.mods', 'Profile mods')}>
        <colgroup>
          <col />
          <col class="w-36" />
          <col class="w-36" />
        </colgroup>
        <thead>
          <tr>
            <th class="sticky top-0 z-10 bg-surface-50-900-token px-2 py-3 text-left" scope="col">
              <T defaultValue="Mod" keyName="profiles.merge.mod" />
            </th>
            <th class="sticky top-0 z-10 bg-surface-50-900-token px-2 py-3 text-center" scope="col">
              <label class="label">
                <span><T defaultValue="Current" keyName="profiles.merge.current" /></span>
                <select
                  class="select text-sm font-normal w-full"
                  aria-label={$t('profiles.merge.current', 'Current')}
                  disabled={merging || !$canModify}
                  title={currentProfile}
                  bind:value={currentProfile}
                  on:change={currentProfileChanged}>
                  {#each currentProfiles as name}
                    <option value={name}>{name}</option>
                  {/each}
                </select>
              </label>
            </th>
            <th class="sticky top-0 z-10 bg-surface-50-900-token px-2 py-3 text-center" scope="col">
              <label class="label">
                <span><T defaultValue="Incoming" keyName="profiles.merge.incoming" /></span>
                <select
                  class="select text-sm font-normal w-full"
                  aria-label={$t('profiles.merge.incoming', 'Incoming')}
                  disabled={merging || !$canModify}
                  title={incomingProfile}
                  bind:value={incomingProfile}
                  on:change={loadComparison}>
                  {#each incomingProfiles as name}
                    <option value={name}>{name}</option>
                  {/each}
                </select>
              </label>
            </th>
          </tr>
        </thead>
        <tbody>
          {#if loading}
            <tr><td class="p-2" colspan="3"><p role="status"><T defaultValue="Loading profiles..." keyName="profiles.merge.loading-profiles" /></p></td></tr>
          {:else if !incomingProfile}
            <tr><td class="p-2" colspan="3"><T defaultValue="Create another profile to compare." keyName="profiles.merge.no-incoming" /></td></tr>
          {:else if !comparisonError && !entries.length}
            <tr><td class="p-2" colspan="3"><p role="status"><T defaultValue="These profiles contain no mods." keyName="profiles.merge.empty" /></p></td></tr>
          {:else}
            {#each entries as entry (entry.mod)}
              <tr>
                <td class="rounded-l-lg bg-surface-200-700-token p-2 align-middle">
                  <div class="flex h-[4.25rem] @md/mods-list:h-[5.5rem] space-x-2 min-w-0">
                    <ModPreview mod={previewFor(entry.mod)} />
                  </div>
                </td>
                <td class="bg-surface-200-700-token p-2 align-middle">
                  <div class="flex flex-col items-center gap-1 min-w-0">
                    <input
                      class="checkbox"
                      aria-label={$t('profiles.merge.current-enabled', 'Enabled in Current: {mod}', { mod: previewFor(entry.mod).name })}
                      disabled={merging || !$canModify}
                      type="checkbox"
                      value={entry.mod}
                      bind:group={enabledMods} />
                    {#if entry.current && entry.incoming && entry.current.version !== entry.incoming.version}
                      <select
                        class="select text-xs max-w-full"
                        aria-label={$t('profiles.merge.choose-version', 'Version in Current: {mod}', { mod: previewFor(entry.mod).name })}
                        disabled={merging || !$canModify}
                        bind:value={selectedVersions[entry.mod]}>
                        <option value={entry.current.version}>{entry.current.version}</option>
                        <option value={entry.incoming.version}>{entry.incoming.version}</option>
                      </select>
                    {:else if entry.current || enabledMods.includes(entry.mod)}
                      <code class="break-all text-center">{selectedVersions[entry.mod]}</code>
                    {:else}
                      <span><T defaultValue="Not present" keyName="profiles.merge.absent" /></span>
                    {/if}
                  </div>
                </td>
                <td class="rounded-r-lg bg-surface-200-700-token p-2 align-middle">
                  <div class="flex flex-col items-center gap-1 min-w-0">
                    <input
                      class="checkbox"
                      aria-label={$t('profiles.merge.incoming-enabled', 'Enabled in Incoming: {mod}', { mod: previewFor(entry.mod).name })}
                      checked={entry.incoming?.enabled ?? false}
                      disabled
                      type="checkbox" />
                    {#if entry.incoming}
                      <code class="break-all text-center">{entry.incoming.version}</code>
                    {:else}
                      <span><T defaultValue="Not present" keyName="profiles.merge.absent" /></span>
                    {/if}
                  </div>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
    {#if merging}
      <p role="status">{$progressMessage}</p>
      <ProgressBar
        class="h-4 w-full"
        max={1}
        meter="bg-primary-600"
        value={$progressPercent} />
    {/if}
  </section>
  <footer class="card-footer flex flex-wrap items-center gap-2">
    <button class="btn" disabled={merging} on:click={parent.onClose}>
      <T defaultValue="Close" keyName="common.close" />
    </button>
    <span class="grow text-sm">
      <T defaultValue={'{count} changes'} keyName="profiles.merge.changes" params={{ count: selections.length }} />
    </span>
    <button class="btn text-primary-600" disabled={loading || merging || !$canModify || !selections.length} on:click={finishMerge}>
      <T defaultValue="Save changes" keyName="profiles.merge.save-changes" />
    </button>
  </footer>
</div>
