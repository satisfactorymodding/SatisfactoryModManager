import { queue } from 'async';
import { derived, get, writable } from 'svelte/store';

import { queueAutoStart } from './settingsStore';

interface QueuedAction<T> {
  mod: string;
  action: 'install' | 'remove' | 'enable' | 'disable';
  func: () => Promise<T>;
}

export const hasPendingProfileChange = writable(false);

const queuedActionsInternal = writable<QueuedAction<unknown>[]>([]);
export const queuedMods = derived(queuedActionsInternal, (actions) => actions.map((a) => ({ ...a, func: undefined })));
const modActionsQueue = queue((task: () => Promise<unknown>, cb) => {
  const complete = (e?: Error) => {
    queuedActionsInternal.set(get(queuedActionsInternal).filter((a) => a.func !== task));
    cb(e);
    hasPendingProfileChange.set(false);
  };
  task().then(() => complete()).catch(complete);
});

modActionsQueue.empty(() => {
  if(!get(queueAutoStart)) {
    modActionsQueue.pause();
  }
});

queueAutoStart.subscribe((val) => {
  if(val) {
    modActionsQueue.resume();
  } else {
    modActionsQueue.pause();
  }
});

export function startQueue() {
  modActionsQueue.resume();
}

export async function addQueuedModAction<T>(mod: string, action: string, func: () => Promise<T>): Promise<T> {
  const queuedAction = { mod, action, func } as QueuedAction<T>;
  queuedActionsInternal.set([
    ...get(queuedActionsInternal),
    queuedAction,
  ]);
  if(get(queueAutoStart)) {
    startQueue();
  }
  return modActionsQueue.pushAsync(func);
}

export function removeQueuedModAction(mod: string) {
  const queuedAction = get(queuedActionsInternal).find((a) => a.mod === mod);
  if(!queuedAction) {
    return;
  }
  modActionsQueue.remove((a) => a.data === queuedAction.func);
  queuedActionsInternal.set(get(queuedActionsInternal).filter((a) => a.mod !== mod));
}
