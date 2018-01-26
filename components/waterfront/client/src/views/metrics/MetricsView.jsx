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
import Grid from 'material-ui/Grid';
import {withStyles} from 'material-ui/styles';

import formatLabel from 'components/chart/format';
import LineChart from 'components/chart/Line';
import BarChart from 'components/chart/Bar';

const styles = {
  chart: {
    border: '1px solid rgba(0, 0, 0, 0.05)',
    backgroundColor: 'rgba(0, 0, 0, 0.01)',
    padding: 8,
    height: 300
  },
};

class MetricsView extends React.Component {
  render() {
    const {classes} = this.props;

    return (
      <Grid container spacing={8}>
        <Grid item xs={12} md={6}>
          <LineChart query={[{
              query: 'avg(sum(kubelet_running_pod_count) by (instance))',
              options: {
                label: 'Avg',
                fill: false,
                borderWidth: 1,
                borderColor: 'rgba(10, 67, 124, 1)'
              }
            }, {
              query: 'min(sum(kubelet_running_pod_count) by (instance))',
              options: {
                label: 'Min',
                fill: 0,
                backgroundColor: 'rgba(10, 67, 124, 0.3)',
                borderColor: 'rgba(0, 0, 0, 0)'
              }
            }, {
              query: 'max(sum(kubelet_running_pod_count) by (instance))',
              options: {
                label: 'Max',
                fill: 0,
                backgroundColor: 'rgba(10, 67, 124, 0.3)',
                borderColor: 'rgba(0, 0, 0, 0)'
              }
            }]}
            title="Pods per node"
            format={{prefix: '', decimals: 0}}
            legend="none"
            options={{
              scales: {
                yAxes: [{
                  afterFit: () => {}
                }]
              }
            }}
            className={classes.chart}
          />
        </Grid>
        <Grid item xs={12} md={6}>
          <LineChart query={[{
              query: 'avg(sum(kubelet_running_container_count) by (instance))',
              options: {
                label: 'Avg',
                fill: false,
                borderWidth: 1,
                borderColor: 'rgba(10, 67, 124, 1)'
              }
            }, {
              query: 'min(sum(kubelet_running_container_count) by (instance))',
              options: {
                label: 'Min',
                fill: 0,
                backgroundColor: 'rgba(10, 67, 124, 0.3)',
                borderColor: 'rgba(0, 0, 0, 0)'
              }
            }, {
              query: 'max(sum(kubelet_running_container_count) by (instance))',
              options: {
                label: 'Max',
                fill: 0,
                backgroundColor: 'rgba(10, 67, 124, 0.3)',
                borderColor: 'rgba(0, 0, 0, 0)'
              }
            }]}
            title="Containers per node"
            format={{prefix: '', decimals: 0}}
            legend="none"
            options={{
              scales: {
                yAxes: [{
                  afterFit: () => {}
                }]
              }
            }}
            className={classes.chart}
          />
        </Grid>

        <Grid item xs={12}>
          <BarChart query={{
              query: 'count(node_boot_time)',
              options: {
                label: '# of nodes',
                backgroundColor: 'rgb(10, 67, 124)'
              }
            }}
            title="Nodes"
            format={{prefix: '', decimals: 1}}
            className={classes.chart}
          />
        </Grid>
        <Grid item xs={12}>
          <LineChart query={[{
              query:'sum(node_load1)',
              step: 10,
              options: {
                label: '1 min',
                fill: 'none',
                showLine: false,
                pointRadius: 1
              }
            }, {
              query: 'sum(node_load5)',
              step: 10,
              options: {
                label: '5 min',
                fill: 'none',
                showLine: false,
                pointRadius: 1
              }
            }, {
              query: 'sum(node_load15)',
              step: 10,
              options: {
                label: '15 min',
                fill: 'none',
                showLine: false,
                pointRadius: 1
              }
            }]}
            title="Load"
            format={{prefix: '', decimals: 2}}
            className={classes.chart}
          />
        </Grid>
        <Grid item xs={12}>
          <LineChart query={[{
              query:'sum(rate(node_cpu{mode!="idle"}[5m])) ',
              step: 10,
              options: {
                label: 'Used',
                borderColor: 'rgb(10, 67, 124)',
                backgroundColor: 'rgba(10, 67, 124, 0.6)',
                borderWidth: 1
              }
            }, {
              query: 'count(node_cpu{mode="idle"})',
              step: 10,
              options: {
                label: 'Available',
                fill: 'none',
                borderColor: 'rgb(137, 15, 2)',
                backgroundColor: 'rgba(137, 15, 2, 0.6)',
                borderWidth: 1
              }
            }]}
            title="CPU"
            format={{prefix: '', decimals: 2}}
            className={classes.chart}
          />
        </Grid>

        <Grid item xs={12}>
          <LineChart query={[{
              query:'sum(node_memory_MemTotal)-sum(node_memory_MemAvailable)',
              step: 10,
              options: {
                label: 'Used',
                borderColor: 'rgb(10, 67, 124)',
                backgroundColor: 'rgba(10, 67, 124, 0.6)',
                borderWidth: 1
              }
            }, {
              query: 'sum(node_memory_MemTotal)',
              step: 10,
              options: {
                label: 'Available',
                fill: 'none',
                borderColor: 'rgb(137, 15, 2)',
                backgroundColor: 'rgba(137, 15, 2, 0.6)',
                borderWidth: 1
              }
            }]}
            title="Memory"
            format={{unit: 'B', decimals: 2, scale: 'binary'}}
            className={classes.chart}
          />
        </Grid>

        <Grid item xs={12}>
          <LineChart query={[{
              query:'sum(max(node_filesystem_size{device=~"/dev/.da."}) by (device)) - sum(max(node_filesystem_avail{device=~"/dev/.da."}) by (device))',
              step: 10,
              options: {
                label: 'Used',
                borderColor: 'rgb(10, 67, 124)',
                backgroundColor: 'rgba(10, 67, 124, 0.6)',
                borderWidth: 1
              }
            }, {
              query: 'sum(max(node_filesystem_size{device=~"/dev/.da."}) by (device))',
              step: 10,
              options: {
                label: 'Available',
                fill: 'none',
                borderColor: 'rgb(137, 15, 2)',
                backgroundColor: 'rgba(137, 15, 2, 0.6)',
                borderWidth: 1
              }
            }]}
            title="Storage"
            format={{unit: 'B', decimals: 2, scale: 'binary'}}
            className={classes.chart}
          />
        </Grid>

        <Grid item xs={12}>
          <LineChart query={[{
              query:'sum(rate(node_network_receive_bytes{device=~"(eth|en).*"}[5m]))',
              step: 10,
              options: {
                label: 'Received',
                borderWidth: 1
              }
            }, {
              query: '-sum(rate(node_network_transmit_bytes{device=~"(eth|en).*"}[5m]))',
              step: 10,
              options: {
                label: 'Sent',
                borderWidth: 1
              }
            }]}
            title="Network"
            format={{unit: 'B', decimals: 2, scale: 'binary'}}
            className={classes.chart}
          />
        </Grid>

        <Grid item xs={12}>
          <LineChart query={[{
              query:'sum(rate(node_disk_writes_completed[1m]))',
              step: 10,
              options: {
                label: 'Writes',
                fill: 'none',
                showLine: false,
                pointRadius: 1
              }
            }, {
              query: 'sum(rate(node_disk_bytes_written[1m]))',
              step: 10,
              options: {
                label: 'Bytes written',
                borderColor: 'rgba(0, 0, 0, 0)',
                backgroundColor: 'rgba(10, 67, 124, 0.6)',
                fill: 'origin',
                yAxisID: 'secondary'
              }
            }]}
            format={[
              {unit: 'ops'},
              {unit: 'Bps', decimals: 1}
            ]}
            options={{
              scales: {
                yAxes: [{}, {
                  id: 'secondary',
                  position: 'right',
                  gridLines: {
                    display: false
                  },
                  ticks: {
                    callback: label => formatLabel(label, {unit: 'Bps', decimals: 1})
                  }
                }]
              }
            }}
            title="Write I/O"
            className={classes.chart}
          />
        </Grid>

        <Grid item xs={12}>
          <LineChart query={[{
              query:'sum(rate(node_disk_reads_completed[1m]))',
              step: 10,
              options: {
                label: 'Reads',
                fill: 'none',
                showLine: false,
                pointRadius: 1
              }
            }, {
              query: 'sum(rate(node_disk_bytes_read[1m]))',
              step: 10,
              options: {
                label: 'Bytes read',
                borderColor: 'rgba(0, 0, 0, 0)',
                backgroundColor: 'rgba(10, 67, 124, 0.6)',
                fill: 'origin',
                yAxisID: 'secondary'
              }
            }]}
            format={[
              {unit: 'ops'},
              {unit: 'Bps', decimals: 1}
            ]}
            options={{
              scales: {
                yAxes: [{}, {
                  id: 'secondary',
                  position: 'right',
                  gridLines: {
                    display: false
                  },
                  ticks: {
                    callback: label => formatLabel(label, {unit: 'Bps', decimals: 1})
                  }
                }]
              }
            }}
            title="Read I/O"
            className={classes.chart}
          />
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(MetricsView);
