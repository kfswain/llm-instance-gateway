# API Reference

## Packages
- [inference.networking.x-k8s.io/v1alpha2](#inferencenetworkingx-k8siov1alpha2)


## inference.networking.x-k8s.io/v1alpha2

Package v1alpha2 contains API Schema definitions for the
inference.networking.x-k8s.io API group.


### Resource Types
- [InferencePool](#inferencepool)


#### EndpointPickerConfig



EndpointPickerConfig specifies the configuration needed by the proxy to discover and connect to the endpoint picker extension.
This type is intended to be a union of mutually exclusive configuration options that we may add in the future.



_Appears in:_
- [InferencePoolSpec](#inferencepoolspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `extensionRef` _[Extension](#extension)_ | Extension configures an endpoint picker as an extension service. |  | Required: \{\} <br /> |


#### Extension



Extension specifies how to configure an extension that runs the endpoint picker.



_Appears in:_
- [EndpointPickerConfig](#endpointpickerconfig)
- [InferencePoolSpec](#inferencepoolspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `group` _[Group](#group)_ | Group is the group of the referent.<br />The default value is "", representing the Core API group. |  | MaxLength: 253 <br />Pattern: `^$\|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$` <br /> |
| `kind` _[Kind](#kind)_ | Kind is the Kubernetes resource kind of the referent.<br />Defaults to "Service" when not specified.<br />ExternalName services can refer to CNAME DNS records that may live<br />outside of the cluster and as such are difficult to reason about in<br />terms of conformance. They also may not be safe to forward to (see<br />CVE-2021-25740 for more information). Implementations MUST NOT<br />support ExternalName Services. | Service | MaxLength: 63 <br />MinLength: 1 <br />Pattern: `^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$` <br /> |
| `name` _[ObjectName](#objectname)_ | Name is the name of the referent. |  | MaxLength: 253 <br />MinLength: 1 <br />Required: \{\} <br /> |
| `portNumber` _[PortNumber](#portnumber)_ | The port number on the service running the extension. When unspecified,<br />implementations SHOULD infer a default value of 9002 when the Kind is<br />Service. |  | Maximum: 65535 <br />Minimum: 1 <br /> |
| `failureMode` _[ExtensionFailureMode](#extensionfailuremode)_ | Configures how the gateway handles the case when the extension is not responsive.<br />Defaults to failClose. | FailClose | Enum: [FailOpen FailClose] <br /> |


#### ExtensionConnection



ExtensionConnection encapsulates options that configures the connection to the extension.



_Appears in:_
- [Extension](#extension)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `failureMode` _[ExtensionFailureMode](#extensionfailuremode)_ | Configures how the gateway handles the case when the extension is not responsive.<br />Defaults to failClose. | FailClose | Enum: [FailOpen FailClose] <br /> |


#### ExtensionFailureMode

_Underlying type:_ _string_

ExtensionFailureMode defines the options for how the gateway handles the case when the extension is not
responsive.

_Validation:_
- Enum: [FailOpen FailClose]

_Appears in:_
- [Extension](#extension)
- [ExtensionConnection](#extensionconnection)

| Field | Description |
| --- | --- |
| `FailOpen` | FailOpen specifies that the proxy should forward the request to an endpoint of its picking when the Endpoint Picker fails.<br /> |
| `FailClose` | FailClose specifies that the proxy should drop the request when the Endpoint Picker fails.<br /> |


#### ExtensionReference



ExtensionReference is a reference to the extension.

Connections to this extension MUST use TLS by default. Implementations MAY
provide a way to customize this connection to use cleartext, a different
protocol, or custom TLS configuration.

If a reference is invalid, the implementation MUST update the `ResolvedRefs`
Condition on the InferencePool's status to `status: False`. A 5XX status code
MUST be returned for the request that would have otherwise been routed to the
invalid backend.



_Appears in:_
- [Extension](#extension)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `group` _[Group](#group)_ | Group is the group of the referent.<br />The default value is "", representing the Core API group. |  | MaxLength: 253 <br />Pattern: `^$\|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$` <br /> |
| `kind` _[Kind](#kind)_ | Kind is the Kubernetes resource kind of the referent.<br />Defaults to "Service" when not specified.<br />ExternalName services can refer to CNAME DNS records that may live<br />outside of the cluster and as such are difficult to reason about in<br />terms of conformance. They also may not be safe to forward to (see<br />CVE-2021-25740 for more information). Implementations MUST NOT<br />support ExternalName Services. | Service | MaxLength: 63 <br />MinLength: 1 <br />Pattern: `^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$` <br /> |
| `name` _[ObjectName](#objectname)_ | Name is the name of the referent. |  | MaxLength: 253 <br />MinLength: 1 <br />Required: \{\} <br /> |
| `portNumber` _[PortNumber](#portnumber)_ | The port number on the service running the extension. When unspecified,<br />implementations SHOULD infer a default value of 9002 when the Kind is<br />Service. |  | Maximum: 65535 <br />Minimum: 1 <br /> |


#### Group

_Underlying type:_ _string_

Group refers to a Kubernetes Group. It must either be an empty string or a
RFC 1123 subdomain.

This validation is based off of the corresponding Kubernetes validation:
https://github.com/kubernetes/apimachinery/blob/02cfb53916346d085a6c6c7c66f882e3c6b0eca6/pkg/util/validation/validation.go#L208

Valid values include:

* "" - empty string implies core Kubernetes API group
* "gateway.networking.k8s.io"
* "foo.example.com"

Invalid values include:

* "example.com/bar" - "/" is an invalid character

_Validation:_
- MaxLength: 253
- Pattern: `^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`

_Appears in:_
- [Extension](#extension)
- [ExtensionReference](#extensionreference)
- [ParentGatewayReference](#parentgatewayreference)
- [PoolObjectReference](#poolobjectreference)



#### InferencePool



InferencePool is the Schema for the InferencePools API.






| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `inference.networking.x-k8s.io/v1alpha2` | | |
| `kind` _string_ | `InferencePool` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[InferencePoolSpec](#inferencepoolspec)_ |  |  |  |
| `status` _[InferencePoolStatus](#inferencepoolstatus)_ | Status defines the observed state of InferencePool. | \{ parent:[map[conditions:[map[lastTransitionTime:1970-01-01T00:00:00Z message:Waiting for controller reason:Pending status:Unknown type:Accepted]] parentRef:map[kind:Status name:default]]] \} |  |






#### InferencePoolSpec



InferencePoolSpec defines the desired state of InferencePool



_Appears in:_
- [InferencePool](#inferencepool)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `selector` _object (keys:[LabelKey](#labelkey), values:[LabelValue](#labelvalue))_ | Selector defines a map of labels to watch model server Pods<br />that should be included in the InferencePool.<br />In some cases, implementations may translate this field to a Service selector, so this matches the simple<br />map used for Service selectors instead of the full Kubernetes LabelSelector type.<br />If specified, it will be applied to match the model server pods in the same namespace as the InferencePool.<br />Cross namesoace selector is not supported. |  | Required: \{\} <br /> |
| `targetPortNumber` _integer_ | TargetPortNumber defines the port number to access the selected model server Pods.<br />The number must be in the range 1 to 65535. |  | Maximum: 65535 <br />Minimum: 1 <br />Required: \{\} <br /> |
| `extensionRef` _[Extension](#extension)_ | Extension configures an endpoint picker as an extension service. |  | Required: \{\} <br /> |


#### InferencePoolStatus



InferencePoolStatus defines the observed state of InferencePool.



_Appears in:_
- [InferencePool](#inferencepool)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `parent` _[PoolStatus](#poolstatus) array_ | Parents is a list of parent resources (usually Gateways) that are<br />associated with the InferencePool, and the status of the InferencePool with respect to<br />each parent.<br />A maximum of 32 Gateways will be represented in this list. When the list contains<br />`kind: Status, name: default`, it indicates that the InferencePool is not<br />associated with any Gateway and a controller must perform the following:<br /> - Remove the parent when setting the "Accepted" condition.<br /> - Add the parent when the controller will no longer manage the InferencePool<br />   and no other parents exist. |  | MaxItems: 32 <br /> |


#### Kind

_Underlying type:_ _string_

Kind refers to a Kubernetes Kind.

Valid values include:

* "Service"
* "HTTPRoute"

Invalid values include:

* "invalid/kind" - "/" is an invalid character

_Validation:_
- MaxLength: 63
- MinLength: 1
- Pattern: `^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$`

_Appears in:_
- [Extension](#extension)
- [ExtensionReference](#extensionreference)
- [ParentGatewayReference](#parentgatewayreference)
- [PoolObjectReference](#poolobjectreference)



#### LabelKey

_Underlying type:_ _string_

LabelKey was originally copied from: https://github.com/kubernetes-sigs/gateway-api/blob/99a3934c6bc1ce0874f3a4c5f20cafd8977ffcb4/apis/v1/shared_types.go#L694-L731
Duplicated as to not take an unexpected dependency on gw's API.

LabelKey is the key of a label. This is used for validation
of maps. This matches the Kubernetes "qualified name" validation that is used for labels.
Labels are case sensitive, so: my-label and My-Label are considered distinct.

Valid values include:

* example
* example.com
* example.com/path
* example.com/path.html

Invalid values include:

* example~ - "~" is an invalid character
* example.com. - can not start or end with "."

_Validation:_
- MaxLength: 253
- MinLength: 1
- Pattern: `^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]{0,61})?[A-Za-z0-9]$`

_Appears in:_
- [InferencePoolSpec](#inferencepoolspec)



#### LabelValue

_Underlying type:_ _string_

LabelValue is the value of a label. This is used for validation
of maps. This matches the Kubernetes label validation rules:
* must be 63 characters or less (can be empty),
* unless empty, must begin and end with an alphanumeric character ([a-z0-9A-Z]),
* could contain dashes (-), underscores (_), dots (.), and alphanumerics between.

Valid values include:

* MyValue
* my.name
* 123-my-value

_Validation:_
- MaxLength: 63
- MinLength: 0
- Pattern: `^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$`

_Appears in:_
- [InferencePoolSpec](#inferencepoolspec)



#### Namespace

_Underlying type:_ _string_

Namespace refers to a Kubernetes namespace. It must be a RFC 1123 label.

This validation is based off of the corresponding Kubernetes validation:
https://github.com/kubernetes/apimachinery/blob/02cfb53916346d085a6c6c7c66f882e3c6b0eca6/pkg/util/validation/validation.go#L187

This is used for Namespace name validation here:
https://github.com/kubernetes/apimachinery/blob/02cfb53916346d085a6c6c7c66f882e3c6b0eca6/pkg/api/validation/generic.go#L63

Valid values include:

* "example"

Invalid values include:

* "example.com" - "." is an invalid character

_Validation:_
- MaxLength: 63
- MinLength: 1
- Pattern: `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`

_Appears in:_
- [ParentGatewayReference](#parentgatewayreference)



#### ObjectName

_Underlying type:_ _string_

ObjectName refers to the name of a Kubernetes object.
Object names can have a variety of forms, including RFC 1123 subdomains,
RFC 1123 labels, or RFC 1035 labels.

_Validation:_
- MaxLength: 253
- MinLength: 1

_Appears in:_
- [Extension](#extension)
- [ExtensionReference](#extensionreference)
- [ParentGatewayReference](#parentgatewayreference)
- [PoolObjectReference](#poolobjectreference)



#### ParentGatewayReference



ParentGatewayReference identifies an API object including its namespace,
defaulting to Gateway.



_Appears in:_
- [PoolStatus](#poolstatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `group` _[Group](#group)_ | Group is the group of the referent. | gateway.networking.k8s.io | MaxLength: 253 <br />Pattern: `^$\|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$` <br /> |
| `kind` _[Kind](#kind)_ | Kind is kind of the referent. For example "Gateway". | Gateway | MaxLength: 63 <br />MinLength: 1 <br />Pattern: `^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$` <br /> |
| `name` _[ObjectName](#objectname)_ | Name is the name of the referent. |  | MaxLength: 253 <br />MinLength: 1 <br /> |
| `namespace` _[Namespace](#namespace)_ | Namespace is the namespace of the referent.  If not present,<br />the namespace of the referent is assumed to be the same as<br />the namespace of the referring object. |  | MaxLength: 63 <br />MinLength: 1 <br />Pattern: `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$` <br /> |


#### PoolStatus



PoolStatus defines the observed state of InferencePool from a Gateway.



_Appears in:_
- [InferencePoolStatus](#inferencepoolstatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `parentRef` _[ParentGatewayReference](#parentgatewayreference)_ | GatewayRef indicates the gateway that observed state of InferencePool. |  |  |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#condition-v1-meta) array_ | Conditions track the state of the InferencePool.<br />Known condition types are:<br />* "Accepted"<br />* "ResolvedRefs" | [map[lastTransitionTime:1970-01-01T00:00:00Z message:Waiting for controller reason:Pending status:Unknown type:Accepted]] | MaxItems: 8 <br /> |


#### PortNumber

_Underlying type:_ _integer_

PortNumber defines a network port.

_Validation:_
- Maximum: 65535
- Minimum: 1

_Appears in:_
- [Extension](#extension)
- [ExtensionReference](#extensionreference)
