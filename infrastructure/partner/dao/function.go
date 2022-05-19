package dao

import (
	"github.com/kamva/mgm/v3"
	"github.com/pkg/errors"
)

func NewNftInfo(name, description, uri string) (*NFT, error) {
	nft := &NFT{
		Name:        name,
		Description: description,
		ImageUri:    uri,
	}
	err := mgm.Coll(nft).Create(nft)
	if err != nil {
		return nil, errors.Wrap(err, "NewNFTInfo")
	}
	return nft, nil

}

func FindId(objectId string) (*NFT, error) {
	nft := &NFT{}
	err := mgm.Coll(&NFT{}).FindByID(objectId, nft)
	if err != nil {
		return nil, errors.Wrap(err, "FindId")
	}
	return nft, nil
}
