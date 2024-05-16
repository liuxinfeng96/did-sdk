/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import "errors"

// SetAdmin 设置管理员
// @params ski 与GetSenderPk()保持一致，采用公钥ski的形式
func (d *DidContract) SetAdmin(ski string) error {
	// 必须是合约创建者才能操作
	ok, err := isSenderCreator()
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("only the creator of the contract has permission")
	}

	err = d.dal.putAdmin(ski)
	if err != nil {
		return err
	}

	return nil
}

// SetAdmin 删除管理员
// @params ski 与GetSenderPk()保持一致，采用公钥ski的形式
func (d *DidContract) DeleteAdmin(ski string) error {
	// 必须是合约创建者才能操作
	ok, err := isSenderCreator()
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("only the creator of the contract has permission")
	}

	err = d.dal.deleteAdmin(ski)
	if err != nil {
		return err
	}

	return nil
}

// IsAdmin 判断是否是管理员
// @params ski 与GetSenderPk()保持一致，采用公钥ski的形式
func (d *DidContract) IsAdmin(ski string) bool {
	return d.dal.isAdmin(ski)
}
