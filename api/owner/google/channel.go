package google

import (
	"backend/infrastructure/owner/dao"
	pb "backend/proto/owner"
	"backend/tool"
	"context"
	"net/http"
)

func (receiver *OwnerService) GetChannel(_ context.Context, req *pb.ChannelRequest) (*pb.ChannelResponse, error) {
	userDB := dao.Original{ID: req.Id}
	result, err := userDB.Read()
	if dao.IsEmpty(err) {
		tool.Logger().Warning("not found channel id", err, "channel id", req.Id)
		return &pb.ChannelResponse{Status: &pb.OwnerStatus{
			Code:    http.StatusNotFound,
			Message: "not found",
		}}, nil
	}
	if err != nil {
		tool.Logger().Warning("fail read db by GetChannel", err)
		return &pb.ChannelResponse{
			Status: &pb.OwnerStatus{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			},
		}, err
	}
	return &pb.ChannelResponse{
		Status: &pb.OwnerStatus{
			Code:    http.StatusOK,
			Message: "internal server error",
		},
		Title:       result.Channel.Name,
		Description: result.Channel.Description,
		Image:       result.Channel.Image,
		Url:         result.Channel.Url,
		Address:     result.Address,
	}, nil
}
