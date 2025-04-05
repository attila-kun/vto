package main

import (
	"github.com/attilakun/crosslist/commongo/shopifyapp"
)

type shopifyCallbacks struct {
	handleShopInstalled func(shopDomain string)
}

var _ shopifyapp.ShopifyCallbacks = (*shopifyCallbacks)(nil)

func (s *shopifyCallbacks) HandleShopInstalled(shopDomain string) {
	s.handleShopInstalled(shopDomain)
}

func (s *shopifyCallbacks) DeleteAccessToken(shopDomain string) error {
	return nil
}

func (s *shopifyCallbacks) GetShop(shopDomain string) (*shopifyapp.Shop, error) {
	return nil, nil
}

func (s *shopifyCallbacks) UpsertAccessToken(shop string, accessToken string) error {
	return nil
}
