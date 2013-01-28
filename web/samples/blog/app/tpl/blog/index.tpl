<h1>Example Blog</h1>
{{range .Posts}}
<div class="post">
	<h2>{{.Title}}</h2>
	<em>{{.Date.Format "2 Jan 2006"}} - {{.Author.Name}}</em>
	<div class="body">
		{{.Body}}
	</div>
</div>
{{end}}