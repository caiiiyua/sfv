(function() {
	var app = angular.module('myapp', []);
	app.controller('MainController', ['$http', function($http) {
		this.message = "hello";
		this.users = users;
		this.user_count = users.length;
		this.email = "";

		this.addUser = function() {
			if (this.email.length > 0) {
			email = this.email;
			var res = $http.post('/sfvers', {'email': email});
			res.success(function(data, status, headers, config) {
				this.user_count = data;
			});
			users.push(this.email);
			this.email = "";
			this.users = users;
			// this.user_count = users.length;
			};

		};
	}]);

	var users = [
		'ckjacket@mm.com',
		'kkk@ls.com'
	];
})();