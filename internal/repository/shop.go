package repository

import (
	"api_shope/dto"
	"api_shope/model"
	"api_shope/utils/helper"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ShopRepo interface {
	IsUserAdminStore(userId, storeId uint) (bool, error)

	//store
	GetMyStore(userId uint) (*dto.StoreAndProduct, error)
	GetAllStore() ([]dto.JustStore, error)
	CreateStore(req *dto.CreateStoreReq) error
	UpdateStore(req *dto.UpdateStoreReq) error
	DeleteStore(id uint) error

	//product
	GetAllProduct() ([]dto.Product, error)
	GetProduct(id uint) (*dto.Product, error)
	CreateProduct(req *dto.CreateProductReq) error
	UpdateProduct(req *dto.UpdateProductReq) error
	DeleteProduct(id uint) error

	// cart item
	GetMyCartItems(userId uint) ([]dto.CartItem, error)
	CreateCartItem(req *dto.CreateCartItemReq) error
	UpdateAmountCartItem(req *dto.UpdateAmountCartItemReq) error
	UpdatePaidCartItem(req *dto.UpdatePaidCartItemReq) error
	DeleteCartItem(userId, id uint) error
	CheckStock(id uint, req int) (bool, error)
}

type shopRepo struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewShopRepo(db *gorm.DB, redis *redis.Client) ShopRepo {
	return &shopRepo{db, redis}
}

var ctx = context.Background()

func (r *shopRepo) IsUserAdminStore(userId, storeId uint) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Store{}).Where("id = ? AND admin_id = ?", storeId, userId).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// penerapan metode caching dengan lazy loading
func (r *shopRepo) GetMyStore(userId uint) (*dto.StoreAndProduct, error) {
	key := fmt.Sprintf("mystore:user:%d", userId)

	cachedData, err := r.redis.Get(ctx, key).Result()
	if err == nil && cachedData != "" {
		var cachedStore dto.StoreAndProduct
		if err := json.Unmarshal([]byte(cachedData), &cachedStore); err == nil {
			log.Println("data dari redis")
			return &cachedStore, nil
		}
	}

	log.Println("data dari mysql")

	var store model.Store
	if err := r.db.Preload("Product").Where("admin_id = ?", userId).First(&store).Error; err != nil {
		return nil, err
	}

	var getProduct []dto.Product
	for _, p := range store.Product {
		getProduct = append(getProduct, dto.Product{
			StoreID:   p.StoreID,
			ID:        p.ID,
			Name:      p.Name,
			Stock:     p.Stock,
			CreatedAt: p.CreatedAt,
		})
	}

	response := dto.StoreAndProduct{
		ID:        store.ID,
		AdminID:   store.AdminID,
		Name:      store.Name,
		CreatedAt: store.CreatedAt,
		Product:   getProduct,
	}

	jsonData, _ := json.Marshal(response)
	err = r.redis.Set(ctx, key, jsonData, 20*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (r *shopRepo) GetAllStore() ([]dto.JustStore, error) {
	key := fmt.Sprintln("store:all")

	cachedData, err := r.redis.Get(ctx, key).Result()
	if err == nil && cachedData != "" {
		var cachedShops []dto.JustStore
		if err := json.Unmarshal([]byte(cachedData), &cachedShops); err == nil {
			log.Println("data dari redis")
			return cachedShops, nil
		}
	}

	log.Println("data dari mysql")
	var shops []dto.JustStore
	if err := r.db.Model(&model.Store{}).Select("id", "name", "admin_id", "created_at").Find(&shops).Error; err != nil {

		return nil, err
	}

	jsonData, _ := json.Marshal(shops)
	if err := r.redis.Set(ctx, key, jsonData, 15*time.Minute).Err(); err != nil {
		return nil, err
	}

	return shops, nil
}

func (r *shopRepo) CreateStore(req *dto.CreateStoreReq) error {
	newStore := model.Store{
		Name:    req.Name,
		AdminID: req.AdminID,
	}

	if err := r.db.Model(&model.Store{}).Create(&newStore).Error; err != nil {
		return err
	}

	return nil
}

func (r *shopRepo) UpdateStore(req *dto.UpdateStoreReq) error {
	if err := r.db.Model(&model.Store{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
		"name": req.Name,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (r *shopRepo) DeleteStore(id uint) error {
	if err := r.db.Model(&model.Store{}).Where("id = ?", id).Delete(&model.Store{}).Error; err != nil {
		return err
	}

	return nil
}

// penerapan write-around caching (penggunaan lazy loading dan write trough yg bersamaan)
func (r *shopRepo) CreateProduct(req *dto.CreateProductReq) error {
	newProduct := model.Product{
		Name:    req.Name,
		StoreID: req.StoreID,
		Stock:   req.Stock,
	}

	if err := r.db.Create(&newProduct).Error; err != nil {
		return err
	}

	key := fmt.Sprintf("product:%d", newProduct.ID)

	_, err := r.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key, map[string]interface{}{
			"name":       newProduct.Name,
			"store_id":   newProduct.StoreID,
			"stock":      newProduct.Stock,
			"created_at": newProduct.CreatedAt,
		})
		pipe.Expire(ctx, key, 30*time.Minute)
		return nil
	})
	if err != nil {
		return err
	}

	r.redis.Del(ctx, "products:all")
	return nil
}

func (r *shopRepo) UpdateProduct(req *dto.UpdateProductReq) error {
	if err := r.db.Model(&model.Product{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
		"name":  req.Name,
		"stock": req.Stock,
	}).Error; err != nil {
		return err
	}
	key := fmt.Sprintf("product:%d", req.ID)

	_, err := r.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key,
			"name", req.Name,
			"stock", req.Stock,
		)
		pipe.Expire(ctx, key, 30*time.Minute)
		return nil
	})
	if err != nil {
		return err
	}

	r.redis.Del(ctx, "products:all")
	return nil
}

