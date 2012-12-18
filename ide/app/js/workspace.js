(function(IDE) {

IDE.Workspace = function() {
	this.canvas = document.getElementById("ide");
	this.textarea = document.getElementById("textarea");
	this.context = this.canvas.getContext("2d");
  this.width = 0;
  this.height = 0;
  this.columns = [];

  this.events = new IDE.Events(this);

  this.add(new IDE.Column(0.5));
  this.add(new IDE.Column(0.5));

	var doc = new IDE.Document();
	doc.body = 'var x = function() {\n\tconsole.log("HELLO WORLD");\n};\n';
	this.columns[0].add(doc);
	this.focused = doc;

	this.render();
};
IDE.Workspace.prototype = {
  "add": function(column) {
    if (column.workspace) {
      column.workspace.remove(column);
    }
    column.workspace = this;
    this.columns.push(column);

    this.render();
  },
  "render": function() {
    var ctx = this.context;

    ctx.save();

    ctx.fillStyle = "#000";
    ctx.fillRect(0, 0, this.width, this.height);

    // Render top bar
    ctx.fillStyle = "#111";
    ctx.fillRect(0, 0, this.width, 19);
    ctx.fillStyle = "#666";
    ctx.fillRect(0, 19, this.width, 1);

    ctx.restore();

    var ox = 0;
    for (var i=0; i<this.columns.length; i++) {
      var col = this.columns[i];
      var cw = Math.floor(this.width * col.width);
      this.columns[i].render(ctx, ox, 20, cw-1, this.height - 20);
      if (i > 0) {
        ctx.save();
        ctx.fillStyle = "#666";
        ctx.fillRect(ox, 20, 1, this.height - 20);
        ctx.restore();
      }
      ox += cw;
    }
  }
};

IDE.Column = function(width) {
  this.documents = [];
  this.width = width;
};
IDE.Column.prototype = {
  "add": function(doc) {
    if (doc.column) {
      doc.column.remove(doc);
    }
    doc.column = this;
    this.documents.push(doc);
  },
  "remove": function(doc) {
    for (var i=0; i<this.documents.length; i++) {
      if (this.documents[i] == doc) {
        this.documents.splice(i, 1);
        doc.column = this;
				this.render();
        return
      }
    }
  },
  "render": function(ctx, x, y, w, h) {
    ctx.save();

    ctx.fillStyle = "#111";
    ctx.fillRect(x, y, w, 19);
    ctx.fillStyle = "#666";
    ctx.fillRect(x, y+19, w, 1);

    ctx.restore();

		for (var i=0; i<this.documents.length; i++) {
			this.documents[i].render(ctx, x, y+20, w, h-20);
		}
  }
};

})(IDE);