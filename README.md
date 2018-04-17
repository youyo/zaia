# zabbix-aws-integration-agent

Used by Zabbix LLD.  
It performs authentication by Assume role.  
  
The role of zabbix-aws-integration-agent.

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "sts:AssumeRole"
            ],
            "Resource": "*",
            "Effect": "Allow"
        }
    ]
}
```

The side to be integrated has a ReadOnlyAccess by aws maneged policy.  
Trust Relationship policy.

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::00000000:root"
      },
      "Action": "sts:AssumeRole",
      "Condition": {}
    }
  ]
}
```

## Usage

### Start process

zabbix-aws-integration-agent is running in docker container.  
https://hub.docker.com/r/youyo/zabbix-aws-integration-agent/

```
$ docker container run -d -p 10050:10050 -e AWS_ACCESS_KEY_ID=your_aws_access_key -e AWS_SECRET_ACCESS_KEY=your_aws_secret_access_key --restart always --name zabbix-aws-integration-agent youyo/zabbix-aws-integration-agent:latest
```

### Request method

Request method is zabbix protocol.

- Discovery of running EC2 instances.

Syntax

```
$ zabbix_get -s zabbix-aws-integration-agent-host -k aws-integration.ec2.discovery[Zabbix-host-group,ARN,Region]
```

Example

```
$ zabbix_get -s 127.0.0.1 -k aws-integration.ec2.discovery[zabbix-host-group,arn:aws:iam::00000000:role/rolename,ap-northeast-1]
{
  "data": [
    {
      "{#INSTANCE_ID}": "i-0000xxxx",
      "{#INSTANCE_NAME}": "Name tag of the instance",
      "{#INSTANCE_ROLE}": "Role tag of the instance",
      "{#INSTANCE_PUBLIC_IP}": "public ip",
      "{#INSTANCE_PRIVATE_IP}": "private ip",
      "{#ZABBIX_HOST_GROUP}": "zabbix-host-group"
    }
  ]
}
```

- Whether the ec2 instance has maintenance

Syntax

```
$ zabbix_get -s zabbix-aws-integration-agent-host -k aws-integration.ec2.maintenance[Instance-id,No maintenance message,ARN,Region]
```

Example

```
$ zabbix_get -s 127.0.0.1 -k aws-integration.ec2.maintenance[i-0000xxxx,No maintenance,arn:aws:iam::00000000:role/rolename,ap-northeast-1]
No maintenance
```

- Fetch CloudWatch data

Syntax

```
$ zabbix_get -s zabbix-aws-integration-agent-host -k aws-integration.cloudwatch.get-metrics[Namespace,Dimension-name,Dimension-value,Metrics,Statistic,ARN,Region]
```

Example

```
$ zabbix_get -s 127.0.0.1 -k aws-integration.cloudwatch.get-metrics[AWS/EC2,InstanceId,i-0000xxxx,CPUUtilization,Average,arn:aws:iam::00000000:role/rolename,ap-northeast-1]
7.696
```

- Discovery of RDS instances.

Syntax

```
$ zabbix_get -s zabbix-aws-integration-agent-host -k aws-integration.rds.discovery[Zabbix-host-group,ARN,Region]
```

Example

```
$ zabbix_get -s 127.0.0.1 -k aws-integration.rds.discovery[zabbix-host-group,arn:aws:iam::00000000:role/rolename,ap-northeast-1]
{
  "data": [
    {
      "{#DB_INSTANCE_ARN}": "arn:aws:rds:ap-northeast-1:00000000:db:identifier",
      "{#DB_INSTANCE_IDENTIFIER}": "identifier",
      "{#ENGINE}": "aurora-mysql",
      "{#ZABBIX_HOST_GROUP}": "zabbix-host-group"
    }
  ]
}
```
