/*
Copyright 2018 Pax Automa Systems, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import {makeExecutableSchema} from 'graphql-tools';
import GraphQLJSON from 'graphql-type-json';

const typeDefs = `
scalar JSON

enum NodeStatus {
  NOT_READY
  READY
}

type Node {
  id: String!
  status: NodeStatus!
  ip: String
  pod_cidr: String
  hardware_info: String
}

type ClusterInfo {
  settings: JSON
}

type User {
  username: String!
  loging_time: String!
}

type LoginInfo {
  logged_in: Boolean!
  user: User
}

type Query {
  nodes: [Node]!
  node(id: String!): Node
  cluster_info: ClusterInfo!
  login_info: LoginInfo!
}

type Mutation {
  login(username: String!, password: String!): LoginInfo!
  logout: LoginInfo!
}

schema {
  query: Query
  mutation: Mutation
}
`;

const resolvers = {
  Query: {
    nodes: (_obj, _args, {apiClient}) => apiClient.getNodes(),
    node: (_obj, args, {apiClient}) => apiClient.getNode(args.id),
    cluster_info: (_obj, _args, {apiClient}) => apiClient.getClusterInfo(),
    login_info: (_obj, _args, {apiClient}) => apiClient.getLoginInfo()
  },
  Node: {
    status: node => node.status || 'NOT_READY'
  },
  Mutation: {
    login: (_obj, args, {apiClient}) => apiClient.login(args.username, args.password),
    logout: (_obj, args, {apiClient}) => apiClient.logout()
  },
  JSON: GraphQLJSON
};

export default makeExecutableSchema({typeDefs, resolvers});
