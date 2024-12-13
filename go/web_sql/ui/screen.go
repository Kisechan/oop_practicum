package ui

import (
	"fmt"
	"web_sql/rep"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var members = map[string][]string{
	"Cart": {
		"ID",
		"UserID",
		"ProductID",
		"Quantity",
		"AddTime",
	},

	" Category ": {
		"ID",
		"Name",
		"ParentID",
	},

	" Coupon ": {
		"ID",
		"Code",
		"Type",
		"Discount",
		"Minimum",
		"UserID",
		"ProductID",
		"ExpirationDate",
		"Status",
	},

	" DeliveryAddress ": {
		"ID",
		"UserID",
		"Phone",
		"Address",
		"Name",
	},

	"OrderItem ": {
		"ID",
		"OrderID",
		"ProductID",
		"Quantity",
		"UnitPrice",
		"TotalPrice",
	},

	" Order ": {
		"ID",
		"UserID",
		"Total",
		"Status",
		"CreatedTime",
		"UpdateTime",
	},

	" Product ": {
		"ID",
		"Name",
		"Description",
		"Price",
		"Stock",
		"Type",
		"CategoryID",
		"Seller",
		"IsActive",
	},

	" Review ": {
		"ID",
		"UserID",
		"ProductID",
		"Rating",
		"Comment",
		"Time",
	},

	"Shipping": {
		"ID",
		"OrderItemID",
		"TrackingNumber",
		"Carrier",
		"Status",
		"EstimatedDeliveredTime",
		"CreateTime",
		"ShippedTime",
		"CompletedTime",
	},

	"User": {
		"ID",
		"Username",
		"Password",
		"Email",
		"Phone",
	},
}

func welcomeScreen(win fyne.Window) fyne.CanvasObject {

	text := widget.NewRichText(
		&widget.TextSegment{
			Text: "欢迎",
			Style: widget.RichTextStyle{
				SizeName: theme.SizeNameHeadingText,
			},
		},
	)
	content := container.NewBorder(text, nil, nil, nil)
	return container.NewCenter(content)
}

func usersScreen(win fyne.Window) fyne.CanvasObject {
	records, _ := rep.GetAll[rep.User](rep.DB)
	table := widget.NewTable(
		func() (int, int) { return len(records) + 1, len(members["User"]) },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell 000, 000")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			records, _ = rep.GetAll[rep.User](rep.DB)
			if i.Row == 0 {
				label.SetText(members["User"][i.Col])
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				record := records[i.Row-1]
				switch i.Col {
				case 0:
					label.SetText(fmt.Sprintf("%d", record.ID))
				case 1:
					label.SetText(record.Username)
				case 2:
					label.SetText(record.Password)
				case 3:
					label.SetText(*record.Email)
				case 4:
					label.SetText(record.Phone)
				}
			}
		},
	)
	table.SetColumnWidth(1, 140)
	table.SetColumnWidth(2, 140)
	table.SetColumnWidth(3, 220)

	var selectedRow int
	table.OnSelected = func(id widget.TableCellID) {

		if id.Row > 0 {
			selectedRow = id.Row - 1
			fmt.Printf("Selected row: %d\n", selectedRow)
		}
	}
	buttons := container.NewGridWithColumns(3,
		widget.NewButtonWithIcon("增加信息", theme.ContentAddIcon(), func() {
			username := widget.NewEntry()
			password := widget.NewPasswordEntry()
			password.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "密码只能包含字母、数字、_和-")
			email := widget.NewEntry()
			email.Validator = validation.NewRegexp(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "电子邮件必须符合规范")
			phone := widget.NewEntry()
			phone.Validator = validation.NewRegexp(`^[0-9_()-]+$`, "电话号码只能含有数字、括号和-")
			items := []*widget.FormItem{
				widget.NewFormItem("用户名", username),
				widget.NewFormItem("密码", password),
				widget.NewFormItem("电子邮箱", email),
				widget.NewFormItem("电话号码", phone),
			}

			form := dialog.NewForm("增加一条信息", "确认", "取消", items, func(b bool) {
				if !b {
					return
				}
				rep.Create[rep.User](rep.DB, &rep.User{
					Username: username.Text,
					Password: password.Text,
					Email:    &email.Text,
					Phone:    phone.Text,
				})
				table.Refresh()
			}, win)
			form.Resize(fyne.NewSize(400, 480))
			form.Show()
		}),
		widget.NewButtonWithIcon("删除信息", theme.ContentRemoveIcon(), func() {
			cnf := dialog.NewConfirm(
				"确认",
				"你真的要删除这条数据吗？\n你会失去它很久的（真的很久）",
				func(b bool) {
					selectedID := records[selectedRow].ID
					fmt.Println("删除的ID是", selectedID)
					rep.DeleteID[rep.User](rep.DB, selectedID)
					table.Refresh()
				},
				win,
			)
			cnf.SetDismissText("取消")
			cnf.SetConfirmText("确认")
			cnf.Show()
		}),
		widget.NewButtonWithIcon("修改信息", theme.DocumentCreateIcon(), func() {
		}),
	)
	entry := widget.NewEntry()
	entry.SetPlaceHolder("查询数据")
	search := widget.NewButtonWithIcon("查询", theme.SearchIcon(), func() {
		rep.GetField[rep.User](rep.DB, "ID", entry.Text)
	})
	searchLine := container.NewGridWithColumns(2, entry, search)
	curd := container.NewVBox(widget.NewSeparator(), buttons, searchLine, widget.NewSeparator())
	// 增删改查按钮

	return container.NewBorder(curd, nil, nil, nil, table)
}
