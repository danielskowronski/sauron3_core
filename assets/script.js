function constructHostRow(id, name){
	return 	"<div class='hostRow' data-hostid='"+id+"'>"+
			"<div class='hostTitle'>"+name+"</div>"+
			"<div class='hostProbes'></div>"+
			"</div>";
}
function constructHostsContainer(data){
	hosts=jQuery.parseJSON(data)
	$.each( hosts, function( key, value ) {
		$("#display").append(constructHostRow(value.id, value.name))
	});
}

function constructLivecheckField(id, name){
	return 	"<div class='hostCheck dead' data-checkid='"+id+"'>"+name+"</div>";
}
function constructLivechecksContainer(data){
	hosts=jQuery.parseJSON(data)
	$.each( hosts, function( key, value ) {
		$("[data-hostid="+value.host_id+"]").append(constructLivecheckField(value.check_id, value.name))
	});
}

function probeLivecheck(){
	$.get( "/probe/", function( data ) {
		livechecks=jQuery.parseJSON(data)
		$.each( livechecks, function( key, value ) {
			obj = $("[data-checkid="+value.check_id+"]")
			if (value.alive) {
				$(obj).removeClass("dead")
				$(obj).addClass("alive")
			} else {
				$(obj).removeClass("alive")
				$(obj).addClass("dead")
			}
		});
	});
}

$(function() {
	$.get( "/hosts/", function( data ) {
		constructHostsContainer(data)

		$.get( "/livechecks/", function( data ) {
			constructLivechecksContainer(data)

			probeLivecheck();
			setInterval(function(){ probeLivecheck(); }, 5000);
		});
	});

});