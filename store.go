package main

import (
	"github.com/attilakun/crosslist/commongo/shopifyapp"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// type Shop struct {
// 	Id          int64  `db:"id" json:"id"`
// 	AccessToken string `db:"access_token" json:"access_token"`
// 	ShopDomain  string `db:"shop_domain" json:"shop_domain"`
// }

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

var _ shopifyapp.ShopifyCallbacks = (*Store)(nil)

func (s *Store) HandleShopInstalled(shopDomain string) {
	// Log when a shop is installed
	log.Info().Str("shop_domain", shopDomain).Msg("Shop installed")
}

func (s *Store) DeleteAccessToken(shopDomain string) error {
	// Delete the access token for a shop from the database
	_, err := s.db.Exec("UPDATE shop SET access_token = '' WHERE shop_domain = $1", shopDomain)
	if err != nil {
		log.Error().Err(err).Str("shop_domain", shopDomain).Msg("Failed to delete access token")
		return err
	}
	return nil
}

func (s *Store) GetShop(shopDomain string) (*shopifyapp.Shop, error) {
	// Get shop information from the database
	var shop shopifyapp.Shop
	err := s.db.Get(&shop, "SELECT id, access_token, shop_domain FROM shop WHERE shop_domain = $1", shopDomain)
	if err != nil {
		log.Error().Err(err).Str("shop_domain", shopDomain).Msg("Failed to get shop")
		return nil, err
	}
	return &shop, nil
}

func (s *Store) UpsertAccessToken(shop string, accessToken string) error {
	// Insert or update the access token for a shop
	_, err := s.db.Exec(`
		INSERT INTO shop (shop_domain, access_token) 
		VALUES ($1, $2) 
		ON CONFLICT (shop_domain) 
		DO UPDATE SET access_token = $2
	`, shop, accessToken)

	if err != nil {
		log.Error().Err(err).Str("shop_domain", shop).Msg("Failed to upsert access token")
		return err
	}
	return nil
}
