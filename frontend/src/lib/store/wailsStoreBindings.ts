import { derived, writable } from 'svelte/store';

import { EventsOn } from '$wailsjs/runtime/runtime';

function oneWayWritableBinding<T, D>(defaultValue: D, mainToRenderer: {
  initialGet?: () => Promise<T>
  updateEvent?: string,
  allowNull?: boolean,
}) {
  const { updateEvent, allowNull, initialGet } = {
    allowNull: true,
    ...mainToRenderer,
  };

  const initialized = writable(false);

  const { subscribe, set } = writable(defaultValue as T | D);

  const setData = (data: T) => {
    if(data === null && !allowNull) {
      set(defaultValue);
    } else {
      set(data);
    }
  };

  if (updateEvent) {
    EventsOn(updateEvent, setData);
  }

  if(initialGet) {
    initialGet().then(setData).then(() => initialized.set(true));
  }

  return {
    subscribe,
    isInit: derived(initialized, (i) => i),
    set: setData,
  };
}

export function binding<T>(defaultValue: T,
  mainToRenderer: {
    initialGet?: () => Promise<T>
    updateEvent?: string,
    allowNull?: boolean,
  },
) {
  const { subscribe, isInit } = oneWayWritableBinding(defaultValue, mainToRenderer);

  return {
    subscribe,
    isInit,
  };
}

export function bindingTwoWay<T, D>(defaultValue: D,
  mainToRenderer: {
    initialGet?: () => Promise<T>
    updateEvent?: string,
    allowNull?: boolean,
  },
  rendererToMain: {
    updateFunction: (value: T) => Promise<void>,
  },
) {
  const { subscribe, isInit, set } = oneWayWritableBinding(defaultValue, mainToRenderer);
  const { updateFunction } = rendererToMain;

  return {
    subscribe,
    isInit,
    asyncSet: async (value: T) => {
      set(value);      
      await updateFunction(value);
    },
  };
}

export function bindingTwoWayNoExcept<T>(defaultValue: T,
  mainToRenderer: {
    initialGet?: () => Promise<T>
    updateEvent?: string,
    allowNull?: boolean,
  },
  rendererToMain: {
    updateFunction: (value: T) => Promise<void>,
  },
) {
  const { subscribe, isInit, set } = oneWayWritableBinding(defaultValue, mainToRenderer);
  const { updateFunction } = rendererToMain;

  return {
    subscribe,
    isInit,
    set: (value: T) => {
      set(value);
      updateFunction(value); // must not throw
    },
  };
}
