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
import Card, {CardContent, CardActions} from 'material-ui/Card';
import Typography from 'material-ui/Typography';
import Button from 'material-ui/Button';
import {withStyles} from 'material-ui/styles';

const styles = {
};

class CredentialsCard extends React.Component {
  onDownloadCredentials() {
    window.open(`${process.env.API_BASE_URL}/clientcert`, '_blank');
  }

  render() {
    const {classes, className} = this.props;

    return (
      <Card className={className}>
        <CardContent>
          <Typography type="headline" component="h2">
            API and kubectl access
          </Typography>
          <Typography component="p">
            Click the button below to generate and download TLS certificate credentials
            that can be used to control the Operos Kubernetes instance via
            the REST API or <code>kubectl</code>.
          </Typography>
        </CardContent>
        <CardActions>
          <Button dense color="primary" onClick={() => this.onDownloadCredentials()}>
            Download credentials
          </Button>
        </CardActions>
      </Card>
    );
  }
}

export default withStyles(styles)(CredentialsCard);
