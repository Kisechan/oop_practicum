package pages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
		// fmt.Println("User:", *user, "Login Successfully")
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
		registerWindow.Resize(fyne.NewSize(680, 400))
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

	var loginResponse LoginResponse
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
		// time.Sleep(5 * time.Second)
		// fyne.CurrentApp().Driver().AllWindows()[1].Close()
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

// 创建用户信息展示区域
func createInfoBox() *fyne.Container {
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
				// 更新服务器上的用户信息
				if err := updateUserInfo(currentUser); err != nil {
					dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
				} else {
					dialog.ShowInformation("更新成功", "用户信息已更新", fyne.CurrentApp().Driver().AllWindows()[0])
					// 重新生成 infoBox
					infoBox.Refresh()
				}
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
				// 更新服务器上的用户信息
				if err := updateUserInfo(currentUser); err != nil {
					dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
				} else {
					dialog.ShowInformation("更新成功", "用户信息已更新", fyne.CurrentApp().Driver().AllWindows()[0])
					// 重新生成 infoBox
					infoBox.Refresh()
				}
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
					// 更新服务器上的用户信息
					if err := updateUserInfo(currentUser); err != nil {
						dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
					} else {
						dialog.ShowInformation("更新成功", "用户信息已更新", fyne.CurrentApp().Driver().AllWindows()[0])
						// 重新生成 infoBox
						infoBox.Refresh()
					}
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
					// 更新服务器上的用户信息
					if err := updateUserInfo(currentUser); err != nil {
						dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
					} else {
						dialog.ShowInformation("更新成功", "用户信息已更新", fyne.CurrentApp().Driver().AllWindows()[0])
						// 重新生成 infoBox
						infoBox.Refresh()
					}
				}
			}, fyne.CurrentApp().Driver().AllWindows()[0])
		})
	}

	// 修改密码按钮
	changePasswordButton := widget.NewButton("修改密码", func() {
		// 弹窗表单
		oldPasswordEntry := widget.NewPasswordEntry()
		oldPasswordEntry.SetPlaceHolder("请输入旧密码")
		oldPasswordEntry.MinSize()                                                    // 设置最小宽度
		oldPasswordEntry.Resize(fyne.NewSize(300, oldPasswordEntry.MinSize().Height)) // 设置宽度为 300

		newPasswordEntry := widget.NewPasswordEntry()
		newPasswordEntry.SetPlaceHolder("请输入新密码")
		newPasswordEntry.Resize(fyne.NewSize(300, newPasswordEntry.MinSize().Height)) // 设置宽度为 300

		confirmPasswordEntry := widget.NewPasswordEntry()
		confirmPasswordEntry.SetPlaceHolder("请确认新密码")
		confirmPasswordEntry.Resize(fyne.NewSize(300, confirmPasswordEntry.MinSize().Height)) // 设置宽度为 300

		dialog.ShowForm("修改密码", "确认", "取消", []*widget.FormItem{
			widget.NewFormItem("", oldPasswordEntry),
			widget.NewFormItem("", newPasswordEntry),
			widget.NewFormItem("", confirmPasswordEntry),
		}, func(ok bool) {
			if ok {
				// 检查新密码和确认密码是否一致
				if newPasswordEntry.Text != confirmPasswordEntry.Text {
					dialog.ShowError(fmt.Errorf("新密码与确认密码不一致"), fyne.CurrentApp().Driver().AllWindows()[0])
					return
				}

				// 发送修改密码请求
				err := changePassword(currentUser.ID, oldPasswordEntry.Text, newPasswordEntry.Text)
				if err != nil {
					dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
				} else {
					dialog.ShowInformation("修改成功", "密码已更新", fyne.CurrentApp().Driver().AllWindows()[0])
				}
			}
		}, fyne.CurrentApp().Driver().AllWindows()[0])
	})
	infoBox.Add(container.NewGridWithColumns(
		2,
		widget.NewLabel("修改密码"),
		changePasswordButton,
	))

	return infoBox
}

// 创建用户信息页面
func createUserInfoPage() fyne.CanvasObject {
	// 用户信息
	infoBox := createInfoBox()

	// 刷新按钮
	refreshButton := widget.NewButtonWithIcon("刷新", theme.ViewRefreshIcon(), func() {
		// 重新获取用户信息
		// user, err := fetchUserInfo(currentUser.ID)
		// if err != nil {
		// 	dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		// 	return
		// }

		// 更新当前用户信息
		// currentUser = user

		// 重新生成 infoBox
		infoBox = createInfoBox()

		// 刷新页面
		tabs := fyne.CurrentApp().Driver().AllWindows()[0].Content().(*container.AppTabs)
		tabs.Items[3].Content = createUserInfoPage()
		tabs.Refresh()
	})

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
		refreshButton,
		infoBox,
	)
}

// 添加信息行
func addInfoRow(box *fyne.Container, labelText, value string, editFunc func()) {
	row := container.NewGridWithColumns(
		2,
		widget.NewLabel(labelText+value),
		widget.NewButton("修改", func() {
			editFunc()
		}),
	)
	box.Add(row)
}

// 更新用户信息
func updateUserInfo(user *User) error {
	// 创建更新请求体
	updateRequest := struct {
		ID       int     `json:"id"`
		Username string  `json:"username"`
		Email    *string `json:"email"`
		Phone    string  `json:"phone"`
		Address  *string `json:"address"`
	}{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Address:  user.Address,
	}

	updateJSON, err := json.Marshal(updateRequest)
	if err != nil {
		return fmt.Errorf("编码更新请求失败: %v", err)
	}

	// 发送 PUT 请求
	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/users/profile", bytes.NewBuffer(updateJSON))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("更新用户信息失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

// 修改密码请求
func changePassword(userID int, oldPassword, newPassword string) error {
	// 创建修改密码请求体
	changePasswordRequest := struct {
		UserID      int    `json:"user_id"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}{
		UserID:      userID,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	changePasswordJSON, err := json.Marshal(changePasswordRequest)
	if err != nil {
		return fmt.Errorf("编码修改密码请求失败: %v", err)
	}

	// 发送 PUT 请求
	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/users/change-password", bytes.NewBuffer(changePasswordJSON))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("修改密码失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

// 获取用户信息
// func fetchUserInfo(userID int) (*User, error) {
// 	// 发送 GET 请求
// 	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/users/profile/%d", userID))
// 	if err != nil {
// 		return nil, fmt.Errorf("发送请求失败: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// 检查响应状态码
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("获取用户信息失败，状态码: %d", resp.StatusCode)
// 	}

// 	// 解析响应体
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("读取响应体失败: %v", err)
// 	}

// 	// 解析用户信息
// 	var user User
// 	if err := json.Unmarshal(body, &user); err != nil {
// 		return nil, fmt.Errorf("解析用户信息失败: %v", err)
// 	}

// 	return &user, nil
// }
