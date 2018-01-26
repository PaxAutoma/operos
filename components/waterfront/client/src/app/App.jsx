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
import PropTypes from 'prop-types';
import {BrowserRouter} from 'react-router-dom';
import {graphql, gql} from 'react-apollo';
import {MuiThemeProvider, createMuiTheme} from 'material-ui/styles';

import Chrome from 'app/Chrome';
import LoginScreen from 'app/LoginScreen';
import nodesModule from 'views/nodes';
import accessModule from 'views/access';
import clusterModule from 'views/cluster';
import metricsModule from 'views/metrics';
import aboutModule from 'views/about';

const theme = createMuiTheme({
  palette: {
    background: {
      appBar: '#193f6d'
    }
  }
});

class App extends React.Component {
  static childContextTypes = {
    apiClient: PropTypes.object
  }

  static propTypes = {
    apiClient: PropTypes.object.isRequired
  }

  getChildContext() {
    return {
      apiClient: this.props.apiClient
    };
  }

  render() {
    if (this.props.data.loading) {
      return <div>Loading...</div>
    } else if (!this.props.data.login_info.logged_in) {
      return <LoginScreen />
    }

    const menu = [
      clusterModule.menuItem,
      metricsModule.menuItem,
      nodesModule.menuItem,
      accessModule.menuItem,
      aboutModule.menuItem
    ];

    const views = [
      ...clusterModule.views,
      ...metricsModule.views,
      ...nodesModule.views,
      ...accessModule.views,
      ...aboutModule.views
    ];

    return (
        <BrowserRouter>
          <MuiThemeProvider theme={theme}>
            <Chrome menu={menu} views={views} />
          </MuiThemeProvider>
        </BrowserRouter>
    );
  }
}

export default graphql(gql`
  query {
    login_info {
      logged_in
      user {
        username
      }
    }
  }
`)(App);
