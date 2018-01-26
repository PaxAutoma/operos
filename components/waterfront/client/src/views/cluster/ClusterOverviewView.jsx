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
import Grid from 'material-ui/Grid';
import classnames from 'classnames';

import Gauge from 'components/chart/Gauge';
import SingleStat from 'components/chart/SingleStat';


const styles = {
  chart: {
    border: '1px solid rgba(0, 0, 0, 0.05)',
    backgroundColor: 'rgba(0, 0, 0, 0.01)',
    padding: 8
  },
  gauge: {
    height: 200,
  },
  full: {
    height: 300,
  },
  singleStat: {
    height: 100
  }
};

class ClusterOverviewView extends React.Component {
  render() {
    const {classes} = this.props;

    return (
      <Grid container spacing={8}>
        <Grid item xs={12} sm={6} md={4}>
          <Gauge query={'(sum(node_memory_MemTotal) - sum(node_memory_MemFree+node_memory_Buffers+node_memory_Cached) ) / sum(node_memory_MemTotal) * 100'}
                  title="Cluster memory usage"
                  suffix="%"
                  format={{decimals: 0}}
                  className={classnames(classes.chart, classes.gauge)}
          />
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <Gauge query={'sum(rate(node_cpu{mode!="idle"}[1m])) / count(node_cpu{mode="idle"}) * 100'}
                  title="Cluster CPU usage"
                  suffix="%"
                  format={{decimals: 0}}
                  className={classnames(classes.chart, classes.gauge)}
          />
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <Gauge query={'(sum(node_filesystem_size{device="rootfs"}) - sum(node_filesystem_free{device="rootfs"}) ) / sum(node_filesystem_size{device="rootfs"}) * 100'}
                  title="Cluster filesystem usage"
                  suffix="%"
                  format={{decimals: 0}}
                  className={classnames(classes.chart, classes.gauge)}
          />
        </Grid>
        <Grid item xs={6} md={2}>
          <SingleStat query={'count(node_boot_time)'}
                      title="Nodes"
                      className={classnames(classes.chart, classes.singleStat)}
          />
        </Grid>
        <Grid item xs={6} md={2}>
          <SingleStat query={'sum(kubelet_running_pod_count)'}
                      title="Running pods"
                      className={classnames(classes.chart, classes.singleStat)}
          />
        </Grid>
        <Grid item xs={6} md={2}>
          <SingleStat query={'sum(kubelet_running_container_count)'}
                      title="Running containers"
                      className={classnames(classes.chart, classes.singleStat)}
          />
        </Grid>
        <Grid item xs={6} md={2}>
          <SingleStat query={'count(node_cpu{mode="idle"})'}
                      title="Total CPU cores"
                      className={classnames(classes.chart, classes.singleStat)}
          />
        </Grid>
        <Grid item xs={6} md={2}>
          <SingleStat query={'sum(node_memory_MemTotal)'}
                      title="Total RAM"
                      format={{unit: 'B', scale: 'binary', decimals: 1}}
                      className={classnames(classes.chart, classes.singleStat)}
          />
        </Grid>
        <Grid item xs={6} md={2}>
          <SingleStat query={'sum(node_filesystem_size{mountpoint="/",fstype!="rootfs"})'}
                      title="Total storage"
                      format={{unit: 'B', scale: 'binary', decimals: 1}}
                      className={classnames(classes.chart, classes.singleStat)}
          />
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(ClusterOverviewView);
