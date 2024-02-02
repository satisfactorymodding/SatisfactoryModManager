import { Client } from '@urql/svelte';
import { coerce, compare, minVersion, satisfies } from 'semver';
import { get } from 'svelte/store';

import { CompatibilityState, ModReportedCompatibilityDocument, ModVersionsCompatibilityDocument, SmlVersionsCompatibilityDocument, type Compatibility, GetModVersionTargetsDocument, TargetName } from '$lib/generated';
import { offline } from '$lib/store/settingsStore';
import { common } from '$lib/generated/wailsjs/go/models';
import { OfflineGetMod, OfflineGetSMLVersions } from '$wailsjs/go/ficsitcli/ficsitCLI';

export interface CompatibilityWithSource extends Compatibility {
  source: 'reported' | 'version';
}

export async function getCompatibility(modReference: string, gameBranch: common.GameBranch, gameVersion: number, gameTarget: TargetName, urqlClient: Client): Promise<CompatibilityWithSource> {
  if(!await modSupportsTarget(modReference, gameTarget, urqlClient)) {
    return { state: CompatibilityState.Broken, note: `This mod does not support ${gameTarget}.`, source: 'version' };
  }
  const reportedCompatibility = await getReportedCompatibility(modReference, gameBranch, urqlClient);
  if(reportedCompatibility) {
    return { ...reportedCompatibility, source: 'reported' };
  }
  const versionCompatibility = await getVersionCompatibility(modReference, gameVersion, urqlClient);
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
    case common.GameBranch.EARLY_ACCESS:
      return mod.compatibility.EA;
    case common.GameBranch.EXPERIMENTAL:
      return mod.compatibility.EXP;
    default:
      throw new Error('Invalid game branch');
    }
  }
  return undefined;
}

interface SMLVersion {
  version: string;
  satisfactory_version: number;
}

async function getSMLVersions(urqlClient: Client): Promise<SMLVersion[] | undefined> {
  if(get(offline)) {
    return OfflineGetSMLVersions();
  }
  
  const smlVersionsQuery = await urqlClient.query(SmlVersionsCompatibilityDocument, {}).toPromise();

  return smlVersionsQuery.data?.getSMLVersions.sml_versions;
}

interface ModVersion {
  version: string;
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
  const versions = await getSMLVersions(urqlClient);
  if(!versions) {
    return { state: CompatibilityState.Broken };
  }

  const modVersions = await getModVersions(modReference, urqlClient);
  if(!modVersions || modVersions.length === 0) {
    return { state: CompatibilityState.Broken, note: 'Mod has no versions uploaded.' };
  }

  versions.sort((a, b) => compare(a.version, b.version));

  const versionConstraints = versions.map((version, idx, arr) => ({
    version: version.version,
    satisfactory_version: `>=${gameVersionToSemver(version.satisfactory_version)}` + (idx !== arr.length - 1 ? ` <${gameVersionToSemver(arr[idx + 1].satisfactory_version)}` : ''),
  }));

  const compatibleSMLVersions = versionConstraints
    .filter((versionConstraint) => satisfies(gameVersionToSemver(gameVersion), versionConstraint.satisfactory_version))
    .map((versionConstraint) => versionConstraint.version);
  
  const compatible = modVersions.some((ver) => ver.dependencies
    .some((dep) => dep.mod_id === 'SML' && compatibleSMLVersions.some((smlVer) => satisfies(smlVer, dep.condition))));
  const possiblyCompatible = modVersions.some((ver) => ver.dependencies
    .some((dep) => dep.mod_id === 'SML' && satisfies(minVersion(dep.condition)!, '>=3.0.0')));

  if(compatible) {
    return { state: CompatibilityState.Works };
  }
  if(possiblyCompatible) {
    return { state: CompatibilityState.Damaged, note: 'This mod is likely incompatible with your game version and may cause crashes.' };
  }
  return { state: CompatibilityState.Broken, note: 'This mod is incompatible with your game version.' };
}
