package main

import "errors"

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

func (d *DidContract) IsAdmin(ski string) bool {
	return d.dal.isAdmin(ski)
}
