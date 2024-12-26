package pages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// 全局变量：当前用户
var currentUser *User

// 个人主页
func CreateProfilePage() fyne.CanvasObject {
	// 初始状态：未登录
	if currentUser == nil {
		return createLoginPage()
	}

	// 登录成功后，显示用户信息
	return createUserInfoPage()
}

// 创建登录页面
func createLoginPage() fyne.CanvasObject {
	// 电话号码输入框
	phoneEntry := widget.NewEntry()
	phoneEntry.SetPlaceHolder("请输入电话号码")

	// 密码输入框
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("请输入密码")

	// 登录按钮
	loginButton := widget.NewButton("登录", func() {
		// 获取输入的电话号码和密码
		phone := phoneEntry.Text
		password := passwordEntry.Text

		// 发送登录请求
		user, err := login(phone, password)
		if err != nil {
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}

		// 登录成功，更新当前用户
		currentUser = user
		fmt.Println("User:", *user, "Logined Successfully")
		// 更新个人主页内容
		tabs := fyne.CurrentApp().Driver().AllWindows()[0].Content().(*container.AppTabs)
		tabs.Items[3].Content = createUserInfoPage()
		tabs.Refresh()
	})

	// 注册按钮
	registerButton := widget.NewButton("注册", func() {
		// 显示注册页面
		registerWindow := fyne.CurrentApp().NewWindow("注册")
		registerWindow.SetContent(createRegisterPage())
		registerWindow.Resize(fyne.NewSize(360, 780))
		registerWindow.Show()
	})

	// 布局
	return container.NewVBox(
		widget.NewRichText(
			&widget.TextSegment{
				Text: "您尚未登录\n请选择登录或注册",
				Style: widget.RichTextStyle{
					SizeName:  theme.SizeNameHeadingText,
					Alignment: fyne.TextAlignCenter,
					TextStyle: fyne.TextStyle{
						Bold: true,
					},
				},
			},
		),
		phoneEntry,
		passwordEntry,
		container.NewGridWithColumns(2, loginButton, registerButton),
	)
}

// 登录请求
func login(phone, password string) (*User, error) {
	// 创建登录请求体
	loginRequest := map[string]string{
		"phone":    phone,
		"password": password,
	}
	loginJSON, err := json.Marshal(loginRequest)
	if err != nil {
		return nil, fmt.Errorf("编码登录请求失败: %v", err)
	}

	// 发送 POST 请求
	resp, err := http.Post("http://localhost:8080/users/login", "application/json", bytes.NewBuffer(loginJSON))
	if err != nil {
		return nil, fmt.Errorf("发送登录请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("登录失败，请检查电话号码和密码")
	}

	// 解析响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}
	fmt.Printf("登录响应: %s\n", body)
	// 解析用户信息
	type LoginResponse struct {
		User User `json:"user"`
	}

	var (
		loginResponse LoginResponse
	)

	if err := json.Unmarshal(body, &loginResponse); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v", err)
	}

	return &loginResponse.User, nil
}

// 创建注册页面
func createRegisterPage() fyne.CanvasObject {
	// 用户名输入框
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("请输入用户名")

	// 密码输入框
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("请输入密码")

	// 电话号码输入框
	phoneEntry := widget.NewEntry()
	phoneEntry.SetPlaceHolder("请输入电话号码")

	// 邮箱输入框
	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("请输入邮箱（可选）")

	// 地址输入框
	addressEntry := widget.NewEntry()
	addressEntry.SetPlaceHolder("请输入地址（可选）")

	// 注册按钮
	registerButton := widget.NewButton("注册", func() {
		// 获取输入的用户信息
		username := usernameEntry.Text
		password := passwordEntry.Text
		phone := phoneEntry.Text
		email := emailEntry.Text
		address := addressEntry.Text

		// 发送注册请求
		err := register(username, password, phone, email, address)
		if err != nil {
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[1])
			return
		}

		// 注册成功，关闭注册窗口

		dialog.ShowInformation("注册成功", "请返回登录页面进行登录", fyne.CurrentApp().Driver().AllWindows()[1])
		time.Sleep(5 * time.Second)
		fyne.CurrentApp().Driver().AllWindows()[1].Close()
	})

	// 布局
	return container.NewVBox(
		widget.NewLabelWithStyle("注册", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		usernameEntry,
		passwordEntry,
		phoneEntry,
		emailEntry,
		addressEntry,
		registerButton,
	)
}

