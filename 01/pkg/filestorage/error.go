package filestorage

import "errors"

var ErrMetaRead = errors.New("can not read meta file")
var ErrMetaCreat = errors.New("can not create meta file")
var ErrMetaWrite = errors.New("can not write meta file")
var ErrMetaInvalid = errors.New("meta file is invalid")
var ErrMetaDouble = errors.New("meta file has duplicates")
var ErrFileSave = errors.New("can not save file")
var ErrFileDelete = errors.New("can not delete file")
var ErrFileRead = errors.New("can not read file")
var ErrStorageFull = errors.New("storage is full")
var ErrIDExist = errors.New("id exists")
var ErrIDNotFound = errors.New("id not found")
