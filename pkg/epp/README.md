# The EndPoint Picker (EPP)
This package provides the reference implementation for the Endpoint Picker (EPP). As demonstrated in the diagram below, it implements the [extension protocol](../../docs/proposals/004-endpoint-picker-protocol), enabling a proxy or gateway to request endpoint hints from an extension, and interacts with the model servers through the defined [model server protocol](../..//docs/proposals/003-model-server-protocol).

![Architecture Diagram](../../docs/endpoint-picker.svg)


## Core Functions

An EPP instance handles a single `InferencePool` (and so for each `InferencePool`, one must create a dedicated EPP deployment), it performs the following core functions:

- Endpoint Selection
  - The EPP determines the appropriate Pod endpoint for the load balancer (LB) to route requests.
  - It selects from the pool of ready Pods designated by the assigned InferencePool's [Selector](https://github.com/kubernetes-sigs/gateway-api-inference-extension/blob/7e3cd457cdcd01339b65861c8e472cf27e6b6e80/api/v1alpha1/inferencepool_types.go#L53) field.
  - EPP optionally rewrites the model name in the request to the model name specified in the [model name rewrite header key](https://github.com/kubernetes-sigs/gateway-api-inference-extension/blob/cc9d7711f2661a09beb1b4cc26fa0bae48c47d4d/pkg/epp/metadata/consts.go#L34)
- Observability
  - The EPP generates metrics to enhance observability.
  - It reports InferenceObjective-level metrics, further broken down by model name.
  - Detailed information regarding metrics can be found on the [website](https://gateway-api-inference-extension.sigs.k8s.io/guides/metrics/).


## Scheduling Algorithm 
The scheduling package implements request scheduling algorithms for load balancing requests across backend pods in an inference gateway. The scheduler ensures efficient resource utilization while maintaining low latency and prioritizing critical requests. It applies a series of filters based on metrics and heuristics to select the best pod for a given request. The following flow chart summarizes the current scheduling algorithm

<img src="../../docs/scheduler-flowchart.png" alt="Scheduling Algorithm" width="400" />