func (r *shopRepo) DeleteProduct(id uint) error {
	tx := r.db.Begin()
	if err := tx.Model(&model.CartItem{}).Where("product_id = ?", id).Update("is_product_deleted", true).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&model.Product{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	key := fmt.Sprintf("product:%d", id)

	r.redis.HDel(ctx, key)
	r.redis.Del(ctx, "products:all")

	return nil
}

func (r *shopRepo) GetProduct(id uint) (*dto.Product, error) {
	key := fmt.Sprintf("product:%d", id)

	var product model.Product
	data, err := r.redis.HGetAll(ctx, key).Result()
	if err == nil {
		stock, _ := strconv.Atoi(data["stock"])
		storeID, _ := strconv.Atoi(data["store_id"])
		createdAt, _ := time.Parse(time.RFC3339, data["created_at"])

		product := dto.Product{
			ID:        id,
			Name:      data["name"],
			Stock:     stock,
			StoreID:   uint(storeID),
			CreatedAt: createdAt,
		}

		fmt.Println("data dari redis")
		return &product, nil
	}

	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	fmt.Println("data dari mysql")
	return &dto.Product{
		ID:        product.ID,
		StoreID:   product.StoreID,
		Name:      product.Name,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
	}, nil
}

func (r *shopRepo) GetAllProduct() ([]dto.Product, error) {
	cached, err := r.redis.Get(ctx, "products:all").Result()
	if err == nil {
		var products []dto.Product
		if json.Unmarshal([]byte(cached), &products) == nil {
			fmt.Println("data dari redis")
			return products, nil
		}
	}

	var products []model.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}

	var result []dto.Product
	for _, p := range products {
		result = append(result, dto.Product{
			ID:        p.ID,
			StoreID:   p.StoreID,
			Name:      p.Name,
			Stock:     p.Stock,
			CreatedAt: p.CreatedAt,
		})
	}

	jsonData, _ := json.Marshal(result)
	r.redis.Set(ctx, "products:all", jsonData, 30*time.Minute)

	fmt.Println("data dari mysql")
	return result, nil
}

