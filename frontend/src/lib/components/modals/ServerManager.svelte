<script lang="ts">
  import { mdiAlert, mdiLoading, mdiServerNetwork, mdiTrashCan } from '@mdi/js';
  import { type PopupSettings, popup } from '@skeletonlabs/skeleton';

  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import Select from '$lib/components/Select.svelte';
  import Tooltip from '$lib/components/Tooltip.svelte';
  import { AddRemoteServer, FetchRemoteServerMetadata, RemoveRemoteServer } from '$lib/generated/wailsjs/go/ficsitcli/ficsitCLI';
  import { ficsitcli } from '$lib/generated/wailsjs/go/models';
  import { installsMetadata, remoteServers } from '$lib/store/ficsitCLIStore';

  export let parent: { onClose: () => void };
  
  type RemoteType = ({ type: 'remote'; protocol: string; } | { type: 'local' }) & { name: string; };

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
      await FetchRemoteServerMetadata(server);
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


<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    Dedicated Servers
  </header>
  <section class="p-4 flex-auto space-y-4 overflow-y-auto flex">
    <div class="flex-auto w-full overflow-x-auto overflow-y-auto">
      <table class="table">
        <tbody>
          {#each $remoteServers as remoteServer}
            <tr>
              <td class="break-all">{remoteServer}</td>
              <td>
                {#if $installsMetadata[remoteServer]?.state === ficsitcli.InstallState.VALID}
                  {$installsMetadata[remoteServer].info?.type}
                {:else}
                  <Tooltip popupId={installWarningPopupId(remoteServer)}>
                    <span class="text-base">
                      {#if $installsMetadata[remoteServer]?.state === ficsitcli.InstallState.LOADING}
                        Loading...
                      {:else if $installsMetadata[remoteServer]?.state === ficsitcli.InstallState.INVALID}
                        SMM cannot manage this install
                      {:else}
                        Failed to connect to server, click to retry
                      {/if}
                    </span>
                  </Tooltip>
                  <div
                    class="h-6 w-full text-sm"
                    use:popup={installWarningPopups[remoteServer]}>
                    {#if $installsMetadata[remoteServer]?.state === ficsitcli.InstallState.LOADING}
                      <SvgIcon
                        class="!p-0 !m-0 !w-full !h-full animate-spin text-primary-600"
                        icon={mdiLoading} />
                    {:else if $installsMetadata[remoteServer]?.state === ficsitcli.InstallState.INVALID}
                      <SvgIcon
                        class="!p-0 !m-0 !w-full !h-full text-red-500"
                        icon={mdiAlert} />
                    {:else}
                      <button
                        class="btn-icon h-6 w-full"
                        on:click={() => retryConnect(remoteServer)}>
                        <SvgIcon
                          class="!p-0 !m-0 !w-full !h-full text-red-500"
                          icon={mdiAlert} />
                      </button>
                    {/if}
                  </div>
                {/if}
              </td>
              <td>
                {#if $installsMetadata[remoteServer]?.info?.version}
                  CL{$installsMetadata[remoteServer].info?.version}
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
    </div>
  </section>
  <section class="p-4 space-y-4">
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 items-start auto-rows-fr">
      <Select
        name="newServerProtocol"
        class="col-span-1 h-full"
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
            class="input px-4 h-full sm:col-start-2 col-span-2"
            placeholder="user:pass@host:port/path"
            type="text"
            bind:value={newServerPath}/>
          <p class="sm:col-start-2 col-span-2">
            Note that you might have to escape certain characters in the username and password
          </p>
        {:else}
          <input
            class="input px-4 h-full sm:col-start-2"
            placeholder="user"
            type="text"
            bind:value={newServerUsername}/>
          <input
            class="input px-4 h-full"
            placeholder="pass"
            type="text"
            bind:value={newServerPassword}/>
          <input
            class="input px-4 h-full sm:col-start-2"
            placeholder="host"
            type="text"
            bind:value={newServerHost}/>
          <input
            class="input px-4 h-full"
            placeholder="port"
            type="text"
            bind:value={newServerPort}/>
          <input
            class="input px-4 h-full sm:col-start-2 col-span-2"
            placeholder="path"
            type="text"
            bind:value={newServerPath}/>
          <p class="sm:col-start-2 col-span-2">
            Complete path: {fullInstallPath}
          </p>
        {/if}
        <button class="btn sm:col-start-2 col-span-2 text-sm whitespace-break-spaces bg-surface-200-700-token" on:click={() => advancedMode = !advancedMode}>
          {#if advancedMode}
            Switch to simple mode
          {:else}
            Switch to advanced mode
          {/if}
        </button>
      {:else}
        <input
          class="input px-4 h-full sm:col-start-2 col-span-2"
          placeholder="C:\Path\To\Server"
          type="text"
          bind:value={newServerPath}/>
      {/if}
      <button
        class="btn h-full text-sm bg-primary-600 text-secondary-900 col-start-2 sm:col-start-4 row-start-1"
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
  </section>
  <footer class="card-footer">
    <button
      class="btn h-8 w-full text-sm bg-surface-200-700-token"
      on:click={parent.onClose}>
      Close
    </button>
  </footer>
</div>
