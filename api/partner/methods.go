package partner

import (
	"backend/infrastructure/partner/dao"
	pb "backend/proto/partner"
	"backend/tool"
	"context"
	"net/http"
)

func (receiver *PartnerService) SaveNFTInfo(_ context.Context, req *pb.SaveRequest) (*pb.SaveResponse, error) {
	tool.Logger().Info("save nft", "token name", req.Name)
	if req.Name == "" || req.Description == "" || req.File == nil {

		return &pb.SaveResponse{
			Status: &pb.Status{
				Code:    http.StatusBadRequest,
				Message: "not found value",
			},
		}, nil
	}
	fileId, err := UploadObject(req.File.Chunk)
	if err != nil {
		tool.Logger().Warning("fail save", err, "fileId", fileId)
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
		tool.Logger().Warning("fail save nft", err, "fileId", fileId)
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

func (receiver PartnerService) LoadNFTInfo(_ context.Context, req *pb.LoadRequest) (*pb.LoadResponse, error) {
	tool.Logger().Info("load nft info", "token id", req.Id)
	nft, err := dao.FindId(req.Id)
	if dao.IsEmpty(err) {
		return &pb.LoadResponse{
			Status: &pb.Status{
				Code:    http.StatusNotFound,
				Message: "not found id",
			},
		}, nil
	} else if err != nil {
		tool.Logger().Warning("fail save", err, "token id", req.Id)
		return &pb.LoadResponse{
			Status: &pb.Status{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			},
		}, nil
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
