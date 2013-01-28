<!DOCTYPE html>
<html>
	<head>
		<title>Some Blog</title>
		<link href='http://fonts.googleapis.com/css?family=Life+Savers' rel='stylesheet' type='text/css'>
		{{range .Styles}}
		<link rel="stylesheet" href="{{.}}" type="text/css" media="all">
		{{end}}		
		{{range .Scripts}}
		<script src="{{.}}"></script>
		{{end}}
	</head>
	<body>
		<div id="outer">
			<div id="inner">
{{.Body}}
			</div>
		</div>
	</body>
</html>