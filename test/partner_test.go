package test

import (
	"backend/infrastructure/partner/dao"
	"github.com/kamva/mgm/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddNFT(t *testing.T) {
	nft1, err := dao.NewNftInfo("test1", "test2", "url")
	defer mgm.Coll(&dao.NFT{}).Delete(nft1)
	assert.Nil(t, err, "등록중 에러발생", err)
	id := nft1.ID.Hex()
	nft2, err := dao.FindId(id)
	assert.Nil(t, err, "조회중 에러발생", err)
	assert.Equal(t, nft1.Name, nft2.Name, "서로다른 nft")
}

func TestNotFoundID(t *testing.T) {
	_, err := dao.FindId("62866313e33c7fd8bb7414ac")
	assert.True(t, dao.IsEmpty(err), "no document에러가 아님", err)

}
