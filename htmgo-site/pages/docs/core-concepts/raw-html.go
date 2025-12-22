package core_concepts

import "htmgo-site/ui"
import "github.com/maddalax/htmgo/framework/h"
import . "htmgo-site/pages/docs"

func RawHtml(ctx *h.RequestContext) *h.Page {
	return DocPage(
		ctx,
		h.Div(
			h.Class("flex flex-col gap-3"),
			Title("Raw HTML"),
			Text(`
				In some cases, you may want to render raw html instead of using htmgo's functions.
		    This can be done by using the following methods:
			`),
			Text("Rendering raw html:"),
			ui.GoCodeSnippetSingleLine(RawHtmlExample),
			Text("Rendering with formatting:"),
			ui.GoCodeSnippetSingleLine(RawHtmlExample2),
			Text("Rendering a script:"),
			ui.GoCodeSnippetSingleLine(RawHtmlExample3),
			Text(`
        Important: Be careful when using these methods, these methods do not escape the HTML content
        and should never be used with user input unless you have sanitized the input.
			`),
			h.P(
				h.Text("Sanitizing input can be done using "),
				Link("html.EscapeString", "https://pkg.go.dev/html#EscapeString"),
				h.Text(" or by using "),
				Link("bluemonday", "https://github.com/microcosm-cc/bluemonday"),
				h.Text(" for more control over sanitization."),
			),
			NextStep(
				"mt-4",
				PrevBlock("Tags and Attributes", DocPath("/core-concepts/tags-and-attributes")),
				NextBlock("Conditionals", DocPath("/control/if-else")),
			),
		),
	)
}

const RawHtmlExample = `h.UnsafeRaw("<div>Raw HTML</div>")`
const RawHtmlExample2 = `h.UnsafeRawF("<div>%s</div>", "Raw HTML")`
const RawHtmlExample3 = `h.UnsafeRawScript("alert('Hello World')")`
