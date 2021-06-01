#!/usr/bin/env bash

set -e

usage() {
  cat << EOF
Usage: REGION=<region> AWS_PROFILE=<profile> $(basename $0) <cluster_name>

Upgrade EKS cluster worker nodes by draining and terminating a single instance at a time.

This script requires ~10.5min to cycle EKS worker nodes.


ENVIRONMENT VARIABLES REQUIRED:

  REGION            the AWS region your cluster is deployed to
  AWS_PROFILE       an AWS profile to access your cluster


EXAMPLE:

  REGION=us-west-2 AWS_PROFILE=nwi-eks-admin ./upgrade-eks-workers.sh <cluster_name>

EOF
    if [[ ! -z $AWS_PROFILE ]]; then
        echo "Specify an EKS cluster you want to upgrade:"
        aws eks list-clusters --region $REGION --profile $AWS_PROFILE --no-cli-pager | jq '.clusters[]' | cut -d '"' -f 2 | xargs -I'{}' echo "  {}"
    fi
}

# check ENV Vars and cluster were provided
if [[ -z $REGION || -z $AWS_PROFILE || -z $1 ]]; then
    usage
    exit 1;
fi

# check that cluster is accessible via provided profile
CLUSTER_NAME=$1
if [[ $(aws eks list-clusters --region $REGION --profile $AWS_PROFILE --no-cli-pager | jq '.clusters[]' | cut -d '"' -f 2 | grep $CLUSTER_NAME | wc -l) -lt 1 ]]; then
    usage
    exit 1;
fi

# Configure EKS cluster kube config
aws eks update-kubeconfig --name "$CLUSTER_NAME" --region $REGION --profile $AWS_PROFILE

################################################################
# Pre-Check: ensure we have a healthy cluster to upgrade
################################################################

node_count=$(kubectl get node | wc -l)
if [[ $node_count -lt 4 ]];
then
    echo "UNABLE TO UPGRADE WORKERS: cluster is missing worker nodes."
    kubectl get node
    exit 1;
fi

node_not_ready_count=$(kubectl get node | egrep "NotReady|SchedulingDisabled" | wc -l)
if [[ $node_not_ready_count -gt 0 ]];
then
    echo "UNABLE TO UPGRADE WORKERS: eks worker nodes not ready"
    kubectl get node
    exit 1;
fi

echo ">>> Current pod count: $(kubectl get pod | wc -l )"
echo ">>> Current Node configuration:"
kubectl get nodes
echo ">>> Current Pod configuration:"
kubectl get pods -o wide



################################################################
# Upgrade EKS worker nodes
################################################################

WORKERS=($(kubectl describe node | grep 'Name:' | awk '{print $2}'))
WORKER_COUNT=${#WORKERS[@]}
looptime="5 minute"

for worker in "${WORKERS[@]}"
do
    POD_COUNT=$(kubectl get pods | wc -l)
    INSTANCE_ID=$(kubectl describe node $worker | grep 'ProviderID:' | awk '{print $2}' | cut -d'/' -f 5)

    kubectl drain $worker --ignore-daemonsets --delete-emptydir-data
    # poll/wait for pods to drain from worker node
    wait_for_pod_drain=true
    wait_for_loop_timeout=true
    endtime=$(date -ud "$looptime" +%s)
    while $wait_for_pod_drain && $wait_for_loop_timeout;
    do
        echo ">>> pods on this worker $worker"
        kubectl get pods -o wide | grep $worker
        if [[ $(kubectl drain $worker --ignore-daemonsets --delete-emptydir-data | grep "drained" | wc -l) -gt 0 ]]; then
            wait_for_pod_drain=false
        fi
        if [[ $(date -u +%s) -ge $endtime ]]; then
            wait_for_loop_timeout=false
            echo ">>> draining $worker timed out (will now terminate the ec2 instance)"
        fi
        sleep 2;
    done

    echo ">>> terminating EC2 instance $INSTANCE_ID"
    aws ec2 terminate-instances --instance-ids $INSTANCE_ID --region $REGION --profile $AWS_PROFILE --no-cli-pager || \
        (echo "The profile we are currently using doesnt have sufficient permission to terminate $INSTANCE_ID" && exit 1)

    # poll/wait for instance termination
    wait_for_instance_term=true
    endtime=$(date -ud "$looptime" +%s)
    while $wait_for_instance_term;
    do
        echo ">>> waiting for instance $INSTANCE_ID to terminate"
        sleep 2;
        if [[ $(aws ec2 terminate-instances --region $REGION --profile $AWS_PROFILE --instance-ids $INSTANCE_ID | jq '.TerminatingInstances[].CurrentState|.Name' | grep 'terminated' | wc -l) == 1 ]]; then
            wait_for_instance_term=false
        fi
        if [[ $(date -u +%s) -ge $endtime ]]; then
            echo ">>> TIMEOUT: ec2 instance termination taking longer than expected...exiting"
            exit 1
        fi
    done


    # poll/wait for replacement node to come online
    wait_for_replacement_node_prov=true
    endtime=$(date -ud "$looptime" +%s)
    while $wait_for_replacement_node_prov;
    do
        echo ">>> waiting for ASG to prov new k8s worker node"
        kubectl get node
        echo ">>> current pod count: $(kubectl get pod | wc -l ) (target pod count >= $POD_COUNT)"
        sleep 2;
        if [[ $(kubectl get node | grep -v 'NotReady' | grep 'Ready' | wc -l ) -ge $WORKER_COUNT && \
            $(kubectl get node | egrep 'NotReady|SchedulingDisabled' | wc -l ) -lt 1 ]]; then
            wait_for_replacement_node_prov=false
        fi
        if [[ $(date -u +%s) -ge $endtime ]]; then
            echo ">>> TIMEOUT: replacement worker node instance provisioning taking longer than expecting...exiting"
            exit 1
        fi
    done
done

echo ">>> EKS WORKERS UPGRADED!"
echo ">>> Current pod count: $(kubectl get pod | wc -l )"
echo ">>> Current Node configuration:"
kubectl get nodes
echo ">>> Current Pod configuration:"
kubectl get pods -o wide
