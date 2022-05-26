package google

import (
	"backend/infrastructure/owner/dao"
	"backend/internal/owner"
	"backend/internal/owner/blockchain"
	pb "backend/proto/owner"
	"backend/tool"
	"context"
	"encoding/hex"
	"net/http"
	"strings"
)

func (receiver *OwnerService) SetAnnoymousUser(_ context.Context, req *pb.NicknameRequest) (*pb.NicknameResponse, error) {
	addr, err := blockchain.UnsignedAddress(req.Address)
	if err != nil {
		tool.Logger().Error("fail unsigned", err, "request signed addr", req.Address)
		return &pb.NicknameResponse{Status: &pb.OwnerStatus{
			Code:    http.StatusInternalServerError,
			Message: "fail unsigned",
		}}, nil
	}
	tool.Logger().Info("set nickname", "address", addr)
	user := owner.User{
		Address:  addr,
		Nickname: req.Nickname,
	}
	err = owner.Validate(&user)
	if err != nil {
		tool.Logger().Warning("invalidate addr", err, "address", addr, "nickname", req.Nickname)
		return &pb.NicknameResponse{
			Status: &pb.OwnerStatus{
				Code: http.StatusBadRequest,
			},
		}, nil
	}
	userDB := owner.User2UserDB(user)
	err = userDB.Read()
	if dao.IsEmpty(err) {
		userDB = owner.User2UserDB(user)
		err = userDB.Save()
		if err != nil {
			tool.Logger().Warning("fail save addr", err, "address", addr, "nickname", req.Nickname)
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
		tool.Logger().Warning("fail search addr or nickname", err, "address", addr, "nickname", req.Nickname)
		return &pb.NicknameResponse{
			Status: &pb.OwnerStatus{
				Code:    http.StatusInternalServerError,
				Message: "db error",
			},
		}, nil
	}
	tool.Logger().Warning("found addr", err, "address",
		addr, "searched addr", userDB.Address, "nickname", req.Nickname, "searched addr", userDB.Nickname)
	return &pb.NicknameResponse{
		Status: &pb.OwnerStatus{
			Code:    http.StatusForbidden,
			Message: "exist addr or nickname",
		},
	}, nil
}

func (receiver *OwnerService) GetAnnoymousUser(_ context.Context, req *pb.NicknameRequest) (*pb.NicknameResponse, error) {
	tool.Logger().Info("get nickname", "address", req.Address)
	addr := strings.TrimPrefix(req.Address, "0x")
	byteArray, err := hex.DecodeString(addr)
	if err != nil || len(byteArray) != 20 {
		tool.Logger().Error("not address", err, "address", req.Address)
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
		tool.Logger().Warning("fail read orignal db", err, "address", req.Address)
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
		tool.Logger().Info("empty address", "address", address)
		return &pb.NicknameResponse{
			Status: &pb.OwnerStatus{
				Code:    http.StatusNotFound,
				Message: "not found",
			},
		}, nil
	} else if err != nil {
		tool.Logger().Warning("fail read user db", err, "address", address)
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
