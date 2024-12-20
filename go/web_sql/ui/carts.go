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

func cartsScreen(win fyne.Window) fyne.CanvasObject {
	var records []rep.Cart
	var userID int

	// 初始化表格
	table := widget.NewTable(
		func() (int, int) { return len(records) + 1, 5 }, // 5 列：ID, UserID, ProductID, Quantity, AddTime
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
					label.SetText("产品ID")
				case 3:
					label.SetText("数量")
				case 4:
					label.SetText("添加时间")
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
					label.SetText(fmt.Sprintf("%d", record.ProductID))
				case 3:
					label.SetText(fmt.Sprintf("%d", record.Quantity))
				case 4:
					label.SetText(record.AddTime.Format("2006-01-02 15:04:05"))
				}
			}
		},
	)
	table.SetColumnWidth(1, 140)
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
		records, err = rep.GetCartsByUserID(rep.DB, userID)

		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		table.Show()
		table.Refresh()
	})

	// 增加按钮
	addButton := widget.NewButtonWithIcon("增加信息", theme.ContentAddIcon(), func() {
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

		items := []*widget.FormItem{
			widget.NewFormItem("产品ID", productID),
			widget.NewFormItem("数量", quantity),
		}

		form := dialog.NewForm("增加一条购物车信息", "确认", "取消", items, func(b bool) {
			if !b {
				return
			}

			productIDInt, _ := strconv.Atoi(productID.Text)
			quantityInt, _ := strconv.Atoi(quantity.Text)

			err := rep.Create[rep.Cart](rep.DB, &rep.Cart{
				UserID:    userID,
				ProductID: productIDInt,
				Quantity:  quantityInt,
				AddTime:   time.Now(),
			})
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// 刷新表格
			records, _ = rep.GetCartsByUserID(rep.DB, userID)
			table.Refresh()
		}, win)

		form.Resize(fyne.NewSize(400, 300))
		form.Show()
	})

	// 删除按钮
	deleteButton := widget.NewButtonWithIcon("删除信息", theme.ContentRemoveIcon(), func() {
		if selectedRow < 0 || selectedRow >= len(records) {
			dialog.ShowInformation("错误", "请先选择要删除的购物车信息", win)
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
				rep.DeleteID[rep.Cart](rep.DB, selectedID)

				// 刷新表格
				records, _ = rep.GetCartsByUserID(rep.DB, userID)
				table.Refresh()
			},
			win,
		)
		cnf.SetDismissText("取消")
		cnf.SetConfirmText("确认")
		cnf.Show()
	})

	// 修改按钮
	updateButton := widget.NewButtonWithIcon("修改信息", theme.DocumentCreateIcon(), func() {
		if selectedRow < 0 || selectedRow >= len(records) {
			dialog.ShowInformation("错误", "请先选择要修改的购物车信息", win)
			return
		}

		productID := widget.NewEntry()
		productID.SetText(fmt.Sprintf("%d", records[selectedRow].ProductID))
		productID.Validator = validation.NewRegexp(`^\d+$`, "产品ID必须是整数")

		quantity := widget.NewEntry()
		quantity.SetText(fmt.Sprintf("%d", records[selectedRow].Quantity))
		quantity.Validator = validation.NewRegexp(`^\d+$`, "数量必须是整数")

		items := []*widget.FormItem{
			widget.NewFormItem("产品ID", productID),
			widget.NewFormItem("数量", quantity),
		}

		form := dialog.NewForm("修改购物车信息", "确认", "取消", items, func(b bool) {
			if !b {
				return
			}

			productIDInt, _ := strconv.Atoi(productID.Text)
			quantityInt, _ := strconv.Atoi(quantity.Text)

			err := rep.UpdateStruct(rep.DB, records[selectedRow].ID, rep.Cart{
				ProductID: productIDInt,
				Quantity:  quantityInt,
				AddTime:   time.Now(),
			})
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// 刷新表格
			records, _ = rep.GetCartsByUserID(rep.DB, userID)
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
