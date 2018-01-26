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
import _ from 'lodash';
import List, {ListItem, ListItemIcon, ListItemText} from 'material-ui/List';
import Icon from 'material-ui/Icon';
import Avatar from 'material-ui/Avatar';
import Collapse from 'material-ui/transitions/Collapse';
import Card, {CardContent, CardHeader} from 'material-ui/Card';
import Typography from 'material-ui/Typography';
import {withStyles} from 'material-ui/styles';
import {graphql, gql} from 'react-apollo';
import {withRouter} from 'react-router';

import withLoader from 'components/withLoader';

const CLASS_ICONS = {
  bridge: 'device_hub',
  bus: 'device_hub',
  communication: 'import_export',
  disk: 'save',
  display: 'desktop_windows',
  generic: 'memory',
  input: 'keyboard',
  memory: 'memory',
  multimedia: 'headset',
  network: 'import_export',
  power: 'battery_charging_full',
  printer: 'print',
  processor: 'memory',
  storage: 'device_hub',
  system: 'computer',
  volume: 'folder',
};

const styles = {
  title: {
    marginBottom: 16
  },
  card: {
    marginBottom: 16
  }
};


class NodeView extends React.Component {
  constructor() {
    super();
    this.state = {
      collapseState: {}
    };
  }

  handleCollapseToggle(hwItemId) {
    this.setState({
      collapseState: {
        ...this.state.collapseState,
        [hwItemId]: !this.state.collapseState[hwItemId]
      }
    });
  }

  renderHardwareItem(hwItem, level=0) {
    let primary = hwItem.description;
    if (!primary) {
      primary = _.capitalize(hwItem.class);
    }
    let secondary = [];
    if (hwItem.vendor) {
      secondary.push(hwItem.vendor);
    }
    if (hwItem.product) {
      secondary.push(hwItem.product);
    }

    const icon = CLASS_ICONS[hwItem.class] ? (
      <Icon>{CLASS_ICONS[hwItem.class]}</Icon>
    ) : null;

    const nestedItems = (hwItem.Nodes || []).map(item => (
      this.renderHardwareItem(item, level + 1)
    ));

    const result = [
      <ListItem
          button
          key={hwItem.id}
          onClick={() => this.handleCollapseToggle(hwItem.id)}
          style={{paddingLeft: 16 + 24 * level}}>
        <ListItemIcon><Icon>{icon}</Icon></ListItemIcon>
        <ListItemText
            primary={primary}
            secondary={secondary.join(' / ')}
        />
        { nestedItems.length > 0 &&
          <Icon>{this.state.collapseState[hwItem.id] ? 'expand_less' : 'expand_more'}</Icon>
        }
      </ListItem>
    ];

    if (nestedItems.length > 0) {
      result.push(
        <Collapse in={this.state.collapseState[hwItem.id]} key={`collapse-${hwItem.id}`}>
          {nestedItems}
        </Collapse>
      );
    }

    return result;
  }

  render() {
    const {data, classes} = this.props;
    const node = data.node;

    const hardware = JSON.parse(node.hardware_info);

    return (
      <div>
        <Card className={classes.card}>
          <CardHeader
            avatar={<Avatar><Icon>memory</Icon></Avatar>}
            title={'Node: ' + node.id}
            subheader={node.status == 'READY' ? 'Ready' : 'Not ready'}
          />
        </Card>
        { hardware &&
          <Card className={classes.card}>
            <CardContent>
              <Typography type="subheading" className={classes.title}>
                System hardware
              </Typography>
              <List dense>
                { this.renderHardwareItem(hardware.system) }
              </List>
            </CardContent>
          </Card>
        }
      </div>
    );
  }
}

export default withRouter(graphql(gql`
  query getSingleNode($nodeId: String!) {
    node(id: $nodeId) {
      id
      status
      ip
      pod_cidr
      hardware_info
    }    
  }
`, {
  options: ({match}) => {
    return {
      variables: {
        nodeId: match.params.nodeId
      }
    };
  }
})(withStyles(styles)(withLoader(NodeView))));
