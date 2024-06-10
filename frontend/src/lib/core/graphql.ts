import { cacheExchange } from '@urql/exchange-graphcache';
import { persistedExchange } from '@urql/exchange-persisted';
import type { Client } from '@urql/svelte';
import { createClient, fetchExchange } from '@urql/svelte';

import { schema } from '$lib/generated';

export function initializeGraphQLClient(apiEndpointURL: string): Client {
  return createClient({
    url: apiEndpointURL,
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
          VersionTarget: () => null,
          Mod: (data) => data.mod_reference as string,
        },
        resolvers: {
          Query: {
            getModByReference: (_parent, args) => {
              return { __typename: 'Mod', mod_reference: args.modReference };
            },
          },
        },
      }),
      persistedExchange({
        preferGetForPersistedQueries: true,
      }),
      fetchExchange,
    ],
  });
}
