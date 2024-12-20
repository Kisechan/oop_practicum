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

func couponsScreen(win fyne.Window) fyne.CanvasObject {
	var records []rep.Coupon
	var userID int

	// 初始化表格
	table := widget.NewTable(
		func() (int, int) { return len(records) + 1, 7 }, // 7 列：ID, Code, Discount, Minimum, UserID, ExpirationDate, Status
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
					label.SetText("优惠码")
				case 2:
					label.SetText("折扣")
				case 3:
					label.SetText("最低消费")
				case 4:
					label.SetText("用户ID")
				case 5:
					label.SetText("过期日期")
				case 6:
					label.SetText("状态")
				}
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				// 数据行
				record := records[i.Row-1]
				switch i.Col {
				case 0:
					label.SetText(fmt.Sprintf("%d", record.ID))
				case 1:
					label.SetText(record.Code)
				case 2:
					label.SetText(fmt.Sprintf("%.2f", record.Discount))
				case 3:
					label.SetText(fmt.Sprintf("%.2f", record.Minimum))
				case 4:
					label.SetText(fmt.Sprintf("%d", record.UserID))
				case 5:
					label.SetText(record.ExpirationDate.Format("2006-01-02"))
				case 6:
					label.SetText(record.Status)
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
		records, _ = rep.GetCouponsByUserID(rep.DB, userID)
		table.Show()
		table.Refresh()
	})

	// 增加按钮
	addButton := widget.NewButtonWithIcon("增加优惠券", theme.ContentAddIcon(), func() {
		if userID == 0 {
			dialog.ShowInformation("错误", "请先查找用户ID", win)
			return
		}

		code := widget.NewEntry()
		code.SetPlaceHolder("优惠码")
		code.Validator = validation.NewRegexp(`^[A-Za-z0-9]{1,32}$`, "优惠码必须是1到32位的字母或数字")

		discount := widget.NewEntry()
		discount.SetPlaceHolder("折扣")
		discount.Validator = validation.NewRegexp(`^\d+(\.\d{1,2})?$`, "折扣必须是数字，最多两位小数")

		minimum := widget.NewEntry()
		minimum.SetPlaceHolder("最低消费")
		minimum.Validator = validation.NewRegexp(`^\d+(\.\d{1,2})?$`, "最低消费必须是数字，最多两位小数")

		expirationDate := widget.NewEntry()
		expirationDate.SetPlaceHolder("过期日期 (YYYY-MM-DD)")
		expirationDate.Validator = validation.NewRegexp(`^\d{4}-\d{2}-\d{2}$`, "日期格式必须是YYYY-MM-DD")

		status := widget.NewSelect(
			[]string{"available", "used", "expired"},
			func(s string) {},
		)
		status.Selected = "available"

		items := []*widget.FormItem{
			widget.NewFormItem("优惠码", code),
			widget.NewFormItem("折扣", discount),
			widget.NewFormItem("最低消费", minimum),
			widget.NewFormItem("过期日期", expirationDate),
			widget.NewFormItem("状态", status),
		}

		form := dialog.NewForm("增加一条优惠券信息", "确认", "取消", items, func(b bool) {
			if !b {
				return
			}

			discountValue, _ := strconv.ParseFloat(discount.Text, 64)
			minimumValue, _ := strconv.ParseFloat(minimum.Text, 64)
			expirationDateValue, _ := time.Parse("2006-01-02", expirationDate.Text)

			err := rep.Create[rep.Coupon](rep.DB, &rep.Coupon{
				Code:           code.Text,
				Discount:       discountValue,
				Minimum:        minimumValue,
				UserID:         userID,
				ExpirationDate: expirationDateValue,
				Status:         status.Selected,
			})
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// 刷新表格
			records, _ = rep.GetCouponsByUserID(rep.DB, userID)
			table.Refresh()
		}, win)

		form.Resize(fyne.NewSize(400, 400))
		form.Show()
	})

	// 删除按钮
	deleteButton := widget.NewButtonWithIcon("删除优惠券", theme.ContentRemoveIcon(), func() {
		if selectedRow < 0 || selectedRow >= len(records) {
			dialog.ShowInformation("错误", "请先选择要删除的优惠券信息", win)
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
				rep.DeleteID[rep.Coupon](rep.DB, selectedID)

				// 刷新表格
				records, _ = rep.GetCouponsByUserID(rep.DB, userID)
				table.Refresh()
			},
			win,
		)
		cnf.SetDismissText("取消")
		cnf.SetConfirmText("确认")
		cnf.Show()
	})

	// 修改按钮
	updateButton := widget.NewButtonWithIcon("修改优惠券", theme.DocumentCreateIcon(), func() {
		if selectedRow < 0 || selectedRow >= len(records) {
			dialog.ShowInformation("错误", "请先选择要修改的优惠券信息", win)
			return
		}

		code := widget.NewEntry()
		code.SetText(records[selectedRow].Code)
		code.Validator = validation.NewRegexp(`^[A-Za-z0-9]{1,32}$`, "优惠码必须是1到32位的字母或数字")

		discount := widget.NewEntry()
		discount.SetText(fmt.Sprintf("%.2f", records[selectedRow].Discount))
		discount.Validator = validation.NewRegexp(`^\d+(\.\d{1,2})?$`, "折扣必须是数字，最多两位小数")

		minimum := widget.NewEntry()
		minimum.SetText(fmt.Sprintf("%.2f", records[selectedRow].Minimum))
		minimum.Validator = validation.NewRegexp(`^\d+(\.\d{1,2})?$`, "最低消费必须是数字，最多两位小数")

		expirationDate := widget.NewEntry()
		expirationDate.SetText(records[selectedRow].ExpirationDate.Format("2006-01-02"))
		expirationDate.Validator = validation.NewRegexp(`^\d{4}-\d{2}-\d{2}$`, "日期格式必须是YYYY-MM-DD")

		status := widget.NewSelect(
			[]string{"available", "used", "expired"},
			func(s string) {},
		)
		status.Selected = records[selectedRow].Status

		items := []*widget.FormItem{
			widget.NewFormItem("优惠码", code),
			widget.NewFormItem("折扣", discount),
			widget.NewFormItem("最低消费", minimum),
			widget.NewFormItem("过期日期", expirationDate),
			widget.NewFormItem("状态", status),
		}

		form := dialog.NewForm("修改优惠券信息", "确认", "取消", items, func(b bool) {
			if !b {
				return
			}

			discountValue, _ := strconv.ParseFloat(discount.Text, 64)
			minimumValue, _ := strconv.ParseFloat(minimum.Text, 64)
			expirationDateValue, _ := time.Parse("2006-01-02", expirationDate.Text)

			err := rep.UpdateStruct(rep.DB, records[selectedRow].ID, rep.Coupon{
				Code:           code.Text,
				Discount:       discountValue,
				Minimum:        minimumValue,
				ExpirationDate: expirationDateValue,
				Status:         status.Selected,
			})
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// 刷新表格
			records, _ = rep.GetCouponsByUserID(rep.DB, userID)
			table.Refresh()
		}, win)

		form.Resize(fyne.NewSize(400, 400))
		form.Show()
	})

	// 布局
	searchLine := container.NewGridWithColumns(2, entry, search)
	buttons := container.NewGridWithColumns(3, addButton, deleteButton, updateButton)
	curd := container.NewVBox(widget.NewSeparator(), searchLine, buttons, widget.NewSeparator())

	return container.NewBorder(curd, nil, nil, nil, table)
}
