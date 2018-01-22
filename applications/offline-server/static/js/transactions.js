function refreshChargerStatus() {
	var uri = document.URL.match(/(^https?:\/\/[^/]+)/)[0] + "/v1/transactions";
	processJSON(uri);
	//var uri = document.URL.match(/(^https?:\/\/[^/]+)/)[0] + "/getChargerStatus";
	//processJSON(uri);
}

function processJSON(jsonPath){
	$.getJSON(jsonPath, function(data, textStatus, jqXHR){
		var items = [];
		//console.log(data);
		//o = document.createElement("div")
		rodzic = document.getElementById("tabelkaTBody");
		//o.innerText = json_array[i];
		//i.id = json_array[i].toString();
		if(data.length == 0){
			var a = rodzic.getElementsByClassName("paymentRecord");
			if(a.length > 0){
				console.log("Empty list received and page holds records - clearing the list");
				for(i=0; i<a.length; i++){
					rodzic.removeChild(a[i]);
				}
			}
			return;
		}
		data.forEach(function (j) {
			o2 = document.getElementById(JSON.stringify(j));
			if(o2 === null || o2 === undefined) {
				o = document.createElement("tr");
				o.classList.add("paymentRecord");
				o.id = JSON.stringify(j);

				s = document.createElement("td");
				//s.classList.add("transactionDateAndTime");
				s.classList.add("transactionGarbledDeformattedEdta");
				//s.name = "transactionDateAndTime";
				s.name = "transactionGarbledDeformattedEdta";
				//s.innerText = j.transactionDateAndTime;
				var m = j.transactionDateAndTime.match(/(2[0-9]{3})-([01][0-9])-([0-3][0-9])T([0-2][0-9]):([0-5][0-9]):([0-5][0-9])([+-][0-9:]+)?/);
				s.innerText = m[3] + '.' + m[2] + '.' + m[1];
				o.appendChild(s);

				s = document.createElement("td");
				s.classList.add("card");
				s.name = "card";
				m = j.maskedCard.match(/^(...).*(....)$/);
				if(m) {
					s.innerText = m[1] + ' ... ' + m[2];
				}else{
					s.innerText = "n/a";
				}
				o.appendChild(s);

				s = document.createElement("td");
				s.classList.add("quantity");
				s.name = "quantity";
				s.innerText = "not implemented\nwould be\n999KW\n[digits][units]";
				o.appendChild(s);

				s = document.createElement("td");
				s.classList.add("unitPrice");
				s.name = "unitPrice";
				s.innerText = "not implemented\nwould be\n[prefix][digits][suffix]";
				o.appendChild(s);

				//s = document.createElement("td");
				//s.classList.add("currencyCode");
				//s.name = "currencyCode";
				//s.innerText = j.currencyCode;
				//o.appendChild(s);

				//s = document.createElement("td");
				//s.classList.add("orderDescription");
				//s.name = "orderDescription";
				//s.innerText = j.orderDescription;
				//o.appendChild(s);

				s = document.createElement("td");
				s.classList.add("total");
				s.name = "total";
				s.innerText = "xx " + (j.amount/100).toFixed(2) + " xx";
				o.appendChild(s);

				s = document.createElement("td");
				s.classList.add("amountPaid");
				s.name = "amountPaid";
				s.innerText = "not implemented\nwould be\n[prefix][digits][suffix]";
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
