<!DOCTYPE html>
<html>
	<head>
		<title>Badgerodon IDE</title>
		{{range .Scripts}}
			<script src="{{.}}"></script>
		{{end}}
		{{range .Styles}}
			<link rel="stylesheet" type="text/css" href="{{.}}">
		{{end}}
	</head>
	<body>
		{{.Body}}
	</body>
</html>