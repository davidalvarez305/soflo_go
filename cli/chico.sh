#!/bin/bash

DATE=$(date '+%Y-%m-%d-%H-%M-%S')
DOMAIN=$2

if [ $1 == "-d" ]
then
    echo $1 $2
fi

# Register Domain
text="$(jq --arg domain "$DOMAIN" '.DomainName = $domain' register-domain.json)" && \
echo -E "${text}" > register-domain.json

aws route53domains register-domain --region us-east-1 --cli-input-json file://register-domain.json

# Create Hosted Zone for Registered Domain
ZONE_ID=$(aws route53 create-hosted-zone --name $DOMAIN --caller-reference $DATE | jq '.HostedZone' | grep Id | grep -Eoh "[A-Z0-9]{2,}")

RAW_NS=$(aws route53 get-hosted-zone --id $ZONE_ID | jq '.DelegationSet' | grep -Eoh "ns-[0-9]+.awsdns-[0-9]+.[a-z]+" | cut -d " " -f 1)

nameservers=""

for ns in $RAW_NS
do
    nameservers+="Name=$ns "
done

echo $nameservers

# Point Domain to Hosted Zone
aws route53domains update-domain-nameservers --region us-east-1 --domain-name $DOMAIN --nameservers $nameservers

#echo $ZONE_ID