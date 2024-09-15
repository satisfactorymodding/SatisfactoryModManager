/* eslint-disable no-restricted-imports */
// This is the one place we're allowed to import popup from skeleton, because we're extending it for everything else to use.
import { type PopupSettings, popup as skeletonPopup } from '@skeletonlabs/skeleton';

export type { PopupSettings } from '@skeletonlabs/skeleton';

/**
 * Because, for whatever reason, skeleton decided to listen to all events in the capture phase,
 * mouseleave events coming from children also trigger the popup to close.
 * So we use this wrapper function to stop the event before it reaches skeleton's event listener,
 * when that event is not coming from the trigger node.
 */
export function popup(triggerNode: HTMLElement, args: PopupSettings) {
  function stopIfNotTrigger(event: MouseEvent) {
    if (event.target !== triggerNode) {
      event.stopImmediatePropagation();
    }
  }

  if (args.event === 'hover') {
    triggerNode.addEventListener('mouseleave', stopIfNotTrigger, true);
  }
  const { update, destroy } = skeletonPopup(triggerNode, args);
  return {
    update,
    destroy() {
      destroy();
      triggerNode.removeEventListener('mouseleave', stopIfNotTrigger, true);
    },
  };
}
