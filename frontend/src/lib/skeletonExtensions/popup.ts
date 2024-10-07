/*
MIT License

Copyright (c) 2024-Preset Skeleton Labs, LLC

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
 */

/**
 * https://github.com/skeletonlabs/skeleton/blob/71cd5b981e408f589fae26feff9da5420218f355/packages/skeleton/src/lib/utilities/Popup/popup.ts
 * Modified to only consider mouseleave events on the trigger node
 */

/* eslint-disable no-restricted-imports */
// This is the one place we're allowed to import popup from skeleton, because we're extending it for everything else to use.
import { type computePosition as computePositionType } from '@floating-ui/dom';
import { type PopupSettings as PopupSettings1, storePopup } from '@skeletonlabs/skeleton';
import { get } from 'svelte/store';

export type PopupSettings = PopupSettings1 & { delay?: number; };

export function popup(triggerNode: HTMLElement, args: PopupSettings) {
  // Floating UI Modules
  const { computePosition, autoUpdate, offset, shift, flip, arrow, size, autoPlacement, hide, inline } = get(storePopup);
  // Local State
  const popupState = {
    open: false,
    autoUpdateCleanup: () => {},
  };
  const focusableAllowedList = ':is(a[href], button, input, textarea, select, details, [tabindex]):not([tabindex="-1"])';
  let focusablePopupElements: HTMLElement[];
  const documentationLink = 'https://www.skeleton.dev/utilities/popups';
  // Elements
  let elemPopup: HTMLElement;
  let elemArrow: HTMLElement;

  function setDomElements(): void {
    elemPopup = document.querySelector(`[data-popup="${args.target}"]`) ?? document.createElement('div');
    elemArrow = elemPopup.querySelector('.arrow') ?? document.createElement('div');
  }
  setDomElements(); // init

  // Render Floating UI Popup
  function render(): void {
    // Error handling for required Floating UI modules
    if (!elemPopup) throw new Error(`The data-popup="${args.target}" element was not found. ${documentationLink}`);
    if (!computePosition) throw new Error(`Floating UI 'computePosition' not found for data-popup="${args.target}". ${documentationLink}`);
    if (!offset) throw new Error(`Floating UI 'offset' not found for data-popup="${args.target}". ${documentationLink}`);
    if (!shift) throw new Error(`Floating UI 'shift' not found for data-popup="${args.target}". ${documentationLink}`);
    if (!flip) throw new Error(`Floating UI 'flip' not found for data-popup="${args.target}". ${documentationLink}`);
    if (!arrow) throw new Error(`Floating UI 'arrow' not found for data-popup="${args.target}". ${documentationLink}`);

    // Bundle optional middleware
    const optionalMiddleware = [];
    // https://floating-ui.com/docs/size
    if (size) optionalMiddleware.push(size(args.middleware?.size));
    // https://floating-ui.com/docs/autoPlacement
    if (autoPlacement) optionalMiddleware.push(autoPlacement(args.middleware?.autoPlacement));
    // https://floating-ui.com/docs/hide
    if (hide) optionalMiddleware.push(hide(args.middleware?.hide));
    // https://floating-ui.com/docs/inline
    if (inline) optionalMiddleware.push(inline(args.middleware?.inline));

    // Floating UI Compute Position
    // https://floating-ui.com/docs/computePosition
    computePosition(triggerNode, elemPopup, {
      placement: args.placement ?? 'bottom',
      // Middleware - NOTE: the order matters:
      // https://floating-ui.com/docs/middleware#ordering
      middleware: [
        // https://floating-ui.com/docs/offset
        offset(args.middleware?.offset ?? 8),
        // https://floating-ui.com/docs/shift
        shift(args.middleware?.shift ?? { padding: 8 }),
        // https://floating-ui.com/docs/flip
        flip(args.middleware?.flip),
        // https://floating-ui.com/docs/arrow
        arrow(args.middleware?.arrow ?? { element: elemArrow || null }),
        // Implement optional middleware
        ...optionalMiddleware,
      ],
    }).then(({ x, y, placement, middlewareData }: Awaited<ReturnType<typeof computePositionType>>) => {
      Object.assign(elemPopup.style, {
        left: `${x}px`,
        top: `${y}px`,
      });
      // Handle Arrow Placement:
      // https://floating-ui.com/docs/arrow
      if (middlewareData.arrow) {
        const { x: arrowX, y: arrowY } = middlewareData.arrow;
        const staticSide = {
          top: 'bottom',
          right: 'left',
          bottom: 'top',
          left: 'right',
        }[placement.split('-')[0]]!;
        Object.assign(elemArrow.style, {
          left: arrowX != null ? `${arrowX}px` : '',
          top: arrowY != null ? `${arrowY}px` : '',
          right: '',
          bottom: '',
          [staticSide]: '-4px',
        });
      }
    });
  }

  let isOpening = false;

  // State Handlers
  function open(): void {
    if (!elemPopup) return;
    // Set open state to on
    popupState.open = true;
    // Return the current state
    if (args.state) args.state({ state: popupState.open });
    // Update render settings
    render();

    isOpening = true;

    setTimeout(() => {
      if (!isOpening) return; // Prevent opening if the mouse has left the trigger node
      // Update the DOM
      elemPopup.style.display = 'block';
      elemPopup.style.opacity = '1';
      elemPopup.style.pointerEvents = 'auto';
      // enable popup interactions
      elemPopup.removeAttribute('inert');
      // Trigger Floating UI autoUpdate (open only)
      // https://floating-ui.com/docs/autoUpdate
      popupState.autoUpdateCleanup = autoUpdate(triggerNode, elemPopup, render);
      // Focus the first focusable element within the popup
      focusablePopupElements = Array.from(elemPopup?.querySelectorAll(focusableAllowedList));
    }, args.delay ?? (args.event === 'hover' ? 400 : 0));
  }
  function close(callback?: () => void): void {
    if (!elemPopup) return;
    // Set transition duration
    const cssTransitionDuration = parseFloat(window.getComputedStyle(elemPopup).transitionDuration.replace('s', '')) * 1000;
    // Set open state to off
    popupState.open = false;
    // Return the current state
    if (args.state) args.state({ state: popupState.open });
    // Update the DOM
    elemPopup.style.opacity = '0';
    // disable popup interactions
    elemPopup.setAttribute('inert', '');

    isOpening = false;

    setTimeout(() => {
      // Cleanup Floating UI autoUpdate (close only)
      if (popupState.autoUpdateCleanup) popupState.autoUpdateCleanup();
      // Trigger callback
      if (callback) callback();
    }, cssTransitionDuration);
  }

  // Event Handlers
  function toggle(): void {
    !popupState.open ? open() : close();
  }
  function onWindowClick(event: MouseEvent): void {
    // Return if the popup is not yet open
    if (!popupState.open) return;
    // Return if click is the trigger element
    if (triggerNode.contains(event.target as HTMLElement)) return;
    // If click it outside the popup
    if (elemPopup && !elemPopup.contains(event.target as HTMLElement)) {
      close();
      return;
    }
    // Handle Close Query State
    const closeQueryString: string = args.closeQuery === undefined ? 'a[href], button' : args.closeQuery;
    // Return if no closeQuery is provided
    if (closeQueryString === '') return;
    const closableMenuElements = elemPopup?.querySelectorAll(closeQueryString);
    closableMenuElements?.forEach((elem) => {
      if (elem.contains(event.target as HTMLElement)) close();
    });
  }

  // Keyboard Interactions for A11y
  const onWindowKeyDown = (event: KeyboardEvent): void => {
    if (!popupState.open) return;
    // Handle keys
    const key: string = event.key;
    // On Esc key
    if (key === 'Escape') {
      event.preventDefault();
      triggerNode.focus();
      close();
      return;
    }
    // Update focusable elements (important for Autocomplete)
    focusablePopupElements = Array.from(elemPopup?.querySelectorAll(focusableAllowedList));
    // On Tab or ArrowDown key
    const triggerMenuFocused: boolean = popupState.open && document.activeElement === triggerNode;
    if (
      triggerMenuFocused &&
        (key === 'ArrowDown' || key === 'Tab') &&
        focusableAllowedList.length > 0 &&
        focusablePopupElements.length > 0
    ) {
      event.preventDefault();
      focusablePopupElements[0].focus();
    }
  };

  function ifTriggerNode(func: () => void) {
    return (event: Event) => {
      if (event.target === triggerNode) {
        func();
      }
    };
  }

  const closeIfTriggerNode = ifTriggerNode(() => close());

  // Event Listeners
  switch (args.event) {
    case 'click':
      triggerNode.addEventListener('click', toggle, true);
      window.addEventListener('click', onWindowClick, true);
      break;
    case 'hover':
      triggerNode.addEventListener('mouseover', open, true);
      triggerNode.addEventListener('mouseleave', closeIfTriggerNode, true);
      break;
    case 'focus-blur':
      triggerNode.addEventListener('focus', toggle, true);
      triggerNode.addEventListener('blur', () => close(), true);
      break;
    case 'focus-click':
      triggerNode.addEventListener('focus', open, true);
      window.addEventListener('click', onWindowClick, true);
      break;
    default:
      throw new Error(`Event value of '${args.event}' is not supported. ${documentationLink}`);
  }
  window.addEventListener('keydown', onWindowKeyDown, true);

  // Render popup on initialization
  render();

  // Lifecycle
  return {
    update(newArgs: PopupSettings) {
      close(() => {
        args = newArgs;
        render();
        setDomElements();
      });
    },
    destroy() {
      // Trigger Events
      triggerNode.removeEventListener('click', toggle, true);
      triggerNode.removeEventListener('mouseover', open, true);
      triggerNode.removeEventListener('mouseleave', closeIfTriggerNode, true);
      triggerNode.removeEventListener('focus', toggle, true);
      triggerNode.removeEventListener('focus', open, true);
      triggerNode.removeEventListener('blur', () => close(), true);
      // Window Events
      window.removeEventListener('click', onWindowClick, true);
      window.removeEventListener('keydown', onWindowKeyDown, true);
    },
  };
}