#! /bin/bash
#

phone=$1
message=$2
server="http://xxx.com/api/v1/notify/sms"

curl -H "Content-Type:application/json" -X POST "$server" --data "{\"app\": \"std\", \"tos\": [\"$phone\"], \"content\": {\"msg\":\"$message\"}}"