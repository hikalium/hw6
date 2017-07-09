<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <title>のりかえ案内</title>
</head>
<body>
    <h1>のりかえ案内</h1>
	{{if .StationFrom}}
	<h2>出発駅: {{.StationFrom.Name}}</h2>
	{{end}}
	{{if .StationTo}}
	<h2>到着駅: {{.StationTo.Name}}</h2>
	{{end}}
	{{if .Result}}
	<h2>ルート</h2>
		{{range .Result}}
		<p><a href="../stainfo/?key={{.}}">{{.}}</a></p>
		{{end}}
	{{end}}
	<h2>しらべる</h2>
	<form action="." method="GET">
		
		<select name="stationFrom">
			{{range .Stations}}
			<option value="{{.Key}}">{{.Name}}</option>
			{{end}}
		</select>
		から
		<select name="stationTo">
			{{range .Stations}}
			<option value="{{.Key}}">{{.Name}}</option>
			{{end}}
		</select>
		まで行く方法を
		<button type="submit">知りたい</button>
	</form>
</body>
</html>
