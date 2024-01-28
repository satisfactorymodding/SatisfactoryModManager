<script lang="ts">
  import { mdiAlert, mdiServerNetwork, mdiTrashCan } from '@mdi/js';
  import { popup, type PopupSettings } from '@skeletonlabs/skeleton';

  import Tooltip from '$lib/components/Tooltip.svelte';
  import Select from '$lib/components/Select.svelte';
  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { installsMetadata, remoteServers } from '$lib/store/ficsitCLIStore';
  import { AddRemoteServer, RemoveRemoteServer } from '$lib/generated/wailsjs/go/ficsitcli/FicsitCLI';

  export let parent: { onClose: () => void };
  
  type RemoteType = ({ type: 'remote'; protocol: string; } | { type: 'local' }) & { name: string; }

  const remoteTypes: RemoteType[] = [
    { type: 'remote', protocol: 'ftp://', name: 'FTP' },
    { type: 'remote', protocol: 'sftp://', name: 'SFTP' },
    { type: 'local', name: 'Path' },
  ];

  let newRemoteType = remoteTypes[0];

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

  let newServerUsername = '';
  let newServerPassword = '';
  let newServerHost = '';
  let newServerPort = '';
  let newServerPath = '';
  let err = '';

  let advancedMode = false;

  let addInProgress = false;

  $: fullInstallPath = (() => {
    if (newRemoteType.type === 'local') {
      return newServerPath;
    }
    if (advancedMode) {
      return newRemoteType.protocol + newServerPath;
    }
    return newRemoteType.protocol + encodeURIComponent(newServerUsername) + ':' + encodeURIComponent(newServerPassword) + '@' + newServerHost + ':' + newServerPort + '/' + newServerPath;
  })();

  async function addNewRemoteServer() {
    if (!newServerPath) {
      return;
    }
    try {
      err = '';
      addInProgress = true;
      await AddRemoteServer(fullInstallPath);
      newServerUsername = '';
      newServerPassword = '';
      newServerHost = '';
      newServerPort = '';
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


<div style="max-height: calc(100vh - 128px); max-width: calc(100vw - 128px);" class="card flex flex-col gap-2 !min-w-[800px] min-h-[400px]">
  <header class="card-header font-bold text-2xl text-center">
    Dedicated Servers
  </header>
  <section class="p-4 flex-auto space-y-4 overflow-y-auto mr-4">
    <table class="table w-full">
      <tbody>
        {#each $remoteServers as remoteServer}
          <tr>
            <td class="break-all">{remoteServer}</td>
            <td>
              {#if $installsMetadata[remoteServer]?.type}
                {$installsMetadata[remoteServer].type}
              {:else}
                <button
                  class="btn-icon h-6 w-full text-sm"
                  on:click={() => retryConnect(remoteServer)}
                  use:popup={installWarningPopups[remoteServer]}>
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
              {#if $installsMetadata[remoteServer]?.version}
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
  </section>
  <section class="p-4 space-y-4">
    <div>
      <div class="flex h-10 items-baseline space-x-4">
        <Select
          name="newServerProtocol"
          class="!h-full w-28 shrink-0"
          buttonClass="bg-surface-200-700-token px-4 text-sm"
          itemActiveClass="!bg-surface-300/20"
          itemClass="bg-surface-50-900-token"
          itemKey="name"
          items={remoteTypes}
          bind:value={newRemoteType}>
          <svelte:fragment slot="item" let:item>
            {item.name}
          </svelte:fragment>
        </Select>
        {#if newRemoteType.type === 'remote'}
          {#if advancedMode}
            <input
              class="input !h-full grow mx-4 px-4"
              placeholder="user:pass@host:port/path"
              type="text"
              bind:value={newServerPath}/>
          {:else}
            <div class="flex h-10 items-baseline space-x-1 [&>input]:px-4">
              <input
                class="input !h-full grow"
                placeholder="user"
                type="text"
                bind:value={newServerUsername}/>
              <span>:</span>
              <input
                class="input !h-full grow mx-4 px-4"
                placeholder="pass"
                type="text"
                bind:value={newServerPassword}/>
              <span>@</span>
              <input
                class="input !h-full grow mx-4 px-4"
                placeholder="host"
                type="text"
                bind:value={newServerHost}/>
              <span>:</span>
              <input
                class="input !h-full grow mx-4 px-4"
                placeholder="port"
                type="text"
                bind:value={newServerPort}/>
              <span>/</span>
              <input
                class="input !h-full grow mx-4 px-4"
                placeholder="path"
                type="text"
                bind:value={newServerPath}/>
            </div>
          {/if}
        {:else}
          <input
            class="input !h-full grow mx-4 px-4"
            placeholder="C:\Path\To\Server"
            type="text"
            bind:value={newServerPath}/>
        {/if}
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
    </div>
    {#if newRemoteType.type === 'remote'}
      <div class="flex items-baseline space-x-4">
        <button class="btn h-8 text-sm bg-surface-200-700-token" on:click={() => advancedMode = !advancedMode}>
          {#if advancedMode}
            Switch to simple mode
          {:else}
            Switch to advanced mode
          {/if}
        </button>
        {#if newRemoteType.type === 'remote' && !advancedMode}
          <p>Complete path: {fullInstallPath}</p>
        {:else}
          <p>Note that you might have to escape certain characters in the username and password</p>
        {/if}
      </div>
    {/if}
    <p>{err}</p>
  </section>
  <footer class="card-footer">
    <button
      class="btn h-8 w-full text-sm bg-surface-200-700-token"
      on:click={parent.onClose}>
      Close
    </button>
  </footer>
</div>
