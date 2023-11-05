// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func progressBar(progress string) templ.Component {
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
		_, err = templBuffer.WriteString("<div class=\"relative pt-1\"><div class=\"flex mb-2 items-center justify-between\"><div><span class=\"text-xs font-semibold inline-block py-1 px-2 uppercase rounded-full text-teal-600 bg-teal-200\">")
		if err != nil {
			return err
		}
		var_2 := `Task Progress`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span></div><div class=\"text-right\"><span class=\"text-xs font-semibold inline-block text-teal-600\">")
		if err != nil {
			return err
		}
		var var_3 string = progress
		_, err = templBuffer.WriteString(templ.EscapeString(var_3))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_4 := `percent`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span></div></div><div class=\"overflow-hidden h-2 mb-4 text-xs flex rounded bg-teal-200\"><div data-progress=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(string(progress)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-teal-500\"></div></div></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
