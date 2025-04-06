package main

import (
	"fmt"
	"net/http"

	"github.com/attilakun/crosslist/commongo"
	"github.com/attilakun/crosslist/commongo/shopifyapp"
	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	godotenv.Load()
	commongo.InitLog()
	logger := log.Logger
	shopifyAppSettings := shopifyapp.ShopifyAppSettings{
		UseStaticFrontend: commongo.GetEnvVariable(logger, "USE_STATIC_FRONTEND") == "true",
		ShopifyAppBaseUrl: commongo.GetEnvVariable(logger, "SHOPIFY_APP_BASE_URL"),
	}
	shopifyApp := &goshopify.App{
		ApiKey:      commongo.GetEnvVariable(logger, "SHOPIFY_KEY"),
		ApiSecret:   commongo.GetEnvVariable(logger, "SHOPIFY_SECRET"),
		RedirectUrl: fmt.Sprintf("%s/auth/callback", shopifyAppSettings.ShopifyAppBaseUrl),
		Scope:       "read_products,write_products",
	}

	shopifyCallback := &shopifyCallbacks{
		handleShopInstalled: func(shopDomain string) {
		},
	}

	router := gin.Default()
	// router.LoadHTMLGlob("templates/*")
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})

	shopifyapp.InitShopifyApp(
		router,
		shopifyApp,
		shopifyAppSettings,
		shopifyCallback,
		getLogPorcessedShopifyRoute("shopifyHandler"),
	)

	shopifyapp.InitFrontendHandler(
		router,
		getLogPorcessedShopifyRoute("frontendHandler"),
		shopifyAppSettings,
		shopifyCallback,
	)

	port := commongo.GetEnvVariable(log.Logger, "PORT")
	portStr := ":" + port

	if commongo.GetEnvVariable(log.Logger, "SERVE_TLS") == "true" {
		router.RunTLS(
			portStr,
			commongo.GetEnvVariable(log.Logger, "TLS_CERT_FILE_PATH"),
			commongo.GetEnvVariable(log.Logger, "TLS_KEY_FILE_PATH"),
		)
	} else {
		router.Run(portStr)
	}
}

func getLogPorcessedShopifyRoute(str string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info().Msgf("Processed Shopify %s route: %s", str, c.Request.URL.Path)
	}
}
