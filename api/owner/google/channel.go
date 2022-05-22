package google

import (
	"backend/infrastructure/owner/dao"
	"backend/internal/owner"
	pb "backend/proto/owner"
	"context"
	"encoding/hex"
	"net/http"
	"strings"
)

func (receiver *OwnerService) SetAnnoymousUser(_ context.Context, req *pb.NicknameRequest) (*pb.NicknameResponse, error) {
	user := owner.User{
		Address:  req.Address,
		Nickname: req.Nickname,
	}
	err := owner.Validate(&user)
	if err != nil {
		return &pb.NicknameResponse{
			Status: &pb.OwnerStatus{
				Code: http.StatusBadRequest,
			},
		}, nil
	}
	userDB := owner.User2UserDB(user)
	err = userDB.Read()
	if dao.IsEmpty(err) {
		err = userDB.Save()
		if err != nil {
			return &pb.NicknameResponse{
				Status: &pb.OwnerStatus{
					Code:    http.StatusInternalServerError,
					Message: "fail save",
				},
			}, nil
		}
		return &pb.NicknameResponse{
			Status: &pb.OwnerStatus{
				Code:    http.StatusOK,
				Message: "success",
			},
		}, nil
	} else if err != nil {
		return &pb.NicknameResponse{
			Status: &pb.OwnerStatus{
				Code:    http.StatusInternalServerError,
				Message: "db error",
			},
		}, nil
	}
	return &pb.NicknameResponse{
		Status: &pb.OwnerStatus{
			Code:    http.StatusForbidden,
			Message: "exist addr or nickname",
		},
	}, nil
}

func (receiver OwnerService) GetAnnoymousUser(_ context.Context, req *pb.NicknameRequest) (*pb.NicknameResponse, error) {
	addr := strings.TrimPrefix(req.Address, "0x")
	byteArray, err := hex.DecodeString(addr)
	if err != nil || len(byteArray) != 20 {
		return &pb.NicknameResponse{
			Status: &pb.OwnerStatus{
				Code:    http.StatusBadRequest,
				Message: "Not Address",
			},
		}, nil
	}
	user := owner.User{Address: req.Address}
	userDB := owner.User2UserDB(user)
	err = userDB.Read()
	if dao.IsEmpty(err) {
		return FindChannelName(req.Address)
	} else if err != nil {
		return &pb.NicknameResponse{
			Status: &pb.OwnerStatus{
				Code:    http.StatusInternalServerError,
				Message: "db error",
			},
		}, nil
	}
	return &pb.NicknameResponse{
		Status: &pb.OwnerStatus{
			Code: http.StatusOK,
		},
		Nickname: userDB.Nickname,
	}, nil
}

func FindChannelName(address string) (*pb.NicknameResponse, error) {
	channeOwner := dao.Original{
		Address: address,
	}
	result, err := channeOwner.Read()
	if dao.IsEmpty(err) {
		return &pb.NicknameResponse{
			Status: &pb.OwnerStatus{
				Code: http.StatusNotFound,
			},
		}, nil
	} else if err != nil {
		return &pb.NicknameResponse{
			Status: &pb.OwnerStatus{
				Code: http.StatusInternalServerError,
			},
		}, nil
	}
	return &pb.NicknameResponse{
		Status: &pb.OwnerStatus{
			Code: http.StatusOK,
		},
		Nickname: result.Channel.Name,
	}, nil
}
