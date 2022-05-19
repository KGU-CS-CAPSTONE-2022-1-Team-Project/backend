package dao

import "github.com/kamva/mgm/v3"

type NFT struct {
	mgm.DefaultModel
	Name        string
	Description string
	ImageUri    string
}
