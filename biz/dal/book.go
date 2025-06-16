package dal

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"hertz_template/biz/dbmodel"
)

func CreateBook(books []*dbmodel.Book) error {
	return DB.Create(books).Error
}

func DeleteBook(bookID uint) error {
	var book dbmodel.Book
	if err := DB.First(&book, "id = ?", bookID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("书不存在或已被删除")
		}
		return err
	}
	return DB.Delete(&book).Error
}

func UpdateBook(book *dbmodel.Book) error {
	if err := DB.First(&book, "id = ?", &book.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("书不存在或已被删除")
		}
		return err
	}
	return DB.Save(book).Error
}

func GetBookByID(bookID uint) (*dbmodel.Book, error) {
	var book dbmodel.Book
	if err := DB.Where("id = ?", bookID).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在时返回 nil
		}
		return nil, err // 其他错误
	}
	return &book, nil
}

func GetBookList(pageSize, offset int, title, author string) ([]*dbmodel.Book, int64, error) {
	var books []*dbmodel.Book
	query := DB.Model(&dbmodel.Book{})

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if author != "" {
		query = query.Where("author LIKE ?", "%"+author+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("id").Limit(pageSize).Offset(offset).Find(&books).Error
	return books, total, err
}

func IsBookExists(title string) (bool, error) {
	var count int64
	err := DB.Model(&dbmodel.Book{}).Where("title = ?", title).Count(&count).Error
	return count > 0, err
}
