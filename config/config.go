package config

var (
	Pwd = "1"
)

var (
	IndexPage = `
<h1></h1>
<form method="post" action="/login">
    <input type="password" id="pwd" name="pwd">
    <button type="submit">login</button>
</form>
`
	HomePage = `
<p></p>
<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`
)
