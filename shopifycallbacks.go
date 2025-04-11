package main

import (
	"github.com/attilakun/crosslist/commongo/shopifyapp"
)

type shopifyCallbacks struct {
	shops map[string]string
}

var _ shopifyapp.ShopifyCallbacks = (*shopifyCallbacks)(nil)

func (s *shopifyCallbacks) HandleShopInstalled(shopDomain string) {
	// s.shops[shopDomain] = nil
}

func (s *shopifyCallbacks) DeleteAccessToken(shopDomain string) error {
	return nil
}

func (s *shopifyCallbacks) GetShop(shopDomain string) (*shopifyapp.Shop, error) {
	return &shopifyapp.Shop{
		Id:          int64(1),
		AccessToken: s.shops[shopDomain],
		ShopDomain:  shopDomain,
	}, nil
}

func (s *shopifyCallbacks) UpsertAccessToken(shop string, accessToken string) error {
	s.shops[shop] = accessToken
	return nil
}
