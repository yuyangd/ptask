# ptask

Associate a record set to the public IP of ECS Fargate Task.

## Usage

Run this binary inside a Fargate task.

```Dockerfile
RUN curl -sL -o /usr/local/bin/ptask \
    https://github.com/yuyangd/ptask/releases/download/v0.1/ptask \
 && chmod +x /usr/local/bin/ptask
ENTRYPOINT ["/usr/local/bin/ptask", "exec", "--"]
```

### IAM Policy

```json

```

### Environment Varibles in TaskDefinition

```
AWS_DEFAULT_REGION=ap-southeast-2
HOSTHEADER=ptask.example.com
HOSTZONEID=GFD977DEOUCCN
```

## Go Build

```bash
GOOS=linux GOARCH=amd64 go build .
```

## Run locally

```bash
# == Start mock

# Change the Cluster and Task ARN in main.go to target your AWS account
go run mock/main.go


# == Start ptask

# Ensure AWS credential or configuration setup locally
HOSTHEADER=ptask.example.com HOSTZONEID=GFD977DEOUCCN go run *.go

# == Expecting log messages

# 2020/10/13 11:12:13 arn:aws:ecs:ap-southeast-2:123456789012:cluster/my-ecs-cluster
# 2020/10/13 11:12:14 arn:aws:ecs:ap-southeast-2:123456789012:task/my-ecs-cluster/dfc8752c12344e17afee8696be98ak78
# 2020/10/13 11:12:15 Task ENI provisioned: eni-12345ea8e38cb281c
# 2020/10/13 11:12:16 Public IP: 14.125.12.4
# 2020/10/13 11:12:17 Create or update record set: ptask.example.com

```



