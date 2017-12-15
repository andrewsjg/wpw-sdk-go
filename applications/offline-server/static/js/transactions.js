function refreshChargerStatus() {
	processJSON("http://localhost:8080/v1/transactions");
	//var uri = document.URL.match(/(^https?:\/\/[^/]+)/)[0] + "/getChargerStatus";
	//processJSON(uri);
}

function processJSON(jsonPath){
	$.getJSON(jsonPath, function(data, textStatus, jqXHR){
		var items = [];
		console.log(data);
		//o = document.createElement("div")
		rodzic = document.getElementById("tabelka");
		//o.innerText = json_array[i];
		//i.id = json_array[i].toString();
		data.forEach(function (j) {
			o2 = document.getElementById(JSON.stringify(j));
			if(o2 === null || o2 === undefined) {
				o=document.createElement("div");
				o.innerText=JSON.stringify(j);
				o.id=JSON.stringify(j);
				rodzic.appendChild(o);
			}
		});
	});
}

$.ajaxSetup({'cache':true, timeout:300});
refreshChargerStatus();
var smartChargerTimer = setInterval(refreshChargerStatus, 750 * 2);
