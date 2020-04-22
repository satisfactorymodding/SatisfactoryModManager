import semver from 'semver';

export function isModInstalled(state, mod) {
  return !!state.selectedSatisfactoryInstall.mods[mod.id];
}

export function isModCompatible(mod) {
  return mod.versions.length !== 0 && semver.satisfies(mod.versions[0].sml_version, '>=2.0.0');
}

export function isModVersionInstalled(state, mod) {
  return mod.versions[0]
  && mod.versions[0].version
  && state.selectedSatisfactoryInstall.mods[mod.id] === mod.versions[0].version;
}

export function isModHasUpdate(state, mod) {
  return isModCompatible(mod)
  && isModInstalled(state, mod)
  && !isModVersionInstalled(state, mod)
  && mod.versions.length > 0
  && (mod.versions[0].sml_version === '2.1.0' || mod.versions[0].sml_version === 'v2.1.0') === state.selectedSatisfactoryInstall.name.toLowerCase().includes('experimental');
}

export function isVersionSML20Compatible(version) {
  return semver.satisfies(version.sml_version, '>=2.0.0');
}
