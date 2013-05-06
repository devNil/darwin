var Entity = function(x,y,size,dir,color){
    this.x = x;
    this.y = y;
    this.size = size;
    this.dir = dir;
    this.color = color;
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
    this.player = new Entity(0,0,16,0,0xFF00FF); 
    this.input = new InputHandler();
};
Game.prototype.startUp = function(ctx, x, y){
    //add walls
    for(var i = 0; i < x;i++){
        ctx.fillStyle = "rgb(200,0,0)";
        ctx.fillRect(i,0,16,16);
        ctx.fillRect(i,y-16,16,16);
    }
    for(var i = 0; i < y;i++){
        ctx.fillStyle = "rgb(200,0,0)";
        ctx.fillRect(0,i,16,16);
        ctx.fillRect(x-16,i,16,16);
    }

}
Game.prototype.render = function(ctx){
    for(var i = 0; i < this.entities.length; i++){
        var ent = this.entities[i];
        ctx.fillStyle = ent.color;
        ctx.fillRect(ent.x,ent.y,ent.size,ent.size);
    }
}
