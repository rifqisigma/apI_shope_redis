package usecase

import (
	"api_shope/dto"
	"api_shope/internal/repository"
	"api_shope/utils/helper"
)

type ShopUsecase interface {

	//store
	GetMyStore(userId, storeId uint) (*dto.StoreAndProduct, error)
	GetAllStore() ([]dto.JustStore, error)
	CreateStore(req *dto.CreateStoreReq) error
	UpdateStore(req *dto.UpdateStoreReq) error
	DeleteStore(storeId, userId uint) error

	//product
	GetAllProduct() ([]dto.Product, error)
	GetProduct(id uint) (*dto.Product, error)
	CreateProduct(req *dto.CreateProductReq) error
	UpdateProduct(req *dto.UpdateProductReq) error
	DeleteProduct(userId, storeId, id uint) error

	//cartItem
	GetMyCartItems(userId uint) ([]dto.CartItem, error)
	CreateCartItem(req *dto.CreateCartItemReq) error
	UpdateAmountCartItem(req *dto.UpdateAmountCartItemReq) error
	UpdatePaidCartItem(req *dto.UpdatePaidCartItemReq) error
	DeleteCartItem(userId, id uint) error
}

type shopUsecase struct {
	shopRepo repository.ShopRepo
}

func NewShopUsecase(shopRepo repository.ShopRepo) ShopUsecase {
	return &shopUsecase{shopRepo}
}

func (u *shopUsecase) GetMyStore(userId, storeId uint) (*dto.StoreAndProduct, error) {
	valid, err := u.shopRepo.IsUserAdminStore(userId, storeId)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, helper.ErrNotAdmin
	}

	return u.shopRepo.GetMyStore(userId)
}

func (u *shopUsecase) GetAllStore() ([]dto.JustStore, error) {
	return u.shopRepo.GetAllStore()
}

func (u *shopUsecase) CreateStore(req *dto.CreateStoreReq) error {
	return u.shopRepo.CreateStore(req)
}

func (u *shopUsecase) UpdateStore(req *dto.UpdateStoreReq) error {
	valid, err := u.shopRepo.IsUserAdminStore(req.UserID, req.ID)
	if err != nil {
		return err
	}
	if !valid {
		return helper.ErrNotAdmin
	}

	return u.shopRepo.UpdateStore(req)
}

func (u *shopUsecase) DeleteStore(storeId, userId uint) error {
	valid, err := u.shopRepo.IsUserAdminStore(userId, storeId)
	if err != nil {
		return err
	}
	if !valid {
		return helper.ErrNotAdmin
	}

	return u.shopRepo.DeleteStore(storeId)

}

// product
func (u *shopUsecase) CreateProduct(req *dto.CreateProductReq) error {
	valid, err := u.shopRepo.IsUserAdminStore(req.UserID, req.StoreID)
	if err != nil {
		return err
	}
	if !valid {
		return helper.ErrNotAdmin
	}

	return u.shopRepo.CreateProduct(req)
}

func (u *shopUsecase) UpdateProduct(req *dto.UpdateProductReq) error {
	return u.shopRepo.UpdateProduct(req)
}

func (u *shopUsecase) DeleteProduct(userId, storeId, id uint) error {
	valid, err := u.shopRepo.IsUserAdminStore(userId, storeId)
	if err != nil {
		return err
	}
	if !valid {
		return helper.ErrNotAdmin
	}

	return u.shopRepo.DeleteProduct(id)
}

func (u *shopUsecase) GetAllProduct() ([]dto.Product, error) {
	return u.shopRepo.GetAllProduct()
}

func (u *shopUsecase) GetProduct(id uint) (*dto.Product, error) {
	return u.shopRepo.GetProduct(id)
}

// cart item
func (u *shopUsecase) GetMyCartItems(userId uint) ([]dto.CartItem, error) {
	return u.shopRepo.GetMyCartItems(userId)
}

func (u *shopUsecase) CreateCartItem(req *dto.CreateCartItemReq) error {
	valid, err := u.shopRepo.CheckStock(req.ProductID, req.PurchaseAmount)
	if err != nil {
		return err
	}
	if !valid {
		return helper.ErrStocknotEnough
	}
	return u.shopRepo.CreateCartItem(req)
}

func (u *shopUsecase) UpdateAmountCartItem(req *dto.UpdateAmountCartItemReq) error {
	valid, err := u.shopRepo.CheckStock(req.ProductID, req.PurchaseAmount)
	if err != nil {
		return err
	}
	if !valid {
		return helper.ErrStocknotEnough
	}
	return u.shopRepo.UpdateAmountCartItem(req)
}

func (u *shopUsecase) UpdatePaidCartItem(req *dto.UpdatePaidCartItemReq) error {
	valid, err := u.shopRepo.CheckStock(req.ProductID, req.PurchaseAmount)
	if err != nil {
		return err
	}
	if !valid {
		return helper.ErrStocknotEnough
	}
	return u.shopRepo.UpdatePaidCartItem(req)
}

func (u *shopUsecase) DeleteCartItem(userId, id uint) error {
	return u.shopRepo.DeleteCartItem(userId, id)
}
