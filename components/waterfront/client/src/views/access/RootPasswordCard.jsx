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
import Card, {CardContent, CardActions} from 'material-ui/Card';
import Typography from 'material-ui/Typography';
import Button from 'material-ui/Button';
import IconButton from 'material-ui/IconButton';
import Icon from 'material-ui/Icon';
import TextField from 'material-ui/TextField';
import Snackbar from 'material-ui/Snackbar';
import {CircularProgress} from 'material-ui/Progress';
import {withStyles} from 'material-ui/styles';

const styles = theme => ({
  fields: {
    display: 'flex',
    flexDirection: 'column'
  },
  textbox: {
    maxWidth: 300,
    marginTop: theme.spacing.unit
  },
  progress: {
    marginLeft: theme.spacing.unit * 2
  }
});

class RootPasswordCard extends React.Component {
  constructor() {
    super();
    this.state = {
      password1: '',
      password2: '',
      disabled: true,
      working: false,
      message: null
    }
  }

  static contextTypes = {
    apiClient: PropTypes.object
  }

  onChange(event) {
    const chg = {
      [event.target.id]: event.target.value
    };
    const tmp = Object.assign({}, this.state, chg);

    chg.disabled = tmp.password1 === "" || tmp.password1 !== tmp.password2;
    this.setState(chg);
  }

  onSubmit(event) {
    event.preventDefault();

    this.setState({
      working: true
    });

    this.context.apiClient.setRootPassword(this.state.password1).then(res => {
      this.setState({
        password1: '',
        password2: '',
        disabled: true,
        working: false,
        message: 'Password set successfully'
      });
    }).catch(err => {
      this.setState({
        working: false,
        message: 'Error setting password: ' + err.toString()
      });
    });
  }

  handleCloseSnack() {
    this.setState({
      message: null
    });
  }
  
  render() {
    const {classes, className} = this.props;

    return (
      <Card className={className}>
        <form noValidate autoComplete="off">
          <CardContent>
            <Typography type="headline" component="h2">
              Root password
            </Typography>
            <Typography component="p">
              In this preview release, the root account allows login to this UI as well
              as to the machine console on the controller and all worker machines. It
              can be changed here.
            </Typography>

            <div className={classes.fields}>
              <TextField
                  id="password1"
                  label="Password"
                  type="password"
                  className={classes.textbox}
                  value={this.state.password1}
                  onChange={this.onChange.bind(this)}
              />
              <TextField
                  id="password2"
                  label="Confirm password"
                  type="password"
                  className={classes.textbox}
                  value={this.state.password2}
                  onChange={this.onChange.bind(this)}
              />
            </div>
          </CardContent>
          <CardActions>
            { this.state.working
              ? <CircularProgress className={classes.progress} size={24} />
              : <Button
                    dense
                    color="primary"
                    type="submit"
                    onClick={this.onSubmit.bind(this)}
                    disabled={this.state.disabled}
                >
                  Change password
                </Button>
            }
          </CardActions>
          <Snackbar
              open={!!this.state.message}
              onRequestClose={this.handleCloseSnack.bind(this)}
              message={<span>{this.state.message}</span>}
              action={
                <IconButton
                  key="close"
                  color="inherit"
                  className={classes.close}
                  onClick={this.handleCloseSnack.bind(this)}
                >
                  <Icon>close</Icon>
                </IconButton>
              }
          />
        </form>
      </Card>
    );
  }
}

export default withStyles(styles)(RootPasswordCard);
