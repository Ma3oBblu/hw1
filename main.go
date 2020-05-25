package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

// генерация графического символа для текущего элемента
func makeSymbol(isLast bool) string {
	if isLast == true { // если элемент последний
		return "└───"
	}
	return "├───"
}

// генерация символа табуляции для текущего элемента
func makePrefix(prevPrefix string, isLastParent bool, level int) string {
	if level == 0 { // если 0 уровень - пустая строка
		return ""
	}
	if isLastParent == true { // если родительский элемент последний
		return prevPrefix + "\t"
	}
	return prevPrefix + "│\t"
}

// генерация имени файла
func makeFileName(fi os.FileInfo, postfix string, prefix string, symbol string) string {
	if fi.Size() > 0 { // если размер ненулевой
		return prefix + symbol + fi.Name() + " (" + strconv.Itoa(int(fi.Size())) + postfix + ")"
	}
	return prefix + symbol + fi.Name() + " (empty)"
}

// генерации имени директории
func makeDirName(fi os.FileInfo, prefix string, symbol string) string {
	return prefix + symbol + fi.Name()
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
func currentDirTree(writer io.Writer, path string, printFiles bool, level int, isLastParent bool, prevSeparator string) error {
	dir, err := os.Open(path) // читаем текущий путь
	if err != nil {
		return err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1) // читаем текущую директорию
	if err != nil {
		return err
	}

	normalized := normalizedFiles(fileInfos, printFiles) // нормализация слайса с текущими элементами
	isLast := false                                      // проверка на последний элемент в слайсе
	for idx, fi := range normalized {
		if idx == len(normalized)-1 { // если элемент последний в текущей слайсе
			isLast = true
		}
		prefix := makePrefix(prevSeparator, isLastParent, level)
		symbol := makeSymbol(isLast)
		if fi.IsDir() { // если директория
			fmt.Fprintln(writer, makeDirName(fi, prefix, symbol))
			currentPath := path + string(os.PathSeparator) + fi.Name() // текущий путь, который передаем дальше рекурсивно
			currentDirTree(writer, currentPath, printFiles, level+1, isLast, prefix)
		} else if !fi.IsDir() && printFiles == true {
			fmt.Fprintln(writer, makeFileName(fi, "b", prefix, symbol))
		} else {
			continue
		}
	}
	return nil
}

// точка входа
func dirTree(writer io.Writer, path string, printFiles bool) error {
	err := currentDirTree(writer, path, printFiles, 0, false, "")
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
