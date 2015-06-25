(function() {
	var app = angular.module('myapp', []);
	app.controller('MainController', function() {
		this.message = "hello";
		this.users = users;
		this.user_count = users.length;
		this.email = "";

		this.addUser = function() {
			if (this.email.length > 0) {				
			users.push(this.email);
			this.email = "";
			this.users = users;
			this.user_count = users.length;
			};
		};
	});

	var users = [
		'ckjacket@mm.com',
		'kkk@ls.com'
	];
})();