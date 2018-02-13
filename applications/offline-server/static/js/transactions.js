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
		rodzic = document.getElementById("tableTBody");
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
				var trType = "div";
				var tdType = "span";

				o = document.createElement(trType);
				o.classList.add("tr");
				o.classList.add("paymentRecord");
				o.id = JSON.stringify(j);

				s = document.createElement(tdType);
				//s.classList.add("transactionDateAndTime");
				s.classList.add("transactionGarbledDeformattedEdta");
				s.classList.add("td");
				//s.name = "transactionDateAndTime";
				s.name = "transactionGarbledDeformattedEdta";
				//s.innerText = j.transactionDateAndTime;
				var m = j.transactionDateAndTime.match(/(2[0-9]{3})-([01][0-9])-([0-3][0-9])T([0-2][0-9]):([0-5][0-9]):([0-5][0-9])([+-][0-9:]+)?/);
				s.innerText = m[4] + ':' + m[5] + ':' + m[6];
				o.appendChild(s);

				s = document.createElement(tdType);
				s.classList.add("card");
				s.classList.add("td");
				s.name = "card";
				m = j.maskedCard.match(/^(...).*(.[0-9][0-9][0-9])$/);
				if(m) {
					s.innerHTML = 'xxx&nbsp;...&nbsp;' + m[2];
				}else{
					s.innerHTML = "n/a";
				}
				o.appendChild(s);

				q = document.createElement(tdType);
				q.classList.add("quantity");
				q.classList.add("td");
				q.name = "quantity";
				//q.innerText = "n/a";


				//s = document.createElement("td");
				//s.classList.add("unitPrice");
				//s.name = "unitPrice";
				//s.innerText = "not implemented\nwould be\n[prefix][digits][suffix]";
				//o.appendChild(s);

				//s = document.createElement("td");
				//s.classList.add("currencyCode");
				//s.name = "currencyCode";
				//s.innerText = j.currencyCode;
				//o.appendChild(s);

				s = document.createElement(tdType);
				s.classList.add("orderDescription");
				s.classList.add("td");
				s.name = "orderDescription";
				m = j.orderDescription.match(/(.*) - (.*) @ (.*)./);
				s.innerText = m[3] + ' - ' + m[1];
				if (m[2] == "1 units") {
					q.innerText = "1 unit";
				}
				else {
					q.innerText = m[2];
				}
				o.appendChild(q);
				o.appendChild(s);

				s = document.createElement(tdType);
				s.classList.add("total");
				s.classList.add("td");
				s.name = "total";
				switch(j.currencyCode) {
					case "GBP":
						s.innerHTML = "&pound;" + (j.amount/100).toFixed(2);
						break;
					case "USD":
						s.innerText = "$" + (j.amount/100).toFixed(2);
						break;
					default:
						s.innerText = (j.amount/100).toFixed(2) + j.currencyCode;
				}
				o.appendChild(s);

				//s = document.createElement("td");
				//s.classList.add("amountPaid");
				//s.name = "amountPaid";
				//s.innerText = "not implemented\nwould be\n[prefix][digits][suffix]";
				//o.appendChild(s);


				//o.innerText=JSON.stringify(j);
				o.style.backgroundColor="green";
				rodzic.appendChild(o);
				rodzic.scrollTo({left:0,top:rodzic.scrollHeight,behavior:"smooth"});
				o.style.backgroundColor="white";
			}
		});
	});
}

$.ajaxSetup({'cache':true, timeout:300});
refreshChargerStatus();
var smartChargerTimer = setInterval(refreshChargerStatus, 750 * 2);
