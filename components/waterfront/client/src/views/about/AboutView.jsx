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
import {graphql, gql} from 'react-apollo';

import Typography from 'material-ui/Typography';
import Table, { TableBody, TableCell, TableHead, TableRow } from 'material-ui/Table';
import {withStyles} from 'material-ui/styles';

const styles = {
  description: {
    margin: '8px 0 16px'
  }
};

class AboutView extends React.Component {
  render() {
    const {data, classes} = this.props;

    if (data.loading) {
      return <div>Loading</div>;
    }

    return (
      <div>
        <Typography type="title">
          About Operos
        </Typography>

        <Typography type="body1" className={classes.description}>
          The following is a list of the settings that apply to this
          installation of Operos.
        </Typography>
        
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Setting name</TableCell>
              <TableCell>Value</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            { Object.keys(data.cluster_info.settings).map(k => (
              <TableRow key={k}>
                <TableCell>{k}</TableCell>
                <TableCell>{data.cluster_info.settings[k]}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    );
  }
}

export default graphql(gql`
  query {
    cluster_info {
      settings
    }
  }
`)(withStyles(styles)(AboutView));
