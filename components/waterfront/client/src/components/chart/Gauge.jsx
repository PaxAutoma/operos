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
import {Doughnut} from 'react-chartjs-2';

import formatLabel from 'components/chart/format';
import {noval, levels} from 'components/chart/colors';
import withQueries from 'components/chart/withQueries';

const plugin = {
  id: 'gaugelabel',
  afterDatasetsDraw(chart, _easing, options) {
    const ctx = chart.ctx;

    let label = 'N/A';
    const value = chart.data.datasets[0].data[0];
    if (value !== undefined && value !== null) {
      label = `${formatLabel(value, options.format)}${options.suffix || ''}`;
    }

    ctx.save();
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillStyle = 'rgba(0, 0, 0, 0.6)';
    ctx.font = `${chart.innerRadius * 0.5}px Roboto`;
    ctx.fillText(label, chart.width / 2, chart.chartArea.top + (chart.chartArea.bottom - chart.chartArea.top) * 0.7);
    ctx.restore();
  }
};

class Gauge extends React.Component {
  static propTypes = {
    thresholds: PropTypes.arrayOf(PropTypes.number),
    colors: PropTypes.array,
    data: PropTypes.shape({
      datasets: PropTypes.arrayOf(PropTypes.shape({
        data: PropTypes.number
      }))
    }),
    max: PropTypes.number,
    title: PropTypes.string,
    suffix: PropTypes.string,
    className: PropTypes.string
  }

  render() {
    const {thresholds, colors, data, max, title, suffix, format, className} = this.props;
    let value = null;
    if (data.datasets.length > 0) {
      value = data.datasets[0].data;
    }

    const options = {
      rotation: -1.1 * Math.PI,
      circumference: 1.2 * Math.PI,
      cutoutPercentage: 60,
      maintainAspectRatio: false,
      tooltips: {
        enabled: false
      },
      events: [],
      title: {
        display: !!title,
        text: title,
        fontSize: 14
      },
      plugins: {
        gaugelabel: {
          suffix,
          format
        }
      }
    };

    const realMax = max || 100;
    const realThresholds = thresholds || [
      65, 90
    ];
    const realColors = colors || levels;

    let valueColor = realColors[0];
    for (let idx = 0; idx < realThresholds.length; idx++) {
      if (value >= realThresholds[idx]) {
        valueColor = realColors[idx+1] || noval;
        break;
      }
    }

    const chartData = {
      datasets: [{
        data: [value, realMax - value],
        backgroundColor: [valueColor, noval],
        label: 'Values'
      }],
    };

    return (
      <div className={className}>
        <Doughnut options={options}
                  data={chartData}
                  plugins={[plugin]} />
      </div>
    );
  }
}

export default withQueries(false)(Gauge);
