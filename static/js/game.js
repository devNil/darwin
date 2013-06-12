window.requestAnimFrame = (function(){
	  return  window.requestAnimationFrame       ||
	          window.webkitRequestAnimationFrame ||
	          window.mozRequestAnimationFrame    ||
	          function( callback ){
			              window.setTimeout(callback, 1000 / 60);
				                };
})();

var game = (function(document, window, undefined){
	//exposed api
	var self = {};
	
	//websocket-object
	var socket;

	//url
	var url = "";

	var entities;

	var ctx;

	var running = false;

	(function init(){
		ctx = document.getElementById("display").getContext("2d");
        ctx.font="30px Arial";
		ctx.fillStyle="#FFF";
		ctx.fillText("Not connected!!!", 20, 20);
	}());

	var render = function(){
		if(entities){
			for(var i = 0; i< entities.length; i++){
				var e = entities[i];
				console.log(e);
                if (e.Gap === false) {
				    ctx.fillStyle = e.Color.toString(16);
                } else {
				    ctx.fillStyle = "#000000";
                }
				ctx.fillRect(e.X, e.Y, 8, 8);

			}
		}
	}

	var wsOnError = function(){
		console.log("error, disconnected");
	}

	var wsOnRegister = function(event){
		var cmd = event.data;
		console.log(cmd);
	}

	var loop = function(){	
		render()
		if(running){
			requestAnimFrame(loop);
		}	
	}

	var message = function(event){
		var cmd = JSON.parse(event.data);

		if(cmd.Id == 0){
			ctx.clearRect(0, 0, 640, 480);
			ctx.fillStyle = "#000";
			ctx.fillRect(0, 0, 640, 480);
		}

		if(cmd.Id == 2){
			ctx.clearRect(310, 195, 30, 30);
			ctx.fillText(cmd.Value, 320, 220);
		}

		if(cmd.Id == 3){
			ctx.clearRect(310, 195, 30, 30)
			running = true;
			loop();
		}

		if(cmd.Id == 4){
			entities = cmd.Value;
		}

		if(cmd.Id == 5){
			running = false
			ctx.clearRect(0, 0, 640, 480);
			ctx.fillStyle = "#FFF";
			ctx.fillText("Spielende :'(", 40, 40);
		}
        if(cmd.Id == 6){
            entity = cmd.Value
            //show color of player
            colorBox = document.getElementById("playerColor");
            colorBox.style.backgroundColor=entity.Color.toString(16); 
            //set start position on map
            ctx.fillStyle = entity.Color.toString(16);
            ctx.fillRect(entity.X, entity.Y, 8, 8);
            ctx.fillSytle = "#FFF"
        }
	}

	var error = function(){
		console.log("error!!!");
	}

	var close = function(){
		ctx.clearRect(0, 0, 640, 480);
		ctx.fillText("Connection closed", 20, 20);
	}

	var prepareKeyListener = function(){
		document.body.addEventListener("keyup", function(e){
			var code = e.keyCode;
			if (code === 39){
				socket.send(JSON.stringify({Id:11, Value:""}));
				e.preventDefault();
			}

			if (code == 37){
				socket.send(JSON.stringify({Id:10, Value:""}));
				e.preventDefault();
			}
		});

		document.body.addEventListener("keydown", function(e){
			var code = e.keyCode;
			if (code === 39 || code == 37){
				e.preventDefault();
			}
		});
	};

	self.connect = function(url){
		document.getElementById("approve").onclick = approve;
		url = url;
		socket = new WebSocket(url);
		socket.onmessage = message;
		socket.onerror = error;

		socket.onclose = close;

		prepareKeyListener();

	}
	var approve = function(){
		socket.send(JSON.stringify({Id:1, Value:""}));
	}

	//expose api
	return self;

}(document, window));
