import { Client } from '@urql/svelte';
import { coerce, satisfies } from 'semver';
import { get } from 'svelte/store';

import {
  type Compatibility,
  CompatibilityState,
  ModReportedCompatibilityDocument,
  ModVersionsCompatibilityDocument,
  TargetName,
  type Version,
  type VersionTarget,
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

export async function getCompatibility(modReference: string, urqlClient: Client, offline = false): Promise<CompatibilityWithSource> {
  const installInfo = get(selectedInstallMetadata).info;
  if(!installInfo) {
    return { state: CompatibilityState.Broken, note: 'No game selected.', source: 'version' };
  }

  const reportedCompatibility = !offline ? await getReportedCompatibility(modReference, installInfo.branch, urqlClient) : undefined;
  if(reportedCompatibility && reportedCompatibility.state !== CompatibilityState.Works) {
    return { ...reportedCompatibility, source: 'reported' };
  }

  const modVersions = await getModVersions(modReference, urqlClient);
  if(!modVersions || modVersions.length === 0) {
    return { state: CompatibilityState.Broken, note: 'Mod has no versions uploaded.', source: 'version' };
  }

  const compatibleModVersions = getVersionCompatibleVersions(modVersions);
  if (compatibleModVersions.length === 0) {
    return { state: CompatibilityState.Broken, note: 'This mod is incompatible with your game version.', source: 'version' };
  }

  const targetCompatibility = getTargetCompatibilityFor(modVersions, compatibleModVersions);
  if (targetCompatibility.state !== CompatibilityState.Works) {
    return targetCompatibility;
  }

  if (reportedCompatibility) {
    return { ...reportedCompatibility, source: 'reported' };
  }

  return { state: CompatibilityState.Works, source: 'version' };
}

function friendlyInstallName(install: common.Installation) {
  const installType = install.type === common.InstallType.WINDOWS ? 'Game' : 'Server';
  return `${install.launcher} (${installType})`;
}

function getTargetCompatibilityFor(modVersions: ModVersion[], versionCompatibleVersions: ModVersion[]): CompatibilityWithSource {
  const requiredTargets = Object.entries(get(selectedProfileTargets));
  if (requiredTargets.length === 0) {
    return { state: CompatibilityState.Works, source: 'version' };
  }

  const clientRequiredTargets = requiredTargets.filter(([target]) => clientTargets.includes(target as TargetName)).map(([target]) => target as TargetName);
  const serverRequiredTargets = requiredTargets.filter(([target]) => serverTargets.includes(target as TargetName)).map(([target]) => target as TargetName);

  const versionCompatibility = modVersions.map((version) => ({ version, missingTargets: checkTargetCompatibility(version, clientRequiredTargets, serverRequiredTargets) }));
  const hasEverBeenCompatible = versionCompatibility.some((comp) => comp.missingTargets.length === 0);

  // The versions should be sorted by newest first
  const latestVersionCompatibility = versionCompatibility.filter(({ version }) => versionCompatibleVersions.includes(version))[0];

  if (latestVersionCompatibility.missingTargets.length === 0) {
    return { state: CompatibilityState.Works, source: 'version' };
  }

  const missingTargets = requiredTargets.filter(([target]) => latestVersionCompatibility.missingTargets.includes(target as TargetName));

  const requiredTargetsWithInstall = missingTargets
    .map(([target, installs]) =>
      `${target} required by: ${installs.map((install) => friendlyInstallName(get(installsMetadata)[install].info!)).join(',')}`)
    .join('\n\n');

  if (hasEverBeenCompatible) {
    return { state: CompatibilityState.Damaged, note: `An older version of this mod supports all your installs using profile \`${get(selectedProfile)}\`, but the latest version does not\n\n${requiredTargetsWithInstall}`, source: 'version' };
  }
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
  required_on_remote: boolean;
  targets: { targetName: TargetName }[];
}

async function getModVersions(modReference: string, urqlClient: Client): Promise<ModVersion[] | undefined> {
  if(get(offline)) {
    try {
      return (await OfflineGetMod(modReference)).versions.map((ver) => ({
        ...ver,
        targets: ver.targets.map((target) => ({ ...target, targetName: target.target_name as TargetName })),
      }));
    } catch {
      return undefined;
    }
  }
  
  const modQuery = await urqlClient.query(ModVersionsCompatibilityDocument, { modReference }, { requestPolicy: 'cache-first' }).toPromise();

  // This cast can potentially cause issues later when ModVersion changes, but the API schema allows null for targets[number], even though it would never happen
  return modQuery.data?.getModByReference?.versions as ModVersion[];
}

function getVersionCompatibleVersions(modVersions: ModVersion[]): ModVersion[] {
  const installsToCheck = Object.values(get(selectedProfileTargets)).flat();
  const metadata = get(installsMetadata);
  return modVersions.filter((version) =>
    version.game_version === '' || installsToCheck.some((install) => satisfies(gameVersionToSemver(metadata[install].info?.version ?? 0), version.game_version)),
  );
}
