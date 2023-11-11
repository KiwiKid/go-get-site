// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "strconv"

func pages(pages []Page, website Website, count LinkCountResult, pageUrl string, prevPage string, nextPage string, addedPagesSet map[string]struct{}, progress string, viewPageSize int, processPageSize int, dripLoad bool, dripLoadCount int, dripLoadFreqMin int, dripLoadStr string, processAll bool, skipNewLinkInsert bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<html>")
		if err != nil {
			return err
		}
		err = header(website.BaseUrl).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<body class=\"bg-gray-100\" id=\"container\">")
		if err != nil {
			return err
		}
		err = nav(website.BaseUrl).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<div class=\"container mx-auto mt-12 p-4\"><a href=\"")
		if err != nil {
			return err
		}
		var var_2 templ.SafeURL = templ.URL(website.websiteURL())
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_2)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\">")
		if err != nil {
			return err
		}
		var_3 := `Refresh`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		err = editWebsite(website).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<h1 class=\"text-4xl font-bold mb-6\">")
		if err != nil {
			return err
		}
		var_4 := `PAGES`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1>")
		if err != nil {
			return err
		}
		if dripLoad {
			_, err = templBuffer.WriteString("<h1 class=\"bg-green-700 w-full px-3 py-2\">")
			if err != nil {
				return err
			}
			var_5 := `Drip Load Mode (Run every 5 minutes - run #`
			_, err = templBuffer.WriteString(var_5)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<span>")
			if err != nil {
				return err
			}
			var var_6 string = strconv.Itoa(dripLoadCount)
			_, err = templBuffer.WriteString(templ.EscapeString(var_6))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</span>")
			if err != nil {
				return err
			}
			var_7 := `)`
			_, err = templBuffer.WriteString(var_7)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h1>")
			if err != nil {
				return err
			}
		}
		err = progressBar(progress).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		if len(addedPagesSet) > 0 {
			_, err = templBuffer.WriteString("<details><summary>")
			if err != nil {
				return err
			}
			var var_8 string = strconv.Itoa(len(addedPagesSet))
			_, err = templBuffer.WriteString(templ.EscapeString(var_8))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var_9 := `new pages added - 				`
			_, err = templBuffer.WriteString(var_9)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<div>")
			if err != nil {
				return err
			}
			err = process(count, website).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div></summary><ul>")
			if err != nil {
				return err
			}
			for link := range addedPagesSet {
				_, err = templBuffer.WriteString("<li class=\"bg-green-400\">")
				if err != nil {
					return err
				}
				var var_10 string = link
				_, err = templBuffer.WriteString(templ.EscapeString(var_10))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</li>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</ul></details>")
			if err != nil {
				return err
			}
		} else {
			_, err = templBuffer.WriteString("<div>")
			if err != nil {
				return err
			}
			var_11 := `No pages added, click "Load more website content" to get started`
			_, err = templBuffer.WriteString(var_11)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("<form hx-post=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(website.websitePagesURL()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-target=\"#container\" class=\"space-y-4\"")
		if err != nil {
			return err
		}
		if dripLoad && dripLoadStr != "dripLoadStr" {
			_, err = templBuffer.WriteString(" hx-trigger=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(dripLoadStr))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\"")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("><div id=\"loading-page-processing\" style=\"display:none\">")
		if err != nil {
			return err
		}
		var_12 := `Loading...`
		_, err = templBuffer.WriteString(var_12)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><button class=\"w-full bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\">")
		if err != nil {
			return err
		}
		if dripLoad {
			var_13 := `[WILL-AUTO-CLICK `
			_, err = templBuffer.WriteString(var_13)
			if err != nil {
				return err
			}
			var var_14 string = dripLoadStr
			_, err = templBuffer.WriteString(templ.EscapeString(var_14))
			if err != nil {
				return err
			}
			var_15 := `] Load more now`
			_, err = templBuffer.WriteString(var_15)
			if err != nil {
				return err
			}
		} else {
			var_16 := `Load more website content`
			_, err = templBuffer.WriteString(var_16)
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</button><div class=\"flex items-center space-x-2\" title=\"when checked, the existing page status and content will be ignored (similar to &#39;#39;#39ag39&#39;)es&#39;) p&#39;)es&#39;)\"><input type=\"checkbox\" name=\"processAll\" id=\"processAll\" class=\"rounded border-gray-300 text-blue-600 focus:border-blue-500 focus:ring focus:ring-blue-200\"")
		if err != nil {
			return err
		}
		if processAll {
			_, err = templBuffer.WriteString(" checked=\"checked\"")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("><label for=\"processAll\" class=\"text-gray-700\">")
		if err != nil {
			return err
		}
		var_17 := `Ignore Status/Existing Content`
		_, err = templBuffer.WriteString(var_17)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label></div><div class=\"flex items-center space-x-2\"><input type=\"checkbox\" name=\"skipNewLinkInsert\" id=\"skipNewLinkInsert\" class=\"rounded border-gray-300 text-blue-600 focus:border-blue-500 focus:ring focus:ring-blue-200\"")
		if err != nil {
			return err
		}
		if skipNewLinkInsert {
			_, err = templBuffer.WriteString(" checked=\"checked\"")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("><label for=\"skipNewLinkInsert\" class=\"text-gray-700\">")
		if err != nil {
			return err
		}
		var_18 := `Dont look for new pages`
		_, err = templBuffer.WriteString(var_18)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label></div><div class=\"flex items-center space-x-2\"><label for=\"dripLoad\" class=\"text-gray-700\">")
		if err != nil {
			return err
		}
		var_19 := `Drip Load`
		_, err = templBuffer.WriteString(var_19)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"checkbox\" name=\"dripLoad\" id=\"dripLoad\" class=\"rounded border-gray-300 text-blue-600 focus:border-blue-500 focus:ring focus:ring-blue-200\"")
		if err != nil {
			return err
		}
		if dripLoad {
			_, err = templBuffer.WriteString(" checked=\"checked\"")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("><div class=\"flex items-center space-x-2\" title=\"How often a set of pages will be loaded\"><label for=\"dripLoadFreqMin\" class=\"text-gray-700\">")
		if err != nil {
			return err
		}
		var_20 := `(`
		_, err = templBuffer.WriteString(var_20)
		if err != nil {
			return err
		}
		var var_21 string = dripLoadStr
		_, err = templBuffer.WriteString(templ.EscapeString(var_21))
		if err != nil {
			return err
		}
		var_22 := `)`
		_, err = templBuffer.WriteString(var_22)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"number\" name=\"dripLoadFreqMin\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(strconv.Itoa(dripLoadFreqMin)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"></div> ")
		if err != nil {
			return err
		}
		var_23 := `[Run#:`
		_, err = templBuffer.WriteString(var_23)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" <input type=\"number\" name=\"dripLoadCount\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(strconv.Itoa(dripLoadCount)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"> ")
		if err != nil {
			return err
		}
		var_24 := `]`
		_, err = templBuffer.WriteString(var_24)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><div class=\"flex items-center space-x-2\"><label for=\"processPageSize\" class=\"text-gray-700\">")
		if err != nil {
			return err
		}
		var_25 := `Process Page Cutoff:`
		_, err = templBuffer.WriteString(var_25)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><select id=\"processPageSize\" name=\"processPageSize\" class=\"border rounded py-1 px-3 bg-white shadow-sm text-gray-700 focus:outline-none focus:ring focus:ring-blue-200\"><option value=\"1\"")
		if err != nil {
			return err
		}
		if processPageSize == 1 {
			_, err = templBuffer.WriteString(" selected")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		var_26 := `1 new page`
		_, err = templBuffer.WriteString(var_26)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"5\"")
		if err != nil {
			return err
		}
		if processPageSize == 5 {
			_, err = templBuffer.WriteString(" selected")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		var_27 := `5 new pages`
		_, err = templBuffer.WriteString(var_27)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"10\"")
		if err != nil {
			return err
		}
		if processPageSize == 10 {
			_, err = templBuffer.WriteString(" selected")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		var_28 := `10 new pages`
		_, err = templBuffer.WriteString(var_28)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"25\"")
		if err != nil {
			return err
		}
		if processPageSize == 25 {
			_, err = templBuffer.WriteString(" selected")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		var_29 := `25 new pages`
		_, err = templBuffer.WriteString(var_29)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"50\"")
		if err != nil {
			return err
		}
		if processPageSize == 50 {
			_, err = templBuffer.WriteString(" selected")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		var_30 := `50 new pages`
		_, err = templBuffer.WriteString(var_30)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"100\"")
		if err != nil {
			return err
		}
		if processPageSize == 100 {
			_, err = templBuffer.WriteString(" selected")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		var_31 := `100 new pages`
		_, err = templBuffer.WriteString(var_31)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option></select></div><div class=\"flex items-center space-x-2\"><label for=\"viewPageSize\" class=\"text-gray-700\">")
		if err != nil {
			return err
		}
		var_32 := `Per Page:`
		_, err = templBuffer.WriteString(var_32)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><select id=\"viewPageSize\" name=\"viewPageSize\" class=\"border rounded py-1 px-3 bg-white shadow-sm text-gray-700 focus:outline-none focus:ring focus:ring-blue-200\"><option value=\"1\"")
		if err != nil {
			return err
		}
		if viewPageSize == 1 {
			_, err = templBuffer.WriteString(" selected")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		var_33 := `1 page`
		_, err = templBuffer.WriteString(var_33)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"5\"")
		if err != nil {
			return err
		}
		if viewPageSize == 5 {
			_, err = templBuffer.WriteString(" selected")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		var_34 := `5 pages`
		_, err = templBuffer.WriteString(var_34)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"10\"")
		if err != nil {
			return err
		}
		if viewPageSize == 10 {
			_, err = templBuffer.WriteString(" selected")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		var_35 := `10 pages`
		_, err = templBuffer.WriteString(var_35)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"50\"")
		if err != nil {
			return err
		}
		if viewPageSize == 50 {
			_, err = templBuffer.WriteString(" selected")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		var_36 := `50 pages`
		_, err = templBuffer.WriteString(var_36)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"100\"")
		if err != nil {
			return err
		}
		if viewPageSize == 100 {
			_, err = templBuffer.WriteString(" selected")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		var_37 := `100 pages`
		_, err = templBuffer.WriteString(var_37)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"300\"")
		if err != nil {
			return err
		}
		if viewPageSize == 300 {
			_, err = templBuffer.WriteString(" selected")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		var_38 := `300 pages`
		_, err = templBuffer.WriteString(var_38)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option></select></div></form><div class=\"flex mb-2 font-bold text-sm text-gray-600 border-b-2\"><div class=\"w-1/4 p-2\">")
		if err != nil {
			return err
		}
		var_39 := `Status`
		_, err = templBuffer.WriteString(var_39)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><div class=\"w-1/4 p-2\">")
		if err != nil {
			return err
		}
		var_40 := `Title`
		_, err = templBuffer.WriteString(var_40)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><div class=\"w-1/4 p-2\">")
		if err != nil {
			return err
		}
		var_41 := `Content + Links`
		_, err = templBuffer.WriteString(var_41)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><div class=\"w-1/3 p-2\">")
		if err != nil {
			return err
		}
		var_42 := `URL`
		_, err = templBuffer.WriteString(var_42)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div></div>")
		if err != nil {
			return err
		}
		for _, item := range pages {
			_, err = templBuffer.WriteString("<div class=\"flex mb-2 border-b border-gray-200 hover:bg-gray-50\"><div class=\"w-1/4 p-2\">")
			if err != nil {
				return err
			}
			if len(item.Warning) > 0 {
				_, err = templBuffer.WriteString("<div class=\"bg-red-400 p-4 w-full\"><details><summary>")
				if err != nil {
					return err
				}
				var_43 := `Error`
				_, err = templBuffer.WriteString(var_43)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</summary>")
				if err != nil {
					return err
				}
				var var_44 string = item.Warning
				_, err = templBuffer.WriteString(templ.EscapeString(var_44))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</details></div>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("<details>")
			if err != nil {
				return err
			}
			if item.ToProcess() {
				_, err = templBuffer.WriteString("<summary>")
				if err != nil {
					return err
				}
				var_45 := `TO_PROCESS`
				_, err = templBuffer.WriteString(var_45)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</summary>")
				if err != nil {
					return err
				}
			} else {
				_, err = templBuffer.WriteString("<summary>")
				if err != nil {
					return err
				}
				var_46 := `DONE`
				_, err = templBuffer.WriteString(var_46)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</summary>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var var_47 string = item.DateCreated.Local().Format(format)
			_, err = templBuffer.WriteString(templ.EscapeString(var_47))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var var_48 string = item.DateUpdated.Local().Format(format)
			_, err = templBuffer.WriteString(templ.EscapeString(var_48))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" <div>")
			if err != nil {
				return err
			}
			var var_49 string = item.PageStatus()
			_, err = templBuffer.WriteString(templ.EscapeString(var_49))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div></details></div><div class=\"w-1/4 p-2\">")
			if err != nil {
				return err
			}
			var var_50 string = item.Title
			_, err = templBuffer.WriteString(templ.EscapeString(var_50))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div><div class=\"w-1/4 p-2\"><details><summary class=\"cursor-pointer\">")
			if err != nil {
				return err
			}
			var_51 := `Content: (`
			_, err = templBuffer.WriteString(var_51)
			if err != nil {
				return err
			}
			var var_52 string = strconv.Itoa(len(item.Content))
			_, err = templBuffer.WriteString(templ.EscapeString(var_52))
			if err != nil {
				return err
			}
			var_53 := `) Links: (`
			_, err = templBuffer.WriteString(var_53)
			if err != nil {
				return err
			}
			var var_54 string = strconv.Itoa(len(item.Links))
			_, err = templBuffer.WriteString(templ.EscapeString(var_54))
			if err != nil {
				return err
			}
			var_55 := `)`
			_, err = templBuffer.WriteString(var_55)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</summary><h2>")
			if err != nil {
				return err
			}
			var_56 := `Content`
			_, err = templBuffer.WriteString(var_56)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h2><pre>")
			if err != nil {
				return err
			}
			var var_57 string = item.Content
			_, err = templBuffer.WriteString(templ.EscapeString(var_57))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</pre><br><details><summary><h2>")
			if err != nil {
				return err
			}
			var_58 := `Links (`
			_, err = templBuffer.WriteString(var_58)
			if err != nil {
				return err
			}
			var var_59 string = strconv.Itoa(len(item.Links))
			_, err = templBuffer.WriteString(templ.EscapeString(var_59))
			if err != nil {
				return err
			}
			var_60 := `)`
			_, err = templBuffer.WriteString(var_60)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h2></summary><ul>")
			if err != nil {
				return err
			}
			for _, link := range item.Links {
				_, err = templBuffer.WriteString("<li>")
				if err != nil {
					return err
				}
				var var_61 string = link
				_, err = templBuffer.WriteString(templ.EscapeString(var_61))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</li>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</ul></details></details></div><div class=\"w-1/3 p-2\"><a target=\"_blank\" href=\"")
			if err != nil {
				return err
			}
			var var_62 templ.SafeURL = templ.SafeURL(item.URL)
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_62)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var var_63 string = item.URL
			_, err = templBuffer.WriteString(templ.EscapeString(var_63))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></div></div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		if nextPage != "" {
			_, err = templBuffer.WriteString("<div id=\"next-space\"><button class=\"btn\" hx-get=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(nextPage))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-taget=\"#next-space\">")
			if err != nil {
				return err
			}
			var_64 := `Load more`
			_, err = templBuffer.WriteString(var_64)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button></div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</body>")
		if err != nil {
			return err
		}
		err = websiteDelete(website).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</html>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
