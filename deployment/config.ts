import * as pulumi from "@pulumi/pulumi";

const config = new pulumi.Config();

const emailMaskingSvcStack = new pulumi.StackReference(
    config.require("email-masking-svc-stack")
);
export const sendgridClientKey = config.requireSecret("sendgridClientKey");
export const rabbitmqUrl = emailMaskingSvcStack.requireOutput("rabbitmqUrl");