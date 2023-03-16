<script lang="ts">
  import Button, { Label } from '@smui/button';
  import { mdiCheckCircle, mdiSync, mdiUpload } from '@mdi/js';
  import Dialog, { Actions, Content, Title } from '@smui/dialog';
  import List, { Item, PrimaryText, SecondaryText, Text } from '@smui/list';
  import { getContextClient, queryStore } from '@urql/svelte';

  import UpdateChangelog from './UpdateChangelog.svelte';
  import SMMUpdateDialog from './SMMUpdateDialog.svelte';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { checkForUpdates, canModify, progress, updates, updateCheckInProgress } from '$lib/store/ficsitCLIStore';
  import { error } from '$lib/store/generalStore';
  import { OfflineGetModsByReferences, UpdateAllMods } from '$wailsjs/go/ficsitcli/FicsitCLI';
  import type { ficsitcli } from '$wailsjs/go/models';
  import { CheckForUpdates as CheckForSMMUpdates } from '$wailsjs/go/bindings/Update';
  import { smmUpdate, smmUpdateReady } from '$lib/store/smmUpdateStore';
  import { GetModNamesDocument } from '$lib/generated';
  import { offline } from '$lib/store/settingsStore';

  const client = getContextClient();

  $: modNamesQuery = queryStore({
    query: GetModNamesDocument,
    client,
    pause: !!$offline,
    variables: {
      modReferences: $updates.map((u) => u.item).filter((u) => u !== 'SML') as string[],
    },
  });

  let modNamesQueryResult: { mod_reference: string; name: string; }[] | undefined;
  
  $: modNames = modNamesQueryResult?.reduce((acc, mod) => {
    if(mod) {
      acc[mod.mod_reference] = mod.name;
    }
    return acc;
  }, {} as Record<string, string>) ?? {};

  $: {
    if($offline) {
      OfflineGetModsByReferences($updates.map((u) => u.item).filter((u) => u !== 'SML') as string[]).then((mods) => { modNamesQueryResult = mods; });
    } else {
      modNamesQueryResult = $modNamesQuery.data?.getMods?.mods;
    }
  }

  let updatesDialog = false;

  async function updateAll() {
    if($updates.length > 0) {
      try {
        await UpdateAllMods();
        $updates = [];
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
  }

  let selectedUpdates: ficsitcli.Update[] = [];

  async function updateSelected() {
    // TODO
    console.log(selectedUpdates);
  }

  function toggleSelected(update: ficsitcli.Update) {
    if(selectedUpdates.includes(update)) {
      selectedUpdates = selectedUpdates.filter((u) => u !== update);
    } else {
      selectedUpdates = [...selectedUpdates, update];
    }
  }
  
  $: () => {
    $updates;
    selectedUpdates = [];
  };

  let changelogUpdate: ficsitcli.Update | null = null;

  function checkForAllUpdates() {
    checkForUpdates();
    if(!$smmUpdate || $smmUpdateReady) {
      CheckForSMMUpdates();
    }
  }

  let smmUpdateDialog: SMMUpdateDialog;

  function showUpdateDialog() {
    if($smmUpdate) {
      smmUpdateDialog.show();
    } else if ($updates.length > 0) {
      updatesDialog = true;
    }
  }
</script>

<Button variant="unelevated" class="w-full mt-2 update-button {$smmUpdate || $updates.length > 0 ? 'has-update' : ''}" on:click={() => showUpdateDialog()}>
  <Label>
    {#if $smmUpdate}
      SMM update available
    {:else if $updates.length > 0}
      {$updates.length} updates available
    {:else}
      No updates right now
    {/if}
  </Label>
  <div class="grow" />
  <SvgIcon icon={mdiCheckCircle} class="h-5 w-5" />
</Button>
<Button variant="unelevated" class="w-full mt-2" on:click={checkForAllUpdates} disabled={!!$progress || $updateCheckInProgress}>
  <Label>
    Check for updates
  </Label>
  <div class="grow" />
  <SvgIcon icon={mdiSync} class="h-5 w-5 {$updateCheckInProgress ? 'update-check' : ''}" />
</Button>

<Dialog
  bind:open={updatesDialog}
  surface$style="width: 500px; max-width: calc(100vw - 32px);"
  class="updates-dialog"
>
  <Title>Updates</Title>
  <Content>
    <List>
      {#each $updates as update}
        <Item 
          on:SMUI:action={() => toggleSelected(update)}
        >
          {#if selectedUpdates.includes(update)}
            <SvgIcon icon={mdiUpload} class="h-5 w-5" />
          {:else}
            <div class="w-5"/>
          {/if}
          <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
            <PrimaryText class="text-base">{ modNames[update.item] ?? update.item }</PrimaryText>
            <SecondaryText>{ update.currentVersion } -> { update.newVersion }</SecondaryText>
          </Text>
          <div class="grow" />
          <div on:click|stopPropagation={() => {/* empty */}}>
            <Button on:click={() => changelogUpdate = update}>Changelog</Button>
          </div>
        </Item>
      {/each}
    </List>
  </Content>
  <Actions>
    <Button on:click={() => updateAll()} disabled={!$canModify || $updateCheckInProgress}>
      <Label>Update All</Label>
    </Button>
    <Button on:click={() => updateSelected()} disabled={!$canModify || $updateCheckInProgress}>
      <Label>Update Selected</Label>
    </Button>
  </Actions>
</Dialog>

<UpdateChangelog bind:update={changelogUpdate} />

<SMMUpdateDialog bind:this={smmUpdateDialog} />