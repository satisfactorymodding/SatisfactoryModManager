import { CompatibilityState, ModVersionsCompatibilityDocument, SmlVersionsCompatibilityDocument, type Compatibility, type Mod } from '$lib/generated';
import { getClient } from '@urql/svelte';
import { coerce, compare, minVersion, satisfies } from 'semver';

function gameVersionToSemver(version: number): string | null {
  return coerce(version)?.format();
}

export function getReportedCompatibility(mod: Pick<Mod, 'compatibility' | 'mod_reference'>, gameBranch: GameBranch): Compatibility | undefined {
  if(mod.compatibility) {
    switch(gameBranch) {
    case 'EA':
    case 'Early Access':
      return mod.compatibility.EA;
    case 'EXP':
    case 'Experimental':
      return mod.compatibility.EA;
    default:
      throw new Error('Invalid game branch');
    }
  }
  return undefined;
}

export async function getVersionCompatibility(modReference: string, gameVersion: number): Promise<Compatibility> {
  const urqlclient = getClient();

  const smlVersionsQuery = await urqlclient.query(SmlVersionsCompatibilityDocument).toPromise();
  const versions = smlVersionsQuery.data?.getSMLVersions.sml_versions;
  if(!versions) {
    return CompatibilityState.Broken;
  }

  const modQuery = await urqlclient.query(ModVersionsCompatibilityDocument, { modReference }).toPromise();
  const modVersions = modQuery.data?.getModByReference?.versions;
  if(!modVersions) {
    return CompatibilityState.Broken;
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
    .some((dep) => dep.mod_id === 'SML' && satisfies(minVersion(dep.condition)), '>=3.0.0'));

  if(compatible) {
    return { state: CompatibilityState.Works };
  }
  if(possiblyCompatible) {
    return { state: CompatibilityState.Damaged, note: 'This mod is likely incompatible with your game version and may cause crashes.' };
  }
  return { state: CompatibilityState.Broken, note: 'This mod is incompatible with your game version.' };
}