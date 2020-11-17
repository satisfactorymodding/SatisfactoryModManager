import { ApolloClient } from 'apollo-client';
import { createHttpLink } from 'apollo-link-http';
import { InMemoryCache } from 'apollo-cache-inmemory';

const httpLink = createHttpLink({
  uri: 'https://api.ficsit.app/v2/query',
});

const cache = new InMemoryCache();

// eslint-disable-next-line import/prefer-default-export
export const apolloClient = new ApolloClient({
  link: httpLink,
  cache,
});
