Zabbix-SQS alertscript
======================

Forward Zabbix alerts to SQS queue by using a custom Zabbix [alertscript](https://www.zabbix.com/documentation/3.4/manual/config/notifications/media/script).

Script installation
-------------------

Simply decompress archive found in [releases section](https://github.com/claranet/zabbix-sqs/releases):

    $ bunzip2 `zabbix-sqs.bz2`

Then mv `zabbix-sqs` bonary and its json config file on zabbix server to the `AlertScriptsPath`
directory which can be found in the `/etc/zabbix/zabbix_server.conf` configuration file:

    $ grep -e '^AlertScriptsPath' /etc/zabbix/zabbix_server.conf
    AlertScriptsPath=/path/to/zabbix/alertscripts
    $ mv `zabbix-sqs` /path/to/zabbix/alertscripts
    $ mv `zabbix-sqs.json` /path/to/zabbix/alertscripts

Be sure zabbix user is able to execute `zabbix-sqs` file

Zabbix configuration
--------------------

To forward Zabbix events to SQS a new media script needs to be created
and associated with a user. Follow the steps below as a Zabbix Admin user...

1/ Create a new media type [Admininstration > Media Types > Create Media Type]

```
Name: SQS
Type: Script
Script name: zabbix-sqs
Script parameters:
    1st: {ALERT.MESSAGE}
    2nd: Nothing (remove field)
    3nd: Nothing (remove field)
Enabled: [x]
```

2/ Modify the Media for the Admin user or create a new one if you prefer filter permissions [Administration > Users]

```
Type: SQS
Send to: sqs-queue
When active: 1-7,00:00-24:00
Use if severity: (all)
Status: Enabled
```

3/ Configure Action [Configuration > Actions > Create Action > Action]

```
Name: zabbix-sqs
```
```
Default subject: Nothing (empty field)
```
```
{"ID":{TRIGGER.ID},"Name":"{TRIGGER.NAME}","Status":"{TRIGGER.STATUS}","Group":"{TRIGGER.HOSTGROUP.NAME}","Severity":"{TRIGGER.SEVERITY}","Hostname":"{HOSTNAME}","IP":"{IPADDRESS}","Items":[{"Name":"{ITEM.NAME1}","Value":"{ITEM.VALUE1}"},{"Name":"{ITEM.NAME2}","Value":"{ITEM.VALUE2}"},{"Name":"{ITEM.NAME3}","Value":"{ITEM.VALUE3}"},{"Name":"{ITEM.NAME4}","Value":"{ITEM.VALUE4}"},{"Name":"{ITEM.NAME5}","Value":"{ITEM.VALUE5}"},{"Name":"{ITEM.NAME6}","Value":"{ITEM.VALUE6}"},{"Name":"{ITEM.NAME7}","Value":"{ITEM.VALUE7}"},{"Name":"{ITEM.NAME8}","Value":"{ITEM.VALUE8}"},{"Name":"{ITEM.NAME9}","Value":"{ITEM.VALUE9}"}]}
```

RECOVERY
```
Default subject: Nothing (empty field)
```
```
{"ID":{TRIGGER.ID},"Name":"{TRIGGER.NAME}","Status":"{TRIGGER.STATUS}","Group":"{TRIGGER.HOSTGROUP.NAME}","Severity":"{TRIGGER.SEVERITY}","Hostname":"{HOSTNAME}","IP":"{IPADDRESS}","Items":[{"Name":"{ITEM.NAME1}","Value":"{ITEM.VALUE1}"},{"Name":"{ITEM.NAME2}","Value":"{ITEM.VALUE2}"},{"Name":"{ITEM.NAME3}","Value":"{ITEM.VALUE3}"},{"Name":"{ITEM.NAME4}","Value":"{ITEM.VALUE4}"},{"Name":"{ITEM.NAME5}","Value":"{ITEM.VALUE5}"},{"Name":"{ITEM.NAME6}","Value":"{ITEM.VALUE6}"},{"Name":"{ITEM.NAME7}","Value":"{ITEM.VALUE7}"},{"Name":"{ITEM.NAME8}","Value":"{ITEM.VALUE8}"},{"Name":"{ITEM.NAME9}","Value":"{ITEM.VALUE9}"}]}
```

https://www.zabbix.com/documentation/3.2/manual/appendix/macros/supported_by_location

To send OK events ...

````
Recovery message: [check]
Enabled [check]
````

At the Conditions tab, to only forward PROBLEM and OK events ...

```
(A)	Maintenance status not in "maintenance"
(B)	Trigger value = "PROBLEM"
```

To forward PROBLEM, ACKNOWLEDGED, OK events ...

```
(A)	Maintenance status not in "maintenance"
```

Finally, add an operation:

```
Send to Users: Admin (or previously created user)
Send only to: zabbix-sqs
```

Script configuration
--------------------

The configuration file `zabbix-sqs.json` should be located next to the script in alertscripts directory

**Configuration parameters**

The following parameters should be configured :

  * QueueURL
  * Region
  * AccessKeyID
  * SecretAccessKey

in `zabbix-sqs.json` file next to the binaray `zabbix-sqs`

Here is a [sample configuration file](zabbix-sqs.json.sample)

AWS configuration
-----------------

The following resources should be created :

  * IAM user
  * IAM policy
  * SQS Queue

First, create the queue and configure it according your needs. Get the SQS `QueueURL` and its `Region` to set it in the configuration file.

Then, create an IAM user for zabbix and generate credentials. Get `AccessKeyID` and `SecretAccessKey` to set it in the configuration file.

Finally, add an inline IAM policy to the previously created user to grant it SQS privileges. Taking example from [the sample configuration file](zabbix-sqs.json.sample) :

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "sqs:ListQueues",
            "Resource": "*",
            "Condition": {
                "IpAddress": {
                    "aws:SourceIp": "#ZABBIX_SERVER_IP#"
                }
            }
        },
        {
            "Effect": "Allow",
            "Action": "sqs:*",
            "Resource": "https://sqs.ap-northeast-1.amazonaws.com/42133769/my-queue",
            "Condition": {
                "IpAddress": {
                    "aws:SourceIp": "#ZABBIX_SERVER_IP#"
                }
            }
        }
    ]
}
```

Troubleshooting
---------------

See the [PagerDuty guide](http://www.pagerduty.com/docs/guides/zabbix-integration-guide/)
to configuring Zabbix integrations for an example installation with
screenshots.

References
----------

  * [Zabbix Custom Alert Scripts](https://www.zabbix.com/documentation/3.4/manual/config/notifications/media/script)
  * [Zabbix Custom User Macros](https://www.zabbix.com/documentation/3.4/manual/config/macros/usermacros)

License
-------

Copyright (c) 2018 Claranet. Available under the MIT License.

