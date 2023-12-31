// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func header(header string) templ.Component {
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
		_, err = templBuffer.WriteString("<head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>")
		if err != nil {
			return err
		}
		var_2 := `❓ Asker -  `
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		var var_3 string = header
		_, err = templBuffer.WriteString(templ.EscapeString(var_3))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</title><script src=\"https://unpkg.com/htmx.org@1.9.6\">")
		if err != nil {
			return err
		}
		var_4 := ``
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><script src=\"https://unpkg.com/htmx.org/dist/ext/sse.js\">")
		if err != nil {
			return err
		}
		var_5 := ``
		_, err = templBuffer.WriteString(var_5)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><script src=\"https://cdn.tailwindcss.com\">")
		if err != nil {
			return err
		}
		var_6 := ``
		_, err = templBuffer.WriteString(var_6)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><script>")
		if err != nil {
			return err
		}
		var_7 := `
				document.addEventListener("DOMContentLoaded", function() {
					const progressElements = document.querySelectorAll("[data-progress]");

					// Function to update element style
					const updateProgress = (element) => {
						const progressValue = element.getAttribute("data-progress");
						element.style.width = ` + "`" + `${progressValue}%` + "`" + `;
						console.log("Updated progress to " + progressValue);
					};

					// Initialize the observer
					const observer = new MutationObserver((mutations) => {
						mutations.forEach((mutation) => {
							if (mutation.type === "attributes" && mutation.attributeName === "data-progress") {
								updateProgress(mutation.target);
							}
						});
					});

					// Observer configuration
					const config = { attributes: true, childList: false, subtree: false };

					// Apply observer and initial styling
					progressElements.forEach((element) => {
						updateProgress(element); // Set initial width
						observer.observe(element, config); // Start observing for changes
					});
				});

		`
		_, err = templBuffer.WriteString(var_7)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script></head>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
