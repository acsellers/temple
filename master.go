package main

const Master = `<!DOCTYPE html>
<html>
<head>
	<title>{{ or .Title "Template Generation Demo"}}</title>
</head>
<body>
	<div id="topbar">
		[[ call .Features "topbar" ]]
	</div>
	<div id="navbar">
		[[ call .Features "navbar" ]]
	</div>
	{{ .Content }}
	<div id="footer">
		[[ call .Features "footer" ]]
	</div>
</body>
</html>`

var (
	Tops = map[string]string{
		"topbar.one":   `<div class="topbar-item">One</div>`,
		"topbar.two":   `<div class="topbar-item">Two</div>`,
		"topbar.three": `<div class="topbar-item">Three</div>`,
	}
	Navs = map[string]string{
		"navbar.one":   `<div class="navbar-item">One</div>`,
		"navbar.two":   `<div class="navbar-item">Two</div>`,
		"navbar.three": `<div class="navbar-item">Three</div>`,
	}
	Foots = map[string]string{
		"footer.one":   `<div class="footer-item">One</div>`,
		"footer.two":   `<div class="footer-item">Two</div>`,
		"footer.three": `<div class="footer-item">Three</div>`,
	}
)
