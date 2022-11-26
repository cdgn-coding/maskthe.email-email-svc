import * as pulumi from "@pulumi/pulumi";

const config = new pulumi.Config();

const clusterSetup = new pulumi.StackReference(
    config.require("cluster-setup")
);

export const sendgridClientKey = config.requireSecret("sendgridClientKey");
export const rabbitmqEndpoint = clusterSetup.requireOutput("rabbitmqEndpoint");