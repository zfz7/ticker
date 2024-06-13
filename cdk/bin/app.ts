#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { LambdaStack } from '../stacks/lambdaStack';
import {stages} from "./config";
import {StackProps} from "aws-cdk-lib";

const app = new cdk.App();

stages.forEach(stage => {
    const stackProps: StackProps = {
        env: {account: stage.account, region: stage.region}
    }
    new LambdaStack(app, 'Ticker-Lambda-Stack', {
        ...stackProps
    });

})
