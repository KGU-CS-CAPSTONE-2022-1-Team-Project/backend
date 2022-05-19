package partner

import (
	"backend/infrastructure/partner/dao"
	pb "backend/proto/partner"
	"context"
	"net/http"
)

func (receiver *PartnerService) SaveNFTInfo(_ context.Context, req *pb.SaveRequest) (*pb.SaveResponse, error) {
	if req.Name == "" || req.Description == "" || req.File == nil {
		return &pb.SaveResponse{
			Status: &pb.Status{
				Code:    http.StatusBadRequest,
				Message: "not found value",
			},
		}, nil
	}
	fileId, err := UploadObject(req.File.Chunk)
	if err != nil || fileId == "" {
		return &pb.SaveResponse{
			Status: &pb.Status{
				Code:    http.StatusBadRequest,
				Message: "not found value",
			},
		}, nil
	}
	nft, err := dao.NewNftInfo(
		req.Name,
		req.Description,
		"https://kr.object.ncloudstorage.com/nft-image/"+fileId)
	if err != nil {
		return &pb.SaveResponse{
			Status: &pb.Status{
				Code:    http.StatusInternalServerError,
				Message: "db error",
			},
		}, nil
	}
	return &pb.SaveResponse{
		Id:      nft.ID.Hex(),
		Success: true,
		Status: &pb.Status{
			Code:    http.StatusOK,
			Message: "success",
		},
	}, nil
}
