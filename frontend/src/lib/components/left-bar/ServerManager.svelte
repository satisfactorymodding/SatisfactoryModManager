<script lang="ts">
  import Button, { Label } from '@smui/button';
  import { mdiServerNetwork, mdiTrashCan } from '@mdi/js';
  import Dialog, { Actions, Content, Title } from '@smui/dialog';
  import DataTable, { Body, Cell, Row } from '@smui/data-table';
  import Select, { Option } from '@smui/select';
  import Textfield from '@smui/textfield';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { remoteServers } from '$lib/store/ficsitCLIStore';
  import { AddRemoteServer, RemoveRemoteServer } from '$lib/generated/wailsjs/go/ficsitcli/FicsitCLI';
  import type { installfinders } from '$lib/generated/wailsjs/go/models';

  let dialogOpen = false;

  const allowedProtocols = ['ftp://'];

  async function removeServer(server: installfinders.Installation) {
    try {
      await RemoveRemoteServer(server.path);
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

  async function addNewRemoteServer() {
    if (!newServerPath) {
      return;
    }
    try {
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
    }
  }
</script>

<Button variant="unelevated" class="w-full mt-2" on:click={() => dialogOpen = true}>
  <Label>
    Manage Servers
  </Label>
  <div class="grow" />
  <SvgIcon icon={mdiServerNetwork} class="h-5 w-5" />
</Button>

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
            <Cell>{remoteServer.path}</Cell>
            <Cell>{remoteServer.type}</Cell>
            <Cell>{remoteServer.version}</Cell>
            <Cell>
              <Button on:click={() => removeServer(remoteServer)}>
                <SvgIcon icon={mdiTrashCan} class="!p-1 !m-0 !w-full !h-full group-hover:!hidden"/>
              </Button>
            </Cell>
          </Row>
        {/each}
      </Body>
    </DataTable>
    <div class="mt-4">
      <div class="flex h-10">
        <Select bind:value={newProtocol} class="!h-full w-32">
          {#each allowedProtocols as protocol}
            <Option value={protocol}>{protocol}</Option>
          {/each}
        </Select>
        <Textfield bind:value={newServerPath} class="!h-full grow mx-4"/>
        <Button on:click={() => addNewRemoteServer()} class="!h-full">Add</Button>
      </div>
      <p>{err}</p>
    </div>
  </Content>
  <Actions>
    <Button on:click={() => dialogOpen = false}>Close</Button>
  </Actions>
</Dialog>
