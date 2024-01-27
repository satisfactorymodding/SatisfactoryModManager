<script lang="ts">
  import { mdiAlert, mdiServerNetwork, mdiTrashCan } from '@mdi/js';
  import { popup, type PopupSettings } from '@skeletonlabs/skeleton';

  import Tooltip from '$lib/components/Tooltip.svelte';
  import Select from '$lib/components/Select.svelte';
  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { installsMetadata, remoteServers } from '$lib/store/ficsitCLIStore';
  import { AddRemoteServer, RemoveRemoteServer } from '$lib/generated/wailsjs/go/ficsitcli/FicsitCLI';

  export let parent: { onClose: () => void };
  
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

  function installWarningPopupId(install: string) {
    return `remote-server-warning-${install}`;
  }

  $: installWarningPopups = $remoteServers.map((i) => [i, {
    event: 'hover',
    target: installWarningPopupId(i),
    middleware: {
      offset: 4,
    },
    placement: 'bottom',
  } as PopupSettings]).reduce((acc, [k, v]) => ({ ...acc, [k as string]: v as PopupSettings }), {} as Record<string, PopupSettings>);
</script>


<div class="card flex flex-col gap-2 !min-w-[800px] min-h-[400px]" style="max-height: calc(100vh - 128px); max-width: calc(100vw - 128px);">
  <header class="card-header font-bold text-2xl text-center">
    Dedicated Servers
  </header>
  <section class="p-4 grow">
    <table class="table w-full">
      <tbody>
        {#each $remoteServers as remoteServer}
          <tr>
            <td class="break-all">{$installsMetadata[remoteServer].path}</td>
            <td>
              {#if $installsMetadata[remoteServer].type}
                {$installsMetadata[remoteServer].type}
              {:else}
                <button
                  use:popup={installWarningPopups[remoteServer]}
                  class="btn-icon h-6 w-full text-sm"
                  on:click={() => retryConnect(remoteServer)}>
                  <SvgIcon
                    class="!p-0 !m-0 !w-full !h-full text-red-500"
                    icon={mdiAlert} />
                </button>
                <Tooltip popupId={installWarningPopupId(remoteServer)}>
                  <span class="text-base">
                    Failed to connect to server, click to retry
                  </span>
                </Tooltip>
              {/if}
            </td>
            <td>
              {#if $installsMetadata[remoteServer].version}
                CL{$installsMetadata[remoteServer].version}
              {/if}
            </td>
            <td>
              <button
                class="btn-icon h-5 w-full"
                on:click={() => removeServer(remoteServer)}>
                <SvgIcon
                  class="!p-0 !m-0 !w-full !h-full hover:text-red-500"
                  icon={mdiTrashCan}/>
              </button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
    <div class="mt-4">
      <div class="flex h-10">
        <Select bind:value={newProtocol} items={allowedProtocols} name="newServerProtocol" class="!h-full w-28 shrink-0"
          buttonClass="bg-surface-200-700-token px-4 text-sm"
          itemClass="bg-surface-50-900-token"
          itemActiveClass="!bg-surface-300/20">
          <svelte:fragment slot="item" let:item>
            {item}
          </svelte:fragment>
        </Select>
        <input type="text" bind:value={newServerPath} class="input !h-full grow mx-4 px-4" placeholder="user:pass@host:port/path"/>
        <button
          class="btn h-full text-sm bg-primary-600 text-secondary-900"
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
  </section>
  <footer class="card-footer">
    <button
      class="btn h-8 w-full text-sm bg-surface-200-700-token"
      on:click={parent.onClose}>
      Close
    </button>
  </footer>
</div>
