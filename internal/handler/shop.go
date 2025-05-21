package handler

import (
	"api_shope/dto"
	"api_shope/internal/usecase"
	"api_shope/utils/helper"
	"api_shope/utils/middleware"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ShopHandler struct {
	shopUsecase usecase.ShopUsecase
}

func NewShopHandler(shopUsecase usecase.ShopUsecase) *ShopHandler {
	return &ShopHandler{shopUsecase}
}

// store
func (h *ShopHandler) GetMyStore(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	claims, ok := claimsRaw.(*helper.JWTCLAIMS)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "protected api")
		return
	}

	params := mux.Vars(r)
	paramsId, err := strconv.Atoi(params["storeId"])
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "params tidak ditemukan")
		return
	}

	response, err := h.shopUsecase.GetMyStore(claims.UserID, uint(paramsId))
	if err != nil {
		switch err {
		case helper.ErrNotAdmin:
			helper.WriteError(w, http.StatusUnauthorized, "bukan admin")
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, response)
}

func (h *ShopHandler) GetAllStore(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	_, ok := claimsRaw.(*helper.JWTCLAIMS)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "protected api")
		return
	}

	response, err := h.shopUsecase.GetAllStore()
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, response)
}

func (h *ShopHandler) CreateStore(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	claims, ok := claimsRaw.(*helper.JWTCLAIMS)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "protected api")
		return
	}

	var req dto.CreateStoreReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	req.AdminID = claims.UserID
	if err := h.shopUsecase.CreateStore(&req); err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, nil)
}

func (h *ShopHandler) UpdateStore(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	claims, ok := claimsRaw.(*helper.JWTCLAIMS)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "protected api")
		return
	}

	var req dto.UpdateStoreReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	params := mux.Vars(r)
	paramsStoreId, err := strconv.Atoi(params["storeId"])
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "params tidak ditemukan")
		return
	}

	req.ID = uint(paramsStoreId)
	req.UserID = claims.UserID
	if err := h.shopUsecase.UpdateStore(&req); err != nil {
		switch err {
		case helper.ErrNotAdmin:
			helper.WriteError(w, http.StatusUnauthorized, "bukan admin")
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, nil)
}

func (h *ShopHandler) DeleteStore(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	claims, ok := claimsRaw.(*helper.JWTCLAIMS)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "protected api")
		return
	}

	params := mux.Vars(r)
	paramsStoreId, err := strconv.Atoi(params["storeId"])
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "params tidak ditemukan")
		return
	}

	if err := h.shopUsecase.DeleteStore(uint(paramsStoreId), claims.UserID); err != nil {
		switch err {
		case helper.ErrNotAdmin:
			helper.WriteError(w, http.StatusUnauthorized, "bukan admin")
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, nil)
}

// product
func (h *ShopHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	claims, ok := claimsRaw.(*helper.JWTCLAIMS)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "protected api")
		return
	}

	var req dto.CreateProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Stock < 1 {
		helper.WriteError(w, http.StatusBadRequest, "invalid stock")
		return
	}

	params := mux.Vars(r)
	paramsStoreId, _ := strconv.Atoi(params["storeId"])

	req.UserID = claims.UserID
	req.StoreID = uint(paramsStoreId)
	if err := h.shopUsecase.CreateProduct(&req); err != nil {
		switch err {
		case helper.ErrNotAdmin:
			helper.WriteError(w, http.StatusUnauthorized, "bukan admin")
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, nil)

}

func (h *ShopHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	claims, ok := claimsRaw.(*helper.JWTCLAIMS)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "protected api")
		return
	}

	var req dto.UpdateProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Stock < 1 {
		helper.WriteError(w, http.StatusBadRequest, "invalid stock")
		return
	}

	params := mux.Vars(r)
	paramsProductId, err := strconv.Atoi(params["productId"])
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "params tidak ditemukan")
		return
	}

	req.ID = uint(paramsProductId)
	req.UserID = claims.UserID
	if err := h.shopUsecase.UpdateProduct(&req); err != nil {
		switch err {
		case helper.ErrNotAdmin:
			helper.WriteError(w, http.StatusUnauthorized, "bukan admin")
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, nil)
}

