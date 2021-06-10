package api

import (
	"context"

	"gitlab.com/siriusfreak/lecture-7-demo/common"
	"google.golang.org/protobuf/types/known/emptypb"

	textService "gitlab.com/siriusfreak/lecture-7-demo/internal/text_service"
	desc "gitlab.com/siriusfreak/lecture-7-demo/pkg/lecture-7-demo"
)

type api struct {
	desc.UnimplementedLecture7DemoServer
	textService textService.Service
}


func NewLecture7DemoAPI() desc.Lecture7DemoServer {
	srv, err := textService.InitService()

	if err != nil {
		panic(err)
	}

	return &api{
		textService: srv,
	}
}

func (a *api)AddV1(ctx context.Context, req *desc.AddRequestV1) (*emptypb.Empty, error) {
	common.IncProcessedByHandler("AddV1")

	err := a.textService.AddV1(ctx, req.Id, req.Text, req.Result, req.CallbackUrl)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, err
}

func (a *api)CallbackFirstV1(ctx context.Context, req *desc.CallbackFirstV1Request) (*emptypb.Empty, error) {
	common.IncProcessedByHandler("CallbackFirstV1")

	err := a.textService.CallbackFirstV1(ctx, req.Id, req.Result)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, err

}
func (a *api)CallbackSecondV1(ctx context.Context, req *desc.CallbackSecondV1Request) (*emptypb.Empty, error) {
	common.IncProcessedByHandler("CallbackSecondV1")

	err := a.textService.CallbackSecondV1(ctx, req.Id, req.Result)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, err

}
func (a *api)StatusV1(ctx context.Context, req *emptypb.Empty) (*desc.StatusResponseV1, error) {
	common.IncProcessedByHandler("StatusV1")

	states, err := a.textService.StatusV1(ctx)
	if err != nil {
		return nil, err
	}

	resArr := make([]*desc.StatusResponseV1_Status, 0)
	for key, val := range states {
		resArr = append(resArr, &desc.StatusResponseV1_Status{
			Id:      key,
			Correct: val,
		})
	}
	resp := &desc.StatusResponseV1{
		StatusList: resArr,
	}

	return resp, nil
}