import App from './App.svelte';

import { GetAPIEndpoint, GetSiteEndpoint } from '$wailsjs/go/bindings/App';

const app = new App({
  target: document.getElementById('app')!,
  props: {
    apiEndpointURL: await GetAPIEndpoint(),
    siteEndpointURL: await GetSiteEndpoint(),
  },
});

export default app;
