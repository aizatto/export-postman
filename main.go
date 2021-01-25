package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

type Postman struct {
	Info struct {
		PostmanID   string `json:"_postman_id"`
		Name        string
		Description string
		Schema      string
	}
	Item []PostmanItem
}

type PostmanItem struct {
	Name     string
	Request  *PostmanRequest
	Response []PostmanResponse
	Item     []PostmanItem
}

type PostmanRequest struct {
	Method string
	Header []struct {
		Key   string
		Name  string
		Value string
		Text  string
	}
	Body struct {
		Mode string
		Raw  string
	}
	URL struct {
		Raw  string
		Path []string
	}
}

type PostmanResponse struct {
	Name            string
	OriginalRequest PostmanRequest
	Status          string
	Code            int
	Body            string
	Header          []struct {
		Key  string
		Name string
	}
}

func main() {
	if len(os.Args) == 1 {
		color.Red("No files passed")
		os.Exit(1)
	}

	for _, arg := range os.Args[1:len(os.Args)] {
		err := processFile(arg)
		if err != nil {
			color.Red(err.Error())
		}
	}
}

func processFile(file string) error {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("File is a directory: %s", file)
	}

	body, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	var postman Postman
	err = json.Unmarshal(body, &postman)
	if err != nil {
		return err
	}

	output := ""
	for _, item := range postman.Item {
		output += processItem(item, 1)
		output += "\n"
	}

	base := strings.TrimSuffix(file, filepath.Ext(file))
	err = ioutil.WriteFile(
		fmt.Sprintf("%s.md", base),
		[]byte(output),
		0744,
	)
	if err != nil {
		return err
	}

	return nil
}

func processItem(item PostmanItem, depth int) string {
	output := ""

	for i := 0; i < depth; i++ {
		output += "#"
	}

	path := ""
	if item.Request != nil &&
		len(item.Request.URL.Path) > 0 {
		path = "/" + strings.Join(item.Request.URL.Path, "/")
	}

	output += fmt.Sprintf(
		" %s %s",
		item.Name,
		path,
	)

	output += "\n"

	for _, item := range item.Item {
		output += "\n"
		output += processItem(item, depth+1)
	}

	if item.Request == nil {
		return output
	}

	output += "\n"

	request := item.Request

	data := [][]string{
		{
			"Url API",
			path,
		},
		{
			"Method",
			request.Method,
		},
		{
			"Format Body",
			"Dikirim lewat format json",
		},
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"", ""})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

	output += tableString.String()

	if len(item.Response) == 0 {
		return output
	}
	response := item.Response[0]

	output += "\n"

	body := prettyprintJson(response.OriginalRequest.Body.Raw)
	if len(body) > 0 {
		output += "Body:\n\n"
		output += fmt.Sprintf(
			"```\n%s\n```\n\n",
			body,
		)
	} else {
		output += "Body: _Empty_\n\n"
	}

	r := prettyprintJson(response.Body)
	if len(r) > 0 {
		output += "Response:\n\n"
		output += fmt.Sprintf(
			"```\n%s\n```\n\n",
			r,
		)
	}

	return output
}

func prettyprintJson(input string) string {
	if len(input) == 0 || len(input) == 1 {
		return ""
	}

	var fields map[string]interface{}
	err := json.Unmarshal([]byte(input), &fields)
	if err != nil {
		color.Red("1")
		color.Red("%d", len(input))
		color.Red(err.Error())
		return ""
	}

	cleanupJson(&fields)

	body, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		color.Red("2")
		color.Red(err.Error())
		return ""
	}

	strbody := string(body)
	// return strings.Join(strings.Split(strbody, "\n"), "\n\n")
	return strbody
}

func cleanupJson(fieldsptr *map[string]interface{}) {
	fields := *fieldsptr

	for key, value := range fields {
		switch value.(type) {
		case string:
			if strings.HasPrefix(value.(string), "data:image/png;base64,") {
				fields[key] = "data:image/png;base64,..."
			} else if strings.HasPrefix(value.(string), "data:image/jpeg;base64,") {
				fields[key] = "data:image/png;base64,..."
			}
		case map[string]interface{}:
			v2 := value.(map[string]interface{})
			cleanupJson(&v2)
		}
	}
}
