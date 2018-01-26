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
import {graphql, gql, compose} from 'react-apollo';

import Avatar from 'material-ui/Avatar';
import Button from 'material-ui/Button';
import Divider from 'material-ui/Divider';
import Icon from 'material-ui/Icon';
import Card, {CardHeader, CardContent, CardActions} from 'material-ui/Card';
import {CircularProgress} from 'material-ui/Progress';
import {withStyles} from 'material-ui/styles';

const styles = {
  card: {
    minWidth: 200
  }
};

class UserCard extends React.Component {
  constructor() {
    super();
    this.state = {
      loggingOut: false
    }
  }

  handleLogout() {
    this.setState({
      loggingOut: true
    })

    this.props.mutate().catch(err => {
      this.setState({
        loggingOut: false
      })  
    })
  }

  render() {
    const {classes, data: {login_info: {user}}} = this.props;

    return (
      <Card elevation={0} className={classes.card}>
        <CardHeader
            avatar={<Avatar><Icon>person</Icon></Avatar>}
            title={user.username}
        />
        <Divider />
        <CardActions>
          { this.state.loggingOut
            ? <CircularProgress size={24} />
            : <Button dense raised color="accent" onClick={() => this.handleLogout()}>
                Log out
              </Button>
          }
        </CardActions>
      </Card>
    );
  }
}

export default withStyles(styles)(compose(
  graphql(gql`
    query {
      login_info {
        user {
          username
        }
      }
    }
  `),
  graphql(gql`
    mutation logout {
      logout {
        logged_in
      }
    }  
  `)
)(UserCard));
