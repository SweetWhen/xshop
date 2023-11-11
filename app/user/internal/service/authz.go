package service

import (
	"context"
	"fmt"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (us *UserService) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	fmt.Printf("user check get req:%s\n", req.String())
	return &authv3.CheckResponse{
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{
				Headers: []*corev3.HeaderValueOption{
					{
						Header: &corev3.HeaderValue{
							Key:   "X-Operator-Key",
							Value: fmt.Sprintf("%d", 123),
						},
						Append: &wrapperspb.BoolValue{},
					},
					{
						Header: &corev3.HeaderValue{
							Key:   "Name",
							Value: "123name",
						},
						Append: &wrapperspb.BoolValue{},
					},
					{
						Header: &corev3.HeaderValue{
							Key:   "CONTENT-TYPE",
							Value: "application/json",
						},
						Append: &wrapperspb.BoolValue{},
					}, {
						Header: &corev3.HeaderValue{
							Key:   "Connection",
							Value: "keep-alive",
						},
						Append: &wrapperspb.BoolValue{},
					},
				},
			},
		},
	}, nil
}
