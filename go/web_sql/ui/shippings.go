package ui

// func shippingsScreen(win fyne.Window) fyne.CanvasObject {
// 	var records []rep.Shipping
// 	records, _ = rep.GetAll[rep.Shipping](rep.DB)
// 	table := widget.NewTable(
// 		func() (int, int) { return len(records) + 1, len(rep.Members["Shipping"]) },
// 		func() fyne.CanvasObject {
// 			return widget.NewLabel("Cell 000, 000")
// 		},
// 		func(i widget.TableCellID, o fyne.CanvasObject) {
// 			label := o.(*widget.Label)
// 			if i.Row == 0 {
// 				label.SetText(rep.Members["Shipping"][i.Col])
// 				label.TextStyle = fyne.TextStyle{Bold: true}
// 			} else {
// 				record := records[i.Row-1]
// 				switch i.Col {
// 				case 0:
// 					label.SetText(fmt.Sprintf("%d", record.ID))
// 				case 1:
// 					label.SetText(record.Shippingname)
// 				case 2:
// 					label.SetText(record.Password)
// 				case 3:
// 					label.SetText(*record.Email)
// 				case 4:
// 					label.SetText(record.Phone)
// 				}
// 			}
// 		},
// 	)
// 	table.SetColumnWidth(1, 140)
// 	table.SetColumnWidth(2, 140)
// 	table.SetColumnWidth(3, 220)

// 	var selectedRow int
// 	table.OnSelected = func(id widget.TableCellID) {

// 		if id.Row > 0 {
// 			selectedRow = id.Row - 1
// 			fmt.Printf("Selected row: %d\n", selectedRow)
// 		}
// 	}
// 	entry := widget.NewEntry()
// 	entry.SetPlaceHolder("查询数据")
// 	buttons := container.NewGridWithColumns(3,
// 		widget.NewButtonWithIcon("增加信息", theme.ContentAddIcon(), func() {
// 			shippingname := widget.NewEntry()
// 			shippingname.Text = "Example Shipping"
// 			shippingname.Validator = validation.NewRegexp(`^.+$`, "请输入用户名")
// 			password := widget.NewPasswordEntry()
// 			password.Validator = validation.NewRegexp(`^[A-Za-z0-9!@#$%^&*()_+\-=[\]{};':"\|,.<>/?]+$`, "密码只能包含字母、数字、_和-")
// 			password.Text = "Example_Password"
// 			email := widget.NewEntry()
// 			email.Validator = validation.NewRegexp(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "电子邮件必须符合规范")
// 			email.Text = "Examplt@Email.com"
// 			phone := widget.NewEntry()
// 			phone.Validator = validation.NewRegexp(`^\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}$`, "电话号码必须符合规范")
// 			phone.Text = "123-456-7890"
// 			validation.NewAllStrings(
// 				shippingname.Validator,
// 				password.Validator,
// 				email.Validator,
// 				phone.Validator,
// 			)
// 			items := []*widget.FormItem{
// 				widget.NewFormItem("用户名", shippingname),
// 				widget.NewFormItem("密码", password),
// 				widget.NewFormItem("电子邮箱", email),
// 				widget.NewFormItem("电话号码", phone),
// 			}
// 			form := dialog.NewForm("增加一条信息", "确认", "取消", items, func(b bool) {
// 				if !b {
// 					return
// 				}
// 				rep.Create[rep.Shipping](rep.DB, &rep.Shipping{
// 					Shippingname: shippingname.Text,
// 					Password: password.Text,
// 					Email:    &email.Text,
// 					Phone:    phone.Text,
// 				})
// 				if entry.Text == "" {
// 					records, _ = rep.GetAll[rep.Shipping](rep.DB)
// 					table.Refresh()
// 					return
// 				}
// 				records, _ = rep.SearchVague[rep.Shipping](rep.DB, "Shipping", entry.Text)
// 				table.Refresh()

