package apis

import (
	. "book_management_system_backend/models"
	"github.com/gofiber/fiber/v2"
)

// GetMeta godoc
// @Summary 获取统计信息
// @Tags Meta Module
// @Produce json
// @Router /meta [get]
// @Success 200 {object} MetaInfo
func GetMeta(c *fiber.Ctx) (err error) {
	var currentUser User
	if err = GetCurrentUser(c, &currentUser); err != nil {
		return
	}
	var metaInfo MetaInfo

	// 统计用户、书籍、购买记录、销售记录、流水的数量
	if err = DB.Model(&User{}).Count(&metaInfo.UserCount).Error; err != nil {
		return
	}
	if err = DB.Model(&Book{}).Count(&metaInfo.BookCount).Error; err != nil {
		return
	}
	if err = DB.Model(&Purchase{}).Count(&metaInfo.PurchaseCount).Error; err != nil {
		return
	}
	if err = DB.Model(&Sale{}).Count(&metaInfo.SaleCount).Error; err != nil {
		return
	}
	if err = DB.Model(&Balance{}).Count(&metaInfo.BalanceCount).Error; err != nil {
		return
	}

	// 统计过去12个月的月份和每个月的销售数量，Postgres 专用
	if err = DB.Raw(`
		SELECT to_char(sale_time, 'YYYY-MM') AS month, COUNT(*) AS count
		FROM sale
		GROUP BY month
		ORDER BY month DESC
		LIMIT 12
	`).Scan(&metaInfo.SaleCountByMonth).Error; err != nil {
		return
	}

	// 统计过去12个月的月份和每个月的购买数量，Postgres 专用
	if err = DB.Raw(`
		SELECT to_char(purchase_time, 'YYYY-MM') AS month, COUNT(*) AS count
		FROM purchase
		GROUP BY month
		ORDER BY month DESC
		LIMIT 12
	`).Scan(&metaInfo.PurchaseCountByMonth).Error; err != nil {
		return
	}

	// 统计过去12个月的月份和每个月的流水数量，Postgres 专用
	if err = DB.Raw(`
		SELECT to_char(time, 'YYYY-MM') AS month, COUNT(*) AS count
		FROM balance
		GROUP BY month
		ORDER BY month DESC
		LIMIT 12
	`).Scan(&metaInfo.BalanceCountByMonth).Error; err != nil {
		return
	}
	return c.JSON(metaInfo)
}
