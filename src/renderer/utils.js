import marked from 'marked';
import sanitizeHtml from 'sanitize-html';
import originalFilenamify from 'filenamify';
import { getSetting, saveSetting } from '~/settings';

export function lastElement(arr) {
  return arr[arr.length - 1];
}

export function markdownAsElement(markdown) {
  const html = sanitizeHtml(marked(markdown), {
    allowedTags: sanitizeHtml.defaults.allowedTags.concat(['img', 'video', 'details', 'summary', 'source', 'h1', 'h2']),
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

/**
 * @param {Date} date The date
 */
export function filenameFriendlyDate(date) {
  const year = date.getUTCFullYear();
  const month = date.getUTCMonth();
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