// penerapan write trough
func (r *shopRepo) CreateCartItem(req *dto.CreateCartItemReq) error {
	newCartItem := model.CartItem{
		ProductID:      &req.ProductID,
		UserID:         req.UserID,
		PurchaseAmount: req.PurchaseAmount}

	if err := r.db.Model(&model.CartItem{}).Create(&newCartItem).Error; err != nil {
		return err
	}

	key := fmt.Sprintf("user:%d:cartitem:%d:", req.UserID, newCartItem.ID)
	userCartItemsKey := fmt.Sprintf("user:%d:cartitems", req.UserID)
	_, err := r.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key, map[string]interface{}{
			"product_id":      newCartItem.ProductID,
			"user_id":         newCartItem.UserID,
			"purchase_amount": newCartItem.PurchaseAmount,
			"is_paid":         false,
		})
		pipe.Expire(ctx, key, 30*time.Minute)
		pipe.SAdd(ctx, userCartItemsKey, newCartItem.ID)
		pipe.Expire(ctx, userCartItemsKey, 30*time.Minute)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *shopRepo) UpdateAmountCartItem(req *dto.UpdateAmountCartItemReq) error {
	if err := r.db.Model(&model.CartItem{}).Where("id = ?", req.ID).Update("purchase_amount", req.PurchaseAmount).Error; err != nil {
		return err
	}

	key := fmt.Sprintf("user:%d:cartitem:%d:", req.UserID, req.ID)
	_, err := r.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key,
			"purchase_amount", req.PurchaseAmount,
		)
		pipe.Expire(ctx, key, 30*time.Minute)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *shopRepo) UpdatePaidCartItem(req *dto.UpdatePaidCartItemReq) error {
	if err := r.db.Model(&model.CartItem{}).Where("id = ?", req.ID).Update("is_paid", true).Error; err != nil {
		return err
	}

	key := fmt.Sprintf("user:%d:cartitem:%d:", req.UserID, req.ID)
	keyQueque := fmt.Sprintf("behind:pending:buy:%d", req.ID)
	message := fmt.Sprintf("pembelian product dengan id %v /n total item %v", req.ProductID, req.PurchaseAmount)
	_, err := r.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key,
			"purchase_amount", req.PurchaseAmount,
			"is_paid", true,
		)
		pipe.Expire(ctx, key, 30*time.Minute)
		pipe.HSet(ctx, keyQueque, map[string]interface{}{
			"id":      req.ID,
			"email":   req.Email,
			"message": message,
			"op":      "buy",
		})
		pipe.Expire(ctx, keyQueque, 10*time.Minute)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *shopRepo) DeleteCartItem(userId, id uint) error {
	if err := r.db.Model(&model.CartItem{}).Where("id = ?", id).Delete(&model.CartItem{}).Error; err != nil {
		return err
	}

	itemKey := fmt.Sprintf("user:%d:cartitem:%d", userId, id)
	itemsKey := fmt.Sprintf("user:%d:cartitems", userId)

	_, err := r.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SRem(ctx, itemsKey, id)
		pipe.Del(ctx, itemKey)
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *shopRepo) GetMyCartItems(userId uint) ([]dto.CartItem, error) {
	itemsKey := fmt.Sprintf("user:%d:cartitems", userId)
	itemIDs, err := r.redis.SMembers(ctx, itemsKey).Result()

	var items []dto.CartItem

	if err == nil || len(itemIDs) != 0 {
		cmds := make([]*redis.MapStringStringCmd, len(itemIDs))
		_, err = r.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			for i, itemID := range itemIDs {
				key := fmt.Sprintf("user:%d:cartitem:%s:", userId, itemID)
				cmds[i] = pipe.HGetAll(ctx, key)
			}
			return nil
		})
		if err == nil {
			for i, cmd := range cmds {
				data, err := cmd.Result()
				if err == nil || len(data) != 0 {
					cartItemId, _ := strconv.ParseUint(itemIDs[i], 10, 64)
					productID, _ := strconv.ParseUint(data["product_id"], 10, 64)
					userIDParsed, _ := strconv.ParseUint(data["user_id"], 10, 64)
					purchaseAmount, _ := strconv.Atoi(data["purchase_amount"])
					isPaid, _ := strconv.ParseBool(data["is_paid"])

					items = append(items, dto.CartItem{
						ID:             uint(cartItemId),
						ProductID:      uint(productID),
						UserID:         uint(userIDParsed),
						PurchaseAmount: purchaseAmount,
						IsPaid:         isPaid,
					})
				}

			}

			if len(items) > 0 {
				fmt.Println("data dari redis")
				return items, nil
			}

		}
	}

	if err := r.db.Where("user_id = ?", userId).Find(&items).Error; err != nil {

		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("user %d does not have any cart items", userId)
	}

	fmt.Println("data dari mysql (fallback)")
	return items, nil
}

func (r *shopRepo) CheckStock(id uint, req int) (bool, error) {
	key := fmt.Sprintf("product:%d", id)

	stockStr, err := r.redis.HGet(ctx, key, "stock").Result()
	if err == nil {
		stock, _ := strconv.Atoi(stockStr)
		if req > stock {
			fmt.Println("check dari redis")
			return false, helper.ErrStocknotEnough
		}
	}

	var stockProduct int
	if err := r.db.Model(&model.Product{}).Select("stock").Where("id = ?", id).Scan(&stockProduct).Error; err != nil {
		return false, err
	}

	if req > stockProduct {
		return false, helper.ErrStocknotEnough
	}

	fmt.Println("check dari mysql")
	return true, nil
}
