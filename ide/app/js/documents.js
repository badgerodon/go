(function(IDE) {

IDE.Document = function() {
	this.column = null;
	this.title = "";
	this.body = "";
	this.offset = {
		x: 0,
		y: 0
	};
	this.cursor = [5, 0];
};
IDE.Document.prototype = {
	onbackspace: function() {
		if (this.cursor[1] > 0) {

		} else {
			var s = this.cursor[0] - 1;
			if (s < 0) s = 0;
			var e = s + 1;
		}

		this.body = this.body.substr(0, s)
			+ this.body.substr(e);

		this.cursor = [s, 0];

		this.column.workspace.render();
	},
	ondelete: function() {
		var s = this.cursor[0];
		var e = s + (this.cursor[1] || 1);

		this.body = this.body.substr(0, s)
			+ this.body.substr(e);

		this.cursor = [s, 0];

		this.column.workspace.render();
	},
	onleft: function() {
		var s = this.cursor[0];
		if (s > 0) {
			s--;
		}
		this.cursor = [s, 0];

		this.column.workspace.render();
	},
	onright: function() {
		var s = this.cursor[0];
		if (s < this.body.length) {
			s++;
		}
		this.cursor = [s, 0];

		this.column.workspace.render();
	},
	oninsert: function(c) {
		var s = this.cursor[0];
		var e = s + this.cursor[1];

		this.body = this.body.substr(0, s)
			+ c
			+ this.body.substr(e);

		this.cursor = [e+1, 0];

		this.column.workspace.render();
	},
	render: function(ctx, x, y, w, h) {
		ctx.save();

		ctx.fillStyle = "#FFF";
		ctx.font = "normal 12px Consolas";
		ctx.textBaseline = "middle";

		var dim = ctx.measureText("A");
		var oh = 16;
		var ow = dim.width;
		var ox = x;
		var oy = y;

		var ln = 0;
		var nextLineNumber = function() {
			ctx.textAlign = "right";
			ctx.fillText((ln + 1) + "", x + ow*3, oy + oh/2);
			ox = x + ow*3 + 5;
			ln++;
		};
		nextLineNumber();

		for (var i=0; i<this.body.length; i++) {
			var c = this.body.charAt(i);

			if (i == this.cursor[0]) {
				ctx.fillRect(ox, oy, 1, oh);
			}

			if (c == '\n') {
				oy += oh;
				nextLineNumber();
				continue;
			}

			if (c == '\t') {
				ox += ow*2;
				continue;
			}

			ctx.textAlign = "center";
			ctx.fillText(c, ox + ow/2, oy + oh/2);
			ox += ow;
		}

		ctx.restore();
	}
};

})(IDE);