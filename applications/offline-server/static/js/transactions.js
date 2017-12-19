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
				o = document.createElement("div");
				o.id = JSON.stringify(j);

				s = document.createElement("span");
				s.classList.add("transactionDateAndTime");
				s.name = "transactionDateAndTime";
				s.innerText = j.transactionDateAndTime;
				o.appendChild(s);

				s = document.createElement("span");
				s.classList.add("card");
				s.name = "card";
				s.innerText = j.maskedCard;
				o.appendChild(s);

				s = document.createElement("span");
				s.classList.add("amount");
				s.name = "amount";
				s.innerText = (j.amount/100).toFixed(2);
				o.appendChild(s);


				s = document.createElement("span");
				s.classList.add("currencyCode");
				s.name = "currencyCode";
				s.innerText = j.currencyCode;
				o.appendChild(s);

				s = document.createElement("span");
				s.classList.add("orderDescription");
				s.name = "orderDescription";
				s.innerText = j.orderDescription;
				o.appendChild(s);

				//o.innerText=JSON.stringify(j);
				rodzic.appendChild(o);
			}
		});
	});
}

$.ajaxSetup({'cache':true, timeout:300});
refreshChargerStatus();
var smartChargerTimer = setInterval(refreshChargerStatus, 750 * 2);
