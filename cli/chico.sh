#!/bin/bash

DATE=$(date '+%Y-%m-%d-%H-%M-%S')
DOMAIN=$2
REGION=us-east-1
AVAILABILITY_ZONE=us-east-1a
AMI_ID=ami-09d56f8956ab235b3
INSTANCE_SIZE=t3a.small
KEY_NAME="$2-Key"
S3_BUCKET=""
SG_GROUP_NAME="$2-SG"

if [ $1 == "-d" ]
then
    echo $1 $2
fi

# Register Domain
text="$(jq --arg domain "$DOMAIN" '.DomainName = $domain' register-domain.json)" && \
echo -E "${text}" > register-domain.json

aws route53domains register-domain --region $REGION --cli-input-json file://register-domain.json

# Create Hosted Zone for Registered Domain
ZONE_ID=$(aws route53 create-hosted-zone --name $DOMAIN --caller-reference $DATE | jq '.HostedZone' | grep Id | grep -Eoh "[A-Z0-9]{2,}")

RAW_NS=$(aws route53 get-hosted-zone --id $ZONE_ID | jq '.DelegationSet' | grep -Eoh "ns-[0-9]+.awsdns-[0-9]+.[a-z]+" | cut -d " " -f 1)

nameservers=""

#for ns in $RAW_NS
do
    nameservers+="Name=$ns "
done

# Point Domain to Hosted Zone
aws route53domains update-domain-nameservers --region $REGION --domain-name $DOMAIN --nameservers $nameservers

# Create EC2 Key Pair
aws ec2 create-key-pair --key-name $KEY_NAME --query 'KeyMaterial' --output text > $KEY_NAME.pem

# Change Key Permissions
chmod 400 $KEY_NAME.pem

# Upload to S3
aws s3 cp $KEY_NAME.pem $S3_BUCKET

# Create VPC
VPC_ID=$(aws ec2 create-vpc --ipv6-cidr-block-network-border-group $REGION-lax-1 --cidr-block 10.0.0.0/16 | jq '.Vpc' | grep VpcId | grep -Eoh "vpc-[a-zA-Z0-9]+" )

# Create Security Group
SG_ID=$(aws ec2 create-security-group --group-name $SG_GROUP_NAME --description "$DOMAIN Security Group" --vpc-id $VPC_ID | jq '.GroupId' | grep -Eoh "sg-[a-zA-Z0-9]+" )

# Create EC2 Instance
aws ec2 run-instances --instance-names $DOMAIN --availability-zone $AVAILABILITY_ZONE \
    --image-id $AMI_ID --instance-type $INSTANCE_SIZE --count 1 \
    --key-name $KEY_NAME --security-group-ids $SG_ID --subnet-id subnet-6e7f829e

# Allocate Static IP
aws lightsail allocate-static-ip --static-ip-name $STATIC_IP_NAME

# Attach Static IP
aws lightsail attach-static-ip --static-ip-name $STATIC_IP_NAME --instance-name $DOMAIN