(function() {
	var app = angular.module('myapp', []);
	app.controller('MainController', function() {
		this.message = "hello";
		this.user_num = 1;
	});
})();