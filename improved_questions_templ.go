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

func newImprovedQuestion(websiteId uint, pageId uint, pageBlockID uint, q Question) templ.Component {
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
		_, err = templBuffer.WriteString("<div><form hx-post=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(questionImprovement(websiteId, pageId, pageBlockID, q.ID, "new")))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-swap=\"#page-block-container\" class=\"space-y-4\"><div><label for=\"relevantContent\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_2 := `Content`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><textarea type=\"text\" rows=\"5\" cols=\"80\" name=\"rawQuestionText\" id=\"rawQuestionText\">")
		if err != nil {
			return err
		}
		var var_3 string = splitQuestion(q.QuestionText, true)
		_, err = templBuffer.WriteString(templ.EscapeString(var_3))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</textarea><textarea type=\"text\" rows=\"5\" cols=\"80\" name=\"rawQuestionAnswer\" id=\"rawQuestionAnswer\">")
		if err != nil {
			return err
		}
		var var_4 string = splitQuestion(q.QuestionText, false)
		_, err = templBuffer.WriteString(templ.EscapeString(var_4))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</textarea></div><input type=\"number\" name=\"pageId\" id=\"pageId\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(strconv.Itoa(int(pageId))))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"hidden\"><button type=\"submit\" class=\"inline-flex justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500\">")
		if err != nil {
			return err
		}
		var_5 := `Improve Question`
		_, err = templBuffer.WriteString(var_5)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button></form></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func improvedQuestions(websiteId uint, pageId uint, pageBlockID uint, iqs []ImprovedQuestion) templ.Component {
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
		_, err = templBuffer.WriteString("<div>")
		if err != nil {
			return err
		}
		for _, i := range iqs {
			_, err = templBuffer.WriteString("<div>")
			if err != nil {
				return err
			}
			var var_7 string = i.QuestionText
			_, err = templBuffer.WriteString(templ.EscapeString(var_7))
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
