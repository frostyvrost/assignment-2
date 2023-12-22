package app

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	received    = "RECEIVED"
	in_progress = "IN PROGRESS"
	shipped     = "SHIPPED"
	delivered   = "DELIVERED"
	cancelled   = "CANCELLED"
)

type NewOrderReq struct {
	CustomerName string    `json:"customer_name"`
	Products     []Product `json:"products"`
}

type OrderResponse struct {
	OrderId      int64  `json:"order_id"`
	ProductId    int64  `json:"product_id"`
	ProductCode   string `json:"item_code"`
	CustomerName string `json:"customer_name"`
}

var dbmap = connectToDB()

func Ping(c *gin.Context) {}

func CreateOrder(c *gin.Context) {
	var orderReq NewOrderReq
	c.Bind(&orderReq)

	if isEmpty(orderReq.CustomerName) {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "error": "Customer name cannot be empty. Pass a valid string value and try again !"})
		return
	}
	if len(orderReq.Products) == 0 {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "error": "Products cannot be empty. An order need to have atleast 1 product. Add a product and try again !"})
		return
	}

	order := &Order{
		CustomerName: orderReq.CustomerName,
		Status:       received,
		CreatedAt:    time.Now().UnixNano(),
		UpdatedAt:    time.Now().UnixNano(),
	}

	err := dbmap.Insert(order)
	checkErr(err, "Add new order failed in orders table")

	for _, product := range orderReq.Products {

		if isEAN(product.ProductCode) {
			orderProduct := &OrderProduct{
				OrderId:      order.Id,
				ProductId:    product.Id,
				ProductCode:  product.ProductCode,
				CustomerName: orderReq.CustomerName,
				CreatedAt:    time.Now().UnixNano(),
				UpdatedAt:    time.Now().UnixNano(),
			}
			err := dbmap.Insert(orderProduct)
			checkErr(err, "Add new order_product mapping failed in order_products table")
		} else {
			errMsg := "Product EAN is invalid : " + product.ProductCode
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": errMsg})
			return
		}
	}

	c.JSON(http.StatusCreated,
		gin.H{"status": http.StatusCreated, "message": "Order Created Successfully!", "resourceId": order.Id})
}

func FetchOrder(c *gin.Context) {
	orderId := c.Params.ByName("id")

	if isNumber(orderId) {
	} else {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "error": "Order ID passed is not a valid number."})
		return
	}

	log.Println("Fetching order details for order ID : " + orderId)
	var orderProducts []OrderProduct
	var query = "SELECT * FROM order_products where order_id=" + orderId + " ORDER BY order_product_id"

	_, err := dbmap.Select(&orderProducts, query)

	if len(orderProducts) == 0 || err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"status": http.StatusNotFound, "error": "No order with requested ID exists in the table. Invalid ID."})
	} else {
		var orderResArr []OrderResponse
		for _, op := range orderProducts {
			orderRes := OrderResponse{
				OrderId:      op.OrderId,
				ProductId:    op.ProductId,
				ProductCode:  op.ProductCode,
				CustomerName: op.CustomerName,
			}
			orderResArr = append(orderResArr, orderRes)
		}

		c.JSON(http.StatusOK,
			gin.H{"status": http.StatusOK, "message": "Order Details Fetched Successfully!", "order": orderResArr})
	}
}

func UpdateOrder(c *gin.Context) {
	orderId := c.Params.ByName("id")

	if isNumber(orderId) {
	} else {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "error": "Order ID passed is not a valid number."})
		return
	}

	log.Println("Updating order details for order ID :", orderId)
	var query_orders = "SELECT * FROM orders where order_id=" + orderId
	var order Order
	err := dbmap.SelectOne(&order, query_orders)

	if err != nil || (Order{}) == order || len(order.CustomerName) == 0 {
		c.JSON(http.StatusNotFound,
			gin.H{"status": http.StatusNotFound, "error": "No order with requested ID exists in the table. Invalid ID."})
	} else {
		var orderReq NewOrderReq
		c.Bind(&orderReq)

		if orderReq.CustomerName == order.CustomerName {
			order.UpdatedAt = time.Now().UnixNano()
		} else {
			order.CustomerName = orderReq.CustomerName
			order.UpdatedAt = time.Now().UnixNano()
		}

		_, err = dbmap.Update(&order)

		var orderProducts []OrderProduct
		var query_order_products = "SELECT * FROM order_products where order_id=" + orderId + " ORDER BY order_product_id"
		_, err := dbmap.Select(&orderProducts, query_order_products)

		if len(orderProducts) == 0 || err != nil {
			c.JSON(http.StatusNotFound,
				gin.H{"status": http.StatusNotFound, "error": "No products found with requested order ID. Aborting!"})
		} else if len(orderProducts) != len(orderReq.Products) {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "v1 Update API supports only product details updation in the existing order. New products addition and existing products deletion from existing order will be supported in future API ver. Aborting!"})
		} else {
			for index, orderProduct := range orderProducts {

				product := orderReq.Products[index]

				orderProduct.CustomerName = orderReq.CustomerName
				orderProduct.ProductId = product.Id
				orderProduct.ProductCode = product.ProductCode
				orderProduct.UpdatedAt = time.Now().UnixNano()

				_, err := dbmap.Update(&orderProduct)
				checkErr(err, "Updating order_product mapping failed in order_products table")
			}

			c.JSON(http.StatusOK,
				gin.H{"status": http.StatusOK, "message": "Order Updated Successfully!", "resourceId": order.Id})
		}
	}
}
