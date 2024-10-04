import { Client } from '@urql/svelte';
import { coerce, satisfies } from 'semver';
import { get } from 'svelte/store';

import {
  type Compatibility,
  CompatibilityState,
  GetModVersionTargetsDocument,
  ModReportedCompatibilityDocument,
  ModVersionsCompatibilityDocument,
  TargetName, type Version, type VersionTarget,
} from '$lib/generated';
import {
  installsMetadata,
  selectedInstallMetadata,
  selectedProfile,
  selectedProfileTargets,
} from '$lib/store/ficsitCLIStore';
import { offline } from '$lib/store/settingsStore';
import { OfflineGetMod } from '$wailsjs/go/ficsitcli/ficsitCLI';
import { common } from '$wailsjs/go/models';

export interface CompatibilityWithSource extends Compatibility {
  source: 'reported' | 'version';
}

const clientTargets = [TargetName.Windows];
const serverTargets = [TargetName.LinuxServer, TargetName.WindowsServer];

export async function getCompatibility(modReference: string, urqlClient: Client): Promise<CompatibilityWithSource> {
  const installInfo = get(selectedInstallMetadata).info;
  if(!installInfo) {
    return { state: CompatibilityState.Broken, note: 'No game selected.', source: 'version' };
  }

  const targetCompatibility = await getTargetCompatibilityFor(modReference, urqlClient);
  if (targetCompatibility.state !== CompatibilityState.Works) {
    return targetCompatibility;
  }

  return await getCompatibilityFor(modReference, installInfo, urqlClient);
}

function friendlyInstallName(install: common.Installation) {
  const installType = install.type === common.InstallType.WINDOWS ? 'Game' : 'Server';
  return `${install.launcher} (${installType})`;
}

async function getTargetCompatibilityFor(modReference: string, urqlClient: Client): Promise<CompatibilityWithSource> {
  const result = await urqlClient.query(GetModVersionTargetsDocument, { modReference }).toPromise();
  if(!result.data?.mod) {
    return { state: CompatibilityState.Broken, note: 'Mod not found.', source: 'version' };
  }

  if (result.data.mod.versions.length === 0) {
    return { state: CompatibilityState.Broken, note: 'Mod has no versions available.', source: 'version' };
  }

  const requiredTargets = Object.entries(get(selectedProfileTargets));
  if (requiredTargets.length === 0) {
    return { state: CompatibilityState.Works, source: 'version' };
  }

  const clientRequiredTargets = requiredTargets.filter(([target]) => clientTargets.includes(target as TargetName)).map(([target]) => target as TargetName);
  const serverRequiredTargets = requiredTargets.filter(([target]) => serverTargets.includes(target as TargetName)).map(([target]) => target as TargetName);

  const versionCompatibility = result.data.mod.versions.map((version) => ({ version, missingTargets: checkTargetCompatibility(version, clientRequiredTargets, serverRequiredTargets) }));
  if(versionCompatibility.some((comp) => comp.missingTargets.length === 0)) {
    return { state: CompatibilityState.Works, source: 'version' };
  }

  // The versions should be sorted by newest first
  const latestVersionCompatibility = versionCompatibility[0];

  const missingTargets = requiredTargets.filter(([target]) => latestVersionCompatibility.missingTargets.includes(target as TargetName));

  const requiredTargetsWithInstall = missingTargets
    .map(([target, installs]) =>
      `${target} required by: ${installs.map((install) => friendlyInstallName(get(installsMetadata)[install].info!)).join(',')}`)
    .join('\n\n');
  return { state: CompatibilityState.Broken, note: `This mod does not support one or more of your installs using profile \`${get(selectedProfile)}\`\n\n${requiredTargetsWithInstall}`, source: 'version' };
}

function checkTargetCompatibility(
  version: Pick<Version, 'required_on_remote'> & { targets: (Pick<VersionTarget, 'targetName'> | null)[] },
  clientRequiredTargets: TargetName[],
  serverRequiredTargets: TargetName[]): string[] {
  const unsupportedClientTargets = clientRequiredTargets.filter((target) =>
    !version.targets.some((verTarget) => verTarget && verTarget.targetName === target));
  const unsupportedServerTargets = serverRequiredTargets.filter((target) =>
    !version.targets.some((verTarget) => verTarget && verTarget.targetName === target));

  if (version.required_on_remote) {
    return [...unsupportedClientTargets, ...unsupportedServerTargets];
  }

  const hasAllClient = clientRequiredTargets.length > 0 && unsupportedClientTargets.length === 0;
  const hasAllServer = serverRequiredTargets.length > 0 && unsupportedServerTargets.length === 0;

  if (hasAllClient) {
    return unsupportedServerTargets;
  }
  if (hasAllServer) {
    return unsupportedClientTargets;
  }
  return [...unsupportedClientTargets, ...unsupportedServerTargets];
}

async function getCompatibilityFor(modReference: string, install: common.Installation, urqlClient: Client): Promise<CompatibilityWithSource> {
  const reportedCompatibility = await getReportedCompatibility(modReference, install.branch, urqlClient);
  if(reportedCompatibility) {
    return { ...reportedCompatibility, source: 'reported' };
  }
  const versionCompatibility = await getVersionCompatibility(modReference, install.version, urqlClient);
  return { ...versionCompatibility, source: 'version' };
}

function gameVersionToSemver(version: number): string {
  return coerce(version)!.format();
}

export async function getReportedCompatibility(modReference: string, gameBranch: common.GameBranch, urqlClient: Client): Promise<Compatibility | undefined> {
  const result = await urqlClient.query(ModReportedCompatibilityDocument, { modReference }).toPromise();
  if(!result.data?.getModByReference) {
    return undefined;
  }
  const mod = result.data.getModByReference;
  if(mod.compatibility) {
    switch(gameBranch) {
      case common.GameBranch.STABLE:
        return mod.compatibility.EA;
      case common.GameBranch.EXPERIMENTAL:
        return mod.compatibility.EXP;
      default:
        throw new Error('Invalid game branch');
    }
  }
  return undefined;
}

interface ModVersion {
  version: string;
  game_version: string;
  dependencies: {
    mod_id: string;
    condition: string;
  }[];
}

async function getModVersions(modReference: string, urqlClient: Client): Promise<ModVersion[] | undefined> {
  if(get(offline)) {
    try {
      return (await OfflineGetMod(modReference)).versions;
    } catch {
      return undefined;
    }
  }
  
  const modQuery = await urqlClient.query(ModVersionsCompatibilityDocument, { modReference }, { requestPolicy: 'cache-first' }).toPromise();
  
  return modQuery.data?.getModByReference?.versions;
}

export async function getVersionCompatibility(modReference: string, gameVersion: number, urqlClient: Client): Promise<Compatibility> {
  const modVersions = await getModVersions(modReference, urqlClient);
  if(!modVersions || modVersions.length === 0) {
    return { state: CompatibilityState.Broken, note: 'Mod has no versions uploaded.' };
  }

  const gameVersionSemver = gameVersionToSemver(gameVersion);

  const compatible = modVersions.some((ver) => ver.game_version === '' || satisfies(gameVersionSemver, ver.game_version));

  if(compatible) {
    return { state: CompatibilityState.Works };
  }
  return { state: CompatibilityState.Broken, note: 'This mod is incompatible with your game version.' };
}
