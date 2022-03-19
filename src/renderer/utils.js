import { marked } from 'marked';
import sanitizeHtml from 'sanitize-html';
import originalFilenamify from 'filenamify';
import {
  valid, coerce, eq, validRange, satisfies, minVersion, compare,
} from 'semver';
import gql from 'graphql-tag';
import { getSetting, saveSetting } from '~/settings';
import { apolloClient } from './graphql';

export function lastElement(arr) {
  return arr[arr.length - 1];
}

export function markdownAsElement(markdown) {
  const html = sanitizeHtml(marked(markdown), {
    allowedTags: sanitizeHtml.defaults.allowedTags.concat(['img', 'video', 'details', 'summary', 'source', 'strike', 'blockquote', 'del', 'ins', 's']),
    allowedAttributes: Object.assign(sanitizeHtml.defaults.allowedAttributes, { img: ['src', 'width', 'height'], video: ['src', 'width', 'height', 'controls'], source: ['src', 'type'] }),
  });
  const el = document.createElement('html');
  el.innerHTML = html;
  return el;
}

export function ignoreUpdate(item, version) {
  const ignoredUpdates = getSetting('ignoredUpdates', []);
  if (!ignoredUpdates.some((ignoredUpdate) => ignoredUpdate.item === item && ignoredUpdate.version === version)) {
    ignoredUpdates.push({ item, version });
  }
  saveSetting('ignoredUpdates', ignoredUpdates);
  return ignoredUpdates;
}

export function unignoreUpdate(item, version) {
  const ignoredUpdates = getSetting('ignoredUpdates', []);
  ignoredUpdates.removeWhere((update) => update.item === item && update.version === version);
  saveSetting('ignoredUpdates', ignoredUpdates);
  return ignoredUpdates;
}

export function validAndEq(v1, v2) {
  const v1Valid = valid(coerce(v1));
  const v2Valid = valid(coerce(v2));
  if (v1Valid && v2Valid) {
    return eq(v1Valid, v2Valid);
  }
  return false;
}

/**
 * @param {Date} date The date
 */
export function filenameFriendlyDate(date) {
  const year = date.getUTCFullYear();
  const month = date.getUTCMonth() + 1;
  const day = date.getUTCDate();
  const hour = date.getUTCHours();
  const minute = date.getUTCMinutes();
  const second = date.getUTCSeconds();
  return `${year}-${month}-${day}_${hour}-${minute}-${second}`;
}

export function roundWithDecimals(number, decimals = 0) {
  return Math.round(number * (10 ** decimals)) / (10 ** decimals);
}

const sizeRanges = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

export function bytesToAppropriate(bytes) {
  let rangeNum = 0;
  while (bytes >= 1024 ** (rangeNum + 1)) {
    rangeNum += 1;
  }
  return `${roundWithDecimals(bytes / (1024 ** rangeNum), 2).toFixed(2)} ${sizeRanges[rangeNum]}`;
}

const timeRanges = {
  sec: 1,
  min: 60,
  h: 60 * 60,
  days: 60 * 60 * 24,
};

export function secondsToAppropriate(seconds) {
  const ranges = Object.keys(timeRanges);
  let rangeNum = 0;
  while (rangeNum < ranges.length - 1 && seconds >= timeRanges[ranges[rangeNum + 1]]) {
    rangeNum += 1;
  }
  return `${roundWithDecimals(seconds / timeRanges[ranges[rangeNum]], 0)}${ranges[rangeNum]}`;
}

export function filenamify(str) {
  return originalFilenamify(str, { replacement: '_' });
}

export function setIntervalImmediately(func, interval) {
  func();
  return setInterval(func, interval);
}

export const COMPATIBILITY_LEVEL = {
  INCOMPATIBLE: 'incompatible',
  POSSIBLY_COMPATIBLE: 'possibly_compatible',
  COMPATIBLE: 'compatible',
};

/* eslint-disable camelcase */
export async function isCompatibleFast(mod, gameVersion) {
  if (!mod || mod.versions.length === 0) {
    return false;
  }
  const smlVersions = [...(await apolloClient.query({
    query: gql`
      query smlVersions {
        getSMLVersions(filter: {limit: 100}) {
          sml_versions {
            id,
            version,
            satisfactory_version,
          }
        }
      }
    `,
  })).data.getSMLVersions.sml_versions];
  smlVersions.sort((a, b) => compare(a.version, b.version));
  const versionConstraints = smlVersions.map((version, idx, arr) => ({
    version: version.version,
    satisfactory_version: idx !== arr.length - 1 ? `>=${valid(coerce(version.satisfactory_version))} <${valid(coerce(arr[idx + 1].satisfactory_version))}` : `>=${valid(coerce(version.satisfactory_version))}`,
  }));
  const compatibleSMLVersions = versionConstraints
    .filter((versionConstraint) => satisfies(valid(coerce(gameVersion)), versionConstraint.satisfactory_version))
    .map((versionConstraint) => versionConstraint.version);
  const SML3_0_0_gameVersion = smlVersions.find((smlVer) => smlVer.version === '3.0.0').satisfactory_version;
  const SML3_2_0_gameVersion = smlVersions.find((smlVer) => smlVer.version === '3.2.0').satisfactory_version;
  let minSMLVersionForGame = '2.0.0';
  if (gameVersion > SML3_2_0_gameVersion) {
    minSMLVersionForGame = '3.2.0';
  } else if (gameVersion > SML3_0_0_gameVersion) {
    minSMLVersionForGame = '3.0.0';
  }
  const compatible = mod.versions.some((ver) => (
    validRange(ver.sml_version)
    && satisfies(minVersion(ver.sml_version), `>=${minSMLVersionForGame}`)
    && compatibleSMLVersions.some((smlVer) => satisfies(smlVer, validRange(ver.sml_version)))
  ));
  const possibleCompatible = mod.versions.some((ver) => (
    validRange(ver.sml_version)
    && satisfies(minVersion(ver.sml_version), '>=3.0.0')
    && compatibleSMLVersions.some((smlVer) => satisfies(smlVer, validRange(ver.sml_version)))
  ));
  if (compatible) {
    return COMPATIBILITY_LEVEL.COMPATIBLE;
  }
  if (possibleCompatible) {
    return COMPATIBILITY_LEVEL.POSSIBLY_COMPATIBLE;
  }
  return COMPATIBILITY_LEVEL.INCOMPATIBLE;
}
/* eslint-enable camelcase */
