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
const outputBaseDir = "../assets/images"
const finalUrlPrefix = "{{ site.baseurl }}/assets/images/"

type ImageToSwap struct {
	PostPath       string
	OrigImageUrl   string
	FinalImagePath string
	FinalImageUrl  string
}

func main() {
	tasks := []ImageToSwap{}
	files, err := ioutil.ReadDir(inputBaseDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		filePath := filepath.Join(inputBaseDir, f.Name())
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}
		r := regexp.MustCompile(`\[!\[[ \w-\.]*\]\(\{\{ site\.baseurl \}\}[0-9A-Za-z\-\._:/]+( \"[0-9A-Za-z-\._:/\ \'\\]+\")?\)\]\(([0-9A-Za-z\-\@\._:/]+)\)`)
		res := r.FindAllStringSubmatch(string(fileContent), -1)
		//fmt.Println(f.Name())
		for i, item := range res {
			r := regexp.MustCompile(`.+/(\d{4})-(\d{2})-\d{2}-.+\.md$`)
			finalImagePathChunks := r.FindStringSubmatch(filePath)

			origImageUrl := res[i][len(item)-1]
			r = regexp.MustCompile(`.+/(.+)$`)
			origImageUrlChunks := r.FindStringSubmatch(origImageUrl)
			imageName := origImageUrlChunks[1]

			imageToSwap := ImageToSwap{
				PostPath:       filePath,
				OrigImageUrl:   res[i][len(item)-1],
				FinalImagePath: filepath.Join(outputBaseDir, finalImagePathChunks[1], finalImagePathChunks[2], imageName),
				FinalImageUrl:  finalUrlPrefix + filepath.Join(finalImagePathChunks[1], finalImagePathChunks[2], imageName),
			}
			//fmt.Printf("%s\n", res[i][len(item)-1])
			tasks = append(tasks, imageToSwap)
		}
	}
	for _, imageToSwap := range tasks {
		//fmt.Printf("Processing item: %+v\n", imageToSwap)
		err := downloadFile(imageToSwap.FinalImagePath, imageToSwap.OrigImageUrl)
		if err != nil {
			log.Printf("Download error for item: %+v: %s\n", imageToSwap, err.Error())
		}
		swapText(imageToSwap.PostPath, imageToSwap.OrigImageUrl, imageToSwap.FinalImageUrl)
		if err != nil {
			log.Printf("SwapText error for item: %+v: %s\n", imageToSwap, err.Error())
		}
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
