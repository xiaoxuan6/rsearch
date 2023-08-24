package common

type Model struct {
	Title string `json:"title"`
	Tag   string `json:"tag"`
	Url   string `json:"url"`
}

func FetchModelByTitle(title string) (model []*Model, err error) {
	err = DB.Where("title LICK %?%", title).Select([]string{"title", "tag", "url"}).Find(&model).Error
	return model, err
}

func FetchModelByTag(tag string) (model []*Model, err error) {
	err = DB.Where("tag = ?", tag).Select([]string{"title", "tag", "url"}).Find(&model).Error
	return model, err
}

func Search(keyword string) (model []*Model, err error) {
	err = DB.Where("title LIKE ? or tag LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&model).Error
	return model, err
}

func All() (model []*Model, err error) {
	err = DB.Find(&model).Error
	return model, err
}

func Count() int64 {
	var count int64
	_ = DB.Model(&Model{}).Count(&count).Error
	return count
}

func (model *Model) Insert() (err error) {
	err = DB.Create(model).Error
	return err
}

func CreateInBatches(models []Model) (err error) {
	err = DB.CreateInBatches(models, 100).Error
	return err
}
