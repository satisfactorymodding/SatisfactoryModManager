<script lang="ts">
  import { mdiAlert, mdiLoading, mdiServerNetwork, mdiTrashCan } from '@mdi/js';
  import { getTranslate } from '@tolgee/svelte';
  import _ from 'lodash';

  import RemoteServerPicker from '$lib/components/RemoteServerPicker.svelte';
  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import Select from '$lib/components/Select.svelte';
  import T from '$lib/components/T.svelte';
  import Tooltip from '$lib/components/Tooltip.svelte';
  import { type PopupSettings, popup } from '$lib/skeletonExtensions';
  import { installsMetadata, remoteServers } from '$lib/store/ficsitCLIStore';
  import { AddRemoteServer, FetchRemoteServerMetadata, GetNextRemoteLauncherName, RemoveRemoteServer } from '$wailsjs/go/ficsitcli/ficsitCLI';
  import { ficsitcli } from '$wailsjs/go/models';

  export let parent: { onClose: () => void };
  
  const { t } = getTranslate();

  type RemoteType = ({ type: 'remote'; protocol: string; defaultPort: string; } | { type: 'local' }) & { name: string; };

  const remoteTypes: RemoteType[] = [
    { type: 'remote', protocol: 'ftp://', name: 'FTP', defaultPort: '21' },
    { type: 'remote', protocol: 'sftp://', name: 'SFTP', defaultPort: '22' },
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

  let defaultRemoteName = '';
  let remoteName = '';

  let advancedMode = false;

  let addInProgress = false;
  let maskPassword = true;

  $: authString = encodeURIComponent(newServerUsername) + (newServerPassword ? ':' + encodeURIComponent(newServerPassword) : '');
  $: actualPort = newRemoteType.type === 'remote' ? (newServerPort.length > 0 ? newServerPort : newRemoteType.defaultPort) : '';

  $: trimmedPath = _.trimStart(newServerPath, '/');

  $: fullInstallPath = (() => {
    if (newRemoteType.type === 'local') {
      return newServerPath;
    }
    if (advancedMode) {
      return newRemoteType.protocol + newServerPath;
    }
    return newRemoteType.protocol + authString + '@' + newServerHost + ':' + actualPort + '/' + trimmedPath;
  })();

  $: baseServerPath = (() => {
    if (newRemoteType.type === 'local') {
      return newServerPath;
    }
    if (advancedMode) {
      return newRemoteType.protocol + newServerPath;
    }
    return newRemoteType.protocol + authString + '@' + newServerHost + ':' + actualPort;
  })();

  $: isBaseValid = (() => {
    if (newRemoteType.type === 'local') {
      return newServerPath.length > 0;
    }
    if (advancedMode) {
      return newServerPath.length > 0;
    }
    return newServerUsername.length > 0 && newServerHost.length > 0;
  })();

  let isPathValid = false;

  $: isValid = (() => {
    if (newRemoteType.type === 'local') {
      return newServerPath.length > 0;
    }
    if (advancedMode) {
      return newServerPath.length > 0;
    }
    return newServerUsername.length > 0 && newServerHost.length > 0 && isPathValid;
  })();

  async function addNewRemoteServer() {
    if (!isValid) {
      return;
    }
    try {
      err = '';
      addInProgress = true;
      await AddRemoteServer(fullInstallPath, remoteName);
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

  $: {
    remoteServers;
    (async () => {
      defaultRemoteName = await GetNextRemoteLauncherName();
    })();
  }

  $: installWarningPopups = $remoteServers.map((i) => [i, {
    event: 'hover',
    target: installWarningPopupId(i),
    middleware: {
      offset: 4,
    },
    placement: 'bottom',
  } as PopupSettings]).reduce((acc, [k, v]) => ({ ...acc, [k as string]: v as PopupSettings }), {} as Record<string, PopupSettings>);

  function toggleMaskPassword() {
    maskPassword = !maskPassword;
  }

  function redactRemoteURL(path: string) {
    return path.replace(/(?<=.+:\/\/)(?:(.+?)(?::.*?)?)?(?=@)/, '$1:********');
  }
</script>


<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    <T defaultValue="Dedicated Servers" keyName="server-manager.title" />
  </header>
  <section class="p-4 flex-auto space-y-4 overflow-y-auto flex">
    <div class="flex-auto w-full overflow-x-auto overflow-y-auto">
      <table class="table">
        <tbody>
          {#each $remoteServers as remoteServer}
            <tr>
              <td class="break-all">{$installsMetadata[remoteServer].info?.launcher}</td>
              <td class="break-all">{redactRemoteURL(remoteServer)}</td>
              <td>
                {#if $installsMetadata[remoteServer]?.state === ficsitcli.InstallState.VALID}
                  {$installsMetadata[remoteServer].info?.type}
                {:else}
                  <Tooltip popupId={installWarningPopupId(remoteServer)}>
                    <span class="text-base">
                      {#if $installsMetadata[remoteServer]?.state === ficsitcli.InstallState.LOADING}
                        <T defaultValue="Loading..." keyName="server-manager.loading" />
                      {:else if $installsMetadata[remoteServer]?.state === ficsitcli.InstallState.INVALID}
                        <T defaultValue="SMM cannot manage this install" keyName="server-manager.invalid" />
                      {:else}
                        <T defaultValue="Failed to connect to server, click to retry" keyName="server-manager.failed-to-connect" />
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
  <section class="p-4 space-y-4 overflow-y-auto">
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 items-start auto-rows-[minmax(2.5rem,_max-content)]">
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
            placeholder={$t('server-manager.advanced-path-placeholder', 'user:pass@host:port/path')}
            type="text"
            bind:value={newServerPath}/>
          <p class="sm:col-start-2 col-span-2">
            <T defaultValue="Note that you might have to escape certain characters in the username and password" keyName="server-manager.advanced-note" />
          </p>
        {:else}
          <input
            class="input px-4 h-full sm:col-start-2"
            placeholder={$t('server-manager.username-placeholder', 'username')}
            type="text"
            bind:value={newServerUsername}/>
          <!-- This is a conditional because the type var cant be dynamic with bind:value -->
          {#if maskPassword}
            <input
              class="input px-4 h-full"
              placeholder={$t('server-manager.password-placeholder', 'password')}
              type="password"
              bind:value={newServerPassword}/>
          {:else}
            <input
              class="input px-4 h-full"
              placeholder={$t('server-manager.password-placeholder', 'password')}
              type="text"
              bind:value={newServerPassword}/>
          {/if}
          <input
            class="input px-4 h-full sm:col-start-2"
            placeholder={$t('server-manager.host-placeholder', 'host')}
            type="text"
            bind:value={newServerHost}/>
          <input
            class="input px-4 h-full"
            placeholder={$t('server-manager.port-placeholder', 'port (default: {default})', { default: newRemoteType.defaultPort })}
            type="text"
            bind:value={newServerPort}/>
          <input
            class="input px-4 h-full sm:col-start-2 col-span-2"
            placeholder={$t('server-manager.path-placeholder', 'path')}
            type="text"
            bind:value={newServerPath}/>
          <div class="sm:col-start-2 col-span-2">
            <RemoteServerPicker
              basePath={baseServerPath}
              disabled={!isBaseValid}
              bind:path={newServerPath}
              bind:valid={isPathValid}
            />
          </div>
        {/if}
        <button class="btn sm:col-start-1 col-span-1 row-start-2 text-sm whitespace-break-spaces bg-surface-200-700-token" on:click={() => advancedMode = !advancedMode}>
          {#if advancedMode}
            <T defaultValue="Switch to simple mode" keyName="server-manager.switch-to-simple" />
          {:else}
            <T defaultValue="Switch to advanced mode" keyName="server-manager.switch-to-advanced" />
          {/if}
        </button>
      {:else}
        <input
          class="input px-4 h-full sm:col-start-2 col-span-2"
          placeholder={$t('server-manager.local-path-placeholder', 'C:\\Path\\To\\Server')}
          type="text"
          bind:value={newServerPath}/>
        <div class="sm:col-start-2 col-span-2">
          <RemoteServerPicker
            basePath=""
            bind:path={newServerPath}
            bind:valid={isPathValid}
          />
        </div>
      {/if}
      <input
        class="input px-4 h-full col-start-4 row-start-1"
        placeholder={$t('server-manager.name-placeholder', 'Name (default: {default})', { default: defaultRemoteName })}
        type="text"
        bind:value={remoteName}/>
      <button
        class="btn h-full text-sm bg-primary-600 text-secondary-900 col-start-2 sm:col-start-4 row-start-2"
        disabled={addInProgress || !isValid}
        on:click={() => addNewRemoteServer()}>
        <span>
          {#if !addInProgress}
            <T defaultValue="Add" keyName="server-manager.add" />
          {:else}
            <T defaultValue="Validating..." keyName="server-manager.validating" />
          {/if}
        </span>
        <div class="grow" />
        <SvgIcon
          class="h-5 w-5"
          icon={mdiServerNetwork} />
      </button>
      <button
        class="btn h-full text-sm bg-primary-600 text-secondary-900 col-start-2 sm:col-start-4 row-start-3"
        on:click={() => toggleMaskPassword()}>
        <span>
          {#if maskPassword}
            <T defaultValue="Show Password" keyName="server-manager.show-password" />
          {:else}
            <T defaultValue="Hide Password" keyName="server-manager.hide-password" />
          {/if}
        </span>
      </button>
    </div>
    <p>{err}</p>
  </section>
  <footer class="card-footer">
    <button
      class="btn h-8 w-full text-sm bg-surface-200-700-token"
      on:click={parent.onClose}>
      <T defaultValue="Close" keyName="common.close" />
    </button>
  </footer>
</div>
