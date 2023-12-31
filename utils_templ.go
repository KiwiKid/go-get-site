// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func spinner() templ.Component {
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
		_, err = templBuffer.WriteString("<style type=\"text/css\">")
		if err != nil {
			return err
		}
		var_2 := `
        @keyframes spin {
            0% { 
                transform: rotate(0deg);
             }
            100% { 
                transform: rotate(360deg);
                 }
        }

        .spinner {
            animation: spin 1s linear infinite;
        }
    `
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</style><svg class=\"spinner w-16 h-16 text-green-600\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\"><circle class=\"opacity-25\" cx=\"12\" cy=\"12\" r=\"10\" stroke-width=\"4\"></circle><path class=\"opacity-75\" d=\"M12 2v10\" stroke-width=\"4\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
