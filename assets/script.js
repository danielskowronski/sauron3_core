var refreshRate = 15000; //ms
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
	if (window.location.href.indexOf("2col")>0){a=$("#display").html(); $("body").append('<div id="display2">'+a+'</div>');$("#display2").children(".hostRow").remove(); len=parseInt($("#display").children(".hostRow").length); for (i=parseInt(len/2)+1; i<len; i++) { $($($("#display").children(".hostRow")[parseInt(len/2)+1]).detach()).appendTo("#display2"); } $("#display,#display2").css("display","inline-block").css("width","500px")}
}

var timeoutHandle = null

function probeLivecheck(){
	clearTimeout(timeoutHandle)
	timeoutHandle = 
		setTimeout(function(){ $("#display").css("background", "orange"); }, refreshRate*5);

	$.get( "/probe/", function( data ) {
		livechecks=jQuery.parseJSON(data);
		
		clearTimeout(timeoutHandle)
			$("#display").css("background", "none");

		$.each( livechecks, function( key, value ) {
			hosts = jQuery.parseJSON(data)

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
