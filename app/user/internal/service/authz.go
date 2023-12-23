package service

import (
	"context"
	"fmt"
	"net/http"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/go-kratos/kratos/v2/log"
	status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

func (us *UserService) IsWithoutAuthPath(p, method string) bool {
	if method == http.MethodPost && p == "/user/v1/users/login" {
		return true
	}

	return false
}

func (us *UserService) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	p := req.Attributes.Request.Http.Path
	method := req.Attributes.Request.Http.Method
	log.Infof("user check get p:%s, method:%s, req:%s\n", p, method, req.String())
	if us.IsWithoutAuthPath(p, method) {
		return &authv3.CheckResponse{HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{},
		}}, nil
	}

	// need check cookie and set jwt header
	cookeVaule := req.Attributes.Request.Http.Headers["cookie"]
	request, err := http.NewRequest(method, p, nil)
	if err != nil {
		log.Errorf("NewRequest err:%v", err)
		return makeAuthDenyResponse("newRequest")
	}
	request.Header.Set("cookie", cookeVaule)
	for i, c := range request.Cookies() {
		log.Debugf("check range cookies i:%d, c:%s", i, c.String())
	}
	sidCookie, err := request.Cookie("sid")
	if err != nil {
		log.Errorf("NewRequest err:%v, cookeVaule:%s", err, cookeVaule)
		return makeAuthDenyResponse(fmt.Sprintf("get sid Cookie err:%v", err))
	}
	h, v, err := us.ubiz.SidCheck(sidCookie.Value)
	if err != nil {
		log.Errorf("SidCheck err:%v", err)
		return makeAuthDenyResponse(fmt.Sprintf("SidCheck err:%v", err))
	}

	//return makeNotOkResponse()
	return makeOkResponse(h, v)
}

func makeAuthDenyResponse(msg string) (*authv3.CheckResponse, error) {
	return &authv3.CheckResponse{
		Status: &status.Status{Code: int32(codes.Unauthenticated), Message: msg},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &v3.HttpStatus{Code: v3.StatusCode_Unauthorized},
				Body:   msg,
			},
		},
	}, nil
}

func makeOkResponse(h, val string) (*authv3.CheckResponse, error) {
	return &authv3.CheckResponse{
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{
				Headers: []*corev3.HeaderValueOption{
					{
						Header: &corev3.HeaderValue{
							Key:   h,
							Value: val,
						},
						AppendAction: corev3.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD,
					},
				},
			},
		},
	}, nil
}

func makeNotOkResponse() (*authv3.CheckResponse, error) {
	return &authv3.CheckResponse{
		Status: &status.Status{Code: int32(codes.InvalidArgument), Message: "just a test"},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &v3.HttpStatus{Code: v3.StatusCode(200)},
				Body:   "denied by xshop-user Check",
			},
		},
	}, nil
}
