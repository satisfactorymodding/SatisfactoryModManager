<script lang="ts">
  import Button, { Label } from '@smui/button';
  import MdiIcon from '$lib/components/MDIIcon.svelte';
  import { mdiCheckCircle, mdiSync, mdiUpload } from '@mdi/js';

  import { checkForUpdates, progress, updates, updateCheckInProgress } from '$lib/store/ficsitCLIStore';
  import Dialog, { Actions, Content, Title } from '@smui/dialog';
  import { UpdateAllMods } from '$wailsjs/go/bindings/FicsitCLI';
  import List, { Item, PrimaryText, SecondaryText, Text } from '@smui/list';
  import type { bindings } from '$wailsjs/go/models';

  $: canInstall = !$progress;

  let updatesDialog = false;

  async function updateAll() {
    if($updates.length > 0) {
      const err = await UpdateAllMods();
      if(!err) {
        $updates = [];
      }
    }
  }

  let selectedUpdates: bindings.Update[] = [];

  async function updateSelected() {
    // TODO
    console.log(selectedUpdates);
  }

  function toggleSelected(update: bindings.Update) {
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
</script>

<Button variant="unelevated" class="w-full mt-2 update-button {$updates.length > 0 ? 'has-update' : ''}" on:click={() => updatesDialog = true} disabled={!canInstall}>
  <Label>
    {#if $updates.length === 0}
      No updates right now
    {:else}
      {$updates.length} updates available
    {/if}
  </Label>
  <div class="grow" />
  <MdiIcon icon={mdiCheckCircle} class="h-5" />
</Button>
<Button variant="unelevated" class="w-full mt-2" on:click={checkForUpdates} disabled={!canInstall}>
  <Label>
    Check for updates
  </Label>
  <div class="grow" />
  <MdiIcon icon={mdiSync} class="h-5 {$updateCheckInProgress ? 'update-check' : ''}" />
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
            <MdiIcon icon={mdiUpload} class="h-5" />
          {:else}
            <div class="w-7"/>
          {/if}
          <Text class="pl-2 h-full flex flex-col content-center mb-1.5">
            <PrimaryText class="text-base">{ update.item }</PrimaryText>
            <SecondaryText>{ update.currentVersion } -> { update.newVersion }</SecondaryText>
          </Text>
        </Item>
      {/each}
    </List>
  </Content>
  <Actions>
    <Button on:click={() => updateAll()}>
      <Label>Update All</Label>
    </Button>
    <Button on:click={() => updateSelected()}>
      <Label>Update Selected</Label>
    </Button>
  </Actions>
</Dialog>