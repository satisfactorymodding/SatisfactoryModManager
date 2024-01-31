import { type ModalStore, getModalStore as getSkeletonModalStore, type ModalSettings } from '@skeletonlabs/skeleton';
import _ from 'lodash';
import { getContext, setContext } from 'svelte';
import { get, writable } from 'svelte/store';

const MODAL_STORE_KEY = 'modalStore-extension';

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
function modalService(skeletonModalStore: ModalStore) {
  const proxyStore = writable<ModalSettings[]>(get(skeletonModalStore));
  
  const propagate = _.debounce(() => skeletonModalStore.set(get(proxyStore)), 10);

  return {
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
    trigger(modal: ModalSettings) {
      proxyStore.update((mStore) => {
        mStore.push(modal);
        return mStore;
      });
      propagate();
    },
    /**  Remove first item in queue. */
    close() {
      proxyStore.update((mStore) => {
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
  } as ModalStore;
}
