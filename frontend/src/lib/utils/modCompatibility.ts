import { Client } from '@urql/svelte';
import { coerce, satisfies } from 'semver';
import { get } from 'svelte/store';

import { type Compatibility, CompatibilityState, GetModVersionTargetsDocument, ModReportedCompatibilityDocument, ModVersionsCompatibilityDocument, TargetName } from '$lib/generated';
import { installsMetadata, selectedInstallMetadata, selectedProfileTargets } from '$lib/store/ficsitCLIStore';
import { offline } from '$lib/store/settingsStore';
import { OfflineGetMod } from '$wailsjs/go/ficsitcli/ficsitCLI';
import { common } from '$wailsjs/go/models';

export interface CompatibilityWithSource extends Compatibility {
  source: 'reported' | 'version';
}

export async function getCompatibility(modReference: string, urqlClient: Client): Promise<CompatibilityWithSource> {
  const installInfo = get(selectedInstallMetadata).info;
  if(!installInfo) {
    return { state: CompatibilityState.Broken, note: 'No game selected.', source: 'version' };
  }

  for await (const [gameTarget, installs] of Object.entries(get(selectedProfileTargets))) {
    const installNames = installs.map((install) => friendlyInstallName(get(installsMetadata)[install].info!));
    if(!await modSupportsTarget(modReference, gameTarget as TargetName, urqlClient)) {
      return { state: CompatibilityState.Broken, note: `This mod does not support ${gameTarget}, required on this profile by ${installNames.join(', ')}.`, source: 'version' };
    }
  }

  return await getCompatibilityFor(modReference, installInfo, urqlClient);
}

function friendlyInstallName(install: common.Installation) {
  const installType = install.type === common.InstallType.WINDOWS ? 'Game' : 'Server';
  return `${install.launcher} (${installType})`;
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

export async function modSupportsTarget(modReference: string, gameTarget: TargetName, urqlClient: Client): Promise<boolean> {
  const result = await urqlClient.query(GetModVersionTargetsDocument, { modReference }).toPromise();
  if(!result.data?.mod) {
    return false;
  }
  const mod = result.data.mod;
  return mod.versions.some((ver) => ver.targets.some((target) => target?.targetName === gameTarget));
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
