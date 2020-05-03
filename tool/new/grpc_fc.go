package new

var (
	grpcfc1 = &FileContent{
		FileName: "demo_controller.go",
		Dir:      "internal/transports/grpc/controllers",
		Content: `package controllers

import (
	"context"
	gp "{{.ProPath}}{{.ServerName}}/internal/infra/third_party/protobuf/passport"
	"{{.ProPath}}{{.ServerName}}/internal/application"
	"{{.ProPath}}{{.ServerName}}/internal/infra"
	"{{.ProPath}}{{.ServerName}}/internal/transports/grpc/dto"
)

type DemoController struct {
	infra *infra.Infra

	userSvc *application.UserService
}


func (this *DemoController) GetUserByUserName(ctx context.Context,
	request *gp.GetUserByUserNameRequest) (*gp.GrpcUserReply, error) {
	grpcReply = &gp.GrpcUserReply{}
	userName = request.GetUsername()

	userInfo = this.userSvc.GetUserInfo(ctx, userName)

	grpcReply.Code = 0;

	grpcReply.Data = dto.NewUserInfo(userInfo)

	return grpcReply, nil
}
`,
	}

	grpcfc2 = &FileContent{
		FileName: "routers.go",
		Dir:      "internal/transports/grpc/routers",
		Content: `package routers

import (
	"{{.ProPath}}{{.ServerName}}/internal/transports/grpc/controllers"
	"{{.ProPath}}{{.ServerName}}/internal/infra/third_party/protobuf/passport"
	"google.golang.org/grpc"
)


func RegisterGrpcServer(s *grpc.Server, controllers *controllers.Controllers)  {
	passport.RegisterUserInfoServer(s, controllers.Demo)
}
`,
	}

	grpcfc3 = &FileContent{
		FileName: "grpc.go",
		Dir:      "internal/transports/grpc",
		Content: `package grpc

import (
	"strings"

	"github.com/jukylin/esim/grpc"
	{{.PackageName}} "{{.ProPath}}{{.ServerName}}/internal"
	"{{.ProPath}}{{.ServerName}}/internal/transports/grpc/routers"
	"{{.ProPath}}{{.ServerName}}/internal/transports/grpc/controllers"
)

func NewGrpcServer(app *{{.PackageName}}.App) *grpc.GrpcServer {

	target = app.Conf.GetString("grpc_server_tcp")

	in = strings.Index(target, ":")
	if in < 0 {
		target = ":"+target
	}

	serverOptions = grpc.ServerOptions{}

	//grpc服务初始化
	grpcServer =  grpc.NewGrpcServer(target,
		serverOptions.WithServerConf(app.Conf),
		serverOptions.WithServerLogger(app.Logger),
		serverOptions.WithUnarySrvItcp(),
		serverOptions.WithGrpcServerOption(),
		serverOptions.WithTracer(app.Tracer),
	)

	//注册grpc路由
	routers.RegisterGrpcServer(grpcServer.Server, controllers.NewControllers(app))

	return grpcServer
}
`,
	}

	grpcfc4 = &FileContent{
		FileName: "component_test.go",
		Dir:      "internal/transports/grpc/component-test",
		Content: `package component_test

import (
	"context"
	"testing"

	egrpc "github.com/jukylin/esim/grpc"
	"github.com/jukylin/esim/log"
	"github.com/stretchr/testify/assert"
	gp "{{.ProPath}}{{.ServerName}}/internal/infra/third_party/protobuf/passport"
)

//go test
func TestUserService_GetUserByUserName(t *testing.T) {
	logger = log.NewLogger()

	ctx = context.Background()

	grpcClient = egrpc.NewClient(egrpc.NewClientOptions())
	conn = grpcClient.DialContext(ctx, ":50055")
	defer conn.Close()

	client = gp.NewUserInfoClient(conn)

	req = &gp.GetUserByUserNameRequest{}
	req.Username = "demo"
	reply, err = client.GetUserByUserName(ctx, req)
	if err != nil {
		logger.Errorf(err.Error())
	} else {
		assert.Equal(t, "demo", reply.Data.UserName)
		assert.Equal(t, int32(0), reply.Code)
	}
}`,
	}

	grpcfc5 = &FileContent{
		FileName: "controllers.go",
		Dir:      "internal/transports/grpc/controllers",
		Content: `package controllers

import (
	{{.PackageName}} "{{.ProPath}}{{.ServerName}}/internal"
	"github.com/google/wire"
	"{{.ProPath}}{{.ServerName}}/internal/application"
)


type Controllers struct {

	App *{{.PackageName}}.App

	Demo *DemoController
}


var controllersSet = wire.NewSet(
	wire.Struct(new(Controllers), "*"),
	provideDemoController,
)


func NewControllers(app *{{.PackageName}}.App) *Controllers {
	controllers = initControllers(app)
	return controllers
}


func provideDemoController(app *{{.PackageName}}.App) *DemoController {

	userSvc = application.NewUserSvc(app.Infra)

	demoController = &DemoController{}
	demoController.infra = app.Infra
	demoController.userSvc = userSvc

	return demoController
}
`,
	}

	grpcfc6 = &FileContent{
		FileName: "wire.go",
		Dir:      "internal/transports/grpc/controllers",
		Content: `//+build wireinject

package controllers

import (
	"github.com/google/wire"
	{{.PackageName}} "{{.ProPath}}{{.ServerName}}/internal"
)



func initControllers(app *{{.PackageName}}.App) *Controllers {
	wire.Build(controllersSet)
	return nil
}
`,
	}

	grpcfc7 = &FileContent{
		FileName: "wire_gen.go",
		Dir:      "internal/transports/grpc/controllers",
		Content: `// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package controllers

import (
	{{.PackageName}} "{{.ProPath}}{{.ServerName}}/internal"
)

// Injectors from wire.go:

func initControllers(app *{{.PackageName}}.App) *Controllers {
	demoController = provideDemoController(app)
	controllers = &Controllers{
		App:  app,
		Demo: demoController,
	}
	return controllers
}
`,
	}

	grpcfc8 = &FileContent{
		FileName: "user_dto.go",
		Dir:      "internal/transports/grpc/dto",
		Content: `package dto

import (
	"{{.ProPath}}{{.ServerName}}/internal/domain/user/entity"
	"{{.ProPath}}{{.ServerName}}/internal/infra/third_party/protobuf/passport"
)

type User struct {

	//用户名称
	UserName string {{.SingleMark}}json:"user_name"{{.SingleMark}}

	//密码
	PassWord string {{.SingleMark}}json:"pass_word"{{.SingleMark}}
}

func NewUserInfo(user entity.User) *passport.Info {
	info = &passport.Info{}
	info.UserName = user.UserName
	info.PassWord = user.PassWord
	return info
}`,
	}

	grpcfc9 = &FileContent{
		FileName: "main_test.go",
		Dir:      "internal/transports/grpc/component-test",
		Content: `package component_test

import (
	"os"
	"testing"
	"context"

	"{{.ProPath}}{{.ServerName}}/internal/transports/grpc"
	"{{.ProPath}}{{.ServerName}}/internal/infra"
	_grpc "google.golang.org/grpc"
	egrpc "github.com/jukylin/esim/grpc"
	"github.com/jukylin/esim/container"
	{{.PackageName}} "{{.ProPath}}{{.ServerName}}/internal"
)

func TestMain(m *testing.M) {
	appOptions = {{.PackageName}}.AppOptions{}
	app = {{.PackageName}}.NewApp(appOptions.WithConfPath("../../../../conf/"))

	setUp(app)

	code = m.Run()

	tearDown(app)

	os.Exit(code)
}

func provideStubsGrpcClient(esim *container.Esim) *egrpc.GrpcClient {
	clientOptional = egrpc.ClientOptionals{}
	clientOptions = egrpc.NewClientOptions(
		clientOptional.WithLogger(esim.Logger),
		clientOptional.WithConf(esim.Conf),
		clientOptional.WithDialOptions(_grpc.WithUnaryInterceptor(
			egrpc.ClientStubs(func(ctx context.Context, method string, req, reply interface{}, cc *_grpc.ClientConn, invoker _grpc.UnaryInvoker, opts ..._grpc.CallOption) error {
				esim.Logger.Infof(method)
				err = invoker(ctx, method, req, reply, cc, opts...)
				return err
			}),
		)),
	)

	grpcClient = egrpc.NewClient(clientOptions)

	return grpcClient
}

func setUp(app *{{.PackageName}}.App) {

	app.Infra = infra.NewStubsInfra(provideStubsGrpcClient(app.Esim))

	app.Trans = append(app.Trans, grpc.NewGrpcServer(app))

	app.Start()

	errs = app.Infra.HealthCheck()
	if len(errs) > 0 {
		for _, err = range errs {
			app.Logger.Errorf(err.Error())
		}
	}
}

func tearDown(app *{{.PackageName}}.App) {
	app.Infra.Close()
}`,
	}
)

func GrpcInit() {
	Files = append(Files, grpcfc1, grpcfc2, grpcfc3, grpcfc4, grpcfc5, grpcfc6, grpcfc7, grpcfc8, grpcfc9)
}
