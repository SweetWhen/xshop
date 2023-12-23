package biz

import (
	"context"
	"strings"
	"testing"
	"time"

	userpb "realworld/api/user/v1"
	"realworld/api/user/v1/common"
	"realworld/app/user/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/durationpb"
)

var _ UserRepo = &myUserRepo{}

type myUserRepo struct {
}

// Create implements UserRepo.
func (mr *myUserRepo) Create(context.Context, *User) (*User, error) {
	panic("unimplemented")
}

// Delete implements UserRepo.
func (mr *myUserRepo) Delete(context.Context, string, int32) error {
	panic("unimplemented")
}

// Get implements UserRepo.
func (mr *myUserRepo) Get(context.Context, string) (*User, error) {
	panic("unimplemented")
}

// ListUser implements UserRepo.
func (mr *myUserRepo) ListUser(ctx context.Context, startId int64, cnt int64, status int) (bus []*User, nextStartId int64, err error) {
	panic("unimplemented")
}

// Update implements UserRepo.
func (mr *myUserRepo) Update(context.Context, *UserUpdate) error {
	panic("unimplemented")
}

var rsaSK = `
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDOUvfYlL7+c6oY5c3Uj84qbJDmX0pL6xycWYxwJUat24FYdAGE
mvIA9JcFf3xwehEEYhJcDmr5saVZkg0x8t1YgRhlz0q6W65AQCH4Ghn4GiWmzcjI
ZolxmflQrlhoPhHMldabuW9wA5WD6UzUaV7IIIb8bWpARAYy/afLgt7dPQIDAQAB
AoGAf9BoF3R2KT1P92KNMwvvBNsCnKQVa5h3vee/l02QTm234CrlMdem6a6by90h
IrCL0DJM+1g3Lv282BMhN3sjYO5ZKRH0H/fyaE1KbAmYgD9fwhu5/fJHxvHXkRLm
TbdD2Gke8dKDaeq9GkyMRDDf5LirRkMXVP2tZM69PGtcImUCQQDyaaBo+hTop8nB
fxSiwVm5exhFdrPPXJq+/bmfXV1RiwVb0heLFjCQnYQi90Dw2xT+17lN5Uvn8P4T
afJ8xfqPAkEA2eOCHlRRDhiKFOmXQ2KE94iKnRppG9VkZIkgIehXv1mM5nkMzRCS
yfXXpfU44+qw895LwM5Oh4NhhRGBhh9BcwJBAMoKGzYjaRXX8qIhJrPX7s5WuA39
NzRW/Gq+0dzvVf3Gnrq+yfyUi/mcLyttZGTaVA9rAPjZaYBxLXJE1WQFJiUCQCZm
YCI0Pey8CmnRGSV5EXIGkFdLtkZ/fyfwusb/CafhgmGD5+ukBhqtxwmqhBI25GS2
QqeCNHjRgLhQ84DNtV0CQQCghkkl6BW9pyq0YM3LqAW2fEKeJGOQATgJYrH1dgpe
NgqTEyKchEresIYfq+XDodNZ9c3Gsn+YxIpczcWT7AHZ
-----END RSA PRIVATE KEY-----
`

var rsaAK = `
-----BEGIN RSA PUBLIC KEY-----
MIGJAoGBAM5S99iUvv5zqhjlzdSPzipskOZfSkvrHJxZjHAlRq3bgVh0AYSa8gD0
lwV/fHB6EQRiElwOavmxpVmSDTHy3ViBGGXPSrpbrkBAIfgaGfgaJabNyMhmiXGZ
+VCuWGg+EcyV1pu5b3ADlYPpTNRpXsgghvxtakBEBjL9p8uC3t09AgMBAAE=
-----END RSA PUBLIC KEY-----
`

func TestJwt(t *testing.T) {
	myRepo := &myUserRepo{}
	confData := &conf.Data{
		RsaPrivate: rsaSK,
		// RsaPublic:  rsaAK,
	}
	confSvr := &conf.Server{
		JwtInterval: durationpb.New(time.Second * 30),
	}
	uc := NewUserUsecase(confSvr, myRepo, log.DefaultLogger, confData)
	ctx := context.Background()
	myClaim := &common.MyClaims{ClaimPayload: userpb.ClaimPayload{Uid: 3, Name: "superAdmin"}}
	jwtID, err := uc.Checkin(ctx, &myClaim.ClaimPayload)
	if err != nil {
		t.Fatalf("Checkin err:%v\n", err)
	}
	h, value, err := uc.SidCheck(jwtID)
	if err != nil {
		t.Fatalf("SidCheck err:%v\n", err)
	}
	vales := strings.Split(value, " ")
	t.Logf("after sidCheck header:%s, vales:%v\n", h, vales)
	claims, err := uc.JwtCheck(vales[1])
	if err != nil {
		t.Fatalf("JwtCheck err:%v", err)
	}

	t.Logf("claims: %+v, %s\n", claims.RegisteredClaims, claims.ClaimPayload.String())
}
