import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { Architecture, Code, Function, Runtime } from 'aws-cdk-lib/aws-lambda';
import { Rule, RuleTargetInput, Schedule } from 'aws-cdk-lib/aws-events';
import { LambdaFunction } from 'aws-cdk-lib/aws-events-targets';

export class LambdaStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const lambda = new Function(this, `Ticker-Lambda`, {
      code: Code.fromAsset('../lambda/lambdaFunction.zip', { deployTime: true }),
      handler: 'bootstrap',
      runtime: Runtime.PROVIDED_AL2,
      architecture: Architecture.ARM_64,
      environment: {
        FINNHUB_API_KEY: process.env.FINNHUB_API_KEY!,
        PUSHOVER_APP_KEY: process.env.PUSHOVER_APP_KEY!,
        PUSHOVER_RECIPIENT: process.env.PUSHOVER_RECIPIENT!,
      },
    });

    const eventRule = new Rule(this, 'scheduleRule', {
      schedule: Schedule.expression('cron(5 15-20 ? * MON-FRI *)'),
    });
    eventRule.addTarget(new LambdaFunction(lambda, { event: RuleTargetInput.fromObject({}) }));
  }
}
