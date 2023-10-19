// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func task(siteUrl string) templ.Component {
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
		_, err = templBuffer.WriteString("<html lang=\"en\">")
		if err != nil {
			return err
		}
		err = header().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<body><div hx-get=\"/site/:siteUrl\" hx-trigger=\"load\">")
		if err != nil {
			return err
		}
		var var_2 string = siteUrl
		_, err = templBuffer.WriteString(templ.EscapeString(var_2))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" <img alt=\"Result loading...\" class=\"htmx-indicator\" width=\"150\" src=\"/img/bars.svg\"></div><!--")
		if err != nil {
			return err
		}
		var_3 := `<div hx-ext="sse" sse-connect="/progress" sse-swap="message">
        Contents of this box will be updated in real time
        with every SSE message received from the chatroom.
    </div>`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("--></body></html>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
