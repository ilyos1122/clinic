package handler

import (
	"context"
	"fmt"
	"market_system/config"
	"market_system/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// @Summary MakePay
// @Description Get List Coming details by its ok.
// @Tags Pay
// @Accept json
// @Produce json
// @Param sale_increment_id query string ture "sale_increment_id"
// @Param money query float64 true "Pay money"
// @Success 200 {object} models.Coming "Coming details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Coming not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /make_pay [post]
func (h *Handler) MakePay(c *gin.Context) {

	var (
		incrementId = c.Query("sale_increment_id")
		money       = (cast.ToFloat64(c.Query("money")))
		Id          string
		ClientID    string
		BranchID    string
		IncrementID string
		TotalPrice  float64
		// Paid        float64
		// Debd        float64
	)

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	saleList, err := h.strg.Sale().GetList(ctx, &models.GetListSaleRequest{Limit: 10000})
	for _, v := range saleList.Sales {
		if v.IncrementID == incrementId {
			if v.TotalPrice/2 < money {

				Id = v.Id
				ClientID = v.ClientID
				BranchID = v.BranchID
				IncrementID = v.IncrementID
				TotalPrice = v.TotalPrice
				// Paid = v.Paid
				// Debd = v.Debd

				_, err = h.strg.Sale().Update(ctx, &models.UpdateSale{
					Id:          Id,
					ClientID:    ClientID,
					BranchID:    BranchID,
					IncrementID: IncrementID,
					TotalPrice:  TotalPrice,
					Status:      "success",
					Paid:        money,
					Debd:        TotalPrice - money,
				})
				if err != nil {
					handleResponse(c, 500, err)
					return
				}
				handleResponse(c, 202, "successful payment")

			} else {
				handleResponse(c, http.StatusBadRequest, "not enough money")
			}
		}
	}
	// fmt.Println(
	// 	Id,
	// 	ClientID,
	// 	BranchID,
	// 	IncrementID,
	// 	TotalPrice)

}

// @Summary Get all sales by branch_id
// @Description  Get all sale by branch
// @Tags Otchet
// @Accept json
// @Produce json
// @Param branch_id query string true "Branch ID"
// @Success 200 {object} models.OtchetSale "Sale details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /otchet/{id} [get]
func (h *Handler) OtchetTwo(c *gin.Context) {
	branchID := c.Query("branch_id")
	fmt.Println(branchID)
	var resp models.OtchetSale

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()
	count, err := h.strg.Sale().GetList(ctx, &models.GetListSaleRequest{})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}
	sales, err := h.strg.Sale().GetList(ctx, &models.GetListSaleRequest{
		Offset: 0,
		Limit:  int64(count.Count),
	})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	var sum float64
	for _, sale := range sales.Sales {
		if sale.BranchID == branchID || sale.Status == "success" {
			sum = sum + sale.TotalPrice
		}
	}

	resp.BranchID = branchID
	resp.SaleCount = sales.Count
	resp.Sum = sum

	handleResponse(c, http.StatusOK, resp)

}

// @Summary Registration
// @Description Get clients.
// @Tags Registration
// @Accept json
// @Produce json
// @Param from query string true "From day"
// @Param to query string true "To day"
// @Success 200 {object} models.GetListClientResponse "Client details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Client not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /otchet [get]
func (h *Handler) Registration(c *gin.Context) {

	var (
		from    = c.Query("from")
		to      = c.Query("to")
		clients = models.GetListClientResponse{}
	)
	fromm, err := time.Parse("2006-01-02", from)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	too, err := time.Parse("2006-01-02", to)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(fromm,too)

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	clientList, err := h.strg.Client().GetList(ctx, &models.GetListClientRequest{Limit: 1000000})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	for _, v := range clientList.Clients {
		created_at, err := time.Parse("2006-01-02", v.CreatedAt[:10])
		if err != nil {
			handleResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		if too.After(created_at) && fromm.Before(created_at) {
			if len(created_at.String()) > 0 {
				clients.Clients = append(clients.Clients, v)
			}
		}
	}

	handleResponse(c, http.StatusOK, clients)
}
