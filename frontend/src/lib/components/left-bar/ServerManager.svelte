<script lang="ts">
  import { mdiAlert, mdiServerNetwork, mdiTrashCan } from '@mdi/js';
  import Dialog, { Actions, Content, Title } from '@smui/dialog';
  import DataTable, { Body, Cell, Row } from '@smui/data-table';
  import Tooltip, { Wrapper } from '@smui/tooltip';
  import Textfield from '@smui/textfield';

  import Select from '$lib/components/Select.svelte';
  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { installsMetadata, remoteServers } from '$lib/store/ficsitCLIStore';
  import { AddRemoteServer, RemoveRemoteServer } from '$lib/generated/wailsjs/go/ficsitcli/FicsitCLI';

  let dialogOpen = false;

  const allowedProtocols = ['ftp://', 'sftp://'];

  async function removeServer(server: string) {
    try {
      await RemoveRemoteServer(server);
    } catch (e) {
      if(e instanceof Error) {
        err = e.message;
      } else if (typeof e === 'string') {
        err = e;
      } else {
        err = 'Unknown error';
      }
    }
  }

  let newProtocol = allowedProtocols[0];
  let newServerPath = '';
  let err = '';

  let addInProgress = false;

  async function addNewRemoteServer() {
    if (!newServerPath) {
      return;
    }
    try {
      addInProgress = true;
      await AddRemoteServer(newProtocol + newServerPath);
      newServerPath = '';
    } catch (e) {
      if(e instanceof Error) {
        err = e.message;
      } else if (typeof e === 'string') {
        err = e;
      } else {
        err = 'Unknown error';
      }
    } finally {
      addInProgress = false;
    }
  }

  async function retryConnect(server: string) {
    try {
      await AddRemoteServer(server);
    } catch (e) {
      if(e instanceof Error) {
        err = e.message;
      } else if (typeof e === 'string') {
        err = e;
      } else {
        err = 'Unknown error';
      }
    }
  }
</script>

<button
  class="btn px-4 h-8 w-full text-sm bg-surface-200-700-token"
  on:click={() => dialogOpen = true}>
  <span>Manage Servers</span>
  <div class="grow" />
  <SvgIcon
    class="h-5 w-5"
    icon={mdiServerNetwork} />
</button>

<Dialog
  bind:open={dialogOpen} scrimClickAction="" escapeKeyAction=""
  surface$style="max-height: calc(100vh - 128px); max-width: calc(100vw - 128px);"
  surface$class="!min-w-[800px] min-h-[400px]"
>
  <Title>Dedicated Servers</Title>
  <Content>
    <DataTable table$aria-label="Todo list" style="width: 100%;">
      <Body>
        {#each $remoteServers as remoteServer}
          <Row>
            <Cell>{$installsMetadata[remoteServer].path}</Cell>
            <Cell>
              {#if $installsMetadata[remoteServer].type}
                {$installsMetadata[remoteServer].type}
              {:else}
                <Wrapper>
                  <button
                    class="btn-icon h-6 w-full text-sm"
                    on:click={() => retryConnect(remoteServer)}>
                    <SvgIcon
                      class="!p-0 !m-0 !w-full !h-full text-red-500"
                      icon={mdiAlert} />
                  </button>
                  <Tooltip surface$class="max-w-lg text-base">
                    Failed to connect to server, click to retry
                  </Tooltip>
                </Wrapper>
              {/if}
            </Cell>
            <Cell>
              {#if $installsMetadata[remoteServer].version}
                {$installsMetadata[remoteServer].version}
              {/if}
            </Cell>
            <Cell>
              <button
                class="btn-icon h-6 w-full text-sm"
                on:click={() => removeServer(remoteServer)}>
                <SvgIcon
                  class="!p-0 !m-0 !w-full !h-full hover:text-red-500"
                  icon={mdiTrashCan}/>
              </button>
            </Cell>
          </Row>
        {/each}
      </Body>
    </DataTable>
    <div class="mt-4">
      <div class="flex h-10">
        <Select bind:value={newProtocol} items={allowedProtocols} name="newServerProtocol" class="!h-full w-28 shrink-0">
          <svelte:fragment slot="item" let:item>
            {item}
          </svelte:fragment>
        </Select>
        <Textfield bind:value={newServerPath} class="!h-full grow mx-4"/>
        <button
          class="btn h-full text-sm bg-surface-500-400-token"
          disabled={addInProgress}
          on:click={() => addNewRemoteServer()}>
          <span>
            {#if !addInProgress}
              Add
            {:else}
              Validating...
            {/if}
          </span>
          <div class="grow" />
          <SvgIcon
            class="h-5 w-5"
            icon={mdiServerNetwork} />
        </button>
      </div>
      <p>{err}</p>
    </div>
  </Content>
  <Actions>
    <button
      class="btn h-8 w-full text-sm bg-surface-200-700-token"
      on:click={() => dialogOpen = false}>
      Close
    </button>
  </Actions>
</Dialog>
