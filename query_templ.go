// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "strconv"

func queryContainer(threads []ChatThread, newThreadURL string, websites []Website) templ.Component {
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
		err = header("chat").Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<body>")
		if err != nil {
			return err
		}
		err = nav("chat").Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		err = queryHome(threads, newThreadURL, websites).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</body></html>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func queryHome(threads []ChatThread, newThreadURL string, websites []Website) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_2 := templ.GetChildren(ctx)
		if var_2 == nil {
			var_2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<main class=\"flex justify-center mt-6\" id=\"query-container\"><aside class=\"w-1/4 border-r border-gray-200 p-4\"><h1>")
		if err != nil {
			return err
		}
		var_3 := `Threads`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1>")
		if err != nil {
			return err
		}
		for _, thread := range threads {
			_, err = templBuffer.WriteString("<div class=\"flex mb-2 border-b border-gray-200 hover:bg-gray-50\"><div><a href=\"")
			if err != nil {
				return err
			}
			var var_4 templ.SafeURL = templ.URL(thread.ChatThreadURL())
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_4)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" class=\"underline\">")
			if err != nil {
				return err
			}
			var var_5 string = thread.FirstMessage
			_, err = templBuffer.WriteString(templ.EscapeString(var_5))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var_6 := `- `
			_, err = templBuffer.WriteString(var_6)
			if err != nil {
				return err
			}
			var var_7 string = strconv.Itoa(int(thread.ThreadId))
			_, err = templBuffer.WriteString(templ.EscapeString(var_7))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></div></div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</aside><!--")
		if err != nil {
			return err
		}
		var_8 := ` Main content area `
		_, err = templBuffer.WriteString(var_8)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("--><div class=\"w-3/4 p-4\" hx-swap=\"outerHTML\"><div class=\"mt-16\"><div class=\"text-gray-700 font-bold mb-4\">")
		if err != nil {
			return err
		}
		var_9 := `Search:`
		_, err = templBuffer.WriteString(var_9)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><form hx-post=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(newThreadURL))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\" hx-target=\"#query-container\" class=\"space-y-4\"><fieldset class=\"space-y-4\"><legend class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_10 := `WebsiteUrl:`
		_, err = templBuffer.WriteString(var_10)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</legend>")
		if err != nil {
			return err
		}
		for _, website := range websites {
			_, err = templBuffer.WriteString("<label class=\"block\"><input type=\"radio\" name=\"websiteId\" value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(strconv.FormatUint(uint64(website.ID), 10)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\"> ")
			if err != nil {
				return err
			}
			var var_11 string = website.BaseUrl
			_, err = templBuffer.WriteString(templ.EscapeString(var_11))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</fieldset><div class=\"space-y-4\"><label for=\"query\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_12 := `Query`
		_, err = templBuffer.WriteString(var_12)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"query\" name=\"query\" value=\"\" required class=\"p-2 border rounded w-full focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"><input type=\"submit\" value=\"Submit\" class=\"bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></div></form></div></div></main>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func query(threadId uint, websiteId uint, newChatUrl string, chats []Chat) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_13 := templ.GetChildren(ctx)
		if var_13 == nil {
			var_13 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div><h1>")
		if err != nil {
			return err
		}
		var_14 := `Chat`
		_, err = templBuffer.WriteString(var_14)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1><div><div>")
		if err != nil {
			return err
		}
		var var_15 string = strconv.Itoa(int(websiteId))
		_, err = templBuffer.WriteString(templ.EscapeString(var_15))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><form hx-post=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(newChatUrl))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\" hx-target=\"#chat-container\" class=\"space-y-4\"><input type=\"text\" id=\"threadId\" name=\"threadId\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(strconv.Itoa(int(threadId))))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hidden=\"hidden\" required><input type=\"text\" id=\"websiteId\" name=\"websiteId\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(strconv.Itoa(int(websiteId))))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hidden=\"hidden\" required><label for=\"message\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_16 := `Ask this site:`
		_, err = templBuffer.WriteString(var_16)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"message\" name=\"message\" value=\"\" required class=\"p-2 border rounded w-full focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"><input type=\"submit\" value=\"Submit\" class=\"bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></form></div><div id=\"chat-container\"></div>")
		if err != nil {
			return err
		}
		for _, item := range chats {
			_, err = templBuffer.WriteString("<div class=\"p-4 flex mb-2 border-b border-gray-200 hover:bg-gray-50\"><div></div><details class=\"p-5 border border-gray-200 rounded-lg shadow-sm\"><summary class=\"cursor-pointer text-lg font-semibold text-gray-700 hover:text-gray-900\">")
			if err != nil {
				return err
			}
			var var_17 string = item.Message
			_, err = templBuffer.WriteString(templ.EscapeString(var_17))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</summary> ")
			if err != nil {
				return err
			}
			var var_18 string = strconv.Itoa(int(item.ThreadId))
			_, err = templBuffer.WriteString(templ.EscapeString(var_18))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var var_19 string = strconv.Itoa(int(item.WebsiteId))
			_, err = templBuffer.WriteString(templ.EscapeString(var_19))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</details></div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func queryResult(pageQueryResults []PageQueryResult, websiteId uint, message string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_20 := templ.GetChildren(ctx)
		if var_20 == nil {
			var_20 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div class=\"p-4\">")
		if err != nil {
			return err
		}
		var_21 := `queryResult`
		_, err = templBuffer.WriteString(var_21)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" <div>")
		if err != nil {
			return err
		}
		var var_22 string = strconv.Itoa(int(websiteId))
		_, err = templBuffer.WriteString(templ.EscapeString(var_22))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><div class=\"container mx-auto\"><h3>")
		if err != nil {
			return err
		}
		var var_23 string = message
		_, err = templBuffer.WriteString(templ.EscapeString(var_23))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h3>")
		if err != nil {
			return err
		}
		for _, queryRes := range pageQueryResults {
			_, err = templBuffer.WriteString("<div class=\"p-4 flex mb-2 border-b border-gray-200 hover:bg-gray-50 max-h-32 m-4\"><a href=\"")
			if err != nil {
				return err
			}
			var var_24 templ.SafeURL = templ.SafeURL(queryRes.URL)
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_24)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" target=\"_blank\" rel=\"noopener noreferrer\" class=\"text-blue-500 hover:text-blue-700 underline transition duration-300 ease-in-out\">")
			if err != nil {
				return err
			}
			var var_25 string = queryRes.TidyTitle
			_, err = templBuffer.WriteString(templ.EscapeString(var_25))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a> ")
			if err != nil {
				return err
			}
			var var_26 string = queryRes.Keywords
			_, err = templBuffer.WriteString(templ.EscapeString(var_26))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" <pre>")
			if err != nil {
				return err
			}
			var var_27 string = queryRes.Content
			_, err = templBuffer.WriteString(templ.EscapeString(var_27))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</pre></div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}