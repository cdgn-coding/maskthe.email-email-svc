import * as kx from "@pulumi/kubernetesx";
import * as k8s from "@pulumi/kubernetes";
import { rabbitmqEndpoint, sendgridClientKey } from "./config";
import {createService} from "./utils";

const env = {
    "GO_ENVIRONMENT": "production",
    "SENDGRID_CLIENTKEY": sendgridClientKey,
    "RABBITMQ_URL": rabbitmqEndpoint,
}

const componentName = "email-svc";
const imageName = "email-svc:v1";

const pb = new kx.PodBuilder({
    containers: [
        {
            env,
            name: componentName,
            image: imageName,
            imagePullPolicy: "IfNotPresent",
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

export const appService = createService({
    name: componentName,
    serviceSpecs: {
        type: kx.types.ServiceType.ClusterIP,
        ports: [{
            protocol: "TCP",
            name: "web",
            port: 8080,
            targetPort: 8080
        }],
        selector: {
            app: componentName,
        },
    },
    metadata: {
        name: componentName,
    }
}, deployment);

export const appEndpoint = `${componentName}.default.svc.cluster.local`;

export const appIngress = new k8s.apiextensions.CustomResource(`${componentName}-ingress`, {
    apiVersion: "traefik.containo.us/v1alpha1",
    kind: "IngressRoute",
    metadata: {
        name: componentName,
        namespace: "default",
    },
    spec: {
        entryPoints: [
            "web",
        ],
        routes: [
            {
                match: "Path(`/api/emails`)",
                kind: "Rule",
                services: [
                    {
                        name: componentName,
                        port: 8080,
                    },
                ],
            },
        ],
    },
}, {
    dependsOn: appService
})