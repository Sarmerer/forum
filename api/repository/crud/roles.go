package crud

import (
	"net/http"
)

func (UserRepoCRUD) GetRole(id int64) (int, int, error) {
	var (
		role int
		err  error
	)
	return role, http.StatusOK, err
}

//UpdateRole updates user role in the database
func (UserRepoCRUD) UpdateRole(uid int64, role int) error {

	return nil
}
