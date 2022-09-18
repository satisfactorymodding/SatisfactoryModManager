import { readable, writable, type Readable, type Writable } from 'svelte/store';
import { EventsOff, EventsOn } from '$wailsjs/runtime/runtime';

export interface ReadableBinding<T> extends Readable<T> {
  isInit: boolean;
}
export interface WritableBinding<T> extends Writable<T> {
  isInit: boolean;
}

export function readableBinding<T>(defaultValue: T,
  options: {
    updateEvent?: string,
    allowNull?: boolean,
    initialGet?: () => Promise<T>
  }
) {
  const { updateEvent, allowNull, initialGet } = {
    allowNull: true,
    ...options
  };
  const store = {
    isInit:false,
    ...readable(defaultValue, (set) => {
      const setData = (data) => {
        if(data === null && !allowNull) {
          set(defaultValue);
        } else {
          set(data);
        }
      };

      EventsOn(updateEvent, setData);

      if(initialGet) {
        initialGet().then(setData).then(() => store.isInit = true);
      }

      return () => {
        EventsOff(updateEvent);
      };
    })
  } as ReadableBinding<T>;

  return store;
}

export function writableBinding<T>(defaultValue: T,
  options: {
    initialGet?: () => Promise<T>
  }
) {
  const { initialGet } = {
    ...options
  };

  const store = {
    isInit: false,
    ...writable(defaultValue, (set) => {
      if(initialGet) {
        initialGet().then(set).then(() => store.isInit = true);
      }
    })
  } as WritableBinding<T>;

  return store;
}