package main

import (
	"errors"
	"os"
	"bufio"
)

var(
	ValidArguments = []string{"find","remove"}
	ErrInvalidArguments = errors.New("Invalid arguments")
	ErrParams = errors.New("A target file must be selected")
	ErrFileNotFound = errors.New("File not found")
)

func main() {
	if err := Init(); err == ErrParams || err == ErrInvalidArguments{
		print(err)
		os.Exit(2)
	} else if err != nil {
		print(err)
		os.Exit(1)
	}
}

func Init() error{
	command, err := DefineCommandStrategy()
	if err != nil{
		return err
	}
	if err := ApplyCommandStrategy(command); err!= nil{
		return err
	}
	return nil
}

func DefineCommandStrategy() (string, error){
	if (len(os.Args[1:]) < 2){
		return "", ErrInvalidArguments
	}
	a := os.Args[1]
	for i := 0; i < len(ValidArguments); i++{
		if (a == ValidArguments[i]){
			return a , nil
		}
	}
	return "", ErrInvalidArguments
}

func ApplyCommandStrategy(command string) error{
	switch command {
	case "find":
		if err := FindEntriesOnTarget(); err != nil{
			return err
		}
		return nil
	case "remove":
		if err := RemoveEntriesFromTarget(); err != nil{
			return err
		}
		return nil
	}
	return nil
}

func FindEntriesOnTarget() error{
	e, t, err := ReadFiles()
	if err != nil{
		return err
	}
	c := RemoveDuplicates(t)
	for i := 0; i < len(e); i++{
		for j := 0; j < len(c); j++{
			if e[i] == c[j]{
				c = append(c, c[j])
			}
		}
	}
	o := CreateOutputFile(c)
	return o
}

func RemoveEntriesFromTarget() error{
	e, t, err := ReadFiles()
	if err != nil{
		return err
	}
	c := RemoveDuplicates(t)
	for i := 0; i < len(e); i++{
		for j := 0; j < len(c); j++{
			if e[i] == c[j]{
				c = append(c[:j], c[j+1:]...)
			}
		}
	}
	o := CreateOutputFile(c)
	return o
}

func ReadFiles() ([] string, [] string, error){
	e, err := ReadEntryFile(os.Args[2])
	if err != nil {
		return []string{}, []string{}, err
	}
	t := []string{}
	for i := 3; i < len(os.Args); i++{
		a, err := ReadEntryFile(os.Args[i])
		if err != nil {
			return []string{}, []string{}, err
		}
		t = append(t, a...)
	}
	return e, t, nil
}

func ReadEntryFile(f string) ([]string, error){
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entryArgs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entryArgs = append(entryArgs, scanner.Text())
	}
	return entryArgs, scanner.Err()
}

func RemoveDuplicates(s []string) []string {
    seen := make(map[string]struct{}, len(s))
    j := 0
    for _, v := range s {
        if _, ok := seen[v]; ok {
            continue
        }
        seen[v] = struct{}{}
        s[j] = v
        j++
    }
    return s[:j]
}

func CreateOutputFile(r []string) error{
	file, err := os.Create("csvtoolresult.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	for _, line := range r {
		println(line)
		w.WriteString(line+ "\n")
	}
	w.Flush()
	file.Close()

	return nil
}