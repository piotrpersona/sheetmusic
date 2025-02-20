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
	"strings"
)

var (
	repoURL = "https://github.com/piotrpersona/sheetmusic"
	branch = "main"
	directory = "pdf"
	readmeFile = "README.md"
	ghPagesIndex = "docs/index.md"
	readmeTemplate = "templates/README.md.tmpl"
	ghPagesTemplate = "templates/index.md.tmpl"
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

func templateDocument(templateFile string, templateData TemplateData, writers ...io.Writer) (err error) {
	tmpl, err := template.New(strings.Split(templateFile, "/")[1]).ParseFiles(templateFile)
	if err != nil {
		return
	}
	for _, writer := range writers {
		err = tmpl.Execute(writer, templateData)
		if err != nil {
			return
		}
	}
	return
}

func main() {
	files, err := listDirectory("pdf")
	panicErr(err)

	pdfs, err := constructPDFs(files)
	panicErr(err)

	readmeWriter, err := os.OpenFile(readmeFile, os.O_WRONLY, os.ModeAppend)
	panicErr(err)
	if readmeWriter != nil {
		defer func() {
			err = readmeWriter.Close()
			panicErr(err)
		}()
	}

	ghPagesWriter, err := os.OpenFile(ghPagesIndex, os.O_WRONLY, os.ModeAppend)
	panicErr(err)
	if ghPagesWriter != nil {
		defer func() {
			err = ghPagesWriter.Close()
			panicErr(err)
		}()
	}

	templateData := TemplateData{PDFs: pdfs}
	templateWriterMap := map[string]io.Writer{
		readmeTemplate: readmeWriter,
		ghPagesTemplate: ghPagesWriter,
	}
	for templateFile, writer := range templateWriterMap {
		err = templateDocument(templateFile, templateData, writer)
		panicErr(err)
	}
}
