import * as kx from "@pulumi/kubernetesx";
import { rabbitmqUrl, sendgridClientKey } from "./config";

const env = {
    "GO_ENVIRONMENT": "production",
    "SENDGRID_CLIENTKEY": sendgridClientKey,
    "RABBITMQ_URL": rabbitmqUrl,
}

const componentName = "email-svc";
const imageName = "email-svc";

const pb = new kx.PodBuilder({
    containers: [
        {
            env,
            name: componentName,
            image: imageName,
            imagePullPolicy: "Never",
            resources: { requests: { cpu: "128m", memory: "256Mi" } },
            ports: { http: 8080 },
            livenessProbe: {
                httpGet: {
                    path: "/health",
                    port: 8080,

                },
            },
        },
    ],
});

const deployment = new kx.Deployment(componentName, {
    spec: pb.asDeploymentSpec({ replicas: 1 }),
});

export const appService = deployment.createService({
    type: kx.types.ServiceType.ClusterIP,
    ports: [{ port: 8080, targetPort: 8080 }],
});
