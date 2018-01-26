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
import {withStyles} from 'material-ui/styles';
import Grid from 'material-ui/Grid';
import color from 'color';
import _ from 'lodash';

import withQueries from 'components/chart/withQueries';
import {palette} from 'components/chart/colors';
import formatLabel from 'components/chart/format';


const plugin = {
  id: 'verticalcursor',
  beforeInit(chart) {
    chart.verticalcursor = {};
  },
  afterDatasetsDraw(chart, _easing, _options) {
    if (chart.verticalcursor.x !== undefined) {
      const ctx = chart.ctx;
      ctx.save();

      ctx.beginPath();
      ctx.strokeStyle = 'rgba(0, 0, 0, 0.8)';
      ctx.setLineDash([4, 2]);
      ctx.moveTo(chart.verticalcursor.x, chart.chartArea.top);
      ctx.lineTo(chart.verticalcursor.x, chart.chartArea.bottom);
      ctx.stroke();
      ctx.closePath();

      ctx.restore();
    }
  },
  afterEvent(chart, event, _options) {
    if (event.type === 'mousemove') {
      if (event.x > chart.chartArea.left &&
          event.x < chart.chartArea.right) {
        chart.verticalcursor.x = event.x;
      } else {
        chart.verticalcursor.x = undefined;
      }
    } else if (event.type === 'mouseout') {
      chart.verticalcursor.x = undefined;
    }
  }
};

const styles = {
  container: {
    display: 'flex',
    height: '100%',
    alignItems: 'stretch'
  },
  legend: {
    fontSize: 12,
    marginTop: 24,
    padding: 0,
    width: '100%'
  },
  legendItem: {
    lineHeight: '16px'
  },
  legendBox: {
    height: '1em',
    width: '1em',
    display: 'inline-block',
    verticalAlign: 'middle',
    marginRight: 4
  },
  legendValue: {
    paddingLeft: 24,
    textAlign: 'right'
  }
};

export default (ChartType, chartTypeOptions, datasetDefaults) => {
  const chartClass = class extends React.Component {
    static propTypes = {
      title: PropTypes.string,
      data: PropTypes.shape({
        datasets: PropTypes.array
      }),
      format: PropTypes.oneOfType([
        PropTypes.object,
        PropTypes.arrayOf(PropTypes.object)
      ]),
      legend: PropTypes.oneOf(['none', 'below', 'right'])
    }

    componentDidMount() {
      this.forceUpdate();
    }

    datasetFormat(datasetIdx) {
      const {format} = this.props;
      return Array.isArray(format) ? format[datasetIdx] : format;
    }

    render() {
      const {title, data, classes, options, legend, className} = this.props;
      const chartOptions = _.merge({
        maintainAspectRatio: false,
        animation: false,
        scales: {
          xAxes: [{
            type: 'time',
            time: {
              unit: 'minute'
            }
          }],
          yAxes: [{
            afterFit: (axes) => {
              if (!axes.isHorizontal()) {
                axes.width = 80;
              }
            },
            ticks: {
              callback: (label) => formatLabel(label, this.datasetFormat(0)),
              beginAtZero: true
            }
          }]
        },
        tooltips: {
          mode: 'index',
          intersect: false,
          callbacks: {
            label: (tooltip, data) => {
              return `${data.datasets[tooltip.datasetIndex].label}: ${formatLabel(tooltip.yLabel, this.datasetFormat(tooltip.datasetIndex))}`;
            }
          }
        },
        title: {
          display: !!title,
          text: title,
          fontSize: 14
        },
        legend: {
          display: false
        }
      }, chartTypeOptions, options || {});

      const chartData = {
        datasets: data.datasets.map((ds, idx) => _.merge({
          backgroundColor: color(palette[idx % palette.length]).alpha(0.6).toString(),
          borderColor: color(palette[idx % palette.length]).alpha(0.6).toString()
        }, datasetDefaults, ds))
      };

      let chartBreaks;
      let legendBreaks;
      switch (legend) {
        case 'right', undefined:
          chartBreaks = {sm: 12, md: 8, lg: 9, xl: 10};
          legendBreaks = {sm: 12, md: 4, lg: 3, xl: 2};
          break;
        case 'bottom':
          chartBreaks = {xs: 12};
          legendBreaks = {xs: 12};
          break;
        default:
          chartBreaks = {xs: 12};
      }

      return (
        <div className={className}>
          <Grid container spacing={8} className={classes.container}>
            <Grid item {...chartBreaks}>
              <ChartType data={chartData}
                        options={chartOptions}
                        plugins={[plugin]}
                        ref="chart"
              />
            </Grid>
            { legendBreaks &&
              <Grid item {...legendBreaks}>
                {this.refs.chart && this.renderLegend(this.refs.chart.chart_instance)}
              </Grid>
            }
          </Grid>
        </div>
      );
    }

    renderLegend(chart) {
      const {classes} = this.props;

      return (
        <table className={classes.legend}>
          <thead>
            <tr className={classes.legendItem}>
              <th></th>
              <th className={classes.legendValue}>current</th>
            </tr>
          </thead>
          <tbody>
          { chart.data.datasets.map((ds, idx) => (
            <tr key={idx} className={classes.legendItem}>
              <td>
                <span style={{backgroundColor: ds.backgroundColor, border: `1px solid ${ds.borderColor}`}} className={classes.legendBox}></span>
                {ds.label}
              </td>
              <td className={classes.legendValue}>
                {formatLabel(ds.data[ds.data.length - 1].y, this.datasetFormat(idx))}
              </td>
            </tr>
          )) }
          </tbody>
        </table>
      );
    }
  };

  return withStyles(styles)(withQueries(true)(chartClass));
};
