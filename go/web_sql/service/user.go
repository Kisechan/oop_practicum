package service

import (
	"encoding/json"
	"fmt"
	"web_sql/rep"
)

func CreateUser(req Request) (json.RawMessage, error) {
	var user rep.User
	err := json.Unmarshal(json.RawMessage(req.Payload), &user)
	if err != nil {
		fmt.Println("JSON unmarshal error:", err)
		return nil, err
	}

	err = rep.Create[rep.User](rep.DB, &user)
	if err != nil {
		fmt.Println("Create user error:", err)
		return nil, err
	}
	return nil, nil
}

func UpdateUser(req Request) (json.RawMessage, error) {
	var user rep.User
	err := json.Unmarshal(json.RawMessage(req.Payload), &user)
	if err != nil {
		fmt.Println("JSON unmarshal error:", err)
		return nil, err
	}

	err = rep.UpdateStruct[rep.User](rep.DB, user.ID, user)
	if err != nil {
		fmt.Println("Update user error:", err)
		return nil, err
	}
	return nil, nil
}

func ReadUser(req Request) (json.RawMessage, error) {
	var user rep.User
	err := json.Unmarshal(json.RawMessage(req.Payload), &user)
	if err != nil {
		fmt.Println("JSON unmarshal error:", err)
		return nil, err
	}

	record, err := rep.GetID[rep.User](rep.DB, user.ID)
	if err != nil {
		fmt.Println("Update user error:", err)
		return nil, err
	}
	recordJSON, err := json.Marshal(*record)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return nil, err
	}
	return recordJSON, nil
}

func DeleteUser(req Request) (json.RawMessage, error) {
	return nil, nil
}
