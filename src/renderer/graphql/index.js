import {
  ApolloClient, createHttpLink, ApolloLink, InMemoryCache,
} from '@apollo/client/core';
import { createPersistedQueryLink } from '@apollo/client/link/persisted-queries';
import { withScalars } from 'apollo-link-scalars';
import sha from 'sha.js';
import { DateTimeResolver } from 'graphql-scalars';
import { buildClientSchema } from 'graphql';
import schema from './__generated__/graphql.schema.json';

const link = ApolloLink.from([
  withScalars({
    schema: buildClientSchema(schema),
    typesMap: {
      Date: {
        ...DateTimeResolver,
        parseValue(value) {
          if (typeof value !== 'string' || value) {
            return DateTimeResolver.parseValue(value);
          }
          return null;
        },
        parseLiteral(value, variables) {
          if (typeof value !== 'string' || value) {
            return DateTimeResolver.parseLiteral(value, variables);
          }
          return null;
        },
        serialize(value) {
          if (value instanceof Date) {
            return value.toISOString();
          }
          return value;
        },
      },
    },
  }),
  createPersistedQueryLink({ useGETForHashedQueries: true, sha256: (...args) => sha('sha256').update(args.toString()).digest('hex') }),
  createHttpLink({
    uri: 'https://api.ficsit.app/v2/query',
  }),
]);

const cache = new InMemoryCache();

// eslint-disable-next-line import/prefer-default-export
export const apolloClient = new ApolloClient({
  link,
  cache,
});
