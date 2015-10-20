package models

func Insert(sec interface{}) error {
	_, err := orm.Insert(sec)
	return err
}

func FindPager(ptr2slices interface{}, start, limit int) error {
	return orm.Limit(limit, start).Find(ptr2slices)
}

func GetById(id interface{}, bean interface{}) (bool, error) {
	return orm.Id(id).Get(bean)
}

func GetByExample(bean interface{}) (bool, error) {
	return orm.Get(bean)
}

func FindByExample(ptr2slices interface{}, bean interface{}) error {
	return orm.Find(ptr2slices, bean)
}

func DelByExample(bean interface{}) error {
	_, err := orm.Delete(bean)
	return err
}

func DelById(id, bean interface{}) error {
	_, err := orm.Id(id).Delete(bean)
	return err
}

func UpdateById(id interface{}, bean interface{}, cols ...string) error {
	_, err := orm.Id(id).Cols(cols...).Update(bean)
	return err
}
