// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.680
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Nav() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<nav style=\"padding: 0 0rem 0rem;\"><a href=\"/main-page\">Home</a> <a href=\"/tag/main\">Tags</a> <a href=\"/profile\">Profile</a> <a href=\"/info\">Info</a></nav>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func SearchBar() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div><input class=\"form-control\" type=\"search\" name=\"search\" placeholder=\"Type tags you are intrested...\" hx-post=\"/tags\" hx-trigger=\"input changed delay:500ms, search\" hx-target=\"#search-results\" hx-indicator=\".htmx-indicator\" style=\"width: 600px;\"><table class=\"table\"><tbody id=\"search-results\"></tbody></table></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func ShowHome() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><link rel=\"stylesheet\" href=\"board/mystyle.css\"><link rel=\"stylesheet\" href=\"https://cdn.simplecss.org/simple.min.css\"><script src=\"https://unpkg.com/htmx.org@1.9.12\"></script><title>Board</title><style>\r\n            .container {\r\n        display: flex; /* Ustawia kontener na flexbox */\r\n        }\r\n        .column {\r\n            font-size:\r\n            border: 1px solid #ccc; /* Dodaje obramowanie dla kolumn */\r\n            margin: 0 10px; /* Dodaje marginesy między kolumnami */\r\n        }\r\n        .main {\r\n            flex: 80%; /* Ustawia pierwszą kolumnę na 80% szerokości */\r\n        }\r\n        .aux {\r\n            font-size: 12px;\r\n            flex: 20%; /* Ustawia drugą kolumnę na 20% szerokości */\r\n        }\r\n        </style></head><body><header><h1><img src=\"img\"> Board</h1><hr>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Nav().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</header><main>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = SearchBar().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"container\"><div class=\"main column\"><h2>Kolumna 1</h2><p>To jest zawartość kolumny 1.</p></div><div class=\"aux column\"><b style=\"text-align: center;\">Admin 👑</b><hr><form action=\"http://localhost:1000/login\" method=\"post\"><label for=\"input_email\">Email:</label><br><input type=\"text\" id=\"input_email\" name=\"input_email\"><br><label for=\"input_password\">Password:</label><br><input type=\"password\" id=\"input_password\" name=\"input_password\"><br><input type=\"submit\" value=\"Wyślij\"></form><hr><form action=\"http://localhost:1000/register\" method=\"post\"><label for=\"input_email\">Email:</label><br><input type=\"text\" id=\"input_email\" name=\"input_email\"><br><label for=\"input_password\">Password:</label><br><input type=\"password\" id=\"input_password\" name=\"input_password\"><br><input type=\"submit\" value=\"Wyślij\"></form></div></div></main></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