// 注册请求
func register(username, password, phone, email, address string) error {
	// 创建注册请求体
	registerRequest := User{
		Username: username,
		Password: password,
		Phone:    phone,
		Email:    &email,
		Address:  &address,
	}
	registerJSON, err := json.Marshal(registerRequest)
	if err != nil {
		return fmt.Errorf("编码注册请求失败: %v", err)
	}

	// 发送 POST 请求
	resp, err := http.Post("http://localhost:8080/users/register", "application/json", bytes.NewBuffer(registerJSON))
	if err != nil {
		return fmt.Errorf("发送注册请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("注册失败，请检查输入信息")
	}

	return nil
}

// 创建用户信息页面
func createUserInfoPage() fyne.CanvasObject {
	// 用户信息
	infoBox := container.NewVBox()

	// 添加用户信息
	addInfoRow(infoBox, "用户名		", currentUser.Username, func() {
		// 修改用户名
		newUsername := widget.NewEntry()
		dialog.ShowForm("修改用户名", "确认", "取消", []*widget.FormItem{
			widget.NewFormItem("新用户名", newUsername),
		}, func(ok bool) {
			if ok {
				currentUser.Username = newUsername.Text
				infoBox.Refresh()
			}
		}, fyne.CurrentApp().Driver().AllWindows()[0])
	})

	addInfoRow(infoBox, "电话号码		", currentUser.Phone, func() {
		// 修改电话号码
		newPhone := widget.NewEntry()
		dialog.ShowForm("修改电话号码", "确认", "取消", []*widget.FormItem{
			widget.NewFormItem("新电话号码", newPhone),
		}, func(ok bool) {
			if ok {
				currentUser.Phone = newPhone.Text
				infoBox.Refresh()
			}
		}, fyne.CurrentApp().Driver().AllWindows()[0])
	})

	if currentUser.Email != nil {
		addInfoRow(infoBox, "邮箱			", *currentUser.Email, func() {
			// 修改邮箱
			newEmail := widget.NewEntry()
			dialog.ShowForm("修改邮箱", "确认", "取消", []*widget.FormItem{
				widget.NewFormItem("新邮箱", newEmail),
			}, func(ok bool) {
				if ok {
					newEmailValue := newEmail.Text
					currentUser.Email = &newEmailValue
					infoBox.Refresh()
				}
			}, fyne.CurrentApp().Driver().AllWindows()[0])
		})
	}

	if currentUser.Address != nil {
		addInfoRow(infoBox, "地址			", *currentUser.Address, func() {
			// 修改地址
			newAddress := widget.NewEntry()
			dialog.ShowForm("修改地址", "确认", "取消", []*widget.FormItem{
				widget.NewFormItem("新地址", newAddress),
			}, func(ok bool) {
				if ok {
					newAddressValue := newAddress.Text
					currentUser.Address = &newAddressValue
					infoBox.Refresh()
				}
			}, fyne.CurrentApp().Driver().AllWindows()[0])
		})
	}

	// 布局
	return container.NewVBox(
		widget.NewRichText(
			&widget.TextSegment{
				Text: "主页",
				Style: widget.RichTextStyle{
					SizeName:  theme.SizeNameHeadingText,
					Alignment: fyne.TextAlignCenter,
					TextStyle: fyne.TextStyle{
						Bold: true,
					},
				},
			},
		),
		infoBox,
	)
}

// 添加信息行
func addInfoRow(box *fyne.Container, labelText, value string, editFunc func()) {
	row := container.NewHBox(
		widget.NewLabel(labelText+": "+value),
		widget.NewButton("修改", editFunc),
	)
	box.Add(row)
}
