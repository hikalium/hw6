<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <title>駅情報</title>
</head>
<body>
	<h1>駅情報: {{.Sta.Name}}</h1>
	<h2>隣接駅</h2>
		{{range .Sta.AdjStations}}
		<a href=".?key={{.}}">{{.}}</a>
		{{end}}
	<h2>しらべる</h2>
	<form action="." method="GET">
		
		<select name="key">
			{{range .Stations}}
			<option value="{{.Key}}">{{.Name}}</option>
			{{end}}
		</select>
		の情報を
		<button type="submit">知りたい</button>
	</form>
</body>
</html>
