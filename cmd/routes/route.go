package routes

import (
	"api_shope/internal/handler"
	"api_shope/utils/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(auth *handler.AuthHandler, shop *handler.ShopHandler) *mux.Router {
	r := mux.NewRouter()

	//auth
	r.HandleFunc("/login", auth.Login).Methods(http.MethodPost)
	r.HandleFunc("/register", auth.Register).Methods(http.MethodPost)

	//shop
	useM := r.PathPrefix("/shop").Subrouter()
	useM.Use(middleware.AuthMiddleware)

	//s_store
	useM.HandleFunc("/get-store", shop.GetAllStore).Methods(http.MethodGet)
	useM.HandleFunc("/get-my-store/{storeId}", shop.GetMyStore).Methods(http.MethodGet)
	useM.HandleFunc("/create", shop.CreateStore).Methods(http.MethodPost)
	useM.HandleFunc("/update/{storeId}", shop.UpdateStore).Methods(http.MethodPut)
	useM.HandleFunc("/delete/{storeId}", shop.DeleteStore).Methods(http.MethodDelete)

	//s_product
	productRouter := useM.PathPrefix("/product").Subrouter()

	productRouter.HandleFunc("/get-all-product", shop.GetAllProduct).Methods(http.MethodGet)
	productRouter.HandleFunc("/get-product/{productId}", shop.GetThisProduct).Methods(http.MethodGet)
	productRouter.HandleFunc("/create/{storeId}", shop.CreateProduct).Methods(http.MethodPost)
	productRouter.HandleFunc("/update/{productId}", shop.UpdateProduct).Methods(http.MethodPut)
	productRouter.HandleFunc("/delete/{storeId}/{productId}", shop.DeleteProduct).Methods(http.MethodDelete)

	//s_cart item
	cartItemRouter := useM.PathPrefix("/cart-item").Subrouter()

	cartItemRouter.HandleFunc("/get-my-cart-item", shop.GetMyCartItems).Methods(http.MethodGet)
	cartItemRouter.HandleFunc("/create/{productId}", shop.CreateCartItem).Methods(http.MethodPost)
	cartItemRouter.HandleFunc("/update-amount/{cartItemId}/{productId}", shop.UpdateAmountCartItem).Methods(http.MethodPut)
	cartItemRouter.HandleFunc("/update-paid/{cartItemId}/{productId}", shop.UpdatePaidCartItem).Methods(http.MethodPut)
	cartItemRouter.HandleFunc("/delete/{cartItemId}", shop.DeleteCartItem).Methods(http.MethodDelete)

	return r
}
