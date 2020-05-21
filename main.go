package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

// генерация символов разделителей
func makeSeparator(level int, isLast bool, isLastParent bool) string {
	separator := "───"
	var result, separatorStart string
	// если последний элемент в списке
	if isLast == true {
		separatorStart = "└"
	} else {
		separatorStart = "├"
	}
	for i := 0; i < level; i++ {
		// todo: разобраться как сделать правильные префиксы у закрывающий элементов
		// последний в уровне; последний в уровне, где родитель последний в своем уровне
		if (isLastParent || isLast) && i == 0 { // если первый элемент в списке
			result += "|\t"
		} else if isLastParent && isLast {
			result += "\t"
		} else {
			result += "│\t"
		}
	}
	return result + separatorStart + separator
}

// генерация имени файла
func makeFileName(fi os.FileInfo, postfix string) string {
	if fi.Size() > 0 { // если размер ненулевой
		return fi.Name() + " (" + strconv.Itoa(int(fi.Size())) + postfix + ")"
	}
	return fi.Name() + " (empty)"
}

// сортировка по имени слайса с файловыми данными
type ByName []os.FileInfo

// дефолтные методы
func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

// нормализация слайса с файловыми данными
func normalizedFiles(fileInfos []os.FileInfo, printFiles bool) []os.FileInfo {
	var result []os.FileInfo
	if printFiles == false { // убираем файлы, если нужны только директории
		result = getOnlyDirs(fileInfos)
	} else {
		result = fileInfos
	}
	// сортируем
	sort.Sort(ByName(result))
	return result
}

// оставляем только директории, файлы убираем
func getOnlyDirs(fileInfos []os.FileInfo) []os.FileInfo {
	var result []os.FileInfo
	for _, fi := range fileInfos {
		if fi.IsDir() {
			result = append(result, fi)
		}
	}
	return result
}

// рекурсивный метод для работы с папкой
func currentDirTree(writer io.Writer, path string, printFiles bool, level int) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	normalized := normalizedFiles(fileInfos, printFiles)
	for idx, fi := range normalized {
		//todo: разобраться с рисованием символов - описать логику
		lastParent := false
		if fi.IsDir() && idx == len(normalized)-1 {
			lastParent = true
		}
		if fi.IsDir() {
			separator := makeSeparator(level, idx == len(normalized)-1, lastParent)
			fmt.Fprintln(writer, separator+fi.Name())
			currentDirTree(writer, path+string(os.PathSeparator)+fi.Name(), printFiles, level+1)
		} else if !fi.IsDir() && printFiles == true {
			separator := makeSeparator(level, idx == len(normalized)-1, false)
			fmt.Fprintln(writer, separator+makeFileName(fi, "b"))
		} else {
			continue
		}
	}
	return nil
}

// точка входа
func dirTree(writer io.Writer, path string, printFiles bool) error {
	err := currentDirTree(writer, path, printFiles, 0)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	out := os.Stdout
	// проверка аргументов командной строки
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	// путь с которым работаем
	path := os.Args[1]
	// показывать файлы или нет
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	// вызов функции
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
