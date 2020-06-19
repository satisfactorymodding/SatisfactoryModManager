import marked from 'marked';
import sanitizeHtml from 'sanitize-html';
import { getSetting, saveSetting } from './settings';

export function lastElement(arr) {
  return arr[arr.length - 1];
}

export function markdownAsElement(markdown) {
  const html = sanitizeHtml(marked(markdown), {
    allowedTags: sanitizeHtml.defaults.allowedTags.concat(['img', 'video', 'details', 'summary', 'h1', 'h2']),
    allowedAttributes: Object.assign(sanitizeHtml.defaults.allowedAttributes, { img: ['src', 'width', 'height'], video: ['src', 'width', 'height', 'controls'] }),
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
