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

func usersScreen(win fyne.Window) fyne.CanvasObject {
	var records []rep.User
	records, _ = rep.GetAll[rep.User](rep.DB)
	table := widget.NewTable(
		func() (int, int) { return len(records) + 1, len(rep.Members["User"]) },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell 000, 000")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			if i.Row == 0 {
				label.SetText(rep.Members["User"][i.Col])
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				record := records[i.Row-1]
				switch i.Col {
				case 0:
					label.SetText(fmt.Sprintf("%d", record.ID))
				case 1:
					label.SetText(record.Username)
				case 2:
					label.SetText("**********")
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
	entry := widget.NewEntry()
	entry.SetPlaceHolder("查询数据")
	buttons := container.NewGridWithColumns(3,
		widget.NewButtonWithIcon("增加信息", theme.ContentAddIcon(), func() {
			username := widget.NewEntry()
			username.Text = "Example User"
			username.Validator = validation.NewRegexp(`^.+$`, "请输入用户名")
			password := widget.NewPasswordEntry()
			password.Validator = validation.NewRegexp(`^[A-Za-z0-9!@#$%^&*()_+\-=[\]{};':"\|,.<>/?]+$`, "密码只能包含字母、数字、_和-")
			password.Text = "Example_Password"
			email := widget.NewEntry()
			email.Validator = validation.NewRegexp(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "电子邮件必须符合规范")
			email.Text = "Examplt@Email.com"
			phone := widget.NewEntry()
			phone.Validator = validation.NewRegexp(`^\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}$`, "电话号码必须符合规范")
			phone.Text = "123-456-7890"
			validation.NewAllStrings(
				username.Validator,
				password.Validator,
				email.Validator,
				phone.Validator,
			)
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
				if entry.Text == "" {
					records, _ = rep.GetAll[rep.User](rep.DB)
					table.Refresh()
					return
				}
				records, _ = rep.SearchVague[rep.User](rep.DB, "User", entry.Text)
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
					if !b {
						return
					}
					selectedID := records[selectedRow].ID
					fmt.Println("删除的ID是", selectedID)
					rep.DeleteID[rep.User](rep.DB, selectedID)

					if entry.Text == "" {
						records, _ = rep.GetAll[rep.User](rep.DB)
						table.Refresh()
						return
					}
					records, _ = rep.SearchVague[rep.User](rep.DB, "User", entry.Text)
					table.Refresh()
				},
				win,
			)
			cnf.SetDismissText("取消")
			cnf.SetConfirmText("确认")
			cnf.Show()
		}),
		widget.NewButtonWithIcon("修改信息", theme.DocumentCreateIcon(), func() {
			username := widget.NewEntry()
			username.Text = "Example User"
			username.Validator = validation.NewRegexp(`^.+$`, "请输入用户名")
			password := widget.NewPasswordEntry()
			password.Validator = validation.NewRegexp(`^[A-Za-z0-9!@#$%^&*()_+\-=[\]{};':"\|,.<>/?]+$`, "密码只能包含字母、数字、_和-")
			password.Text = "Example_Password"
			email := widget.NewEntry()
			email.Validator = validation.NewRegexp(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "电子邮件必须符合规范")
			email.Text = "Examplt@Email.com"
			phone := widget.NewEntry()
			phone.Validator = validation.NewRegexp(`^\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}$`, "电话号码必须符合规范")
			phone.Text = "123-456-7890"
			items := []*widget.FormItem{
				widget.NewFormItem("用户名", username),
				widget.NewFormItem("密码", password),
				widget.NewFormItem("电子邮箱", email),
				widget.NewFormItem("电话号码", phone),
			}

			username.Text = records[selectedRow].Username
			password.Text = records[selectedRow].Password
			email.Text = *records[selectedRow].Email
			phone.Text = records[selectedRow].Phone

			form := dialog.NewForm("修改信息", "确认", "取消", items, func(b bool) {
				if !b {
					return
				}
				rep.UpdateStruct[rep.User](rep.DB, records[selectedRow].ID, rep.User{
					Username: username.Text,
					Password: password.Text,
					Email:    &email.Text,
					Phone:    phone.Text,
				})

				if entry.Text == "" {
					records, _ = rep.GetAll[rep.User](rep.DB)
					table.Refresh()
					return
				}
				records, _ = rep.SearchVague[rep.User](rep.DB, "User", entry.Text)
				table.Refresh()
			}, win)
			form.Resize(fyne.NewSize(400, 480))
			form.Show()
		}),
	)
	search := widget.NewButtonWithIcon("查询", theme.SearchIcon(), func() {
		if entry.Text == "" {
			records, _ = rep.GetAll[rep.User](rep.DB)
			table.Refresh()
			return
		}
		records, _ = rep.SearchVague[rep.User](rep.DB, "User", entry.Text)
		table.Refresh()
	})
	searchLine := container.NewGridWithColumns(2, entry, search)
	curd := container.NewVBox(widget.NewSeparator(), buttons, searchLine, widget.NewSeparator())
	// 增删改查按钮

	return container.NewBorder(curd, nil, nil, nil, table)
}
