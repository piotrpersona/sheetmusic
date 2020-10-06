package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

var (
	repoURL = "https://github.com/piotrpersona/sheetmusic"
	branch = "main"
	directory = "pdf"
	readmeFile = "README.md"
	readmeTemplate = "templates/README.md.tmpl"
)

type PDF struct {
	Name        string
	DownloadURL string
}

type TemplateData struct {
	PDFs []PDF
}

func panicErr(err error) {
	if err != nil {
		log.Printf("Panic: %v\n", err)
		panic(err)
	}
}

func listDirectory(rootPath string) (files []string, err error) {
	fileInfos, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return
	}
	for _, info := range fileInfos {
		if info.IsDir() {
			continue
		}
		fileName := info.Name()
		if matched, err := filepath.Match("*.pdf", fileName); matched && err == nil {
			files = append(files, fileName)
		}
	}
	return
}

func createRawURL(name string) (raw string) {
	return fmt.Sprintf("%s/raw/%s/%s/%s", repoURL, branch, directory, name)
}

func constructPDFs(names []string) (pdfs []PDF, err error) {
	for _, name := range names {
		rawURL := createRawURL(name)
		parsedURL, parseErr := url.Parse(rawURL)
		if parseErr != nil {
			err = parseErr
			return
		}
		pdf := PDF{
			Name:        name,
			DownloadURL: parsedURL.String(),
		}
		pdfs = append(pdfs, pdf)
	}
	return
}

func templateREADME(writer io.Writer, templateData TemplateData) (err error) {
	tmpl, err := template.New("README.md.tmpl").ParseFiles(readmeTemplate)
	if err != nil {
		return
	}
	err = tmpl.Execute(writer, templateData)
	return
}

func main() {
	files, err := listDirectory("pdf")
	panicErr(err)

	pdfs, err := constructPDFs(files)
	panicErr(err)

	fileWriter, err := os.OpenFile(readmeFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	panicErr(err)
	if fileWriter != nil {
		defer func() {
			err = fileWriter.Close()
			panicErr(err)
		}()
	}
	templateData := TemplateData{PDFs: pdfs}
	err = templateREADME(fileWriter, templateData)
	panicErr(err)
}
