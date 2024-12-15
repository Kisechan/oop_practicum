package ui

import (
	"fmt"
	"strconv"
	"web_sql/rep"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func productsScreen(win fyne.Window) fyne.CanvasObject {
	var records []rep.Product
	records, _ = rep.GetAll[rep.Product](rep.DB)
	table := widget.NewTable(
		func() (int, int) { return len(records) + 1, len(rep.Members["Product"]) },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell 000, 000")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			if i.Row == 0 {
				label.SetText(rep.Members["Product"][i.Col])
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				record := records[i.Row-1]
				switch i.Col {
				case 0:
					label.SetText(fmt.Sprintf("%d", record.ID))
				case 1:
					label.SetText(record.Name)
				case 2:
					label.SetText(CutStr(record.Description, 20))
				case 3:
					label.SetText(fmt.Sprintf("%.2f", record.Price))
				case 4:
					label.SetText(fmt.Sprintf("%d", record.Stock))
				case 5:
					label.SetText(record.Type)
				case 6:
					label.SetText(record.Category)
				case 7:
					label.SetText(record.Seller)
				case 8:
					label.SetText(record.IsActive)
				case 9:
					if record.Icon != nil {
						label.SetText(CutStr(*record.Icon, 15))
					} else {
						label.SetText("")
					}
				}
			}
		},
	)
	table.SetColumnWidth(1, 140)
	table.SetColumnWidth(2, 180)
	table.SetColumnWidth(7, 280)

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
		widget.NewButtonWithIcon("增加产品", theme.ContentAddIcon(), func() {
			productName := widget.NewEntry()
			productName.Text = "Example Product"
			productName.Validator = validation.NewRegexp(`^.+$`, "请输入产品名称")

			description := widget.NewEntry()
			description.Text = "This is an example product"
			description.Validator = validation.NewRegexp(`^.+$`, "请输入产品描述")

			price := widget.NewEntry()
			price.Text = "99.99"
			price.Validator = validation.NewRegexp(`^\d+(\.\d{2})?$`, "价格必须为数字，最多两位小数")

			stock := widget.NewEntry()
			stock.Text = "100"
			stock.Validator = validation.NewRegexp(`^\d+$`, "库存必须为整数")

			typ := widget.NewSelect(
				[]string{"normal", "presale"},
				func(s string) {},
			)
			typ.Selected = "normal"

			category := widget.NewEntry()
			category.Text = "Example Category"
			category.Validator = validation.NewRegexp(`^.+$`, "请输入产品类别")

			seller := widget.NewEntry()
			seller.Text = "Example Seller"
			seller.Validator = validation.NewRegexp(`^.+$`, "请输入卖家名称")

			isactive := widget.NewSelect(
				[]string{"true", "false"},
				func(s string) {},
			)
			isactive.Selected = "true"

			icon := widget.NewEntry()

			validation.NewAllStrings(
				productName.Validator,
				description.Validator,
				price.Validator,
				stock.Validator,
				category.Validator,
				seller.Validator,
			)

			items := []*widget.FormItem{
				widget.NewFormItem("产品名称", productName),
				widget.NewFormItem("产品描述", description),
				widget.NewFormItem("价格", price),
				widget.NewFormItem("库存", stock),
				widget.NewFormItem("售出类型", typ),
				widget.NewFormItem("类别", category),
				widget.NewFormItem("卖家", seller),
				widget.NewFormItem("是否可用", isactive),
				widget.NewFormItem("图标", icon),
			}

			form := dialog.NewForm("增加一条产品信息", "确认", "取消", items, func(b bool) {
				if !b {
					return
				}

				// 将输入转换为 Product 结构体
				priceValue, _ := strconv.ParseFloat(price.Text, 64)
				stockValue, _ := strconv.Atoi(stock.Text)

				// 创建产品
				err := rep.Create[rep.Product](rep.DB, &rep.Product{
					Name:        productName.Text,
					Description: description.Text,
					Price:       priceValue,
					Stock:       stockValue,
					Category:    category.Text,
					Seller:      seller.Text,
				})
				if err != nil {
					dialog.ShowError(err, win)
					return
				}

				// 刷新表格
				records, _ = rep.GetAll[rep.Product](rep.DB)
				table.Refresh()
			}, win)

			form.Resize(fyne.NewSize(400, 600))
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
					rep.DeleteID[rep.Product](rep.DB, selectedID)

					if entry.Text == "" {
						records, _ = rep.GetAll[rep.Product](rep.DB)
						table.Refresh()
						return
					}
					records, _ = rep.SearchVague[rep.Product](rep.DB, "Product", entry.Text)
					table.Refresh()
				},
				win,
			)
			cnf.SetDismissText("取消")
			cnf.SetConfirmText("确认")
			cnf.Show()
		}),
		widget.NewButtonWithIcon("修改信息", theme.DocumentCreateIcon(), func() {
			// 检查是否有选中的行
			if selectedRow < 0 || selectedRow >= len(records) {
				dialog.ShowInformation("错误", "请先选择要修改的产品", win)
				return
			}

			// 创建输入框并绑定验证器
			productName := widget.NewEntry()
			productName.Validator = validation.NewRegexp(`^.+$`, "请输入产品名称")

			description := widget.NewEntry()
			description.Validator = validation.NewRegexp(`^.+$`, "请输入产品描述")

			price := widget.NewEntry()
			price.Validator = validation.NewRegexp(`^\d+(\.\d{2})?$`, "价格必须为数字，最多两位小数")

			stock := widget.NewEntry()
			stock.Validator = validation.NewRegexp(`^\d+$`, "库存必须为整数")

			typ := widget.NewSelect(
				[]string{"normal", "presale"},
				func(s string) {},
			)
			typ.Selected = "normal"

			category := widget.NewEntry()
			category.Validator = validation.NewRegexp(`^.+$`, "请输入产品类别")

			seller := widget.NewEntry()
			seller.Validator = validation.NewRegexp(`^.+$`, "请输入卖家名称")

			isactive := widget.NewSelect(
				[]string{"true", "false"},
				func(s string) {},
			)

			icon := widget.NewEntry()

			// 动态填充表单内容
			productName.Text = records[selectedRow].Name
			description.Text = records[selectedRow].Description
			price.Text = fmt.Sprintf("%.2f", records[selectedRow].Price)
			stock.Text = fmt.Sprintf("%d", records[selectedRow].Stock)
			category.Text = records[selectedRow].Category
			seller.Text = records[selectedRow].Seller
			isactive.Selected = records[selectedRow].IsActive
			if records[selectedRow].Icon != nil {
				icon.Text = *records[selectedRow].Icon
			}
			items := []*widget.FormItem{
				widget.NewFormItem("产品名称", productName),
				widget.NewFormItem("产品描述", description),
				widget.NewFormItem("价格", price),
				widget.NewFormItem("库存", stock),
				widget.NewFormItem("售出类型", typ),
				widget.NewFormItem("类别", category),
				widget.NewFormItem("卖家", seller),
				widget.NewFormItem("是否可用", isactive),
				widget.NewFormItem("图标", icon),
			}

			form := dialog.NewForm("修改信息", "确认", "取消", items, func(b bool) {
				if !b {
					return
				}

				// 将输入转换为 Product 结构体
				priceValue, _ := strconv.ParseFloat(price.Text, 64)
				stockValue, _ := strconv.Atoi(stock.Text)

				err := rep.UpdateStruct(rep.DB, records[selectedRow].ID, rep.Product{
					Name:        productName.Text,
					Description: description.Text,
					Price:       priceValue,
					Stock:       stockValue,
					Category:    category.Text,
					Seller:      seller.Text,
					IsActive:    isactive.Selected,
					Icon:        &icon.Text,
				})
				if err != nil {
					dialog.ShowError(err, win)
					return
				}

				// 刷新表格
				if entry.Text == "" {
					records, _ = rep.GetAll[rep.Product](rep.DB)
				} else {
					records, _ = rep.SearchVague[rep.Product](rep.DB, "Product", entry.Text)
				}
				table.Refresh()
			}, win)

			form.Resize(fyne.NewSize(400, 600))
			form.Show()
		}),
	)
	search := widget.NewButtonWithIcon("查询", theme.SearchIcon(), func() {
		if entry.Text == "" {
			records, _ = rep.GetAll[rep.Product](rep.DB)
			table.Refresh()
			return
		}
		records, _ = rep.SearchVague[rep.Product](rep.DB, "Product", entry.Text)
		table.Refresh()
	})
	searchLine := container.NewGridWithColumns(2, entry, search)
	curd := container.NewVBox(widget.NewSeparator(), buttons, searchLine, widget.NewSeparator())
	// 增删改查按钮

	return container.NewBorder(curd, nil, nil, nil, table)
}