func (h *ShopHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	claims, ok := claimsRaw.(*helper.JWTCLAIMS)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "protected api")
		return
	}

	params := mux.Vars(r)
	paramsStoreId, err := strconv.Atoi(params["storeId"])
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "params tidak ditemukan")
		return
	}
	paramsProductId, err := strconv.Atoi(params["productId"])
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "params tidak ditemukan")
		return
	}

	if err := h.shopUsecase.DeleteProduct(claims.UserID, uint(paramsStoreId), uint(paramsProductId)); err != nil {
		switch err {
		case helper.ErrNotAdmin:
			helper.WriteError(w, http.StatusUnauthorized, "bukan admin")
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, nil)
}

func (h *ShopHandler) GetAllProduct(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	_, ok := claimsRaw.(*helper.JWTCLAIMS)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "protected api")
		return
	}
	response, err := h.shopUsecase.GetAllProduct()
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, response)
}

func (h *ShopHandler) GetThisProduct(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	_, ok := claimsRaw.(*helper.JWTCLAIMS)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "protected api")
		return
	}

	params := mux.Vars(r)
	paramsProductId, err := strconv.Atoi(params["productId"])
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	response, err := h.shopUsecase.GetProduct(uint(paramsProductId))
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, response)
}

// cart item
func (h *ShopHandler) GetMyCartItems(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	claims, ok := claimsRaw.(*helper.JWTCLAIMS)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "protected api")
		return
	}

	response, err := h.shopUsecase.GetMyCartItems(claims.UserID)
	if err != nil {
		switch err {
		case helper.ErrUnavaible:
			helper.WriteError(w, http.StatusOK, "kau belum ada cart items")
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, response)
}

func (h *ShopHandler) CreateCartItem(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	claims, ok := claimsRaw.(*helper.JWTCLAIMS)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "protected api")
		return
	}

	params := mux.Vars(r)
	paramsProductId, _ := strconv.Atoi(params["productId"])

	var req dto.CreateCartItemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if req.PurchaseAmount < 1 {
		helper.WriteError(w, http.StatusBadRequest, "invalid stock")
		return
	}

	req.UserID = claims.UserID
	req.ProductID = uint(paramsProductId)
	if err := h.shopUsecase.CreateCartItem(&req); err != nil {
		switch err {
		case helper.ErrStocknotEnough:
			helper.WriteError(w, http.StatusBadRequest, "stock tak cukup")
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, nil)
}

func (h *ShopHandler) UpdateAmountCartItem(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	claims, ok := claimsRaw.(*helper.JWTCLAIMS)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "protected api")
		return
	}

	params := mux.Vars(r)
	paramsCartItemId, _ := strconv.Atoi(params["cartItemId"])
	paramsProductId, _ := strconv.Atoi(params["productId"])

	var req dto.UpdateAmountCartItemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid")
		return
	}

	if req.PurchaseAmount < 1 {
		helper.WriteError(w, http.StatusBadRequest, "invalid stock")
		return
	}

	req.UserID = claims.UserID
	req.ID = uint(paramsCartItemId)
	req.ProductID = uint(paramsProductId)
	if err := h.shopUsecase.UpdateAmountCartItem(&req); err != nil {
		switch err {
		case helper.ErrStocknotEnough:
			helper.WriteError(w, http.StatusBadRequest, "stock tak cukup")
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

	}

	helper.WriteJSON(w, http.StatusOK, nil)
}

func (h *ShopHandler) UpdatePaidCartItem(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	claims, ok := claimsRaw.(*helper.JWTCLAIMS)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "protected api")
		return
	}

	params := mux.Vars(r)
	paramsId, _ := strconv.Atoi(params["cartItemId"])
	paramsProductId, _ := strconv.Atoi(params["productId"])

	var req dto.UpdatePaidCartItemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid")
		return
	}
	if req.PurchaseAmount < 1 {
		helper.WriteError(w, http.StatusBadRequest, "invalid stock")
		return
	}

	req.UserID = claims.UserID
	req.Email = claims.Email
	req.ID = uint(paramsId)
	req.ProductID = uint(paramsProductId)
	if err := h.shopUsecase.UpdatePaidCartItem(&req); err != nil {
		switch err {
		case helper.ErrStocknotEnough:
			helper.WriteError(w, http.StatusBadRequest, "stock tak cukup")
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

	}

	helper.WriteJSON(w, http.StatusOK, nil)

}

func (h *ShopHandler) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	claims, ok := claimsRaw.(*helper.JWTCLAIMS)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "protected api")
		return
	}

	params := mux.Vars(r)
	paramsId, _ := strconv.Atoi(params["cartItemId"])

	if err := h.shopUsecase.DeleteCartItem(claims.UserID, uint(paramsId)); err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, nil)
}