// 			}, win)
// 			form.Resize(fyne.NewSize(400, 480))
// 			form.Show()
// 		}),
// 		widget.NewButtonWithIcon("删除信息", theme.ContentRemoveIcon(), func() {
// 			cnf := dialog.NewConfirm(
// 				"确认",
// 				"你真的要删除这条数据吗？\n你会失去它很久的（真的很久）",
// 				func(b bool) {
// 					if !b {
// 						return
// 					}
// 					selectedID := records[selectedRow].ID
// 					fmt.Println("删除的ID是", selectedID)
// 					rep.DeleteID[rep.Shipping](rep.DB, selectedID)

// 					if entry.Text == "" {
// 						records, _ = rep.GetAll[rep.Shipping](rep.DB)
// 						table.Refresh()
// 						return
// 					}
// 					records, _ = rep.SearchVague[rep.Shipping](rep.DB, "Shipping", entry.Text)
// 					table.Refresh()
// 				},
// 				win,
// 			)
// 			cnf.SetDismissText("取消")
// 			cnf.SetConfirmText("确认")
// 			cnf.Show()
// 		}),
// 		widget.NewButtonWithIcon("修改信息", theme.DocumentCreateIcon(), func() {
// 			shippingname := widget.NewEntry()
// 			shippingname.Text = "Example Shipping"
// 			shippingname.Validator = validation.NewRegexp(`^.+$`, "请输入用户名")
// 			password := widget.NewPasswordEntry()
// 			password.Validator = validation.NewRegexp(`^[A-Za-z0-9!@#$%^&*()_+\-=[\]{};':"\|,.<>/?]+$`, "密码只能包含字母、数字、_和-")
// 			password.Text = "Example_Password"
// 			email := widget.NewEntry()
// 			email.Validator = validation.NewRegexp(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "电子邮件必须符合规范")
// 			email.Text = "Examplt@Email.com"
// 			phone := widget.NewEntry()
// 			phone.Validator = validation.NewRegexp(`^\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}$`, "电话号码必须符合规范")
// 			phone.Text = "123-456-7890"
// 			items := []*widget.FormItem{
// 				widget.NewFormItem("用户名", shippingname),
// 				widget.NewFormItem("密码", password),
// 				widget.NewFormItem("电子邮箱", email),
// 				widget.NewFormItem("电话号码", phone),
// 			}

// 			shippingname.Text = records[selectedRow].Shippingname
// 			password.Text = records[selectedRow].Password
// 			email.Text = *records[selectedRow].Email
// 			phone.Text = records[selectedRow].Phone

// 			form := dialog.NewForm("修改信息", "确认", "取消", items, func(b bool) {
// 				if !b {
// 					return
// 				}
// 				rep.UpdateStruct[rep.Shipping](rep.DB, records[selectedRow].ID, rep.Shipping{
// 					Shippingname: shippingname.Text,
// 					Password: password.Text,
// 					Email:    &email.Text,
// 					Phone:    phone.Text,
// 				})

// 				if entry.Text == "" {
// 					records, _ = rep.GetAll[rep.Shipping](rep.DB)
// 					table.Refresh()
// 					return
// 				}
// 				records, _ = rep.SearchVague[rep.Shipping](rep.DB, "Shipping", entry.Text)
// 				table.Refresh()
// 			}, win)
// 			form.Resize(fyne.NewSize(400, 480))
// 			form.Show()
// 		}),
// 	)
// 	search := widget.NewButtonWithIcon("查询", theme.SearchIcon(), func() {
// 		if entry.Text == "" {
// 			records, _ = rep.GetAll[rep.Shipping](rep.DB)
// 			table.Refresh()
// 			return
// 		}
// 		records, _ = rep.SearchVague[rep.Shipping](rep.DB, "Shipping", entry.Text)
// 		table.Refresh()
// 	})
// 	searchLine := container.NewGridWithColumns(2, entry, search)
// 	curd := container.NewVBox(widget.NewSeparator(), buttons, searchLine, widget.NewSeparator())
// 	// 增删改查按钮

// 	return container.NewBorder(curd, nil, nil, nil, table)
// }
