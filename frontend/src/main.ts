import App from './App.svelte';

import { GetAPIEndpoint } from '$wailsjs/go/bindings/App';

const app = new App({
  target: document.getElementById('app')!,
  props: {
    apiEndpointURL: await GetAPIEndpoint(),
  },
});

export default app;
