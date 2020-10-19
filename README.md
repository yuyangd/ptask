# ptask

Associate a Route53 record set to the public IP of an ECS Fargate Task.

## Usage

### Run this binary inside a Fargate task.

```Dockerfile
RUN curl -sL -o /usr/local/bin/ptask \
    https://github.com/yuyangd/ptask/releases/download/v0.2/ptask \
 && chmod +x /usr/local/bin/ptask
ENTRYPOINT ["/usr/local/bin/ptask", "exec", "--"]
```

### Run as Fargate task sidecar container

Build the ptask image

```Dockerfile
FROM alpine:3 AS downloader

# install curl
RUN apk --no-cache add curl ca-certificates

# Download ptask
RUN curl -sL -o /ptask \
  https://github.com/yuyangd/ptask/releases/download/v0.2/ptask \
  && chmod +x /ptask

FROM scratch

COPY --from=downloader /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=downloader /ptask /

ENTRYPOINT [ "/ptask" ]
```

Example Task Definition in CFN

```yaml
  EcsTaskDefinition:
    Type: AWS::ECS::TaskDefinition
    Properties:
      ContainerDefinitions:
        - Name: ptask
          Image: <PTASK-IMAGE>
          Essential: false
          LogConfiguration:
            LogDriver: awslogs
            Options:
              awslogs-group:
                Ref: CloudWatchLogsGroup
              awslogs-region:
                Ref: AWS::Region
              awslogs-stream-prefix: fargate/bastion
          Environment:
            - Name: AWS_DEFAULT_REGION
              Value: ap-southeast-2
            - Name: HOSTHEADER
              Value: <PTASK.EXAMPLE.COM>
            - Name: HOSTZONE
              Value: <EXAMPLE.COM.>
```

### IAM Policy

```yaml
Policies:
  - PolicyName: ECSTaskWithDNS
    PolicyDocument:
      Statement:
        - Effect: Allow
          Action:
            - ecs:DescribeTasks
            - ec2:DescribeNetworkInterfaces
            - route53:ChangeResourceRecordSets
            - route53:GetChange
            - route53:ListHostedZones
            - route53:ListResourceRecordSets
            - cloudformation:CreateStack
            - cloudformation:DescribeStacks
            - cloudformation:UpdateStack
          Resource: "*"
```

### Environment Varibles in TaskDefinition

```
AWS_DEFAULT_REGION=ap-southeast-2
HOSTHEADER=ptask.example.com.
HOSTZONE=example.com.
```

## Go Build

```bash
# Run locally
GOOS=linux GOARCH=amd64 go build .

# Via docker
export DOCKER_BUILDKIT=1 ## Enable Buildkit
docker build -o disk .

# Build docker image
docker build -t "ptask:v1" .
```

## Run locally

```bash
# == Start mock

# Change the Cluster and Task ARN in main.go to target your AWS account
go run mock/main.go


# == Start ptask

# Ensure AWS credential or configuration setup locally
HOSTHEADER=ptask.example.com HOSTZONE=example.com. go run *.go

# == Expecting log messages

# 2020/10/13 11:12:13 arn:aws:ecs:ap-southeast-2:123456789012:cluster/my-ecs-cluster
# 2020/10/13 11:12:14 arn:aws:ecs:ap-southeast-2:123456789012:task/my-ecs-cluster/dfc8752c12344e17afee8696be98ak78
# 2020/10/13 11:12:15 Task ENI provisioned: eni-12345ea8e38cb281c
# 2020/10/13 11:12:16 Public IP: 14.125.12.4
# 2020/10/13 11:12:17 Create or update record set: ptask.example.com

```



