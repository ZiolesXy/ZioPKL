package models

type CheckoutItem struct {
	ID uint
	CheckoutID uint
	ProductID uint
	Quantity int
	Price float64
	Product Product
}