import { type ModalSettings, type ModalStore as SkeletonModalStore, getModalStore as getSkeletonModalStore } from '@skeletonlabs/skeleton';
import _ from 'lodash';
import { getContext, setContext } from 'svelte';
import { get, writable } from 'svelte/store';

const MODAL_STORE_KEY = 'modalStore-extension';

type ModalStore = SkeletonModalStore & {
  trigger: (modal: ModalSettings, top?: boolean) => void;
  triggerUnique: (modal: ModalSettings, top?: boolean) => void;
  close: (component: string) => void;
};

export function getModalStore(): ModalStore {
  const modalStore = getContext<ModalStore | undefined>(MODAL_STORE_KEY);

  if (!modalStore) {
    throw new Error(
      'modalStore is not initialized. Please ensure that `initializeModalStore()` is invoked in the root layout file of this app!',
    );
  }

  return modalStore;
}

/**
 * Initializes the `modalStore`.
 */
export function initializeModalStore(): ModalStore {
  const modalStore = modalService(getSkeletonModalStore());

  return setContext(MODAL_STORE_KEY, modalStore);
}

// For some reason setting the skeleton modalStore too often causes an unreferenced modal to exist in the DOM,
// while the actual modal to be displayed is missing. The modal would show up on the next rerender, but that's weird.
// So we'll use a proxy store that only flushes to the skeleton modalStore at most at 10ms intervals, which seems to not cause the issue.
// 1ms also seems to work, but 10ms is not a noticeable delay.
function modalService(skeletonModalStore: SkeletonModalStore) {
  const proxyStore = writable<ModalSettings[]>(get(skeletonModalStore));
  
  const propagate = _.debounce(() => skeletonModalStore.set(get(proxyStore)), 10);

  return {
    // proxies
    subscribe: proxyStore.subscribe,
    set(mStore: ModalSettings[]) {
      proxyStore.set(mStore);
      propagate();
    },
    update(fn: (mStore: ModalSettings[]) => ModalSettings[]) {
      proxyStore.update(fn);
      propagate();
    },
    /** Append to end of queue. */
    trigger(modal: ModalSettings, top = false) {
      proxyStore.update((mStore) => {
        if (top) {
          mStore.unshift(modal);
        } else {
          mStore.push(modal);
        }
        return mStore;
      });
      propagate();
    },
    /**  Remove first item in queue. */
    close(component = '') {
      proxyStore.update((mStore) => {
        if(component) {
          return mStore.filter((m) => m.component !== component);
        }
        if (mStore.length > 0) mStore.shift();
        return mStore;
      });
      propagate();
    },
    /** Remove all items from queue. */
    clear() {
      proxyStore.set([]);
      propagate();
    },

    // extensions
    triggerUnique(modal: ModalSettings, top = false) {
      const index = get(proxyStore).findIndex((m) => _.isEqual(m, modal));
      if (index === -1) {
        this.trigger(modal, top);
      }
    },
  } as ModalStore;
}
