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
import {withStyles} from 'material-ui/styles';

import CredentialsCard from 'views/access/CredentialsCard';
import RootPasswordCard from 'views/access/RootPasswordCard';


const styles = {
  container: {
    display: 'flex',
    flexDirection: 'column'
  },
  card: {
    maxWidth: 800,
    marginBottom: 16
  }
};

class AccessManagementView extends React.Component {
  render() {
    const {classes} = this.props;

    return (
      <div className={classes.container}>
        <CredentialsCard className={classes.card} />
        <RootPasswordCard className={classes.card} />
      </div>
    );
  }
}

export default withStyles(styles)(AccessManagementView);
