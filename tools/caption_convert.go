package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const inputBaseDir = "../_posts"

type swap struct {
	path string
	from string
	to   string
}

func main() {
	tasks := []swap{}
	files, err := ioutil.ReadDir(inputBaseDir)
	if err != nil {
		log.Fatal(err)
	}
	total := 0
	for _, f := range files {
		filePath := filepath.Join(inputBaseDir, f.Name())
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}
		r := regexp.MustCompile(`\[caption.+caption=\"(.+)\"\](\[.+)\[\/caption\]`)
		res := r.FindAllStringSubmatch(string(fileContent), -1)
		fmt.Println(f.Name())
		for _, item := range res {
			//fmt.Printf("  T1 *%s* *%s*\n", item[1], item[2])
			tasks = append(tasks, swap{
				path: filePath,
				from: item[0],
				to:   fmt.Sprintf("%s\n*%s*", item[2], item[1]),
			})
			total = total + 1
		}

		r = regexp.MustCompile(`\[caption.+\](\[.+\)) (.+)\[\/caption\]`)
		res = r.FindAllStringSubmatch(string(fileContent), -1)
		fmt.Println(f.Name())
		for _, item := range res {
			//fmt.Printf("  T2 *%s* *%s*\n", item[2], item[1])
			tasks = append(tasks, swap{
				path: filePath,
				from: item[0],
				to:   fmt.Sprintf("%s\n*%s*", item[1], item[2]),
			})
			total = total + 1
		}
	}
	fmt.Println("total:", total)
	for _, task := range tasks {
		fmt.Printf("%+v\n", task)
		swapText(task.path, task.from, task.to)
	}
}

func swapText(filePath, from, to string) error {
	read, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("ReadFile error on file: %s: %s", filePath, err.Error())
	}
	//fmt.Println(filePath)

	newContents := strings.Replace(string(read), from, to, -1)

	err = ioutil.WriteFile(filePath, []byte(newContents), 0)
	if err != nil {
		return fmt.Errorf("Writefile error on file: %s: %s", filePath, err.Error())
	}
	return nil
}

func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		fmt.Printf("downloadFile non200 error for file: %s: %d\n", filepath, resp.StatusCode)
	}
	defer resp.Body.Close()
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Write the body to file
	//_, err = io.Copy(ioutil.Discard, resp.Body)
	_, err = io.Copy(out, resp.Body)
	return err
}
