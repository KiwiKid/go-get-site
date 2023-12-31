// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func attributeSetResult(attributeResults []AttributeResult) templ.Component {
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
		for _, ar := range attributeResults {
			_, err = templBuffer.WriteString("<div id=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(AttributeSetResultId(ar.WebsiteID, ar.PageID, ar.AttributeSetID)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-swap-oob=\"true\" class=\"bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative\" role=\"alert\">")
			if err != nil {
				return err
			}
			var var_2 string = ar.AttributeResult
			_, err = templBuffer.WriteString(templ.EscapeString(var_2))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func attributeSetError(attributeResultsError AttributeResultsError) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_3 := templ.GetChildren(ctx)
		if var_3 == nil {
			var_3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div id=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(AttributeSetResultId(attributeResultsError.WebsiteID, attributeResultsError.PageID, attributeResultsError.AttributeSetID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-swap-oob=\"true\" class=\"bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative\" role=\"alert\">")
		if err != nil {
			return err
		}
		var var_4 string = attributeResultsError.Message
		_, err = templBuffer.WriteString(templ.EscapeString(var_4))
		if err != nil {
			return err
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

func attributeSetErrorGeneral(WebsiteID uint, AttributeSetId uint, Message string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_5 := templ.GetChildren(ctx)
		if var_5 == nil {
			var_5 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div id=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(AttributeSetResultGeneralId(WebsiteID, AttributeSetId)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-swap-oob=\"true\" class=\"bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative\" role=\"alert\">")
		if err != nil {
			return err
		}
		var var_6 string = Message
		_, err = templBuffer.WriteString(templ.EscapeString(var_6))
		if err != nil {
			return err
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
