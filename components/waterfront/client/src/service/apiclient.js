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

import rest from 'rest';
import mime from 'rest/interceptor/mime';
import params from 'rest/interceptor/params';
import errorCode from 'rest/interceptor/errorCode';
import moment from 'moment';

export default class ApiClient {
  constructor(baseUrl) {
    this.baseUrl = baseUrl;
    this.client = rest.wrap(mime).wrap(errorCode).wrap(params);
  }

  doRequest(req) {
    return this
      .client(req)
      .catch(err => {
        if (err.error) {
          // Connection errors look like this
          throw err.error;
        }
        if (err.status.code) {
          // Errors from the server look like this
          throw err.status.text;
        }
        throw err;
      });
  }

  get(resource, params) {
    return this.doRequest({
      path: `${this.baseUrl}/${resource}`,
      mixin: {
        withCredentials: true
      },
      params
    })
  }

  post(resource, entity) {
    const req = {
      path: `${this.baseUrl}/${resource}`,
      method: 'POST',
      mixin: {
        withCredentials: true
      }
    };

    if (entity) {
      req.entity = JSON.stringify(entity)
    }

    return this.doRequest(req);
  }

  getNodes() {
    return this.get('nodes').then(res => res.entity.nodes || []);
  }

  getNode(nodeId) {
    return this.get(`nodes/${nodeId}`).then(res => res.entity.node);
  }

  getClusterInfo() {
    return this.get('cluster_info').then(res => res.entity);
  }

  metricsQuery(query) {
    return this.get('metrics/query', {
      query
    });
  }

  metricsQueryRange(query, step) {
    return this.get('metrics/query_range', {
      query,
      start: moment().subtract(1, 'hours').unix(),
      end: moment().unix(),
      step
    });
  }

  getLoginInfo() {
    return this.get('login').then(res => res.entity);
  }

  login(username, password) {
    return this.post('login', {
      username, password
    }).then(res => res.entity);
  }

  logout() {
    return this.post('logout').then(res => res.entity);
  }

  setRootPassword(password) {
    return this.post('rootpass', {
      password
    });
  }
}
