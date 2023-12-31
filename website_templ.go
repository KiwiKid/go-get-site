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

func websiteList(websites []Website) templ.Component {
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
		_, err = templBuffer.WriteString("<div><ul class=\"grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4\">")
		if err != nil {
			return err
		}
		for _, webUrl := range websites {
			_, err = templBuffer.WriteString("<li class=\"mb-2\"><div class=\"flex flex-col bg-gray-100 rounded-lg shadow hover:shadow-md transition-shadow\"><a class=\"bg-blue-500 text-white py-2 px-4 rounded-t-lg hover:bg-blue-600 focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300 block\" href=\"")
			if err != nil {
				return err
			}
			var var_2 templ.SafeURL = templ.URL(webUrl.websitePagesURL())
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_2)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var var_3 string = webUrl.BaseUrl
			_, err = templBuffer.WriteString(templ.EscapeString(var_3))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a><div class=\"p-2\">")
			if err != nil {
				return err
			}
			err = editWebsite(webUrl).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div></div></li>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</ul></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func editWebsite(website Website) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_4 := templ.GetChildren(ctx)
		if var_4 == nil {
			var_4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div><div class=\"mt-8\"><details><summary><div class=\"text-gray-700 font-bold mb-4\">")
		if err != nil {
			return err
		}
		var_5 := `Edit `
		_, err = templBuffer.WriteString(var_5)
		if err != nil {
			return err
		}
		var var_6 string = website.BaseUrl
		_, err = templBuffer.WriteString(templ.EscapeString(var_6))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div></summary><form hx-put=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(website.websiteURL()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\" hx-target=\"#container\" class=\"space-y-4 border border-gray-300 rounded p-4\"><label for=\"websiteUrl\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_7 := `Website Url:`
		_, err = templBuffer.WriteString(var_7)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"url\" id=\"websiteUrl\" name=\"websiteUrl\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(website.BaseUrl))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" pattern=\"^https?://.*\" required class=\"p-2 border rounded w-full focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"><input type=\"text\" id=\"websiteId\" name=\"websiteId\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(strconv.Itoa(int(website.ID))))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hidden=\"hidden\" required><div class=\"w-1/2\" title=\"(optional) click to make before saving each page\"><label for=\"preLoadPageClickSelector\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_8 := `PreLoadPageClickSelector:`
		_, err = templBuffer.WriteString(var_8)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"preLoadPageClickSelector\" name=\"preLoadPageClickSelector\" placeholder=\"input[id=&#39;expand-all-the-things&#39;]\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(website.PreLoadPageClickSelector))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" title=\"an (optional) CSS selector for a click before the content is loaded for each page \" class=\"w-full p-2 border rounded focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></div><div class=\"w-1/2\" title=\"(optional) replace text in the page title\"><label for=\"titleReplace\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_9 := `Title Replace:`
		_, err = templBuffer.WriteString(var_9)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"titleReplace\" name=\"titleReplace\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(website.TitleReplace))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" placeholder=\"- title on every page\" class=\"p-2 border rounded w-full focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></div><details><summary>")
		if err != nil {
			return err
		}
		var_10 := `Login options`
		_, err = templBuffer.WriteString(var_10)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</summary><div class=\"space-y-4\"><div class=\"flex space-x-4\"><div class=\"w-1/2\"><label for=\"customQueryParam\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_11 := `customQueryParam:`
		_, err = templBuffer.WriteString(var_11)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"customQueryParam\" name=\"customQueryParam\" placeholder=\"auth=DontWorryAboutIt\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(website.CustomQueryParam))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"w-full p-2 border rounded focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></div><div class=\"w-1/2\" title=\"the successIndicatorSelector is checked for existance before the content of the page is read to avoid saving loading states\"><label for=\"successIndicatorSelector\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_12 := `SuccessIndicatorSelector:`
		_, err = templBuffer.WriteString(var_12)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"successIndicatorSelector\" name=\"successIndicatorSelector\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(website.SuccessIndicatorSelector))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"w-full p-2 border rounded focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></div><div class=\"w-1/2\" title=\"the (optiona ) can trigger a visit to this url as an entrypoint for the site\"><label for=\"startUrl\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_13 := `Auth Url:`
		_, err = templBuffer.WriteString(var_13)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"url\" id=\"startUrl\" name=\"startUrl\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(website.StartUrl))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" pattern=\"^https?://.*\" placeholder=\"https://example.com?pre-authed-link=true\" class=\"p-2 border rounded w-full focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></div></div><br><h3>")
		if err != nil {
			return err
		}
		var_14 := `OR / AND`
		_, err = templBuffer.WriteString(var_14)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h3><div class=\"flex space-x-4\">")
		if err != nil {
			return err
		}
		if website.LoginName == "dont-login" {
			_, err = templBuffer.WriteString("<h2>")
			if err != nil {
				return err
			}
			var_15 := `(No login steps)`
			_, err = templBuffer.WriteString(var_15)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h2>")
			if err != nil {
				return err
			}
		} else {
			if website.LoginName != "" {
				_, err = templBuffer.WriteString("<h2>")
				if err != nil {
					return err
				}
				var_16 := `Loggin in with `
				_, err = templBuffer.WriteString(var_16)
				if err != nil {
					return err
				}
				var var_17 string = website.LoginName
				_, err = templBuffer.WriteString(templ.EscapeString(var_17))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</h2>")
				if err != nil {
					return err
				}
			} else {
				_, err = templBuffer.WriteString("<h2>")
				if err != nil {
					return err
				}
				var_18 := `(No login steps)`
				_, err = templBuffer.WriteString(var_18)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</h2>")
				if err != nil {
					return err
				}
			}
		}
		_, err = templBuffer.WriteString("<div class=\"w-1/2\"><label for=\"loginName\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_19 := `loginName:`
		_, err = templBuffer.WriteString(var_19)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"loginName\" name=\"loginName\" placeholder=\"bobby_tables\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(website.LoginName))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"w-full p-2 border rounded focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></div><div class=\"w-1/2\"><label for=\"loginPass\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_20 := `loginPass:`
		_, err = templBuffer.WriteString(var_20)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"loginPass\" name=\"loginPass\" placeholder=\"Password123\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(website.LoginPass))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"w-full p-2 border rounded focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></div></div><div class=\"flex space-x-4\"><div class=\"w-1/2\"><label for=\"loginNameSelector\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_21 := `loginNameSelector:`
		_, err = templBuffer.WriteString(var_21)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"loginNameSelector\" placeholder=\"input[id=&#39;username&#39;]\" name=\"loginNameSelector\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(website.LoginNameSelector))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" title=\"a CSS selector for the Login Name input box\" class=\"w-full p-2 border rounded focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></div><div class=\"w-1/2\"><label for=\"loginPassSelector\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_22 := `loginPassSelector:`
		_, err = templBuffer.WriteString(var_22)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"loginPassSelector\" name=\"loginPassSelector\" placeholder=\"input[id=&#39;password&#39;]\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(website.LoginPassSelector))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" title=\"a CSS selector for the Password input box\" class=\"w-full p-2 border rounded focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></div></div><label for=\"submitButtonSelector\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_23 := `SubmitButtonSelector:`
		_, err = templBuffer.WriteString(var_23)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"submitButtonSelector\" name=\"submitButtonSelector\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(website.SubmitButtonSelector))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" title=\"a CSS selector for the submit button\" class=\"w-full p-2 border rounded focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></div></details><div class=\"mt-6\"><input type=\"submit\" value=\"Update website\" class=\"w-full bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></div>")
		if err != nil {
			return err
		}
		err = login(website).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		err = websiteDelete(website).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</form></details></div></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func newWebsite() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_24 := templ.GetChildren(ctx)
		if var_24 == nil {
			var_24 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div><div class=\"mt-8\"><div class=\"text-gray-700 font-bold mb-4\">")
		if err != nil {
			return err
		}
		var_25 := `A new site url:`
		_, err = templBuffer.WriteString(var_25)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><form hx-post=\"/\" hx-swap=\"outerHTML\" hx-target=\"#container\" class=\"space-y-4\"><label for=\"websiteUrl\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_26 := `Website Url:`
		_, err = templBuffer.WriteString(var_26)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"url\" id=\"websiteUrl\" name=\"websiteUrl\" pattern=\"^https?://.*\" required class=\"p-2 border rounded w-full focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"><label for=\"startUrl\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_27 := `Start Url:`
		_, err = templBuffer.WriteString(var_27)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"url\" id=\"startUrl\" name=\"startUrl\" pattern=\"^https?://.*\" class=\"p-2 border rounded w-full focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"><label for=\"successIndicatorSelector\" class=\"block text-sm font-medium text-gray-600\">")
		if err != nil {
			return err
		}
		var_28 := `SuccessIndicatorSelector:`
		_, err = templBuffer.WriteString(var_28)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"successIndicatorSelector\" default=\"body\" name=\"successIndicatorSelector\" required class=\"w-full p-2 border rounded focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"><input type=\"submit\" value=\"Submit\" class=\"bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300\"></form></div></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func websiteDelete(website Website) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_29 := templ.GetChildren(ctx)
		if var_29 == nil {
			var_29 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div><details><summary>")
		if err != nil {
			return err
		}
		var_30 := `Delete & Reset`
		_, err = templBuffer.WriteString(var_30)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</summary><button hx-target=\"#container\" hx-loading=\"#loading\" hx-delete=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(website.websiteURLWithPostFix("pages")))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"bg-red-500 text-white py-2 px-4 rounded hover:bg-red-600 focus:ring focus:ring-opacity-50 focus:ring-red-300 focus:border-red-300\">")
		if err != nil {
			return err
		}
		var_31 := `reset pages`
		_, err = templBuffer.WriteString(var_31)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button><button hx-target=\"#container\" hx-loading=\"#loading\" hx-confirm=\"This will delete all info related to this website\" hx-delete=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(website.websiteURLWithPostFix("all")))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"bg-red-500 text-white py-2 px-4 rounded hover:bg-red-600 focus:ring focus:ring-opacity-50 focus:ring-red-300 focus:border-red-300\">")
		if err != nil {
			return err
		}
		var_32 := `DELETE`
		_, err = templBuffer.WriteString(var_32)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button><div id=\"loading\" style=\"display:none\">")
		if err != nil {
			return err
		}
		var_33 := `Loading...`
		_, err = templBuffer.WriteString(var_33)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div></details></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
