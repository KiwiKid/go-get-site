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

func newQuestions(website Website, pageId uint) templ.Component {
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
		_, err = templBuffer.WriteString(templ.EscapeString(website.websiteURLWithPostFixAndPageId("pages/question", pageId)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"><label for=\"relevantContent&gt;&lt;\" label></label><input type=\"text\" name=\"relevantContent\" id=\"relevantContent\"><button>")
		if err != nil {
			return err
		}
		var_2 := `Get Questions`
		_, err = templBuffer.WriteString(var_2)
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

func questionResult(website Website, pageId uint, questions []Question, raw string) templ.Component {
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
		_, err = templBuffer.WriteString("<div hx-swap=\"outerHTML\">")
		if err != nil {
			return err
		}
		err = newQuestions(website, pageId).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<h1>")
		if err != nil {
			return err
		}
		var_4 := `questions: `
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		var var_5 string = strconv.Itoa(len(questions))
		_, err = templBuffer.WriteString(templ.EscapeString(var_5))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1>")
		if err != nil {
			return err
		}
		for _, q := range questions {
			_, err = templBuffer.WriteString("<h1>")
			if err != nil {
				return err
			}
			var_6 := `QuestionText:`
			_, err = templBuffer.WriteString(var_6)
			if err != nil {
				return err
			}
			var var_7 string = strconv.Itoa(len(q.QuestionText))
			_, err = templBuffer.WriteString(templ.EscapeString(var_7))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var_8 := `- `
			_, err = templBuffer.WriteString(var_8)
			if err != nil {
				return err
			}
			var var_9 string = raw
			_, err = templBuffer.WriteString(templ.EscapeString(var_9))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h1> <h1>")
			if err != nil {
				return err
			}
			var_10 := `RelevantContent:`
			_, err = templBuffer.WriteString(var_10)
			if err != nil {
				return err
			}
			var var_11 string = strconv.Itoa(len(q.RelevantContent))
			_, err = templBuffer.WriteString(templ.EscapeString(var_11))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var_12 := `- `
			_, err = templBuffer.WriteString(var_12)
			if err != nil {
				return err
			}
			var var_13 string = q.RelevantContent
			_, err = templBuffer.WriteString(templ.EscapeString(var_13))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h1>")
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

func questionsFailedResult(website Website, pageId uint, errorMsg string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_14 := templ.GetChildren(ctx)
		if var_14 == nil {
			var_14 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div hx-swap=\"outerHTML\">")
		if err != nil {
			return err
		}
		err = newQuestions(website, pageId).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<div>")
		if err != nil {
			return err
		}
		var var_15 string = errorMsg
		_, err = templBuffer.WriteString(templ.EscapeString(var_15))
		if err != nil {
			return err
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