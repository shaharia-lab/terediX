---
sidebar_position: 2
title: "AWS EC2"
---

# AWS EC2

| option	                                      | type     | description	                                                                                                                                                                                                                                                        |
|----------------------------------------------|----------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 	  [source_name].type                        | text     | `aws_ec2`	                                                                                                                                                                                                                                                          |
| 	  [source_name].configuration.access_key    | secret   | AWS Access Key	                                                                                                                                                                                                                                                     |
| 	  [source_name].configuration.secret_key    | secret   | AWS Secret Key	                                                                                                                                                                                                                                                     |
| 	  [source_name].configuration.region        | text     | AWS Region. e.g: us-west-1	                                                                                                                                                                                                                                         |
| 	  [source_name].configuration.session_token | secret   | AWS session token	                                                                                                                                                                                                                                                  |
| 	  [source_name].configuration.account_id    | text     | AWS account ID	                                                                                                                                                                                                                                                     |
| 	  [source_name].fields                      | list     | List of available fields to add with resource add meta data.<br/> <br/>**Available fields:**<br/> - instance_id<br/>- image_id<br/>- private_dns_name<br/>- instance_type<br/>- architecture<br/>- instance_lifecycle<br/>- instance_state<br/>- vpc_id<br/>- tags	 |
| 	  [source_name].schedule                    | interval | Set interval to schedule this source. e.g: `@every 10s`, `@every <br/>24h` or any valid cron expression `*/10 * * * * *`	                                                                                                                                           |

