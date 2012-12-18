(function(IDE) {

IDE.Events = function(workspace) {
	this.workspace = workspace;

	window.onresize = _.bind(this.onresize, this);
	window.onkeypress = _.bind(this.onkeypress, this);
	window.onkeydown = _.bind(this.onkeydown, this);

	var ta = this.workspace.textarea;
	ta.addEventListener("blur", function(evt) {
		ta.focus();
	});
	ta.focus();

	this.onresize();
};
IDE.Events.prototype = {
	"onkeypress": function(evt) {
		if (evt.charCode) {
			this.workspace.focused.oninsert(String.fromCharCode(evt.charCode));
		}
	},
	"onkeydown": function(evt) {
		switch (evt.keyCode) {
		case 8:
			this.workspace.focused.onbackspace();
			break;
		case 13:
			this.workspace.focused.oninsert('\n');
			break;
		case 37:
			this.workspace.focused.onleft();
			break;
		case 39:
			this.workspace.focused.onright();
			break;
		case 46:
			this.workspace.focused.ondelete();
			break;
		default:
			console.log(evt.keyCode);
			return;
		}

		evt.preventDefault();
	},

	"onresize": function(evt) {
		var w = document.body.clientWidth;
		var h = document.body.clientHeight;
		this.workspace.width = w;
		this.workspace.height = h;
		this.workspace.canvas.setAttribute('width', w);
		this.workspace.canvas.setAttribute('height', h);
		this.workspace.canvas.style.width = w + 'px';
		this.workspace.canvas.style.height = h + 'px';
		this.workspace.render();
	}
};

})(IDE);