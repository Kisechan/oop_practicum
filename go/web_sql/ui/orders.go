package ui

import (
	"fmt"
	"strconv"
	"time"
	"web_sql/rep"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ordersScreen(win fyne.Window) fyne.CanvasObject {
	var records []rep.Order
	var userID int

	// 初始化表格
	table := widget.NewTable(
		func() (int, int) { return len(records) + 1, 8 }, // 8 列：ID, UserID, Total, Status, CreatedTime, UpdateTime, ProductID, Quantity
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell 000, 000")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			if i.Row == 0 {
				// 表头
				switch i.Col {
				case 0:
					label.SetText("ID")
				case 1:
					label.SetText("用户ID")
				case 2:
					label.SetText("总金额")
				case 3:
					label.SetText("状态")
				case 4:
					label.SetText("创建时间")
				case 5:
					label.SetText("更新时间")
				case 6:
					label.SetText("产品ID")
				case 7:
					label.SetText("数量")
				}
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				// 数据行
				record := records[i.Row-1]
				switch i.Col {
				case 0:
					label.SetText(fmt.Sprintf("%d", record.ID))
				case 1:
					label.SetText(fmt.Sprintf("%d", record.UserID))
				case 2:
					if record.Total != nil {
						label.SetText(fmt.Sprintf("%.2f", *record.Total))
					} else {
						label.SetText("N/A")
					}
				case 3:
					label.SetText(record.Status)
				case 4:
					label.SetText(record.CreatedTime.Format("2006-01-02 15:04:05"))
				case 5:
					label.SetText(record.UpdateTime.Format("2006-01-02 15:04:05"))
				case 6:
					label.SetText(fmt.Sprintf("%d", record.ProductID))
				case 7:
					label.SetText(fmt.Sprintf("%d", record.Quantity))
				}
			}
		},
	)
	table.SetColumnWidth(1, 140)
	table.SetColumnWidth(4, 240)
	table.SetColumnWidth(5, 240)
	table.Hide()
	var selectedRow int
	table.OnSelected = func(id widget.TableCellID) {
		if id.Row > 0 {
			selectedRow = id.Row - 1
			fmt.Printf("Selected row: %d\n", selectedRow)
		}
	}

	// 查找用户ID的输入框
	entry := widget.NewEntry()
	entry.SetPlaceHolder("查找用户ID")

	// 查找按钮
	search := widget.NewButtonWithIcon("查找用户", theme.SearchIcon(), func() {
		userIDStr := entry.Text
		userIDInt, err := strconv.Atoi(userIDStr)
		if err != nil {
			dialog.ShowError(fmt.Errorf("用户ID必须是整数"), win)
			return
		}
		userID = userIDInt
		records, _ = rep.GetOrdersByUserID(rep.DB, userID)
		table.Show()
		table.Refresh()
	})

	// 增加按钮
	addButton := widget.NewButtonWithIcon("增加订单", theme.ContentAddIcon(), func() {
		if userID == 0 {
			dialog.ShowInformation("错误", "请先查找用户ID", win)
			return
		}

		productID := widget.NewEntry()
		productID.SetPlaceHolder("产品ID")
		productID.Validator = validation.NewRegexp(`^\d+$`, "产品ID必须是整数")

		quantity := widget.NewEntry()
		quantity.SetPlaceHolder("数量")
		quantity.Validator = validation.NewRegexp(`^\d+$`, "数量必须是整数")

		total := widget.NewEntry()
		total.SetPlaceHolder("总金额")
		total.Validator = validation.NewRegexp(`^\d+(\.\d{1,2})?$`, "总金额必须是数字，最多两位小数")

		status := widget.NewSelect(
			[]string{"pending", "paid", "shipping", "completed"},
			func(s string) {},
		)
		status.Selected = "pending"

		items := []*widget.FormItem{
			widget.NewFormItem("产品ID", productID),
			widget.NewFormItem("数量", quantity),
			widget.NewFormItem("总金额", total),
			widget.NewFormItem("状态", status),
		}

		form := dialog.NewForm("增加一条订单信息", "确认", "取消", items, func(b bool) {
			if !b {
				return
			}

			productIDInt, _ := strconv.Atoi(productID.Text)
			quantityInt, _ := strconv.Atoi(quantity.Text)
			totalValue, _ := strconv.ParseFloat(total.Text, 64)

			err := rep.Create[rep.Order](rep.DB, &rep.Order{
				UserID:      userID,
				ProductID:   productIDInt,
				Quantity:    quantityInt,
				Total:       &totalValue,
				Status:      status.Selected,
				CreatedTime: time.Now(),
				UpdateTime:  time.Now(),
			})
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// 刷新表格
			records, _ = rep.GetOrdersByUserID(rep.DB, userID)
			table.Refresh()
		}, win)

		form.Resize(fyne.NewSize(400, 300))
		form.Show()
	})

	// 删除按钮
	deleteButton := widget.NewButtonWithIcon("删除订单", theme.ContentRemoveIcon(), func() {
		if selectedRow < 0 || selectedRow >= len(records) {
			dialog.ShowInformation("错误", "请先选择要删除的订单信息", win)
			return
		}

		cnf := dialog.NewConfirm(
			"确认",
			"你真的要删除这条数据吗？\n你会失去它很久的（真的很久）",
			func(b bool) {
				if !b {
					return
				}
				selectedID := records[selectedRow].ID
				rep.DeleteID[rep.Order](rep.DB, selectedID)

				// 刷新表格
				records, _ = rep.GetOrdersByUserID(rep.DB, userID)
				table.Refresh()
			},
			win,
		)
		cnf.SetDismissText("取消")
		cnf.SetConfirmText("确认")
		cnf.Show()
	})

	// 修改按钮
	updateButton := widget.NewButtonWithIcon("修改订单", theme.DocumentCreateIcon(), func() {
		if selectedRow < 0 || selectedRow >= len(records) {
			dialog.ShowInformation("错误", "请先选择要修改的订单信息", win)
			return
		}

		productID := widget.NewEntry()
		productID.SetText(fmt.Sprintf("%d", records[selectedRow].ProductID))
		productID.Validator = validation.NewRegexp(`^\d+$`, "产品ID必须是整数")

		quantity := widget.NewEntry()
		quantity.SetText(fmt.Sprintf("%d", records[selectedRow].Quantity))
		quantity.Validator = validation.NewRegexp(`^\d+$`, "数量必须是整数")

		total := widget.NewEntry()
		if records[selectedRow].Total != nil {
			total.SetText(fmt.Sprintf("%.2f", *records[selectedRow].Total))
		} else {
			total.SetText("")
		}
		total.Validator = validation.NewRegexp(`^\d+(\.\d{1,2})?$`, "总金额必须是数字，最多两位小数")

		status := widget.NewSelect(
			[]string{"pending", "paid", "shipping", "completed"},
			func(s string) {},
		)
		status.Selected = records[selectedRow].Status

		items := []*widget.FormItem{
			widget.NewFormItem("产品ID", productID),
			widget.NewFormItem("数量", quantity),
			widget.NewFormItem("总金额", total),
			widget.NewFormItem("状态", status),
		}

		form := dialog.NewForm("修改订单信息", "确认", "取消", items, func(b bool) {
			if !b {
				return
			}

			productIDInt, _ := strconv.Atoi(productID.Text)
			quantityInt, _ := strconv.Atoi(quantity.Text)
			totalValue, _ := strconv.ParseFloat(total.Text, 64)

			err := rep.UpdateOrderStruct(rep.DB, records[selectedRow].ID, rep.Order{
				ProductID:  productIDInt,
				Quantity:   quantityInt,
				Total:      &totalValue,
				Status:     status.Selected,
				UpdateTime: time.Now(),
			})
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// 刷新表格
			records, _ = rep.GetOrdersByUserID(rep.DB, userID)
			table.Refresh()
		}, win)

		form.Resize(fyne.NewSize(400, 300))
		form.Show()
	})

	// 布局
	searchLine := container.NewGridWithColumns(2, entry, search)
	buttons := container.NewGridWithColumns(3, addButton, deleteButton, updateButton)
	curd := container.NewVBox(widget.NewSeparator(), searchLine, buttons, widget.NewSeparator())

	return container.NewBorder(curd, nil, nil, nil, table)
}
