package service

import "net/http"
import 	"html/template"

func home(w http.ResponseWriter, r *http.Request) {

	homeTemplate.Execute(w, "ws://"+r.Host+"/logs?namespace=demo&name=hello-world")

}

func ADDWS() {
	http.HandleFunc("/", home)
}


var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <script>
        window.addEventListener("load", function(evt) {
            var oDiv = document.getElementById('float_banner');
            oDiv.style.position = 'fixed';
            oDiv.style.top = '20px';
            oDiv.style.left = '20px';
            oDiv.style.backgroundColor = "#0000FF";
            oDiv.style.color="#F8F8FF";
            var float = function(message) {
                var d = document.createElement("div");
                d.innerHTML = message;
                oDiv.appendChild(d);
            };


            var output = document.getElementById("output");
            var logtext = document.getElementById("logtext");
            logtext.style.backgroundColor = "#000000";
            logtext.style.color="#F8F8FF";
            //var ws;

            var print = function(message) {
                var d = document.createElement("div");
                d.innerHTML = message;
                output.appendChild(d);
            };

//            var ws = new WebSocket("ws://localhost:8080/echo");
            ws = new WebSocket("{{.}}");
            ws.onopen = function(evt) {
                console.log("Connection open ...");
                ws.send("Hello WebSockets!");
            };
            var i = 0;
            ws.onmessage = function(evt) {
                if (i=0) {
                    float(evt.data);
                    i=1
                }


                print(evt.data)
                //ws.close();
            };

            ws.onclose = function(evt) {
                console.log("Connection closed.");
            };

        });
    </script>
</head>
<div id="logtext">
<div id="float_banner"><p></p></div>
<div id="output"></div>
</div>
</html>
`))
