#!/bin/bash

DATE=$(date '+%Y-%m-%d-%H-%M-%S')
DOMAIN=$2
DATABASE=$4
USER=$6
REGION=us-east-1
AVAILABILITY_ZONE=us-east-1a
AMI_ID=ami-09d56f8956ab235b3
INSTANCE_SIZE=t3a.small
KEY_NAME="$2-Key"
S3_BUCKET=""
SG_GROUP_NAME="$2-SG"

if [[ $1 == "-d" && $3 == "-db" && $5 == "-u" ]];
then
    # Register Domain
    Content="$(jq --arg domain "$DOMAIN" '.DomainName = $domain' register-domain.json)" && echo -E "${Content}" > register-domain.json

    aws route53domains register-domain --region $REGION --cli-input-json file://register-domain.json

    # Create Hosted Zone for Registered Domain
    ZONE_ID=$(aws route53 create-hosted-zone --name $DOMAIN --caller-reference $DATE | jq '.HostedZone' | grep Id | grep -Eoh "[A-Z0-9]{2,}")

    RAW_NS=$(aws route53 get-hosted-zone --id $ZONE_ID | jq '.DelegationSet' | grep -Eoh "ns-[0-9]+.awsdns-[0-9]+.[a-z]+" | cut -d " " -f 1)

    nameservers=""

    for ns in $RAW_NS
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

    # Create EC2 Instance
    INSTANCE_ID=$(aws ec2 run-instances --image-id $AMI_ID --instance-type $INSTANCE_SIZE \
        --count 1 --associate-public-ip-address \
        --key-name $KEY_NAME | grep InstanceId | grep -Eoh "i-[a-z0-9]+")

    # Get Instance Public Id
    EC2_PUBLIC_ID=$(aws ec2 describe-instances --instance-ids $INSTANCE_ID | grep PublicIpAddress | greg -Eoh "[0-9.]+")

    # Update Hosted Zone A Record to EC2 Public Id
    Text="$(jq \
        --arg ip "$EC2_PUBLIC_ID" \
        --arg dns "$DOMAIN" \
        '.Changes[].ResourceRecordSet.ResourceRecords = [{ Value: $ip }] | \
        .Changes[].ResourceRecordSet.Name = $dns' \
        change-zone.json)" && echo -E "${Text}" > change-zone.json

    aws route53 change-resource-record-sets --hosted-zone-id $ZONE_ID --change-batch file://change-zone.json

    # Update Repository
    apt-get update

    apt-get install \
        ca-certificates \
        curl \
        gnupg \
        lsb-release \
        apt-transport-https \
        software-properties-common \
        unzip

    # Docker Installation

    # Add Docker Official GPG Key
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

    # Add Docker Repo
    add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable"

    # Install Docker Engine
    apt install docker-ce

    # Get Docker Compose
    curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

    # Make File Executable
    chmod +x /usr/local/bin/docker-compose

    # Installation PSQL
    apt-get install postgresql-client

    # Install AWS CLI
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
    unzip awscliv2.zip && ./aws/install

    # Configure AWS CLI
    scp -i $DOMAIN-Key.pem ~/.aws/credentials ubuntu@$EC2_PUBLIC_ID:~/.aws/credentials

    # Copy Env File
    scp -i $DOMAIN-Key.pem .env ubuntu@$EC2_PUBLIC_ID:~/soflo_go/server

    # Append DB Name to Env File
    echo "POSTGRES_DB=$DATABASE" >> .env

    # Download Server Files
    cd && git clone https://github.com/davidalvarez305/soflo_go.git

    # Download SQL File
    aws s3 cp ~/db/$DATABASE.sql ~/soflo_go/postgres/

    # Copy SQL to Database
    psql -h localhost -U $USER $DATABASE -f ~/soflo_go/postgres/$DATABASE.sql

    # Start Server
    cd soflo_go && sudo docker-compose -f docker-compose.yml up --build

else
    echo "Missing either -d or -db or -u flag"
    exit 1
fi