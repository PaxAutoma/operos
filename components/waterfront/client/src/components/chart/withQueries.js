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
import _ from 'lodash';
import moment from 'moment';


const queryType = PropTypes.oneOfType([
  PropTypes.shape({
    query: PropTypes.string.isRequired,
    options: PropTypes.object
  }),
  PropTypes.string
]);

export default (range=false) => (Chart) => {
  return class extends React.Component {
    static propTypes = {
      query: PropTypes.oneOfType([
        queryType,
        PropTypes.arrayOf(queryType)
      ]).isRequired,
      refresh: PropTypes.number
    }

    static contextTypes = {
      apiClient: PropTypes.object
    }

    constructor() {
      super();
      this.state = {};
      this.timeout = null;
    }

    componentWillMount() {
      this.doQueries();
    }

    componentWillUnmount() {
      if (this.timeout) {
        clearTimeout(this.timeout);
        this.timeout = null;
      }
    }

    doQueries() {
      const queries = Array.isArray(this.props.query) ? this.props.query : [this.props.query];
      const promises = [];
      for (const q of queries) {
        const query = typeof q === 'string' ? q : q.query;
        if (range) {
          promises.push(this.context.apiClient.metricsQueryRange(query, q.step || 60));
        } else {
          promises.push(this.context.apiClient.metricsQuery(query, q.step || 60));
        }
      }

      Promise.all(promises).then(res => {
        let data = res.map((r, queryIdx) => {
          const queryOptions = queries[queryIdx].options || {};

          return r.entity.data.result.map(resultData => {
            let values;
            switch (r.entity.data.resultType) {
              case 'matrix':
                values = parseTimeseries(resultData.values, queries[queryIdx].step || 60);
                break;
              case 'vector':
                values = parseFloat(resultData.value[1]);
                break;
            }
  
            const dataset = {
              label: resultData.metric['instance'],
              data: values
            };

            return _.merge(dataset, queryOptions);
          });
        });
        
        // Flatten the list of lists
        data = data.reduce((a, b) => a.concat(b), []);

        this.setState({
          data,
          error: null
        });
      })
      .catch(err => {
        this.setState({
          error: err
        });
      })
      .then(() => {
        if (this.props.refresh !== 0) {
          const interval = this.props.refresh === undefined ? 10000 : this.props.refresh;
          this.timeout = setTimeout(this.doQueries.bind(this), interval);
        }
      });
    }

    render() {
      const {queries, className, ...otherProps} = this.props;
      
      if (this.state.error) {
        return <div className={className} style={{textAlign: 'center', verticalAlign: 'middle'}}>Failed to load chart data</div>;
      } else if (this.state.data) {
        return <Chart data={{datasets: this.state.data}} className={className} {...otherProps} />;
      } else {
        return <div className={className} style={{textAlign: 'center', verticalAlign: 'middle'}}></div>;
      }
    }
  };
};

const parseTimeseries = (values, step) => {
  let lastT = null;
  const result = [];

  values.forEach(([t, v]) => {
    // Fill in blank values with nulls
    if (lastT !== null && t > lastT + step) {
      for (let spanT = lastT + step; spanT < t; spanT += step) {
        result.push({t: moment.unix(spanT), y: null});
      }
    }

    result.push({t: moment.unix(t), y: parseFloat(v)});
    lastT = t;
  });

  return result;
};
