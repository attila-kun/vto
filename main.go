package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/attila-kun/vto/frontend"
	"github.com/attilakun/crosslist/commongo"
	"github.com/attilakun/crosslist/commongo/shopifyapp"
	goshopify "github.com/bold-commerce/go-shopify/v3"
	esbuildapi "github.com/evanw/esbuild/pkg/api"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	godotenv.Load()
	commongo.InitLog()
	logger := log.Logger
	frontendDevPort, err := strconv.ParseUint(commongo.GetEnvVariable(logger, "SHOPIFY_FRONTEND_DEV_PORT"), 10, 16)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse frontend dev port")
	}
	shopifyAppSettings := shopifyapp.ShopifyAppSettings{
		UseStaticFrontend: commongo.GetEnvVariable(logger, "USE_STATIC_FRONTEND") == "true",
		ShopifyAppBaseUrl: commongo.GetEnvVariable(logger, "SHOPIFY_APP_BASE_URL"),
		FrontendDevPort:   uint16(frontendDevPort),
	}
	shopifyApp := &goshopify.App{
		ApiKey:      commongo.GetEnvVariable(logger, "SHOPIFY_KEY"),
		ApiSecret:   commongo.GetEnvVariable(logger, "SHOPIFY_SECRET"),
		RedirectUrl: fmt.Sprintf("%s/auth/callback", shopifyAppSettings.ShopifyAppBaseUrl),
		Scope:       "read_products,write_products",
	}

	shopifyCallback := &shopifyCallbacks{
		shops: make(map[string]string),
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

	setupApi(router, shopifyApp.ApiSecret)

	initFrontend(shopifyAppSettings)

	shopifyapp.InitFrontendHandler(
		router,
		getLogPorcessedShopifyRoute("frontendHandler"),
		shopifyAppSettings,
		shopifyCallback,
		frontend.Index(shopifyApp.ApiKey),
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

func setupApi(
	router *gin.Engine,
	apiSecret string,
) {
	apiGroup := router.Group(
		"/api",
		shopifyapp.HandleAuthToken(apiSecret),
	)

	apiGroup.POST("/test", func(c *gin.Context) {
		user := shopifyapp.GetUserFromContext(c)
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello, %d!", user.Id)})
	})
}

func initFrontend(shopifyAppSettings shopifyapp.ShopifyAppSettings) {
	ctx, esbuildErr := esbuildapi.Context(esbuildapi.BuildOptions{
		EntryPoints: []string{"frontend/src/main.tsx"},
		Outdir:      "frontend/dist",
		Bundle:      true,
		JSX:         esbuildapi.JSXAutomatic,
	})
	if esbuildErr != nil {
		log.Fatal().Err(esbuildErr).Msgf("Failed to create esbuild context: %s", esbuildErr)
	}

	_, err2 := ctx.Serve(esbuildapi.ServeOptions{
		Port: uint16(shopifyAppSettings.FrontendDevPort),
		Host: "localhost",
	})
	if err2 != nil {
		log.Fatal().Err(err2).Msgf("Failed to serve frontend: %s", err2)
	}
}
