package routes

import (
	"ventas-app/controllers"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.POST("/login", controllers.Login)
	//r.POST("/usuarios", middleware.AuthRequired("precio", "comprador"), controllers.CrearUsuario)
	r.POST("/usuarios", controllers.CrearUsuario)

	r.GET("/productos", controllers.ListarProductos)
	//r.POST("/productos", middleware.AuthRequired("vendedor", "comprador"), controllers.CrearProducto)
	r.POST("/productos", controllers.CrearProducto)

	//r.POST("/compras", middleware.AuthRequired("comprador", "vendedor"), controllers.RegistrarCompra)
	//r.POST("/ventas", middleware.AuthRequired("vendedor"), controllers.RegistrarVenta)

	r.POST("/compras", controllers.RegistrarCompra)
	r.POST("/ventas", controllers.RegistrarVenta)
}
