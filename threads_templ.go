// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "strconv"

func threads(threads []ChatThread, newThreadURL string, websites []Website) templ.Component {
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
		err = header("chat2").Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<body id=\"container\">")
		if err != nil {
			return err
		}
		err = nav("chat2").Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<main class=\"flex justify-center mt-6\"><aside class=\"w-1/4 border-r border-gray-200 p-4\"><h1>")
		if err != nil {
			return err
		}
		var_2 := `Threads`
		_, err = templBuffer.WriteString(var_2)
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
			var var_3 templ.SafeURL = templ.URL(thread.ChatThreadURL())
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_3)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" class=\"underline\">")
			if err != nil {
				return err
			}
			var var_4 string = thread.FirstMessage
			_, err = templBuffer.WriteString(templ.EscapeString(var_4))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var_5 := `- `
			_, err = templBuffer.WriteString(var_5)
			if err != nil {
				return err
			}
			var var_6 string = strconv.Itoa(int(thread.ThreadId))
			_, err = templBuffer.WriteString(templ.EscapeString(var_6))
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
		var_7 := ` Main content area `
		_, err = templBuffer.WriteString(var_7)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("--><div class=\"w-3/4 p-4\" hx-swap=\"outerHTML\"><div class=\"mt-16\"><div class=\"text-gray-700 font-bold mb-4\">")
		if err != nil {
			return err
		}
		var_8 := `Search:`
		_, err = templBuffer.WriteString(var_8)
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
		_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\" hx-target=\"#container\" class=\"space-y-4\"><fieldset class=\"space-y-4\"><legend class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_9 := `WebsiteUrl:`
		_, err = templBuffer.WriteString(var_9)
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
			var var_10 string = website.BaseUrl
			_, err = templBuffer.WriteString(templ.EscapeString(var_10))
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
		var_11 := `Query`
		_, err = templBuffer.WriteString(var_11)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"query\" name=\"query\" value=\"\" required class=\"p-2 border rounded w-full focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"><input type=\"submit\" value=\"Submit\" class=\"bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></div></form></div></div></main></body></html>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
