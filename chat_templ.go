// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "strconv"

func chat(threadId string, websiteId string, newChatUrl string, chats []Chat) templ.Component {
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
		err = header().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<body id=\"container\">")
		if err != nil {
			return err
		}
		err = nav().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<div><h1>")
		if err != nil {
			return err
		}
		var_2 := `Chat`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1><div><div>")
		if err != nil {
			return err
		}
		var var_3 string = websiteId
		_, err = templBuffer.WriteString(templ.EscapeString(var_3))
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
		_, err = templBuffer.WriteString(templ.EscapeString(threadId))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hidden=\"hidden\" required><input type=\"text\" id=\"websiteId\" name=\"websiteId\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(websiteId))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hidden=\"hidden\" required><label for=\"message\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_4 := `Ask this site:`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"message\" name=\"message\" value=\"\" required class=\"p-2 border rounded w-full focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"><input type=\"submit\" value=\"Submit\" class=\"bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></form></div><div id=\"chat-container\"></div>")
		if err != nil {
			return err
		}
		for _, item := range chats {
			_, err = templBuffer.WriteString("<div class=\"p-4 flex mb-2 border-b border-gray-200 hover:bg-gray-50\"><div></div><details><summary class=\"cursor-pointer\">")
			if err != nil {
				return err
			}
			var var_5 string = item.Message
			_, err = templBuffer.WriteString(templ.EscapeString(var_5))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</summary> ")
			if err != nil {
				return err
			}
			var var_6 string = strconv.Itoa(int(item.ThreadId))
			_, err = templBuffer.WriteString(templ.EscapeString(var_6))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var var_7 string = strconv.Itoa(int(item.WebsiteId))
			_, err = templBuffer.WriteString(templ.EscapeString(var_7))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</details></div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div></body></html>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
