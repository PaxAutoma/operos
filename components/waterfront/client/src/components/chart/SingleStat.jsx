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
import classNames from 'classnames';

import formatLabel from 'components/chart/format';
import withQueries from 'components/chart/withQueries';

const styles = {
  title: {
    fontSize: 14,
    marginBottom: 18
  },
  value: {
    fontSize: 30
  },
  container: {
    textAlign: 'center',
    padding: 8,
    color: 'rgba(0, 0, 0, 0.6)'
  }
};

class SingleStat extends React.Component {
  render() {
    const {title, data, classes, format, className} = this.props;

    return (
      <div className={classNames(classes.container, className)}>
        <h3 className={classes.title}>{title}</h3>
        <div className={classes.value}>{formatLabel(data.datasets[0].data, format)}</div>
      </div>
    );
  }
}

export default withStyles(styles)(withQueries(false)(SingleStat));
