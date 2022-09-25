import type { Client } from '@urql/svelte';
import { createClient } from '@urql/svelte';
import { cacheExchange } from '@urql/exchange-graphcache';
import { persistedFetchExchange } from '@urql/exchange-persisted-fetch';
import schema from '$lib/generated/graphql.schema.urql.json';

export function initializeGraphQLClient(): Client {
  return createClient({
    url: 'https://api.ficsit.app/v2/query',
    exchanges: [
      cacheExchange({
        schema,
        keys: {
          GetMods: () => null,
          GetSMLVersions: () => null,
          LatestVersions: () => null,
          UserMod: () => null,
          GetGuides: () => null,
          OAuthOptions: () => null,
          UserRoles: () => null,
          Compatibility: () => null,
          CompatibilityInfo: () => null,
          VersionDependency: () => null,
        },
      }),
      persistedFetchExchange({
        preferGetForPersistedQueries: true,
      }),
    ],
  });
}