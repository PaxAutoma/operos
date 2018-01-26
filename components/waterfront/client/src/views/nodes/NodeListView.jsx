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
import {Link, withRouter} from 'react-router-dom';
import Table, {TableBody, TableRow, TableCell, TableHead, TableSortLabel} from 'material-ui/Table';
import {graphql, gql} from 'react-apollo';
import {withStyles} from 'material-ui/styles';
import queryString from 'query-string';

import withLoader from 'components/withLoader';

const styles = theme => ({
  empty: {
    textAlign: 'center',
    lineHeight: '5em',
    color: theme.palette.text.disabled
  }
});

const sortKeys = {
  id: node => node.id,
  ip: node => (
    node.ip.split('.').reverse().reduce((acc, val, idx) => {
      return acc + parseInt(val) * Math.pow(256, idx);
    }, 0)
  ),
  status: node => node.status,
  pod_cidr: node => node.pod_cidr,
};

class NodeListView extends React.Component {

  sortNodes(nodes, sortMode) {
    const key = sortKeys[sortMode];

    return nodes.slice().sort((a, b) => {
      a = key(a);
      b = key(b);
      if (a < b) {
        return -1;
      } else if (a > b) {
        return 1;
      }
      return 0;
    });
  }

  onSortChange(sortMode) {
    const {history, location} = this.props;
    const search = queryString.parse(location.search);

    history.push({
      ...location,
      search: queryString.stringify({
        ...search,
        sort: sortMode
      })
    });
  }

  render() {
    const {data, classes, location} = this.props;
    const search = queryString.parse(location.search);

    const sortMode = search['sort'] || 'ip';
    const sortedNodes = this.sortNodes(data.nodes, sortMode);

    if (sortedNodes.length > 0) {
      return (
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>
                <TableSortLabel
                    onClick={() => this.onSortChange('status')}
                    active={sortMode === 'status'}
                >Status</TableSortLabel>
              </TableCell>
              <TableCell>
                <TableSortLabel
                    onClick={() => this.onSortChange('id')}
                    active={sortMode === 'id'}
                >ID</TableSortLabel>
              </TableCell>
              <TableCell>
                <TableSortLabel
                    onClick={() => this.onSortChange('ip')}
                    active={sortMode === 'ip'}
                >IP address</TableSortLabel>
              </TableCell>
              <TableCell>
                <TableSortLabel
                    onClick={() => this.onSortChange('pod_cidr')}
                    active={sortMode === 'pod_cidr'}
                >Pod CIDR</TableSortLabel>
                </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            { sortedNodes.map(node => this.renderNode(node)) }
          </TableBody>
        </Table>
      );
    } else {
      return (
        <div className={classes.empty}>No nodes are currently available.</div>
      );
    }
  }

  renderNode(node) {
    return (
      <TableRow key={node.id}>
        <TableCell>{node.status == 'READY' ? 'Ready' : 'Not ready'}</TableCell>
        <TableCell><Link to={`/nodes/${node.id}`}>{node.id}</Link></TableCell>
        <TableCell>{node.ip}</TableCell>
        <TableCell>{node.pod_cidr}</TableCell>
      </TableRow>
    );
  }
}

export default graphql(gql`{
  nodes {
    id
    status
    ip
    pod_cidr
  }
}`)(withLoader(withStyles(styles)(withRouter(NodeListView))));
