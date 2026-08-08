package receiptutils

import (
	"fmt"
	"kopibang/domain/entity"
)

// GenerateReceiptText membuat teks struk yang rapi untuk dikirim via WA
func GenerateReceiptText(order *entity.Order, products map[string]string) string {
	text := fmt.Sprintf("*kopibang COFFEE - RECEIPT*\nOrder ID: %s\nDate: %s\n\n", order.ID.String()[:8], order.CreatedAt.Format("02 Jan 2006 15:04"))
	
	for _, item := range order.Items {
		productName := products[item.ProductID.String()]
		text += fmt.Sprintf("- %dx %s (@%d) : %d\n", item.Quantity, productName, item.PriceAtTime, item.Quantity*item.PriceAtTime)
	}

	text += fmt.Sprintf("\nSubtotal: %d\nDiscount: %d\n*Total: %d*\n", order.TotalAmount, order.Discount, order.FinalAmount)
	
	if order.IsRedeem {
		text += "\n*(Redeemed 100 Points)*"
	}

	text += "\n\nTerima kasih telah berbelanja!"
	return text
}