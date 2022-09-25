import { readable, writable, type Readable, type Writable } from 'svelte/store';
import { EventsOff, EventsOn } from '$wailsjs/runtime/runtime';

export type InitCallback = () => void;
export interface ReadableBinding<T> extends Readable<T> {
  isInit: boolean;
  subscribeInit(callback: InitCallback): void;
  waitForInit: Promise<void>;
}
export interface WritableBinding<T> extends Writable<T> {
  isInit: boolean;
  subscribeInit(callback: InitCallback): void;
  waitForInit: Promise<void>;
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
  
  const initCallbacks: InitCallback[] = [];
  let resolveInit: () => void = () => {/* empty */};
  const initPromise = new Promise<void>((resolve) => {
    resolveInit = resolve;
  });

  const store = {
    isInit: false,
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
        initialGet().then(setData).then(() => store.isInit = true).then(() => initCallbacks.forEach((cb) => cb())).then(resolveInit);
      }

      return () => {
        EventsOff(updateEvent);
      };
    }),
    waitForInit: initPromise,
    subscribeInit(callback: () => void) {
      if(this.isInit) {
        callback();
      } else {
        initCallbacks.push(callback);
      }
    }
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
  
  const initCallbacks: InitCallback[] = [];
  let resolveInit: () => void = () => {/* empty */};
  const initPromise = new Promise<void>((resolve) => {
    resolveInit = resolve;
  });

  const store = {
    isInit: false,
    ...writable(defaultValue, (set) => {
      if(initialGet) {
        initialGet().then(set).then(() => store.isInit = true).then(() => initCallbacks.forEach((cb) => cb())).then(resolveInit);
      }
    }),
    waitForInit: initPromise,
    subscribeInit(callback: () => void) {
      if(this.isInit) {
        callback();
      } else {
        initCallbacks.push(callback);
      }
    }
  } as WritableBinding<T>;

  return store;
}

export function writableBindingSync<T>(defaultValue: T,
  options: {
    initialGet?: () => Promise<T>,
    updateFunction: (value: T) => Promise<void>
  }
) {
  const { initialGet } = {
    ...options
  };
  
  const initCallbacks: InitCallback[] = [];
  let resolveInit: () => void = () => {/* empty */};
  const initPromise = new Promise<void>((resolve) => {
    resolveInit = resolve;
  });

  const store = {
    isInit: false,
    ...writable(defaultValue, (set) => {
      if(initialGet) {
        initialGet().then(set).then(() => store.isInit = true).then(() => initCallbacks.forEach((cb) => cb())).then(resolveInit);
      }
    }),
    waitForInit: initPromise,
    subscribeInit(callback: () => void) {
      if(this.isInit) {
        callback();
      } else {
        initCallbacks.push(callback);
      }
    }
  } as WritableBinding<T>;

  store.subscribe((value) => {
    options.updateFunction(value);
  });

  return store;
}