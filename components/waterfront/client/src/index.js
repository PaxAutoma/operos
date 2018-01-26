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

import React from 'react';
import ReactDOM from 'react-dom';
import {ApolloClient, ApolloProvider} from 'react-apollo';
import 'chart.js';

import LocalNetworkInterface from 'service/LocalNetworkInterface';
import schema from 'service/schema';
import ApiClient from 'service/apiclient';
import App from 'app/App';

// Needed for onTouchTap
// http://stackoverflow.com/a/34015469/988941
import injectTapEventPlugin from 'react-tap-event-plugin';
injectTapEventPlugin();

// Default chart options
Chart.defaults.global.defaultFontFamily = 'Roboto';

const apiClient = new ApiClient(process.env.API_BASE_URL);
const apolloClient = new ApolloClient({
  networkInterface: new LocalNetworkInterface(schema, {apiClient}),
  dataIdFromObject: object => {
    if (object.__typename) {
      if (object.id !== undefined) {
        return `${object.__typename}:${object.id}`;
      } else {
        return `${object.__typename}:<singleton>`;
      }
    }
    return null;
  }
});

ReactDOM.render(
  <ApolloProvider client={apolloClient}>
    <App apiClient={apiClient} />
  </ApolloProvider>,
  document.getElementById('app')
);
