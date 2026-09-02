package routes

import "github.com/gin-gonic/gin"

const scalarReferenceHTML = `<!doctype html>
<html lang="pt-BR">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Cidadon API Reference</title>
  </head>
  <body>
    <div id="app"></div>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
    <script>
      Scalar.createApiReference('#app', {
        url: '/openapi.yaml',
        theme: 'purple',
        layout: 'modern',
        hideClientButton: false,
        pageTitle: 'Cidadon API',
      })
    </script>
  </body>
</html>`

// RegisterOpenAPI exposes the source contract and a modern interactive reference.
func RegisterOpenAPI(router *gin.Engine, spec []byte) {
	router.GET("/openapi.yaml", func(c *gin.Context) {
		c.Data(200, "application/yaml; charset=utf-8", spec)
	})
	router.GET("/docs", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", []byte(scalarReferenceHTML))
	})
}
