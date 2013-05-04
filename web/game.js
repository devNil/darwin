window.requestAnimFrame = (function(){
  return  window.requestAnimationFrame       ||
          window.webkitRequestAnimationFrame ||
          window.mozRequestAnimationFrame    ||
          function( callback ){
            window.setTimeout(callback, 1000 / 60);
          };
})();

var Entity = function(){
    this.x = 0;
    this.y = 0;
    this.s = 16;
    this.dir = 0;
    this.color = "#FF00FF";
}

var InputHandler = function(){
    this.l = false;
    this.r = false;
}

InputHandler.prototype.left = function(){
    this.l = true;
    this.r = false;
}

InputHandler.prototype.right = function(){
    this.r = true;
    this.l = false;
}

InputHandler.prototype.reset = function(){
    this.l = false;
    this.r = false;
}

var Game = function(){
    this.ticks = 0;
    this.entities = new Array();
    this.player = new Entity(); 
    //this.entities.push();
    this.input = new InputHandler();
};

Game.prototype.render = function(ctx){
    /*for(var i = 0; i < this.entities.length; i++){
        var ent = this.entities[i];
        ctx.fillStyle = ent.color;
        ctx.fillRect(ent.x,ent.y,ent.s,ent.s);
    }*/
    
    var p = this.player;
    
    console.log(p.x);
    
    ctx.fillStyle = "black";
    ctx.fillRect(p.x,p.y,p.s,p.s);
}

Game.prototype.tick = function(){
    this.ticks++;
    
    if(this.ticks % 20 == 0){
        for(var i = 0; i < this.entities.length; i++){
            this.entities[i].x += 16;
        }
    
    if(this.input.r){
            if(this.player.dir === 3){
                this.player.dir = 0;
            }else{
                this.player.dir += 1;
            }
            
    }
    
    if(this.input.l){
        if(this.player.dir === 0){
            this.player.dir = 3;
        }else{
            this.player.dir -= 1; 
        }
    }
    
    console.log(this.player.dir);
        
    //if(this.player.dir === 4) this.player.dir = 0;
    //if(this.player.dir === -1) this.player.dir = 3;
    
    if(this.player.dir === 0){
        this.player.x += 16;
    }
    
    if(this.player.dir === 1){
        this.player.y += 16;
    }
    
    if(this.player.dir === 2){
        this.player.x -= 16;
    }
    
    if(this.player.dir === 3){
        this.player.y -= 16;
    }
    
    this.input.reset();
    }
}