var refreshRate = 5000; //ms
function constructHostRow(name){
	return 	"<div class='hostRow' data-host='"+name+"' >"+
			"<div class='hostTitle'>"+name+"</div>"+
			"<div class='hostProbes'></div>"+
			"</div>";
}

function constructLivecheckField(host,probe){
	return 	"<div class='hostCheck' data-host='"+host+"' data-probe='"+probe+"'>"+probe+"</div>";
}

function parseDefinitions(data){
	hosts = jQuery.parseJSON(data)
	$.each( hosts , function( key, host ) {
		console.log(host)
		$("#display").append(constructHostRow(host.Title))
		$.each( host.Probes, function( key, probe ) {
			$("[data-host='"+host.Title+"'].hostRow .hostProbes").append(
				constructLivecheckField(host.Title, probe.Title)
			);

		});
	});
}

function probeLivecheck(){
	var timeoutHandle = 
		setTimeout(function(){ $("#display").css("background", "orange"); }, refreshRate*5);

	$.get( "/probe/", function( data ) {
		livechecks=jQuery.parseJSON(data);
		$.each( livechecks, function( key, value ) {
			hosts = jQuery.parseJSON(data)

			clearTimeout(timeoutHandle)
			$("#display").css("background", "none");

			$.each( hosts , function( key, host ) {
				$.each( host.Probes, function( key, probe ) {			
					var obj = $("[data-host='"+host.Title+"'] [data-probe='"+probe.Title+"']")
					if (probe.Alive) {
						$(obj).removeClass("dead")
						$(obj).addClass("alive")
					} else {
						$(obj).removeClass("alive")
						$(obj).addClass("dead")
					}
				});
			});
		});
	});
}

$(function() {
	$.get( "/definitions/", function( data ) {
		parseDefinitions(data)

		probeLivecheck();
		setInterval(function(){ probeLivecheck(); }, refreshRate);
	});
});