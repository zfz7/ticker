import {Stage} from "./types";

export const AWS_ACCOUNT = "722763650436"
export const prod: Stage = {
    isProd: true,
    name: 'prod',
    region: 'us-west-2',
    account: AWS_ACCOUNT
}
export const stages: Stage[] = [prod];
