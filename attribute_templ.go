// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"strconv"
)

func newAttribute(attributeModels []AttributeModel) templ.Component {
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
		_, err = templBuffer.WriteString("<div class=\"p-4 shadow rounded-lg bg-white border border-gray-400 rounded-md w-1/4\"><form hx-post=\"/attribute\" hx-swap=\"#attribute-container\" class=\"space-y-4 \"><div><label for=\"attributeSeedQuery\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_2 := `Attribute Seed Query`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><textarea name=\"attributeSeedQuery\" id=\"attributeSeedQuery\" rows=\"3\" class=\"mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm\" placeholder=\"Enter attribute seed query\"></textarea></div><div><label for=\"attributeModelID\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_3 := `Attribute Model ID`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><select name=\"attributeModelID\" id=\"attributeModelID\" class=\"mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm\">")
		if err != nil {
			return err
		}
		for _, model := range attributeModels {
			_, err = templBuffer.WriteString("<option value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(strconv.Itoa(int(model.ID))))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var var_4 string = model.Name
			_, err = templBuffer.WriteString(templ.EscapeString(var_4))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</select></div><button type=\"submit\" class=\"inline-flex justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500\">")
		if err != nil {
			return err
		}
		var_5 := `Create Attribute`
		_, err = templBuffer.WriteString(var_5)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" <span class=\"htmx-indicator\">")
		if err != nil {
			return err
		}
		err = spinner().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span></button></form></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func attributes(attributes []Attribute, message string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_6 := templ.GetChildren(ctx)
		if var_6 == nil {
			var_6 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div class=\" border border-gray-400 rounded-md w-1/4\">")
		if err != nil {
			return err
		}
		var var_7 string = message
		_, err = templBuffer.WriteString(templ.EscapeString(var_7))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		for _, attr := range attributes {
			_, err = templBuffer.WriteString("<div>")
			if err != nil {
				return err
			}
			var var_8 string = strconv.Itoa(int(attr.ID))
			_, err = templBuffer.WriteString(templ.EscapeString(var_8))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div> <div>")
			if err != nil {
				return err
			}
			var var_9 string = attr.AISeedQuery
			_, err = templBuffer.WriteString(templ.EscapeString(var_9))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div> <div>")
			if err != nil {
				return err
			}
			var var_10 string = attr.AITask
			_, err = templBuffer.WriteString(templ.EscapeString(var_10))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div> <div>")
			if err != nil {
				return err
			}
			var var_11 string = strconv.Itoa(int(attr.AIArgs.MinLength))
			_, err = templBuffer.WriteString(templ.EscapeString(var_11))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div> <div>")
			if err != nil {
				return err
			}
			var var_12 string = strconv.Itoa(int(attr.AIArgs.MaxLength))
			_, err = templBuffer.WriteString(templ.EscapeString(var_12))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div> <div>")
			if err != nil {
				return err
			}
			var var_13 string = strconv.Itoa(int(attr.AttributeModelID))
			_, err = templBuffer.WriteString(templ.EscapeString(var_13))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div> <div>")
			if err != nil {
				return err
			}
			var var_14 string = strconv.Itoa(int(attr.AttributeSetID))
			_, err = templBuffer.WriteString(templ.EscapeString(var_14))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
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
