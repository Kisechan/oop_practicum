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

func reviewsScreen(win fyne.Window) fyne.CanvasObject {
	var records []rep.Review
	var productID int

	// 初始化表格
	table := widget.NewTable(
		func() (int, int) { return len(records) + 1, 6 }, // 6 列：ID, UserID, ProductID, Rating, Comment, Time
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
					label.SetText("评分")
				case 4:
					label.SetText("评论")
				case 5:
					label.SetText("时间")
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
					label.SetText(record.Rating)
				case 4:
					label.SetText(record.Comment)
				case 5:
					label.SetText(record.Time.Format("2006-01-02 15:04:05"))
				}
			}
		},
	)
	table.Hide()
	table.SetColumnWidth(1, 140)

	var selectedRow int
	table.OnSelected = func(id widget.TableCellID) {
		if id.Row > 0 {
			selectedRow = id.Row - 1
			fmt.Printf("Selected row: %d\n", selectedRow)
		}
	}

	// 查找产品ID的输入框
	entry := widget.NewEntry()
	entry.SetPlaceHolder("查找产品ID")

	// 查找按钮
	search := widget.NewButtonWithIcon("查找产品", theme.SearchIcon(), func() {
		productIDStr := entry.Text
		productIDInt, err := strconv.Atoi(productIDStr)
		if err != nil {
			dialog.ShowError(fmt.Errorf("产品ID必须是整数"), win)
			return
		}
		productID = productIDInt
		records, _ = rep.GetReviewsByProductID(rep.DB, productID)
		table.Show()
		table.Refresh()
	})

	// 增加按钮
	addButton := widget.NewButtonWithIcon("增加评论", theme.ContentAddIcon(), func() {
		if productID == 0 {
			dialog.ShowInformation("错误", "请先查找产品ID", win)
			return
		}

		userID := widget.NewEntry()
		userID.SetPlaceHolder("用户ID")
		userID.Validator = validation.NewRegexp(`^\d+$`, "用户ID必须是整数")

		rating := widget.NewSelect(
			[]string{"1", "2", "3", "4", "5"},
			func(s string) {},
		)
		rating.Selected = "5"

		comment := widget.NewEntry()
		comment.SetPlaceHolder("评论内容")
		comment.Validator = validation.NewRegexp(`^.+$`, "评论内容不能为空")

		items := []*widget.FormItem{
			widget.NewFormItem("用户ID", userID),
			widget.NewFormItem("评分", rating),
			widget.NewFormItem("评论", comment),
		}

		form := dialog.NewForm("增加一条评论信息", "确认", "取消", items, func(b bool) {
			if !b {
				return
			}

			userIDInt, _ := strconv.Atoi(userID.Text)

			err := rep.Create[rep.Review](rep.DB, &rep.Review{
				UserID:    userIDInt,
				ProductID: productID,
				Rating:    rating.Selected,
				Comment:   comment.Text,
				Time:      time.Now(),
			})
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// 刷新表格
			records, _ = rep.GetReviewsByProductID(rep.DB, productID)
			table.Refresh()
		}, win)

		form.Resize(fyne.NewSize(400, 300))
		form.Show()
	})

	// 删除按钮
	deleteButton := widget.NewButtonWithIcon("删除评论", theme.ContentRemoveIcon(), func() {
		if selectedRow < 0 || selectedRow >= len(records) {
			dialog.ShowInformation("错误", "请先选择要删除的评论信息", win)
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
				rep.DeleteID[rep.Review](rep.DB, selectedID)

				// 刷新表格
				records, _ = rep.GetReviewsByProductID(rep.DB, productID)
				table.Refresh()
			},
			win,
		)
		cnf.SetDismissText("取消")
		cnf.SetConfirmText("确认")
		cnf.Show()
	})

	// 修改按钮
	updateButton := widget.NewButtonWithIcon("修改评论", theme.DocumentCreateIcon(), func() {
		if selectedRow < 0 || selectedRow >= len(records) {
			dialog.ShowInformation("错误", "请先选择要修改的评论信息", win)
			return
		}

		userID := widget.NewEntry()
		userID.SetText(fmt.Sprintf("%d", records[selectedRow].UserID))
		userID.Validator = validation.NewRegexp(`^\d+$`, "用户ID必须是整数")

		rating := widget.NewSelect(
			[]string{"1", "2", "3", "4", "5"},
			func(s string) {},
		)
		rating.Selected = records[selectedRow].Rating

		comment := widget.NewEntry()
		comment.SetText(records[selectedRow].Comment)
		comment.Validator = validation.NewRegexp(`^.+$`, "评论内容不能为空")

		items := []*widget.FormItem{
			widget.NewFormItem("用户ID", userID),
			widget.NewFormItem("评分", rating),
			widget.NewFormItem("评论", comment),
		}

		form := dialog.NewForm("修改评论信息", "确认", "取消", items, func(b bool) {
			if !b {
				return
			}

			userIDInt, _ := strconv.Atoi(userID.Text)

			err := rep.UpdateStruct(rep.DB, records[selectedRow].ID, rep.Review{
				UserID:  userIDInt,
				Rating:  rating.Selected,
				Comment: comment.Text,
				Time:    time.Now(),
			})
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// 刷新表格
			records, _ = rep.GetReviewsByProductID(rep.DB, productID)
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
