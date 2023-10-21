// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func nav() templ.Component {
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
		_, err = templBuffer.WriteString("<div class=\"bg-blue-600 p-4\"><div class=\"container mx-auto\"><nav class=\"flex items-center justify-between\"><!--")
		if err != nil {
			return err
		}
		var_2 := ` Logo or Brand Name `
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("--><a href=\"/\" class=\"text-white text-2xl font-bold\">")
		if err != nil {
			return err
		}
		var_3 := `Go-Get-Site`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a><!--")
		if err != nil {
			return err
		}
		var_4 := ` Navigation Links `
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("--><div class=\"flex space-x-4\"><a href=\"/\" class=\"text-white hover:text-blue-400\">")
		if err != nil {
			return err
		}
		var_5 := `Home`
		_, err = templBuffer.WriteString(var_5)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a><!--")
		if err != nil {
			return err
		}
		var_6 := `  <a href="#" class="text-white hover:text-blue-400">About</a>
                    <a href="#" class="text-white hover:text-blue-400">Services</a>
                    <a href="#" class="text-white hover:text-blue-400">Contact</a>`
		_, err = templBuffer.WriteString(var_6)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("--></div></nav></div></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}