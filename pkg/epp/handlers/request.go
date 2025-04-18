/*
Copyright 2025 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	extProcPb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"
	schedulingtypes "sigs.k8s.io/gateway-api-inference-extension/pkg/epp/scheduling/types"
	errutil "sigs.k8s.io/gateway-api-inference-extension/pkg/epp/util/error"
	logutil "sigs.k8s.io/gateway-api-inference-extension/pkg/epp/util/logging"
)

// HandleRequestBody always returns the requestContext even in the error case, as the request context is used in error handling.
func (s *StreamingServer) HandleRequestBody(
	ctx context.Context,
	reqCtx *RequestContext,
	req *extProcPb.ProcessingRequest,
	requestBodyMap map[string]interface{},
) (*RequestContext, error) {
	var requestBodyBytes []byte
	var err error
	logger := log.FromContext(ctx)

	// Resolve target models.
	reqCtx, err := s.translateModel(ctx, reqCtx, requestBodyMap)
	if err != nil {
		return reqCtx, err
	}

	target, err := s.scheduler.Schedule(ctx, &schedulingtypes.LLMRequest{
		InferenceModelName:  reqCtx.InferenceModelName,
		ResolvedTargetModel: reqCtx.ResolvedTargetModel,
		Critical:            modelObj.Spec.Criticality != nil && *modelObj.Spec.Criticality == v1alpha2.Critical,
	})
	if err != nil {
		return reqCtx, errutil.Error{Code: errutil.InferencePoolResourceExhausted, Msg: fmt.Errorf("failed to find target pod: %w", err).Error()}
	}
	targetPod := target.GetPod()

	// Insert target endpoint to instruct Envoy to route requests to the specified target pod.
	// Attach the port number
	pool, err := s.datastore.PoolGet()
	if err != nil {
		return reqCtx, err
	}
	endpoint := targetPod.Address + ":" + strconv.Itoa(int(pool.Spec.TargetPortNumber))

	logger.V(logutil.DEFAULT).Info("Request handled",
		"model", reqCtx.InferenceModelName, "targetModel", reqCtx.ResolvedTargetModel, "endpoint", targetPod, "endpoint metrics",
		fmt.Sprintf("%+v", target))

	requestBodyMap["model"] = reqCtx.ResolvedTargetModel
	requestBodyBytes, err = json.Marshal(requestBodyMap)
	if err != nil {
		logger.V(logutil.DEFAULT).Error(err, "Error marshaling request body")
		return reqCtx, errutil.Error{Code: errutil.Internal, Msg: fmt.Sprintf("error marshaling request body: %v", err)}
	}

	reqCtx.RequestSize = len(requestBodyBytes)
	reqCtx.InferenceModelName = reqCtx.InferenceModelName
	reqCtx.ResolvedTargetModel = reqCtx.ResolvedTargetModel
	reqCtx.TargetPod = targetPod.NamespacedName.String()
	reqCtx.TargetEndpoint = endpoint

	s.populateRequestHeaderResponse(reqCtx, endpoint, len(requestBodyBytes))

	reqCtx.reqBodyResp = &extProcPb.ProcessingResponse{
		// The Endpoint Picker supports two approaches to communicating the target endpoint, as a request header
		// and as an unstructure ext-proc response metadata key/value pair. This enables different integration
		// options for gateway providers.
		Response: &extProcPb.ProcessingResponse_RequestBody{
			RequestBody: &extProcPb.BodyResponse{
				Response: &extProcPb.CommonResponse{
					BodyMutation: &extProcPb.BodyMutation{
						Mutation: &extProcPb.BodyMutation_StreamedResponse{
							StreamedResponse: &extProcPb.StreamedBodyResponse{
								Body:        requestBodyBytes,
								EndOfStream: true,
							},
						},
					},
				},
			},
		},
	}
	return reqCtx, nil
}

// translateModel fetches the
func (s *StreamingServer) translateModel(ctx context.Context, reqCtx RequestContext, requestBody map[string]interface{}) (RequestContext, error) {
	var ok bool

	reqCtx.InferenceModelName, ok = requestBody["model"].(string)
	if !ok {
		return reqCtx, errutil.Error{Code: errutil.BadRequest, Msg: "model not found in request"}
	}

	// NOTE: The nil checking for the modelObject means that we DO allow passthrough currently.
	// This might be a security risk in the future where adapters not registered in the InferenceModel
	// are able to be requested by using their distinct name.
	modelObj := s.datastore.ModelGet(reqCtx.InferenceModelName)
	if modelObj == nil {
		return reqCtx, errutil.Error{Code: errutil.BadConfiguration, Msg: fmt.Sprintf("error finding a model object in InferenceModel for input %v", reqCtx.InferenceModelName)}
	} else if len(modelObj.Spec.TargetModels) > 0 {
		reqCtx.ResolvedTargetModel = RandomWeightedDraw(ctx, modelObj, 0)
		if reqCtx.ResolvedTargetModel == "" {
			return reqCtx, errutil.Error{Code: errutil.BadConfiguration, Msg: fmt.Sprintf("error getting target model name for model %v", modelObj.Name)}
		}
	}

	return reqCtx, nil
}

func (s *StreamingServer) HandleRequestHeaders(ctx context.Context, reqCtx *RequestContext, req *extProcPb.ProcessingRequest_RequestHeaders) error {
	reqCtx.RequestReceivedTimestamp = time.Now()

	// an EoS in the request headers means this request has no body or trailers.
	if req.RequestHeaders.EndOfStream {
		// We will route this request to a random pod as this is assumed to just be a GET
		// More context: https://github.com/kubernetes-sigs/gateway-api-inference-extension/pull/526
		// The above PR will address endpoint admission, but currently any request without a body will be
		// routed to a random upstream pod.
		pod := GetRandomPod(s.datastore)
		pool, err := s.datastore.PoolGet()
		if err != nil {
			return err
		}
		endpoint := pod.Address + ":" + strconv.Itoa(int(pool.Spec.TargetPortNumber))
		s.populateRequestHeaderResponse(reqCtx, endpoint, 0)
	}
	return nil
}
