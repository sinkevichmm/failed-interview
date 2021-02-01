package filestorage

import "failed-interview/01/internal/models"

func metaToFile(meta *models.Meta) (file *models.File) {
	file = &models.File{}

	file.DateUpload = meta.DateUpload
	file.Extension = meta.Extension
	file.ID = meta.ID
	file.Name = meta.Name

	return file
}

func fileToMeta(file *models.File) (meta *models.Meta) {
	meta = &models.Meta{}

	meta.DateUpload = file.DateUpload
	meta.Extension = file.Extension
	meta.ID = file.ID
	meta.Name = file.Name

	return meta
}
