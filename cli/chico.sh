#!/bin/bash

DATE=$(date '+%Y-%m-%d-%H-%M-%S')
DOMAIN=$2.com
DATABASE=$4
USER=$6
REGION=us-east-1
AVAILABILITY_ZONE=us-east-1a
AMI_ID=ami-09d56f8956ab235b3
INSTANCE_SIZE=t3a.small
KEY_NAME=$2.pem
S3_BUCKET=$8

if [[ $1 == "-d" && $3 == "-db" && $5 == "-u" && $7 == "-b" ]];
then
    echo "Registering domain..."
    # Register Domain
    Content="$(jq --arg domain "$DOMAIN" '.DomainName = $domain' register-domain.json)" && echo -E "${Content}" > register-domain.json

    aws route53domains register-domain --region $REGION --cli-input-json file://register-domain.json

    echo "Creating hosted zone..."
    # Create Hosted Zone for Registered Domain
    ZONE_ID=$(aws route53 create-hosted-zone --name $DOMAIN --caller-reference $DATE | jq '.HostedZone' | grep Id | grep -Eoh "[A-Z0-9]{2,}")
    ZONE_ID=Z09517553MH38FF1I5MCZ

    echo "Changing name servers..."
    RAW_NS=$(aws route53 get-hosted-zone --id $ZONE_ID | jq '.DelegationSet' | grep -Eoh "ns-[0-9]+.awsdns-[0-9]+.[a-z]+" | cut -d " " -f 1)

    nameservers=""

    for ns in $RAW_NS
    do
        nameservers+="Name=$ns "
    done

    # Point Domain to Hosted Zone
    aws route53domains update-domain-nameservers --region $REGION --domain-name $DOMAIN --nameservers $nameservers

    echo "Creating key pairs for EC2..."
    # Create EC2 Key Pair
    aws ec2 create-key-pair --key-name $KEY_NAME --query 'KeyMaterial' --output text > $KEY_NAME

    echo "Changing key permissions..."
    # Change Key Permissions
    sudo chmod 400 $KEY_NAME

    # Upload to S3
    aws s3 cp $KEY_NAME $S3_BUCKET/keys/

    echo "Creating EC2 Instance..."
    # Create EC2 Instance
    INSTANCE_ID=$(aws ec2 run-instances --image-id $AMI_ID --instance-type $INSTANCE_SIZE \
        --count 1 --associate-public-ip-address \
        --key-name $KEY_NAME | grep InstanceId | grep -Eoh "i-[a-z0-9]+")
    
    INSTANCE_ID=i-04e62e05a49443c9f

    # Get Instance Public Id
    EC2_PUBLIC_ID=$(aws ec2 describe-instances --instance-ids $INSTANCE_ID | grep PublicIpAddress | grep -Eoh "[0-9.]+")

    # Update Hosted Zone A Record to EC2 Public Id
    echo "Updating A Record to Point to EC2 Instance..."
    Text="$(jq \
        --arg ip "$EC2_PUBLIC_ID" \
        --arg dns "$DOMAIN" \
        '.Changes[].ResourceRecordSet.ResourceRecords = [{ Value: $ip }] | .Changes[].ResourceRecordSet.Name = $dns' \
        change-hosted-zone.json)" && echo -E "${Text}" > change-hosted-zone.json

    aws route53 change-resource-record-sets --hosted-zone-id $ZONE_ID --change-batch file://change-hosted-zone.json

    scp -r -i $KEY_NAME ./prep ubuntu@$EC2_PUBLIC_ID:/home/ubuntu/

    echo "Initiating server setup..."
    ssh -i $KEY_NAME ubuntu@$EC2_PUBLIC_ID "yes y | chmod +x ./prep/server.sh && sudo ./prep/server.sh $DATABASE $USER $S3_BUCKET"

else
    echo "Missing either -d or -db or -u flag"
    exit 1
fi