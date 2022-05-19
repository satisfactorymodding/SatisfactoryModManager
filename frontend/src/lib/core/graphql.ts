import type { Client } from '@urql/svelte';
import { createClient } from '@urql/svelte';
import { cacheExchange } from '@urql/exchange-graphcache';
import { persistedFetchExchange } from '@urql/exchange-persisted-fetch';
import type { LoadInput } from '@sveltejs/kit/types/page';
import schema from '$lib/generated/graphql.schema.urql.json';

export const initializeGraphQLClient = (fetch?: LoadInput['fetch']): Client => createClient({
  url: 'https://api.ficsit.app/v2/query',
  fetch,
  exchanges: [
    cacheExchange({
      schema,
      keys: {
        GetMods: () => null,
        UserMod: () => null,
      },
    }),
    persistedFetchExchange({
      preferGetForPersistedQueries: true,
    }),
  ],
});
