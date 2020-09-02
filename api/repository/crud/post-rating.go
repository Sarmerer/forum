package crud

// func (PostRepoCRUD) UpdateRating(upORdown string, pid uint64) error {
// 	var (
// 		result sql.Result
// 		err    error
// 	)
// 	switch upORdown {
// 	case "up":
// 		if result, err = repository.DB.Exec(
// 			"UPDATE posts SET rating = rating + 1 WHERE id = ?",
// 			pid,
// 		); err != nil {

// 		}
// 	case "down":
// 		if result, err = repository.DB.Exec(
// 			"UPDATE posts SET rating = rating - 1 WHERE id = ?",
// 			pid,
// 		); err != nil {

// 		}
// 	default:
// 	}
// 	return nil
// }
