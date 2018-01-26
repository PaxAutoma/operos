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

import AppBar from 'material-ui/AppBar';
import Drawer from 'material-ui/Drawer';
import Divider from 'material-ui/Divider';
import Icon from 'material-ui/Icon';
import IconButton from 'material-ui/IconButton';
import List, {ListItem, ListItemIcon, ListItemText} from 'material-ui/List';
import Toolbar from 'material-ui/Toolbar';
import Typography from 'material-ui/Typography';
import { withStyles } from 'material-ui/styles';
import Popover from 'material-ui/Popover';
import classNames from 'classnames';
import {Redirect, Route, Switch} from 'react-router-dom';
import {withRouter} from 'react-router';

import UserCard from 'components/UserCard';

const drawerWidth = 240;
const smallWidth = 56;

const styles = theme => ({
  root: {
    height: '100%',
    width: '100%',
    position: 'relative',
    display: 'flex'
  },
  appBar: {
    position: 'absolute',
    zIndex: theme.zIndex.navDrawer + 1,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    })
  },
  appBarShift: {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  menuButton: {
    marginLeft: 4,
    marginRight: 16,
  },
  hide: {
    display: 'none',
  },
  drawerPaper: {
    position: 'relative',
    height: '100%',
    overflow: 'hidden',
    width: drawerWidth,
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  drawerPaperClose: {
    width: smallWidth,
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
  },
  drawerInner: {
    // Make the items inside not wrap when transitioning:
    width: drawerWidth,
  },
  drawerHeader: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
    padding: '0 8px',
    height: 56,
    [theme.breakpoints.up('sm')]: {
      height: 64,
    },
  },
  content: {
    marginTop: 64,
    flexGrow: 1,
    padding: 16,
    overflow: 'auto'
  },
  activeMenuItem: {
    backgroundColor: theme.palette.text.divider
  },
  title: {
    flex: 1,
    display: 'flex',
    alignItems: 'center'
  },
  accountMenuButton: {
    marginRight: 4
  },
  list: {
    paddingTop: 0
  },
  logo: {
    marginRight: theme.spacing.unit
  }
});

class Chrome extends React.Component {

  constructor() {
    super();
    this.state = {
      open: false,
      accountMenuAnchor: null,
      accountMenuOpen: false
    };
  }

  handleDrawerOpen() {
    this.setState({ open: true });
  }

  handleDrawerClose() {
    this.setState({ open: false });
  }

  handleMenuAction(link) {
    this.props.history.push(link);
  }

  handleAccountMenuOpen(event) {
    this.setState({
      accountMenuAnchor: event.target,
      accountMenuOpen: true
    });
  }

  handleAccountMenuClose() {
    this.setState({
      accountMenuAnchor: null,
      accountMenuOpen: false
    });
  }

  openKubeDashboard() {
    window.open(process.env.KUBE_DASHBOARD_URL, '_blank');
  }
  
  renderMain() {
    return (
      <Switch>
        <Route path="/" exact render={() => (
          <Redirect to="/cluster" />
        )} />
        { this.props.views.map(view => (
          <Route
              key={view.path}
              path={view.path}
              exact={view.exact}
              render={() => {
                return <view.component />;
              }}
          />
        )) }
      </Switch>
    );
  }

  render() {
    const classes = this.props.classes;

    return (
      <div className={classes.root}>
        <AppBar className={classNames(classes.appBar, this.state.open && classes.appBarShift)} color="default">
          <Toolbar disableGutters={!this.state.open}>
            <IconButton
                color="contrast"
                onClick={() => this.handleDrawerOpen()}
                className={classNames(classes.menuButton, this.state.open && classes.hide)}
            >
              <Icon>menu</Icon>
            </IconButton>
            <div className={classes.title}>
              <img src="/static/operos-logo.svg" alt="Operos logo" className={classes.logo} />
            </div>
            <IconButton
                color="contrast"
                aria-label="Account"
                onClick={this.handleAccountMenuOpen.bind(this)}
                className={classes.accountMenuButton}
            >
              <Icon>account_box</Icon>
            </IconButton>
            <Popover
                anchorEl={this.state.accountMenuAnchor}
                open={this.state.accountMenuOpen}
                onRequestClose={this.handleAccountMenuClose.bind(this)}
                anchorOrigin={{vertical: 'bottom', horizontal: 'left'}}
            >
              <UserCard />
            </Popover>
          </Toolbar>
        </AppBar>
        <Drawer
          type="permanent"
          classes={{
            paper: classNames(classes.drawerPaper, !this.state.open && classes.drawerPaperClose),
          }}
          open={this.state.open}
        >
          <div className={classes.drawerInner}>
            <div className={classes.drawerHeader}>
              <IconButton onClick={() => this.handleDrawerClose()}>
                <Icon>chevron_left</Icon>
              </IconButton>
            </div>
            <Divider />
            <List className={classes.list}>
              { this.props.menu.map(item => (
                <Route
                    key={item.name}
                    path={item.link}
                    exact
                    children={({match}) => (
                  <ListItem
                      button
                      onTouchTap={() => this.handleMenuAction(item.link)}
                      className={match ? classes.activeMenuItem : null}
                  >
                    <ListItemIcon><Icon>{item.icon}</Icon></ListItemIcon>
                    <ListItemText primary={item.name} />
                  </ListItem>
                )} />
              ))}
            </List>
            <Divider />
            <List className={classes.list}>
              <ListItem button onTouchTap={() => this.openKubeDashboard()}>
                <ListItemIcon><img src="/static/kube-logo.svg" /></ListItemIcon>
                <ListItemText primary="Kubernetes" />
              </ListItem>
            </List>
          </div>
        </Drawer>
        <main className={classes.content}>
          { this.renderMain() }
        </main>
      </div>
    );
  }
}

export default withRouter(withStyles(styles, {withTheme: true})(Chrome));
