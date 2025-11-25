# Kafka Deployment Issue - Detailed Analysis

## Problem Summary
The distributed autoscaler started logging `failed to enqueue job` errors because the backing Kafka cluster was flapping and occasionally crashed outright. The instability came from conflicting Strimzi resources and an incomplete KRaft configuration, which in turn triggered stale client connections inside the API service.

## Root Causes

### 1. Duplicate `Kafka` Resources
- Two `Kafka` resources named `my-cluster` were applied in the same namespace.
- The duplicate CRDs caused competing KRaft voter lists, which produced `Received request or response with leader OptionalInt[0] and epoch 1 ...` and broker CrashLoopBackOff events.
- The cluster only worked temporarily until the KRaft quorum logic tried to reconcile the mismatched epochs.

### 2. Missing YAML Separator
- The original `kafka.yaml` omitted the `---` separator between the `Kafka` resource and the `KafkaNodePool`.
- Without the separator, `kubectl` could merge manifests unpredictably, further confusing the Strimzi operator.

### 3. Incorrect Node Pool Configuration
- Early patches defined separate broker/controller pools but forgot to set `kraftMetadata: shared`.
- In KRaft mode, brokers that are split from controllers must share metadata through that flag; without it, brokers never join the active quorum.

### 4. Stale Kafka Connections in the API
- After the cluster reconfiguration, the API server still held producers that pointed at the dead bootstrap service.
- Symptoms included `topic partition has no leader`, `Unknown Topic Or Partition`, and persistent `failed to enqueue job` errors.

## Technical Context
- Strimzi v0.49.0, Kafka 4.1.1 (KRaft).
- KRaft removes ZooKeeper, so controller/broker voter sets must be perfectly consistent.
- The Strimzi operator can only reconcile node pools when the parent `Kafka` resource is healthy and resource boundaries are clear.

## Final Working Configuration (`k8s/kafka/kafka.yaml`)
```yaml
apiVersion: kafka.strimzi.io/v1
kind: KafkaNodePool
metadata:
  name: dual-role
  namespace: kafka
  labels:
    strimzi.io/cluster: my-cluster
spec:
  replicas: 1
  roles:
    - controller
    - broker
  storage:
    type: ephemeral
    kraftMetadata: shared

---
apiVersion: kafka.strimzi.io/v1
kind: Kafka
metadata:
  name: my-cluster
  namespace: kafka
spec:
  kafka:
    version: 4.1.1
    metadataVersion: 4.1-IV1
    listeners:
      - name: plain
        port: 9092
        type: internal
        tls: false
    config:
      offsets.topic.replication.factor: 1
      transaction.state.log.replication.factor: 1
      transaction.state.log.min.isr: 1
      default.replication.factor: 1
      min.insync.replicas: 1
  entityOperator:
    topicOperator: {}
    userOperator: {}
```

### Key Fixes
1. Single dual-role node pool combines broker and controller responsibilities for local development.
2. `kraftMetadata: shared` allows the lone broker to read controller metadata.
3. All replication/min ISR settings reduced to `1` to match single-node requirements.
4. Removed the `strimzi.io/node-pools: enabled` annotation so Strimzi auto-discovers pools.

## Recovery Steps
1. Delete all existing `Kafka`/`KafkaNodePool` resources linked to `my-cluster`.
2. Apply the corrected manifest above.
3. Wait for the pods (`my-cluster-dual-role-0`, `my-cluster-entity-operator-0`) to reach `Running`.
4. Restart API pods to refresh Kafka producer connections.
5. Confirm topic leadership and connectivity with the verification commands below.

## Verification Commands
```bash
kubectl get kafka -n kafka
kubectl get kafkanodepool -n kafka
kubectl get pods -n kafka

kubectl exec my-cluster-dual-role-0 -n kafka -- \
  /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list

kubectl exec my-cluster-dual-role-0 -n kafka -- \
  /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 \
  --describe --topic jobs

kubectl logs <api-pod> -n default --tail=20
```

## Lessons Learned
1. Validate multi-resource YAML files with `yamllint` or `kubectl apply --dry-run=client`.
2. Ensure unique resource names to avoid Strimzi reconciliation loops.
3. Understand KRaft-specific requirements (voter configuration, `kraftMetadata` usage).
4. Restart or reconfigure Kafka clients when recreating clusters.
5. Make incremental changes and observe operator logs before moving on.

## Prevention Strategies
1. Add CI checks for YAML formatting and duplicate resource detection.
2. Watch Strimzi operator logs/events whenever applying Kafka manifests.
3. Improve application readiness checks around Kafka availability.
4. Implement exponential backoff and reconnection logic inside Kafka producers/consumers.
5. Keep working configuration examples in the repo to bootstrap future debugging.

## Prompt for Future Troubleshooting
Use this snippet when asking for help on similar incidents:

> "I'm experiencing Kafka deployment issues in Kubernetes using the Strimzi operator. Symptoms: [CrashLoopBackOff, connection errors, etc.] \
> Strimzi version: [version] \
> Kafka version: [version] \
> Deployment mode: [KRaft/ZooKeeper] \
> Node pool setup: [single/multiple, roles] \
> Please analyze YAML structure, KRaft quorum config, node pool metadata settings, Strimzi operator logs, and client connection handling."
