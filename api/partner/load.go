package partner

import (
	"backend/infrastructure/partner/dao"
	pb "backend/proto/partner"
	"context"
	"github.com/pkg/errors"
	"net/http"
)

func (receiver PartnerService) LoadNFTInfo(_ context.Context, req *pb.LoadRequest) (*pb.LoadResponse, error) {
	nft, err := dao.FindId(req.Id)
	if dao.IsEmpty(err) {
		return &pb.LoadResponse{
			Status: &pb.Status{
				Code:    http.StatusNotFound,
				Message: "not found id",
			},
		}, nil
	} else if err != nil {
		return &pb.LoadResponse{
			Status: &pb.Status{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			},
		}, errors.Wrap(err, "LoadNFTInfo")
	}
	return &pb.LoadResponse{
		Name:        nft.Name,
		Description: nft.Description,
		Url:         nft.ImageUri,
		Status: &pb.Status{
			Code:    http.StatusOK,
			Message: "success",
		},
	}, nil
}
