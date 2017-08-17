package wslog

import "html/template"

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<head>
<meta charset="utf-8">
<script>
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        console.log("open");
        if (ws != null) {
            console.log("ws open so nothing further to do");
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            console.log("Should have openened ws");
            print("OPEN");
        }
        ws.onclose = function(evt) {
            console.log("Should now close ws");
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            console.log("should message the ws");
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            console.log("should ouput error on ws");
            print("ERROR: " + evt.data);
        }
        console.log("should return from open onclick method");
        return false;
    };
    document.getElementById("justclose").onclick = function(evt) {
        console.log("justclose");
        if (!ws) {
             console.log("!ws so returning false");
             return false;
        }
        console.log("Actually closing ws");
        ws.close();
        ws.onclose();
        return false;
    };
});
</script>
<style>
.mainFont {
    font-family: 'Century Gothic', CenturyGothic, AppleGothic, sans-serif;
}
.codeFont {
    font-family: 'Lucida Sans Typewriter', 'Lucida Console', monaco, 'Bitstream Vera Sans Mono', monospace;
    color: white;
    background-color: black;
}
</style>
</head>
<body class="mainFont">
<table>
<tr><td valign="top" width="50%">
<p class="mainFont">Click "Open" to create a connection to the Rapsberry Pi Log. Click "Close" to close the connection.
<p>
<button id="open">Open</button>
<button id="justclose">Close</button>
</td><td valign="top" width="50%">
<div id="output" class="codeFont"></div>
</td></tr></table>
</body>
</html>
`))
