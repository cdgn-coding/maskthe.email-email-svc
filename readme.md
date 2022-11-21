# Email Service

This service is a bridge from SendGrid to the rest of the system.
It subscribes to a topic and sends emails to the recipients,
and receives emails using the SendGrid Webhook and publishes them to a topic.

![Architecture](/docs/diagrams/email-svc.png)

## How to run

## Prerequisites

- [Docker](https://docs.docker.com/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/doc/install)
- [Pulumi](https://www.pulumi.com/docs/get-started/install/)

### Minikube

If you want to run the service locally, you can use [Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/).

Build the docker image onto the cluster.

```bash
minikube image build -t email-svc .
```

Use pulumi to deploy the service.

```bash
pulumi up -y
```

## License

    Copyright 2022 Carlos David Gonzalez Nexans
    
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    
      http://www.apache.org/licenses/LICENSE-2.0
    
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.